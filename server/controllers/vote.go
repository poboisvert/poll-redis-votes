package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/poboisvert/poll-redis-pubsub/models"
	"github.com/poboisvert/poll-redis-pubsub/services"
)

func Vote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var vote models.Vote

	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := services.UpdateVoteCount(vote.PollID, vote.OptionIndex)
	if err != nil {
		http.Error(w, "Could not update vote count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(http.StatusOK)
}
