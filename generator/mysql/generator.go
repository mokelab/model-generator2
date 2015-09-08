package mysql

import (
	g "../"
	"../../model"
	"fmt"
	"io"
)

const (
	OPTION_PRIMARY_KEYS = "primary_keys"
)

type mysqlDDLGenerator struct {
}

func (g *mysqlDDLGenerator) Generate(table *model.Table, options g.Options, w io.Writer) {
	out := "create table " + table.TableName() + "(\n"
	for i, field := range table.Fields {
		if i > 0 {
			out = out + ",\n"
		}
		out = out + fmt.Sprintf("  %s %s", field.SnakeName(), field.TableType)
	}
	// created and modified time
	out = out + ",\n"
	out = out + "  created_time bigint,\n"
	out = out + "  modified_time bigint"
	if primaryKeys, ok := options[OPTION_PRIMARY_KEYS]; ok {
		out = out + ",\n"
		out = out + fmt.Sprintf("  primary key(%s)", primaryKeys)
	}

	out = out + "\n)\n"
	fmt.Fprintf(w, out)
}
