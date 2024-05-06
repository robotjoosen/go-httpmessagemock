package translator_test

import (
	"net/http"
	"testing"

	"github.com/robotjoosen/go-httpmessagemock"
	"github.com/robotjoosen/go-httpmessagemock/model"
	"github.com/robotjoosen/go-httpmessagemock/translator"
	"github.com/stretchr/testify/suite"
)

type HttpRequestSuite struct {
	suite.Suite
}

func TestHttpRequestSuite(t *testing.T) {
	suite.Run(t, new(HttpRequestSuite))
}

func (s *HttpRequestSuite) TestParse() {
	tcs := map[string]struct {
		withMessage     string
		expectedError   error
		expectEndpoints map[string]model.Message
	}{
		"basic message": {
			withMessage: "### Simple Request\nGET http://api.test\n\nHTTP/1.1 200 OK",
			expectEndpoints: map[string]model.Message{
				"Simple Request": {
					Request: model.Request{
						URL:    "http://api.test",
						Method: http.MethodGet,
					},
					Response: map[int]model.Response{
						200: {Headers: make(map[string]string)},
					},
				},
			},
		},
		"header only message": {
			withMessage: "### Simple Request\nGET http://api.test\n\nHTTP/1.1 200 OK\nContent-Type: application/json",
			expectEndpoints: map[string]model.Message{
				"Simple Request": {
					Request: model.Request{
						URL:    "http://api.test",
						Method: http.MethodGet,
					},
					Response: map[int]model.Response{
						200: {Headers: map[string]string{"Content-Type": "application/json"}},
					},
				},
			},
		},
		"json body only message": {
			withMessage: "### Simple Request\nGET http://api.test\n\nHTTP/1.1 200 OK\n\n{\"msg\":\"hello world\"}",
			expectEndpoints: map[string]model.Message{
				"Simple Request": {
					Request: model.Request{
						URL:    "http://api.test",
						Method: http.MethodGet,
					},
					Response: map[int]model.Response{
						200: {
							Headers: map[string]string{},
							Body:    "{\"msg\":\"hello world\"}",
						},
					},
				},
			},
		},
		"string body only message": {
			withMessage: "### Simple Request\nGET http://api.test\n\nHTTP/1.1 200 OK\n\nHello World",
			expectEndpoints: map[string]model.Message{
				"Simple Request": {
					Request: model.Request{
						URL:    "http://api.test",
						Method: http.MethodGet,
					},
					Response: map[int]model.Response{
						200: {
							Headers: make(map[string]string),
							Body:    "Hello World",
						},
					},
				},
			},
		},
		"header and body message": {
			withMessage: "### Simple Request\nGET http://api.test\n\nHTTP/1.1 200 OK\nContent-Type: application/json\n\n{\"msg\":\"hello world\"}",
			expectEndpoints: map[string]model.Message{
				"Simple Request": {
					Request: model.Request{
						URL:    "http://api.test",
						Method: http.MethodGet,
					},
					Response: map[int]model.Response{
						200: {
							Headers: map[string]string{"Content-Type": "application/json"},
							Body:    "{\"msg\":\"hello world\"}",
						},
					},
				},
			},
		},
		"multi response": {
			withMessage: "### Simple Request\nGET http://api.test\nHTTP/1.1 200 OK\n\nHTTP/1.1 401 Unauthorized",
			expectEndpoints: map[string]model.Message{
				"Simple Request": {
					Request: model.Request{
						URL:    "http://api.test",
						Method: http.MethodGet,
					},
					Response: map[int]model.Response{
						200: {Headers: make(map[string]string)},
						401: {Headers: make(map[string]string)},
					},
				},
			},
		},
		"regex request": {
			withMessage: "### Regex URL\nGET =~^http://api.test/items/\\d+\nHTTP/1.1 200 OK\n\nHTTP/1.1 404 Not Found",
			expectEndpoints: map[string]model.Message{
				"Regex URL": {
					Request: model.Request{
						URL:    `=~^http://api.test/items/\d+`,
						Method: http.MethodGet,
					},
					Response: map[int]model.Response{
						200: {Headers: make(map[string]string)},
						404: {Headers: make(map[string]string)},
					},
				},
			},
		},
	}

	httpmessagemock.DefaultTranslator(translator.NewHTTPTranslator())

	for name, tc := range tcs {
		s.Run(name, func() {
			mock, err := httpmessagemock.New([]byte(tc.withMessage))
			s.Assert().Equal(tc.expectedError, err)
			s.Assert().Equal(tc.expectEndpoints, mock.Messages)
		})
	}
}
