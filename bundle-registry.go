package obcore

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-utils/ufs"
	"github.com/go-utils/ugo"
)

//	Bundle package registry, accessed from `Ctx.Bundles`.
type BundleRegistry struct {
	Ctx *Ctx

	initialLoad  bool
	mx           ugo.MutexIf
	allBundles   map[string]*Bundle
	cachedByKind map[string]Bundles
	watched      map[string]bool
}

func (me *BundleRegistry) init(ctx *Ctx) {
	me.Ctx, me.cachedByKind, me.watched = ctx, map[string]Bundles{}, map[string]bool{}
}

func (me *BundleRegistry) ensureLoaded() {
	defer me.mx.UnlockIf(me.mx.Lock())
	if loadNow := (me.allBundles == nil); loadNow {
		me.allBundles = map[string]*Bundle{}
		me.initialLoad = true
		me.Ctx.Hive.Subs.WatchIn(me.reloadBundle, true, "pkg")
		me.refreshCachesAndMeta()
		me.initialLoad = false
	}
}

func (me *BundleRegistry) fileName(name, kind string) string {
	return strf("%s.%s.ob-pkg", name, kind)
}

func (me *BundleRegistry) refreshCachesAndMeta() {
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
						me.Ctx.Log.Errorf("[BUNDLE] Bad dependency: '%s' requires '%s', which was not found.", bundle.NameFull, req)
						bundle.Diag.BadDeps = append(bundle.Diag.BadDeps, req)
					}
				}
			}
		}
		sort.Sort(bundles)
		me.cachedByKind[kind] = bundles
	}
}

func (me *BundleRegistry) reloadBundle(bundleDirPath string) {
	defer me.mx.UnlockIf(me.mx.LockIf(!me.initialLoad))
	addsOrDels, dirName := false, filepath.Base(bundleDirPath)
	for key, bundle := range me.allBundles {
		if key != bundle.NameFull || !me.Ctx.Hive.Subs.FileExists("pkg", bundle.NameFull, me.fileName(bundle.Name, bundle.Kind)) {
			me.Ctx.Log.Warningf("[BUNDLE] Removing '%s': bundle directory or file no longer exists or renamed", key)
			addsOrDels = true
			delete(me.allBundles, key)
		}
	}
	if pos := strings.Index(dirName, "-"); pos > 0 && pos == strings.LastIndex(dirName, "-") {
		if !me.watched[bundleDirPath] {
			addsOrDels = true
			me.watched[bundleDirPath] = true
			me.Ctx.Hive.Subs.WatchIn(func(dp string) { me.reloadBundle(filepath.Dir(dp)) }, false, "pkg", dirName)
		}
		kind, name := dirName[:pos], dirName[pos+1:]
		if cfgFilePath := filepath.Join(bundleDirPath, me.fileName(name, kind)); ufs.FileExists(cfgFilePath) {
			me.Ctx.Log.Infof("[BUNDLE] Loading '%s' from '%s'", dirName, cfgFilePath)
			bundle := me.allBundles[dirName]
			if bundle == nil {
				addsOrDels, bundle = true, newBundle(me)
				me.allBundles[dirName] = bundle
			}
			bundle.reload(kind, name, dirName, cfgFilePath)
		}
	} else {
		me.Ctx.Log.Warningf("[BUNDLE] Skipping '%s': expected directory name format '{kind}-{name}'", bundleDirPath)
	}
	if addsOrDels && !me.initialLoad {
		me.refreshCachesAndMeta()
	}
}

//	If `deps` is empty, returns all `Bundles` of the specified `kind`.
//
//	Otherwise, out of all `Bundles` specified in `deps` or directly
//	or indirectly required by them, returns those of the specified `kind`.
func (me *BundleRegistry) ByKind(kind string, deps []string) (all Bundles) {
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

//	Returns the `*Bundle` with the specified fully-qualified identifier.
//	`kindAndName` can be a single string such as `"webuilib-jquery"`, or 2 strings for `kind` and `name`, such as `"webuilib", "jquery"`.
func (me *BundleRegistry) ByName(kindAndName ...string) *Bundle {
	me.ensureLoaded()
	return me.allBundles[strings.Join(kindAndName, "-")]
}
