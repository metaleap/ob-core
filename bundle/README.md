# obpkg
--
    import "github.com/openbase/ob-core/bundle"

Bundle management: loads .ob-pkg files inside {hive}/{sub}/pkg/{kind}-{name}/ directories.

## Usage

```go
var BundleCfgLoaders = map[string]BundleCfgReloader{}
```
Contains one BundleCfgReloader handler per bundle Kind. When a Bundle gets
(re)loaded, after populating its CfgRaw hash-maps, it calls the appropriate
BundleCfgReloader associated with its Kind to notify it of its potentially new
or changed BundleCfg settings.

#### type Bundle

```go
type Bundle struct {
	//	The kind of bundle, according to its directory name, for example "webuilib"
	Kind string

	//	The name of this bundle, not including its Kind
	Name string

	//	The full identifier of this bundle, which is Kind and Name joined by a dash, for example "webuilib-jquery"
	NameFull string

	//	Diagnostic info
	Diag struct {
		//	Full bundle names of all Info.Require entries that are not currently installed inside {hive}/{sub}/pkg/
		BadDeps []string

		//	The error that occurred when loading the .ob-pkg file, if any.
		//	Outside of unlikely hard-disk crashes, this is most likely a TOML syntax error in the file.
		LoadErr error
	}

	//	Information from the '[bundle]' section of the .ob-pkg bundle configuration file.
	Info struct {
		//	Human-readable bundle title
		Title string

		//	Human-readable, comprehensive bundle description
		Desc string

		//	Web address for more information, in the case of 3rd-party bundles
		Www string

		//	Denotes (fully-qualified) bundles required for this bundle to function
		Require []string
	}

	//	A value or struct that represents CfgRaw in a 'native', bundle-specific way.
	//	To be set by a BundleCfgReloader registered in BundleCfgLoaders.
	Cfg interface{}

	//	Information from the '[default]' and other sections of the .ob-pkg bundle configuration file.
	CfgRaw struct {
		//	Information from the '[default]' section of the .ob-pkg bundle configuration file.
		Default BundleCfg

		//	Information from any other sections of the .ob-pkg bundle configuration file.
		More map[string]BundleCfg
	}
}
```

Represents a bundle found in a {hive}/{sub}/pkg/{kind-name}/{name.kind.ob-pkg}
file

#### type BundleCfg

```go
type BundleCfg map[string]interface{}
```

Used in Bundle.CfgRaw

#### type BundleCfgReloader

```go
type BundleCfgReloader func(*Bundle)
```

Used by Bundle.Kind-specific imports to register their reload handlers with
BundleCfgLoaders.

#### type Bundles

```go
type Bundles []*Bundle
```

A collection of *Bundle pointers

#### func (Bundles) Len

```go
func (me Bundles) Len() int
```
Implements sort.Interface.Len()

#### func (Bundles) Less

```go
func (me Bundles) Less(i, j int) bool
```
Implements sort.Interface.Less()

#### func (Bundles) Swap

```go
func (me Bundles) Swap(i, j int)
```
Implements sort.Interface.Swap()

#### type Registry

```go
type Registry struct {
}
```

Bundle package registry, used for Reg

```go
var (
	//	The global bundle package registry
	Reg Registry
)
```

#### func (*Registry) ByKind

```go
func (me *Registry) ByKind(kind string, deps []string) (all Bundles)
```
If deps is empty, returns all Bundles of the specified kind. Otherwise, out of
all Bundles specified in deps or directly or indirectly required by them,
returns those of the specified kind.

#### func (*Registry) ByName

```go
func (me *Registry) ByName(kindAndName ...string) *Bundle
```
Returns the *Bundle with the specified fully-qualified identifier. kindAndName
can be a single string such as "webuilib-jquery", or 2 strings for kind and
name, such as "webuilib" and "jquery".

#### func (*Registry) Init

```go
func (me *Registry) Init()
```
Initializes a few internal hash-maps to non-nil, NO need to call this for Reg.
Does not load the bundles just yet, this is done via any of the ByXYZ() methods.

--
**godocdown** http://github.com/robertkrimen/godocdown