package database

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) {
	log.Println("Running database migrations...")

	// Folder tempat file migration berada
	migrationDir := "./migrations"
	
	// Jika dijalankan dari Docker, folder mungkin berbeda, kita cek
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		migrationDir = "server/api/migrations" // fallback untuk local dev
	}

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Printf("Warning: Could not read migration directory: %v", err)
		return
	}

	// Ambil semua file .up.sql dan urutkan
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		log.Printf("Applying migration: %s", fileName)
		
		filePath := filepath.Join(migrationDir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading migration file %s: %v", fileName, err)
			continue
		}

		// Jalankan SQL
		_, err = db.Exec(string(content))
		if err != nil {
			// Jika error karena tabel sudah ada, kita abaikan saja
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Migration %s skipped (tables already exist)", fileName)
				continue
			}
			log.Printf("Error executing migration %s: %v", fileName, err)
		} else {
			log.Printf("Migration %s applied successfully", fileName)
		}
	}

	log.Println("Database migrations completed.")
}
