package obwebui

import (
	"fmt"
	"io"

	obpkg_webuilib "github.com/openbase/ob-core/pkg/webuilib"
)

type WebUi struct {
	Libs []*obpkg_webuilib.PkgCfg
}

func errf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func outf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
