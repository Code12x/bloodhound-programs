package actions

import (
	"context"
	"example/bloodhound/models"
	"example/bloodhound/scopes"
	"fmt"
	"time"

	"github.com/alteamc/minequery/v2"
	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.mongodb.org/mongo-driver/bson"
)

func updateServer(address string, isTest bool) serverResponse {
	res, err := minequery.Ping17(address, 25565)
	if err != nil {
		Database.MasscanIPs.DeleteOne(context.TODO(), bson.D{{"ip", address}})
		currentNumOfRoutines--
		fmt.Println("nope:", address)
		return serverNotFound()
	}

	// Logging the address into the db.
	version := res.VersionName
	var players []string
	for _, player := range res.SamplePlayers {
		players = append(players, player.Nickname)

		if isTest {
			if player.Nickname == "LiveOverflow" {
				sendSmsMessage(player.Nickname, address, players)
			}

			if player.Nickname == "LiveUnderflow" {
				sendSmsMessage(player.Nickname, address, players)
			}
		}
	}

	currentTime := time.Now().Format(time.RFC822)
	// Checking if the server is already in the db
	filter := bson.D{{"address", address}}

	dbSearchRes := Database.SmellyServers.FindOne(context.TODO(), filter)

	var server models.SmellyServer
	dbSearchRes.Decode(&server)

	var updatedServer models.SmellyServer

	updatedServer.Address = address
	updatedServer.Version = version
	updatedServer.Players = players
	updatedServer.DateUpdated = currentTime

	if server.Address == "" {
		updatedServer.DateCreated = currentTime
		Database.SmellyServers.InsertOne(context.TODO(), updatedServer)

	} else {
		newFilter := bson.D{{"address", server.Address}}

		if server.DateCreated != "" {
			updatedServer.DateCreated = server.DateCreated
			Database.SmellyServers.ReplaceOne(context.TODO(), newFilter, updatedServer)
		} else {
			updatedServer.DateCreated = currentTime
			Database.SmellyServers.ReplaceOne(context.TODO(), newFilter, updatedServer)
		}
	}

	for _, player := range updatedServer.Players {
		filter := bson.D{{"player", player}}

		var playerServerHistory models.PlayerServerHistory
		Database.PlayerServerHistory.FindOne(context.TODO(), filter).Decode(&playerServerHistory)

		if len(playerServerHistory.Servers) <= 0 {
			Database.PlayerServerHistory.InsertOne(context.TODO(), models.PlayerServerHistory{Player: player, Servers: []string{server.Address}})
		} else {
			hasCurrentServer := false

			for _, currentAddress := range playerServerHistory.Servers {
				if currentAddress == updatedServer.Address {
					hasCurrentServer = true
					break
				}
			}

			if !hasCurrentServer {
				playerServerHistory.Servers = append(playerServerHistory.Servers, updatedServer.Address)
			}
		}
	}

	if isTest {
		fmt.Printf("=============== Minecraft Server Found!\tAddress: %v\tPlayers Online: %v ========================\n", address, res.OnlinePlayers)
		Database.MasscanIPs.DeleteOne(context.TODO(), bson.D{{"ip", address}})
		currentNumOfRoutines--
	}
	return serverUpdateSuccess(updatedServer)
}

func authenticate(c *gin.Context, callScope string) bool {
	keyPassed := c.Query("apikey")

	if keyPassed == "" {
		return false
	}

	var apiKey models.ApiKeys
	filter := bson.D{{"key", keyPassed}}

	Database.ApiKeys.FindOne(context.TODO(), filter).Decode(&apiKey)

	inScope := false

	for _, scope := range apiKey.Scope {
		if scope == scopes.ALL {
			inScope = true
			break
		}
		if callScope == scope {
			inScope = true
			break
		}
	}

	if inScope {
		return true
	} else {
		return false
	}
}

func sendSmsMessage(playerName string, ip string, players []string) {
	client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}
	params.SetBody(fmt.Sprintf("%v FOUND!\nIP: %v\nPlayers Online: %v", playerName, ip, players))
	params.SetFrom("<removed>")
	params.SetTo("<removed>")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}
}
