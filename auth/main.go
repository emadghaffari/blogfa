package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	"github.com/emadghaffari/blogfa/auth/handler"
	pb "github.com/emadghaffari/blogfa/auth/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("auth"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterAuthHandler(srv.Server(), new(handler.Auth))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
