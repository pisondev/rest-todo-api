package app

import (
	"database/sql"
	"fmt"
	"os"
	"rest-todo-api/helper"
	"time"

	"github.com/joho/godotenv"
)

func NewDB() *sql.DB {
	err := godotenv.Load()
	helper.PanicIfError(err)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"),
	)

	db, err := sql.Open("mysql", dsn)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}
