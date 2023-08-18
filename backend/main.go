package main

import (
	"github.com/RugeFX/ruge-chat-app/database"
	"github.com/RugeFX/ruge-chat-app/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if envErr := godotenv.Load(); envErr != nil {
		panic(envErr)
	}
	database.ConnectDB()

	r := gin.Default()
	r.Use(cors.Default())

	routes.SetupRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	r.Run(":3000")
}
