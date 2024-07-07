package main

import (
	"bet-app/config"
	"bet-app/handlers"
	"bet-app/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/add-money", handlers.AddMoney)
		protected.POST("/roll-dice", handlers.RollDice)
	}

	r.Run(":8080")
}
