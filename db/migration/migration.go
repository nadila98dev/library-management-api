package migration

import (
	"fmt"
	"library-management-api/models"
	"log"
)


func Migration() {
	err := db.DB.AutoMigrate(
		&models.Users();
	)
	if err != nil {
		log.Fatal("Failed to migrate..")
	}
	fmt.Println("Migrated successfully")
}