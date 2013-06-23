package obcore

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-utils/ufs"
)

//	Provides access to a specified `Hive`-directory.
//
//	A `Hive` is the root directory with `dist` and `cust` sub-directories,
//	which contain configuration files, static web-served files, "template schema"
//	files, bundle manifests and possibly data-base files depending on setup.
type HiveRoot struct {
	//	The current `Hive`-directory path, set via `HiveRoot.Init`
	Dir string

	fsWatcher *ufs.Watcher

	//	Paths to some well-known `Hive` sub-directories
	Paths struct {
		//	{hive}/logs
		Logs string
	}

	//	Represents the `Hive` sub-directories `dist` and `cust`
	Subs HiveSubs
}

//	Creates a new log file at `{me.Dir}/logs/{date-time}.log`.
func (me *HiveRoot) CreateLogFile() (fullPath string, newOutFile *os.File, err error) {
	if err = ufs.EnsureDirExists(me.Paths.Logs); err == nil {
		now := time.Now()
		fullPath = filepath.Join(me.Paths.Logs, strf("%s__%v.log", now.Format("2006-01-02_15-04-05"), now.UnixNano()))
		newOutFile, err = os.Create(fullPath)
	}
	return
}

func (me *HiveRoot) dispose() (err error) {
	if me.fsWatcher != nil {
		err = me.fsWatcher.Close()
		me.fsWatcher = nil
	}
	return
}

func (me *HiveRoot) init(dir string) (err error) {
	me.Dir = dir
	me.Subs.init(me)
	me.Paths.Logs = me.Path("logs")
	if me.fsWatcher == nil {
		if me.fsWatcher, err = ufs.NewWatcher(); err == nil && me.fsWatcher != nil {
			go me.fsWatcher.Go()
		}
	}
	return
}

//	Returns whether the specified `dirPath` points to a valid `Hive`-directory.
func IsHive(dirPath string) bool {
	return ufs.DirsOrFilesExistIn(dirPath, "cust", "dist")
}

//	Returns a cleaned, `me.Dir`-joined full path for the specified `Hive`-relative path segments.
//
//	For example, if `me.Dir` is `obtest/hive`, then `me.Path("logs", "unknowable.log")` returns `obtest/hive/logs/unknowable.log`.
func (me *HiveRoot) Path(relPath ...string) (fullFsPath string) {
	if fullFsPath = filepath.Join(relPath...); len(me.Dir) > 0 {
		fullFsPath = filepath.Join(me.Dir, fullFsPath)
	}
	return
}
