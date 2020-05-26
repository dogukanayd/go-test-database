package mysql_unit

import (
	"log"
	"testing"
)

func TestNewUnitV2(t *testing.T) {
	db, def := NewUnitV2()

	// _, _ = connection.Exec("CREATE TABLE `test_table`(`id` int(11),`name` varchar(500)) ENGINE = InnoDB  DEFAULT CHARSET = utf8;")

	q := `CREATE TABLE test_table(id int(11),name varchar(500)) ENGINE = InnoDB  DEFAULT CHARSET = utf8;`
	_, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: connection should close

	defer def()

	t.Log("ping success")
}
