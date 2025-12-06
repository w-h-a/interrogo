package v1alpha1

type ChatResponse struct {
	Response  string   `json:"response"`
	ToolCalls []string `json:"tool_calls"`
	Error     string   `json:"error,omitempty"`
}
