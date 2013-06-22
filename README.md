# obcore
--
    import "github.com/openbase/ob-core"

Core ('kernel'-level, but server-less) functionality package

## Usage

```go
const (
	//	The name of the environment variable storing the Hive-directory path, if set.
	//	Used as a fall-back by Hive.GuessDir().
	ENV_OBHIVE = "OBHIVE"
)
```

```go
const (
	//	Framework/platform title. Who knows, it might change..
	OB_TITLE = "OpenBase"
)
```

```go
var (
	//	Runtime options
	Opt struct {

		//	Set this to true before calling Init() if the runtime is a sandboxed environment (such
		//	as Google App Engine) with security restrictions (no syscall, no unsafe, no file-writes)
		Sandboxed bool

		//	Set to true before Init() in cmd/ob-server/main.go.
		//	Should remain false in practically all other scenarios.
		//	(If true, much additional logic is executed and server-related resources allocated that
		//	are unneeded when importing this package in a "server-side but server-less client" scenario.)
		Server bool
	}
)
```

#### func  Dispose

```go
func Dispose()
```
Clean-up. Call this when you're done working with this package and all allocated
resources should be released.

#### func  Init

```go
func Init(hiveDirPath string, logger Logger) (err error)
```
Initialization. Call this before working with this package. Before calling
Init(), you may need to set Opt.Sandboxed, see Opt for details. If logger is
nil, Log is set to a no-op dummy and logging is disabled. In any event, Init()
doesn't log the err being returned, so be sure to check it. If err is not nil,
this package is not in a usable state and must not be used.

#### type HiveRoot

```go
type HiveRoot struct {
	//	The current Hive-directory path, set via Hive.Init()
	Dir string

	//	Paths to some well-known HiveRoot directories
	Paths struct {
		//	{hive}/logs
		Logs string
	}

	//	Sub-hives
	Subs HiveSubs
}
```

Provides access to a specified Hive directory.

```go
var (
	//	Provides access to the 'Hive-directory' used throughout the package.
	//	The 'Hive' is the root directory with the 'dist' and 'cust' sub-directories,
	//	which contain configuration files, static web-served files, "template schema"
	//	files, bundle manifests and possibly data-base files depending on setup.
	Hive HiveRoot
)
```

#### func (*HiveRoot) CreateLogFile

```go
func (me *HiveRoot) CreateLogFile() (fullPath string, newOutFile *os.File, err error)
```
Creates a new log file at: {me.Dir}/logs/{date-time}.log

#### func (*HiveRoot) GuessHiveRootDir

```go
func (me *HiveRoot) GuessHiveRootDir(userSpecified string) (guess string)
```
Returns userSpecified if that is a valid Hive-directory path as per
HiveRoot.IsHive(), else returns the value of the OBHIVE environment variable
(regardless of path validity).

#### func (*HiveRoot) Init

```go
func (me *HiveRoot) Init(dir string)
```
Initializes me.Dir to the specified dir (without checking it, call IsHive()
beforehand to do so). Then initializes me.Subs and me.Paths based on me.Dir.

#### func (*HiveRoot) IsHive

```go
func (_ *HiveRoot) IsHive(dir string) bool
```
Returns true if the specified directory path points to a valid Hive-directory.

#### func (*HiveRoot) Path

```go
func (me *HiveRoot) Path(relPath ...string) (fullFsPath string)
```
Returns a cleaned, me.Dir-joined full path for the specified Hive-relative path
segments. For example, if me.Dir is "obtest/hive", then me.Path("pkg", "mysql")
returns "obtest/hive/pkg/mysql"

#### type HiveSub

```go
type HiveSub struct {

	//	Paths to some well-known HiveSub directories
	Paths struct {
		//	{hive}/{sub}/client
		Client string

		//	{hive}/{sub}/client/pub
		ClientPub string

		//	{hive}/{sub}/pkg
		Pkg string
	}
}
```

Represents either the dist or the cust directory in a hive directory

#### func (*HiveSub) DirExists

```go
func (me *HiveSub) DirExists(relPath ...string) bool
```
Returns whether the specified {hive}/{sub}-relative directory exists

#### func (*HiveSub) DirPath

```go
func (me *HiveSub) DirPath(relPath ...string) (dirPath string)
```
Returns a {hive}/{sub}-joined representation of the specified
{hive}/{sub}-relative path, if it represents an existing directory

#### func (*HiveSub) FileExists

```go
func (me *HiveSub) FileExists(relPath ...string) bool
```
Returns whether the specified {hive}/{sub}-relative file exists

#### func (*HiveSub) FilePath

