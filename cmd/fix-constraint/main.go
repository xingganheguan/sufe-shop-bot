package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection parameters
	host := "192.168.110.68"
	port := "5432"
	dbname := "sufeshopbot"
	user := "sufeshopbot"
	password := "ooaK3Wj4XLW9x&bXfRRv##"
	sslmode := "disable"

	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get the underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	log.Println("Connected to database successfully")

	// Step 1: Check current indexes
	var indexes []struct {
		IndexName string `gorm:"column:indexname"`
		TableName string `gorm:"column:tablename"`
		IndexDef  string `gorm:"column:indexdef"`
	}

	err = db.Raw(`
		SELECT indexname, tablename, indexdef 
		FROM pg_indexes 
		WHERE tablename = 'message_templates' 
		AND schemaname = 'public'
	`).Scan(&indexes).Error

	if err != nil {
		log.Fatalf("Failed to query indexes: %v", err)
	}

	log.Println("\nCurrent indexes on message_templates table:")
	for _, idx := range indexes {
		log.Printf("- %s: %s", idx.IndexName, idx.IndexDef)
	}

	// Step 2: Drop the old unique constraint if it exists
	err = db.Exec("DROP INDEX IF EXISTS idx_message_templates_code").Error
	if err != nil {
		log.Printf("Warning: Failed to drop old index (may not exist): %v", err)
	} else {
		log.Println("\nDropped old index 'idx_message_templates_code' (if existed)")
	}

	// Step 3: Create new composite unique index
	// First, let's check if idx_code_lang already exists
	var existingIndex struct {
		Count int
	}
	err = db.Raw(`
		SELECT COUNT(*) as count
		FROM pg_indexes 
		WHERE tablename = 'message_templates' 
		AND indexname = 'idx_code_lang'
		AND schemaname = 'public'
	`).Scan(&existingIndex).Error

	if err != nil {
		log.Fatalf("Failed to check for existing index: %v", err)
	}

	if existingIndex.Count > 0 {
		log.Println("\nComposite index 'idx_code_lang' already exists!")
	} else {
		err = db.Exec("CREATE UNIQUE INDEX idx_code_lang ON message_templates (code, language)").Error
		if err != nil {
			log.Fatalf("Failed to create new composite index: %v", err)
		}
		log.Println("\nCreated new composite unique index 'idx_code_lang' on (code, language)")
	}

	// Step 4: Verify the changes
	err = db.Raw(`
		SELECT indexname, tablename, indexdef 
		FROM pg_indexes 
		WHERE tablename = 'message_templates' 
		AND schemaname = 'public'
	`).Scan(&indexes).Error

	if err != nil {
		log.Fatalf("Failed to query indexes after changes: %v", err)
	}

	log.Println("\nFinal indexes on message_templates table:")
	for _, idx := range indexes {
		log.Printf("- %s: %s", idx.IndexName, idx.IndexDef)
	}

	// Check for any duplicate (code, language) combinations
	var duplicates []struct {
		Code     string
		Language string
		Count    int
	}

	err = db.Raw(`
		SELECT code, language, COUNT(*) as count
		FROM message_templates
		GROUP BY code, language
		HAVING COUNT(*) > 1
	`).Scan(&duplicates).Error

	if err != nil {
		log.Printf("Warning: Failed to check for duplicates: %v", err)
	} else if len(duplicates) > 0 {
		log.Println("\nWARNING: Found duplicate (code, language) combinations:")
		for _, dup := range duplicates {
			log.Printf("- Code: %s, Language: %s, Count: %d", dup.Code, dup.Language, dup.Count)
		}
		log.Println("You may need to resolve these duplicates before the unique constraint can be properly enforced.")
	} else {
		log.Println("\nNo duplicate (code, language) combinations found. Constraint is properly enforced.")
	}

	log.Println("\nConstraint fix completed successfully!")
}