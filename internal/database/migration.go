package database

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, dir string) error {
	dryRun := os.Getenv("DB_MIGRATE_DRY_RUN") == "true"

	// ensure schema_migrations
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`).Error; err != nil {
		return err
	}

	// advisory lock
	if err := db.Exec(`SELECT pg_advisory_lock(123456789)`).Error; err != nil {
		return err
	}
	defer db.Exec(`SELECT pg_advisory_unlock(123456789)`)

	var applied []string
	if err := db.Raw(`SELECT version FROM schema_migrations`).Scan(&applied).Error; err != nil {
		return err
	}

	appliedMap := map[string]bool{}
	for _, v := range applied {
		appliedMap[v] = true
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		version := strings.Split(filepath.Base(file), "_")[0]
		if appliedMap[version] {
			continue
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		log.Printf("[migrate] %s (dry-run=%v)", file, dryRun)

		if dryRun {
			continue
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(string(sqlBytes)).Error; err != nil {
				return err
			}
			return tx.Exec(
				`INSERT INTO schema_migrations(version) VALUES (?)`,
				version,
			).Error
		}); err != nil {
			return err
		}
	}

	return nil
}
