package handler

type Event struct {
}

type Response struct {
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}
