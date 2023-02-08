package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dreamscached/minequery/v2"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Host struct {
	Ip string `bson:"ip"`
}

type SmellyServer struct {
	Address     string   `bson:address`
	Version     string   `bson:version`
	Players     []string `bson:players`
	DateCreated string   `bson:datecreated`
	DateUpdated string   `bson:dateupdated`
}

var (
	masscanIps           *mongo.Collection
	smellyServers        *mongo.Collection
	wg                   = &sync.WaitGroup{}
	currentNumOfRoutines = 0
)

const (
	MAX_NUM_OF_ROUTINES = 10000
)

func main() {
	godotenv.Load()

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DATABASE_URI")))
	if err != nil {
		panic(err)
	}
	masscanIps = mongoClient.Database("Bloodhound").Collection("MasscanIPs")
	smellyServers = mongoClient.Database("Bloodhound").Collection("SmellyServers")

	batchNum := 0

	for cursor, err := masscanIps.Find(context.TODO(), bson.D{}); cursor.RemainingBatchLength() >= 0; {
		fmt.Printf("Batch Number: %v\n", batchNum)
		if err != nil {
			panic(err)
		}
		cursor.Next(context.TODO())

		for cursor.RemainingBatchLength() > 0 {
			fmt.Println("Remaining:", cursor.RemainingBatchLength())
			if currentNumOfRoutines < MAX_NUM_OF_ROUTINES {
				fmt.Printf("IP: %v\n", cursor.Current.Lookup("ip").StringValue())
				go testIp(cursor.Current.Lookup("ip").StringValue())
				wg.Add(1)
				currentNumOfRoutines++
			} else {
				fmt.Println("waiting")
				time.Sleep(time.Duration(time.Millisecond) * 250)
			}
			cursor.Next(context.TODO())
		}
		batchNum++
	}
	wg.Wait()
}

func testIp(ip string) {
	defer wg.Done()

	pinger := minequery.NewPinger(
		minequery.WithTimeout(5*time.Second),
		minequery.WithUseStrict(true),
	)

	response, err := pinger.Ping17(ip, 25565)
	if err != nil {
		masscanIps.DeleteOne(context.TODO(), bson.D{{"ip", ip}})
		currentNumOfRoutines--
		return
	}

	version := response.VersionName
	minequeryPlayers := response.SamplePlayers
	var players []string
	dateCreated := time.Now().Format(time.RFC822)
	dateUpdated := time.Now().Format(time.RFC822)

	for _, player := range minequeryPlayers {
		players = append(players, player.Nickname)

		if player.Nickname == "LiveOverflow" {
			sendSmsMessage(player.Nickname, ip, players)
		}

		if player.Nickname == "LiveUnderflow" {
			sendSmsMessage(player.Nickname, ip, players)
		}
	}

	var smellyServer SmellyServer = SmellyServer{Address: ip, Version: version, Players: players, DateCreated: dateCreated, DateUpdated: dateUpdated}

	smellyServers.InsertOne(context.TODO(), smellyServer)
	fmt.Printf("=============== Minecraft Server Found!\tAddress: %v ========================\n", ip)
	masscanIps.DeleteOne(context.TODO(), bson.D{{"ip", ip}})
	currentNumOfRoutines--
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
