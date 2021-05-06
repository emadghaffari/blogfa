package app

import (
	"blogfa/auth/service/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
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
	options = append(options, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		grpc_auth.UnaryServerInterceptor(middleware.M.JWT),
		grpc_prometheus.UnaryServerInterceptor,
	),
	))

	options = append(options, grpc.StreamInterceptor(
		grpc_auth.StreamServerInterceptor(middleware.M.JWT),
	))

	return options
}
