package httpjson

import "errors"

type Error struct {
	StatusCodeValue int `json:"-"`

	ID      string   `json:"id"`
	Message string   `json:"message,omitempty"`
	Fields  []string `json:"fields,omitempty"`
}

func (this Error) StatusCode() int {
	return this.StatusCodeValue
}
func (this Error) Error() string {
	panic(notImplemented)
}

var notImplemented = errors.New("return value not necessary, just implementing a marker interface")
