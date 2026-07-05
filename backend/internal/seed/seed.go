package seed

import (
	"log"

	"gorm.io/gorm"
)

// Run seeds all initial data if not already present
func Run(db *gorm.DB) error {
	log.Println("Seeding database...")

	if err := seedCareerGoals(db); err != nil {
		return err
	}
	if err := seedCompanies(db); err != nil {
		return err
	}
	if err := seedKnowledgeGraph(db); err != nil {
		return err
	}
	if err := seedCompanyNodeImportance(db); err != nil {
		return err
	}

	log.Println("Seeding complete.")
	return nil
}
