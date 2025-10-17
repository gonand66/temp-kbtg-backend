package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Create users table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		membership_id TEXT UNIQUE NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		membership_level TEXT DEFAULT 'Bronze',
		points INTEGER DEFAULT 0,
		joined_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = DB.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("✅ Database initialized successfully")

	// Insert sample data if table is empty
	insertSampleData()

	return nil
}

func insertSampleData() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)

	if count == 0 {
		insertQuery := `
		INSERT INTO users (membership_id, first_name, last_name, phone_number, email, membership_level, points, joined_date)
		VALUES 
			('LBK001234', 'สมชาย', 'ใจดี', '081-234-5678', 'somchai@example.com', 'Gold', 15420, '2023-06-15'),
			('LBK001235', 'สมหญิง', 'รักดี', '082-345-6789', 'somying@example.com', 'Silver', 8500, '2023-07-20'),
			('LBK001236', 'สมศักดิ์', 'มีสุข', '083-456-7890', 'somsak@example.com', 'Bronze', 2100, '2023-08-10');
		`

		if _, err := DB.Exec(insertQuery); err != nil {
			log.Printf("Failed to insert sample data: %v", err)
		} else {
			log.Println("✅ Sample data inserted")
		}
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
