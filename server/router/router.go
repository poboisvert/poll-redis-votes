package router

import (
	"github.com/poboisvert/poll-redis-pubsub/controllers"
	_ "github.com/poboisvert/poll-redis-pubsub/services"

	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	r := mux.NewRouter()

	// Poll endpoints
	r.HandleFunc("/polls", controllers.CreatePoll).Methods("POST")
	r.HandleFunc("/polls/{pollID}", controllers.GetPollByID).Methods("GET")
	r.HandleFunc("/polls/{pollID}", controllers.UpdatePoll).Methods("PUT")
	r.HandleFunc("/polls/{pollID}", controllers.DeletePoll).Methods("DELETE")

	// Vote endpoints
	r.HandleFunc("/vote", controllers.Vote).Methods("POST")

	// Websocket endpoint for real-time updates
	//r.HandleFunc("/ws/{pollID}", controllers.WsHandler).Methods("GET")

	// Other real-time endpoints (Example: Get live poll results)
	r.HandleFunc("/results", controllers.GetAllPolls).Methods("GET")

	return r
}
