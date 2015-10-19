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
	expected := `package mysql

import (
	m "../"
	"database/sql"
	"testing"
)

func createUserDAO(db *sql.DB) *userDAO{
	return NewUserDAO(NewConnection(db))
}

func assertUser(t *testing.T, item *m.User, name string, age int) {
	caller := getCaller()
	if item.Name != name {
		t.Errorf("[%s] Name must be %s but %s", caller, name, item.Name)
	}
	if item.Age != age {
		t.Errorf("[%s] Age must be %d but %d", caller, age, item.Age)
	}

}

func hardDeleteUser(db *sql.DB, id string) {
	s, _ := db.Prepare("DELETE FROM user WHERE id=?")
	defer s.Close()
	s.Exec(id)
}

func TestUser_All(t *testing.T) {
	db, err := connect()
	if err != nil {
		t.Errorf("Failed to connect")
		return
	}
	defer db.Close()

	dao := createUserDAO(db)

	var name string = "NameVal"
	var age int = 5

	hardDeleteUser(db, id)

	item, err := dao.Create(name, age)
	if err != nil {
		t.Errorf("Failed to Create : %s", err)
		return
	}

	assertUser(t, &item, "NameVal", 5)

	item2, err := dao.Get(id)
	if err != nil {
		t.Errorf("Failed to Get : %s", err)
		return
	}
	assertUser(t, &item2, "NameVal", 5)

	name = "NewNameVal"
	age = 8

	item3, err := dao.Update(id, name, age)
	if err != nil {
		t.Errorf("Failed to Update : %s", err)
		return
	}
	assertUser(t, &item3, "NewNameVal", 8)

	item4, err := dao.Get(id)
	if err != nil {
		t.Errorf("Failed to Get : %s", err)
		return
	}
	assertUser(t, &item4, "NewNameVal", 8)

	err = dao.Delete(id)
	if err != nil {
		t.Errorf("Failed to Delete : %s", err)
		return
	}
}
`
	if result != expected {
		t.Errorf("Unexpected result : %s", result)
	}
}
