package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePanic(c *gin.Context) {
	if err := recover(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
}
