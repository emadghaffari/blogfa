package app

import (
	"blogfa/auth/broker"
	"blogfa/auth/config"
	"blogfa/auth/database/etcd"
	"blogfa/auth/database/mysql"
	"blogfa/auth/database/redis"
	"blogfa/auth/pkg/jtrace"
	zapLogger "blogfa/auth/pkg/logger"
	pb "blogfa/auth/proto"
	"blogfa/auth/service"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	group "github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logger *zap.Logger

// StartApplication func
func StartApplication() {
	fmt.Println("--------------------------------")

	// if go code crashed we get error and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init zap logger
	initLogger()

	// init configs
	if err := initConfigs(); err != nil {
		return
	}

	//  Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency
	closer, err := initJaeger()
	if err != nil {
		return
	}
	defer closer.Close()

	if viper.GetString("environment") == "production" {
		initConfigServer()
		defer etcd.Storage.GetClient().Close()
	}

	if err := initDatabase(); err != nil {
		return
	}

	if err := initRedis(); err != nil {
		return
	}

	if err := initMessageBroker(); err != nil {
		return
	}
	defer broker.Nats.Conn().Close()
	defer broker.Nats.ECConn().Close()

	g := createService()
	initMetricsEndpoint(g)
	initCancelInterrupt(g)

	fmt.Println("--------------------------------")
	if err := g.Run(); err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit("server stopped")
	}
}

// init zap logger
func initLogger() {
	defer fmt.Printf("zap logger is available \n")
	zapLogger.SetLogPath("logs")
	logger = zapLogger.GetZapLogger(false)
}

// init configs
func initConfigs() error {
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
func initGRPCHandler(g *group.Group) {
	defer fmt.Printf("grpc connected port:%s \n", config.Global.Service.GRPC.Port)

	options := defaultGRPCOptions(logger, jtrace.Tracer.GetTracer())
	// Add your GRPC options here

	lis, err := net.Listen("tcp", config.Global.Service.GRPC.Port)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.ErrorLevel).Commit(err.Error())
	}

	g.Add(func() error {
		baseServer := grpc.NewServer(options...)

		// reflection for evans
		reflection.Register(baseServer)

		pb.RegisterAuthServer(baseServer, new(service.Auth))
		return baseServer.Serve(lis)
	}, func(error) {
		lis.Close()
	})
}

// init metrics
func initMetricsEndpoint(g *group.Group) {
	defer fmt.Printf("metrics started port:%s \n", config.Global.Service.HTTP.Port)

	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	debugListener, err := net.Listen("tcp", config.Global.Service.HTTP.Port)
	if err != nil {
		zapLogger.Prepare(logger).Development().Level(zap.InfoLevel).Add("msg", "transport debug/HTTP during Listen err").Commit(err.Error())
	}
	g.Add(func() error {
		return http.Serve(debugListener, http.DefaultServeMux)
	}, func(error) {
		debugListener.Close()
	})
}

// init cancle Interrupt
func initCancelInterrupt(g *group.Group) {
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
func initJaeger() (io.Closer, error) {
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
func initConfigServer() {
	fmt.Printf("etcd storage loaded successfully \n")
	if err := etcd.Storage.Connect(); err != nil {
		return
	}

	// loop over watchList
	for _, key := range config.Global.ETCD.WatchList {

		// get configs for first time on app starts
		etcd.Storage.GetKey(context.Background(), key, func(kv *mvccpb.KeyValue) {
			// set configs from storage to struct - if exists in Set method
			config.Global.Set(string(kv.Key), kv.Value)

		}, clientv3.WithPrefix())

		// start to watch keys
		etcd.Storage.WatchKey(context.Background(), key, func(e *clientv3.Event) {

			// set configs from storage to struct - if exists in Set method
			config.Global.Set(string(e.Kv.Key), e.Kv.Value)

		}, clientv3.WithPrefix())
	}

	// apply service discovery - put service details
	etcd.Storage.Put(context.Background(), config.Global.Service.Name, config.Global.GetService())
}

// init mysql database
func initDatabase() error {
	fmt.Printf("mysql storage loaded successfully \n")
	return mysql.Storage.Connect(config.Global)
}

// init message broker
func initMessageBroker() error {
	fmt.Printf("nats message broker loaded successfully \n")
	if err := broker.Nats.Connect(); err != nil {
		return err
	}

	if err := broker.Nats.EncodedConn(); err != nil {
		return err
	}

	return nil

}

// init Redis database
func initRedis() error {
	fmt.Printf("redis database loaded successfully \n")
	if err := redis.Storage.Connect(config.Global); err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
