package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
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

	user := User{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": user.Name,
		})
	})
	r.Run()
}
