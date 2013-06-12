package obwebui

import (
	"html/template"
	"path/filepath"
	"strings"
	"sync"

	uio "github.com/metaleap/go-util/io"

	ob "github.com/openbase/ob-core"
)

var (
	pageTemplateCache struct {
		sync.Mutex
		m map[string]*PageTemplate
	}
)

func init() {
	pageTemplateCache.m = map[string]*PageTemplate{}
}

func GetPageTemplate(name string) *PageTemplate {
	pageTemplateCache.Lock()
	defer pageTemplateCache.Unlock()
	pt, ok := pageTemplateCache.m[name]
	if !ok {
		pt = newPageTemplate(name)
		pt.load()
		pageTemplateCache.m[name] = pt
	}
	return pt
}

type PageTemplate struct {
	*template.Template
	name string
}

func newPageTemplate(name string) (me *PageTemplate) {
	me = &PageTemplate{name: name}
	return
}

func (me *PageTemplate) load() {
	//	currently NOT proper "dual-dir" handling!
	loader := func(dirPath string) {
		if strings.Contains(dirPath, "cust") {
			return
		}
		fileNames := []string{filepath.Join(dirPath, me.name+".html")}
		uio.WalkFilesIn(dirPath, func(fullPath string) bool {
			if !strings.HasSuffix(fullPath, string(filepath.Separator)+me.name+".html") {
				fileNames = append(fileNames, fullPath)
			}
			return true
		})
		var err error
		me.Template, err = template.ParseFiles(fileNames...)
		if err != nil {
			me.Template, err = template.New(me.name).Parse(strf("ERROR loading templates at '%s': %+v", dirPath, err))
		}
		return
	}
	ob.Hive.WatchDualDir(loader, "client", "tmpl", me.name)
}
