package obwebui

import (
	"html/template"
	"path/filepath"

	"github.com/go-utils/ufs"
	"github.com/go-utils/ugo"

	ob "github.com/openbase/ob-core"
)

var (
	pageTemplateCache struct {
		ugo.MutexIf
		m map[string]*PageTemplate
	}
)

func init() {
	pageTemplateCache.m = map[string]*PageTemplate{}
}

func getPageTemplate(ctx *ob.Ctx, subRelDirPath string) *PageTemplate {
	pageTemplateCache.Lock()
	pt, ok := pageTemplateCache.m[subRelDirPath]
	if !ok {
		pt = newPageTemplate(ctx, subRelDirPath)
		pageTemplateCache.m[subRelDirPath] = pt
	}
	pageTemplateCache.Unlock()
	if !ok {
		pt.load()
	}
	return pt
}

type PageTemplate struct {
	ugo.MutexIf
	*ob.Ctx
	*template.Template
	subRelDirPath string
}

func newPageTemplate(ctx *ob.Ctx, subRelDirPath string) (me *PageTemplate) {
	me = &PageTemplate{Ctx: ctx, subRelDirPath: subRelDirPath}
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
		defer me.UnlockIf(me.Lock())
		if me.Template, err = template.ParseFiles(fileNames...); err != nil {
			me.Template, err = template.New(me.subRelDirPath).Parse(strf("ERROR loading template at '%s': %+v", dirPath, err))
		}
		return
	}
	me.Ctx.Hive.Subs.WatchIn(loader, true, me.subRelDirPath)
}
