package obsrv

import (
	"bytes"
	"net/http"

	webctx "github.com/gorilla/context"

	ob "github.com/openbase/ob-core"
	obpkg "github.com/openbase/ob-core/pkg"
	obwebui "github.com/openbase/ob-core/webui"
)

var (
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

func serveRequest(w http.ResponseWriter, r *http.Request) {
	rc := newRequestContext(w, r)
	for _, on := range On.Request.Serving {
		on(rc)
	}
	rc.serveRequest()
	for _, on := range On.Request.Served {
		on(rc)
	}
}

//	Encapsulates and provides context for a (non-static) web request
type RequestContext struct {
	//	Context related to the current Page, if any.
	obwebui.PageContext

	//	The http.ResponseWriter for this RequestContext
	Out http.ResponseWriter

	//	The http.Request for this RequestContext
	Req *http.Request

	//	Defaults to ob.Log
	Log ob.Logger

	//	Not used in the default stand-alone implementation (cmd/ob-server).
	//	May be used in sandboxed mode (eg. the GAE package uses it for the current appengine.Context)
	Ctx interface{}

	pkg *obpkg.Registry
}

func newRequestContext(httpResponse http.ResponseWriter, httpRequest *http.Request) (me *RequestContext) {
	me = &RequestContext{Out: httpResponse, Req: httpRequest, Log: ob.Log}
	me.PageContext.Init()
	return
}

//	http://www.gorillatoolkit.org/pkg/context#Get
func (me *RequestContext) Get(key interface{}) interface{} {
	return webctx.Get(me.Req, key)
}

func (me *RequestContext) serveRequest() {
	var w bytes.Buffer
	err := me.WebUi.SkinTemplate.Execute(&w, me)
	if err == nil {
		me.Out.Write(w.Bytes())
	} else {
		me.Out.Write([]byte(err.Error()))
	}
}

//	http://www.gorillatoolkit.org/pkg/context#Set
func (me *RequestContext) Set(key, val interface{}) {
	webctx.Set(me.Req, key, val)
}

//	A function that accepts a *RequestContext
type RequestContextEventHandler func(*RequestContext)

//	A collection of RequestContextEventHandler function-values
type RequestContextEventHandlers []RequestContextEventHandler

//	Adds the specified eventHandlers to me
func (me *RequestContextEventHandlers) Add(eventHandlers ...RequestContextEventHandler) {
	*me = append(*me, eventHandlers...)
}
