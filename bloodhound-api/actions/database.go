package actions

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConf struct {
	MongoClient         *mongo.Client
	MongoDatabase       *mongo.Database
	SmellyServers       *mongo.Collection
	MasscanIPs          *mongo.Collection
	PlayerServerHistory *mongo.Collection
	ApiKeys             *mongo.Collection
}

var Database *DatabaseConf

func GetDatabase() *DatabaseConf {
	var err error

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DATABASE_URI")))
	if err != nil {
		panic(err)
	}

	mongoDatabase := mongoClient.Database("Bloodhound")

	smellyServers := mongoDatabase.Collection("SmellyServers")
	masscanIPs := mongoDatabase.Collection("MasscanIPs")
	playersServerHistory := mongoDatabase.Collection("PlayersServerHistory")
	apiKeys := mongoDatabase.Collection("ApiKeys")

	return &DatabaseConf{MongoClient: mongoClient, MongoDatabase: mongoDatabase, SmellyServers: smellyServers, MasscanIPs: masscanIPs, PlayerServerHistory: playersServerHistory, ApiKeys: apiKeys}
}
