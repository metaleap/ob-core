package obpkg_webuiskin

import (
	"path/filepath"

	ob "github.com/openbase/ob-core"
)

func init() {
	ob.BundleCfgReloaders["webuiskin"] = reloadBundleCfg
}

func reloadBundleCfg(bundle *ob.Bundle) {
	var cfg *BundleCfg
	if cfg, _ = bundle.Cfg.(*BundleCfg); cfg == nil {
		cfg = newBundleCfg()
		bundle.Cfg = cfg
	}
	cfg.bundle = bundle
	cfg.SubRelTemplateDirPath = filepath.Join("pkg", bundle.NameFull, "template")
}

//	Represents the Bundle.Cfg of a Bundle of Kind "webuiskin"
type BundleCfg struct {
	//	A HiveSub-relative directory path for the template files of this webuiskin.
	//	For example, for a BundleCfg with Name "fluid", this would be:
	//	pkg/webuiskin-fluid/template
	SubRelTemplateDirPath string

	bundle *ob.Bundle
}

func newBundleCfg() (me *BundleCfg) {
	me = &BundleCfg{}
	return
}
