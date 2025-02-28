package main

import (
	"bytes"
	"errors"
	"log"
	"strings"
)

type HttpRequestParserBytes struct {
	Headers [10]string
}

type RequestFirstLine struct {
	HttpMethod  string
	URI         string
	HttpVersion string
}

func (firstLine RequestFirstLine) ToString() string {
	return "HttpMethod:" + firstLine.HttpMethod + ", URI:" + firstLine.URI + ", HttpVersion: " + firstLine.HttpVersion
}

func (request HttpRequestParserBytes) ParseHttpRequest(requestBytes []byte) HttpRequest {
	_, err := request.parseHttpRequestFirstLine(requestBytes)
	if err != nil {
		log.Println("error in reading first line of request: ", err)
	}

	// log.Println("correctly parse first line of request:" + firstLine.ToString())
	return HttpRequest{}
}

func (request HttpRequestParserBytes) parseHttpRequestFirstLine(requestBytes []byte) (RequestFirstLine, error) {
	// Split the first line from the request by \r\n (CRLF)
	// HTTP request lines are typically ended with \r\n (carriage return and newline)
	dummyResult := RequestFirstLine{"", "", ""}
	lines := bytes.SplitN(requestBytes, []byte("\r\n"), 2)
	if len(lines) < 1 {
		log.Println("invalid HTTP request")
		return dummyResult, errors.New("invalid HTTP request")
	}

	// Split the first line by spaces (to get method, URL, and version)
	parts := strings.Fields(string(lines[0]))
	if len(parts) != 3 {
		log.Println("invalid HTTP request line")
		return dummyResult, errors.New("invalid HTTP request line")
	}

	return RequestFirstLine{parts[0], parts[1], parts[2]}, nil
}
