package obpkg_webuilib

import (
	"fmt"
	"path/filepath"
	"sort"

	usl "github.com/go-utils/uslice"

	ob "github.com/openbase/ob-core"
)

func init() {
	ob.BundleCfgLoaders["webuilib"] = reloadBundleCfg
}

func reloadBundleCfg(bundle *ob.Bundle) {
	var cfg *BundleCfg
	if cfg, _ = bundle.Cfg.(*BundleCfg); cfg == nil {
		cfg = newBundleCfg()
		bundle.Cfg = cfg
	}
	cfg.bundle = bundle
	css, _ := bundle.CfgRaw.Default["css"].([]interface{})
	cfg.Css = usl.StrConvert(css, true)
	js, _ := bundle.CfgRaw.Default["js"].([]interface{})
	cfg.Js = usl.StrConvert(js, true)
	cfg.Versions = []string{}
	bundle.Ctx().Hive.Subs.WalkDirsIn(func(dirPath string) bool {
		usl.StrAppendUnique(&cfg.Versions, filepath.Base(dirPath))
		return true
	}, "pkg", bundle.NameFull)
	sort.Sort(sort.Reverse(sort.StringSlice(cfg.Versions)))
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

//	Represents the Bundle.Cfg of a Bundle of Kind "webuilib"
type BundleCfg struct {
	//	Extension-less, path-less CSS file names, from the "css" setting
	Css []string

	//	Extension-less, path-less JS file names, from the "js" setting
	Js []string

	//	All versions found as separate folders in the bundle directory,
	//	sorted descending-alphabetically ("from newest to oldest")
	Versions []string

	bundle *ob.Bundle
}

func newBundleCfg() (me *BundleCfg) {
	me = &BundleCfg{}
	return
}

//	Returns a server-relative URL for the specified webuilib file.
//	For example, for a BundleCfg with name "bootstrap2", Url("bootstrap-responsive","css") returns:
//	/_pkg/webuilib-bootstrap2/{Versions[0]}/css/bootstrap-responsive.css
func (me *BundleCfg) Url(fileBaseName, dotLessExtExt string) string {
	return strf("/_pkg/%s/%s/%s/%s.%s", me.bundle.NameFull, me.Versions[0], dotLessExtExt, fileBaseName, dotLessExtExt)
}
