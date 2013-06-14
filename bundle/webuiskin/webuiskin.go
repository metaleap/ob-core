package obpkg_webuiskin

import (
	"path/filepath"

	obpkg "github.com/openbase/ob-core/bundle"
)

func init() {
	obpkg.BundleCfgLoaders["webuiskin"] = reloadBundleCfg
}

func reloadBundleCfg(bundle *obpkg.Bundle) {
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
	//	For example, for a BundleCfg with name "fluid", this would be:
	//	pkg/webuiskin-fluid/template
	SubRelTemplateDirPath string

	bundle *obpkg.Bundle
}

func newBundleCfg() (me *BundleCfg) {
	me = &BundleCfg{}
	return
}
