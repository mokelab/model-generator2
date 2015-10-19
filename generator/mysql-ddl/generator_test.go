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
				Name:      "Name",
				Type:      "string",
				TableType: "varchar(32)",
			},
		},
	}
	options := map[string]string{
		OPTION_PRIMARY_KEYS: "Name",
	}
	g := &mysqlDDLGenerator{}
	w := &bytes.Buffer{}
	g.Generate(table, options, w)
	result := w.String()
	expected := `create table if not exists user(
  name varchar(32),
  created_time bigint,
  modified_time bigint,
  primary key(name)
) engine=InnoDB;
`
	if result != expected {
		t.Errorf("Unexpected result : %s", result)
	}
}

func Test_0001_multi_primary_key(t *testing.T) {
	table := &model.Table{
		Name: "User",
		Fields: []*model.Type{
			&model.Type{
				Name:      "Name",
				Type:      "string",
				TableType: "varchar(32)",
			},
			&model.Type{
				Name:      "UserId",
				Type:      "string",
				TableType: "varchar(32)",
			},
		},
	}
	options := map[string]string{
		OPTION_PRIMARY_KEYS: "Name,UserId",
	}
	g := &mysqlDDLGenerator{}
	w := &bytes.Buffer{}
	g.Generate(table, options, w)
	result := w.String()
	expected := `create table if not exists user(
  name varchar(32),
  user_id varchar(32),
  created_time bigint,
  modified_time bigint,
  primary key(name,user_id)
) engine=InnoDB;
`
	if result != expected {
		t.Errorf("Unexpected result : %s", result)
	}
}
