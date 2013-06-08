package obcore

import (
	"fmt"
	"log"
	"path/filepath"
)

const (
	OB_TITLE = "OpenBase"
)

var (
	Hive      ObHive
	Sandboxed bool
)

func Dispose() {
	if Hive.Watch != nil {
		Hive.Watch.Close()
		Hive.Watch = nil
	}
}

func Init(hiveDirPath string) (err error) {
	log.Printf("[INIT]\t@ hive = '%s', sandboxed = %v", hiveDirPath, Sandboxed)
	if !Sandboxed {
		if hiveDirPath, err = filepath.Abs(hiveDirPath); (err == nil) && !Hive.IsHive(hiveDirPath) {
			err = fmt.Errorf("Not a valid %s Hive directory installation: '%s'.", OB_TITLE, hiveDirPath)
		}
	}
	if err == nil {
		Hive.DirPath = hiveDirPath
		err = Hive.init()
	}
	return
}
