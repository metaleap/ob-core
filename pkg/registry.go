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
	//	The global package registry
	Reg Registry
)

func init() {
	Reg.Init()
}

//	Package registry, used for Reg
type Registry struct {
	initialBulkLoading bool
	mutex              sync.Mutex
	allPackages        map[string]*Package
	cachedByKind       map[string]Packages
	watched            map[string]bool
}

func (me *Registry) fileName(name, kind string) string {
	return strf("%s.%s.ob-pkg", name, kind)
}

//	Initializes a few internal hash-maps to non-nil, NO need to call this for Reg.
//	Does not load the packages just yet, this is done via any of the ByXYZ() methods.
func (me *Registry) Init() {
	Reg.cachedByKind, Reg.watched = map[string]Packages{}, map[string]bool{}
}

func (me *Registry) reloadPackage(pkgDirPath string) {
	if !me.initialBulkLoading {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}
	addsOrDels, dirName := false, filepath.Base(pkgDirPath)
	for key, pkg := range me.allPackages {
		if key != pkg.NameFull || !ob.Hive.SubFileExists("pkg", pkg.NameFull, me.fileName(pkg.Name, pkg.Kind)) {
			ob.Log.Warningf("[PKG] Removing '%s': package directory or file no longer exists or renamed", key)
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
			ob.Log.Infof("[PKG] Loading '%s' from '%s'", dirName, cfgFilePath)
			pkg := me.allPackages[dirName]
			if pkg == nil {
				addsOrDels, pkg = true, newPackage()
				me.allPackages[dirName] = pkg
			}
			pkg.reload(kind, name, dirName, cfgFilePath)
		}
	} else {
		ob.Log.Warningf("[PKG] Skipping '%s': expected directory name format '{pkgkind}-{pkgname}'", pkgDirPath)
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
						ob.Log.Errorf("[PKG] Bad dependency: '%s' requires '%s', which was not found.", pkg.NameFull, req)
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
	me.mutex.Lock()
	defer me.mutex.Unlock()
	if loadNow := (me.allPackages == nil); loadNow {
		me.allPackages = map[string]*Package{}
		me.initialBulkLoading = true
		ob.Hive.WatchDualDir(me.reloadPackage, true, "pkg")
		me.refreshCachesAndMeta()
		me.initialBulkLoading = false
	}
}

//	If deps is empty, returns all Packages of the specified kind.
//	Otherwise, out of all Packages specified in deps or directly
//	or indirectly required by them, returns those of the specified kind.
func (me *Registry) ByKind(kind string, deps []string) (all Packages) {
	me.ensureLoaded()
	if len(deps) == 0 {
		all = me.cachedByKind[kind]
	} else {
		var (
			byDep func(string)
			pkg   *Package
		)
		m := map[string]bool{}
		byDep = func(dep string) {
			if pkg = me.allPackages[dep]; pkg != nil {
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
			if pkg = me.allPackages[k]; pkg.Kind == kind {
				all = append(all, me.allPackages[k])
			}
		}
	}
	return
}

//	Returns the Package with the specified fully-qualified identifier.
//	kindAndName can be a single string such as "webuilib-jquery", or 2 strings for kind and name, such as "webuilib" and "jquery".
func (me *Registry) ByName(kindAndName ...string) *Package {
	me.ensureLoaded()
	return me.allPackages[strings.Join(kindAndName, "-")]
}
