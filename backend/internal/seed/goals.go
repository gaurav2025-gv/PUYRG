package seed

import (
	"puyrg/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedCareerGoals(db *gorm.DB) error {
	goals := []models.CareerGoal{
		{Name: "Software Engineer", Slug: "swe", Description: "General SWE track", Icon: "💻"},
		{Name: "Backend Engineer", Slug: "backend", Description: "Go, databases, APIs", Icon: "⚙️"},
		{Name: "AI/ML Engineer", Slug: "ai-ml", Description: "Machine learning and AI systems", Icon: "🤖"},
		{Name: "Competitive Programmer", Slug: "cp", Description: "ICPC, Codeforces, AtCoder", Icon: "🏆"},
		{Name: "Quant Developer", Slug: "quant", Description: "HFT, math, algorithms", Icon: "📈"},
		{Name: "DevOps/SRE", Slug: "devops", Description: "Infrastructure and reliability", Icon: "🔧"},
	}
	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).
		Create(&goals).Error
}
