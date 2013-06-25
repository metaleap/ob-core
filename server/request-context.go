package obsrv

import (
	"bytes"
	"net/http"

	webctx "github.com/gorilla/context"

	ob "github.com/openbase/ob-core"
)

//	A function that accepts a `*RequestContext`.
type RequestContextHandler func(*RequestContext)

//	A collection of `RequestContextHandler` functions.
type RequestContextHandlers []RequestContextHandler

//	Appends all the specified `handlers` to `me`.
func (me *RequestContextHandlers) Add(handlers ...RequestContextHandler) {
	*me = append(*me, handlers...)
}

//	Provides context for a non-static web request.
type RequestContext struct {
	Ctx *Ctx

	//	Context related to the current `Page`, if any.
	Page *PageContext

	//	Defaults to `Ctx.Log`.
	Log ob.Logger

	//	The `http.Request` for this `RequestContext`.
	Req *http.Request

	out http.ResponseWriter
}

func newRequestContext(ctx *Ctx, o http.ResponseWriter, r *http.Request) (me *RequestContext) {
	me = &RequestContext{Ctx: ctx, out: o, Req: r, Log: ctx.Log}
	me.Page = newPageContext(ctx)
	return
}

//	http://www.gorillatoolkit.org/pkg/context#Get
func (me *RequestContext) Get(key interface{}) interface{} {
	return webctx.Get(me.Req, key)
}

func (me *RequestContext) serveRequest() {
	var w bytes.Buffer
	err := me.Page.WebUI.Skin.exec(&w, me)
	if err == nil {
		me.out.Write(w.Bytes())
	} else {
		me.out.Write([]byte(err.Error()))
	}
}

//	http://www.gorillatoolkit.org/pkg/context#Set
func (me *RequestContext) Set(key, val interface{}) {
	webctx.Set(me.Req, key, val)
}
