package services

import (
	"encoding/json"
	"time"

	"github.com/poboisvert/poll-redis-pubsub/db"
	"github.com/poboisvert/poll-redis-pubsub/models"
)

func CreatePoll(poll *models.Poll) (*models.Poll, error) {
	var newPoll models.Poll // Changed to a non-pointer type

	// No need to connect to the database here, we will use GetDB
	dbInstance := db.GetDB()

	// Create table if not exists
	_, err := dbInstance.Exec(`CREATE TABLE IF NOT EXISTS polls (
		ID SERIAL PRIMARY KEY,
		question TEXT NOT NULL,
		options JSONB NOT NULL,
		total_votes INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}

	// Create new poll
	stmt, err := dbInstance.Prepare("INSERT INTO polls (question, options, total_votes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING ID, created_at, updated_at")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	optionsJSON, err := json.Marshal(poll.Options)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(poll.Question, optionsJSON, 0, time.Now().Format("2006-01-02T15:04:05Z"), time.Now().Format("2006-01-02T15:04:05Z")).Scan(&newPoll.ID, &newPoll.CreatedAt, &newPoll.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &newPoll, nil // Return a pointer to newPoll
}

func GetPollByID(pollID int) (*models.Poll, error) {
	var poll models.Poll

	// No need to connect to the database here, we will use GetDB
	dbInstance := db.GetDB()

	// Get poll by ID
	stmt, err := dbInstance.Prepare("SELECT ID, question, options, total_votes, created_at, updated_at FROM polls WHERE ID = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(pollID).Scan(&poll.ID, &poll.Question, &poll.Options, &poll.TotalVotes, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

func GetAllPolls() ([]models.Poll, error) {
	var polls []models.Poll

	// No need to connect to the database here, we will use GetDB
	dbInstance := db.GetDB()

	// Get all polls
	rows, err := dbInstance.Query("SELECT ID, question, options, total_votes, created_at, updated_at FROM polls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var poll models.Poll
		var optionsJSON []byte
		err = rows.Scan(&poll.ID, &poll.Question, &optionsJSON, &poll.TotalVotes, &poll.CreatedAt, &poll.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(optionsJSON, &poll.Options); err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		polls = append(polls, poll)
	}

	return polls, nil
}

func UpdatePoll(pollID int, poll *models.Poll) (*models.Poll, error) {
	var updatedPoll models.Poll // Changed to a non-pointer type

	// No need to connect to the database here, we will use GetDB
	dbInstance := db.GetDB()

	// Update poll
	stmt, err := dbInstance.Prepare("UPDATE polls SET question = $1, options = $2, total_votes = $3, updated_at = $4 WHERE ID = $5 RETURNING ID, created_at, updated_at")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	optionsJSON, err := json.Marshal(poll.Options)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(poll.Question, optionsJSON, poll.TotalVotes, time.Now().Format("2006-01-02T15:04:05Z"), pollID).Scan(&updatedPoll.ID, &updatedPoll.CreatedAt, &updatedPoll.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedPoll, nil // Return a pointer to updatedPoll
}

func DeletePoll(pollID int) error {
	// No need to connect to the database here, we will use GetDB
	dbInstance := db.GetDB()

	// Delete poll
	stmt, err := dbInstance.Prepare("DELETE FROM polls WHERE ID = $1")
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
