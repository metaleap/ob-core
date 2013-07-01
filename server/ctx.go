package obsrv

import (
	"github.com/go-utils/ugo"

	ob "github.com/openbase/ob-core"
)

//	Server-side aware wrapper around `ob.Ctx`.
type Ctx struct {
	//	The underlying `Ctx` being wrapped.
	ob.Ctx

	Http struct {
		Handler HttpHandler
	}

	pageTemplateCache struct {
		mx ugo.MutexIf
		m  map[string]*PageTemplate
	}
}

//	Only valid method to create and initialize a new `*Ctx`.
func NewCtx(hiveDir string, logger ob.Logger) (me *Ctx, err error) {
	me = &Ctx{}
	if err = me.Ctx.Init(hiveDir, logger); err != nil {
		me.Dispose()
		me = nil
	} else {
		me.init()
	}
	return
}

func (me *Ctx) init() {
	me.pageTemplateCache.m = map[string]*PageTemplate{}
	me.Http.Handler.initRouter(me)
}
