package models

type Poll struct {
	ID         int      `json:"id"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	TotalVotes int      `json:"total_votes"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}
