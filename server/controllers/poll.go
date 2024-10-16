package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/poboisvert/poll-redis-pubsub/models"
	"github.com/poboisvert/poll-redis-pubsub/services"
)

var ErrPollNotFound = errors.New("poll not found") // Define the error

func CreatePoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization, content-type")

	var poll models.Poll
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newPoll, err := services.CreatePoll(&poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newPoll)
}

func GetAllPolls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	polls, err := services.GetAllPolls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(polls)
}

func GetPollByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pollID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	poll, err := services.GetPollByID(pollID)
	if err != nil {
		if errors.Is(err, ErrPollNotFound) { // Use the defined error
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(poll)
}

func UpdatePoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pollID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	var updatedPoll models.Poll
	if err := json.NewDecoder(r.Body).Decode(&updatedPoll); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedPollResult, err := services.UpdatePoll(pollID, &updatedPoll)
	if err != nil {
		if errors.Is(err, ErrPollNotFound) { // Use the defined error
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(updatedPollResult)
}

func DeletePoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pollID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	err = services.DeletePoll(pollID)
	if err != nil {
		if errors.Is(err, ErrPollNotFound) { // Use the defined error
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
