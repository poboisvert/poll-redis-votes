package router

import (
	"github.com/poboisvert/poll-redis-pubsub/controllers"
	"github.com/poboisvert/poll-redis-pubsub/utils"

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
	r.HandleFunc("/votes", controllers.Vote).Methods("POST")

	// Websocket endpoint for real-time updates
	r.HandleFunc("/ws/{pollID}", utils.WsHandler).Methods("GET")

	// Other real-time endpoints (Example: Get live poll results)
	r.HandleFunc("/polls", controllers.GetAllPolls).Methods("GET")

	return r
}
