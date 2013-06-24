package obwebui

import (
	ob "github.com/openbase/ob-core"
	obpkg_webuilib "github.com/openbase/ob-core/bundle/webuilib"
	obpkg_webuiskin "github.com/openbase/ob-core/bundle/webuiskin"
)

type WebUi struct {
	Libs         []*obpkg_webuilib.BundleCfg
	SkinTemplate *PageTemplate
}

type PageContext struct {
	ctx *ob.Ctx

	WebUi WebUi
}

func NewPageContext(ctx *ob.Ctx) (me *PageContext) {
	me = &PageContext{ctx: ctx}
	var cfg *obpkg_webuilib.BundleCfg
	reg := ctx.Bundles()
	skinBundle := reg.ByName("webuiskin", "fluid")
	me.WebUi.SkinTemplate = getPageTemplate(ctx, skinBundle.Cfg.(*obpkg_webuiskin.BundleCfg).SubRelTemplateDirPath)
	for _, bundle := range reg.ByKind("webuilib", skinBundle.Info.Require) {
		if cfg, _ = bundle.Cfg.(*obpkg_webuilib.BundleCfg); cfg != nil {
			me.WebUi.Libs = append(me.WebUi.Libs, cfg)
		}
	}
	return
}
