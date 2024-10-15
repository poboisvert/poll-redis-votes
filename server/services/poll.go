package services

import (
	"encoding/json"
	"time"

	"github.com/poboisvert/poll-redis-pubsub/db"
	"github.com/poboisvert/poll-redis-pubsub/models"
)

func CreatePoll(poll *models.Poll) (*models.Poll, error) {
	var newPoll *models.Poll

	err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Add logic for validating the poll data (e.g., options, question)

	// Create new poll
	stmt, err := db.db.Prepare("INSERT INTO polls (question, options, total_votes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING ID, created_at, updated_at")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(poll.Question, json.Marshal(poll.Options), 0, time.Now().Format("2006-01-02T15:04:05Z"), time.Now().Format("2006-01-02T15:04:05Z")).Scan(&newPoll.ID, &newPoll.CreatedAt, &newPoll.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return newPoll, nil
}

func GetPollByID(pollID int) (*models.Poll, error) {
	var poll models.Poll

	err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Get poll by ID
	stmt, err := db.db.Prepare("SELECT ID, question, options, total_votes, created_at, updated_at FROM polls WHERE ID = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(pollID).Scan(&poll.ID, &poll.Question, &poll.Options, &poll.TotalVotes, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Unmarshal options from JSON string
	err = json.Unmarshal(poll.Options, &poll.Options)
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

func GetAllPolls() ([]models.Poll, error) {
	var polls []models.Poll

	err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Get all polls
	rows, err := db.db.Query("SELECT ID, question, options, total_votes, created_at, updated_at FROM polls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var poll models.Poll
		err = rows.Scan(&poll.ID, &poll.Question, &poll.Options, &poll.TotalVotes, &poll.CreatedAt, &poll.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Unmarshal options from JSON string
		err = json.Unmarshal(poll.Options, &poll.Options)
		if err != nil {
			return nil, err
		}

		polls = append(polls, poll)
	}

	return polls, nil
}

func UpdatePoll(pollID int, poll *models.Poll) (*models.Poll, error) {
	var updatedPoll *models.Poll

	err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Update poll
	stmt, err := db.db.Prepare("UPDATE polls SET question = $1, options = $2, total_votes = $3, updated_at = $4 WHERE ID = $5 RETURNING ID, created_at, updated_at")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(poll.Question, json.Marshal(poll.Options), poll.TotalVotes, time.Now().Format("2006-01-02T15:04:05Z"), pollID).Scan(&updatedPoll.ID, &updatedPoll.CreatedAt, &updatedPoll.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return updatedPoll, nil
}

func DeletePoll(pollID int) error {
	err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete poll
	stmt, err := db.db.Prepare("DELETE FROM polls WHERE ID = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pollID)
	if err != nil {
		return err
	}

	return nil
}
