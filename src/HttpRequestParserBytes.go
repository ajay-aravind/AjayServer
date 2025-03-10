package main

import (
	"bytes"
	"errors"
	"log"
	"strings"
)

type HttpRequestParserBytes struct {
	RequestFirstLine RequestFirstLine
	Headers          [10]string
	RequestLines     [][]byte
	RequestBody      []byte
}

type RequestFirstLine struct {
	HttpMethod  string
	URI         string
	HttpVersion string
}

func (firstLine RequestFirstLine) ToString() string {
	return "HttpMethod:" + firstLine.HttpMethod + ", URI:" + firstLine.URI + ", HttpVersion: " + firstLine.HttpVersion
}

func (parser HttpRequestParserBytes) ParseHttpRequest(requestBytes []byte) HttpRequest {
	_, err := parser.parseHttpRequestFirstLine(requestBytes)
	if err != nil {
		log.Println("error in reading first line of request: ", err)
	}

	// log.Println("correctly parse first line of request:" + firstLine.ToString())
	return HttpRequest{}
}

func (parsor HttpRequestParserBytes) parseHttpRequestFirstLine(requestBytes []byte) (RequestFirstLine, error) {
	// Split the first line from the request by \r\n (CRLF)
	// HTTP request lines are typically ended with \r\n (carriage return and newline)
	defaultResult := RequestFirstLine{"", "", ""}
	parsor.RequestLines = bytes.SplitN(requestBytes, []byte("\r\n"), 2)
	if len(parsor.RequestLines) < 1 {
		log.Println("invalid HTTP request")
		return defaultResult, errors.New("invalid HTTP request")
	}

	// Split the first line by spaces (to get method, URL, and version)
	parts := strings.Fields(string(parsor.RequestLines[0]))
	if len(parts) != 3 {
		log.Println("invalid HTTP request line")
		return defaultResult, errors.New("invalid HTTP request line")
	}

	parsor.RequestFirstLine = RequestFirstLine{parts[0], parts[1], parts[2]}
	return RequestFirstLine{parts[0], parts[1], parts[2]}, nil
}

func (parsor HttpRequestParserBytes) parseHttpRequestHeaders() error {
	return errors.New("not implemented yet")
}

func (parsor HttpRequestParserBytes) parseHttpRequestBody() error {
	return errors.New("not implemented yet")
}
