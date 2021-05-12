package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// check service health
func (a *Auth) Health(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}
