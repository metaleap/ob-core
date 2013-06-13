package obwebui

//	WebUi.Libs[..].CssUrl

import (
	obpkg "github.com/openbase/ob-core/pkg"
	obpkg_webuilib "github.com/openbase/ob-core/pkg/webuilib"
	obpkg_webuiskin "github.com/openbase/ob-core/pkg/webuiskin"
)

type PageContextWebUi struct {
	Libs         []*obpkg_webuilib.PkgCfg
	SkinTemplate *PageTemplate
}

type PageContext struct {
	WebUi PageContextWebUi
}

func (me *PageContext) Init() {
	var cfg *obpkg_webuilib.PkgCfg
	skinPkg := obpkg.Reg.ByName("webuiskin", "fluid")
	me.WebUi.SkinTemplate = getPageTemplate(skinPkg.Cfg.(*obpkg_webuiskin.PkgCfg).SubRelTemplateDirPath)
	for _, pkg := range obpkg.Reg.ByKind("webuilib", skinPkg.Info.Require) {
		if cfg, _ = pkg.Cfg.(*obpkg_webuilib.PkgCfg); cfg != nil {
			me.WebUi.Libs = append(me.WebUi.Libs, cfg)
		}
	}
}
