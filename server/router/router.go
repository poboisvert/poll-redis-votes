package router

import (
	// Import log package
	"net/http"

	"github.com/poboisvert/poll-redis-pubsub/controllers"
	"github.com/poboisvert/poll-redis-pubsub/utils"

	"github.com/gorilla/mux"
)

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the original handler
		h.ServeHTTP(w, r)
	}
}

func CreateRouter() *mux.Router {
	r := mux.NewRouter()

	// Poll endpoints
	r.HandleFunc("/polls", corsHandler(controllers.CreatePoll)).Methods("POST")
	r.HandleFunc("/polls/{pollID}", corsHandler(controllers.GetPollByID)).Methods("GET")
	r.HandleFunc("/polls/{pollID}", corsHandler(controllers.UpdatePoll)).Methods("PUT")
	r.HandleFunc("/polls/{pollID}", corsHandler(controllers.DeletePoll)).Methods("DELETE")

	// Vote endpoints
	r.HandleFunc("/votes", corsHandler(controllers.Vote)).Methods("POST")

	// Websocket endpoint for real-time updates
	r.HandleFunc("/ws/{pollID}", utils.WsHandler).Methods("GET")

	// Other real-time endpoints (Example: Get live poll results)
	r.HandleFunc("/polls", corsHandler(controllers.GetAllPolls)).Methods("GET")

	return r
}
