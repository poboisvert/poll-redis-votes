package models

type Vote struct {
	PollID      int `json:"poll_id"`
	OptionIndex int `json:"option_index"`
}

// Example Postman payload for voting on a poll
/*
{
    "poll_id": 1,
    "option_index": 0
}
*/
