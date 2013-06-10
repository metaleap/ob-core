package obcore

import (
	"sync"
)

var (
	pageTemplateCache struct {
		sync.Mutex
		m map[string]*PageTemplate
	}
)

func GetPageTemplate(name string) *PageTemplate {
	pageTemplateCache.Lock()
	defer pageTemplateCache.Unlock()
	pt, ok := pageTemplateCache.m[name]
	if !ok {
		// Hive.Watch.WatchDir(Hive.Subs.Cust.Path(...), runHandlerNow, handler)
		pt = newPageTemplate(name)
		pageTemplateCache.m[name] = pt
	}
	return pt
}

type PageTemplate struct {
}

func newPageTemplate(name string) (me *PageTemplate) {
	me = &PageTemplate{}
	return
}
