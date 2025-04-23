package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
+---------+--------------------------------------------------------------------------+-----------------------------+
| Method  | Description                                                              | Payload in General          |
+---------+--------------------------------------------------------------------------+-----------------------------+
| GET     | Requests a representation of the specified resource.                     | Optional (URL parameters)   |
|         | Data is typically passed through URL parameters.                         |                             |
+---------+--------------------------------------------------------------------------+-----------------------------+
| POST    | Submits data to be processed to a specified resource.                    | Yes                         |
|         | Often used to create or update resources. The data is sent in body.      |                             |
+---------+--------------------------------------------------------------------------+-----------------------------+
| PUT     | Replaces the current representation of the target resource.              | Yes                         |
|         | Used for updating or creating resources.                                 |                             |
+---------+--------------------------------------------------------------------------+-----------------------------+
| DELETE  | Deletes the specified resource.                                          | Optional                    |
+---------+--------------------------------------------------------------------------+-----------------------------+
| HEAD    | Like GET, but without the response body. Used to retrieve metadata.      | Optional                    |
+---------+--------------------------------------------------------------------------+-----------------------------+
| OPTIONS | Describes communication options for the target resource.                 | Optional                    |
|         | Allows the client to determine supported HTTP methods.                   |                             |
+---------+--------------------------------------------------------------------------+-----------------------------+
| PATCH   | Applies partial modifications to a resource.                             | Yes                         |
+---------+--------------------------------------------------------------------------+-----------------------------+
| CONNECT | Establishes a tunnel to the server, often for secure HTTPS connections.  | Optional                    |
+---------+--------------------------------------------------------------------------+-----------------------------+
| TRACE   | Performs a message loop-back test along the request path.                | No                          |
|         | Useful for debugging.                                                    |                             |
+---------+--------------------------------------------------------------------------+-----------------------------+
| PRI     | Initiates a prioritized request (HTTP/2).                                | Optional                    |
|         | Not widely supported.                                                    |                             |
+---------+--------------------------------------------------------------------------+-----------------------------+


+-------------------------------+---------------------------------------------------------------+
| Content-Type                  | Description                                                   |
+-------------------------------+---------------------------------------------------------------+
| application/json              | 📄 JSON needs to be deserialized into objects or data         |
|                               | structures.                                                   |
+-------------------------------+---------------------------------------------------------------+                                              |
+-------------------------------+---------------------------------------------------------------+
| multipart/form-data           | 📎 Used for forms with file uploads. Requires boundary        |
|                               | parsing and file handling.                                   |
+-------------------------------+---------------------------------------------------------------+
| application/x-www-form-urlencoded | 📝 Standard HTML forms. Must be decoded into key-value pairs. |
+-------------------------------+---------------------------------------------------------------+

*/

/*
Valid methods support by the server
*/
const (
	GetMethod     = "GET"
	PostMethod    = "POST"
	PutMethod     = "PUT"
	DeleteMethod  = "DELETE"
	HeadMethod    = "HEAD"
	OptionsMethod = "OPTIONS"
	PatchMethod   = "PATCH"
	ConnectMethod = "CONNECT"
	TraceMethod   = "TRACE"
)

var validHttpMethods = [9]string{
	GetMethod,
	PostMethod,
	PutMethod,
	DeleteMethod,
	HeadMethod,
	OptionsMethod,
	PatchMethod,
	ConnectMethod,
	TraceMethod,
}

/*
	Valid protocol versions
*/

const (
	HTTP1     = "HTTP/1.0"
	HTTP1DOT1 = "HTTP/1.1"
	HTTP2     = "HTTP/2"
	HTTP3     = "HTTP/3"
)

var validHttpVersions = [4]string{
	HTTP1,
	HTTP1DOT1,
	HTTP2,
	HTTP3,
}

var httpMethodsWithoutBody = [6]string{
	GetMethod,
	DeleteMethod,
	HeadMethod,
	OptionsMethod,
	ConnectMethod,
	TraceMethod,
}

/*
Content types that are auto parsed by server
*/
const (
	applicationJson    = "application/json"
	multiPartFormData  = "multipart/form-data"
	formDataUrlEncoded = "application/x-www-form-urlencoded"
)

var contentTypesToParse = [3]string{
	applicationJson,
	multiPartFormData,
	formDataUrlEncoded,
}

const (
	ContentTypeHeader = "content-type"
	HttpMethod        = "Method"
)

type HttpRequest struct {
	Headers                   map[string]string
	ContentType               string
	CustomerHeaders           map[string]string
	IsHttps                   bool
	Content                   []byte
	Protocol                  string
	Method                    string
	Location                  string
	FormContentBody           map[string]string      // only gets populated if the content-type is multipart/form-data  or application/x-www-form-urlencoded
	JsonContentBody           map[string]interface{} // only gets populated if the content-type is application/json
	FormUrlEncodedContentBody map[string][]string    // only get populated if the content-type is application/x-www-form-urlencoded
}

func (request HttpRequest) ParseHttpRequest(requestBytes []byte) HttpRequest {
	parser := HttpRequestParserString{}
	return parser.ParseTillRequestHeaders(requestBytes)
}

func (request HttpRequest) ValidateRequestTillHeaders() error {

	isValidMethod := false

	for _, methodName := range validHttpMethods {
		if request.Method == methodName {
			isValidMethod = true
		}
	}

	if !isValidMethod {
		logrus.Debug("Invalid http method")
		return errors.New("Invalid http method")
	}

	isValidHttpVersion := false

	for _, httpVersion := range validHttpVersions {
		if request.Protocol == httpVersion {
			isValidHttpVersion = true
		}
	}

	if !isValidHttpVersion {
		logrus.Debug("Invalid http protocol")
		return errors.New("Invalid http protocol")
	}

	return nil
}

func (request HttpRequest) IsRequestAllowedToHaveBody() bool {
	for _, methodName := range httpMethodsWithoutBody {
		if methodName == request.Method {
			logrus.Debug("Request not supposed to contain body, so skip reading content-length and parsing body")
			return false
		}
	}
	return true
}

func (request HttpRequest) PrintContent() {
	switch request.ContentType {
	case applicationJson:
		jsonBytes, _ := json.MarshalIndent(request.JsonContentBody, "", "  ")
		logrus.Debug("content: ", string(jsonBytes))
	case multiPartFormData:
		result := ""
		for k, v := range request.FormContentBody {
			result += k + "=" + v + "\n"
		}
		logrus.Debug("content: ", result)

	case formDataUrlEncoded:
		result := ""
		for k, values := range request.FormUrlEncodedContentBody {
			result += k + "=" + strings.Join(values, ",") + "\n"
		}

		logrus.Debug("content: ", result)
	default:
		logrus.Debug("content: ", request.Content)
	}
}
