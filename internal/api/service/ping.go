package service

import (
	"encoding/json"
	"net/http"

	http_pkg "github.com/vishenosik/CherryWatch/pkg/http"
	"github.com/vishenosik/CherryWatch/pkg/versions"
)

func (srv server) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiVersion, err := http_pkg.ApiVersionFromContext[versions.DoubleVersion](r.Context())
		if err != nil {
			http_pkg.SendErrors(w, http.StatusBadRequest, err.Error())
		}

		type response struct {
			ApiVersion string `json:"api_version"`
			Status     string `json:"status"`
		}

		w.Header().Set("Content-Type", "application/json")
		resps := response{ApiVersion: apiVersion.String()}

		switch apiVersion.String() {
		case "1.0":
			resps.Status = "ok buddy"
		case "1.1":
			resps.Status = "ok maam"
		default:
			http_pkg.SendErrors(w, http.StatusBadRequest, "version unsupported")
			return
		}

		if err := json.NewEncoder(w).Encode(resps); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

	}
}
