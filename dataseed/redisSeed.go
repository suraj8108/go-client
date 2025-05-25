package dataseed

import (
	"fmt"
	"sync"
	"time"

	"github.com/suraj8108/clientApp/model"
	"github.com/suraj8108/clientApp/redisutils"
)

const MAX_DATA_TO_SEED int = 500000

func seed10KData(start int, redisHandler redisutils.RedisClient, sc *sync.WaitGroup) {
	for i := 1; i <= 10000; i++ {
		currentKey := start + i
		student := model.Student{
			StudentName:     fmt.Sprintf("User%v", currentKey),
			StudentEmail:    fmt.Sprintf("user%v@gmail.com", currentKey),
			StudentRedisKey: fmt.Sprintf("record%v", currentKey),
		}
		err := redisHandler.InsertDataInRedis(currentKey, student)
		if err != nil {
			return
		}
	}
	sc.Done()
}

func SeedDataInRedis(redisHandler redisutils.RedisClient) {
	seedDataStart := 0
	sc := sync.WaitGroup{}
	start := time.Now()
	for {
		dataToSeed := 10000
		sc.Add(1)
		go seed10KData(seedDataStart, redisHandler, &sc)
		seedDataStart += dataToSeed

		sc.Add(1)
		go seed10KData(seedDataStart, redisHandler, &sc)
		seedDataStart += dataToSeed

		sc.Add(1)
		go seed10KData(seedDataStart, redisHandler, &sc)
		seedDataStart += dataToSeed

		sc.Add(1)
		go seed10KData(seedDataStart, redisHandler, &sc)
		seedDataStart += dataToSeed

		sc.Add(1)
		go seed10KData(seedDataStart, redisHandler, &sc)
		seedDataStart += dataToSeed

		// sc.Add(1)
		// go seed10KData(seedDataStart, redisHandler, &sc)
		// seedDataStart += dataToSeed

		if seedDataStart >= MAX_DATA_TO_SEED {
			break
		}
		fmt.Println("----------------------------------------------------------------")
	}
	sc.Wait()
	elapsed := time.Since(start)
	fmt.Println(seedDataStart)
	fmt.Printf("Total time taken is %v \n", elapsed)

}
