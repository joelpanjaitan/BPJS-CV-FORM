package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB*sql.DB

func InitDatabase(){
	var err error
	DB, err = sql.Open("mysql","root:password@tcp(localhost:3307)/api_db")
	if err!= nil{
		log.Fatal("Connection to database failed:", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal("Timeout in database", err)
	}

	log.Println("Database is connected successfully")
}