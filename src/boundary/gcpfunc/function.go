package gcpfunc

import (
	"github.com/webdevwilson/serverless-multicloud/src/boundary"
	"net/http"
)

var handler http.HandlerFunc = boundary.NewRequestHandler("/execute").ServeHTTP

func Handler(w http.ResponseWriter, req *http.Request) {
	handler(w, req)
}
