package middleware

import "context"

var (
	M Middleware = &middle{}
)

type Middleware interface {
	JWT(ctx context.Context) (context.Context, error)
}

type middle struct{}

func (m *middle) JWT(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
