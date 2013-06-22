# obpkg_webuilib
--
    import "github.com/openbase/ob-core/bundle/webuilib"

Server-side web UI: 3rd-party JS/CSS libs

## Usage

#### type BundleCfg

```go
type BundleCfg struct {
	//	Extension-less, path-less CSS file names, from the "css" setting
	Css []string

	//	Extension-less, path-less JS file names, from the "js" setting
	Js []string

	//	All versions found as separate folders in the bundle directory,
	//	sorted descending-alphabetically ("from newest to oldest")
	Versions []string
}
```

Represents the Bundle.Cfg of a Bundle of Kind "webuilib"

#### func (*BundleCfg) Url

```go
func (me *BundleCfg) Url(fileBaseName, dotLessExtExt string) string
```
Returns a server-relative URL for the specified webuilib file. For example, for
a BundleCfg with name "bootstrap2", Url("bootstrap-responsive","css") returns:
/_pkg/webuilib-bootstrap2/{Versions[0]}/css/bootstrap-responsive.css

--
**godocdown** http://github.com/robertkrimen/godocdown