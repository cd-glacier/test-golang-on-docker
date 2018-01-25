package main

import (
	"database/sql"
	"fmt"
	"g-hyoga/kyuko/go/model"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:password@/kyuko")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("SELECT * FROM canceled_class")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}

	canceledclasses := []model.CanceledClass{}
	canceledclasses, err = model.ScanCanceledClass(rows)
	if err != nil {
		fmt.Println(err.Error)
	}

	fmt.Println(canceledclasses)
}
