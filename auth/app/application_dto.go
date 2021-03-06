package app

import (
	router "blogfa/auth/app/router/http"
	"blogfa/auth/app/router/middleware"
	"blogfa/auth/client/broker"
	"blogfa/auth/client/consul"
	"blogfa/auth/client/etcd"
	"blogfa/auth/client/mysql"
	"blogfa/auth/client/redis"
	"blogfa/auth/config"
	controller "blogfa/auth/controller/grpc"
	"blogfa/auth/pkg/jtrace"
	zapLogger "blogfa/auth/pkg/logger"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	group "github.com/oklog/oklog/pkg/group"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartApplication func
func (a *App) StartApplication() {
	fmt.Println("\n\n--------------------------------")

	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init zap logger
	a.initLogger()

	// init configs
	if err := a.initConfigs(); err != nil {
		return
	}

	//  Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency
	closer, err := a.initJaeger()
	if err != nil {
		return
	}
	defer closer.Close()

	if viper.GetString("environment") == "production" {
		if err := a.initConfigServer(); err != nil {
			fmt.Println(err.Error())
		}
		defer etcd.Storage.GetClient().Close()
	}

	if err := a.initDatabase(); err != nil {
		return
	}

	if err := a.initRedis(); err != nil {
		return
	}

	if err := a.initConsul(); err != nil {
		return
	}

	if err := a.initMessageBroker(); err != nil {
		return
	}
	defer broker.Nats.Conn().Close()

	// create service
	g := a.createService()
	fmt.Printf("--------------------------------\n\n")
	if err := g.Run(); err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit("server stopped")
	}
}

// init zap logger
func (a *App) initLogger() {
	defer fmt.Printf("zap logger is available \n")
	zapLogger.SetLogPath("logs")
	logger = zapLogger.GetZapLogger(config.Global.Debug())
}

// init configs
func (a *App) initConfigs() error {
	defer fmt.Printf("configs loaded from file successfully \n")

	// Current working directory
	dir, err := os.Getwd()
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit("init configs")
	}

	// read from file
	return config.Load(dir + "/config.yaml")
}

// init grpc connection
func (a *App) initGRPCHandler(g *group.Group) {
	defer fmt.Printf("grpc connected port:%s \n", config.Global.Service.GRPC.Port)

	options := a.defaultGRPCOptions(logger, jtrace.Tracer.GetTracer())
	// Add your GRPC options here

	lis, err := net.Listen("tcp", config.Global.Service.GRPC.Port)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit(err.Error())
	}

	g.Add(func() error {
		baseServer := grpc.NewServer(options...)

		// reflection for evans
		reflection.Register(baseServer)

		pb.RegisterAuthServer(baseServer, new(controller.Auth))
		return baseServer.Serve(lis)
	}, func(error) {
		lis.Close()
	})
}

// init HTTP Endpoint
// add rest endpoints
func (a *App) initHTTPEndpoint(g *group.Group) {
	defer fmt.Printf("metrics started port:%s \n", config.Global.Service.HTTP.Port)

	g.Add(func() error {
		if err := router.Router.GetRouter().Run(config.Global.Service.HTTP.Port); err != nil {
			zapLogger.Prepare(logger).Development().Level(zap.InfoLevel).Add("msg", "transport debug/HTTP during Listen err").Commit(err.Error())
			return err
		}
		return nil
	}, func(error) {})
}

// init cancle Interrupt
func (a *App) initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

// init jaeger tracer
func (a *App) initJaeger() (io.Closer, error) {
	defer fmt.Printf("Jaeger loaded successfully \n")
	closer, err := jtrace.Tracer.Connect()
	if err != nil {
		return nil, err
	}

	return closer, nil
}

// in production you load envs from etcd storage
// you can change, add or delete watch keys
// watches example: key: redis - value: {"password":"****","address":"***:6985","db":"0",....}
func (a *App) initConfigServer() error {
	defer fmt.Printf("etcd storage loaded successfully \n")
	if err := etcd.Storage.Connect(); err != nil {
		return err
	}

	// loop over watchList
	for _, key := range config.Global.ETCD.WatchList {

		// get configs for first time on app starts
		err := etcd.Storage.GetKey(context.Background(), key, func(kv *mvccpb.KeyValue) {
			// set configs from storage to struct - if exists in Set method
			config.Global.Set(string(kv.Key), kv.Value)

		}, clientv3.WithPrefix())
		if err != nil {
			return err
		}

		// start to watch keys
		etcd.Storage.WatchKey(context.Background(), key, func(e *clientv3.Event) {

			// set configs from storage to struct - if exists in Set method
			config.Global.Set(string(e.Kv.Key), e.Kv.Value)

		}, clientv3.WithPrefix())
	}

	// apply service discovery - put service details
	return etcd.Storage.Put(context.Background(), config.Global.Service.Name, config.Global.GetService())
}

// init mysql database
func (a *App) initDatabase() error {
	fmt.Printf("mysql storage loaded successfully \n")
	return mysql.Storage.Connect(config.Global)
}

// init message broker
func (a *App) initMessageBroker() error {
	fmt.Printf("nats message broker loaded successfully \n")
	if err := broker.Nats.Connect(); err != nil {
		return err
	}

	return nil

}

// init Redis database
func (a *App) initRedis() error {
	fmt.Printf("redis database loaded successfully \n")
	if err := redis.Storage.Connect(config.Global); err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

// init Consul
func (a *App) initConsul() error {
	fmt.Printf("Consul loaded successfully \n")
	if err := consul.Storage.Connect(config.Global); err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

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
