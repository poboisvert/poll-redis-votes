package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/poboisvert/poll-redis-pubsub/config"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func ConnectRedis() error {
	redisConfig := config.GetConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisConfig.RedisHost,
		Password: redisConfig.RedisPassword,
		DB:       0, // default database
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	}
	fmt.Println("Connected to Redis")
	return nil
}

func UpdateVoteCount(pollId int, optionIndex int) error {
	key := fmt.Sprintf("vote_count:%d", pollId)

	// Check if the vote count key exists in Redis.
	exists, err := rdb.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}

	// If the key does not exist, it means the poll does not exist, so we need to return an error.
	if exists == 0 {
		return errors.New("poll not found")
	}

	// Increment the vote count for the specified option.
	err = rdb.HIncrBy(context.Background(), key, strconv.Itoa(optionIndex), 1).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetVoteCount(pollId int) (map[int]int64, error) {
	key := fmt.Sprintf("vote_count:%d", pollId)

	// Get vote counts for each option.
	votes, err := rdb.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	voteCount := make(map[int]int64)
	for optionIndexString, voteString := range votes {
		optionIndex, err := strconv.Atoi(optionIndexString)
		if err != nil {
			return nil, err
		}
		vote, err := strconv.ParseInt(voteString, 10, 64)
		if err != nil {
			return nil, err
		}
		voteCount[optionIndex] = vote
	}

	return voteCount, nil
}
