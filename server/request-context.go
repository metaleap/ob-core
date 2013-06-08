package obsrv

import (
	"net/http"

	ob "github.com/openbase/ob-core"
)

func serveRequest(w http.ResponseWriter, r *http.Request) {
	rc := newRequestContext(w, r)
	for _, on := range On.Request.Serving {
		on(rc)
	}
	w.Write([]byte("Hello World!"))
	for _, on := range On.Request.Served {
		on(rc)
	}
}

//	A function that accepts a *RequestContext
type RequestContextEventHandler func(*RequestContext)

//	A collection of RequestContextEventHandler function-values
type RequestContextEventHandlers []RequestContextEventHandler

//	Adds the specified eventHandlers to me
func (me *RequestContextEventHandlers) Add(eventHandlers ...RequestContextEventHandler) {
	*me = append(*me, eventHandlers...)
}

//	Encapsulates and provides context for a (non-static) web request
type RequestContext struct {
	//	The http.ResponseWriter for this RequestContext
	Out http.ResponseWriter

	//	The http.Request for this RequestContext
	Req *http.Request

	//	Defaults to Opt.Log
	Log ob.Logger

	//	Not used in stand-alone mode.
	//	May be used in sandboxed mode (eg. the GAE package uses it for the current appengine.Context)
	Ctx interface{}
}

func newRequestContext(httpResponse http.ResponseWriter, httpRequest *http.Request) (me *RequestContext) {
	me = &RequestContext{Out: httpResponse, Req: httpRequest, Log: ob.Opt.Log}
	return
}
