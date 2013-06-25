// Web server functionality, used by `openbase/ob-gae` and `openbase/ob-core/server/standalone`.
package obsrv

import (
	"fmt"
)

func strf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
