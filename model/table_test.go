package model

import (
	"testing"
)

func TestTable0000_Name(t *testing.T) {
	table := Table{
		Name: "MokeName",
	}
	name := table.TableName()
	if name != "moke_name" {
		t.Errorf("TableName must be moke_name but %s", name)
		return
	}
}
