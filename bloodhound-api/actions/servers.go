package actions

import (
	"context"
	"example/bloodhound/models"
	"example/bloodhound/scopes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func ServersFindPlayer(c *gin.Context) {
	authenticated := authenticate(c, scopes.FIND_PLAYER)

	if !authenticated {
		c.JSON(http.StatusUnauthorized, unauthorized())
		return
	}

	player := c.Query("player")

	if player == "" {
		c.JSON(http.StatusBadRequest, badRequestMessage())
		return
	}

	filter := bson.D{{"player", player}}
	var playerServerHistory models.PlayerServerHistory

	Database.PlayerServerHistory.FindOne(context.TODO(), filter).Decode(&playerServerHistory)

	if playerServerHistory.Player == "" {
		c.JSON(http.StatusOK, playerNotFound())
		return
	} else {
		var servers []models.SmellyServer

		for _, address := range playerServerHistory.Servers {
			updatedServer := updateServer(address, false)

			for _, playerInUpdatedServer := range updatedServer.Payload.Players {
				if playerInUpdatedServer == player {
					servers = append(servers, updatedServer.Payload)
				}
			}
		}
		c.JSON(http.StatusOK, servers)
	}
}

func ServersFindAddress(c *gin.Context) {
	authenticated := authenticate(c, scopes.FIND_ADDRESS)

	if !authenticated {
		c.JSON(http.StatusUnauthorized, unauthorized())
		return
	}

	address := c.Query("address")

	if address == "" {
		c.JSON(http.StatusBadRequest, badRequestMessage())
		return
	}

	status := updateServer(address, false)

	c.JSON(http.StatusOK, status.Payload)
}

func ServersFindVersion(c *gin.Context) {
	authenticated := authenticate(c, scopes.FIND_VERSION)

	if !authenticated {
		c.JSON(http.StatusUnauthorized, unauthorized())
		return
	}

	version := c.Query("version")

	if version == "" {
		c.JSON(http.StatusBadRequest, badRequestMessage())
		return
	}

	filter := bson.D{{"version", version}}

	cursor, err := Database.SmellyServers.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}

	var servers []models.SmellyServer = []models.SmellyServer{}

	for cursor.RemainingBatchLength() > 0 {
		var server models.SmellyServer
		if cursor.Next(context.TODO()) {
			cursor.Decode(&server)
			servers = append(servers, server)
		}
	}

	c.JSON(http.StatusOK, servers)
}

func ServersFindPlayerHistory(c *gin.Context) {
	authenticated := authenticate(c, scopes.FIND_PLAYER_HISTORY)

	if !authenticated {
		c.JSON(http.StatusUnauthorized, unauthorized())
		return
	}

	player := c.Query("player")

	if player == "" {
		c.JSON(http.StatusBadRequest, badRequestMessage())
		return
	}

	filter := bson.D{{"player", player}}
	var playerServerHistory models.PlayerServerHistory

	Database.PlayerServerHistory.FindOne(context.TODO(), filter).Decode(&playerServerHistory)

	if playerServerHistory.Player == "" {
		c.JSON(http.StatusOK, playerNotFound())
	} else {
		c.JSON(http.StatusOK, playerServerHistory.Servers)
	}
}
