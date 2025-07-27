// Package database provides database connection and initialization utilities.
package database

import (
	"log"
	"money-go/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
)

type Database struct {
	db *sqlx.DB
}

func NewDB() *Database {
	dbPath := "./db.sqlite"
	log.Printf("Attempting to connect to database at: %s", dbPath)
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection created successfully")
	return &Database{
		db: db,
	}
}

func (d *Database) CreateTables() error {
	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        balance TEXT NOT NULL DEFAULT '0',
        created_at TEXT NOT NULL
    );`

	if _, err := d.db.Exec(userTable); err != nil {
		return err
	}

	transactionTable := `
    CREATE TABLE IF NOT EXISTS transactions (
        id TEXT PRIMARY KEY,
        from_user_id TEXT NOT NULL,
        to_user_id TEXT NOT NULL,
        amount TEXT NOT NULL,
        currency TEXT NOT NULL DEFAULT 'USD',
        status TEXT NOT NULL DEFAULT 'pending',
        created_at TEXT NOT NULL,
        FOREIGN KEY (from_user_id) REFERENCES users(id),
        FOREIGN KEY (to_user_id) REFERENCES users(id)
    );`

	if _, err := d.db.Exec(transactionTable); err != nil {
		return err
	}

	return nil
}

func (d *Database) SeedUsers() error {
	seedUsers := []models.User{
		{
			ID:        uuid.New().String(),
			Name:      "Alice Johnson",
			Email:     "alice@example.com",
			Balance:   decimal.NewFromFloat(1000.00),
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Bob Smith",
			Email:     "bob@example.com",
			Balance:   decimal.NewFromFloat(750.50),
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Charlie Brown",
			Email:     "charlie@example.com",
			Balance:   decimal.NewFromFloat(2500.75),
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Diana Prince",
			Email:     "diana@example.com",
			Balance:   decimal.NewFromFloat(500.00),
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}

	insertQuery := `
        INSERT OR REPLACE INTO users (id, name, email, balance, created_at)
        VALUES (?, ?, ?, ?, ?)`

	for _, user := range seedUsers {
		_, err := d.db.Exec(insertQuery,
			user.ID,
			user.Name,
			user.Email,
			user.Balance.String(),
			user.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	log.Printf("Seeded %d users successfully", len(seedUsers))
	return nil
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}

func (d *Database) Close() error {
	return d.db.Close()
}

func NewConnection() *sqlx.DB {
	log.Printf("Tyring to set up DB")
	database := NewDB()

	if err := database.CreateTables(); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	if err := database.SeedUsers(); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	if err := database.db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")
	return database.GetDB()
}
