package mysql

import (
	"../../model"
	"bytes"
	"testing"
)

func Test_0000_simple(t *testing.T) {
	table := &model.Table{
		Name: "User",
		Fields: []*model.Type{
			&model.Type{
				Name:      "Id",
				Type:      "string",
				TableType: "varchar(48)",
			},
			&model.Type{
				Name:      "Name",
				Type:      "string",
				TableType: "varchar(32)",
			},
			&model.Type{
				Name:      "Age",
				Type:      "int",
				TableType: "int",
			},
		},
	}
	options := map[string]string{
		OPTION_PRIMARY_KEYS:         "Id",
		OPTION_AUTO_GENERATE:        "Id",
		OPTION_AUTO_GENERATE_LENGTH: "48",
	}
	g := &mysqlMockDAOGenerator{}
	w := &bytes.Buffer{}
	g.Generate(table, options, w)
	result := w.String()
	expected := `package mock

import (
	m "../"
)

type UserDAO struct {
	getResult    m.User
	createResult m.User
	updateResult m.User

	getError    error
	createError error
	updateError error
	deleteError error
}

func (d *userDAO) Get(id string) (m.User, error) {
	return d.getResult, d.getError
}

func (d *userDAO) Create(name string, age int) (m.User, error) {
	return d.createResult, d.createError
}

func (d *userDAO) Update(id string, name string, age int) (m.User, error) {
	return d.updateResult, d.updateError
}

func (d *userDAO) Delete(id string) error {
	return d.deleteError
}
`
	if result != expected {
		t.Errorf("Unexpected result : %s", result)
	}
}
