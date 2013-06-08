package obsrv

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	ob "github.com/openbase/ob-core"
)

var (
	//	The http.Server used to serve web requests
	Http http.Server

	//	The http.Handler that maps to the [hive]/client/pub directory
	StaticFiles http.Handler

	//	Request-routing multiplexer
	Mux *mux.Router

	//	Custom event handlers
	On struct {
		//	Request-related event handlers
		Request struct {
			//	Event handlers to be invoked before serving a web request (except _static files)
			Serving RequestContextEventHandlers

			//	Event handlers to be invoked immediately after serving a web request (except _static files)
			Served RequestContextEventHandlers
		}
	}
)

//	Initializes the package for serving web requests
func Init() {
	StaticFiles = http.FileServer(http.Dir(ob.Hive.Path("client", "pub")))
	Mux = mux.NewRouter()
	Mux.PathPrefix("/_static/").Handler(http.StripPrefix("/_static/", StaticFiles))
	Mux.PathPrefix("/").HandlerFunc(serveRequest)
	Http.Handler = Mux
	Http.ReadTimeout = 2 * time.Minute
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
