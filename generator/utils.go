package generator

import (
	"../model"
	"strings"
)

func RemoveAutoGenerateKey(fields []*model.Type, key string) ([]*model.Type, *model.Type) {
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

func CreateArgsSignature(fields []*model.Type) string {
	return ToList(fields, func(field *model.Type) string {
		return ToLowerCamel(field.Name) + " " + field.Type
	})
}

func CreateExecArgs(fields []*model.Type) string {
	return ToList(fields, func(field *model.Type) string {
		return ToLowerCamel(field.Name)
	})
}

func CreateTableRows(fields []*model.Type) string {
	return ToList(fields, func(field *model.Type) string {
		return field.SnakeName()
	})
}

func CreateUpdateSet(fields []*model.Type) string {
	return ToList(fields, func(field *model.Type) string {
		return field.SnakeName() + "=?"
	})
}

func CreatePlaceholders(fields []*model.Type) string {
	return ToList(fields, func(field *model.Type) string {
		return "?"
	})
}

func ToList(fields []*model.Type, f func(field *model.Type) string) string {
	out := ""
	for i, field := range fields {
		if i > 0 {
			out += ", "
		}
		out += f(field)
	}
	return out
}

func ToLowerCamel(src string) string {
	if len(src) == 0 {
		return ""
	}
	return strings.ToLower(src[0:1]) + src[1:]
}
