package structs

type Task struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Status  string `json:"status,omitempty"`
	Result  string `json:"result,omitempty"`
}
