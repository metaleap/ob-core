package obsrv

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	//	The http.Server used to serve web requests
	Http http.Server

	//	The http.Handler that maps to the [hive]/client/pub directory
	StaticFiles http.Handler

	//	Request-routing multiplexer
	Mux *mux.Router
)

func init() {
	Mux = mux.NewRouter()
	Mux.PathPrefix("/").HandlerFunc(serveDefault)
	Http.Handler = Mux
	Http.ReadTimeout = 2 * time.Minute
}

func serveDefault(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

//	Starts listening to and serving web requests.
//	Uses Transport Layer Security only if both tls arguments are specified.
func ListenAndServe(addr, tlsCertFile, tlsKeyFile string) (err error) {
	Http.Addr = addr
	if len(tlsCertFile) > 0 && len(tlsKeyFile) > 0 {
		return Http.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	}
	return Http.ListenAndServe()
}
