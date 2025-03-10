package migrate

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/msw-x/moon/tabtable"
	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ustring"
)

type Schema struct {
	Tables []Table
}

func (o *Schema) String() string {
	return ufmt.JoinSliceFuncWith("\n", o.Tables, func(v Table) string {
		return v.String()
	})
}

func (o *Schema) Pretty() string {
	return ufmt.JoinSliceFuncWith("\n", o.Tables, func(v Table) string {
		return v.Pretty()
	})
}

func (o *Schema) AddTable(name string) error {
	err := ValidateName(name)
	if err != nil {
		return fmt.Errorf("table '%s' %v", name, err)
	}
	t := o.table(name)
	if t != nil {
		return fmt.Errorf("table[%s] already exists", name)
	}
	o.Tables = append(o.Tables, Table{
		Name: name,
	})
	return nil
}

func (o *Schema) DropTable(name string) error {
	return uerr.Unimplemented()
}

func (o *Schema) AddColumn(tableName, columnName, columnType, columnConstraints string) error {
	err := ValidateName(columnName)
	if err != nil {
		return fmt.Errorf("table '%s' column '%s' %v", tableName, columnName, err)
	}
	t, err := o.Table(tableName)
	if err != nil {
		return err
	}
	c := t.column(columnName)
	if c != nil {
		return fmt.Errorf("column[%s] already exists", columnName)
	}
	t.Columns = append(t.Columns, Column{
		Name:        columnName,
		Type:        columnType,
		Constraints: columnConstraints,
	})
	return nil
}

func (o *Schema) DromColumn(tableName, columnName string) error {
	return uerr.Unimplemented()
}

func (o *Schema) AlterColumnType(tableName, columnName, columnType string) error {
	c, err := o.Column(tableName, columnName)
	if err == nil {
		c.Type = columnType
	}
	return err
}

func (o *Schema) RemoveColumnKeyConstraint(tableName, columnName string) (constraint string, err error) {
	var c *Column
	c, err = o.Column(tableName, columnName)
	if err == nil {
		l := strings.Fields(c.Constraints)
		r := slices.DeleteFunc(l, func(s string) bool {
			if slices.Contains([]string{"UNIQUE"}, strings.ToUpper(s)) {
				constraint = s
				return true
			}
			return false
		})
		if constraint == "" {
			err = fmt.Errorf("table[%s] column[%s] key constraint not found", tableName, columnName)
		} else {
			c.Constraints = strings.Join(r, " ")
		}
	}
	return
}

func (o *Schema) Table(name string) (r *Table, err error) {
	r = o.table(name)
	if r == nil {
		err = fmt.Errorf("table[%s] not found", name)
	}
	return
}

func (o *Schema) Column(tableName, columnName string) (r *Column, err error) {
	var t *Table
	t, err = o.Table(tableName)
	if err == nil {
		var c *Column
		c, err = t.Column(columnName)
		if err == nil {
			r = c
		}
	}
	return
}

func (o *Schema) ColumnType(tableName, columnName string) (columnType string, columnConstraints string, err error) {
	var c *Column
	c, err = o.Column(tableName, columnName)
	if err == nil {
		columnType = c.Type
		columnConstraints = c.Constraints
	}
	return
}

func (o *Schema) table(name string) (r *Table) {
	for i, v := range o.Tables {
		if v.Name == name {
			r = &o.Tables[i]
			break
		}
	}
	return
}

type Table struct {
	Name    string
	Columns []Column
}

const columnTab = "  "

func (o *Table) String() string {
	return o.Name + "\n" + ufmt.JoinSliceFuncWith("\n", o.Columns, func(v Column) string {
		return columnTab + v.String()
	})
}

func (o *Table) Pretty() string {
	t := tabtable.New()
	for _, v := range o.Columns {
		t.Write(columnTab+v.Name, v.Type, v.Constraints)
	}
	return o.Name + "\n" + t.String()
}

func (o *Table) Column(name string) (r *Column, err error) {
	r = o.column(name)
	if r == nil {
		err = fmt.Errorf("column[%s] not found", name)
	}
	return
}

func (o *Table) column(name string) (r *Column) {
	for i, v := range o.Columns {
		if v.Name == name {
			r = &o.Columns[i]
			break
		}
	}
	return
}

type Column struct {
	Name        string
	Type        string
	NotNull     bool
	Default     string
	Constraints string
}

func (o *Column) String() string {
	return ufmt.NotableJoin(o.Name, o.Type, o.Constraints)
}

func ValidateName(s string) error {
	if s == "" {
		return errors.New("is empty")
	}
	if !ustring.IsLower(s) {
		return errors.New("contains capital letters")
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '.' {
			return errors.New("contains invalid symbols")
		}
	}
	return nil
}
