package spartan

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
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
	p := req.URL.Path
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = p[:len(p)-1]
	}

	dst := h.Routes[p]
	if dst == "" {
		log.Printf("unknown path: %s", p)
		dst = h.Routes["/"]
	}
	if dst == "" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	log.Printf("[%d] %s -> %s", http.StatusTemporaryRedirect, p, dst)
	http.Redirect(w, req, dst, http.StatusTemporaryRedirect)
}
