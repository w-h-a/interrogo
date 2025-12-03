package mcp

import (
	"context"

	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/w-h-a/interrogo/internal/server"
)

type toolMiddlewareKey struct{}

func WithToolMiddleware(ms ...mcpserver.ToolHandlerMiddleware) server.Option {
	return func(o *server.Options) {
		o.Context = context.WithValue(o.Context, toolMiddlewareKey{}, ms)
	}
}

func getToolMiddlewareFromCtx(ctx context.Context) ([]mcpserver.ToolHandlerMiddleware, bool) {
	ms, ok := ctx.Value(toolMiddlewareKey{}).([]mcpserver.ToolHandlerMiddleware)
	return ms, ok
}

type resourceMiddlewareKey struct{}

func WithResourceMiddleware(ms ...mcpserver.ResourceHandlerMiddleware) server.Option {
	return func(o *server.Options) {
		o.Context = context.WithValue(o.Context, resourceMiddlewareKey{}, ms)
	}
}

func getResourceMiddlewareFromCtx(ctx context.Context) ([]mcpserver.ResourceHandlerMiddleware, bool) {
	ms, ok := ctx.Value(resourceMiddlewareKey{}).([]mcpserver.ResourceHandlerMiddleware)
	return ms, ok
}
