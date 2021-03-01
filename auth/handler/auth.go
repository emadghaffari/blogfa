package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	auth "github.com/emadghaffari/blogfa/auth/proto"
)

type Auth struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Auth) Register(ctx context.Context, in *auth.RegisterRequest, out *auth.RegisterResponse) error {
	log.Info("Received Auth.Call request")
	return nil
}
