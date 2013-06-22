package obcore

import (
	"github.com/go-utils/ufs"
)

//	Used for Hive.Subs
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

//	me.Dist.FileExists(subRelPath...) || me.Dist.FileExists(subRelPath...)
func (me *HiveSubs) FileExists(subRelPath ...string) bool {
	return me.Dist.FileExists(subRelPath...) || me.Cust.FileExists(subRelPath...)
}

//	Returns either me.Cust.FilePath(subRelPath ...) or me.Dist.FilePath(subRelPath ...)
func (me *HiveSubs) FilePath(subRelPath ...string) (filePath string) {
	if filePath = me.Cust.FilePath(subRelPath...); len(filePath) == 0 {
		filePath = me.Dist.FilePath(subRelPath...)
	}
	return
}

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

func (me *HiveSubs) WatchIn(handler ufs.WatcherHandler, runHandlerNow bool, subRelPath ...string) {
	me.root.fsWatcher.WatchIn(me.Dist.Path(subRelPath...), "*", runHandlerNow, handler)
	me.root.fsWatcher.WatchIn(me.Cust.Path(subRelPath...), "*", runHandlerNow, handler)
	return
}
