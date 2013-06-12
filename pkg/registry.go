package obpkg

import (
	"path/filepath"
	"strings"
	"sync"

	uio "github.com/metaleap/go-util/io"

	ob "github.com/openbase/ob-core"
)

var (
	Reg Registry
)

type Registry struct {
	sync.Mutex
	m map[string]*Package
}

func (me *Registry) reloadPackages(subDirPath string) {
	me.Lock()
	defer me.Unlock()
	uio.WalkDirsIn(filepath.Dir(subDirPath), func(pkgDirPath string) bool {
		dirName := filepath.Base(pkgDirPath)
		if pos := strings.Index(dirName, "-"); pos > 0 && pos == strings.LastIndex(dirName, "-") {
			name, kind := dirName[:pos], dirName[pos+1:]
			pkg := me.m[dirName]
			if pkg == nil {
				pkg = newPackage()
				me.m[dirName] = pkg
			}
			pkg.reload(kind, name, dirName, filepath.Join(pkgDirPath, strf("%s.%s.obpkg", name, kind)))
		} else {
			ob.Opt.Log.Warningf("[PKG] Skipping '%s': expected directory name format '{pkgkind}-{pkgname}'", pkgDirPath)
		}
		return true
	})

}

func (me *Registry) ensureLoaded() {
	me.Lock()
	defer me.Unlock()
	if me.m == nil {
		me.m = map[string]*Package{}
		//	GO_1_0
		ob.Hive.WatchDualDir(func(subDirPath string) { me.reloadPackages(subDirPath) }, "pkg")
	}
}

/*
	uio.WalkAllFiles(hiveSubPath, func(fullPath string) bool {
		//	is .obpkg file?
		if fileNameExt := filepath.Ext(fullPath); fileNameExt == ".obpkg" {
			//	parent dir should be pkg full name: '{kind}-{name}', eg. 'webuilib-jquery'
			dirName := filepath.Base(filepath.Dir(fullPath))
			if pos := strings.Index(dirName, "-"); pos > 0 && pos == strings.LastIndex(dirName, "-") {
				kind, name, fileName := dirName[:pos], dirName[pos+1:], filepath.Base(fullPath)
				//	file name should be '{name}.{kind}.obpkg', eg. jquery.webuilib.obpkg
				if fileName == strf("%s.%s%s", name, kind, fileNameExt) {
					fullName := kind + "-" + name
					pkg := me.m[fullName]
					if pkg == nil {
						pkg = newPackage(kind, name, fullName)
						me.m[fullName] = pkg
					}
					if loader := kinds[kind]; loader == nil {
						pkg
					}
				} else {
					me.log.Warningf("[PKG] Skipping '%s': expected file name format '{pkgname}.{pkgkind}.obpkg' consistent with parent directory.")
				}
			} else {
				me.log.Warningf("[PKG] Skipping '%s': expected directory name format '{pkgkind}-{pkgname}'")
			}
		}
		return true
	})
*/

func (me *Registry) AllOfKind(kind string) (pkgs Packages) {
	me.ensureLoaded()
	for _, pkg := range me.m {
		if pkg.Kind == kind {
			pkgs = append(pkgs, pkg)
		}
	}
	return
}

func (me *Registry) ByFullName(fullName string) *Package {
	me.ensureLoaded()
	return me.m[fullName]
}
