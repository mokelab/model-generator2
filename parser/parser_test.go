package parser

import (
	"go/parser"
	"go/token"
	"testing"
)

func Test0000_simple(t *testing.T) {
	src := "package model\n" +
		"type Moke struct{\n" +
		"  Name string `varchar(32)`\n" +
		"  Age int `int`\n" +
		"}\n"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Errorf("Failed to parse input : %s", err)
		return
	}
	tables, err := Parse(f)
	if err != nil {
		t.Errorf("Failed to parse : %s", err)
		return
	}
	if tables == nil {
		t.Errorf("return must not be nil")
		return
	}
	if len(tables) != 1 {
		t.Errorf("tables must be 1 but %d", len(tables))
		return
	}
	table := tables[0]
	if table.Name != "Moke" {
		t.Errorf("Name must be Moke but %s", table.Name)
		return
	}
	if len(table.Fields) != 2 {
		t.Errorf("Fields must be 2 but %d", len(table.Fields))
		return
	}
	field := table.Fields[0]
	if field == nil {
		t.Errorf("Field[0] must not be nil")
		return
	}
	if field.Name != "Name" {
		t.Errorf("Field[0].Name must be Name but %s", field.Name)
		return
	}
	if field.Type != "string" {
		t.Errorf("Field[0].Type must be string but %s", field.Type)
		return
	}
	if field.TableType != "varchar(32)" {
		t.Errorf("Field[0].TableType must be varchar(32) but %s", field.TableType)
		return
	}
}
