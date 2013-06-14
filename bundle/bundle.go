package obpkg

import (
	"github.com/goforks/toml"

	usl "github.com/metaleap/go-util/slice"

	ob "github.com/openbase/ob-core"
)

//	Represents a bundle found in a {hive}/{sub}/pkg/{kind-name}/{name.kind.ob-pkg} file
type Bundle struct {
	//	The kind of bundle, according to its directory name, for example "webuilib"
	Kind string

	//	The name of this bundle, not including its Kind
	Name string

	//	The full identifier of this bundle, which is Kind and Name joined by a dash, for example "webuilib-jquery"
	NameFull string

	//	Diagnostic info
	Diag struct {
		//	Full bundle names of all Info.Require entries that are not currently installed inside {hive}/{sub}/pkg/
		BadDeps []string

		//	The error that occurred when loading the .ob-pkg file, if any.
		//	Outside of unlikely hard-disk crashes, this is most likely a TOML syntax error in the file.
		LoadErr error
	}

	//	Information from the '[bundle]' section of the .ob-pkg bundle configuration file.
	Info struct {
		//	Human-readable bundle title
		Title string

		//	Human-readable, comprehensive bundle description
		Desc string

		//	Web address for more information, in the case of 3rd-party bundles
		Www string

		//	Denotes (fully-qualified) bundles required for this bundle to function
		Require []string
	}

	//	A value or struct that represents CfgRaw in a 'native', bundle-specific way.
	//	To be set by a BundleCfgReloader registered in BundleCfgLoaders.
	Cfg interface{}

	//	Information from the '[default]' and other sections of the .ob-pkg bundle configuration file.
	CfgRaw struct {
		//	Information from the '[default]' section of the .ob-pkg bundle configuration file.
		Default BundleCfg

		//	Information from any other sections of the .ob-pkg bundle configuration file.
		More map[string]BundleCfg
	}
}

func newBundle() (me *Bundle) {
	me = &Bundle{}
	me.CfgRaw.Default, me.CfgRaw.More = BundleCfg{}, map[string]BundleCfg{}
	return
}

//	This may load from the primary dist .ob-pkg file, or just partial additions/overrides from cust.
//	But it's only called if .ob-pkg file (filePath) does exist.
func (me *Bundle) reload(kind, name, fullName, filePath string) {
	me.Kind, me.Name, me.NameFull = kind, name, fullName
	config := map[string]interface{}{}
	str := func(m map[string]interface{}, name string) (s string) {
		s, _ = m[name].(string)
		return
	}
	if _, me.Diag.LoadErr = toml.DecodeFile(filePath, config); me.Diag.LoadErr != nil {
		ob.Log.Errorf("[BUNDLE] %s", me.Diag.LoadErr.Error())
	} else {
		var (
			ok  bool
			cfg map[string]interface{}
			key string
			val interface{}
		)
		if cfg, ok = config["bundle"].(map[string]interface{}); ok {
			me.Info.Title, me.Info.Desc, me.Info.Www = str(cfg, "title"), str(cfg, "desc"), str(cfg, "www")
			if req, _ := cfg["require"].([]interface{}); len(req) > 0 {
				usl.StrAppendUniques(&me.Info.Require, usl.StrConvert(req, true)...)
			}
		}
		if cfg, ok = config["default"].(map[string]interface{}); ok {
			for key, val = range cfg {
				me.CfgRaw.Default[key] = val
			}
		}
		for key, val = range config {
			if key != "bundle" && key != "default" {
				// println("MORE:" + key)
			}
		}
		if loader := BundleCfgLoaders[kind]; loader != nil {
			loader(me)
		}
	}
	if len(me.Info.Title) == 0 {
		me.Info.Title = me.NameFull
	}
}
