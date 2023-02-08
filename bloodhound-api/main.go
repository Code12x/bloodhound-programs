package main

import (
	"example/bloodhound/actions"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	actions.Database = actions.GetDatabase()

	router := gin.Default()
	v1 := router.Group("/api/v1.0/servers")

	servers := v1.GET("/servers") // USE AUTH MIDDLEWARE
	servers.GET("/find-player", actions.ServersFindPlayer)
	servers.GET("/find-address", actions.ServersFindAddress)
	servers.GET("/find-version", actions.ServersFindVersion)
	servers.GET("/find-player-history", actions.ServersFindPlayerHistory)

	operation := v1.GET("/utilities") // USE AUTH MIDDLEWARE
	operation.GET("/test-ips-in-queue", actions.TestIpsInQueue)

	router.Run("localhost:3000")
}
