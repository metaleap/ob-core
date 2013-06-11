package obwebui

//	WebUI.Libs[..].CssUrl

type PageContext struct {
	WebUI WebUI
}

func (me *PageContext) Init() {
	me.WebUI.Libs = append(me.WebUI.Libs, WebUILib{CssUrls: []string{"/foo.css"}, JsUrls: []string{"/bar.js"}})
}
