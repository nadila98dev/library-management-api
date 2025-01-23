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
		fmt.Fatal("Failed to migrate..")
	}
	fmt.Println("Migrated successfully")
}