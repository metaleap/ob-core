# obwebui
--
    import "github.com/openbase/ob-core/webui"

Server-side web UI

## Usage

#### type PageContext

```go
type PageContext struct {
	WebUi WebUi
}
```


#### func  NewPageContext

```go
func NewPageContext(ctx *ob.Ctx) (me *PageContext)
```

#### type PageTemplate

```go
type PageTemplate struct {
	ugo.MutexIf
	*ob.Ctx
	*template.Template
}
```


#### type WebUi

```go
type WebUi struct {
	Libs         []*obpkg_webuilib.BundleCfg
	SkinTemplate *PageTemplate
}
```

--
**godocdown** http://github.com/robertkrimen/godocdown
