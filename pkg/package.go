package obpkg

import (
	"github.com/goforks/toml"

	usl "github.com/metaleap/go-util/slice"

	ob "github.com/openbase/ob-core"
)

//	Represents a package found in a {hive}/{sub}/pkg/{kind-name}/{name.kind.ob-pkg} file
type Package struct {
	//	The kind of package, according to its directory name, for example "webuilib"
	Kind string

	//	The name of this package, not including its Kind
	Name string

	//	The full identifier of this package, which is Kind and Name joined by a dash, for example "webuilib-jquery"
	NameFull string

	//	Diagnostic info
	Diag struct {
		//	Full package names of all Info.Require entries that are not currently installed inside {hive}/{sub}/pkg
		BadDeps []string

		//	The error that occurred when loading the .ob-pkg file, if any.
		//	Outside of unlikely hard-disk crashes, this is most likely a TOML syntax error in the file.
		LoadErr error
	}

	//	Information from the '[pkg]' section of the .ob-pkg package configuration file.
	Info struct {
		//	Human-readable package title
		Title string

		//	Human-readable, comprehensive package description
		Desc string

		//	Web address for more information, in the case of 3rd-party packages
		Www string

		//	Denotes (fully-qualified) packages required for this package to function
		Require []string
	}

	//	A value or struct that represents CfgRaw in a 'native', package-specific way.
	//	To be set by a PkgCfgReloader registered in PkgCfgLoaders.
	Cfg interface{}

	//	Information from the '[default]' and other sections of the .ob-pkg package configuration file.
	CfgRaw struct {
		//	Information from the '[default]' section of the .ob-pkg package configuration file.
		Default PkgCfg

		//	Information from any other sections of the .ob-pkg package configuration file.
		More map[string]PkgCfg
	}
}

func newPackage() (me *Package) {
	me = &Package{}
	me.CfgRaw.Default, me.CfgRaw.More = PkgCfg{}, map[string]PkgCfg{}
	return
}

//	This may load from the primary dist .ob-pkg file, or just partial additions/overrides from cust.
//	But it's only called if .ob-pkg file (filePath) does exist.
func (me *Package) reload(kind, name, fullName, filePath string) {
	me.Kind, me.Name, me.NameFull = kind, name, fullName
	cfg := map[string]interface{}{}
	str := func(m map[string]interface{}, name string) (s string) {
		s, _ = m[name].(string)
		return
	}
	if _, me.Diag.LoadErr = toml.DecodeFile(filePath, cfg); me.Diag.LoadErr != nil {
		ob.Log.Errorf("[PKG] %s", me.Diag.LoadErr.Error())
	} else {
		var (
			ok                 bool
			cfgPkg, cfgDefault map[string]interface{}
			key                string
			val                interface{}
		)
		if cfgPkg, ok = cfg["pkg"].(map[string]interface{}); ok {
			me.Info.Title, me.Info.Desc, me.Info.Www = str(cfgPkg, "title"), str(cfgPkg, "desc"), str(cfgPkg, "www")
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
