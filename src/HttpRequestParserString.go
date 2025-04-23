package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type HttpRequestParserString struct {
}

func (parser HttpRequestParserString) ParseTillRequestHeaders(requestBytes []byte) HttpRequest {
	var err error
	var request HttpRequest = HttpRequest{}
	logrus.Debug("Parsing request till request headers")
	requestFirstLine, remainingRequestLines, err := parser.parseHttpRequestFirstLine(requestBytes)

	if err != nil {
		logrus.Warn("Bad formatted request", err)
	}
	logrus.Debug(requestFirstLine.ToString())
	request.Method = requestFirstLine.HttpMethod
	request.Location = requestFirstLine.URI
	request.Protocol = requestFirstLine.HttpVersion

	requestHeaders, remainingRequestLines, err := parser.parseHttpRequestHeaders(remainingRequestLines)

	if err != nil {
		logrus.Warn("Bad formatted request", err)
	}

	request.Headers = requestHeaders

	for k, v := range request.Headers {
		logrus.Debug("headerKey: ", k, ", headerValue: ", v)
	}

	request.ContentType = requestHeaders[ContentTypeHeader]
	// parser.parseHttpRequestBody(remainingRequestLines, request.Headers[ContentTypeHeader], request.Method)
	return request
}

func (parsor HttpRequestParserString) parseHttpRequestFirstLine(requestBytes []byte) (RequestFirstLine, [][]byte, error) {
	// Split request into lines
	lines := bytes.Split(requestBytes, []byte("\n"))
	if len(lines) == 0 {
		return RequestFirstLine{"", "", ""}, lines, errors.New("empty request")
	}

	// Trim any trailing carriage return (\r)
	firstLine := strings.TrimSpace(string(lines[0]))
	parts := strings.Split(firstLine, " ")

	// HTTP Request line must have exactly 3 parts: Method, Location, Protocol
	if len(parts) != 3 {
		return RequestFirstLine{"", "", ""}, lines, errors.New("malformed request line")
	}

	method := parts[0]
	location := parts[1]
	protocol := parts[2]

	// Basic validation
	if method == "" || location == "" || protocol == "" {
		return RequestFirstLine{"", "", ""}, lines, errors.New("invalid request line format")
	}

	return RequestFirstLine{method, location, protocol}, lines[1:], nil
}

func (parsor HttpRequestParserString) parseHttpRequestHeaders(requestLines [][]byte) (map[string]string, [][]byte, error) {
	if len(requestLines) == 0 {
		return nil, nil, errors.New("empty request lines")
	}

	headers := make(map[string]string)
	var remainingLines [][]byte

	for index, line := range requestLines {
		trimmedLine := strings.TrimSpace(string(line))

		// End of headers is signaled by an empty line
		if trimmedLine == "" {
			remainingLines = requestLines[index:]

			logrus.Debug("request lines: ", string(bytes.Join(requestLines, []byte("\n"))))
			logrus.Debug("remaining lines: ", string(bytes.Join(remainingLines, []byte("\n"))))
			break
		}

		// Split header into key and value
		parts := strings.SplitN(trimmedLine, ":", 2)
		if len(parts) != 2 {
			return nil, nil, errors.New("malformed header line: " + trimmedLine)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Validate header key and value
		if key == "" || value == "" {
			return nil, nil, errors.New("invalid header format: key or value is empty")
		}

		// Ensure no duplicate headers (case insensitive comparison)
		lowerKey := strings.ToLower(key)
		if _, exists := headers[lowerKey]; exists {
			return nil, nil, errors.New("duplicate header detected: " + key)
		}

		headers[lowerKey] = value
	}

	return headers, remainingLines, nil
}

func (parser HttpRequestParserString) parseHttpRequestBody(requestLines []byte, request *HttpRequest) error {
	if len(requestLines) == 0 {
		return errors.New("empty request body")
	}
	body := requestLines

	mediaType, params, err := mime.ParseMediaType(request.ContentType)
	if err != nil {
		logrus.Debug("Error parsing content type:", err)
		return err
	}

	logrus.Debug("HTTP Method:", request.Method)
	logrus.Debug("Content-Type:", mediaType)

	switch mediaType {
	case applicationJson:
		if jsonContentBody, err := parser.parseJSONBody(body); err != nil {
			return err
		} else {
			request.JsonContentBody = jsonContentBody
		}

	case multiPartFormData:
		if formContentBody, err := parser.parseMultipartBody(body, params); err != nil {
			return err
		} else {
			request.FormContentBody = formContentBody
		}

	case formDataUrlEncoded:
		if formUrlEncodedContentBody, err := parser.parseFormURLEncodedBody(body); err != nil {
			return err
		} else {
			request.FormUrlEncodedContentBody = formUrlEncodedContentBody
		}
	default:
		request.Content = requestLines
		logrus.Debug("The content type can't be auto parsed, just assinging byte[] to HttpRequest.Content")
	}

	logrus.Debug("Content is parsed without any issues")
	return nil
}

func (parser HttpRequestParserString) parseJSONBody(body []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	var err error
	if err := json.Unmarshal(body, &jsonData); err != nil {
		logrus.Debug("Invalid JSON:", err)
	} else {
		logrus.Debug("Parsed JSON:", jsonData)
	}
	return jsonData, err
}

// won't work for file uploads so far
func (parser HttpRequestParserString) parseMultipartBody(body []byte, params map[string]string) (map[string]string, error) {
	boundary := params["boundary"]
	formData := make(map[string]string)

	if boundary == "" {
		logrus.Debug("Missing boundary in multipart/form-data")
		return formData, nil
	}

	var err error
	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Debug("Error reading multipart part:", err)
			break
		}
		partData, _ := io.ReadAll(part)
		logrus.Trace("Part: %s = %s\n", part.FormName(), string(partData))
		formData[part.FormName()] = string(partData)
	}

	return formData, err
}

func (parser HttpRequestParserString) parseFormURLEncodedBody(body []byte) (map[string][]string, error) {
	formData, err := url.ParseQuery(string(body))
	formDataMap := make(map[string][]string)
	if err != nil {
		logrus.Debug("Invalid form data:", err)
	} else {
		logrus.Debug("Parsed Form Data:")
		for key, values := range formData {
			fmt.Printf("%s = %v\n", key, values)
			formDataMap[key] = values
		}
	}

	return formDataMap, err
}
