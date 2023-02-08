package actions

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

var (
	currentNumOfRoutines = 0
	MAX_NUM_OF_ROUTINES  = 10
)

func TestIpsInQueue(c *gin.Context) {
	go testIps()
	c.JSON(http.StatusOK, operationReceived())
}

func testIps() {
	cursor, err := Database.MasscanIPs.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	for cursor.RemainingBatchLength() >= 0 {
		cursor.Next(context.TODO())

		for cursor.RemainingBatchLength() >= 0 {
			//fmt.Println("Remaining:", cursor.RemainingBatchLength())
			if currentNumOfRoutines < MAX_NUM_OF_ROUTINES {
				ip := cursor.Current.Lookup("ip").StringValue()
				go updateServer(ip, true)
				currentNumOfRoutines++
			} else {
				//fmt.Println("waiting")
				time.Sleep(time.Duration(time.Millisecond) * 250)
			}
			cursor.Next(context.TODO())
		}
	}
}
