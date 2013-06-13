package obpkg

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"

	uio "github.com/metaleap/go-util/io"

	ob "github.com/openbase/ob-core"
)

var (
	Reg Registry
)

func init() {
	Reg.cachedByKind, Reg.watched = map[string]Packages{}, map[string]bool{}
}

type Registry struct {
	sync.Mutex
	initialBulkLoading bool
	allPackages        map[string]*Package
	cachedByKind       map[string]Packages
	watched            map[string]bool
}

func (me *Registry) fileName(name, kind string) string {
	return strf("%s.%s.ob-pkg", name, kind)
}

func (me *Registry) reloadPackage(pkgDirPath string) {
	if !me.initialBulkLoading {
		me.Lock()
		defer me.Unlock()
	}
	addsOrDels, dirName := false, filepath.Base(pkgDirPath)
	for key, pkg := range me.allPackages {
		if key != pkg.NameFull || !ob.Hive.FileExists("pkg", pkg.NameFull, me.fileName(pkg.Name, pkg.Kind)) {
			ob.Opt.Log.Warningf("[PKG] Removing '%s': package directory or file no longer exists or renamed", key)
			addsOrDels = true
			delete(me.allPackages, key)
		}
	}
	if pos := strings.Index(dirName, "-"); pos > 0 && pos == strings.LastIndex(dirName, "-") {
		if !me.watched[pkgDirPath] {
			addsOrDels = true
			me.watched[pkgDirPath] = true
			ob.Hive.WatchDualDir(func(dp string) { me.reloadPackage(filepath.Dir(dp)) }, false, "pkg", dirName)
		}
		kind, name := dirName[:pos], dirName[pos+1:]
		if cfgFilePath := filepath.Join(pkgDirPath, me.fileName(name, kind)); uio.FileExists(cfgFilePath) {
			ob.Opt.Log.Infof("[PKG] Loading '%s' from '%s'", dirName, cfgFilePath)
			pkg := me.allPackages[dirName]
			if pkg == nil {
				addsOrDels, pkg = true, newPackage()
				me.allPackages[dirName] = pkg
			}
			pkg.reload(kind, name, dirName, cfgFilePath)
		}
	} else {
		ob.Opt.Log.Warningf("[PKG] Skipping '%s': expected directory name format '{pkgkind}-{pkgname}'", pkgDirPath)
	}
	if addsOrDels && !me.initialBulkLoading {
		me.refreshCachesAndMeta()
	}
}

func (me *Registry) refreshCachesAndMeta() {
	allKinds := map[string]bool{}
	for _, pkg := range me.allPackages {
		allKinds[pkg.Kind] = true
	}
	me.cachedByKind = map[string]Packages{}
	for kind, _ := range allKinds {
		pkgs := Packages{}
		for _, pkg := range me.allPackages {
			if pkg.Kind == kind {
				pkgs = append(pkgs, pkg)
				for _, req := range pkg.Info.Require {
					if me.allPackages[req] == nil {
						ob.Opt.Log.Errorf("[PKG] Bad dependency: '%s' requires '%s', which was not found.", pkg.NameFull, req)
						pkg.Diag.BadDeps = append(pkg.Diag.BadDeps, req)
					}
				}
			}
		}
		sort.Sort(pkgs)
		me.cachedByKind[kind] = pkgs
	}
}

func (me *Registry) ensureLoaded() {
	me.Lock()
	defer me.Unlock()
	if loadNow := (me.allPackages == nil); loadNow {
		me.allPackages = map[string]*Package{}
		me.initialBulkLoading = true
		ob.Hive.WatchDualDir(me.reloadPackage, true, "pkg")
		me.refreshCachesAndMeta()
		me.initialBulkLoading = false
	}
}

func (me *Registry) ByKind(kind string, deps []string) (all Packages) {
	me.ensureLoaded()
	if len(deps) == 0 {
		all = me.cachedByKind[kind]
	} else {
		var byDep func(string)
		m := map[string]bool{}
		byDep = func(dep string) {
			if pkg := me.allPackages[dep]; pkg != nil {
				m[pkg.NameFull] = true
				for _, d := range pkg.Info.Require {
					byDep(d)
				}
			}
		}
		for _, d := range deps {
			byDep(d)
		}
		all = make(Packages, 0, len(m))
		for k, _ := range m {
			all = append(all, me.allPackages[k])
		}
	}
	return
}

func (me *Registry) ByName(kind, name string) *Package {
	me.ensureLoaded()
	return me.allPackages[strf("%s-%s", kind, name)]
}
