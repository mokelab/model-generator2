package mysql

import (
	g "../"
	"../../model"
	"io"
	"strconv"
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
	"testing"
)

func create{{.TypeName}}DAO(db *sql.DB) *{{.TypeNameLocal}}DAO{
	return New{{.TypeName}}DAO(NewConnection(db))
}

func assert{{.TypeName}}(t *testing.T, item *m.{{.TypeName}}, {{.CreateArgs}}) {
	caller := getCaller()
{{.AssertStatements}}
}

func hardDelete{{.TypeName}}(db *sql.DB, {{.PrimaryKeyArgsSignature}}) {
	s, _ := db.Prepare("DELETE FROM {{.TableName}} WHERE {{.PrimaryKeyWhere}}")
	defer s.Close()
	s.Exec({{.PrimaryKeyArgs}})
}

func Test{{.TypeName}}_All(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Errorf("Failed to connect")
		return
	}
	defer db.Close()

	dao := create{{.TypeName}}DAO(db)

{{.TestVars}}
	hardDelete{{.TypeName}}(db, {{.PrimaryKeyArgs}})

	item, err := dao.Create({{.Args}})
	if err != nil {
		t.Errorf("Failed to Create : %s", err)
		return
	}

	assert{{.TypeName}}(t, &item, {{.AssertVars}})

	item2, err := dao.Get({{.PrimaryKeyArgs}})
	if err != nil {
		t.Errorf("Failed to Get : %s", err)
		return
	}
	assert{{.TypeName}}(t, &item2, {{.AssertVars}})

{{.TestVars2}}
	item3, err := dao.Update({{.ArgsUpdate}})
	if err != nil {
		t.Errorf("Failed to Update : %s", err)
		return
	}
	assert{{.TypeName}}(t, &item3, {{.AssertVars2}})

	item4, err := dao.Get({{.PrimaryKeyArgs}})
	if err != nil {
		t.Errorf("Failed to Get : %s", err)
		return
	}
	assert{{.TypeName}}(t, &item4, {{.AssertVars2}})

	err = dao.Delete({{.PrimaryKeyArgs}})
	if err != nil {
		t.Errorf("Failed to Delete : %s", err)
		return
	}
}
`
)

type mysqlDAOTestGenerator struct {
}

func NewGenerator() *mysqlDAOTestGenerator {
	return &mysqlDAOTestGenerator{}
}

func (o *mysqlDAOTestGenerator) Generate(table *model.Table, options g.Options, w io.Writer) {
	primaryKeyStr, ok := options[OPTION_PRIMARY_KEYS]
	if !ok {
		return
	}
	primaryKeys := findPrimaryKeys(table.Fields, primaryKeyStr)

	var fields []*model.Type
	autoGenerateKeyStr, _ := options[OPTION_AUTO_GENERATE]
	autoGenerateKeyLength, _ := options[OPTION_AUTO_GENERATE_LENGTH]
	if len(autoGenerateKeyStr) > 0 && len(autoGenerateKeyLength) > 0 {
		fields, _ = g.RemoveAutoGenerateKey(table.Fields, autoGenerateKeyStr)
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
		"Args":                    g.CreateExecArgs(fields),
		"ArgsUpdate":              g.CreateExecArgs(table.Fields),
		"Rows":                    g.CreateTableRows(table.Fields),
		"UpdateSet":               g.CreateUpdateSet(table.Fields),
		"PrimaryKeyArgsSignature": g.CreateArgsSignature(primaryKeys),
		"PrimaryKeyArgs":          g.CreateExecArgs(primaryKeys),
		"PrimaryKeyWhere":         g.CreateUpdateSet(primaryKeys),
		"Placeholders":            g.CreatePlaceholders(table.Fields),
		"AssertStatements":        createAssertStatements(fields),
		"TestVars":                createTestVars(fields),
		"TestVars2":               createTestVars2(fields),
		"AssertVars":              createAssertVars(fields),
		"AssertVars2":             createAssertVars2(fields),
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

func createReturn(fields []*model.Type) string {
	out := ""
	for _, field := range fields {
		out += "\t\t" + field.Name + " : " + g.ToLowerCamel(field.Name) + ",\n"
	}
	return out
}

func createAssertStatements(fields []*model.Type) string {
	out := ""
	for _, field := range fields {
		varName := g.ToLowerCamel(field.Name)
		var placeholder string
		if field.Type == "int" {
			placeholder = "%d"
		} else {
			placeholder = "%s"
		}
		out += "\tif item." + field.Name + " != " + varName + " {\n" +
			"\t\tt.Errorf(\"[%s] " + field.Name + " must be " + placeholder + " but " + placeholder + "\", caller, " + varName + ", item." + field.Name + ")\n" +
			"\t}\n"
	}
	return out
}

func createTestVars(fields []*model.Type) string {
	out := ""
	for _, field := range fields {
		varName := g.ToLowerCamel(field.Name)
		out += "\tvar " + varName + " " + field.Type + " = " + createValue(field, "") + "\n"
	}
	return out
}

func createTestVars2(fields []*model.Type) string {
	out := ""
	for _, field := range fields {
		varName := g.ToLowerCamel(field.Name)
		out += "\t" + varName + " = " + createValue(field, "New") + "\n"
	}
	return out
}

func createAssertVars(fields []*model.Type) string {
	return g.ToList(fields, func(field *model.Type) string {
		return createValue(field, "")
	})
}

func createAssertVars2(fields []*model.Type) string {
	return g.ToList(fields, func(field *model.Type) string {
		return createValue(field, "New")
	})
}

func createValue(field *model.Type, prefix string) string {
	switch field.Type {
	case "int":
		return strconv.Itoa(len(prefix) + 5)
	case "string":
		return "\"" + prefix + field.Name + "Val\""
	default:
		return "\"\""
	}
}
