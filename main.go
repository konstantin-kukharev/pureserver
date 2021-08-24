package pureserver

import "github.com/konstantin-kukharev/pureserver/internal/http"

type PatHandler http.PatHandler
type HttpRequestInterface http.HttpRequestInterface
type ResponseWriter http.ResponseWriter
type PatternServeMuxInterface http.PatternServeMuxInterface
type PatternServeMux http.PatternServeMux
type Router http.RouterInterface

type Server interface {
	Serve() error
	SetPort(...int)
	SetUnixSocket(...string)
	SetLoops(int)
}

type HandlerFunc func(w ResponseWriter, r HttpRequestInterface)

func (a HandlerFunc) Handle(w http.ResponseWriter, r http.HttpRequestInterface) {
	a(w, r)
}

func NewHttp(mux PatternServeMuxInterface) Server {
	server := &http.Server{}
	server.SetHandler(mux)
	return server
}

// NewMux returns a new PatternServeMuxInterface
func NewMux() PatternServeMuxInterface {
	mux := http.PatternServeMux{Handlers: make(map[string][]*http.PatHandler)}

	return &mux
}
