package respond

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EmptySuccess(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func Error(c *gin.Context, code int, label, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"label":   label,
		"message": message,
	})
}
