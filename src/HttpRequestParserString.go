package main

import (
	"bytes"
	"errors"
	"log"
	"strings"
)

type HttpRequestParserString struct {
}

func (parser HttpRequestParserString) ParseRequest(requestBytes []byte) HttpRequest {
	var request HttpRequest = HttpRequest{}
	requestFirstLine, remainingRequestLines, err := parser.parseHttpRequestFirstLine(requestBytes)

	if err != nil {
		log.Fatal("Bad formatted request", err)
	}
	log.Println(requestFirstLine.ToString())
	request.Method = requestFirstLine.HttpMethod
	request.Location = requestFirstLine.URI
	request.Protocol = requestFirstLine.HttpVersion

	requestHeaders, remainingRequestLines, err := parser.parseHttpRequestHeaders(remainingRequestLines)

	if err != nil {
		log.Fatal("Bad formatted request", err)
	}

	request.Headers = requestHeaders

	for k, v := range request.Headers {
		log.Println("headerKey", k, "headerValue", v)
	}

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

func (parser HttpRequestParserString) parseHttpRequestBody(requestLines [][]byte) (string, error) {
	if len(requestLines) == 0 {
		return "", errors.New("empty request body")
	}

	var bodyLines [][]byte

	for _, line := range requestLines {
		bodyLines = append(bodyLines, line)
	}

	if len(bodyLines) == 0 {
		return "", errors.New("request body is missing or empty")
	}

	return string(bytes.Join(bodyLines, []byte("\n"))), nil
}
