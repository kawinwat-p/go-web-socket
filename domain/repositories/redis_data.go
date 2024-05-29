package repositories

import (
	"context"
	"encoding/json"
	"bn-survey-point/domain/datasources"
	"bn-survey-point/domain/entities"
	"log"

	"github.com/go-redis/redis/v8"
)

type redisConnectionRepository struct {
	Context   context.Context
	RedisWR   *redis.Client
	RedisRead *redis.Client
}

type IRedisConnectionRepository interface {
	GetRedisAlertMessageData() *entities.AlertMessage
	SetRedisData(dataByte []byte) bool
}

func NewRedisRepository(redis *datasources.RedisConnection) IRedisConnectionRepository {
	return &redisConnectionRepository{
		Context:   redis.Context,
		RedisWR:   redis.RedisWR,
		RedisRead: redis.RedisRead,
	}
}

func (repo redisConnectionRepository) GetRedisAlertMessageData() *entities.AlertMessage {
	dataSpeakerRedis, err := repo.RedisRead.Get(repo.RedisRead.Context(), "AlertMessage").Result()
	if err != nil {
		log.Println("error GetUsersData ", err.Error())
		return nil
	}
	var data entities.AlertMessage
	json.Unmarshal([]byte(dataSpeakerRedis), &data)
	log.Println("Get users data to redis success!")
	return &data
}

func (repo redisConnectionRepository) SetRedisData(dataByte []byte) bool {
	err := repo.RedisWR.Set(repo.RedisWR.Context(), "AlertMessage", dataByte, 0).Err()
	if err != nil {
		log.Println("error SetUsersName ", err.Error())
		return false
	}
	log.Println("Set new users data to redis success!")
	return true
}
