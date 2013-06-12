package obpkg

import (
	"path/filepath"

	"github.com/goforks/toml"

	usl "github.com/metaleap/go-util/slice"

	ob "github.com/openbase/ob-core"
)

type Package struct {
	Info struct {
		Title   string
		Desc    string
		Www     string
		Require []string
	}

	Cfg interface{}

	CfgRaw struct {
		Default PkgCfg
		More    map[string]PkgCfg
	}

	Diag struct {
		BadDeps []string
		LoadErr error
	}

	Dir, Kind, Name, NameFull string
}

func newPackage() (me *Package) {
	me = &Package{}
	me.CfgRaw.Default, me.CfgRaw.More = PkgCfg{}, map[string]PkgCfg{}
	return
}

//	This may load from the primary dist .ob-pkg file, or just partial additions/overrides from cust
func (me *Package) reload(kind, name, fullName, filePath string) {
	me.Dir, me.Kind, me.Name, me.NameFull = filepath.Dir(filePath), kind, name, fullName
	cfg := map[string]interface{}{}
	s := func(m map[string]interface{}, name string) (s string) {
		s, _ = m[name].(string)
		return
	}
	if _, me.Diag.LoadErr = toml.DecodeFile(filePath, cfg); me.Diag.LoadErr != nil {
		ob.Opt.Log.Errorf("[PKG] %s", me.Diag.LoadErr.Error())
	} else {
		var (
			ok                 bool
			cfgPkg, cfgDefault map[string]interface{}
			key                string
			val                interface{}
		)
		if cfgPkg, ok = cfg["pkg"].(map[string]interface{}); ok {
			me.Info.Title, me.Info.Desc, me.Info.Www = s(cfgPkg, "title"), s(cfgPkg, "desc"), s(cfgPkg, "www")
			if req, _ := cfgPkg["require"].([]interface{}); len(req) > 0 {
				usl.StrAppendUniques(&me.Info.Require, usl.StrConvert(req, true)...)
			}
		}
		if cfgDefault, ok = cfg["default"].(map[string]interface{}); ok {
			for key, val = range cfgDefault {
				me.CfgRaw.Default[key] = val
			}
		}
		for key, val = range cfg {
			if key != "pkg" && key != "default" {
				// println("MORE:" + key)
			}
		}
		if loader := PkgCfgLoaders[kind]; loader != nil {
			loader(me)
		}
	}
	if len(me.Info.Title) == 0 {
		me.Info.Title = me.NameFull
	}
}
