package mysql

import (
	g "../"
	"../../model"
	"io"
	"strings"
	"text/template"
)

const (
	OPTION_PRIMARY_KEYS         = "primary_keys"
	OPTION_AUTO_GENERATE        = "auto_generate"
	OPTION_AUTO_GENERATE_LENGTH = "auto_generate_length"
	TEMPLATE_DAO                = `package mysql

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

	{{if .AutoGenerateKey}}
	{{.AutoGenerateVar}}, err := insertWithUUID({{.AutoGenerateKeyLength}}, func(id string) error {
		_, err = st.Exec({{.Args}})
		return err
	})
	{{else}}
	_, err = st.Exec({{.Args}})
	{{end}}
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
	if !rows.Next() {
		return m.{{.TypeName}}{}, err
	}

	return d.scan(rows)
}

func (d *{{.TypeNameLocal}}DAO) Update({{.UpdateArgs}}) (m.{{.TypeName}}, error) {
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

func NewGenerator() *mysqlDAOGenerator {
	return &mysqlDAOGenerator{}
}

func (o *mysqlDAOGenerator) Generate(table *model.Table, options g.Options, w io.Writer) {
	primaryKeyStr, ok := options[OPTION_PRIMARY_KEYS]
	if !ok {
		return
	}
	primaryKeys := findPrimaryKeys(table.Fields, primaryKeyStr)

	var fields []*model.Type
	autoGenerateKeyStr, _ := options[OPTION_AUTO_GENERATE]
	autoGenerateKeyLength, _ := options[OPTION_AUTO_GENERATE_LENGTH]
	if len(autoGenerateKeyStr) > 0 && len(autoGenerateKeyLength) > 0 {
		fields, _ = removeAutoGenerateKey(table.Fields, autoGenerateKeyStr)
	} else {
		autoGenerateKeyStr = ""
		fields = table.Fields
	}
	t, _ := template.New("dao").Parse(TEMPLATE_DAO)
	args := map[string]string{
		"TableName":               table.TableName(),
		"TypeNameLocal":           g.ToLowerCamel(table.Name),
		"TypeName":                table.Name,
		"AutoGenerateKey":         autoGenerateKeyStr,
		"AutoGenerateVar":         g.ToLowerCamel(autoGenerateKeyStr),
		"AutoGenerateKeyLength":   autoGenerateKeyLength,
		"CreateArgs":              g.CreateArgsSignature(fields),
		"UpdateArgs":              g.CreateArgsSignature(table.Fields),
		"Args":                    g.CreateExecArgs(table.Fields),
		"Rows":                    g.CreateTableRows(table.Fields),
		"UpdateSet":               g.CreateUpdateSet(table.Fields),
		"PrimaryKeyArgsSignature": g.CreateArgsSignature(primaryKeys),
		"PrimaryKeyArgs":          g.CreateExecArgs(primaryKeys),
		"PrimaryKeyWhere":         g.CreateUpdateSet(primaryKeys),
		"Placeholders":            g.CreatePlaceholders(table.Fields),
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
		for _, field := range fields {
			if field.Name == key {
				out = append(out, field)
				break
			}
		}
	}
	return out
}

func removeAutoGenerateKey(fields []*model.Type, key string) ([]*model.Type, *model.Type) {
	out := make([]*model.Type, 0)
	var autoGenerateKey *model.Type
	for _, field := range fields {
		if field.Name == key {
			autoGenerateKey = field
		} else {
			out = append(out, field)
		}
	}
	return out, autoGenerateKey
}

func createScanVars(fields []*model.Type) string {
	out := ""
	for i, field := range fields {
		if i > 0 {
			out += "\n"
		}
		out += "\tvar " + g.ToLowerCamel(field.Name) + " " + field.Type
	}
	return out
}

func createScanArgs(fields []*model.Type) string {
	return g.ToList(fields, func(field *model.Type) string {
		return "&" + g.ToLowerCamel(field.Name)
	})
}

func createReturn(fields []*model.Type) string {
	out := ""
	for _, field := range fields {
		out += "\t\t" + field.Name + " : " + g.ToLowerCamel(field.Name) + ",\n"
	}
	return out
}
