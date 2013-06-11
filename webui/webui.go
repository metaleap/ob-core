package obwebui

import (
	"fmt"
	"io"
)

type WebUI struct {
	Libs []WebUILib
}

type WebUILib struct {
	CssUrls, JsUrls []string
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
