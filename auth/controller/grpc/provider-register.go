package grpc

import (
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
)

// RegisterProvider, for create new provider
func (a *Auth) RegisterProvider(ctx context.Context, req *pb.ProviderRegisterRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("register-provider")
	defer span.Finish()
	span.SetTag("register", "register provider")

	response, err := grpc.Service.RegisterProvider(jtrace.Tracer.ContextWithSpan(ctx, span), req)
	if err != nil {
		return response, err
	}

	child := jtrace.Tracer.ChildOf(span, "register")
	child.SetTag("register", "after register provider")
	defer child.Finish()

	// return successfully message
	return response, nil
}
