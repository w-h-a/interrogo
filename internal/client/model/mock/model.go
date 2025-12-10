package mock

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/w-h-a/interrogo/internal/client/model"
)

type mockModel struct {
	options  model.Options
	callFunc func(prompt string) (string, error)
}

func (m *mockModel) Call(ctx context.Context, p string, o ...llms.CallOption) (string, error) {
	return m.callFunc(p)
}

func (m *mockModel) GenerateContent(ctx context.Context, msgs []llms.MessageContent, o ...llms.CallOption) (*llms.ContentResponse, error) {
	return nil, nil
}

func NewModel(opts ...model.Option) *mockModel {
	options := model.NewOptions(opts...)

	m := &mockModel{
		options: options,
	}

	if cf, ok := getCallFuncFromCtx(options.Context); ok {
		m.callFunc = cf
	}

	return m
}
