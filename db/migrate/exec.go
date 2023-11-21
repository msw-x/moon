package migrate

import (
	"context"
	"strings"

	"github.com/uptrace/bun"
)

func Exec(db *bun.DB, sql string, isTx bool, splitter string) (err error) {
	if sql == "" {
		return
	}
	ctx := context.Background()
	var queries []string
	if splitter == "" {
		queries = append(queries, sql)
	} else {
		for _, q := range strings.Split(sql, splitter) {
			queries = append(queries, q)
		}
	}
	var idb bun.IConn
	if isTx {
		var tx bun.Tx
		tx, err = db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		idb = tx
	} else {
		var conn bun.Conn
		conn, err = db.Conn(ctx)
		if err != nil {
			return err
		}
		idb = conn
	}
	for _, q := range queries {
		_, err = idb.ExecContext(ctx, q)
		if err != nil {
			break
		}
	}
	switch v := idb.(type) {
	case bun.Tx:
		if err == nil {
			err = v.Commit()
		} else {
			v.Rollback()
		}
	case bun.Conn:
		if err == nil {
			err = v.Close()
		} else {
			v.Close()
		}
	}
	return
}
