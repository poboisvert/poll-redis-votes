package models

// Poll represents the structure of a voting poll.
type Poll struct {
	ID         int      `json:"id"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	TotalVotes int      `json:"total_votes"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

// Example Postman payload for creating a new poll
/*
{
    "question": "What is your favorite programming language?",
    "options": ["Go", "Python", "JavaScript", "Java"],
    "total_votes": 0,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-01T12:00:00Z"
}
*/
