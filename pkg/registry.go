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
	m            map[string]*Package
	cachedByKind map[string]Packages
	watched      map[string]bool
}

func (me *Registry) fileName(name, kind string) string {
	return strf("%s.%s.ob-pkg", name, kind)
}

func (me *Registry) reloadPackages(subDirPath string) {
	me.Lock()
	defer me.Unlock()
	pkgsDirPath := subDirPath
	if filepath.Base(pkgsDirPath) != "pkg" {
		pkgsDirPath = filepath.Dir(subDirPath)
	}
	allKinds, addsOrDels := map[string]bool{}, false
	for key, pkg := range me.m {
		if key != pkg.NameFull || !uio.FileExists(filepath.Join(pkgsDirPath, pkg.NameFull, me.fileName(pkg.Name, pkg.Kind))) {
			delete(me.m, key)
		}
	}
	uio.WalkDirsIn(pkgsDirPath, func(pkgDirPath string) bool {
		dirName := filepath.Base(pkgDirPath)
		if !me.watched[pkgDirPath] {
			me.watched[pkgDirPath] = true
			ob.Hive.WatchDualDir(func(dp string) { me.reloadPackages(filepath.Dir(dp)) }, false, "pkg", dirName)
		}
		if pos := strings.Index(dirName, "-"); pos > 0 && pos == strings.LastIndex(dirName, "-") {
			kind, name := dirName[:pos], dirName[pos+1:]
			if cfgFilePath := filepath.Join(pkgDirPath, me.fileName(name, kind)); uio.FileExists(cfgFilePath) {
				allKinds[kind] = true
				pkg := me.m[dirName]
				if pkg == nil {
					addsOrDels = true
					pkg = newPackage()
					me.m[dirName] = pkg
				}
				pkg.reload(kind, name, dirName, cfgFilePath)
			}
		} else {
			ob.Opt.Log.Warningf("[PKG] Skipping '%s': expected directory name format '{pkgkind}-{pkgname}'", pkgDirPath)
		}
		return true
	})
	if addsOrDels {
		me.refreshCachesAndPkgInfos(allKinds)
	}
}

func (me *Registry) refreshCachesAndPkgInfos(allKinds map[string]bool) {
	for kind, _ := range allKinds {
		pkgs := Packages{}
		for _, pkg := range me.m {
			if pkg.Kind == kind {
				pkgs = append(pkgs, pkg)
				for _, req := range pkg.Info.Require {
					if me.m[req] == nil {
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
	var load bool
	me.Lock()
	if load = me.m == nil; load {
		me.m = map[string]*Package{}
	}
	me.Unlock()
	if load {
		ob.Hive.WatchDualDir(me.reloadPackages, true, "pkg")
	}
}

func (me *Registry) AllOfKind(kind string) Packages {
	me.ensureLoaded()
	return me.cachedByKind[kind]
}
