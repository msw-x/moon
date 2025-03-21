package migrate

import (
	"context"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"time"

	"github.com/msw-x/moon/tabtable"
	"github.com/uptrace/bun"
)

type Migrator struct {
	db                   *bun.DB
	tableName            string
	splitter             string
	saveDownSql          bool
	autoDownSql          bool
	autoDownSqlTrace     func(*Migration, *Context)
	rollbackLost         bool
	markAppliedOnSuccess bool
	ro                   bool
	l                    Migrations
}

func NewMigrator(db *bun.DB) *Migrator {
	o := new(Migrator)
	o.db = db
	o.splitter = "\n\n"
	o.tableName = "migrations"
	o.saveDownSql = true
	o.autoDownSql = true
	o.rollbackLost = true
	return o
}

func (o *Migrator) WithReadonly(ro bool) *Migrator {
	o.ro = ro
	return o
}

func (o *Migrator) WithTableName(v string) *Migrator {
	o.tableName = v
	return o
}

func (o *Migrator) WithSplitter(v string) *Migrator {
	o.splitter = v
	return o
}

func (o *Migrator) WithSaveDownSql(v bool) *Migrator {
	o.saveDownSql = v
	return o
}

func (o *Migrator) WithAutoDownSql(v bool) *Migrator {
	o.autoDownSql = v
	return o
}

func (o *Migrator) WithAutoDownSqlTrace(v func(*Migration, *Context)) *Migrator {
	o.autoDownSqlTrace = v
	return o
}

func (o *Migrator) WithRollbackLost(v bool) *Migrator {
	o.rollbackLost = v
	return o
}

func (o *Migrator) WithMarkAppliedOnSuccess(v bool) *Migrator {
	o.markAppliedOnSuccess = v
	return o
}

func (o *Migrator) Init() (err error) {
	if o.ro {
		return
	}
	_, err = o.db.NewCreateTable().
		Model((*Migration)(nil)).
		ModelTableExpr(o.tableName).
		IfNotExists().
		Exec(context.Background())
	return
}

func (o *Migrator) Reset() (err error) {
	_, err = o.db.NewDropTable().
		Model((*Migration)(nil)).
		ModelTableExpr(o.tableName).
		IfExists().
		Exec(context.Background())
	return
}

func (o *Migrator) Forget() (err error) {
	if o.ro {
		return
	}
	_, err = o.db.NewTruncateTable().TableExpr(o.tableName).Exec(context.Background())
	return
}

func (o *Migrator) Load(fs fs.FS, final string) (err error) {
	err = o.db.NewSelect().
		ColumnExpr("*").
		Model(&o.l).
		ModelTableExpr(o.tableName).
		Scan(context.Background())
	if err != nil {
		return
	}
	err = o.l.Load(fs, final)
	if err != nil {
		return
	}
	if o.autoDownSql {
		err = o.l.AutoGenerateDownTrace(o.autoDownSqlTrace)
	}
	if err != nil {
		return
	}
	o.l.SetFuncs(o.db, o.splitter)
	o.l.SortAsc()
	return
}

func (o *Migrator) Migrations() Migrations {
	return o.l
}

func (o *Migrator) AppliedMigrations() (l Migrations) {
	for _, m := range o.l {
		if m.IsApplied() {
			l = append(l, m)
		}
	}
	return
}

func (o *Migrator) UnappliedMigrations() (l Migrations) {
	for _, m := range o.l {
		if !m.IsApplied() && !m.Lost() {
			l = append(l, m)
		}
	}
	return
}

func (o *Migrator) RolledbackMigrations() (l Migrations) {
	if !o.saveDownSql || !o.rollbackLost {
		return
	}
	a := o.AppliedMigrations()
	slices.Reverse(a)
	for _, m := range a {
		if m.Lost() {
			l = append(l, m)
		} else {
			break
		}
	}
	slices.Reverse(l)
	return
}

func (o *Migrator) MigrationsCount() int {
	return len(o.Migrations())
}

func (o *Migrator) AppliedMigrationsCount() int {
	return len(o.AppliedMigrations())
}

