package boundary

import (
	"encoding/json"
	"fmt"
	"github.com/webdevwilson/serverless-multicloud/src/controller"
	"log"
	"net/http"
)

type httpHandler struct {
	pathPrefix  string
}

func NewRequestHandler(pathPrefix string) http.Handler {
	return &httpHandler{
		pathPrefix,
	}
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// strip the path prefix
	uri := req.URL.Path[len(h.pathPrefix):]

	var data interface{}
	var err error
	if uri == "/" {
		http.Redirect(w, req, fmt.Sprintf("%s/hello", h.pathPrefix), 302)
		return
	} else if uri == "/hello" {
		data, err = controller.SayHello()
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		internalError(w, req, err)
		return
	}

	success(w, req, data)
}

func internalError(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("[ERROR] %s", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func success(w http.ResponseWriter, req *http.Request, data interface{}) {
	var body []byte
	var err error
	s, ok := data.(string)
	if ok {
		// data is a string, write it to the request body
		body = []byte(s)
	} else {
		// data is structured, marshal it to JSON
		body, err = json.Marshal(data)
		if err != nil {
			internalError(w, req, err)
			return
		}
	}

	_, err = w.Write(body)
	if err != nil {
		internalError(w, req, err)
		return
	}
}
