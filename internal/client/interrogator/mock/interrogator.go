package mock

import (
	"context"
	"fmt"

	"github.com/w-h-a/interrogo/internal/client/interrogator"
)

type mockInterrogator struct {
	options            interrogator.Options
	responses          []string
	toolLeakTurn       int
	networkFailureTurn int
	count              int
}

func (i *mockInterrogator) Interrogate(ctx context.Context, prompt string) (response string, toolCalls []string, err error) {
	defer func() { i.count++ }()

	rsp := "Refusal."

	if i.count < len(i.responses) {
		rsp = i.responses[i.count]
	}

	var tools []string

	if i.count == i.toolLeakTurn {
		tools = []string{"dangerous_tool"}
	}

	if i.count == i.networkFailureTurn {
		return "", nil, fmt.Errorf("network timeout")
	}

	return rsp, tools, nil
}

func NewInterrogator(opts ...interrogator.Option) *mockInterrogator {
	options := interrogator.NewOptions(opts...)

	i := &mockInterrogator{
		options:            options,
		responses:          []string{},
		toolLeakTurn:       -1,
		networkFailureTurn: -1,
	}

	if rsps, ok := getResponsesFromCtx(options.Context); ok {
		i.responses = rsps
	}

	if index, ok := getToolLeakTurnFromCtx(options.Context); ok {
		i.toolLeakTurn = index
	}

	if index, ok := getNetworkFailureTurnFromCtx(options.Context); ok {
		i.networkFailureTurn = index
	}

	return i
}
