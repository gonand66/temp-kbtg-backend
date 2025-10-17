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
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Create transfers table
	createTransfersTable := `
	CREATE TABLE IF NOT EXISTS transfers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_user_id INTEGER NOT NULL,
		to_user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK (amount > 0),
		status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
		note TEXT,
		idempotency_key TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		completed_at TEXT,
		fail_reason TEXT,
		FOREIGN KEY (from_user_id) REFERENCES users(id),
		FOREIGN KEY (to_user_id) REFERENCES users(id)
	);`

	if _, err = DB.Exec(createTransfersTable); err != nil {
		return fmt.Errorf("failed to create transfers table: %v", err)
	}

	// Create indexes for transfers
	transferIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_transfers_from ON transfers(from_user_id);",
		"CREATE INDEX IF NOT EXISTS idx_transfers_to ON transfers(to_user_id);",
		"CREATE INDEX IF NOT EXISTS idx_transfers_created ON transfers(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_transfers_idem_key ON transfers(idempotency_key);",
	}

	for _, indexSQL := range transferIndexes {
		if _, err = DB.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create transfer index: %v", err)
		}
	}

	// Create point_ledger table
	createLedgerTable := `
	CREATE TABLE IF NOT EXISTS point_ledger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		change INTEGER NOT NULL,
		balance_after INTEGER NOT NULL,
		event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
		transfer_id INTEGER,
		reference TEXT,
		metadata TEXT,
		created_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (transfer_id) REFERENCES transfers(id)
	);`

	if _, err = DB.Exec(createLedgerTable); err != nil {
		return fmt.Errorf("failed to create point_ledger table: %v", err)
	}

	// Create indexes for point_ledger
	ledgerIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_ledger_user ON point_ledger(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_ledger_transfer ON point_ledger(transfer_id);",
		"CREATE INDEX IF NOT EXISTS idx_ledger_created ON point_ledger(created_at);",
	}

	for _, indexSQL := range ledgerIndexes {
		if _, err = DB.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create ledger index: %v", err)
		}
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
