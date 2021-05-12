package http

import "github.com/gin-gonic/gin"

var (
	Controller auth = &Auth{}
)

type auth interface {
	Health(c *gin.Context)
	Metrics(c *gin.Context)
}

// Auth service
type Auth struct{}
