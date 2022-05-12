package routes

import (
	"zate/pkg/auth/handler"
	"zate/pkg/auth/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/v1/auth/signup", handler.SignUp)
	router.Use(middlewares.VerifyEmailAuthentication())
	router.POST("/v1/auth/signup/verify", handler.VerifyEmail)
}
