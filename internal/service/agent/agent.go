package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tmc/langchaingo/llms"
	toolprovider "github.com/w-h-a/interrogo/internal/client/tool_provider"
)

type Agent struct {
	model        llms.Model
	tools        toolprovider.ToolProvider
	instructions string
	isRunning    bool
	mtx          sync.RWMutex
}

func (a *Agent) Run(stop chan struct{}) error {
	a.mtx.RLock()
	if a.isRunning {
		a.mtx.RUnlock()
		return errors.New("agent is already running")
	}
	a.mtx.RUnlock()

	if err := a.Start(); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}

	<-stop

	return a.Stop()
}

func (a *Agent) Start() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.isRunning {
		return errors.New("agent is already started")
	}

	a.isRunning = true

	return nil
}

func (a *Agent) Stop() error {
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	return a.stop(stopCtx)
}

func (a *Agent) stop(ctx context.Context) error {
	a.mtx.Lock()

	if !a.isRunning {
		a.mtx.Unlock()
		return errors.New("agent not running")
	}

	a.isRunning = false

	a.mtx.Unlock()

	gracefulStopDone := make(chan struct{})
	go func() {
		// TODO: close clients gracefully
		close(gracefulStopDone)
	}()

	var stopErr error

	select {
	case <-gracefulStopDone:
	case <-ctx.Done():
		stopErr = ctx.Err()
	}

	return stopErr
}

func (a *Agent) TakeTurns(ctx context.Context, input string) (string, []string, error) {
	availableTools, err := a.tools.List(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to list tools: %w", err)
	}

	var lcTools []llms.Tool
	for _, t := range availableTools {
		params := map[string]any{
			"type":       t.Schema.Type,
			"properties": t.Schema.Properties,
		}
		if len(t.Schema.Required) > 0 {
			params["required"] = t.Schema.Required
		}

		lcTools = append(lcTools, llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		})
	}

	history := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, a.instructions),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	var toolCallsLog []string
	maxTurns := 5

	fmt.Println("\n--- Starting Turns ---")

	for range maxTurns {
		// Ask LLM what to do
		rsp, err := a.model.GenerateContent(ctx, history, llms.WithTools(lcTools))
		if err != nil {
			return "", nil, fmt.Errorf("llm generation failed: %w", err)
		}

		msg := rsp.Choices[0]

		// If no tool calls, return
		if len(msg.ToolCalls) == 0 {
			return msg.Content, toolCallsLog, nil
		}

		// Append the LLM's response to history
		var parts []llms.ContentPart
		for _, tc := range msg.ToolCalls {
			parts = append(parts, llms.ToolCall{ID: tc.ID, Type: tc.Type, FunctionCall: tc.FunctionCall})
		}
		history = append(history, llms.MessageContent{
			Role:  llms.ChatMessageTypeAI,
			Parts: parts,
		})

		// Execute the tool call
		for _, tc := range msg.ToolCalls {
			fnName := tc.FunctionCall.Name
			fnArgsStr := tc.FunctionCall.Arguments

			toolCallsLog = append(toolCallsLog, fnName)

			fmt.Printf("Agent calling tool: %s\n", fnName)

			var args map[string]any
			var toolResult string

			if err := json.Unmarshal([]byte(fnArgsStr), &args); err != nil {
				toolResult = fmt.Sprintf("Error: failed to parse arguments JSON: %v", err)
			} else {
				toolResult, err = a.tools.Call(ctx, fnName, args)
				if err != nil {
					toolResult = fmt.Sprintf("Error: %v", err)
				}
			}

			// Feed the result back to the LLM
			history = append(history, llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: tc.ID,
						Name:       fnName,
						Content:    toolResult,
					},
				},
			})
		}
	}

	return "Error: max turns exceeded", toolCallsLog, nil
}

func New(model llms.Model, tools toolprovider.ToolProvider, instructions string) *Agent {
	return &Agent{
		model:        model,
		tools:        tools,
		instructions: instructions,
		isRunning:    false,
		mtx:          sync.RWMutex{},
	}
}
