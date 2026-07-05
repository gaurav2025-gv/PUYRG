package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"puyrg/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConfigFromEnv reads database config from environment variables
func ConfigFromEnv() Config {
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}
	return Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  sslMode,
	}
}

// DSN returns the PostgreSQL connection string
func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// Connect opens a PostgreSQL connection and assigns it to DB
func Connect(cfg Config) (*gorm.DB, error) {
	logLevel := logger.Warn
	if os.Getenv("APP_ENV") == "development" {
		logLevel = logger.Info
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return db, nil
}

// Migrate runs GORM auto-migration for all models
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		// Identity
		&models.User{},
		&models.CareerGoal{},
		&models.UserGoal{},

		// Companies & Roles
		&models.Company{},
		&models.Role{},

		// Knowledge Graph
		&models.KnowledgeNode{},
		&models.NodePrerequisite{},
		&models.CompanyNodeImportance{},

		// Problems & Learning Ledger
		&models.Problem{},
		&models.Attempt{},

		// Revision System
		&models.RevisionSchedule{},
		&models.RevisionSession{},

		// Mastery & Progress
		&models.MasteryRecord{},
		&models.ProgressSnapshot{},
		&models.ReadinessSnapshot{},

		// Resources
		&models.Resource{},

		// Roadmaps
		&models.Roadmap{},
		&models.RoadmapItem{},

		// Projects
		&models.Project{},
		&models.ProjectTask{},

		// AI
		&models.AIAnalysis{},
		&models.AIRecommendation{},

		// Achievements
		&models.Achievement{},
		&models.UserAchievement{},

		// Interview Experience
		&models.InterviewExperience{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	log.Println("Database migrations complete.")
	return nil
}

// Ping verifies the database connection is alive
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
