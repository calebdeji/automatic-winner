package main

import (
	"os"
	"zate/pkg/auth/routes"
	"zate/services"

	"github.com/gin-gonic/gin"
)

func main() {

	services.LoadEnv()
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()

	router.Use(gin.Logger())

	routes.AuthRoutes(router)

	router.Run(":" + port)

}
