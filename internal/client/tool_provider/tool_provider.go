package toolprovider

import (
	"context"

	"github.com/w-h-a/interrogo/api/tools/v1alpha1"
)

type ToolProvider interface {
	Start(ctx context.Context) error
	List(ctx context.Context) ([]v1alpha1.ToolDefinition, error)
	Call(ctx context.Context, name string, args map[string]any) (string, error)
}
