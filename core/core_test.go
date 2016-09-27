package core

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {

	//	tx, _ := db.Begin()
	//	for i := 0; i < 393160+10000000; i++ {
	//		stmt.Exec(i, strconv.Itoa(i))
	//		if i%10000 == 0 {
	//			tx.Commit()
	//		}
	//	}
	//	tx.Commit()
}

func TestInsertData(t *testing.T) {
	Execution()
}
