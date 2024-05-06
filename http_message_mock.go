package httpmessagemock

import (
	"net/http"
	"os"

	"github.com/jarcoal/httpmock"
	"github.com/robotjoosen/go-httpmessagemock/model"
)

var defaultTranslator TranslatorAware

type TranslatorAware interface {
	Parse(data []byte) (map[string]model.Message, error)
}

type MessageMock struct {
	translator TranslatorAware
	Messages   map[string]model.Message
}

func DefaultTranslator(Translator TranslatorAware) {
	defaultTranslator = Translator
}

func New(data []byte, options ...OptionFunc) (*MessageMock, error) {
	var err error
	m := &MessageMock{}

	for _, option := range options {
		option(m)
	}

	if m.translator == nil {
		m.translator = defaultTranslator
	}

	m.Messages, err = m.translator.Parse(data)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func NewFromFile(filePath string, options ...OptionFunc) (*MessageMock, error) {
	data, err := getData(filePath)
	if err != nil {
		return nil, err
	}

	return New(data, options...)
}

func (m *MessageMock) RegisterResponders(responses model.ResponseConfig) error {
	for label, message := range m.Messages {
		statusCode, exists := responses[label]
		if !exists {
			continue
		}

		httpmock.RegisterResponder(message.Request.Method, message.Request.URL, func(request *http.Request) (*http.Response, error) {
			response, ok := message.Response[statusCode]
			if !ok {
				return httpmock.NewStringResponse(http.StatusInternalServerError, "response not found"), nil
			}

			httpResponse := httpmock.NewStringResponse(statusCode, response.Body)

			for key, value := range response.Headers {
				httpResponse.Header.Set(key, value)
			}

			return httpResponse, nil
		})
	}

	return nil
}

func getData(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err // todo: join with custom error
	}

	return data, nil
}
