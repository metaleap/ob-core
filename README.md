# obcore
--
    import "github.com/openbase/ob-core"

Core ('kernel'-level, but server-less) functionality package

## Usage

```go
const (
	//	Framework/platform title. Who knows, it might change..
	OB_TITLE = "OpenBase"
)
```

```go
var (
	//	Contains one `BundleCfgReloader` handler per bundle kind.
	//	When a `Bundle` gets (re)loaded, after populating its `CfgRaw` hash-maps,
	//	it calls the appropriate `BundleCfgReloader` associated with its `Kind` to
	//	refresh its `Cfg` according to its potentially new or changed `CfgRaw`.
	BundleCfgLoaders = map[string]BundleCfgReloader{}
)
```

#### func  IsHive

```go
func IsHive(dirPath string) bool
```
Returns whether the specified `dirPath` points to a valid `Hive`-directory.

#### type Bundle

```go
type Bundle struct {

	//	The kind of this `Bundle`, according to its directory name,
	//	for example `webuilib`.
	Kind string

	//	The name of this `Bundle`, not including its `Kind`.
	Name string

	//	The full identifier of this `Bundle`, which is `Kind` and `Name`
	//	joined by a dash, for example `webuilib-jquery`.
	NameFull string

	//	Diagnostic info
	Diag struct {
		//	Full `Bundle` names of all `Info.Require` entries that
		//	are not currently installed inside `{hive}/{sub}/pkg/`.
		BadDeps []string

		//	The `error` that occurred when loading the `.ob-pkg` file, if any.
		//	Outside of unlikely file-system issues, this is most likely a TOML syntax error in the file.
		LoadErr error
	}

	//	Information from the `[bundle]` section of the `.ob-pkg` bundle configuration file.
	Info struct {
		//	Human-readable bundle title
		Title string

		//	Human-readable, comprehensive bundle description
		Desc string

		//	Web address for more information, in the case of 3rd-party bundles
		Www string

		//	Names fully-qualified bundles required for this bundle to function
		Require []string
	}

	//	Represents `CfgRaw` in a 'native'/'parsed'/'processed', bundle-specific way.
	//	To be set by a `BundleCfgReloader` registered in `BundleCfgLoaders`.
	Cfg interface{}

	//	Unprocessed information loaded from the `.ob-pkg` bundle configuration file.
	CfgRaw struct {
		//	Information from the `[default]` section.
		Default BundleCfg

		//	Information from any other sections.
		More map[string]BundleCfg
	}
}
```

Represents a bundle package found in a
`{hive}/{sub}/pkg/{kind-name}/{name.kind.ob-pkg}` file.

#### func (*Bundle) Ctx

```go
func (me *Bundle) Ctx() *Ctx
```

#### type BundleCfg

```go
type BundleCfg map[string]interface{}
```

Type for `Bundle.CfgRaw`.

#### type BundleCfgReloader

```go
type BundleCfgReloader func(*Bundle)
```

Used by `Bundle.Kind`-specific imports to register their reload handlers with
`BundleCfgLoaders`.

#### type BundleRegistry

```go
type BundleRegistry struct {
}
```

Bundle package registry, accessed from `Ctx.Bundles`.

#### func (*BundleRegistry) ByKind

```go
func (me *BundleRegistry) ByKind(kind string, deps []string) (all Bundles)
```
If `deps` is empty, returns all `Bundles` of the specified `kind`.

Otherwise, out of all `Bundles` specified in `deps` or directly or indirectly
required by them, returns those of the specified `kind`.

#### func (*BundleRegistry) ByName

```go
func (me *BundleRegistry) ByName(kindAndName ...string) *Bundle
```
Returns the `*Bundle` with the specified fully-qualified identifier.
`kindAndName` can be a single string such as `"webuilib-jquery"`, or 2 strings
for `kind` and `name`, such as `"webuilib", "jquery"`.

#### type Bundles

```go
type Bundles []*Bundle
```

A collection of `*Bundle`s.

#### func (Bundles) Len

```go
func (me Bundles) Len() int
```
Implements `sort.Interface.Len()`.

#### func (Bundles) Less

```go
func (me Bundles) Less(i, j int) bool
```
Implements `sort.Interface.Less()`.

