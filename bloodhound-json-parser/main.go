package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Host struct {
	Ip string `json:"ip"`
}

var (
	database *mongo.Collection
	wg       = &sync.WaitGroup{}
)

func main() {
	godotenv.Load()

	fmt.Print("Enter the json scan file to test the ips: ")
	var fileName string
	fmt.Scan(&fileName)

	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DATABASE_URI")))
	if err != nil {
		panic(err)
	}
	database = mongoClient.Database("Bloodhound").Collection("MasscanIPs")

	hosts := []Host{}
	json.Unmarshal(file, &hosts)

	if len(hosts) > 20 {
		var hostsList [][]Host
		const NUM_OF_GOROUTINES int = 100
		var hostsPerRoutine int = len(hosts) / NUM_OF_GOROUTINES

		var goroutineIndex int = 0
		for goroutineIndex < int(NUM_OF_GOROUTINES) {
			var hostList []Host
			for i, host := range hosts {
				if goroutineIndex == NUM_OF_GOROUTINES-1 && i >= goroutineIndex*int(hostsPerRoutine) {
					hostList = append(hostList, host)
				} else if i >= goroutineIndex*int(hostsPerRoutine) && i < (goroutineIndex+1)*int(hostsPerRoutine) {
					hostList = append(hostList, host)
				}
			}
			hostsList = append(hostsList, hostList)
			goroutineIndex++
		}
		for _, hostList := range hostsList {
			wg.Add(1)
			go addData(hostList)
		}
	} else {
		wg.Add(1)
		go addData(hosts)
	}
	wg.Wait()
	fmt.Println("Success!")
}

func addData(data []Host) {
	defer wg.Done()
	for _, i := range data {
		database.InsertOne(context.TODO(), i)
		fmt.Println("Inserted One!")
	}
}
