package obpkg

import (
	"fmt"
	"io"

	usl "github.com/metaleap/go-util/slice"
)

//	Used by Package.Kind-specific imports to register their reload handlers with PkgCfgLoaders.
type PkgCfgReloader func(pkg *Package)

//	Used in Package.CfgRaw
type PkgCfg map[string]interface{}

//	Contains one PkgCfgReloader handler per package kind.
//	When a Package gets (re)loaded, after populating its CfgRaw hash-maps,
//	it calls the appropriate PkgCfgReloader associated with its Kind to
//	notify it of its potentially new or changed PkgCfg settings.
var PkgCfgLoaders = map[string]PkgCfgReloader{}

//	A collection of *Package pointers
type Packages []*Package

//	Implements sort.Interface.Len()
func (me Packages) Len() int { return len(me) }

//	Implements sort.Interface.Less()
func (me Packages) Less(i, j int) bool {
	pi, pj := me[i], me[j]
	//	If i requires j, than j<i
	if usl.StrHas(pi.Info.Require, pj.NameFull) {
		return false
	}
	//	If j requires i, than i<j
	if usl.StrHas(pj.Info.Require, pi.NameFull) {
		return true
	}
	return pi.NameFull < pj.NameFull
}

//	Implements sort.Interface.Swap()
func (me Packages) Swap(i, j int) { me[i], me[j] = me[j], me[i] }

func errf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func outf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
