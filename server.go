package sparta

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

func NewServer(ctx context.Context, raw []byte) *http.Server {
	return &http.Server{
		Addr: ":8000",
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		Handler: newHandler(raw),
	}
}

type Handler struct {
	Routes map[string]string `json:"routes"`
}

func newHandler(raw []byte) http.Handler {
	var conf Handler
	if err := json.Unmarshal(raw, &conf); err != nil {
		panic(err)
	}

	return conf
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dst := h.Routes[req.URL.Path]
	if dst == "" {
		log.Printf("unknown path: %s", req.URL.Path)
		dst = h.Routes["/"]
	}
	if dst == "" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	log.Printf("[%d] %s -> %s", http.StatusTemporaryRedirect, req.URL.Path, dst)
	http.Redirect(w, req, dst, http.StatusTemporaryRedirect)
}
