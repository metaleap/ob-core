package obcore

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	uio "github.com/metaleap/go-util/io"
)

const (
	//	Framework/platform title. Who knows, it might change..
	OB_TITLE = "OpenBase"
)

var (
	//	Provides access to the 'Hive-directory', the root directory
	//	containing configuration files, static web-served files, "template schema"
	//	files, package manifests and possibly data-base files depending on setup.
	Hive ObHive

	//	Runtime options
	Opt struct {
		//	Created in Init() and never nil, even if logging is disabled.
		Log *log.Logger

		//	Set this to true before calling Init() if the runtime is a sandboxed environment (such
		//	as Google App Engine) with security restrictions (no syscall, no unsafe, no file-writes)
		Sandboxed bool

		//	Set to true before Init() in cmd/ob-server/main.go.
		//	Should remain false in practically all other scenarios.
		//	(If true, much additional logic is executed and server-related resources allocated that
		//	are unneeded when importing this package in a "server-side, server-less client" scenario.)
		Server bool
	}
)

//	Clean-up. Call this when you're done working with this package and all allocated resources should be released.
func Dispose() {
	if Hive.Watch != nil {
		Hive.Watch.Close()
		Hive.Watch = nil
	}
}

//	Initialization. Call this before working with this package.
//	Before calling Init(), you may need to set Opt.Sandboxed, see Opt for details.
//	If logWriter is nil, logging is disabled and Opt.Log uses a uio.DiscardWriter.
//	In any event, Opt.Log doesn't log the err being returned, so be sure to check it.
//	If err is not nil, this package is not in a usable state and must not be used.
func Init(hiveDirPath string, logWriter io.Writer) (err error) {
	if logWriter == nil {
		logWriter = &uio.DiscardWriter{}
	}
	Opt.Log = log.New(logWriter, "", log.LstdFlags)
	if Opt.Server {
		log.Printf("[INIT]\t@ hive = '%s', server = %v, sandboxed = %v", hiveDirPath, Opt.Server, Opt.Sandboxed)
	}
	if !Opt.Sandboxed {
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
