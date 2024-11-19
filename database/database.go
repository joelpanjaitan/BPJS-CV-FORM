package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB*sql.DB

func InitDatabase(){
	var err error
	DB, err = sql.Open("mysql","localhost:root@tcp(localhost:3306)/cv_app")
	if err!= nil{
		log.Fatal("Connection to database failed:", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal("Timeout in database", err)
	}

	log.Println("Database Connection successful")
}