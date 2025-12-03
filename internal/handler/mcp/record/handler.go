package record

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type recordHandler struct{}

func (h *recordHandler) ListRecords(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("arguments not a map"), nil
	}

	resource, ok := args["resource"].(string)
	if !ok {
		return mcp.NewToolResultError("resource argument must be a string"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found 1 record: %s", resource)), nil
}

func (h *recordHandler) ListRecordsTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.Tool{
			Name:        "list_records",
			Description: "List records",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"resource": map[string]string{
						"type":        "string",
						"description": "The resource name",
					},
				},
				Required: []string{"resource"},
			},
		},
		Handler: h.ListRecords,
	}
}

func New() *recordHandler {
	return &recordHandler{}
}
