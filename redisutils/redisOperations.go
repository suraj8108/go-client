package redisutils

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	client *redis.Client
}

func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// pingResp, err := client.Ping().Result()
	return client
}

func NewRedis() RedisClient {
	client := createClient()
	return RedisClient{client: client}
}

func (rc *RedisClient) InsertDataInRedis(currentKey int, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	// ctx := context.Background()

	key := fmt.Sprintf("record%v", currentKey)
	err = rc.client.Set(key, json, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Seeding data %v \n", currentKey)
	return nil
}

func (rc *RedisClient) InsertDataInRedisBySelfKey(key string, value interface{}, sc *sync.WaitGroup) error {
	defer sc.Done()
	json, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	// ctx := context.Background()

	err = rc.client.Set(key, json, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Seeding data %v \n", key)
	fmt.Println("Inserting data in redis", value)
	return nil
}

func (rc *RedisClient) FetchDataFromRedis(currentKey int) (string, error) {

	key := fmt.Sprintf("record%v", currentKey)
	val, err := rc.client.Get(key).Result()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(val)
	return val, nil
}

func (rc *RedisClient) FetchBulkDataFromRedis(keys []string) ([]interface{}, error) {
	vals, err := rc.client.MGet(keys...).Result()
	fmt.Println("Redis data from Redis Data Store", vals)
	return vals, err
}
