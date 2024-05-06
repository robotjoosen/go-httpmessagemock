package translator

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/robotjoosen/go-httpmessagemock/model"
)

type HttpTranslator struct{}

func NewHTTPTranslator() (t HttpTranslator) {
	return
}

func (t HttpTranslator) Parse(data []byte) (map[string]model.Message, error) {
	messages := make(map[string]model.Message)

	sections := strings.Split(string(data), "###") // split sections
	for _, section := range sections {
		section = strings.TrimSpace(section)
		label, contents, ok := strings.Cut(section, "\n")
		if !ok {
			continue
		}

		httpMessage := regexp.MustCompile(`(\w*) (\S*https?:+\S+\w*)|(HTTP)/[0-9.]* (\d*) [\w ]*`)
		lines := httpMessage.FindAllStringSubmatch(contents, -1)
		index := httpMessage.FindAllStringSubmatchIndex(contents, -1)
		message := model.Message{
			Response: make(map[int]model.Response),
		}

		for i, line := range lines {
			lineIdentifier := line[3]
			if lineIdentifier == "" {
				lineIdentifier = line[1]
			}

			// setup message request or responses
			switch lineIdentifier {
			case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
				http.MethodConnect, http.MethodOptions, http.MethodTrace, http.MethodHead:
				message.Request.Method = line[1]
				message.Request.URL = line[2]
				// message.Request.Body = body
				// message.Request.Headers = headers
			case "HTTP":
				statusCode, err := strconv.Atoi(line[4])
				if err != nil {
					continue
				}

				// find the start and end of the content
				startLine := index[i][1]
				endLine := len(contents)
				if len(index) > i+1 {
					endLine = index[i+1][0] - 1
				}

				// split headers and body
				headersRaw, body, _ := strings.Cut(contents[startLine:endLine], "\n\n")

				// find headers
				headers := make(map[string]string)
				headerRegex := regexp.MustCompile(`([\w\S]*): ?([\w\S]*)`)
				headerLines := headerRegex.FindAllStringSubmatch(headersRaw, -1)
				for _, header := range headerLines {
					headers[header[1]] = header[2]
				}

				message.Response[statusCode] = model.Response{
					Headers: headers,
					Body:    body,
				}
			}
		}

		messages[label] = message
	}

	return messages, nil
}
