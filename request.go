package httpjson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Bind(response http.ResponseWriter, request *http.Request, model interface{}) bool {
	if !strings.Contains(request.Header.Get("content-type"), "json") {
		http.Error(response,
			fmt.Sprintf("%s (json content-type required)", http.StatusText(http.StatusUnsupportedMediaType)),
			http.StatusUnsupportedMediaType,
		)
		return false
	}

	err := json.NewDecoder(request.Body).Decode(model)
	if err != nil {
		http.Error(response,
			fmt.Sprintf("%s (json decode failure: [%s])", http.StatusText(http.StatusBadRequest), err),
			http.StatusBadRequest,
		)
		return false
	}

	return true
}
