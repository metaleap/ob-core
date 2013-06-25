package obsrv

import (
	obpkg_webuilib "github.com/openbase/ob-core/bundle/webuilib"
	obpkg_webuiskin "github.com/openbase/ob-core/bundle/webuiskin"
)

//	Created during server-side rendering of a `PageTemplate`.
type PageContext struct {
	WebUI struct {
		Libs []*obpkg_webuilib.BundleCfg
		Skin *PageTemplate
	}
}

func newPageContext(ctx *Ctx) (me *PageContext) {
	me = &PageContext{}
	reg := ctx.Bundles()
	skinBundle := reg.ByName("webuiskin", "fluid")
	me.WebUI.Skin = ctx.getPageTemplate(skinBundle.Cfg.(*obpkg_webuiskin.BundleCfg).SubRelTemplateDirPath)
	var cfg *obpkg_webuilib.BundleCfg
	for _, bundle := range reg.ByKind("webuilib", skinBundle.Info.Require) {
		if cfg, _ = bundle.Cfg.(*obpkg_webuilib.BundleCfg); cfg != nil {
			me.WebUI.Libs = append(me.WebUI.Libs, cfg)
		}
	}
	return
}
