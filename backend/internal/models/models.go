package models

import (
	"time"

	"gorm.io/gorm"
)

// ── Base ─────────────────────────────────────────────────────────────────────

type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ── Users ────────────────────────────────────────────────────────────────────

type User struct {
	Base
	Name           string     `gorm:"not null" json:"name"`
	Email          string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string     `gorm:"not null" json:"-"`
	CurrentLevel   string     `gorm:"default:'beginner'" json:"currentLevel"`
	TargetDeadline *time.Time `json:"targetDeadline,omitempty"`
	CFHandle       string     `json:"cfHandle,omitempty"`
	CFRating       int        `json:"cfRating,omitempty"`
	LCHandle       string     `json:"lcHandle,omitempty"`
	LCSolved       int        `json:"lcSolved,omitempty"`
	AvatarURL      string     `json:"avatarUrl,omitempty"`
	IsAdmin        bool       `gorm:"default:false" json:"isAdmin"`

	// Relations
	Goals           []UserGoal          `gorm:"foreignKey:UserID" json:"goals,omitempty"`
	Attempts        []Attempt           `gorm:"foreignKey:UserID" json:"attempts,omitempty"`
	RevisionSchedules []RevisionSchedule `gorm:"foreignKey:UserID" json:"revisionSchedules,omitempty"`
}

// ── Goals ────────────────────────────────────────────────────────────────────

// CareerGoal defines what the user is targeting (e.g. SWE, Backend, AI)
type CareerGoal struct {
	Base
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IsActive    bool   `gorm:"default:true" json:"isActive"`
}