#### func (Bundles) Swap

```go
func (me Bundles) Swap(i, j int)
```
Implements `sort.Interface.Swap()`.

#### type Ctx

```go
type Ctx struct {
	//	Represents access to the `Hive`-directory.
	Hive HiveRoot

	//	Set via `NewCtx`, never `nil` (even if logging is disabled).
	Log Logger
}
```

Global access. ONLY valid when initialized via `NewCtx`.

#### func  NewCtx

```go
func NewCtx(hiveDir string, server, sandboxed bool, logger Logger) (me *Ctx, err error)
```
Initializes and returns a new `*Ctx` providing access to the specified
`hiveDir`.

- `hiveDir`: the `Hive`-directory path accessed by `me`.

- `server` and `sandboxed` should be `false`, unless the caller is from
`ob-core/server/standalone` or `ob-gae`.

- If `logger` is `nil`, `me.Log` is set to a no-op dummy and logging is
disabled. In any event, `Init` doesn't log the `err` being returned (if any), so
be sure to handle it.

Whenever `err` is `nil`, `me` is non-`nil` and vice versa.

#### func (*Ctx) Bundles

```go
func (me *Ctx) Bundles() *BundleRegistry
```
Returns the bundle package registry for `me`.

#### func (*Ctx) Dispose

```go
func (me *Ctx) Dispose() (err error)
```
Clean-up when you're shutting down.

#### type HiveRoot

```go
type HiveRoot struct {
	//	The current `Hive`-directory path, set via `HiveRoot.Init`
	Dir string

	//	Paths to some well-known `Hive` sub-directories
	Paths struct {
		//	{hive}/logs
		Logs string
	}

	//	Represents the `Hive` sub-directories `dist` and `cust`
	Subs HiveSubs
}
```

Provides access to a specified `Hive`-directory.

A `Hive` is the root directory with `dist` and `cust` sub-directories, which
contain configuration files, static web-served files, "template schema" files,
bundle manifests and possibly data-base files depending on setup.

#### func (*HiveRoot) CreateLogFile

```go
func (me *HiveRoot) CreateLogFile() (fullPath string, newOutFile *os.File, err error)
```
Creates a new log file at `{me.Dir}/logs/{date-time}.log`.

#### func (*HiveRoot) Path

```go
func (me *HiveRoot) Path(relPath ...string) (fullFsPath string)
```
Returns a cleaned, `me.Dir`-joined full path for the specified `Hive`-relative
path segments.

For example, if `me.Dir` is `obtest/hive`, then `me.Path("logs",
"unknowable.log")` returns `obtest/hive/logs/unknowable.log`.

#### type HiveSub

