package mock

import (
	"context"

	"github.com/w-h-a/interrogo/internal/client/interrogator"
)

type responsesKey struct{}

func WithResponses(rsps []string) interrogator.Option {
	return func(o *interrogator.Options) {
		o.Context = context.WithValue(o.Context, responsesKey{}, rsps)
	}
}

func getResponsesFromCtx(ctx context.Context) ([]string, bool) {
	rsps, ok := ctx.Value(responsesKey{}).([]string)
	return rsps, ok
}

type toolLeakTurnKey struct{}

func WithToolLeakTurn(index int) interrogator.Option {
	return func(o *interrogator.Options) {
		o.Context = context.WithValue(o.Context, toolLeakTurnKey{}, index)
	}
}

func getToolLeakTurnFromCtx(ctx context.Context) (int, bool) {
	index, ok := ctx.Value(toolLeakTurnKey{}).(int)
	return index, ok
}

type networkFailureTurnKey struct{}

func WithNetworkFailureTurn(index int) interrogator.Option {
	return func(o *interrogator.Options) {
		o.Context = context.WithValue(o.Context, networkFailureTurnKey{}, index)
	}
}

func getNetworkFailureTurnFromCtx(ctx context.Context) (int, bool) {
	index, ok := ctx.Value(networkFailureTurnKey{}).(int)
	return index, ok
}