```go
func (me *HiveSub) FilePath(relPath ...string) (filePath string)
```
Returns a {hive}/{sub}-joined representation of the specified
{hive}/{sub}-relative path, if it represents an existing file

#### func (*HiveSub) Path

```go
func (me *HiveSub) Path(relPath ...string) string
```
Returns a {hive}/{sub}-joined representation of the specified
{hive}/{sub}-relative path

#### func (*HiveSub) ReadFile

```go
func (me *HiveSub) ReadFile(relPath ...string) (data []byte, err error)
```
Reads the file at the specified {hive}/{sub}-relative path

#### type HiveSubs

```go
type HiveSubs struct {

	//	{hive}/dist/
	Dist HiveSub

	//	{hive}/cust/
	Cust HiveSub
}
```

Used for Hive.Subs

#### func (*HiveSubs) FileExists

```go
func (me *HiveSubs) FileExists(subRelPath ...string) bool
```
me.Dist.FileExists(subRelPath...) || me.Dist.FileExists(subRelPath...)

#### func (*HiveSubs) FilePath

```go
func (me *HiveSubs) FilePath(subRelPath ...string) (filePath string)
```
Returns either me.Cust.FilePath(subRelPath ...) or me.Dist.FilePath(subRelPath
...)

#### func (*HiveSubs) WalkAllDirs

```go
func (me *HiveSubs) WalkAllDirs(visitor uio.WalkerVisitor, relPath ...string) (errs []error)
```

#### func (*HiveSubs) WalkAllFiles

```go
func (me *HiveSubs) WalkAllFiles(visitor uio.WalkerVisitor, relPath ...string) (errs []error)
```

#### func (*HiveSubs) WalkDirsIn

```go
func (me *HiveSubs) WalkDirsIn(visitor uio.WalkerVisitor, relPath ...string) (errs []error)
```

#### func (*HiveSubs) WalkFilesIn

```go
func (me *HiveSubs) WalkFilesIn(visitor uio.WalkerVisitor, relPath ...string) (errs []error)
```

#### func (*HiveSubs) WatchIn

```go
func (me *HiveSubs) WatchIn(handler uio.WatcherHandler, runHandlerNow bool, subRelPath ...string)
```

#### type Hub

```go
type Hub struct {
}
```


#### func  GetHub

```go
func GetHub(path string) (hub *Hub)
```

#### func  RootHub

```go
func RootHub() *Hub
```

#### func (*Hub) Parent

```go
func (me *Hub) Parent() (parent *Hub)
```

#### type Logger

```go
type Logger interface {
	// Debugf formats its arguments according to the format, analogous to fmt.Printf,
	// and records the text as a log message at Debug level
	Debugf(format string, args ...interface{})

	// Infof is like Debugf, but at Info level
	Infof(format string, args ...interface{})

	// Warningf is like Debugf, but at Warning level
	Warningf(format string, args ...interface{})

	// Error is like Debugf, but at Error level
	Error(err error) error

	// Errorf is like Debugf, but at Error level
	Errorf(format string, args ...interface{})

	// Criticalf is like Debugf, but at Critical level
	Criticalf(format string, args ...interface{})
}
```

An interface for log output. ObLogger provides the canonical implementation

```go
var (
	//	Set via Init(), never nil (even if logging is disabled)
	Log Logger
)
```

#### type ObLogger

```go
type ObLogger struct {
}
```

The canonical implementation of the Logger interface, using a standard
log.Logger

#### func  NewLogger

```go
func NewLogger(out io.Writer) (me *ObLogger)
```
Creates and returns a new ObLogger with the specified out io.Writer

#### func (*ObLogger) Criticalf

```go
func (me *ObLogger) Criticalf(format string, args ...interface{})
```
Criticalf is like Debugf, but at Critical level

#### func (*ObLogger) Debugf

```go
func (me *ObLogger) Debugf(format string, args ...interface{})
```
Debugf formats its arguments according to the format, analogous to fmt.Printf,
and records the text as a log message at Debug level

#### func (*ObLogger) Error

```go
func (me *ObLogger) Error(err error) error
```
Error is like Debugf, but at Error level. Returns err.

#### func (*ObLogger) Errorf

```go
func (me *ObLogger) Errorf(format string, args ...interface{})
```
Errorf is like Debugf, but at Error level

#### func (*ObLogger) Infof

```go
func (me *ObLogger) Infof(format string, args ...interface{})
```
Infof is like Debugf, but at Info level

#### func (*ObLogger) Warningf

```go
func (me *ObLogger) Warningf(format string, args ...interface{})
```
Warningf is like Debugf, but at Warning level

--
**godocdown** http://github.com/robertkrimen/godocdown