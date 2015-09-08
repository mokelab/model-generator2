package mysql

import (
	g "../"
	"../../model"
	"fmt"
	"io"
	"strings"
	"text/template"
)

const (
	OPTION_PRIMARY_KEYS = "primary_keys"
	TEMPLATE_DAO        = `package mysql

import (
	m "../"
	"database/sql"
)

type {{.TypeNameLocal}}DAO struct {
	connection *Connection
}

func New{{.TypeName}}DAO(connection *Connection) *{{.TypeNameLocal}}DAO {
	return &{{.TypeNameLocal}}DAO{
		connection: connection,
	}
}

func (d *{{.TypeNameLocal}}DAO) Create({{.CreateArgs}}) (m.{{.TypeName}}, error) {
	tr, err := d.connection.Begin()
	if err != nil {
		return m.{{.TypeName}}{}, err
	}
	defer tr.Rollback()
	st, err := tr.Prepare("INSERT INTO {{.TableName}}({{.Rows}},created_time,modified_time) VALUES({{.Placeholders}},unix_timestamp(now()),unix_timestamp(now()))")
	if err != nil {
		return m.{{.TypeName}}{}, err
	}
	defer st.Close()

	_, err = st.Exec({{.Args}})
	if err != nil {
		return m.{{.TypeName}}{}, err
	}

	tr.Commit()
	return m.{{.TypeName}}{
{{.Return}}	}, nil
}

func (d *{{.TypeNameLocal}}DAO) Get({{.PrimaryKeyArgsSignature}}) (m.{{.TypeName}}, error) {
	db := d.connection.Connect()

	st, err := db.Prepare("SELECT {{.Rows}} FROM {{.TableName}} WHERE {{.PrimaryKeyWhere}}")
	if err != nil {
		return m.{{.TypeName}}{}, err
	}
	defer st.Close()

	rows, err := st.Query({{.PrimaryKeyArgs}})
	if err != nil {
		return m.{{.TypeName}}{}, err
	}
	defer rows.Close()

	return d.scan(rows)
}

func (d *{{.TypeNameLocal}}DAO) Update({{.CreateArgs}}) (m.{{.TypeName}}, error) {
	db := d.connection.Connect()

	st, err := db.Prepare("UPDATE {{.TableName}} SET {{.UpdateSet}},modified_time=unix_timestamp(now()) WHERE {{.PrimaryKeyWhere}}")
	if err != nil {
		return m.{{.TypeName}}{}, err
	}
	defer st.Close()

	_, err = st.Exec({{.Args}}, {{.PrimaryKeyArgs}})
	if err != nil {
		return m.{{.TypeName}}{}, err
	}
	return m.{{.TypeName}}{
{{.Return}}	}, nil
}

func (d *{{.TypeNameLocal}}DAO) Delete({{.PrimaryKeyArgsSignature}}) error {
	db := d.connection.Connect()

	st, err := db.Prepare("DELETE FROM {{.TableName}} WHERE {{.PrimaryKeyWhere}}")
	if err != nil {
		return err
	}
	defer st.Close()

	_, err = st.Exec({{.PrimaryKeyArgs}})
	return err
}

func (d *{{.TypeNameLocal}}DAO) scan(rows *sql.Rows) (m.{{.TypeName}}, error) {
{{.ScanVars}}
	err := rows.Scan({{.ScanArgs}})
	if err != nil {
		return m.{{.TypeName}}{}, err
	}
	return m.{{.TypeName}}{
{{.Return}}	}, nil
}
`
)

type mysqlDAOGenerator struct {
}

func (g *mysqlDAOGenerator) Generate(table *model.Table, options g.Options, w io.Writer) {
	primaryKeyStr, ok := options[OPTION_PRIMARY_KEYS]
	if !ok {
		return
	}
	primaryKeys := findPrimaryKeys(table.Fields, primaryKeyStr)
	t, _ := template.New("dao").Parse(TEMPLATE_DAO)
	args := map[string]string{
		"TableName":               table.TableName(),
		"TypeNameLocal":           toLowerCamel(table.Name),
		"TypeName":                table.Name,
		"CreateArgs":              createArgsSignature(table.Fields),
		"Args":                    createExecArgs(table.Fields),
		"Rows":                    createTableRows(table.Fields),
		"UpdateSet":               createUpdateSet(table.Fields),
		"PrimaryKeyArgsSignature": createArgsSignature(primaryKeys),
		"PrimaryKeyArgs":          createExecArgs(primaryKeys),
		"PrimaryKeyWhere":         createUpdateSet(primaryKeys),
		"Placeholders":            createPlaceholders(table.Fields),
		"ScanVars":                createScanVars(table.Fields),
		"ScanArgs":                createScanArgs(table.Fields),
		"Return":                  createReturn(table.Fields),
	}
	t.Execute(w, args)
}

func findPrimaryKeys(fields []*model.Type, keyStr string) []*model.Type {
	out := make([]*model.Type, 0)
	keys := strings.SplitN(keyStr, ",", -1)
	for _, key := range keys {
		fmt.Println(key)
		for _, field := range fields {
			if field.Name == key {
				out = append(out, field)
				break
			}
		}
	}
	return out
}

func createArgsSignature(fields []*model.Type) string {
	return toList(fields, func(field *model.Type) string {
		return toLowerCamel(field.Name) + " " + field.Type
	})
}

func createExecArgs(fields []*model.Type) string {
	return toList(fields, func(field *model.Type) string {
		return toLowerCamel(field.Name)
	})
}

func createTableRows(fields []*model.Type) string {
	return toList(fields, func(field *model.Type) string {
		return field.SnakeName()
	})
}

func createUpdateSet(fields []*model.Type) string {
	return toList(fields, func(field *model.Type) string {
		return field.SnakeName() + "=?"
	})
}

func createPlaceholders(fields []*model.Type) string {
	return toList(fields, func(field *model.Type) string {
		return "?"
	})
}

func createScanVars(fields []*model.Type) string {
	out := ""
	for i, field := range fields {
		if i > 0 {
			out += "\n"
		}
		out += "\tvar " + toLowerCamel(field.Name) + " " + field.Type
	}
	return out
}

func createScanArgs(fields []*model.Type) string {
	return toList(fields, func(field *model.Type) string {
		return "&" + toLowerCamel(field.Name)
	})
}

func createReturn(fields []*model.Type) string {
	out := ""
	for _, field := range fields {
		out += "\t\t" + field.Name + " : " + toLowerCamel(field.Name) + ",\n"
	}
	return out
}

func toList(fields []*model.Type, f func(field *model.Type) string) string {
	out := ""
	for i, field := range fields {
		if i > 0 {
			out += ", "
		}
		out += f(field)
	}
	return out
}

func toLowerCamel(src string) string {
	return strings.ToLower(src[0:1]) + src[1:]
}
