package controller

import context "context"

type (
	options struct {
		ctx context.Context
	}
	OptionFunc func(*options)
)

func WithContext(ctx context.Context) OptionFunc {
	return OptionFunc(func(o *options) {
		o.ctx = ctx
	})
}
