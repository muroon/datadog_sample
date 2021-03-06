package jsonmodel

type Message struct {
	ID   uint32 `json:"id"`
	Text string `json:"text"`
}

type PostMessageResult struct {
	Status bool `json:status`
}

type CustomResult struct {
	Status bool `json:status`
}

type ErrorResult struct {
	Message string `json:message`
}
