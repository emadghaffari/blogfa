package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// check service health
func Health(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}
