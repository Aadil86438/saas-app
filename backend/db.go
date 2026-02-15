package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var lDB *sql.DB

func InitDB() {
	log.Println("InitDB(+)")
	
	lDBPassword := os.Getenv("DB_PASSWORD")
	if lDBPassword == "" {
		lDBPassword = "2466"
	}
	
	lConnectionString := os.Getenv("DATABASE_URL")

if lConnectionString == "" {
    lConnectionString = "host=localhost port=5432 user=postgres password=" + lDBPassword + " dbname=todos_db sslmode=disable"
}

	
	lDBInstance, lErr := sql.Open("postgres", lConnectionString)
	if lErr != nil {
		log.Fatal("Failed to connect to database:", lErr)
	}
	
	lErr = lDBInstance.Ping()
	if lErr != nil {
		log.Fatal("Failed to ping database:", lErr)
	}
	
	lDB = lDBInstance
	
	lErr = CreateTables()
	if lErr != nil {
		log.Fatal("Failed to create tables:", lErr)
	}
	
	log.Println("InitDB(-)")
}

func CreateTables() error {
	log.Println("CreateTables(+)")
	
	lUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	
	lTodosTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(200) NOT NULL,
		content TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	
	lSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL
	);`
	
	_, lErr := lDB.Exec(lUsersTable)
	if lErr != nil {
		log.Println("CreateTables(-) error:", lErr)
		return lErr
	}
	
	_, lErr = lDB.Exec(lTodosTable)
	if lErr != nil {
		log.Println("CreateTables(-) error:", lErr)
		return lErr
	}
	
	_, lErr = lDB.Exec(lSessionsTable)
	if lErr != nil {
		log.Println("CreateTables(-) error:", lErr)
		return lErr
	}
	
	log.Println("CreateTables(-)")
	return nil
}

func GetDB() *sql.DB {
	return lDB
}

