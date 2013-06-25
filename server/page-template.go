package obsrv

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/go-utils/ufs"
	"github.com/go-utils/ugo"
)

func (me *Ctx) getPageTemplate(subRelDirPath string) *PageTemplate {
	defer me.pageTemplateCache.mx.UnlockIf(me.pageTemplateCache.mx.Lock())
	pt, ok := me.pageTemplateCache.m[subRelDirPath]
	if !ok {
		pt = newPageTemplate(subRelDirPath)
		me.pageTemplateCache.m[subRelDirPath] = pt
		pt.load(me)
	}
	return pt
}

//	Wraps a `html/template.Template` defined in a `webuiskin` `Kind` of `Bundle`.
type PageTemplate struct {
	mx            ugo.MutexIf
	tmpl          *template.Template
	subRelDirPath string
}

func newPageTemplate(subRelDirPath string) (me *PageTemplate) {
	me = &PageTemplate{subRelDirPath: subRelDirPath}
	return
}

func (me *PageTemplate) exec(w io.Writer, rc *RequestContext) error {
	me.mx.Lock()
	tmpl := me.tmpl
	me.mx.Unlock()
	return tmpl.Execute(w, rc)
}

func (me *PageTemplate) load(ctx *Ctx) {
	const mainFileName = "main.ob-tmpl"
	loader := func(fullPath string) {
		dirPath := filepath.Dir(fullPath)
		fileNames := []string{filepath.Join(dirPath, mainFileName)}
		ufs.WalkFilesIn(dirPath, func(fullPath string) bool {
			if filepath.Base(fullPath) != mainFileName {
				fileNames = append(fileNames, fullPath)
			}
			return true
		})
		tmpl, err := template.ParseFiles(fileNames...)
		if err != nil {
			tmpl, err = template.New(me.subRelDirPath).Parse(strf("ERROR loading template at '%s': %+v", dirPath, err))
		}
		defer me.mx.UnlockIf(me.mx.Lock())
		me.tmpl = tmpl
	}
	loader(ctx.Hive.Subs.Path(me.subRelDirPath, mainFileName))
	ctx.Hive.Subs.WatchIn(loader, false, me.subRelDirPath)
}
