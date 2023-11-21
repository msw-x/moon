package migrate

import (
	"io"
	"io/fs"
	"strings"
)

func ReadSql(f fs.FS, path string) (sql string, err error) {
	var r io.Reader
	r, err = f.Open(path)
	if err != nil {
		return
	}
	b := new(strings.Builder)
	_, err = io.Copy(b, r)
	if err != nil {
		return
	}
	sql = b.String()
	return
}
