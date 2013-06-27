package obsrv_daemon

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-utils/ugo"
	"github.com/go-utils/unet"

	ob "github.com/openbase/ob-core"
	obsrv "github.com/openbase/ob-core/server"
)

const (
	//	The name of the environment variable storing the `Hive`-directory path, if set.
	ENV_OBHIVE = "OBHIVE"
)

//	Provided just in case you need to customize the `http.Server`
//	being used prior to calling `InitThenListenAndServe`.
//
//	Its `ReadTimeout` is `init`ialized to `2 * time.Minute`.
var HttpServer http.Server

func init() {
	HttpServer.ReadTimeout = 2 * time.Minute
}

//	Options for `InitThenListenAndServe`.
type Opt struct {
	//	TCP address for `HttpServer.Addr`.
	HttpAddr string

	//	Set to `true` to have `InitThenListenAndServe` write all log output to a new log file in `{hive}/logs/`.
	LogToFile bool

	//	Set to `true` to suppress any and all writes to `stdout`.
	Silent bool

	//	If not 0, schedules a single emulated (in-process, transport-less) `GET /` "warm-up request"
	//	right from within this process, n `time.Duration` after HTTP-serving was initiated.
	//	(To be useful at all, this should probably be between 50-100ms and 1s.)
	WarmupRequestAfter time.Duration

	//	HTTPS via Transport Layer Security is supported only when both a
	//	`CertFile` and a `KeyFile` --as per `http.ListenAndServeTLS`-- are specified.
	TLS struct {
		//	File name containing a certificate for HTTPS-serving via TLS. If the
		//	certificate is signed by a certificate authority, this should be the
		//	concatenation of the server's certificate followed by the CA's certificate.
		CertFile string

		//	File name containing a matching private key for TLS serving.
		KeyFile string
	}
}

//	If `dirPath` indicates a valid `Hive`-directory path, returns its `filepath.Abs` equivalent;
//	otherwise returns the `$OBHIVE` environment variable regardless of `Hive`-directory validity.
func HiveDir(dirPath string) string {
	if dirPath, _ = filepath.Abs(dirPath); !ob.IsHive(dirPath) {
		dirPath = filepath.Clean(os.Getenv(ENV_OBHIVE))
	}
	return dirPath
}

//	Called by `func main` in `openbase/ob-core/cmd/ob-server`.
//
//	Sanitizes the specified `hiveDir` via the `HiveDir` function. Overrides `HttpServer.Addr` and `HttpServer.Handler`.
//
//	(Do note, this function does all initializations, defers all clean-ups and then runs 'forever'.)
func InitThenListenAndServe(hiveDir string, opt *Opt) (logFilePath string, err error) {
	var ctx *obsrv.Ctx
	//	pre-init
	hiveDir = HiveDir(hiveDir)

	//	init
	logger := ob.NewLogger(ugo.Ifw(opt.Silent, nil, os.Stdout))
	if ctx, err = obsrv.NewCtx(hiveDir, logger); err == nil {
		defer ctx.Dispose()
		logger.Infof("INIT @ %s", ctx.Hive.Dir)

		//	create logger file?
		if opt.LogToFile {
			var logFile *os.File
			if logFilePath, logFile, err = ctx.Hive.CreateLogFile(); err != nil {
				return
			} else {
				defer logFile.Close()
				logger.Infof("LOG @ %s", logFilePath)
				logger = ob.NewLogger(logFile)
				ctx.Log = logger
			}
		}

		//	all systems go!
		HttpServer.Handler, HttpServer.Addr = obsrv.NewHttpHandler(ctx), opt.HttpAddr
		https := len(opt.TLS.CertFile) > 0 && len(opt.TLS.KeyFile) > 0
		logger.Infof("LIVE @ %s", unet.Addr(ugo.Ifs(https, "https", "http"), HttpServer.Addr))
		if opt.WarmupRequestAfter > 0 {
			go localWarmup(ctx, opt.WarmupRequestAfter, 4)
		}
		if https {
			err = HttpServer.ListenAndServeTLS(opt.TLS.CertFile, opt.TLS.KeyFile)
		} else {
			err = HttpServer.ListenAndServe()
		}
	}
	return
}

func localWarmup(ctx *obsrv.Ctx, after time.Duration, num int) {
	req := func() {
		var w unet.ResponseBuffer
		if r, err := http.NewRequest("GET", "/", nil); r != nil {
			r.Header["User-Agent"] = []string{"LocalWarmup"}
			now := time.Now()
			HttpServer.Handler.ServeHTTP(&w, r)
			ctx.Log.Infof("Warmup `GET /` took %v", time.Now().Sub(now))
		} else {
			ctx.Log.Error(err)
		}
	}
	time.Sleep(after)
	//	cheap way to allow for -race detection?
	for i := 0; i < num; i++ {
		go req()
	}
}
