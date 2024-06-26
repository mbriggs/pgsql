// Package pgsql helps build SQL queries.
package pgsql

import (
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

type SQLWriter interface {
	WriteSQL(*strings.Builder, *Args)
}

type rawSQL string

func (r rawSQL) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString(string(r))
}

type Ident pgx.Identifier

func (i Ident) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString(pgx.Identifier(i).Sanitize())
}

type Param struct {
	Value interface{}
}

func (p Param) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString(args.Use(p.Value).String())
}

func Build(ab SQLWriter) (string, []interface{}) {
	sb := &strings.Builder{}
	args := &Args{}

	ab.WriteSQL(sb, args)

	return sb.String(), args.Values()
}

type binaryExpr struct {
	left  SQLWriter
	op    string
	right SQLWriter
}

func (be *binaryExpr) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteByte('(')
	be.left.WriteSQL(sb, args)
	sb.WriteByte(' ')
	sb.WriteString(be.op)
	sb.WriteByte(' ')
	be.right.WriteSQL(sb, args)
	sb.WriteByte(')')
}

type RowMap map[string]interface{}

func (rm RowMap) InsertData() ([]string, *ValuesStatement) {
	keys := rm.sortedKeys()

	values := make([]interface{}, len(keys))
	for i, k := range keys {
		values[i] = rm[k]
	}

	vs := Values()
	vs.Row(values...)

	return keys, vs
}

func (rm RowMap) UpdateData() []*Assignment {
	keys := rm.sortedKeys()

	assignments := make([]*Assignment, len(keys))
	for i, k := range keys {
		assignments[i] = &Assignment{Left: rawSQL(k), Right: &Param{Value: rm[k]}}
	}

	return assignments
}

func (rm RowMap) sortedKeys() []string {
	keys := make([]string, 0, len(rm))
	for k := range rm {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

type FormatString struct {
	s    string
	args []interface{}
}

func (fs *FormatString) WriteSQL(sb *strings.Builder, args *Args) {
	sb.WriteString(args.Format(fs.s, fs.args...))
}

type whereList []SQLWriter

func (wl whereList) WriteSQL(sb *strings.Builder, args *Args) {
	if len(wl) == 0 {
		return
	}

	sb.WriteString(" where ")

	for i, expr := range wl {
		if i > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteByte('(')
		expr.WriteSQL(sb, args)
		sb.WriteByte(')')
	}
}

type returningList []SQLWriter

func (rl returningList) WriteSQL(sb *strings.Builder, args *Args) {
	if len(rl) == 0 {
		return
	}

	sb.WriteString(" returning ")

	for i, expr := range rl {
		if i > 0 {
			sb.WriteString(", ")
		}
		expr.WriteSQL(sb, args)
	}
}
