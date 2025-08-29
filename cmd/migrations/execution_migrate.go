package migrations

import (
	"fmt"
	"log"

	"github.com/universal-go-service/boilerplate/internal/domain/entities"
	"github.com/universal-go-service/boilerplate/pkg/providers/database"
)

func ExecuteMigration(db database.DatabaseProvider) {
	err := db.GetDB().AutoMigrate(&entities.Item{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Migration executed successfully")
}
