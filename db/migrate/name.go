package migrate

import (
	"path"
	"strings"

	"github.com/msw-x/moon/ustring"
)

func Name(s string) (name string, comment string, isUp bool, isTx bool, ok bool) {
	const upSuffix = ".up.sql"
	const downSuffix = ".down.sql"
	const txSuffix = ".tx"
	s = path.Base(s)
	ok = strings.HasSuffix(s, upSuffix)
	if ok {
		isUp = true
		s = strings.TrimSuffix(s, upSuffix)
	} else {
		ok = strings.HasSuffix(s, downSuffix)
		if ok {
			s = strings.TrimSuffix(s, downSuffix)
		}
	}
	if ok {
		isTx = strings.HasSuffix(s, txSuffix)
		if isTx {
			s = strings.TrimSuffix(s, txSuffix)
		}
		name, comment = ustring.SplitPair(s, "_")
	}
	return
}
