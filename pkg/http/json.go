// TODO: move to web
package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	BlankRoute = "/"
)

func Decode[Type any](r *http.Request) (Type, error) {
	var elem Type
	if r == nil {
		return elem, errors.New("nil request")
	}
	if err := json.NewDecoder(r.Body).Decode(&elem); err != nil {
		return elem, err
	}
	if err := r.Body.Close(); err != nil {
		return elem, err
	}
	return elem, nil
}

func MethodFunc(prefix string) func(string) string {
	return func(method string) string {
		return fmt.Sprintf("/%s.%s", prefix, method)
	}
}
