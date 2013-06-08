package obcore

import (
	"os"
	"path/filepath"

	uio "github.com/metaleap/go-util/io"
)

const (
	//	The name of the environment variable storing the Hive-directory path, if set.
	//	Used as a fall-back by ObHive.GuessDirPath().
	ENV_OBHIVE = "OBHIVE"
)

var (
	//	Provides access to the 'Hive-directory', the root directory
	//	containing configuration files, static web-served files, "template schema"
	//	files, package manifests and possibly data-base files depending on setup.
	Hive ObHive
)

//	Singleton type. Only valid use is the exported Hive variable.
type ObHive struct {
	//	The current Hive-directory path
	DirPath string

	//	The uio.Watcher used for any and all Hive-directory change-event notifications.
	//	This is nil if Opt.Sandboxed is true.
	Watch *uio.Watcher
}

func (_ ObHive) init() (err error) {
	if !Opt.Sandboxed {
		if Hive.Watch, err = uio.NewWatcher(); err == nil && Hive.Watch != nil {
			go Hive.Watch.Go()
		}
	}
	return
}

//	Creates a new log file at: [Hive.DirPath]/log/[date-time].log
func (_ ObHive) CreateLogFile() (fullPath string, newOutFile *os.File, err error) {
	logDirPath := Hive.Path("log")
	if err = uio.EnsureDirExists(logDirPath); err == nil {
		fullPath = filepath.Join(logDirPath, strf("%s.log", Opt.initTime.Format("2006-01-02_15-04-05")))
		newOutFile, err = os.Create(fullPath)
	}
	return
}

//	Returns userSpecified if that is a valid Hive-directory path as per ObHive.IsHive(),
//	else returns the value of the OBHIVE environment variable (regardless of path validity).
func (_ ObHive) GuessDirPath(userSpecified string) (guess string) {
	if guess = userSpecified; !Hive.IsHive(guess) {
		guess = os.Getenv(ENV_OBHIVE)
	}
	return
}

//	Returns true if the specified path points to a valid Hive-directory.
func (_ ObHive) IsHive(dirPath string) bool {
	return uio.DirsFilesExist(dirPath, "client")
}

//	Returns a cleaned Hive-relative path for the specified path segments.
//	For example, if Hive.DirPath is "obtest/test2", then Hive.Path("pkg", "mysql") returns "obtest/test2/pkg/mysql"
func (_ ObHive) Path(segments ...string) (fullFsPath string) {
	fullFsPath = filepath.Join(segments...)
	if len(Hive.DirPath) > 0 {
		fullFsPath = filepath.Join(Hive.DirPath, fullFsPath)
	}
	return
}
