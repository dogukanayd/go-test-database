package mysqlunit

import (
	"github.com/dogukanayd/go-test-database/databases/testmysql"
	"testing"
	"upper.io/db.v3"
)

type TestTable struct {
	Name string `db:"name"`
}

func TestNewUnit(t *testing.T) {
	var u Unit

	u.databaseName = testmysql.TestDatabaseName
	u.t = t
	connection, tearDown := u.Start()

	defer tearDown()

	t.Run("it should create user with name 'go awesome' and return", func(t *testing.T) {
		u := TestTable{
			Name: "go awesome",
		}

		connection.Collection("test_table").Insert(u)

		err := connection.Collection("test_table").Find(db.Cond{
			"name": "go awesome",
		}).One(&u)

		if err != nil {
			t.Fatalf("error in unit_test: %v", err)
		}

		if u.Name != "go awesome" {
			t.Error("can not find user name: go awesome")
		}
	})
}
