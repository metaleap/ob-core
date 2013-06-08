package obsrv

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"

	ob "github.com/openbase/ob-core"
)

var (
	//	The http.Server used to serve web requests
	Http http.Server

	//	Multi-plexing request router
	Router *mux.Router

	//	Custom event handlers
	On struct {
		//	Request-related event handlers
		Request struct {
			//	Event handlers to be invoked before serving a web request (except static files)
			Serving RequestContextEventHandlers

			//	Event handlers to be invoked immediately after serving a web request (except static files)
			Served RequestContextEventHandlers
		}
	}
)

//	Initializes the package for serving web requests
func Init() {
	staticDist := http.FileServer(http.Dir(ob.Hive.Path("dist", "client", "pub")))
	staticCust := http.FileServer(http.Dir(ob.Hive.Path("cust", "pub")))
	custRoot := ob.Hive.Path("cust", "pub", "root")
	staticCustRoot := http.StripPrefix("/", http.FileServer(http.Dir(custRoot)))
	Router = mux.NewRouter()
	Router.PathPrefix("/_dist/").Handler(http.StripPrefix("/_dist/", staticDist))
	Router.PathPrefix("/_cust/").Handler(http.StripPrefix("/_cust/", staticCust))
	ob.Hive.Watch.WatchFiles(custRoot, "*.*", true, func(filePath string) {
		Router.Path("/" + filepath.Base(filePath)).Handler(staticCustRoot)
	})
	Router.PathPrefix("/").HandlerFunc(serveRequest)
	Http.Handler = Router
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
