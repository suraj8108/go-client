package main

import (
	"github.com/suraj8108/clientApp/client"
	"github.com/suraj8108/clientApp/redisutils"
)

func main() {

	// Create Redis client
	redisHandler := redisutils.NewRedis()

	// Data 5L data in Redis
	// dataseed.SeedDataInRedis(redisHandler)

	// rawValue, _ := redisHandler.FetchDataFromRedis(20)
	// utils.UnMarshalRedisData(rawValue)

	// keys := []string{"record20", "record10"}
	// rawValues, _ := redisHandler.FetchBulkDataFromRedis(keys)
	// utils.UnMarshalRedisBulkData(rawValues)

	connectionHandler := client.NewConnHandler(redisHandler)
	connectionHandler.ClientOperation()

}
