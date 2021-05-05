package app

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	group "github.com/oklog/oklog/pkg/group"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (a *App) createService() (g *group.Group) {
	g = &group.Group{}

	// init GRPC Handlers
	Base.initGRPCHandler(g)

	// init http endpoints
	Base.initHTTPEndpoint(g)

	// init cancel
	Base.initCancelInterrupt(g)
	return g
}

// defaultGRPCOptions
// add options for grpc connection
func (a *App) defaultGRPCOptions(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption {
	options := []grpc.ServerOption{}

	// UnaryInterceptor and OpenTracingServerInterceptor for tracer
	options = append(options, grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
	))

	return options
}
