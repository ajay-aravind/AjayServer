package main

type HttpRequest struct {
	Headers         [10]string
	ContentType     string
	CustomerHeaders map[string]string
	IsHttps         bool
}
