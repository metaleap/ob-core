package obcore

import (
	"github.com/go-forks/toml"

	"github.com/go-utils/uslice"
)

//	Represents a bundle package found in a `{hive}/{sub}/pkg/{kind-name}/{name.kind.ob-pkg}` file.
type Bundle struct {
	Ctx *Ctx

	//	The kind of this `Bundle`, according to its directory name,
	//	for example `webuilib`.
	Kind string

	//	The name of this `Bundle`, not including its `Kind`.
	Name string

	//	The full identifier of this `Bundle`, which is `Kind` and `Name`
	//	joined by a dash, for example `webuilib-jquery`.
	NameFull string

	//	Diagnostic info
	Diag struct {
		//	Full `Bundle` names of all `Info.Require` entries that
		//	are not currently installed inside `{hive}/{sub}/pkg/`.
		BadDeps []string

		//	The `error` that occurred when loading the `.ob-pkg` file, if any.
		//	Could be e.g. a file-system I/O issue or a TOML syntax error.
		LoadErr error
	}

	//	Information from the `[bundle]` section of the `.ob-pkg` bundle configuration file.
	Info struct {
		//	Human-readable bundle title
		Title string

		//	Human-readable, comprehensive bundle description
		Desc string

		//	Web address for more information, in the case of 3rd-party bundles
		Www string

		//	Names fully-qualified bundles required for this bundle to function
		Require []string
	}

	//	Unprocessed information loaded from the `.ob-pkg` bundle configuration file.
	CfgRaw struct {
		//	Information from the `[default]` section.
		Default BundleCfg

		//	Information from any other sections.
		More map[string]BundleCfg
	}

	//	Represents `CfgRaw` in a 'native'/'parsed'/'processed', `Kind`-specific way.
	//	To be set by a `BundleCfgReloader` registered in `BundleCfgLoaders`, if any.
	Cfg interface{}
}

func newBundle(reg *BundleRegistry) (me *Bundle) {
	me = &Bundle{Ctx: reg.Ctx}
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
		me.Ctx.Log.Errorf("[BUNDLE] %s", me.Diag.LoadErr.Error())
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
				uslice.StrAppendUniques(&me.Info.Require, uslice.StrConvert(req, true)...)
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
		if loader := BundleCfgReloaders[kind]; loader != nil {
			loader(me)
		}
	}
	if len(me.Info.Title) == 0 {
		me.Info.Title = me.NameFull
	}
}
