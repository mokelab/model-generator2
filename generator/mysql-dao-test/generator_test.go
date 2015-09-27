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
	g := &mysqlDAOTestGenerator{}
	w := &bytes.Buffer{}
	g.Generate(table, options, w)
	result := w.String()
	expected := `create table user(
  name varchar(32),
  created_time bigint,
  modified_time bigint,
  primary key(name)
)
`
	if result != expected {
		t.Errorf("Unexpected result : %s", result)
	}
}
