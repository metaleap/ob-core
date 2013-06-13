package obpkg_webuiskin

import (
	"path/filepath"

	obpkg "github.com/openbase/ob-core/pkg"
)

func init() {
	obpkg.PkgCfgLoaders["webuiskin"] = reloadPkgCfg
}

func reloadPkgCfg(pkg *obpkg.Package) {
	var cfg *PkgCfg
	if cfg, _ = pkg.Cfg.(*PkgCfg); cfg == nil {
		cfg = newPkgCfg()
		pkg.Cfg = cfg
	}
	cfg.pkg = pkg
	cfg.SubRelTemplateDirPath = filepath.Join("pkg", pkg.NameFull, "template")
}

//	Represents the Package.Cfg of a Package of Kind "webuiskin"
type PkgCfg struct {
	//	A HiveSub-relative directory path for the template files of this webuiskin.
	//	For example, for a PkgCfg with name "fluid", this would be:
	//	pkg/webuiskin-fluid/template
	SubRelTemplateDirPath string

	pkg *obpkg.Package
}

func newPkgCfg() (me *PkgCfg) {
	me = &PkgCfg{}
	return
}
