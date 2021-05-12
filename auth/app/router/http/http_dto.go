package http

import (
	"blogfa/auth/controller/http"

	"github.com/gin-gonic/gin"
)

func (a *Auth) GetRouter() *gin.Engine {
	router := gin.Default()

	// metrics
	router.GET("/metrics", http.Controller.Metrics)

	// health check
	router.GET("/health", http.Controller.Health)

	return router
}
