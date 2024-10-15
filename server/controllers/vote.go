package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/poboisvert/poll-redis-pubsub/models"
	"github.com/poboisvert/poll-redis-pubsub/services"
)

func Vote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var vote models.Vote

	// Example JSON payload that will be decoded
	/*
		{
		    "poll_id": 1,
		    "option_index": 0
		}
	*/

	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("pollId: %d\n", vote.PollID)
	fmt.Printf("OptionIndex: %d\n", vote.OptionIndex)

	err := services.UpdateVoteCount(vote.PollID, vote.OptionIndex)
	if err != nil {
		http.Error(w, "Could not update vote count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(http.StatusOK)
}
