package obsrv

import (
	"os"

	ugo "github.com/metaleap/go-util"
	ustr "github.com/metaleap/go-util/str"

	ob "github.com/openbase/ob-core"
)

//	Called by func main() in cmd/ob-server/main.go package.
//	In theory, you could 'abuse' this by writing your own custom server
//	executable instead of using ob-server, but this isn't really intended.
//	(Do note, this function does run 'forever' and thus defers all cleanups.)
func Main(hiveDir, httpAddr, tlsCertFile, tlsKeyFile string, logToFile, silent bool) {
	var err error

	//	pre-init
	if len(hiveDir) == 0 {
		hiveDir, _ = os.Getwd()
		hiveDir = ob.Hive.GuessDir(hiveDir)
	}
	ob.Opt.Server = true

	//	init
	log := ob.NewLogger(os.Stdout)
	if silent {
		log.Out = nil
	}
	if err = ob.Init(hiveDir, log); err != nil {
		log.Fatal(err)
	}
	defer ob.Dispose()

	//	create log file?
	if logToFile {
		var (
			logFile     *os.File
			logFilePath string
		)
		if logFilePath, logFile, err = ob.Hive.CreateLogFile(); err != nil {
			log.Fatal(err)
		} else {
			defer logFile.Close()
			log.Infof("LOG @ %s", logFilePath)
			log.Out = logFile
		}
	}

	//	all systems go!
	proto := ugo.Ifs(len(tlsCertFile) > 0 && len(tlsKeyFile) > 0, "https", "http")
	log.Infof("LIVE @ %s://%s%s", proto, ugo.HostName(), ustr.StripSuffix(httpAddr, ":"+proto))
	Init()
	log.Fatal(ListenAndServe(httpAddr, tlsCertFile, tlsKeyFile))
}
