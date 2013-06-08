package obcore

import (
	"os"
	"path/filepath"

	uio "github.com/metaleap/go-util/io"
)

const (
	ENV_OBHIVE = "OBHIVE"
)

type ObHive struct {
	DirPath string
	Watch   *uio.Watcher
}

func (_ ObHive) init() (err error) {
	if !Sandboxed {
		if Hive.Watch, err = uio.NewWatcher(); err == nil && Hive.Watch != nil {
			go Hive.Watch.Go()
		}
	}
	return
}

func (_ ObHive) GuessDirPath(userSpecified string) (guess string) {
	if guess = userSpecified; !Hive.IsHive(guess) {
		guess = os.Getenv(ENV_OBHIVE)
	}
	return
}

func (_ ObHive) IsHive(dirPath string) bool {
	return uio.DirsFilesExist(dirPath, "client")
}

func (_ ObHive) Path(names ...string) (fullFsPath string) {
	fullFsPath = filepath.Join(names...)
	if len(Hive.DirPath) > 0 {
		fullFsPath = filepath.Join(Hive.DirPath, fullFsPath)
	}
	return
}
