package obpkg_webuilib

import (
	"fmt"
	"path/filepath"

	uio "github.com/metaleap/go-util/io"
	usl "github.com/metaleap/go-util/slice"
	usort "github.com/metaleap/go-util/slice/sort"

	obpkg "github.com/openbase/ob-core/pkg"
)

func init() {
	obpkg.PkgCfgLoaders["webuilib"] = reloadPkgCfg
}

func reloadPkgCfg(pkg *obpkg.Package) {
	var cfg *PkgCfg
	if cfg, _ = pkg.Cfg.(*PkgCfg); cfg == nil {
		cfg = newPkgCfg()
		pkg.Cfg = cfg
	}
	cfg.pkg = pkg
	css, _ := pkg.CfgRaw.Default["css"].([]interface{})
	cfg.Css = usl.StrConvert(css, true)
	js, _ := pkg.CfgRaw.Default["js"].([]interface{})
	cfg.Js = usl.StrConvert(js, true)
	cfg.Versions = []string{}
	uio.WalkDirsIn(pkg.Dir, func(dirPath string) bool {
		usl.StrAppendUnique(&cfg.Versions, filepath.Base(dirPath))
		return true
	})
	cfg.Versions = usort.StrSortDesc(cfg.Versions)
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

type PkgCfg struct {
	Css      []string
	Js       []string
	Versions []string

	pkg *obpkg.Package
}

func newPkgCfg() (me *PkgCfg) {
	me = &PkgCfg{}
	return
}

func (me *PkgCfg) Url(fileBaseName, dotLessExtExt string) string {
	return strf("/_pkg/%s/%s/%s/%s.%s", me.pkg.NameFull, me.Versions[0], dotLessExtExt, fileBaseName, dotLessExtExt)
}
