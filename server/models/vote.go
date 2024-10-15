package models

type Vote struct {
	PollID      int `json:"poll_id"`
	OptionIndex int `json:"option_index"`
}
