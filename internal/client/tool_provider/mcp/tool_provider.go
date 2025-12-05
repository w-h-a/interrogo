package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/w-h-a/interrogo/api/tools/v1alpha1"
	toolprovider "github.com/w-h-a/interrogo/internal/client/tool_provider"
)

type mcpToolProvider struct {
	options toolprovider.Options
	client  *client.Client
}

func (tp *mcpToolProvider) Start(ctx context.Context) error {
	if err := tp.client.Start(ctx); err != nil {
		return err
	}

	if _, err := tp.client.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo:      mcp.Implementation{Name: "vertex-agent", Version: "0.1.0"},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (tp *mcpToolProvider) List(ctx context.Context) ([]v1alpha1.ToolDefinition, error) {
	rsp, err := tp.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}

	var tools []v1alpha1.ToolDefinition
	for _, t := range rsp.Tools {
		tools = append(tools, v1alpha1.ToolDefinition{
			Name:        t.Name,
			Description: t.Description,
			Schema: &v1alpha1.Schema{
				Type:       t.InputSchema.Type,
				Properties: t.InputSchema.Properties,
				Required:   t.InputSchema.Required,
			},
		})
	}

	return tools, nil
}

func (tp *mcpToolProvider) Call(ctx context.Context, name string, args map[string]any) (string, error) {
	result, err := tp.client.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		},
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result.Content), nil
}

func NewToolProvider(opts ...toolprovider.Option) toolprovider.ToolProvider {
	options := toolprovider.NewOptions(opts...)

	tp := &mcpToolProvider{
		options: options,
	}

	c, _ := client.NewSSEMCPClient(options.Location)

	tp.client = c

	return tp
}
