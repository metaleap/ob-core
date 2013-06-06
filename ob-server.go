package core

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ObServer struct {
	Http        http.Server
	StaticFiles http.Handler
	Mux         *mux.Router
}

func serveDefault(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (_ *ObServer) init() {
	Server.Mux = mux.NewRouter()
	Server.Mux.PathPrefix("/").HandlerFunc(serveDefault)
	Server.Http.Handler = Server.Mux
	Server.Http.ReadTimeout = 2 * time.Minute
	if Sandboxed {
		http.Handle("/", Server.Mux)
	}
}
