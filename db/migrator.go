package db

import (
	"io/fs"

	"github.com/msw-x/moon/db/migrate"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ulog"
)

type Migrator struct {
	log *ulog.Log
	ro  bool
	m   *migrate.Migrator
}

func NewMigrator(d *Db) *Migrator {
	o := new(Migrator)
	o.log = d.log.Branch("migrator")
	o.ro = d.ro
	o.m = migrate.NewMigrator(d.db).WithReadonly(o.ro)
	return o
}

func (o *Migrator) WithSaveDownSql(v bool) *Migrator {
	o.m.WithSaveDownSql(v)
	return o
}

func (o *Migrator) WithAutoDownSql(v bool) *Migrator {
	o.m.WithAutoDownSql(v)
	return o
}

func (o *Migrator) WithRollbackLost(v bool) *Migrator {
	o.m.WithRollbackLost(v)
	return o
}

func (o *Migrator) WithMarkAppliedOnSuccess(v bool) *Migrator {
	o.m.WithMarkAppliedOnSuccess(v)
	return o
}

type MigratorOptions struct {
	Lock         bool
	Rollback     bool
	PreviewDown  bool
	RepairDown   bool
	TraceSchema  bool
	PrettySchema bool
}

func (o *Migrator) Exec(fs fs.FS, opt MigratorOptions) (ok bool) {
	if opt.TraceSchema {
		o.m.WithAutoDownSqlTrace(func(m *migrate.Migration, c *migrate.Context) {
			var s string
			if opt.PrettySchema {
				s = c.Pretty()
			} else {
				s = c.String()
			}
			if s != "" {
				s = "\n" + s
			}
			o.log.Trace(m.String(), s)
		})
	}
	err := o.m.Init()
	if err == nil {
		if o.load(fs) {
			if o.m.MigrationsCount() == 0 {
				o.log.Info("no migrations")
				ok = true
			} else {
				o.log.Info(o.m.Status())
				if opt.Lock {
					o.log.Info("locked")
				} else {
					switch {
					case opt.Rollback:
						o.rollbackLast()
					case opt.PreviewDown:
						o.previewDown()
					case opt.RepairDown:
						o.repairDown()
					default:
						ok = o.migrate()
					}
				}
			}
		}
	} else {
		o.log.Error("init:", err)
	}
	return
}

func (o *Migrator) load(fs fs.FS) bool {
	err := o.m.Load(fs)
	if err == nil {
		o.log.Info("migrations loaded")
		return true
	} else {
		o.log.Error("load migrations:", err)
		return false
	}
}

func (o *Migrator) migrate() bool {
	if o.m.HasIncompleteMigrations() {
		if o.ro {
			return true
		}
		o.log.Info("migrate")
		err := o.m.Migrate()
		if err == nil {
			o.log.Info("migrate successfully")
			o.log.Info(o.m.Status())
		} else {
			o.log.Error("migrate:", err)
			return false
		}
	} else {
		o.log.Info("migrations are relevant")
	}
	return true
}

func (o *Migrator) rollbackLast() {
	if o.ro {
		return
	}
	if o.m.UnappliedMigrationsCount() == 0 {
		err := o.m.RollbackLast()
		if err == nil {
			o.log.Info("rollback last migrations successfully")
			o.log.Info(o.m.Status())
		} else {
			o.log.Error("rollback last migrations:", err)
		}
	} else {
		o.log.Info("last migrations already unapplied")
	}
}

func (o *Migrator) previewDown() {
	name, down, err := o.m.Migrations().PreviewDown()
	if err == nil {
		o.log.Infof("preview migration[%s] down:\n%s", name, down)
	} else {
		o.log.Error(err)
	}
}

func (o *Migrator) repairDown() {
	o.log.Info("repair down")
	l, err := o.m.RepairDown()
	if err == nil {
		s := ufmt.JoinSlice(l)
		if s == "" {
			s = "no"
		}
		o.log.Info("repaired down:", s)
	} else {
		o.log.Error(err)
	}
}
