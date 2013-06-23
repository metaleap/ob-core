package obsrv

import (
	"bytes"
	"net/http"

	webctx "github.com/gorilla/context"

	ob "github.com/openbase/ob-core"
	obwebui "github.com/openbase/ob-core/webui"
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
	*ob.Ctx

	//	Context related to the current `Page`, if any.
	*obwebui.PageContext

	//	The `http.ResponseWriter` for this `RequestContext`.
	Out http.ResponseWriter

	//	The `http.Request` for this `RequestContext`.
	Req *http.Request

	//	Defaults to `ob.Log`.
	Log ob.Logger
}

func newRequestContext(ctx *ob.Ctx, httpResponse http.ResponseWriter, httpRequest *http.Request) (me *RequestContext) {
	me = &RequestContext{Ctx: ctx, Out: httpResponse, Req: httpRequest, Log: ctx.Log}
	me.PageContext = obwebui.NewPageContext(ctx)
	return
}

//	http://www.gorillatoolkit.org/pkg/context#Get
func (me *RequestContext) Get(key interface{}) interface{} {
	return webctx.Get(me.Req, key)
}

func (me *RequestContext) serveRequest() {
	var w bytes.Buffer
	err := me.PageContext.WebUi.SkinTemplate.Execute(&w, me)
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
