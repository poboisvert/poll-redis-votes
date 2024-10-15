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

	// Enable CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	fmt.Printf("Starting server on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
