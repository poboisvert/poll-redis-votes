package utils

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/poboisvert/poll-redis-pubsub/services" // Import the services package for GetVoteCount
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development; tighten in production
	},
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

// 1. Establish a Websocket connection for each pollID requested by the frontend.
// 2. Continuously monitor and send real-time vote count updates to the connected client.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pollID, err := strconv.Atoi(vars["pollID"])

	if err != nil {
		log.Println("Invalid poll ID:", vars["pollID"])
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	conn, err := Upgrade(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	// Continuously send real-time updates with vote counts
	log.Println("Valid poll ID:", pollID)
	for {
		voteCount, err := services.GetVoteCount(pollID) // Use GetVoteCount from the services package
		if err != nil {
			log.Println("Error getting vote count:", err)
			break // Exit the loop on error
		}

		message := fmt.Sprintf(`{"voteCount": %d}`, voteCount)
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error writing to WebSocket:", err)
			return
		}

	}
}
