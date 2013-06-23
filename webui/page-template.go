package obwebui

import (
	"html/template"
	"path/filepath"
	"sync"

	"github.com/go-utils/ufs"

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

func getPageTemplate(ctx *ob.Ctx, subRelDirPath string) *PageTemplate {
	pageTemplateCache.Lock()
	defer pageTemplateCache.Unlock()
	pt, ok := pageTemplateCache.m[subRelDirPath]
	if !ok {
		pt = newPageTemplate(ctx, subRelDirPath)
		pt.load()
		pageTemplateCache.m[subRelDirPath] = pt
	}
	return pt
}

type PageTemplate struct {
	ctx *ob.Ctx
	*template.Template
	subRelDirPath string
}

func newPageTemplate(ctx *ob.Ctx, subRelDirPath string) (me *PageTemplate) {
	me = &PageTemplate{ctx: ctx, subRelDirPath: subRelDirPath}
	return
}

func (me *PageTemplate) load() {
	loader := func(fullPath string) {
		dirPath := filepath.Dir(fullPath)
		fileNames := []string{filepath.Join(dirPath, "main.ob-tmpl")}
		ufs.WalkFilesIn(dirPath, func(fullPath string) bool {
			if filepath.Base(fullPath) != "main.ob-tmpl" {
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
	me.ctx.Hive.Subs.WatchIn(loader, true, me.subRelDirPath)
}
