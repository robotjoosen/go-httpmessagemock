package httpmessagemock_test

import (
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/robotjoosen/go-httpmessagemock"
	"github.com/robotjoosen/go-httpmessagemock/model"
	"github.com/robotjoosen/go-httpmessagemock/translator"
)

func ExampleDefaultTranslator() {
	httpmock.Activate()
	httpmessagemock.DefaultTranslator(translator.NewHTTPTranslator())

	messageMock, err := httpmessagemock.NewFromFile("./translator/http.example.http")
	if err != nil {
		panic(err)
	}

	err = messageMock.RegisterResponders(model.ResponseConfig{
		"Get items": http.StatusOK,
		"Get item":  http.StatusNotFound,
	})
	if err != nil {
		panic(err)
	}

	// output:
}

func ExampleWithTranslatorOption() {
	httpmock.Activate()

	messageMock, err := httpmessagemock.NewFromFile(
		"./translator/mockoon.example.json",
		httpmessagemock.WithTranslatorOption(translator.NewMockoonTranslator()),
	)
	if err != nil {
		panic(err)
	}

	err = messageMock.RegisterResponders(model.ResponseConfig{
		"items":    http.StatusOK,
		"item/:id": http.StatusNotFound,
	})
	if err != nil {
		panic(err)
	}

	// output:
}
