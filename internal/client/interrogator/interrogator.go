package interrogator

import "context"

type Interrogator interface {
	Interrogate(ctx context.Context, prompt string) (response string, toolCalls []string, err error)
}
