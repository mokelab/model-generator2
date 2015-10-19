package parser

import (
	"../model"
	"errors"
	"go/ast"
	_ "go/parser"
	"go/token"
)

func Parse(f *ast.File) ([]*model.Table, error) {
	tables := make([]*model.Table, 0)
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		table, err := parseSpecs(genDecl.Specs)
		if err != nil {
			continue
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func parseSpecs(specs []ast.Spec) (*model.Table, error) {
	for _, spec := range specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		fields, err := parseFields(structType.Fields)
		if err != nil {
			return nil, err
		}
		return &model.Table{
			Name:   typeSpec.Name.Name,
			Fields: fields,
		}, nil
	}
	return nil, errors.New("Failed to find struct")
}

func parseFields(fields *ast.FieldList) ([]*model.Type, error) {
	list := make([]*model.Type, 0)
	for _, field := range fields.List {
		if field.Tag == nil {
			return nil, errors.New("Tag not found")
		}
		tag := newTag(field.Tag.Value)
		tagValue := tag.get("ct")
		identType, ok := field.Type.(*ast.Ident)
		if !ok {
			return nil, errors.New("Type is not *ast.Ident")
		}
		list = append(list, &model.Type{
			Name:      field.Names[0].Name,
			Type:      identType.Name,
			TableType: tagValue,
		})
	}
	return list, nil
}
