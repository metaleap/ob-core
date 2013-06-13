package obpkg_webuiskin

import (
	"fmt"
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

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

type PkgCfg struct {
	SubRelTemplateDirPath string
	pkg                   *obpkg.Package
}

func newPkgCfg() (me *PkgCfg) {
	me = &PkgCfg{}
	return
}
