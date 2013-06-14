package obwebui

//	WebUi.Libs[..].CssUrl

import (
	obpkg "github.com/openbase/ob-core/bundle"
	obpkg_webuilib "github.com/openbase/ob-core/bundle/webuilib"
	obpkg_webuiskin "github.com/openbase/ob-core/bundle/webuiskin"
)

type PageContextWebUi struct {
	Libs         []*obpkg_webuilib.BundleCfg
	SkinTemplate *PageTemplate
}

type PageContext struct {
	WebUi PageContextWebUi
}

func (me *PageContext) Init() {
	var cfg *obpkg_webuilib.BundleCfg
	skinBundle := obpkg.Reg.ByName("webuiskin", "fluid")
	me.WebUi.SkinTemplate = getPageTemplate(skinBundle.Cfg.(*obpkg_webuiskin.BundleCfg).SubRelTemplateDirPath)
	for _, bundle := range obpkg.Reg.ByKind("webuilib", skinBundle.Info.Require) {
		if cfg, _ = bundle.Cfg.(*obpkg_webuilib.BundleCfg); cfg != nil {
			me.WebUi.Libs = append(me.WebUi.Libs, cfg)
		}
	}
}
