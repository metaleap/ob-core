package obwebui

//	WebUi.Libs[..].CssUrl

import (
	obpkg "github.com/openbase/ob-core/pkg"
	obpkg_webuilib "github.com/openbase/ob-core/pkg/webuilib"
)

type PageContext struct {
	WebUi WebUi
}

func (me *PageContext) Init() {
	var cfg *obpkg_webuilib.PkgCfg
	for _, pkg := range obpkg.Reg.AllOfKind("webuilib") {
		if cfg, _ = pkg.Cfg.(*obpkg_webuilib.PkgCfg); cfg != nil {
			me.WebUi.Libs = append(me.WebUi.Libs, cfg)
		}
	}
}
