// Server-side web UI
package obwebui

import (
	"fmt"
)

func errf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
