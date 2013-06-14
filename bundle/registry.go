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
	//	The global bundle package registry
	Reg Registry
)

func init() {
	Reg.Init()
}

//	Bundle package registry, used for Reg
type Registry struct {
	initialBulkLoading bool
	mutex              sync.Mutex
	allBundles         map[string]*Bundle
	cachedByKind       map[string]Bundles
	watched            map[string]bool
}

func (me *Registry) fileName(name, kind string) string {
	return strf("%s.%s.ob-pkg", name, kind)
}

//	Initializes a few internal hash-maps to non-nil, NO need to call this for Reg.
//	Does not load the bundles just yet, this is done via any of the ByXYZ() methods.
func (me *Registry) Init() {
	Reg.cachedByKind, Reg.watched = map[string]Bundles{}, map[string]bool{}
}

func (me *Registry) reloadBundle(bundleDirPath string) {
	if !me.initialBulkLoading {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}
	addsOrDels, dirName := false, filepath.Base(bundleDirPath)
	for key, bundle := range me.allBundles {
		if key != bundle.NameFull || !ob.Hive.SubFileExists("pkg", bundle.NameFull, me.fileName(bundle.Name, bundle.Kind)) {
			ob.Log.Warningf("[BUNDLE] Removing '%s': bundle directory or file no longer exists or renamed", key)
			addsOrDels = true
			delete(me.allBundles, key)
		}
	}
	if pos := strings.Index(dirName, "-"); pos > 0 && pos == strings.LastIndex(dirName, "-") {
		if !me.watched[bundleDirPath] {
			addsOrDels = true
			me.watched[bundleDirPath] = true
			ob.Hive.WatchDualDir(func(dp string) { me.reloadBundle(filepath.Dir(dp)) }, false, "pkg", dirName)
		}
		kind, name := dirName[:pos], dirName[pos+1:]
		if cfgFilePath := filepath.Join(bundleDirPath, me.fileName(name, kind)); uio.FileExists(cfgFilePath) {
			ob.Log.Infof("[BUNDLE] Loading '%s' from '%s'", dirName, cfgFilePath)
			bundle := me.allBundles[dirName]
			if bundle == nil {
				addsOrDels, bundle = true, newBundle()
				me.allBundles[dirName] = bundle
			}
			bundle.reload(kind, name, dirName, cfgFilePath)
		}
	} else {
		ob.Log.Warningf("[BUNDLE] Skipping '%s': expected directory name format '{kind}-{name}'", bundleDirPath)
	}
	if addsOrDels && !me.initialBulkLoading {
		me.refreshCachesAndMeta()
	}
}

func (me *Registry) refreshCachesAndMeta() {
	allKinds := map[string]bool{}
	for _, bundle := range me.allBundles {
		allKinds[bundle.Kind] = true
	}
	me.cachedByKind = map[string]Bundles{}
	for kind, _ := range allKinds {
		bundles := Bundles{}
		for _, bundle := range me.allBundles {
			if bundle.Kind == kind {
				bundles = append(bundles, bundle)
				for _, req := range bundle.Info.Require {
					if me.allBundles[req] == nil {
						ob.Log.Errorf("[BUNDLE] Bad dependency: '%s' requires '%s', which was not found.", bundle.NameFull, req)
						bundle.Diag.BadDeps = append(bundle.Diag.BadDeps, req)
					}
				}
			}
		}
		sort.Sort(bundles)
		me.cachedByKind[kind] = bundles
	}
}

func (me *Registry) ensureLoaded() {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	if loadNow := (me.allBundles == nil); loadNow {
		me.allBundles = map[string]*Bundle{}
		me.initialBulkLoading = true
		ob.Hive.WatchDualDir(me.reloadBundle, true, "pkg")
		me.refreshCachesAndMeta()
		me.initialBulkLoading = false
	}
}

//	If deps is empty, returns all Bundles of the specified kind.
//	Otherwise, out of all Bundles specified in deps or directly
//	or indirectly required by them, returns those of the specified kind.
func (me *Registry) ByKind(kind string, deps []string) (all Bundles) {
	me.ensureLoaded()
	if len(deps) == 0 {
		all = me.cachedByKind[kind]
	} else {
		var (
			byDep  func(string)
			bundle *Bundle
		)
		m := map[string]bool{}
		byDep = func(dep string) {
			if bundle = me.allBundles[dep]; bundle != nil {
				m[bundle.NameFull] = true
				for _, d := range bundle.Info.Require {
					byDep(d)
				}
			}
		}
		for _, d := range deps {
			byDep(d)
		}
		all = make(Bundles, 0, len(m))
		for k, _ := range m {
			if bundle = me.allBundles[k]; bundle.Kind == kind {
				all = append(all, me.allBundles[k])
			}
		}
		sort.Sort(all)
	}
	return
}

//	Returns the *Bundle with the specified fully-qualified identifier.
//	kindAndName can be a single string such as "webuilib-jquery", or 2 strings for kind and name, such as "webuilib" and "jquery".
func (me *Registry) ByName(kindAndName ...string) *Bundle {
	me.ensureLoaded()
	return me.allBundles[strings.Join(kindAndName, "-")]
}
