package core

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
	Server    ObServer
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
		if err = Hive.init(); err == nil {
			Server.init()
		}
	}
	return
}

func ListenAndServe(addr, tlsCertFile, tlsKeyFile string) (err error) {
	if Sandboxed {
		err = fmt.Errorf("Cannot call ListenAndServe() in Sandboxed mode")
		return
	}
	Server.Http.Addr = addr
	if len(tlsCertFile) > 0 && len(tlsKeyFile) > 0 {
		return Server.Http.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	}
	return Server.Http.ListenAndServe()
}
