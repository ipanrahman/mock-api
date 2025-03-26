package domain

type Response struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}
type MockResponse struct {
	Conditions []Condition `json:"conditions"`
	Response   Response    `json:"response"`
}
