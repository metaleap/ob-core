package obsrv_daemon

import (
	"net/http"
	"os"
	"time"

	"github.com/go-utils/ugo"
	"github.com/go-utils/unet"

	ob "github.com/openbase/ob-core"
	obsrv "github.com/openbase/ob-core/server"
)

//	Provided just in case you need to customize the WriteTimeout, MaxHeaderBytes,
//	TLSConfig or TLSNextProto options before calling InitThenListenAndServe()
var HttpServer http.Server

//	Options for Main()
type Opt struct {
	//	TCP address as per http.Server
	HttpAddr string

	//	Set to true to have Main() write all log output to a new log file in {hive}/logs/
	LogToFile bool

	//	Set to true to suppress any and all writes to "standard-output"
	Silent bool

	//	If not 0, schedules a single emulated "GET /" warm-up "request" right
	//	from within this process, n Duration after HTTP-serving was initiated.
	//	(To be useful at all, this should probably be between 50-100ms and 1s.)
	WarmupRequestAfter time.Duration

	//	HTTPS via Transport Layer Security is supported only when both a
	//	CertFile and a KeyFile, as per http.ListenAndServeTLS(), are specified.
	TLS struct {
		//	File name containing a certificate for HTTPS-serving via TLS. If the
		//	certificate is signed by a certificate authority, this should be the
		//	concatenation of the server's certificate followed by the CA's certificate.
		CertFile string

		//	File name containing a matching private key for TLS serving.
		KeyFile string
	}
}

//	Called by func main() in cmd/ob-server/main.go package.
//
//	(Do note, this function does all initializations, defers all clean-ups and then runs 'forever'.)
func InitThenListenAndServe(hiveDir string, opt *Opt) (logFilePath string, err error) {
	//	pre-init
	ob.Opt.Server = true
	if len(hiveDir) == 0 {
		hiveDir, _ = os.Getwd()
		hiveDir = ob.Hive.GuessHiveRootDir(hiveDir)
	}

	//	init
	logger := ob.NewLogger(ugo.Ifw(opt.Silent, nil, os.Stdout))
	if err = ob.Init(hiveDir, logger); err != nil {
		return
	}
	defer ob.Dispose()

	//	create logger file?
	if opt.LogToFile {
		var (
			logFile *os.File
		)
		if logFilePath, logFile, err = ob.Hive.CreateLogFile(); err != nil {
			return
		} else {
			defer logFile.Close()
			logger.Infof("LOG @ %s", logFilePath)
			logger = ob.NewLogger(logFile)
			ob.Log = logger
		}
	}

	//	all systems go!
	obsrv.Init()
	HttpServer.Handler, HttpServer.Addr, HttpServer.ReadTimeout = obsrv.Router, opt.HttpAddr, 2*time.Minute
	https := len(opt.TLS.CertFile) > 0 && len(opt.TLS.KeyFile) > 0
	logger.Infof("LIVE @ %s", unet.Addr(ugo.Ifs(https, "https", "http"), HttpServer.Addr))
	if opt.WarmupRequestAfter.Nanoseconds() > 0 {
		go localWarmupRequest(opt.WarmupRequestAfter)
	}
	if https {
		err = HttpServer.ListenAndServeTLS(opt.TLS.CertFile, opt.TLS.KeyFile)
	} else {
		err = HttpServer.ListenAndServe()
	}
	return
}

func localWarmupRequest(after time.Duration) {
	time.Sleep(after)
	var w unet.ResponseBuffer
	if r, err := http.NewRequest("GET", "/", nil); r != nil {
		r.Header["User-Agent"] = []string{"LocalWarmup"}
		now := time.Now()
		HttpServer.Handler.ServeHTTP(&w, r)
		ob.Log.Infof("Warmup request served in %v", time.Now().Sub(now))
	} else {
		ob.Log.Error(err)
	}
}
