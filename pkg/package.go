package obpkg

import (
	"github.com/goforks/toml"

	ob "github.com/openbase/ob-core"

	usl "github.com/metaleap/go-util/slice"
)

type PkgCfg map[string]interface{}

type Packages []*Package

type Package struct {
	Info struct {
		Title   string
		Desc    string
		Www     string
		Require []string
	}

	Config struct {
		Default PkgCfg
		More    map[string]PkgCfg
	}

	Error error
	// IsCust, HasCust      bool
	Kind, Name, NameFull string
}

func newPackage() (me *Package) {
	me = &Package{}
	me.Config.Default, me.Config.More = PkgCfg{}, map[string]PkgCfg{}
	return
}

//	This may load from the primary dist .obpkg file, or just partial additions/overrides from cust
func (me *Package) reload(kind, name, fullName, filePath string) {
	me.Kind, me.Name, me.NameFull = kind, name, fullName
	cfg := map[string]interface{}{}
	s := func(m map[string]interface{}, name string) (s string) {
		s, _ = m[name].(string)
		return
	}
	if _, me.Error = toml.DecodeFile(filePath, cfg); me.Error != nil {
		ob.Opt.Log.Error(me.Error)
	} else {
		var (
			ok                 bool
			cfgPkg, cfgDefault map[string]interface{}
			key                string
			val                interface{}
		)
		if cfgPkg, ok = cfg["pkg"].(map[string]interface{}); ok {
			me.Info.Title, me.Info.Desc, me.Info.Www = s(cfgPkg, "title"), s(cfgPkg, "desc"), s(cfgPkg, "www")
			if req, _ := cfgPkg["require"].([]string); len(req) > 0 {
				usl.StrAppendUniques(&me.Info.Require, req...)
			}
		}
		if cfgDefault, ok = cfg["default"].(map[string]interface{}); ok {
			for key, val = range cfgDefault {
				me.Config.Default[key] = val
			}
		}
		for key, val = range cfg {
			if key != "pkg" && key != "default" {
				println("MORE:" + key)
			}
		}
	}
	if len(me.Info.Title) == 0 {
		me.Info.Title = me.NameFull
	}
}
