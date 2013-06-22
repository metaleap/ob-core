# obwebui
--
    import "github.com/openbase/ob-core/webui"

Server-side web UI

## Usage

#### type PageContext

```go
type PageContext struct {
	WebUi PageContextWebUi
}
```


#### func (*PageContext) Init

```go
func (me *PageContext) Init()
```

#### type PageContextWebUi

```go
type PageContextWebUi struct {
	Libs         []*obpkg_webuilib.BundleCfg
	SkinTemplate *PageTemplate
}
```


#### type PageTemplate

```go
type PageTemplate struct {
	*template.Template
}
```

--
**godocdown** http://github.com/robertkrimen/godocdown