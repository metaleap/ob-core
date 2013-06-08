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
