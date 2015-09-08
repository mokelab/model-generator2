package model

type Type struct {
	Name      string
	Type      string
	TableType string
}

func (t *Type) SnakeName() string {
	return toSnake(t.Name)
}
