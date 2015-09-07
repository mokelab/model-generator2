package model

type Table struct {
	Name   string
	Fields []*Type
}

func (t Table) TableName() string {
	return toSnake(t.Name)
}
