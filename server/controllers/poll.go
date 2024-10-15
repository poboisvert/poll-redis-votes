package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/poboisvert/poll-redis-pubsub/models"
	"github.com/poboisvert/poll-redis-pubsub/services"
)

func CreatePoll(w http.ResponseWriter, r *http.Request) {
	var poll *models.Poll
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newPoll, err := services.CreatePoll(poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newPoll)
}

func GetPollByID(w http.ResponseWriter, r *http.Request) {
	pollID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	poll, err := services.GetPollByID(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(poll)
}

func GetAllPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := services.GetAllPolls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(polls)
}

func UpdatePoll(w http.ResponseWriter, r *http.Request) {
	pollID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	var updatedPoll *models.Poll
	if err := json.NewDecoder(r.Body).Decode(&updatedPoll); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = services.UpdatePoll(pollID, updatedPoll)
	if err != nil {
		if errors.Is(err, models.ErrPollNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePoll(w http.ResponseWriter, r *http.Request) {
	pollID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	err = services.DeletePoll(pollID)
	if err != nil {
		if errors.Is(err, models.ErrPollNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
