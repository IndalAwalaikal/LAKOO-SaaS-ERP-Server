package database

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"lakoo/backend/pkg/config"
)

func NewMySQLConnection(cfg *config.Config) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&multiStatements=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var db *sqlx.DB
	var err error

	// Retry logic for DB connection
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect to MySQL (attempt %d/10): %v. Retrying in 5s...", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to MySQL after 10 attempts: %v", err)
	}

	log.Println("Connected to MySQL successfully")
	return db
}
