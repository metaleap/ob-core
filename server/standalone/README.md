# obsrv_daemon
--
    import "github.com/openbase/ob-core/server/standalone"

Used by cmd/ob-server/main.go

## Usage

```go
const (
	//	The name of the environment variable storing the `Hive`-directory path, if set.
	//	Used as a fall-back by `HiveRoot.GuessDir`.
	ENV_OBHIVE = "OBHIVE"
)
```

```go
var HttpServer http.Server
```
Provided just in case you need to customize the `WriteTimeout`,
`MaxHeaderBytes`, `TLSConfig` or `TLSNextProto` options prior to calling
`InitThenListenAndServe`.

#### func  InitThenListenAndServe

```go
func InitThenListenAndServe(hiveDir string, opt *Opt) (logFilePath string, err error)
```
Called by `func main` in `cmd/ob-server/main.go` package.

(Do note, this function does all initializations, defers all clean-ups and then
runs 'forever'.)

#### type Opt

```go
type Opt struct {
	//	TCP address as per `http.Server.Addr`.
	HttpAddr string

	//	Set to `true` to have `InitThenListenAndServe` write all log output to a new log file in `{hive}/logs/`.
	LogToFile bool

	//	Set to `true` to suppress any and all writes to `stdout`.
	Silent bool

	//	If not 0, schedules a single emulated (transport-less) `GET /` "warm-up request"
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
```

Options for `InitThenListenAndServe`.

--
**godocdown** http://github.com/robertkrimen/godocdown
