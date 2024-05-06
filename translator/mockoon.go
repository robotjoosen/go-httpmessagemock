package translator

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/robotjoosen/go-httpmessagemock/model"
)

type MockoonData struct {
	Uuid           string        `json:"uuid"`
	LastMigration  int           `json:"lastMigration"`
	Name           string        `json:"name"`
	EndpointPrefix string        `json:"endpointPrefix"`
	Latency        int           `json:"latency"`
	Port           int           `json:"port"`
	Hostname       string        `json:"hostname"`
	Folders        []interface{} `json:"folders"`
	Routes         []struct {
		Uuid          string `json:"uuid"`
		Type          string `json:"type"`
		Documentation string `json:"documentation"`
		Method        string `json:"method"`
		Endpoint      string `json:"endpoint"`
		Responses     []struct {
			Uuid       string `json:"uuid"`
			Body       string `json:"body"`
			Latency    int    `json:"latency"`
			StatusCode int    `json:"statusCode"`
			Label      string `json:"label"`
			Headers    []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headers"`
			BodyType       string `json:"bodyType"`
			FilePath       string `json:"filePath"`
			DatabucketID   string `json:"databucketID"`
			SendFileAsBody bool   `json:"sendFileAsBody"`
			Rules          []struct {
				Target   string `json:"target"`
				Modifier string `json:"modifier"`
				Value    string `json:"value"`
				Invert   bool   `json:"invert"`
				Operator string `json:"operator"`
			} `json:"rules"`
			RulesOperator     string        `json:"rulesOperator"`
			DisableTemplating bool          `json:"disableTemplating"`
			FallbackTo404     bool          `json:"fallbackTo404"`
			Default           bool          `json:"default"`
			CrudKey           string        `json:"crudKey"`
			Callbacks         []interface{} `json:"callbacks"`
		} `json:"responses"`
		ResponseMode interface{} `json:"responseMode"`
	} `json:"routes"`
	RootChildren []struct {
		Type string `json:"type"`
		Uuid string `json:"uuid"`
	} `json:"rootChildren"`
	ProxyMode         bool   `json:"proxyMode"`
	ProxyHost         string `json:"proxyHost"`
	ProxyRemovePrefix bool   `json:"proxyRemovePrefix"`
	TlsOptions        struct {
		Enabled    bool   `json:"enabled"`
		Type       string `json:"type"`
		PfxPath    string `json:"pfxPath"`
		CertPath   string `json:"certPath"`
		KeyPath    string `json:"keyPath"`
		CaPath     string `json:"caPath"`
		Passphrase string `json:"passphrase"`
	} `json:"tlsOptions"`
	Cors    bool `json:"cors"`
	Headers []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"headers"`
	ProxyReqHeaders []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"proxyReqHeaders"`
	ProxyResHeaders []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"proxyResHeaders"`
	Data      []interface{} `json:"data"`
	Callbacks []interface{} `json:"callbacks"`
}

type MockoonTranslator struct{}

func NewMockoonTranslator() (t MockoonTranslator) {
	return
}

// Parse translates a Mockoon json into a message model.
// note: currently this only supports inline body content.
func (t MockoonTranslator) Parse(data []byte) (map[string]model.Message, error) {
	var mockoonData MockoonData

	if err := json.Unmarshal(data, &mockoonData); err != nil {
		return nil, err
	}

	messages := make(map[string]model.Message)
	for _, route := range mockoonData.Routes {
		responses := make(map[int]model.Response)
		for _, response := range route.Responses {
			headers := make(map[string]string)
			for _, header := range response.Headers {
				headers[header.Key] = header.Value
			}

			responses[response.StatusCode] = model.Response{
				Headers: headers,
				Body:    response.Body,
			}
		}

		url := mockoonData.Hostname + "/" + route.Endpoint
		if strings.Contains(route.Endpoint, ":") {
			regex := regexp.MustCompile(`:\S+?(\/|\z)`)
			url = "=~^" + mockoonData.Hostname + "/" + regex.ReplaceAllStringFunc(route.Endpoint, func(s string) string {
				rsp := `\S+\z`
				if strings.HasSuffix(s, "/") {
					rsp += "/"
				}

				return rsp
			})
		}

		messages[route.Endpoint] = model.Message{
			Request: model.Request{
				Method: strings.ToUpper(route.Method),
				URL:    url,
			},
			Response: responses,
		}
	}

	return messages, nil
}