```go
type HiveSub struct {

	//	Paths to some well-known `HiveSub` directories
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

Represents either the `dist` or the `cust` sub-directory inside a `HiveRoot`.

#### func (*HiveSub) DirExists

```go
func (me *HiveSub) DirExists(relPath ...string) bool
```
Returns whether the specified `{hive}/{sub}`-relative directory exists.

#### func (*HiveSub) DirPath

```go
func (me *HiveSub) DirPath(relPath ...string) (dirPath string)
```
Returns a `{hive}/{sub}`-joined representation of the specified
`{hive}/{sub}`-relative path, if it represents an existing directory.

#### func (*HiveSub) FileExists

```go
func (me *HiveSub) FileExists(relPath ...string) bool
```
Returns whether the specified `{hive}/{sub}`-relative file exists.

#### func (*HiveSub) FilePath

```go
func (me *HiveSub) FilePath(relPath ...string) (filePath string)
```
Returns a `{hive}/{sub}`-joined representation of the specified
`{hive}/{sub}`-relative path, if it represents an existing file.

#### func (*HiveSub) Path

```go
func (me *HiveSub) Path(relPath ...string) string
```
Returns a `{hive}/{sub}`-joined representation of the specified
`{hive}/{sub}`-relative path (regardless of whether it exists).

#### func (*HiveSub) ReadFile

```go
func (me *HiveSub) ReadFile(relPath ...string) (data []byte, err error)
```
Reads the file at the specified `{hive}/{sub}`-relative path.

#### type HiveSubs

```go
type HiveSubs struct {

	//	{hive}/dist/
	Dist HiveSub

	//	{hive}/cust/
	Cust HiveSub
}
```

Only used for `Hive.Subs`.

#### func (*HiveSubs) FileExists

```go
func (me *HiveSubs) FileExists(subRelPath ...string) bool
```
Returns whether `me.Dist` or `me.Cust` contains the specified file.

#### func (*HiveSubs) FilePath

```go
func (me *HiveSubs) FilePath(subRelPath ...string) (filePath string)
```
Returns either `me.Cust.FilePath(subRelPath...)` or
`me.Dist.FilePath(subRelPath...)`.

#### func (*HiveSubs) WalkAllDirs

```go
func (me *HiveSubs) WalkAllDirs(visitor ufs.WalkerVisitor, relPath ...string) (errs []error)
```
`ufs.WalkAllDirs` for `me.Dist` and `me.Cust`.

#### func (*HiveSubs) WalkAllFiles

```go
func (me *HiveSubs) WalkAllFiles(visitor ufs.WalkerVisitor, relPath ...string) (errs []error)
```
`ufs.WalkAllFiles` for `me.Dist` and `me.Cust`.

#### func (*HiveSubs) WalkDirsIn

```go
func (me *HiveSubs) WalkDirsIn(visitor ufs.WalkerVisitor, relPath ...string) (errs []error)
```
`ufs.WalkDirsIn` for `me.Dist` and `me.Cust`.

#### func (*HiveSubs) WalkFilesIn

```go
func (me *HiveSubs) WalkFilesIn(visitor ufs.WalkerVisitor, relPath ...string) (errs []error)
```
`ufs.WalkFilesIn` for `me.Dist` and `me.Cust`.

#### func (*HiveSubs) WatchIn

```go
func (me *HiveSubs) WatchIn(handler ufs.WatcherHandler, runHandlerNow bool, subRelPath ...string)
```
`ufs.DirWatcher.WatchIn` for `me.Dist` and `me.Cust`.

#### type Logger

```go
type Logger interface {
	// `Debugf` formats its arguments according to the `format`, analogous to `fmt.Printf`,
	// and records the text as a log message at Debug level.
	Debugf(format string, args ...interface{})

	// `Infof` is like `Debugf`, but at Info level.
	Infof(format string, args ...interface{})

	// `Warningf` is like `Debugf`, but at Warning level.
	Warningf(format string, args ...interface{})

	// `Error` records the specified `error` message at Error level,
	//	then should return the same specified `error` for more convenient in-place handling.
	Error(error) error

	// `Errorf` is like `Debugf`, but at Error level.
	Errorf(format string, args ...interface{})

	// `Criticalf` is like `Debugf`, but at Critical level.
	Criticalf(format string, args ...interface{})
}
```

An interface for log output. `ObLogger` provides the canonical implementation.

#### type ObLogger

```go
type ObLogger struct {
}
```

The canonical implementation of the `Logger` interface, using a standard
`log.Logger`.

#### func  NewLogger

```go
func NewLogger(out io.Writer) (me *ObLogger)
```
Creates and returns a new `*ObLogger`; `out` is optional and if `nil`, this
disables logging.

#### func (*ObLogger) Criticalf

```go
func (me *ObLogger) Criticalf(format string, args ...interface{})
```
Implements `Logger` interface.

#### func (*ObLogger) Debugf

```go
func (me *ObLogger) Debugf(format string, args ...interface{})
```
Implements `Logger` interface.

#### func (*ObLogger) Error

```go
func (me *ObLogger) Error(err error) error
```
Implements `Logger` interface.

#### func (*ObLogger) Errorf

```go
func (me *ObLogger) Errorf(format string, args ...interface{})
```
Implements `Logger` interface.

#### func (*ObLogger) Infof

```go
func (me *ObLogger) Infof(format string, args ...interface{})
```
Implements `Logger` interface.

#### func (*ObLogger) Warningf

```go
func (me *ObLogger) Warningf(format string, args ...interface{})
```
Implements `Logger` interface.

--
**godocdown** http://github.com/robertkrimen/godocdown
