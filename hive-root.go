package obcore

import (
	"os"
	"path/filepath"

	uio "github.com/metaleap/go-util/io"
)

const (
	//	The name of the environment variable storing the Hive-directory path, if set.
	//	Used as a fall-back by Hive.GuessDir().
	ENV_OBHIVE = "OBHIVE"
)

var (
	//	Provides access to the 'Hive-directory' used throughout the package.
	//	The 'Hive' is the root directory with the 'dist' and 'cust' sub-directories,
	//	which contain configuration files, static web-served files, "template schema"
	//	files, package manifests and possibly data-base files depending on setup.
	Hive HiveRoot
)

//	Provides access to a specified Hive directory.
type HiveRoot struct {
	//	The current Hive-directory path, set via Hive.Init()
	Dir string

	fsWatcher *uio.Watcher

	Paths struct {
		Logs string
	}

	Subs struct {
		Cust, Dist HiveSub
	}
}

//	Creates a new log file at: {me.Dir}/logs/{date-time}.log
func (me *HiveRoot) CreateLogFile() (fullPath string, newOutFile *os.File, err error) {
	if err = uio.EnsureDirExists(me.Paths.Logs); err == nil {
		fullPath = filepath.Join(me.Paths.Logs, strf("%s__%v.log", Opt.initTime.Format("2006-01-02_15-04-05"), Opt.initTime.UnixNano()))
		newOutFile, err = os.Create(fullPath)
	}
	return
}

func (me *HiveRoot) dispose() {
	if me.fsWatcher != nil {
		me.fsWatcher.Close()
		me.fsWatcher = nil
	}
}

func (me *HiveRoot) FileExists(subRelPath ...string) bool {
	return len(me.FilePath(subRelPath...)) > 0
}

func (me *HiveRoot) FilePath(subRelPath ...string) (filePath string) {
	if filePath = me.Subs.Cust.FilePath(subRelPath...); len(filePath) == 0 {
		filePath = me.Subs.Dist.FilePath(subRelPath...)
	}
	return
}

//	Returns userSpecified if that is a valid Hive-directory path as per HiveRoot.IsHive(),
//	else returns the value of the OBHIVE environment variable (regardless of path validity).
func (me *HiveRoot) GuessDir(userSpecified string) (guess string) {
	if guess = userSpecified; !me.IsHive(guess) {
		guess = os.Getenv(ENV_OBHIVE)
	}
	return
}

//	Initializes me.Dir to the specified dir (without checking it, call IsHive() beforehand to do so).
//	Then initializes me.Subs and me.Paths based on me.Dir.
func (me *HiveRoot) Init(dir string) {
	me.Dir = dir
	me.Subs.Cust.init(me, "cust")
	me.Subs.Dist.init(me, "dist")
	p := &me.Paths
	p.Logs = me.Path("logs")
}

func (me *HiveRoot) init(dir string) (err error) {
	me.Init(dir)
	if me.fsWatcher == nil {
		if me.fsWatcher, err = uio.NewWatcher(); err == nil && me.fsWatcher != nil {
			go me.fsWatcher.Go()
		}
	}
	return
}

//	Returns true if the specified directory path points to a valid Hive-directory.
func (_ *HiveRoot) IsHive(dir string) bool {
	return uio.DirsFilesExist(dir, "cust", "dist")
}

//	Returns a cleaned, me.Dir-joined full path for the specified Hive-relative path segments.
//	For example, if me.Dir is "obtest/hive", then me.Path("pkg", "mysql") returns "obtest/hive/pkg/mysql"
func (me *HiveRoot) Path(relPath ...string) (fullFsPath string) {
	if fullFsPath = filepath.Join(relPath...); len(me.Dir) > 0 {
		fullFsPath = filepath.Join(me.Dir, fullFsPath)
	}
	return
}

func (me *HiveRoot) WatchDualDir(handler uio.WatcherHandler, runHandlerNow bool, subRelPath ...string) {
	me.fsWatcher.WatchIn(me.Subs.Dist.Path(subRelPath...), "*", runHandlerNow, handler)
	me.fsWatcher.WatchIn(me.Subs.Cust.Path(subRelPath...), "*", runHandlerNow, handler)
	return
}
