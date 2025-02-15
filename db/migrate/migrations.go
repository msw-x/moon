package migrate

import (
	"io/fs"
	"sort"

	"github.com/uptrace/bun"
)

type Migrations []*Migration

func (o Migrations) Count() int {
	return len(o)
}

func (o Migrations) Applied() (applied Migrations) {
	for _, m := range o {
		if m.IsApplied() {
			applied = append(applied, m)
		}
	}
	applied.SortDesc()
	return
}

func (o Migrations) Unapplied() (unapplied Migrations) {
	for _, v := range o {
		if !v.IsApplied() {
			unapplied = append(unapplied, v)
		}
	}
	unapplied.SortAsc()
	return
}

func (o Migrations) Last() *Migration {
	return o[o.Count()-1]
}

func (o Migrations) LastGroupId() (lastGroupId int64) {
	for _, v := range o {
		if v.GroupId > lastGroupId {
			lastGroupId = v.GroupId
		}
	}
	return
}

func (o Migrations) LastGroup() (l Migrations) {
	groupId := o.LastGroupId()
	if groupId > 0 {
		for _, v := range o {
			if v.GroupId == groupId {
				l = append(l, v)
			}
		}
	}
	return
}

func (o Migrations) SetFuncs(db *bun.DB, splitter string) {
	for _, v := range o {
		v.SetFuncs(db, splitter)
	}
}

func (o Migrations) SortAsc() {
	sort.Slice(o, func(i, j int) bool {
		return o[i].Name < o[j].Name
	})
}

func (o Migrations) SortDesc() {
	sort.Slice(o, func(i, j int) bool {
		return o[i].Name > o[j].Name
	})
}

func (o *Migrations) Load(f fs.FS, final string) error {
	var skip bool
	return fs.WalkDir(f, ".", func(path string, e fs.DirEntry, err error) error {
		if skip {
			return nil
		}
		if err != nil {
			return err
		}
		if e.IsDir() {
			return nil
		}
		name, comment, isUp, isTx, ok := Name(path)
		if !ok {
			return nil
		}
		m := o.get(name)
		m.Comment = comment
		m.IsTx = isTx
		var sql string
		sql, err = ReadSql(f, path)
		if err != nil {
			return err
		}
		if isUp {
			m.UpSql = sql
			m.hasUpSql = true
		} else {
			m.DownSql = sql
			m.hasDownSql = true
		}
		skip = name == final
		return nil
	})
}

func (o Migrations) PreviewDown(specific string) (name string, down string, err error) {
	c := NewContext()
	for i := range o {
		name = o[i].Name
		down, err = o[i].PreviewDown(c)
		if err != nil {
			name = ""
			down = ""
			break
		}
		if name == specific {
			break
		}
	}
	return
}

func (o *Migrations) AutoGenerateDown() error {
	return o.AutoGenerateDownTrace(nil)
}

func (o *Migrations) AutoGenerateDownTrace(f func(*Migration, *Context)) (err error) {
	c := NewContext()
	for i := range *o {
		m := (*o)[i]
		err = m.AutoGenerateDown(c)
		if err != nil {
			break
		}
		if f != nil {
			f(m, c)
		}
	}
	return
}

func (o *Migrations) RepairDown(update func(*Migration) error) (names []string, err error) {
	c := NewContext()
	for i := range *o {
		m := (*o)[i]
		var ok bool
		ok, err = m.RepairDown(c)
		if ok {
			err = update(m)
			names = append(names, m.Name)
		}
		if err != nil {
			names = []string{}
			break
		}
	}
	return
}

func (o *Migrations) ViewSchema() (s string, err error) {
	c := NewContext()
	for i := range *o {
		m := (*o)[i]
		_, err = m.PreviewDown(c)
		if err != nil {
			break
		}
	}
	s = c.String()
	return
}

func (o *Migrations) get(name string) *Migration {
	for i := range *o {
		m := (*o)[i]
		if m.Name == name {
			return m
		}
	}
	*o = append(*o, &Migration{Name: name})
	return (*o)[len(*o)-1]
}
