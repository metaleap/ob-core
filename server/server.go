package obsrv

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	webmux "github.com/gorilla/mux"

	ob "github.com/openbase/ob-core"
)

//	Must be initialized via `NewHttpHandler`.
type HttpHandler struct {
	http.Handler
	*ob.Ctx

	//	Custom event handlers
	On struct {
		//	Request-related event handlers
		Request struct {
			//	Event handlers to be invoked before
			//	serving a web request (except static files).
			PreServe RequestContextHandlers

			//	Event handlers to be invoked immediately after
			//	serving a web request (except static files).
			PostServe RequestContextHandlers
		}
	}
}

//	Initializes a new `*HttpHandler` to host the specified `*ob.Ctx`.
func NewHttpHandler(ctx *ob.Ctx) (router *HttpHandler) {
	if ctx != nil {
		router = &HttpHandler{Ctx: ctx}
		mux := webmux.NewRouter()
		router.Handler = mux
		mux.PathPrefix("/_dist/").Handler(http.StripPrefix("/_dist/", http.FileServer(http.Dir(ctx.Hive.Subs.Dist.Paths.ClientPub))))
		mux.PathPrefix("/_cust/").Handler(http.StripPrefix("/_cust/", http.FileServer(http.Dir(ctx.Hive.Subs.Cust.Paths.ClientPub))))
		dual := newHiveSubsStaticHandler(ctx, ctx.Hive.Subs.Dist.Paths.ClientPub, ctx.Hive.Subs.Cust.Paths.ClientPub)
		mux.PathPrefix("/_static/").Handler(http.StripPrefix("/_static/", dual))
		mux.Path("/{name}.{ext}").Handler(http.StripPrefix("/", dual))
		dual = newHiveSubsStaticHandler(ctx, ctx.Hive.Subs.Dist.Paths.Pkg, ctx.Hive.Subs.Cust.Paths.Pkg)
		mux.PathPrefix("/_pkg/").Handler(http.StripPrefix("/_pkg/", dual))
		mux.PathPrefix("/").HandlerFunc(router.serveRequest)
	}
	return
}

func (me *HttpHandler) serveRequest(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	rc := newRequestContext(me.Ctx, w, r)
	for _, on := range me.On.Request.PreServe {
		on(rc)
	}
	rc.serveRequest()
	for _, on := range me.On.Request.PostServe {
		on(rc)
	}
	if false {
		me.Ctx.Log.Infof("Request served in %v", time.Now().Sub(now))
	}
}

type hiveSubsStaticHandler struct {
	ctx              *ob.Ctx
	distSrv, custSrv http.Handler
}

func newHiveSubsStaticHandler(ctx *ob.Ctx, distDir, custDir string) (me *hiveSubsStaticHandler) {
	me = &hiveSubsStaticHandler{
		ctx:     ctx,
		distSrv: http.FileServer(http.Dir(distDir)),
		custSrv: http.FileServer(http.Dir(custDir)),
	}
	return
}

func (me *hiveSubsStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ext, dir := filepath.Ext(r.URL.Path), filepath.Base(filepath.Dir(r.URL.Path))
	if strings.HasPrefix(ext, ".ob-") || (strings.HasPrefix(dir, "__") && strings.HasSuffix(dir, "__")) {
		http.Error(w, "Forbidden", 403)
	} else if me.ctx.Hive.Subs.Cust.FileExists(r.URL.Path) {
		me.custSrv.ServeHTTP(w, r)
	} else {
		me.distSrv.ServeHTTP(w, r)
	}
}