func (o *Migrator) UnappliedMigrationsCount() int {
	return len(o.UnappliedMigrations())
}

func (o *Migrator) RolledbackMigrationsCount() int {
	return len(o.RolledbackMigrations())
}

func (o *Migrator) HasIncompleteMigrations() bool {
	return o.RolledbackMigrationsCount() != 0 || o.UnappliedMigrationsCount() != 0
}

func (o *Migrator) Migrate() (err error) {
	err = o.rollback(o.RolledbackMigrations())
	if err != nil {
		return
	}
	groupId := o.l.LastGroupId() + 1
	for _, m := range o.UnappliedMigrations() {
		err = m.Up()
		if err == nil || !o.markAppliedOnSuccess {
			m.SetApplied(groupId)
			o.markApplied(m)
		}
		if err != nil {
			break
		}
	}
	return
}

func (o *Migrator) Rollback() error {
	return o.rollback(o.l.LastGroup())
}

func (o *Migrator) RollbackLast() (err error) {
	if o.UnappliedMigrationsCount() == 0 {
		return o.Rollback()
	}
	return nil
}

func (o *Migrator) RepairDown() ([]string, error) {
	return o.l.RepairDown(o.updateApplied)
}

func (o *Migrator) ViewSchema() (string, error) {
	return o.l.ViewSchema()
}

func (o *Migrator) Status() string {
	w := tabtable.New()
	w.Write("name", "comment", "group", "migrated")
	for _, m := range o.l {
		var migratedAt string
		if !m.MigratedAt.IsZero() {
			migratedAt = m.MigratedAt.In(time.Now().Location()).Format("02-Jan-2006 15:04:05")
		}
		var notes []string
		note := func(s string) {
			notes = append(notes, s)
		}
		if m.IsTx {
			note("tx")
		}
		if m.Lost() {
			note("lost")
		}
		if m.Error() != nil {
			note(m.Error().Error())
		}
		var group string
		if m.GroupId > 0 {
			group = fmt.Sprint(m.GroupId)
		}
		w.Write(m.Name, m.Comment, group, migratedAt, strings.Join(notes, ", "))
	}
	s := fmt.Sprintf("migrations: applied[%d/%d]", o.AppliedMigrationsCount(), o.MigrationsCount())
	if o.RolledbackMigrationsCount() > 0 {
		s = fmt.Sprintf("%s rolledback[%d]", s, o.RolledbackMigrationsCount())
	}
	if o.UnappliedMigrationsCount() > 0 {
		s = fmt.Sprintf("%s unapplied[%d]", s, o.UnappliedMigrationsCount())
	}
	if len(o.l) > 0 {
		s = s + "\n" + w.String()
	}
	return s
}

func (o *Migrator) writeApplied(m *Migration, update bool) (err error) {
	if o.ro {
		return
	}
	var downSql string
	if !o.saveDownSql {
		downSql, m.DownSql = m.DownSql, downSql
	}
	if update {
		_, err = o.db.NewUpdate().Model(m).
			ModelTableExpr(o.tableName).
			Where("id = ?", m.Id).
			Exec(context.Background())
	} else {
		_, err = o.db.NewInsert().Model(m).
			ModelTableExpr(o.tableName).
			Exec(context.Background())
	}
	if !o.saveDownSql {
		downSql, m.DownSql = m.DownSql, downSql
	}
	return
}

func (o *Migrator) updateApplied(m *Migration) error {
	return o.writeApplied(m, true)
}

func (o *Migrator) markApplied(m *Migration) error {
	return o.writeApplied(m, false)
}

func (o *Migrator) markUnapplied(m *Migration) (err error) {
	if o.ro {
		return
	}
	_, err = o.db.NewDelete().
		Model(m).
		ModelTableExpr(o.tableName).
		Where("id = ?", m.Id).
		Exec(context.Background())
	return
}

func (o *Migrator) rollback(l Migrations) (err error) {
	slices.Reverse(l)
	for _, m := range l {
		err = m.Down()
		if err == nil {
			o.markUnapplied(m)
			m.SetUnapplied()
		} else {
			break
		}
	}
	return
}
