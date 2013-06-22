# obsrv
--
    import "github.com/openbase/ob-core/server"

Web server functionality, used by the cmd/ob-server main package

## Usage

```go
var (
	//	Multi-plexing request router
	Router *webmux.Router

	//	Custom event handlers
	On struct {
		//	Request-related event handlers
		Request struct {
			//	Event handlers to be invoked before serving a web request (except static files)
			Serving RequestContextEventHandlers

			//	Event handlers to be invoked immediately after serving a web request (except static files)
			Served RequestContextEventHandlers
		}
	}
)
```

#### func  Init

```go
func Init()
```
Initializes the package for serving web requests. To be called after ob.Init()

#### type RequestContext

```go
type RequestContext struct {
	//	Context related to the current Page, if any.
	obwebui.PageContext

	//	The http.ResponseWriter for this RequestContext
	Out http.ResponseWriter

	//	The http.Request for this RequestContext
	Req *http.Request

	//	Defaults to ob.Log
	Log ob.Logger

	//	Not used in the default stand-alone implementation (cmd/ob-server).
	//	May be used in sandboxed mode (eg. the GAE package uses it for the current appengine.Context)
	Ctx interface{}
}
```

Encapsulates and provides context for a (non-static) web request

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

#### type RequestContextEventHandler

```go
type RequestContextEventHandler func(*RequestContext)
```

A function that accepts a *RequestContext

#### type RequestContextEventHandlers

```go
type RequestContextEventHandlers []RequestContextEventHandler
```

A collection of RequestContextEventHandler function-values

#### func (*RequestContextEventHandlers) Add

```go
func (me *RequestContextEventHandlers) Add(eventHandlers ...RequestContextEventHandler)
```
Adds the specified eventHandlers to me

--
**godocdown** http://github.com/robertkrimen/godocdown