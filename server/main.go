package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/poboisvert/poll-redis-pubsub/db"
	"github.com/poboisvert/poll-redis-pubsub/router"
	"github.com/poboisvert/poll-redis-pubsub/services"
)

func main() {
	// Connect to database
	err := db.Connect()
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}
	defer db.Close()

	// Connect to redis
	err = services.ConnectRedis()
	if err != nil {
		log.Fatal("Could not connect to redis: ", err)
	}

	r := router.CreateRouter()

	fmt.Printf("Starting server on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
