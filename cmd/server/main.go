package main

import (
	"airbook-backend/config"
	"airbook-backend/database"
	"airbook-backend/middleware"
	"airbook-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	routes.RegisterRoutes(r)

	r.Run(":8080")
}