type UserGoal struct {
	Base
	UserID       uint       `gorm:"not null;index" json:"userId"`
	CareerGoalID uint       `gorm:"not null" json:"careerGoalId"`
	CompanyID    *uint      `json:"companyId,omitempty"`
	RoleID       *uint      `json:"roleId,omitempty"`
	Priority     int        `gorm:"default:1" json:"priority"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"isActive"`

	CareerGoal CareerGoal `gorm:"foreignKey:CareerGoalID" json:"careerGoal,omitempty"`
	Company    *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

// ── Companies ────────────────────────────────────────────────────────────────

type CompanyTier int

const (
	TierOne   CompanyTier = 1 // Rubrik
	TierTwo   CompanyTier = 2 // MAANG + Elite Product
	TierThree CompanyTier = 3 // HFT
	TierFour  CompanyTier = 4 // CP / ICPC
)

type Company struct {
	Base
	Name             string      `gorm:"uniqueIndex;not null" json:"name"`
	Slug             string      `gorm:"uniqueIndex;not null" json:"slug"`
	Tier             CompanyTier `gorm:"not null;default:2" json:"tier"`
	Description      string      `json:"description"`
	LogoURL          string      `json:"logoUrl,omitempty"`
	Website          string      `json:"website,omitempty"`
	HQLocation       string      `json:"hqLocation,omitempty"`
	IsActive         bool        `gorm:"default:true" json:"isActive"`
	IsCustom         bool        `gorm:"default:false" json:"isCustom"`

	// Readiness weights (must sum to 100)
	WeightDSA        int `gorm:"default:35" json:"weightDsa"`
	WeightCoreCS     int `gorm:"default:20" json:"weightCoreCs"`
	WeightDevelopment int `gorm:"default:20" json:"weightDevelopment"`
	WeightProjects   int `gorm:"default:15" json:"weightProjects"`
	WeightBehavioral int `gorm:"default:10" json:"weightBehavioral"`

	// Compensation data (public data only)
	CompRangeMinLPA  int    `json:"compRangeMinLpa,omitempty"`
	CompRangeMaxLPA  int    `json:"compRangeMaxLpa,omitempty"`
	CompCurrency     string `gorm:"default:'INR'" json:"compCurrency"`

	// Relations
	NodeImportance []CompanyNodeImportance `gorm:"foreignKey:CompanyID" json:"nodeImportance,omitempty"`
	Roles          []Role                  `gorm:"many2many:company_roles;" json:"roles,omitempty"`
}

type Role struct {
	Base
	Name        string `gorm:"not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `gorm:"default:true" json:"isActive"`
}

// ── Knowledge Graph ───────────────────────────────────────────────────────────

type NodeType string

const (
	NodeTypeDomain   NodeType = "domain"
	NodeTypeTopic    NodeType = "topic"
	NodeTypeSubtopic NodeType = "subtopic"
	NodeTypeConcept  NodeType = "concept"
	NodeTypePattern  NodeType = "pattern"
)

type DifficultyLevel int

const (
	DifficultyBeginner     DifficultyLevel = 1
	DifficultyEasy         DifficultyLevel = 2
	DifficultyMedium       DifficultyLevel = 3
	DifficultyHard         DifficultyLevel = 4
	DifficultyExpert       DifficultyLevel = 5
)

type KnowledgeNode struct {
	Base
	ParentID         *uint           `gorm:"index" json:"parentId,omitempty"`
	DomainID         *uint           `gorm:"index" json:"domainId,omitempty"`
	Type             NodeType        `gorm:"not null;default:'concept'" json:"type"`
	Name             string          `gorm:"not null" json:"name"`
	Slug             string          `gorm:"uniqueIndex;not null" json:"slug"`
	Description      string          `json:"description"`
	Difficulty       DifficultyLevel `gorm:"default:3" json:"difficulty"`
	EstimatedMinutes int             `gorm:"default:60" json:"estimatedMinutes"`
	RevisionIntervalDays int         `gorm:"default:7" json:"revisionIntervalDays"`
	IsActive         bool            `gorm:"default:true" json:"isActive"`
	SortOrder        int             `gorm:"default:0" json:"sortOrder"`

	// Mastery level metadata
	MasteryLevels    int `gorm:"default:7" json:"masteryLevels"` // 0-7 scale

	// Relations
	Parent           *KnowledgeNode          `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children         []KnowledgeNode         `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Prerequisites    []KnowledgeNode         `gorm:"many2many:node_prerequisites;joinForeignKey:NodeID;joinReferences:PrerequisiteID" json:"prerequisites,omitempty"`
	CompanyImportance []CompanyNodeImportance `gorm:"foreignKey:NodeID" json:"companyImportance,omitempty"`
	Resources        []Resource              `gorm:"foreignKey:NodeID" json:"resources,omitempty"`
}

// node_prerequisites join table
type NodePrerequisite struct {
	NodeID         uint `gorm:"primaryKey"`
	PrerequisiteID uint `gorm:"primaryKey"`
}

// CompanyNodeImportance — maps company to knowledge node with targets
type CompanyNodeImportance struct {
	Base
	CompanyID          uint   `gorm:"not null;index" json:"companyId"`
	NodeID             uint   `gorm:"not null;index" json:"nodeId"`
	ImportanceScore    int    `gorm:"default:3" json:"importanceScore"` // 1-5
	InterviewFrequency int    `gorm:"default:5" json:"interviewFrequency"` // 1-10
	ImportanceLabel    string `gorm:"default:'Medium'" json:"importanceLabel"` // Critical/High/Medium/Low
	MinimumEasy        int    `gorm:"default:0" json:"minimumEasy"`
	MinimumMedium      int    `gorm:"default:5" json:"minimumMedium"`
	MinimumHard        int    `gorm:"default:0" json:"minimumHard"`
	MinimumTotal       int    `gorm:"default:10" json:"minimumTotal"`
	RecommendedTotal   int    `gorm:"default:15" json:"recommendedTotal"`
	Notes              string `json:"notes"`

	Company Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Node    KnowledgeNode  `gorm:"foreignKey:NodeID" json:"node,omitempty"`
}

// ── Problems / Learning Ledger ────────────────────────────────────────────────

type Platform string

const (
	PlatformLeetCode  Platform = "LeetCode"
	PlatformCodeforces Platform = "Codeforces"
	PlatformAtCoder   Platform = "AtCoder"
	PlatformCSES      Platform = "CSES"
	PlatformGFG       Platform = "GFG"
	PlatformOther     Platform = "Other"
)

type ProblemDifficulty string

const (
	DiffEasy   ProblemDifficulty = "Easy"
	DiffMedium ProblemDifficulty = "Medium"
	DiffHard   ProblemDifficulty = "Hard"
)

// Problem — canonical problem catalog
type Problem struct {
	Base
	Platform        Platform          `gorm:"not null" json:"platform"`
	ExternalID      string            `json:"externalId,omitempty"`
	Title           string            `gorm:"not null" json:"title"`
	URL             string            `json:"url,omitempty"`
	Difficulty      ProblemDifficulty `gorm:"not null" json:"difficulty"`
	CFRating        int               `json:"cfRating,omitempty"`
	QualityWeight   int               `gorm:"default:2" json:"qualityWeight"`
	IsActive        bool              `gorm:"default:true" json:"isActive"`

	Tags []KnowledgeNode `gorm:"many2many:problem_node_tags;" json:"tags,omitempty"`
}

// Attempt — user problem solving attempt (Learning Ledger entry)
type Attempt struct {
	Base
	UserID           uint              `gorm:"not null;index" json:"userId"`
	ProblemID        *uint             `gorm:"index" json:"problemId,omitempty"`
	NodeID           *uint             `gorm:"index" json:"nodeId,omitempty"`
	Platform         Platform          `gorm:"not null" json:"platform"`
	ProblemTitle     string            `gorm:"not null" json:"problemTitle"`
	ProblemURL       string            `json:"problemUrl,omitempty"`
	Difficulty       ProblemDifficulty `gorm:"not null" json:"difficulty"`
	CFRating         int               `json:"cfRating,omitempty"`
	Topic            string            `json:"topic"`
	Subtopic         string            `json:"subtopic"`
	Pattern          string            `json:"pattern"`
	TimeTakenMinutes int               `json:"timeTakenMinutes"`
	Result           string            `gorm:"default:'Accepted'" json:"result"`
	ConfidenceScore  int               `gorm:"default:7" json:"confidenceScore"` // 1-10
	RevisionNeeded   bool              `gorm:"default:true" json:"revisionNeeded"`
	MistakeType      string            `json:"mistakeType,omitempty"` // Logic/Syntax/Observation/EdgeCases/DPTransition
	Notes            string            `json:"notes,omitempty"`
	QualityWeight    int               `gorm:"default:2" json:"qualityWeight"`
	AttemptedAt      time.Time         `gorm:"not null" json:"attemptedAt"`

	User    User           `gorm:"foreignKey:UserID" json:"-"`
	Problem *Problem       `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
	Node    *KnowledgeNode `gorm:"foreignKey:NodeID" json:"node,omitempty"`
}

// ── Revision System ───────────────────────────────────────────────────────────

type RevisionStatus string

const (
	RevisionStatusScheduled RevisionStatus = "scheduled"
	RevisionStatusDue       RevisionStatus = "due"
	RevisionStatusOverdue   RevisionStatus = "overdue"
	RevisionStatusMastered  RevisionStatus = "mastered"
	RevisionStatusSkipped   RevisionStatus = "skipped"
)

type RevisionMode string

const (
	RevisionModeNormal     RevisionMode = "normal"
	RevisionModeBlind      RevisionMode = "blind_resolve"
	RevisionModeOral       RevisionMode = "oral_concept"
)

type RevisionSchedule struct {
	Base
	UserID               uint           `gorm:"not null;index" json:"userId"`
	AttemptID            uint           `gorm:"not null;index" json:"attemptId"`
	ProblemTitle         string         `gorm:"not null" json:"problemTitle"`
	Pattern              string         `json:"pattern"`
	CurrentRevisionNum   int            `gorm:"default:1" json:"currentRevisionNum"`
	RequiredRevisions    int            `gorm:"default:3" json:"requiredRevisions"`
	NextRevisionAt       time.Time      `gorm:"not null;index" json:"nextRevisionAt"`
	LastRevisionAt       *time.Time     `json:"lastRevisionAt,omitempty"`
	MemoryEstimate       int            `gorm:"default:70" json:"memoryEstimate"` // 0-100
	Status               RevisionStatus `gorm:"default:'scheduled';index" json:"status"`
	LastRevisionAccuracy int            `json:"lastRevisionAccuracy"`

	User     User              `gorm:"foreignKey:UserID" json:"-"`
	Attempt  Attempt           `gorm:"foreignKey:AttemptID" json:"attempt,omitempty"`
	Sessions []RevisionSession `gorm:"foreignKey:RevisionScheduleID" json:"sessions,omitempty"`
}

type RevisionSession struct {
	Base
	UserID             uint         `gorm:"not null;index" json:"userId"`
	RevisionScheduleID uint         `gorm:"not null;index" json:"revisionScheduleId"`
	AttemptID          uint         `gorm:"not null" json:"attemptId"`
	RevisionNumber     int          `gorm:"not null" json:"revisionNumber"`
	Mode               RevisionMode `gorm:"default:'normal'" json:"mode"`
	Result             string       `gorm:"default:'Solved'" json:"result"`
	TimeTakenMinutes   int          `json:"timeTakenMinutes"`
	ConfidenceScore    int          `json:"confidenceScore"` // 1-10
	AccuracyScore      int          `json:"accuracyScore"`   // 0-100
	NeededHint         bool         `json:"neededHint"`
	Notes              string       `json:"notes,omitempty"`
	RevisedAt          time.Time    `gorm:"not null" json:"revisedAt"`
}

// ── Mastery Records ───────────────────────────────────────────────────────────

type MasteryRecord struct {
	Base
	UserID       uint      `gorm:"not null;index" json:"userId"`
	NodeID       uint      `gorm:"not null;index" json:"nodeId"`
	MasteryLevel int       `gorm:"default:0" json:"masteryLevel"` // 0-7
	MasteryScore float64   `gorm:"default:0" json:"masteryScore"` // 0-100
	LastAccuracy int       `json:"lastAccuracy"`
	MasteredAt   *time.Time `json:"masteredAt,omitempty"`
	Source       string    `json:"source"` // manual/ai/revision

	User User          `gorm:"foreignKey:UserID" json:"-"`
	Node KnowledgeNode `gorm:"foreignKey:NodeID" json:"node,omitempty"`
}

// ── Resources ─────────────────────────────────────────────────────────────────

type ResourceType string

const (
	ResourceTypeArticle  ResourceType = "article"
	ResourceTypeVideo    ResourceType = "video"
	ResourceTypeProblemList ResourceType = "problem_list"
	ResourceTypeNotes    ResourceType = "notes"
	ResourceTypeCourse   ResourceType = "course"
)

type Resource struct {
	Base
	NodeID           uint         `gorm:"not null;index" json:"nodeId"`
	Type             ResourceType `gorm:"not null" json:"type"`
	Title            string       `gorm:"not null" json:"title"`
	URL              string       `gorm:"not null" json:"url"`
	Source           string       `json:"source"` // cp-algorithms, youtube, leetcode, etc.
	EstimatedMinutes int          `json:"estimatedMinutes"`
	QualityScore     int          `gorm:"default:5" json:"qualityScore"` // 1-10
	IsActive         bool         `gorm:"default:true" json:"isActive"`

	Node KnowledgeNode `gorm:"foreignKey:NodeID" json:"node,omitempty"`
}

// ── Readiness ─────────────────────────────────────────────────────────────────

// ReadinessSnapshot — computed daily readiness score per user per company
type ReadinessSnapshot struct {
	Base
	UserID           uint    `gorm:"not null;index" json:"userId"`
	CompanyID        uint    `gorm:"not null;index" json:"companyId"`
	OverallScore     float64 `json:"overallScore"`
	DSAScore         float64 `json:"dsaScore"`
	CoreCSScore      float64 `json:"coreCsScore"`
	DevelopmentScore float64 `json:"developmentScore"`
	ProjectsScore    float64 `json:"projectsScore"`
	BehavioralScore  float64 `json:"behavioralScore"`
	SnapshotDate     time.Time `gorm:"not null;index" json:"snapshotDate"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

// ── Progress Snapshots ────────────────────────────────────────────────────────

type ProgressSnapshot struct {
	Base
	UserID        uint      `gorm:"not null;index" json:"userId"`
	Date          time.Time `gorm:"not null;index" json:"date"`
	TotalSolved   int       `json:"totalSolved"`
	TotalMastered int       `json:"totalMastered"`
	QualityScore  int       `json:"qualityScore"`
	EasySolved    int       `json:"easySolved"`
	MediumSolved  int       `json:"mediumSolved"`
	HardSolved    int       `json:"hardSolved"`
}

// ── Roadmaps ──────────────────────────────────────────────────────────────────

type RoadmapStatus string

const (
	RoadmapStatusActive    RoadmapStatus = "active"
	RoadmapStatusCompleted RoadmapStatus = "completed"
	RoadmapStatusArchived  RoadmapStatus = "archived"
)

type RoadmapItemType string

const (
	RoadmapItemTypeProblem  RoadmapItemType = "problem"
	RoadmapItemTypeRevision RoadmapItemType = "revision"
	RoadmapItemTypeProject  RoadmapItemType = "project"
	RoadmapItemTypeReading  RoadmapItemType = "reading"
	RoadmapItemTypeMock     RoadmapItemType = "mock"
)

type Roadmap struct {
	Base
	UserID      uint          `gorm:"not null;index" json:"userId"`
	Title       string        `gorm:"not null" json:"title"`
	Description string        `json:"description"`
	Status      RoadmapStatus `gorm:"default:'active'" json:"status"`
	GeneratedBy string        `gorm:"default:'ai'" json:"generatedBy"` // ai/manual
	StartsAt    time.Time     `json:"startsAt"`
	EndsAt      *time.Time    `json:"endsAt,omitempty"`

	Items []RoadmapItem `gorm:"foreignKey:RoadmapID" json:"items,omitempty"`
}

type RoadmapItem struct {
	Base
	RoadmapID    uint            `gorm:"not null;index" json:"roadmapId"`
	NodeID       *uint           `json:"nodeId,omitempty"`
	Type         RoadmapItemType `gorm:"not null" json:"type"`
	Title        string          `gorm:"not null" json:"title"`
	Description  string          `json:"description"`
	IsCompleted  bool            `gorm:"default:false" json:"isCompleted"`
	DueDate      *time.Time      `json:"dueDate,omitempty"`
	SortOrder    int             `gorm:"default:0" json:"sortOrder"`
	ImpactScore  float64         `json:"impactScore"`

	Node *KnowledgeNode `gorm:"foreignKey:NodeID" json:"node,omitempty"`
}

// ── Projects ──────────────────────────────────────────────────────────────────

type ProjectStatus string

const (
	ProjectStatusPlanned    ProjectStatus = "planned"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusDeployed   ProjectStatus = "deployed"
)

type Project struct {
	Base
	UserID      uint          `gorm:"not null;index" json:"userId"`
	Title       string        `gorm:"not null" json:"title"`
	Description string        `json:"description"`
	RepoURL     string        `json:"repoUrl,omitempty"`
	LiveURL     string        `json:"liveUrl,omitempty"`
	Status      ProjectStatus `gorm:"default:'planned'" json:"status"`
	TechStack   string        `json:"techStack"` // JSON array stored as string
	StartDate   *time.Time    `json:"startDate,omitempty"`
	EndDate     *time.Time    `json:"endDate,omitempty"`

	Tasks []ProjectTask `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

type ProjectTask struct {
	Base
	ProjectID   uint      `gorm:"not null;index" json:"projectId"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `gorm:"default:false" json:"isCompleted"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	SortOrder   int       `gorm:"default:0" json:"sortOrder"`
}

// ── AI Tables ─────────────────────────────────────────────────────────────────

type AIAnalysis struct {
	Base
	UserID      uint   `gorm:"not null;index" json:"userId"`
	Type        string `gorm:"not null" json:"type"` // roadmap/weakness/resume/interview_sim/code_review
	CompanyID   *uint  `json:"companyId,omitempty"`
	InputData   string `json:"inputData"`  // JSON
	OutputData  string `json:"outputData"` // JSON
	ModelUsed   string `json:"modelUsed"`
	TokensUsed  int    `json:"tokensUsed"`

	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

type AIRecommendation struct {
	Base
	UserID      uint    `gorm:"not null;index" json:"userId"`
	NodeID      *uint   `json:"nodeId,omitempty"`
	CompanyID   *uint   `json:"companyId,omitempty"`
	Title       string  `gorm:"not null" json:"title"`
	Reason      string  `json:"reason"`
	ImpactScore float64 `json:"impactScore"`
	IsActedOn   bool    `gorm:"default:false" json:"isActedOn"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}

// ── Achievements ──────────────────────────────────────────────────────────────

type Achievement struct {
	Base
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Condition   string `json:"condition"` // JSON rule
	Points      int    `gorm:"default:10" json:"points"`
}

type UserAchievement struct {
	Base
	UserID        uint      `gorm:"not null;index" json:"userId"`
	AchievementID uint      `gorm:"not null" json:"achievementId"`
	EarnedAt      time.Time `gorm:"not null" json:"earnedAt"`

	Achievement Achievement `gorm:"foreignKey:AchievementID" json:"achievement,omitempty"`
}

// ── Interview Experience ──────────────────────────────────────────────────────

type InterviewExperience struct {
	Base
	UserID    uint   `gorm:"not null;index" json:"userId"`
	CompanyID uint   `gorm:"not null;index" json:"companyId"`
	Round     string `json:"round"` // OA/Phone/Technical/Final
	Date      time.Time `json:"date"`
	Result    string `json:"result"` // Cleared/Rejected/Pending
	Notes     string `json:"notes"`
	Questions string `json:"questions"` // JSON array

	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}
