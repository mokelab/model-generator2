package mysql

import (
	g "../"
	"../../model"
	"fmt"
	"io"
	"strings"
)

const (
	OPTION_PRIMARY_KEYS = "primary_keys"
)

type mysqlDDLGenerator struct {
}

func NewGenerator() *mysqlDDLGenerator {
	return &mysqlDDLGenerator{}
}

func (g *mysqlDDLGenerator) Generate(table *model.Table, options g.Options, w io.Writer) {
	out := "create table if not exists " + table.TableName() + "(\n"
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
		out = out + fmt.Sprintf("  primary key(%s)", g.toSnakeNames(primaryKeys))
	}

	out = out + "\n) engine=InnoDB;\n"
	fmt.Fprintf(w, out)
}

func (g *mysqlDDLGenerator) toSnakeNames(src string) string {
	list := strings.SplitN(src, ",", -1)
	out := ""
	for i, key := range list {
		if i > 0 {
			out = out + ","
		}
		field := model.Type{Name: key}
		out = out + field.SnakeName()
	}
	return out
}
