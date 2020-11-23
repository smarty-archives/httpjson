package httpjson

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func Render(response http.ResponseWriter, options ...RenderOption) {
	renderer := NewResponseRenderer()

	Options.RenderIndent("", "  ")(renderer)

	for _, option := range options {
		option(renderer)
	}

	renderer.Render(response)
}

type ResponseRenderer struct {
	Indent string `json:"-"`
	Prefix string `json:"-"`

	StatusCodeValue int         `json:"-"`
	Headers         http.Header `json:"-"`

	Errors []Error     `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func NewResponseRenderer() *ResponseRenderer {
	return &ResponseRenderer{Indent: "  "}
}

func (this *ResponseRenderer) IncludeError(errs ...Error) {
	for _, err := range errs {
		if err.StatusCode() >= http.StatusInternalServerError {
			this.Errors = append([]Error{err}, this.Errors...)
		} else {
			this.Errors = append(this.Errors, err)
		}
	}
}

func (this *ResponseRenderer) AddHeader(values []string, key string) {
	if this.Headers == nil {
		this.Headers = make(http.Header)
	}
	for _, value := range values {
		this.Headers.Add(key, value)
	}
}

func (this *ResponseRenderer) Render(response http.ResponseWriter) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	if this.Prefix != "" || this.Indent != "" {
		encoder.SetIndent(this.Prefix, this.Indent)
	}
	err := encoder.Encode(this)
	if err != nil {
		http.Error(response, "json encode error", http.StatusInternalServerError)
		return
	}

	header := response.Header()
	for key, values := range this.Headers {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	header.Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(this.statusCode())
	_, _ = io.Copy(response, buffer)
}

func (this *ResponseRenderer) statusCode() int {
	errorCode := this.errorStatusCode()
	if errorCode >= http.StatusInternalServerError {
		return errorCode
	}
	if this.StatusCodeValue > 0 { // Perhaps this value shouldn't precede non-500 error code status codes?
		return this.StatusCodeValue
	}
	if errorCode > 0 {
		return errorCode
	}
	return http.StatusOK
}

func (this *ResponseRenderer) errorStatusCode() int {
	for _, err := range this.Errors {
		code := err.StatusCode()
		if code > 0 {
			return code
		}
	}
	return 0
}
