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
	TEMPLATE_DAO                = `package mock

import (
	m "../"
)

type {{.TypeName}}DAO struct {
	getResult    m.{{.TypeName}}
	createResult m.{{.TypeName}}
	updateResult m.{{.TypeName}}

	getError    error
	createError error
	updateError error
	deleteError error
}

func (d *{{.TypeNameLocal}}DAO) Get({{.PrimaryKeyArgsSignature}}) (m.{{.TypeName}}, error) {
	return d.getResult, d.getError
}

func (d *{{.TypeNameLocal}}DAO) Create({{.CreateArgs}}) (m.{{.TypeName}}, error) {
	return d.createResult, d.createError
}

func (d *{{.TypeNameLocal}}DAO) Update({{.UpdateArgs}}) (m.{{.TypeName}}, error) {
	return d.updateResult, d.updateError
}

func (d *{{.TypeNameLocal}}DAO) Delete({{.PrimaryKeyArgsSignature}}) error {
	return d.deleteError
}
`
)

type mysqlMockDAOGenerator struct {
}

func NewGenerator() *mysqlMockDAOGenerator {
	return &mysqlMockDAOGenerator{}
}

func (o *mysqlMockDAOGenerator) Generate(table *model.Table, options g.Options, w io.Writer) {
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
		"UpdateArgs":              g.CreateArgsSignature(table.Fields),
		"Args":                    g.CreateExecArgs(fields),
		"ArgsUpdate":              g.CreateExecArgs(table.Fields),
		"Rows":                    g.CreateTableRows(table.Fields),
		"UpdateSet":               g.CreateUpdateSet(table.Fields),
		"PrimaryKeyArgsSignature": g.CreateArgsSignature(primaryKeys),
		"PrimaryKeyArgs":          g.CreateExecArgs(primaryKeys),
		"PrimaryKeyWhere":         g.CreateUpdateSet(primaryKeys),
		"Placeholders":            g.CreatePlaceholders(table.Fields),
		"AssertStatements":        createAssertStatements(fields),
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
	case "int64":
		return strconv.Itoa(len(prefix) + 50)
	case "string":
		return "\"" + prefix + field.Name + "Val\""
	default:
		return "\"\""
	}
}
