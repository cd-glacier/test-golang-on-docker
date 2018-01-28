package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
}

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(db-server:3306)/test_db")
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	rows, err := db.Query("SELECT * FROM test_table")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
	}

	for {

	}

}
