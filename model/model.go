package model

type Message struct {
	Request  Request
	Response map[int]Response `json:"response"`
}

type Request struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type Response struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type ResponseConfig map[string]int
