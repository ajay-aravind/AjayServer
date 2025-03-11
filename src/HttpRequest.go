package main

type HttpRequest struct {
	Headers         map[string]string
	ContentType     string
	CustomerHeaders map[string]string
	IsHttps         bool
	Content         string
	Protocol        string
	Method          string
	Location        string
}

func (request HttpRequest) ParseHttpRequest(requestBytes []byte) HttpRequest {
	parser := HttpRequestParserString{}
	return parser.ParseRequest(requestBytes)
}
