package obpkg

import (
	"fmt"
	"io"

	usl "github.com/metaleap/go-util/slice"
)

type PkgCfgReloader func(pkg *Package)

type PkgCfg map[string]interface{}

var PkgCfgLoaders = map[string]PkgCfgReloader{}

type Packages []*Package

//	Implements sort.Interface.Len()
func (me Packages) Len() int { return len(me) }

//	Implements sort.Interface.Less()
func (me Packages) Less(i, j int) bool {
	p1, p2 := me[i], me[j]
	if usl.StrHas(p1.Info.Require, p2.NameFull) {
		return false
	}
	if usl.StrHas(p2.Info.Require, p1.NameFull) {
		return true
	}
	return me[i].NameFull < me[j].NameFull
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
