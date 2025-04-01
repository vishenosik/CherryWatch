// TODO: move to web-tools
package httpjson

import (
	"encoding/json"
	"net/http"
)

func Decode[Type any](r *http.Request) (Type, error) {
	var elem Type
	if err := json.NewDecoder(r.Body).Decode(&elem); err != nil {
		return elem, err
	}
	defer r.Body.Close()
	return elem, nil
}
