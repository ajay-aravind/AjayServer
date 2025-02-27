package main

type HttpRequest struct {
	Headers         [10]string
	ContentType     string
	CustomerHeaders map[string]string
	IsHttps         bool
	Content         string
}

func (request HttpRequest) ParseHttpRequest(requestBytes []byte) HttpRequest {
	parser := HttpRequestParserBytes{}
	return parser.ParseHttpRequest(requestBytes)
}
