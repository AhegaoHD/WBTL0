package infrastructure

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	// Инициализируем соединение с БД
	connStr := "user=postgres password=qweQWE dbname=WB_Tech_level_0 sslmode=disable port=5555"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Проверяем соединение с БД
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
