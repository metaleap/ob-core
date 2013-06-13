package obpkg_webuilib

import (
	"fmt"
	"path/filepath"

	usl "github.com/metaleap/go-util/slice"
	usort "github.com/metaleap/go-util/slice/sort"

	ob "github.com/openbase/ob-core"
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
	ob.Hive.WalkDirsIn(func(dirPath string) bool {
		usl.StrAppendUnique(&cfg.Versions, filepath.Base(dirPath))
		return true
	}, "pkg", pkg.NameFull)
	cfg.Versions = usort.StrSortDesc(cfg.Versions)
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

//	Represents the Package.Cfg of a Package of Kind "webuilib"
type PkgCfg struct {
	//	Extension-less, path-less CSS file names, from the "css" setting
	Css []string

	//	Extension-less, path-less JS file names, from the "js" setting
	Js []string

	//	All versions found as separate folders in the package directory,
	//	sorted descending-alphabetically ("from newest to oldest")
	Versions []string

	pkg *obpkg.Package
}

func newPkgCfg() (me *PkgCfg) {
	me = &PkgCfg{}
	return
}

//	Returns a server-relative URL for the specified webuilib file.
//	For example, for a PkgCfg with name "bootstrap2", Url("bootstrap-responsive","css") returns:
//	/_pkg/webuilib-bootstrap2/{Versions[0]}/css/bootstrap-responsive.css
func (me *PkgCfg) Url(fileBaseName, dotLessExtExt string) string {
	return strf("/_pkg/%s/%s/%s/%s.%s", me.pkg.NameFull, me.Versions[0], dotLessExtExt, fileBaseName, dotLessExtExt)
}
