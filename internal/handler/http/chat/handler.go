package chat

import (
	"encoding/json"
	"log"
	"net/http"

	httphandler "github.com/w-h-a/interrogo/internal/handler/http"
	"github.com/w-h-a/interrogo/internal/service/agent"
)

type chatHandler struct {
	agent *agent.Agent
}

func (h *chatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	ctx := httphandler.ReqToCtx(r)

	var req struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	log.Printf("Received Prompt: %s", req.Message)

	rsp, toolCalls, err := h.agent.TakeTurns(ctx, req.Message)
	if err != nil {
		log.Printf("Agent error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"response":   rsp,
		"tool_calls": toolCalls,
	})
}

func New(a *agent.Agent) *chatHandler {
	return &chatHandler{agent: a}
}
