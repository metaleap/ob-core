# obsrv
--
    import "github.com/openbase/ob-core/server"

Web server functionality, used by `openbase/ob-gae` and
`openbase/ob-core/server/standalone`.

## Usage

#### type Ctx

```go
type Ctx struct {
	ob.Ctx
}
```


#### func  NewCtx

```go
func NewCtx(hiveDir string, logger ob.Logger) (me *Ctx, err error)
```

#### type HttpHandler

```go
type HttpHandler struct {
	http.Handler
	*Ctx

	//	Custom event handlers
	On struct {
		//	Request-related event handlers
		Request struct {
			//	Event handlers to be invoked before
			//	serving a web request (except static files).
			PreServe RequestContextHandlers

			//	Event handlers to be invoked immediately after
			//	serving a web request (except static files).
			PostServe RequestContextHandlers
		}
	}
}
```

Must be initialized via `NewHttpHandler`.

#### func  NewHttpHandler

```go
func NewHttpHandler(ctx *Ctx) (router *HttpHandler)
```
Initializes a new `*HttpHandler` to host the specified `*Ctx`.

#### type PageContext

```go
type PageContext struct {
	WebUI struct {
		Libs []*obpkg_webuilib.BundleCfg
		Skin *PageTemplate
	}
}
```

Created during server-side rendering of a `PageTemplate`.

#### type PageTemplate

```go
type PageTemplate struct {
}
```


#### func (*PageTemplate) Execute

```go
func (me *PageTemplate) Execute(w io.Writer, rc *RequestContext) error
```

#### type RequestContext

```go
type RequestContext struct {
	Ctx *Ctx

	//	Context related to the current `Page`, if any.
	Page *PageContext

	//	The `http.ResponseWriter` for this `RequestContext`.
	Out http.ResponseWriter

	//	The `http.Request` for this `RequestContext`.
	Req *http.Request

	//	Defaults to `ob.Log`.
	Log ob.Logger
}
```

Provides context for a non-static web request.

#### func (*RequestContext) Get

```go
func (me *RequestContext) Get(key interface{}) interface{}
```
http://www.gorillatoolkit.org/pkg/context#Get

#### func (*RequestContext) Set

```go
func (me *RequestContext) Set(key, val interface{})
```
http://www.gorillatoolkit.org/pkg/context#Set

#### type RequestContextHandler

```go
type RequestContextHandler func(*RequestContext)
```

A function that accepts a `*RequestContext`.

#### type RequestContextHandlers

```go
type RequestContextHandlers []RequestContextHandler
```

A collection of `RequestContextHandler` functions.

#### func (*RequestContextHandlers) Add

```go
func (me *RequestContextHandlers) Add(handlers ...RequestContextHandler)
```
Appends all the specified `handlers` to `me`.

--
**godocdown** http://github.com/robertkrimen/godocdown
