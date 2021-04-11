package provider

import (
	"blogfa/auth/database/mysql"
	"blogfa/auth/pkg/jtrace"
	"context"
)

// Register method, register a provider
func (p *Provider) Register(ctx context.Context, prov Provider) error {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "register_provider")
	defer span.Finish()
	span.SetTag("register", "register provider model")

	tx := mysql.Storage.GetDatabase().Begin()

	if err := tx.Create(&prov).Error; err != nil {
		tx.Rollback()
		return err
	}
	defer tx.Commit()

	return nil
}
