## poll-redis-pubsub

This project is a Go application that allows users to create polls and submit votes in real-time using Redis for pub/sub messaging. It provides a simple API for managing polls and recording votes.

- Routes:
  -- go get github.com/go-redis/redis/v8
  -- go get github.com/lib/pq
  -- go get github.com/gorilla/websocket

### Vote

```
{
    "poll_id": 1,
    "option_index": 0
}
```

### Create Pool

```
{
"question": "What is your favorite programming language?",
"options": ["Go", "Python", "JavaScript", "Java"],
"total_votes": 0,
"created_at": "2023-10-01T12:00:00Z",
"updated_at": "2023-10-01T12:00:00Z"
}
```
