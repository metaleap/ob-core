package obpkg

import (
	"fmt"
	"io"

	"github.com/go-utils/uslice"
)

//	Used by Bundle.Kind-specific imports to register their reload handlers with BundleCfgLoaders.
type BundleCfgReloader func(*Bundle)

//	Used in Bundle.CfgRaw
type BundleCfg map[string]interface{}

//	Contains one BundleCfgReloader handler per bundle Kind.
//	When a Bundle gets (re)loaded, after populating its CfgRaw hash-maps,
//	it calls the appropriate BundleCfgReloader associated with its Kind to
//	notify it of its potentially new or changed BundleCfg settings.
var BundleCfgLoaders = map[string]BundleCfgReloader{}

//	A collection of *Bundle pointers
type Bundles []*Bundle

//	Implements sort.Interface.Len()
func (me Bundles) Len() int { return len(me) }

//	Implements sort.Interface.Less()
func (me Bundles) Less(i, j int) bool {
	pi, pj := me[i], me[j]
	//	If i requires j, than j<i
	if uslice.StrHas(pi.Info.Require, pj.NameFull) {
		return false
	}
	//	If j requires i, than i<j
	if uslice.StrHas(pj.Info.Require, pi.NameFull) {
		return true
	}
	return pi.NameFull < pj.NameFull
}

//	Implements sort.Interface.Swap()
func (me Bundles) Swap(i, j int) { me[i], me[j] = me[j], me[i] }

func errf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func outf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
