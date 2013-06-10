package obsrv

import (
	"net/http"
	"path/filepath"
	"time"

	uio "github.com/metaleap/go-util/io"

	webmux "github.com/gorilla/mux"

	ob "github.com/openbase/ob-core"
)

var (
	//	The http.Server used to serve web requests
	Http http.Server

	//	Multi-plexing request router
	Router *webmux.Router
)

type dualStaticHandler struct {
	distDir, custDir string
	distSrv, custSrv http.Handler
}

func newDualStaticHandler(distDir, custDir string) (me *dualStaticHandler) {
	me = &dualStaticHandler{distDir: distDir, custDir: custDir}
	me.distSrv, me.custSrv = http.FileServer(http.Dir(distDir)), http.FileServer(http.Dir(custDir))
	return
}

func (me *dualStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if uio.FileExists(filepath.Join(me.custDir, r.URL.Path)) {
		me.custSrv.ServeHTTP(w, r)
	} else {
		me.distSrv.ServeHTTP(w, r)
	}
}

//	Initializes the package for serving web requests
func Init() {
	Router = webmux.NewRouter()
	Router.PathPrefix("/_dist/").Handler(http.StripPrefix("/_dist/", http.FileServer(http.Dir(ob.Hive.Subs.Dist.Paths.ClientPub))))
	Router.PathPrefix("/_cust/").Handler(http.StripPrefix("/_cust/", http.FileServer(http.Dir(ob.Hive.Subs.Cust.Paths.ClientPub))))
	dual := newDualStaticHandler(ob.Hive.Subs.Dist.Paths.ClientPub, ob.Hive.Subs.Cust.Paths.ClientPub)
	Router.PathPrefix("/_static/").Handler(http.StripPrefix("/_static/", dual))
	Router.Path("/{name}.{ext}").Handler(http.StripPrefix("/", dual))
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
