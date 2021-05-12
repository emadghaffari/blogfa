package http

import "github.com/gin-gonic/gin"

var (
	Router auth = &Auth{}
)

// Auth service
type Auth struct{}

type auth interface {
	GetRouter() *gin.Engine
}
