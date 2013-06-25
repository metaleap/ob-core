package obcore

import (
	"github.com/go-utils/ufs"
)

//	Only used for `Hive.Subs`.
type HiveSubs struct {
	root *HiveRoot

	//	{hive}/dist/
	Dist HiveSub

	//	{hive}/cust/
	Cust HiveSub
}

func (me *HiveSubs) init(root *HiveRoot) {
	me.root = root
	me.Dist.init(root, "dist")
	me.Cust.init(root, "cust")
}

//	Returns whether `me.Dist` or `me.Cust` contains the specified directory.
func (me *HiveSubs) DirExists(subRelPath ...string) bool {
	return me.Dist.DirExists(subRelPath...) || me.Cust.DirExists(subRelPath...)
}

//	Returns whether `me.Dist` or `me.Cust` contains the specified file.
func (me *HiveSubs) FileExists(subRelPath ...string) bool {
	return me.Dist.FileExists(subRelPath...) || me.Cust.FileExists(subRelPath...)
}

//	Returns either `me.Cust.DirPath(subRelPath...)` or `me.Dist.DirPath(subRelPath...)`.
func (me *HiveSubs) DirPath(subRelPath ...string) (dirPath string) {
	if dirPath = me.Cust.DirPath(subRelPath...); len(dirPath) == 0 {
		dirPath = me.Dist.DirPath(subRelPath...)
	}
	return
}

//	Returns either `me.Cust.FilePath(subRelPath...)` or `me.Dist.FilePath(subRelPath...)`.
func (me *HiveSubs) FilePath(subRelPath ...string) (filePath string) {
	if filePath = me.Cust.FilePath(subRelPath...); len(filePath) == 0 {
		filePath = me.Dist.FilePath(subRelPath...)
	}
	return
}

//	Returns either `me.Cust.FilePath` or `me.Cust.DirPath` or `me.Dist.Path` for the specified `subRelPath...`.
func (me *HiveSubs) Path(subRelPath ...string) (path string) {
	if path = me.Cust.FilePath(subRelPath...); len(path) == 0 {
		if path = me.Cust.DirPath(subRelPath...); len(path) == 0 {
			path = me.Dist.Path(subRelPath...)
		}
	}
	return
}

//	`ufs.WalkAllDirs` for `me.Dist` and `me.Cust`.
func (me *HiveSubs) WalkAllDirs(visitor ufs.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, ufs.WalkAllDirs(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, ufs.WalkAllDirs(dp, visitor)...)
	}
	return
}

//	`ufs.WalkAllFiles` for `me.Dist` and `me.Cust`.
func (me *HiveSubs) WalkAllFiles(visitor ufs.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, ufs.WalkAllFiles(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, ufs.WalkAllFiles(dp, visitor)...)
	}
	return
}

//	`ufs.WalkDirsIn` for `me.Dist` and `me.Cust`.
func (me *HiveSubs) WalkDirsIn(visitor ufs.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, ufs.WalkDirsIn(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, ufs.WalkDirsIn(dp, visitor)...)
	}
	return
}

//	`ufs.WalkFilesIn` for `me.Dist` and `me.Cust`.
func (me *HiveSubs) WalkFilesIn(visitor ufs.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, ufs.WalkFilesIn(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, ufs.WalkFilesIn(dp, visitor)...)
	}
	return
}

//	`ufs.DirWatcher.WatchIn` for `me.Dist` and `me.Cust`.
func (me *HiveSubs) WatchIn(handler ufs.WatcherHandler, runHandlerNow bool, subRelPath ...string) {
	me.root.fsWatcher.WatchIn(me.Dist.Path(subRelPath...), "*", runHandlerNow, handler)
	me.root.fsWatcher.WatchIn(me.Cust.Path(subRelPath...), "*", runHandlerNow, handler)
	return
}
