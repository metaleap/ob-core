# obpkg_webuiskin
--
    import "github.com/openbase/ob-core/bundle/webuiskin"

Server-side web UI: "skins" (main HTML templates)

## Usage

#### type BundleCfg

```go
type BundleCfg struct {
	//	A HiveSub-relative directory path for the template files of this webuiskin.
	//	For example, for a BundleCfg with Name "fluid", this would be:
	//	pkg/webuiskin-fluid/template
	SubRelTemplateDirPath string
}
```

Represents the Bundle.Cfg of a Bundle of Kind "webuiskin"

--
**godocdown** http://github.com/robertkrimen/godocdown
