package generator

import (
	"../model"
	"io"
)

type Generator interface {
	Generate(table *model.Table, options Options, w io.Writer)
}

type Options map[string]string
