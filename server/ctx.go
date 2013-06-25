package obsrv

import (
	"github.com/go-utils/ugo"

	ob "github.com/openbase/ob-core"
)

type Ctx struct {
	ob.Ctx

	pageTemplateCache struct {
		mx ugo.MutexIf
		m  map[string]*PageTemplate
	}
}

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
}
