package seed

import (
	"puyrg/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// seedKnowledgeGraph seeds all domains, topics, and subtopics
func seedKnowledgeGraph(db *gorm.DB) error {
	// Step 1: seed domains
	domains := domainNodes()
	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).
		Create(&domains).Error; err != nil {
		return err
	}

	// Build slug->ID map
	nodeMap := map[string]uint{}
	var allNodes []models.KnowledgeNode
	db.Find(&allNodes)
	for _, n := range allNodes {
		nodeMap[n.Slug] = n.ID
	}

	// Step 2: seed topics under domains
	topics := topicNodes(nodeMap)
	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).
		Create(&topics).Error; err != nil {
		return err
	}

	// Refresh map
	db.Find(&allNodes)
	for _, n := range allNodes {
		nodeMap[n.Slug] = n.ID
	}

	// Step 3: seed subtopics/patterns
	subtopics := subtopicNodes(nodeMap)
	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).
		Create(&subtopics).Error; err != nil {
		return err
	}

	return nil
}
