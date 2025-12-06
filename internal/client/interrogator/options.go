package interrogator

import (
	"context"
)

type Option func(*Options)

type Options struct {
	Target  string
	Context context.Context
}

func WithTarget(target string) Option {
	return func(o *Options) {
		o.Target = target
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
