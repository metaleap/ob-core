package obsrv

import (
	"net/http"
	"path/filepath"
	"strings"

	webmux "github.com/gorilla/mux"

	ob "github.com/openbase/ob-core"
)

var (
	//	Multi-plexing request router
	Router *webmux.Router

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

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	rc := newRequestContext(w, r)
	for _, on := range On.Request.Serving {
		on(rc)
	}
	rc.serveRequest()
	for _, on := range On.Request.Served {
		on(rc)
	}
}

type hiveSubsStaticHandler struct {
	distSrv, custSrv http.Handler
}

func newHiveSubsStaticHandler(distDir, custDir string) (me *hiveSubsStaticHandler) {
	me = &hiveSubsStaticHandler{
		distSrv: http.FileServer(http.Dir(distDir)),
		custSrv: http.FileServer(http.Dir(custDir)),
	}
	return
}

func (me *hiveSubsStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ext, dir := filepath.Ext(r.URL.Path), filepath.Base(filepath.Dir(r.URL.Path))
	if strings.HasPrefix(ext, ".ob-") || (strings.HasPrefix(dir, "__") && strings.HasSuffix(dir, "__")) {
		http.Error(w, "Forbidden", 403)
	} else if ob.Hive.Subs.Cust.FileExists(r.URL.Path) {
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
	dual := newHiveSubsStaticHandler(ob.Hive.Subs.Dist.Paths.ClientPub, ob.Hive.Subs.Cust.Paths.ClientPub)
	Router.PathPrefix("/_static/").Handler(http.StripPrefix("/_static/", dual))
	Router.Path("/{name}.{ext}").Handler(http.StripPrefix("/", dual))
	dual = newHiveSubsStaticHandler(ob.Hive.Subs.Dist.Paths.Pkg, ob.Hive.Subs.Cust.Paths.Pkg)
	Router.PathPrefix("/_pkg/").Handler(http.StripPrefix("/_pkg/", dual))
	Router.PathPrefix("/").HandlerFunc(defaultHandler)
}
