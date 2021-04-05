package provider

import (
	"blogfa/auth/pkg/jtrace"
	"context"
)

func (p *Provider) Register(ctx context.Context, user Provider) error {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "register_provider")
	defer span.Finish()
	span.SetTag("register", "register provider model")

	return nil
}
