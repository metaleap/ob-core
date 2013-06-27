# obsrv_daemon
--
    import "github.com/openbase/ob-core/server/standalone"

Used by `openbase/ob-core/cmd/ob-server`

## Usage

```go
const (
	//	The name of the environment variable storing the `Hive`-directory path, if set.
	ENV_OBHIVE = "OBHIVE"
)
```

```go
var HttpServer http.Server
```
Provided just in case you need to customize the `http.Server` being used prior
to calling `InitThenListenAndServe`.

Its `ReadTimeout` is `init`ialized to `2 * time.Minute`.

#### func  HiveDir

```go
func HiveDir(dirPath string) string
```
If `dirPath` indicates a valid `Hive`-directory path, returns its `filepath.Abs`
equivalent; otherwise returns the `$OBHIVE` environment variable regardless of
`Hive`-directory validity.

#### func  InitThenListenAndServe

```go
func InitThenListenAndServe(hiveDir string, opt *Opt) (logFilePath string, err error)
```
Called by `func main` in `openbase/ob-core/cmd/ob-server`.

Sanitizes the specified `hiveDir` via the `HiveDir` function. Overrides
`HttpServer.Addr` and `HttpServer.Handler`.

(Do note, this function does all initializations, defers all clean-ups and then
runs 'forever'.)

#### type Opt

```go
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
```

Options for `InitThenListenAndServe`.

--
**godocdown** http://github.com/robertkrimen/godocdown
