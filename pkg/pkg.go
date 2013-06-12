package obpkg

import (
	"fmt"
	"io"
)

type PkgCfgReloader func(pkg *Package)

var PkgCfgLoaders = map[string]PkgCfgReloader{}

func errf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func outf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
