package obwebui

import (
	"html/template"
	"path/filepath"
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

func getPageTemplate(subRelDirPath string) *PageTemplate {
	pageTemplateCache.Lock()
	defer pageTemplateCache.Unlock()
	pt, ok := pageTemplateCache.m[subRelDirPath]
	if !ok {
		pt = newPageTemplate(subRelDirPath)
		pt.load()
		pageTemplateCache.m[subRelDirPath] = pt
	}
	return pt
}

type PageTemplate struct {
	*template.Template
	subRelDirPath string
}

func newPageTemplate(subRelDirPath string) (me *PageTemplate) {
	me = &PageTemplate{subRelDirPath: subRelDirPath}
	return
}

func (me *PageTemplate) load() {
	loader := func(fullPath string) {
		dirPath := filepath.Dir(fullPath)
		fileNames := []string{filepath.Join(dirPath, "main.html")}
		uio.WalkFilesIn(dirPath, func(fullPath string) bool {
			if filepath.Base(fullPath) != "main.html" {
				fileNames = append(fileNames, fullPath)
			}
			return true
		})
		var err error
		me.Template, err = template.ParseFiles(fileNames...)
		if err != nil {
			me.Template, err = template.New(me.subRelDirPath).Parse(strf("ERROR loading template at '%s': %+v", dirPath, err))
		}
		return
	}
	ob.Hive.WatchDualDir(loader, true, me.subRelDirPath)
}
