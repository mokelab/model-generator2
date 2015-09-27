package main

import (
	g "./generator"
	mysql_dao "./generator/mysql-dao"
	mysql_dao_test "./generator/mysql-dao-test"
	mysql_ddl "./generator/mysql-ddl"
	_ "./model"
	p "./parser"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

const (
	TYPE_MYSQL_DDL      = "mysql_ddl"
	TYPE_MYSQL_DAO      = "mysql_dao"
	TYPE_MYSQL_DAO_TEST = "mysql_dao_test"
)

func main() {
	var outType *string = flag.String("outType", TYPE_MYSQL_DDL, "output type")
	var srcFile *string = flag.String("src", "", "input file")
	var primaryKeys *string = flag.String("primaryKeys", "", "Primary Key Fields")
	flag.Parse()

	if len(*srcFile) == 0 {
		fmt.Errorf("src file must not be empty.\n")
		return
	}
	options := map[string]string{
		"primary_keys": *primaryKeys,
	}

	generator := createGenerator(*outType)
	if generator == nil {
		fmt.Errorf("Unknown type : %s", *outType)
		return
	}

	// parse
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *srcFile, nil, 0)
	if err != nil {
		fmt.Errorf("Failed to parse input : %s", err)
		return
	}
	tables, err := p.Parse(f)
	if tables == nil {
		return
	}
	generator.Generate(tables[0], options, os.Stdout)
}

func createGenerator(outType string) g.Generator {
	switch outType {
	case TYPE_MYSQL_DDL:
		return mysql_ddl.NewGenerator()
	case TYPE_MYSQL_DAO:
		return mysql_dao.NewGenerator()
	case TYPE_MYSQL_DAO_TEST:
		return mysql_dao_test.NewGenerator()
	default:
		return nil
	}
}
