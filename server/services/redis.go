package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/poboisvert/poll-redis-pubsub/config"
	"github.com/poboisvert/poll-redis-pubsub/db"

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
		// Create the key for the poll in Redis
		rdb.HSet(context.Background(), key).Err()

	}

	// If the key does not exist, it means the poll does not exist, so we need to return an error.
	if exists == 0 {
		// Create the key for the poll in Redis
		rdb.HSet(context.Background(), key).Err()
	}

	// Increment the vote count for the specified option.
	err = rdb.HIncrBy(context.Background(), key, strconv.Itoa(optionIndex), 1).Err()
	if err != nil {
		return err // Return the error instead of nil
	}
	voteCount, err := GetVoteCount(pollId)
	if err != nil {
		return err // Return the error if getting vote count fails
	}

	totalVotes := int64(0)
	for _, count := range voteCount {
		totalVotes += count
	}
	fmt.Printf("totalVotes: %d\n", totalVotes)

	// Update the total_votes in the PostgreSQL polls table
	dbInstance := db.GetDB()
	stmt, err := dbInstance.Prepare("UPDATE polls SET total_votes = $1 WHERE ID = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(totalVotes, pollId)
	if err != nil {
		return err // Return the error if the update fails
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
