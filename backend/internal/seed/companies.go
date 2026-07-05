package seed

import (
	"puyrg/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func seedCompanies(db *gorm.DB) error {
	companies := []models.Company{
		// Tier 1
		{Name: "Rubrik", Slug: "rubrik", Tier: models.TierOne,
			Description:       "Enterprise cloud data management. Strong Go, PostgreSQL, DSA, LLD.",
			WeightDSA: 30, WeightCoreCS: 25, WeightDevelopment: 25, WeightProjects: 10, WeightBehavioral: 10,
			CompRangeMinLPA: 40, CompRangeMaxLPA: 95},

		// Tier 2
		{Name: "Google", Slug: "google", Tier: models.TierTwo,
			Description:       "Algorithm-heavy. Strong DSA required. System design for L4+.",
			WeightDSA: 45, WeightCoreCS: 20, WeightDevelopment: 15, WeightProjects: 10, WeightBehavioral: 10,
			CompRangeMinLPA: 30, CompRangeMaxLPA: 120},
		{Name: "Meta", Slug: "meta", Tier: models.TierTwo,
			Description:       "Coding speed and implementation focused. Strong DSA + behavioral.",
			WeightDSA: 42, WeightCoreCS: 18, WeightDevelopment: 18, WeightProjects: 12, WeightBehavioral: 10,
			CompRangeMinLPA: 30, CompRangeMaxLPA: 110},
		{Name: "Microsoft", Slug: "microsoft", Tier: models.TierTwo,
			Description:       "Balanced DSA, design, and behavioral. CS fundamentals important.",
			WeightDSA: 38, WeightCoreCS: 22, WeightDevelopment: 18, WeightProjects: 12, WeightBehavioral: 10,
			CompRangeMinLPA: 25, CompRangeMaxLPA: 90},
		{Name: "Amazon", Slug: "amazon", Tier: models.TierTwo,
			Description:       "Leadership principles + DSA. Practical problem solving.",
			WeightDSA: 35, WeightCoreCS: 15, WeightDevelopment: 20, WeightProjects: 15, WeightBehavioral: 15,
			CompRangeMinLPA: 25, CompRangeMaxLPA: 90},
		{Name: "Apple", Slug: "apple", Tier: models.TierTwo,
			Description:       "Strong DSA + system design. Quality-focused engineering.",
			WeightDSA: 40, WeightCoreCS: 20, WeightDevelopment: 18, WeightProjects: 12, WeightBehavioral: 10,
			CompRangeMinLPA: 25, CompRangeMaxLPA: 100},
		{Name: "Atlassian", Slug: "atlassian", Tier: models.TierTwo,
			Description:       "DSA + practical system design. Strong engineering culture.",
			WeightDSA: 38, WeightCoreCS: 20, WeightDevelopment: 20, WeightProjects: 12, WeightBehavioral: 10,
			CompRangeMinLPA: 25, CompRangeMaxLPA: 80},
		{Name: "Snowflake", Slug: "snowflake", Tier: models.TierTwo,
			Description:       "Database internals + distributed systems. Strong DSA.",
			WeightDSA: 40, WeightCoreCS: 22, WeightDevelopment: 18, WeightProjects: 10, WeightBehavioral: 10,
			CompRangeMinLPA: 30, CompRangeMaxLPA: 100},
		{Name: "Databricks", Slug: "databricks", Tier: models.TierTwo,
			Description:       "Distributed systems + strong DSA + Spark internals.",
			WeightDSA: 40, WeightCoreCS: 22, WeightDevelopment: 18, WeightProjects: 10, WeightBehavioral: 10,
			CompRangeMinLPA: 30, CompRangeMaxLPA: 100},
		{Name: "Stripe", Slug: "stripe", Tier: models.TierTwo,
			Description:       "Strong engineering + API design + distributed systems.",
			WeightDSA: 38, WeightCoreCS: 20, WeightDevelopment: 22, WeightProjects: 10, WeightBehavioral: 10,
			CompRangeMinLPA: 30, CompRangeMaxLPA: 110},
		{Name: "Uber", Slug: "uber", Tier: models.TierTwo,
			Description:       "Distributed systems + strong DSA + real-time systems.",
			WeightDSA: 40, WeightCoreCS: 20, WeightDevelopment: 20, WeightProjects: 10, WeightBehavioral: 10,
			CompRangeMinLPA: 25, CompRangeMaxLPA: 90},
		{Name: "LinkedIn", Slug: "linkedin", Tier: models.TierTwo,
			Description:       "Graph algorithms + DSA + system design.",
			WeightDSA: 40, WeightCoreCS: 20, WeightDevelopment: 18, WeightProjects: 12, WeightBehavioral: 10,
			CompRangeMinLPA: 25, CompRangeMaxLPA: 85},

		// Tier 3 — HFT
		{Name: "Jane Street", Slug: "jane-street", Tier: models.TierThree,
			Description:       "Pure math, probability, algorithms. OCaml preferred. Very competitive.",
			WeightDSA: 55, WeightCoreCS: 10, WeightDevelopment: 10, WeightProjects: 5, WeightBehavioral: 5,
			CompRangeMinLPA: 60, CompRangeMaxLPA: 200},
		{Name: "Hudson River Trading", Slug: "hrt", Tier: models.TierThree,
			Description:       "Advanced CP + math + low latency systems.",
			WeightDSA: 55, WeightCoreCS: 10, WeightDevelopment: 10, WeightProjects: 5, WeightBehavioral: 5,
			CompRangeMinLPA: 50, CompRangeMaxLPA: 180},
		{Name: "Citadel Securities", Slug: "citadel", Tier: models.TierThree,
			Description:       "Quant trading + math + advanced algorithms.",
			WeightDSA: 50, WeightCoreCS: 10, WeightDevelopment: 15, WeightProjects: 5, WeightBehavioral: 5,
			CompRangeMinLPA: 50, CompRangeMaxLPA: 200},
		{Name: "Tower Research Capital", Slug: "tower-research", Tier: models.TierThree,
			Description:       "Algo trading, low-latency, advanced math and CP.",
			WeightDSA: 55, WeightCoreCS: 10, WeightDevelopment: 10, WeightProjects: 5, WeightBehavioral: 5,
			CompRangeMinLPA: 50, CompRangeMaxLPA: 200},
		{Name: "Optiver", Slug: "optiver", Tier: models.TierThree,
			Description:       "Math, probability, trading algorithms.",
			WeightDSA: 50, WeightCoreCS: 10, WeightDevelopment: 12, WeightProjects: 5, WeightBehavioral: 5,
			CompRangeMinLPA: 45, CompRangeMaxLPA: 150},
		{Name: "IMC", Slug: "imc", Tier: models.TierThree,
			Description:       "Market making, math, CP, low latency.",
			WeightDSA: 50, WeightCoreCS: 10, WeightDevelopment: 12, WeightProjects: 5, WeightBehavioral: 5,
			CompRangeMinLPA: 40, CompRangeMaxLPA: 140},

		// Tier 4 — CP
		{Name: "ICPC", Slug: "icpc", Tier: models.TierFour,
			Description:       "International Collegiate Programming Contest. Pure algorithms.",
			WeightDSA: 80, WeightCoreCS: 10, WeightDevelopment: 5, WeightProjects: 3, WeightBehavioral: 2},
	}

	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).
		Create(&companies).Error
}
