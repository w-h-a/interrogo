package mock

import (
	"context"

	"github.com/w-h-a/interrogo/internal/client/model"
)

type callFuncKey struct{}

func WithCallFunc(cf func(prompt string) (string, error)) model.Option {
	return func(o *model.Options) {
		o.Context = context.WithValue(o.Context, callFuncKey{}, cf)
	}
}

func getCallFuncFromCtx(ctx context.Context) (func(prompt string) (string, error), bool) {
	cf, ok := ctx.Value(callFuncKey{}).(func(prompt string) (string, error))
	return cf, ok
}
