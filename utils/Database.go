package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "gorm.io/driver/mysql"
)

var Db *sql.DB

func init() {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	db := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		db,
	)

	if user == "" || password == "" || host == "" || port == "" || db == "" {
		panic("Check enviromental variables setup for database\n" + dsn)
	}

	// client, err := mongo.NewClient(options.Client().ApplyURI(dburl))

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	Db = database

}
