package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	for {
		db, err := sql.Open("mysql", "root:@/test-db")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		fmt.Println("runnig")
		time.Sleep(3 * time.Second)
	}
}
