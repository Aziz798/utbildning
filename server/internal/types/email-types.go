package types

type Email struct {
	EmailTo  string `json:"email_to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	UserName string `json:"user_name"`
}
