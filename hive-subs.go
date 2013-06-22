package obcore

import (
	"github.com/go-utils/uio"
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

func (me *HiveSubs) WalkAllDirs(visitor uio.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, uio.WalkAllDirs(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, uio.WalkAllDirs(dp, visitor)...)
	}
	return
}

func (me *HiveSubs) WalkAllFiles(visitor uio.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, uio.WalkAllFiles(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, uio.WalkAllFiles(dp, visitor)...)
	}
	return
}

func (me *HiveSubs) WalkDirsIn(visitor uio.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, uio.WalkDirsIn(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, uio.WalkDirsIn(dp, visitor)...)
	}
	return
}

func (me *HiveSubs) WalkFilesIn(visitor uio.WalkerVisitor, relPath ...string) (errs []error) {
	dp := me.Dist.DirPath(relPath...)
	if len(dp) > 0 {
		errs = append(errs, uio.WalkFilesIn(dp, visitor)...)
	}
	if dp = me.Cust.DirPath(relPath...); len(dp) > 0 {
		errs = append(errs, uio.WalkFilesIn(dp, visitor)...)
	}
	return
}

func (me *HiveSubs) WatchIn(handler uio.WatcherHandler, runHandlerNow bool, subRelPath ...string) {
	me.root.fsWatcher.WatchIn(me.Dist.Path(subRelPath...), "*", runHandlerNow, handler)
	me.root.fsWatcher.WatchIn(me.Cust.Path(subRelPath...), "*", runHandlerNow, handler)
	return
}
