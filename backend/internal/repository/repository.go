package repository

import (
	"puyrg/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// ── User ──────────────────────────────────────────────────────────────────────

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepo{db} }

func (r *userRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}
func (r *userRepo) GetByID(id uint) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, id).Error
	return &u, err
}
func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ?", email).First(&u).Error
	return &u, err
}
func (r *userRepo) Update(u *models.User) error { return r.db.Save(u).Error }
func (r *userRepo) Delete(id uint) error        { return r.db.Delete(&models.User{}, id).Error }

// ── Company ───────────────────────────────────────────────────────────────────

type CompanyRepository interface {
	Create(company *models.Company) error
	GetByID(id uint) (*models.Company, error)
	GetBySlug(slug string) (*models.Company, error)
	GetAll() ([]models.Company, error)
	GetByTier(tier models.CompanyTier) ([]models.Company, error)
	Update(company *models.Company) error
	Delete(id uint) error
	GetWithNodeImportance(id uint) (*models.Company, error)
}

type companyRepo struct{ db *gorm.DB }

func NewCompanyRepository(db *gorm.DB) CompanyRepository { return &companyRepo{db} }

func (r *companyRepo) Create(c *models.Company) error { return r.db.Create(c).Error }
func (r *companyRepo) GetByID(id uint) (*models.Company, error) {
	var c models.Company
	err := r.db.First(&c, id).Error
	return &c, err
}
func (r *companyRepo) GetBySlug(slug string) (*models.Company, error) {
	var c models.Company
	err := r.db.Where("slug = ?", slug).First(&c).Error
	return &c, err
}
func (r *companyRepo) GetAll() ([]models.Company, error) {
	var companies []models.Company
	err := r.db.Where("is_active = true").Order("tier, name").Find(&companies).Error
	return companies, err
}
func (r *companyRepo) GetByTier(tier models.CompanyTier) ([]models.Company, error) {
	var companies []models.Company
	err := r.db.Where("tier = ? AND is_active = true", tier).Order("name").Find(&companies).Error
	return companies, err
}
func (r *companyRepo) Update(c *models.Company) error { return r.db.Save(c).Error }
func (r *companyRepo) Delete(id uint) error           { return r.db.Delete(&models.Company{}, id).Error }
func (r *companyRepo) GetWithNodeImportance(id uint) (*models.Company, error) {
	var c models.Company
	err := r.db.Preload("NodeImportance.Node").First(&c, id).Error
	return &c, err
}

// ── KnowledgeNode ─────────────────────────────────────────────────────────────

type KnowledgeNodeRepository interface {
	Create(node *models.KnowledgeNode) error
	GetByID(id uint) (*models.KnowledgeNode, error)
	GetBySlug(slug string) (*models.KnowledgeNode, error)
	GetAll() ([]models.KnowledgeNode, error)
	GetByType(nodeType models.NodeType) ([]models.KnowledgeNode, error)
	GetChildren(parentID uint) ([]models.KnowledgeNode, error)
	GetDomains() ([]models.KnowledgeNode, error)
	GetWithChildren(id uint) (*models.KnowledgeNode, error)
	GetWithPrerequisites(id uint) (*models.KnowledgeNode, error)
	Update(node *models.KnowledgeNode) error
	Delete(id uint) error
	BulkCreate(nodes []models.KnowledgeNode) error
	GetByCompany(companyID uint) ([]models.KnowledgeNode, error)
}

type knowledgeNodeRepo struct{ db *gorm.DB }

func NewKnowledgeNodeRepository(db *gorm.DB) KnowledgeNodeRepository {
	return &knowledgeNodeRepo{db}
}

func (r *knowledgeNodeRepo) Create(n *models.KnowledgeNode) error { return r.db.Create(n).Error }
func (r *knowledgeNodeRepo) GetByID(id uint) (*models.KnowledgeNode, error) {
	var n models.KnowledgeNode
	err := r.db.First(&n, id).Error
	return &n, err
}
func (r *knowledgeNodeRepo) GetBySlug(slug string) (*models.KnowledgeNode, error) {
	var n models.KnowledgeNode
	err := r.db.Where("slug = ?", slug).First(&n).Error
	return &n, err
}
func (r *knowledgeNodeRepo) GetAll() ([]models.KnowledgeNode, error) {
	var nodes []models.KnowledgeNode
	err := r.db.Where("is_active = true").Order("sort_order, name").Find(&nodes).Error
	return nodes, err
}
func (r *knowledgeNodeRepo) GetByType(t models.NodeType) ([]models.KnowledgeNode, error) {
	var nodes []models.KnowledgeNode
	err := r.db.Where("type = ? AND is_active = true", t).Order("sort_order, name").Find(&nodes).Error
	return nodes, err
}
func (r *knowledgeNodeRepo) GetChildren(parentID uint) ([]models.KnowledgeNode, error) {
	var nodes []models.KnowledgeNode
	err := r.db.Where("parent_id = ? AND is_active = true", parentID).Order("sort_order, name").Find(&nodes).Error
	return nodes, err
}
func (r *knowledgeNodeRepo) GetDomains() ([]models.KnowledgeNode, error) {
	return r.GetByType(models.NodeTypeDomain)
}
func (r *knowledgeNodeRepo) GetWithChildren(id uint) (*models.KnowledgeNode, error) {
	var n models.KnowledgeNode
	err := r.db.Preload("Children").First(&n, id).Error
	return &n, err
}
func (r *knowledgeNodeRepo) GetWithPrerequisites(id uint) (*models.KnowledgeNode, error) {
	var n models.KnowledgeNode
	err := r.db.Preload("Prerequisites").First(&n, id).Error
	return &n, err
}
func (r *knowledgeNodeRepo) Update(n *models.KnowledgeNode) error { return r.db.Save(n).Error }
func (r *knowledgeNodeRepo) Delete(id uint) error {
	return r.db.Delete(&models.KnowledgeNode{}, id).Error
}
func (r *knowledgeNodeRepo) BulkCreate(nodes []models.KnowledgeNode) error {
	return r.db.CreateInBatches(nodes, 100).Error
}
func (r *knowledgeNodeRepo) GetByCompany(companyID uint) ([]models.KnowledgeNode, error) {
	var nodes []models.KnowledgeNode
	err := r.db.
		Joins("JOIN company_node_importances ON company_node_importances.node_id = knowledge_nodes.id").
		Where("company_node_importances.company_id = ? AND knowledge_nodes.is_active = true", companyID).
		Order("company_node_importances.importance_score DESC").
		Find(&nodes).Error
	return nodes, err
}

// ── CompanyNodeImportance ─────────────────────────────────────────────────────

type CompanyNodeImportanceRepository interface {
	Create(mapping *models.CompanyNodeImportance) error
	GetByCompany(companyID uint) ([]models.CompanyNodeImportance, error)
	GetByCompanyAndNode(companyID, nodeID uint) (*models.CompanyNodeImportance, error)
	Update(mapping *models.CompanyNodeImportance) error
	Delete(id uint) error
	BulkCreate(mappings []models.CompanyNodeImportance) error
	BulkUpsert(mappings []models.CompanyNodeImportance) error
}

type companyNodeImportanceRepo struct{ db *gorm.DB }

func NewCompanyNodeImportanceRepository(db *gorm.DB) CompanyNodeImportanceRepository {
	return &companyNodeImportanceRepo{db}
}

func (r *companyNodeImportanceRepo) Create(m *models.CompanyNodeImportance) error {
	return r.db.Create(m).Error
}
func (r *companyNodeImportanceRepo) GetByCompany(companyID uint) ([]models.CompanyNodeImportance, error) {
	var mappings []models.CompanyNodeImportance
	err := r.db.Preload("Node").Where("company_id = ?", companyID).
		Order("importance_score DESC").Find(&mappings).Error
	return mappings, err
}
func (r *companyNodeImportanceRepo) GetByCompanyAndNode(companyID, nodeID uint) (*models.CompanyNodeImportance, error) {
	var m models.CompanyNodeImportance
	err := r.db.Where("company_id = ? AND node_id = ?", companyID, nodeID).First(&m).Error
	return &m, err
}
func (r *companyNodeImportanceRepo) Update(m *models.CompanyNodeImportance) error {
	return r.db.Save(m).Error
}
func (r *companyNodeImportanceRepo) Delete(id uint) error {
	return r.db.Delete(&models.CompanyNodeImportance{}, id).Error
}
func (r *companyNodeImportanceRepo) BulkCreate(mappings []models.CompanyNodeImportance) error {
	return r.db.CreateInBatches(mappings, 100).Error
}
func (r *companyNodeImportanceRepo) BulkUpsert(mappings []models.CompanyNodeImportance) error {
	return r.db.Save(mappings).Error
}

// ── Attempt ───────────────────────────────────────────────────────────────────

type AttemptRepository interface {
	Create(attempt *models.Attempt) error
	GetByID(id uint) (*models.Attempt, error)
	GetByUser(userID uint, limit, offset int) ([]models.Attempt, error)
	GetByUserAndTopic(userID uint, topic string) ([]models.Attempt, error)
	GetByUserAndNode(userID, nodeID uint) ([]models.Attempt, error)
	GetRecentByUser(userID uint, limit int) ([]models.Attempt, error)
	CountByUser(userID uint) (int64, error)
	CountByUserAndDifficulty(userID uint, difficulty models.ProblemDifficulty) (int64, error)
	GetTopicStatsForUser(userID uint) ([]TopicStat, error)
	GetQualityScoreForUser(userID uint) (int, error)
	GetMonthlyTrendForUser(userID uint) ([]MonthlyTrend, error)
}

type TopicStat struct {
	Topic        string
	Solved       int
	QualityScore int
}

type MonthlyTrend struct {
	Month  string
	Solved int
}

type attemptRepo struct{ db *gorm.DB }

func NewAttemptRepository(db *gorm.DB) AttemptRepository { return &attemptRepo{db} }

func (r *attemptRepo) Create(a *models.Attempt) error { return r.db.Create(a).Error }
func (r *attemptRepo) GetByID(id uint) (*models.Attempt, error) {
	var a models.Attempt
	err := r.db.Preload("Node").First(&a, id).Error
	return &a, err
}
func (r *attemptRepo) GetByUser(userID uint, limit, offset int) ([]models.Attempt, error) {
	var attempts []models.Attempt
	err := r.db.Where("user_id = ?", userID).
		Order("attempted_at DESC").
		Limit(limit).Offset(offset).Find(&attempts).Error
	return attempts, err
}
func (r *attemptRepo) GetByUserAndTopic(userID uint, topic string) ([]models.Attempt, error) {
	var attempts []models.Attempt
	err := r.db.Where("user_id = ? AND topic = ?", userID, topic).
		Order("attempted_at DESC").Find(&attempts).Error
	return attempts, err
}
func (r *attemptRepo) GetByUserAndNode(userID, nodeID uint) ([]models.Attempt, error) {
	var attempts []models.Attempt
	err := r.db.Where("user_id = ? AND node_id = ?", userID, nodeID).
		Order("attempted_at DESC").Find(&attempts).Error
	return attempts, err
}
func (r *attemptRepo) GetRecentByUser(userID uint, limit int) ([]models.Attempt, error) {
	return r.GetByUser(userID, limit, 0)
}
func (r *attemptRepo) CountByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Attempt{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
func (r *attemptRepo) CountByUserAndDifficulty(userID uint, diff models.ProblemDifficulty) (int64, error) {
	var count int64
	err := r.db.Model(&models.Attempt{}).
		Where("user_id = ? AND difficulty = ?", userID, diff).Count(&count).Error
	return count, err
}
func (r *attemptRepo) GetTopicStatsForUser(userID uint) ([]TopicStat, error) {
	var stats []TopicStat
	err := r.db.Model(&models.Attempt{}).
		Select("topic, COUNT(*) as solved, SUM(quality_weight) as quality_score").
		Where("user_id = ?", userID).
		Group("topic").
		Order("solved DESC").
		Scan(&stats).Error
	return stats, err
}
func (r *attemptRepo) GetQualityScoreForUser(userID uint) (int, error) {
	var total int
	err := r.db.Model(&models.Attempt{}).
		Select("COALESCE(SUM(quality_weight), 0)").
		Where("user_id = ?", userID).
		Scan(&total).Error
	return total, err
}
func (r *attemptRepo) GetMonthlyTrendForUser(userID uint) ([]MonthlyTrend, error) {
	var trend []MonthlyTrend
	err := r.db.Model(&models.Attempt{}).
		Select("TO_CHAR(attempted_at, 'Mon YYYY') as month, COUNT(*) as solved").
		Where("user_id = ?", userID).
		Group("TO_CHAR(attempted_at, 'Mon YYYY')").
		Order("MIN(attempted_at)").
		Scan(&trend).Error
	return trend, err
}

// ── RevisionSchedule ──────────────────────────────────────────────────────────

type RevisionRepository interface {
	Create(schedule *models.RevisionSchedule) error
	GetByID(id uint) (*models.RevisionSchedule, error)
	GetDueToday(userID uint) ([]models.RevisionSchedule, error)
	GetOverdue(userID uint) ([]models.RevisionSchedule, error)
	GetByUser(userID uint) ([]models.RevisionSchedule, error)
	GetPending(userID uint) ([]models.RevisionSchedule, error)
	Update(schedule *models.RevisionSchedule) error
	CountMastered(userID uint) (int64, error)
	CountPending(userID uint) (int64, error)
	CreateSession(session *models.RevisionSession) error
}

type revisionRepo struct{ db *gorm.DB }

func NewRevisionRepository(db *gorm.DB) RevisionRepository { return &revisionRepo{db} }

func (r *revisionRepo) Create(s *models.RevisionSchedule) error { return r.db.Create(s).Error }
func (r *revisionRepo) GetByID(id uint) (*models.RevisionSchedule, error) {
	var s models.RevisionSchedule
	err := r.db.Preload("Sessions").First(&s, id).Error
	return &s, err
}
func (r *revisionRepo) GetDueToday(userID uint) ([]models.RevisionSchedule, error) {
	var schedules []models.RevisionSchedule
	now := time.Now().UTC()
	err := r.db.Where("user_id = ? AND status != ? AND next_revision_at <= ?",
		userID, models.RevisionStatusMastered, now.Add(24*time.Hour)).
		Order("next_revision_at ASC").Find(&schedules).Error
	return schedules, err
}
func (r *revisionRepo) GetOverdue(userID uint) ([]models.RevisionSchedule, error) {
	var schedules []models.RevisionSchedule
	err := r.db.Where("user_id = ? AND status != ? AND next_revision_at < ?",
		userID, models.RevisionStatusMastered, time.Now().UTC()).
		Order("next_revision_at ASC").Find(&schedules).Error
	return schedules, err
}
func (r *revisionRepo) GetByUser(userID uint) ([]models.RevisionSchedule, error) {
	var schedules []models.RevisionSchedule
	err := r.db.Where("user_id = ?", userID).
		Order("next_revision_at ASC").Find(&schedules).Error
	return schedules, err
}
func (r *revisionRepo) GetPending(userID uint) ([]models.RevisionSchedule, error) {
	var schedules []models.RevisionSchedule
	err := r.db.Where("user_id = ? AND status != ?", userID, models.RevisionStatusMastered).
		Find(&schedules).Error
	return schedules, err
}
func (r *revisionRepo) Update(s *models.RevisionSchedule) error { return r.db.Save(s).Error }
func (r *revisionRepo) CountMastered(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.RevisionSchedule{}).
		Where("user_id = ? AND status = ?", userID, models.RevisionStatusMastered).
		Count(&count).Error
	return count, err
}
func (r *revisionRepo) CountPending(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.RevisionSchedule{}).
		Where("user_id = ? AND status != ?", userID, models.RevisionStatusMastered).
		Count(&count).Error
	return count, err
}
func (r *revisionRepo) CreateSession(session *models.RevisionSession) error {
	return r.db.Create(session).Error
}

// ── MasteryRecord ─────────────────────────────────────────────────────────────

type MasteryRepository interface {
	Upsert(record *models.MasteryRecord) error
	GetByUserAndNode(userID, nodeID uint) (*models.MasteryRecord, error)
	GetByUser(userID uint) ([]models.MasteryRecord, error)
	GetWeakTopics(userID uint, limit int) ([]models.MasteryRecord, error)
}

type masteryRepo struct{ db *gorm.DB }

func NewMasteryRepository(db *gorm.DB) MasteryRepository { return &masteryRepo{db} }

func (r *masteryRepo) Upsert(record *models.MasteryRecord) error {
	return r.db.Save(record).Error
}
func (r *masteryRepo) GetByUserAndNode(userID, nodeID uint) (*models.MasteryRecord, error) {
	var m models.MasteryRecord
	err := r.db.Where("user_id = ? AND node_id = ?", userID, nodeID).First(&m).Error
	return &m, err
}
func (r *masteryRepo) GetByUser(userID uint) ([]models.MasteryRecord, error) {
	var records []models.MasteryRecord
	err := r.db.Preload("Node").Where("user_id = ?", userID).
		Order("mastery_level DESC").Find(&records).Error
	return records, err
}
func (r *masteryRepo) GetWeakTopics(userID uint, limit int) ([]models.MasteryRecord, error) {
	var records []models.MasteryRecord
	err := r.db.Preload("Node").Where("user_id = ?", userID).
		Order("mastery_level ASC, mastery_score ASC").
		Limit(limit).Find(&records).Error
	return records, err
}

// ── ReadinessSnapshot ─────────────────────────────────────────────────────────

type ReadinessRepository interface {
	Save(snapshot *models.ReadinessSnapshot) error
	GetLatest(userID, companyID uint) (*models.ReadinessSnapshot, error)
	GetHistory(userID, companyID uint, days int) ([]models.ReadinessSnapshot, error)
	GetAllLatestForUser(userID uint) ([]models.ReadinessSnapshot, error)
}

type readinessRepo struct{ db *gorm.DB }

func NewReadinessRepository(db *gorm.DB) ReadinessRepository { return &readinessRepo{db} }

func (r *readinessRepo) Save(s *models.ReadinessSnapshot) error { return r.db.Create(s).Error }
func (r *readinessRepo) GetLatest(userID, companyID uint) (*models.ReadinessSnapshot, error) {
	var s models.ReadinessSnapshot
	err := r.db.Where("user_id = ? AND company_id = ?", userID, companyID).
		Order("snapshot_date DESC").First(&s).Error
	return &s, err
}
func (r *readinessRepo) GetHistory(userID, companyID uint, days int) ([]models.ReadinessSnapshot, error) {
	var snapshots []models.ReadinessSnapshot
	since := time.Now().UTC().AddDate(0, 0, -days)
	err := r.db.Where("user_id = ? AND company_id = ? AND snapshot_date >= ?",
		userID, companyID, since).
		Order("snapshot_date ASC").Find(&snapshots).Error
	return snapshots, err
}
func (r *readinessRepo) GetAllLatestForUser(userID uint) ([]models.ReadinessSnapshot, error) {
	var snapshots []models.ReadinessSnapshot
	err := r.db.Raw(`
		SELECT DISTINCT ON (company_id) *
		FROM readiness_snapshots
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY company_id, snapshot_date DESC
	`, userID).Scan(&snapshots).Error
	return snapshots, err
}

// ── Resource ──────────────────────────────────────────────────────────────────

type ResourceRepository interface {
	Create(resource *models.Resource) error
	GetByNode(nodeID uint) ([]models.Resource, error)
	GetByNodeAndType(nodeID uint, resourceType models.ResourceType) ([]models.Resource, error)
	Update(resource *models.Resource) error
	Delete(id uint) error
	BulkCreate(resources []models.Resource) error
}

type resourceRepo struct{ db *gorm.DB }

func NewResourceRepository(db *gorm.DB) ResourceRepository { return &resourceRepo{db} }

func (r *resourceRepo) Create(res *models.Resource) error { return r.db.Create(res).Error }
func (r *resourceRepo) GetByNode(nodeID uint) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Where("node_id = ? AND is_active = true", nodeID).
		Order("quality_score DESC").Find(&resources).Error
	return resources, err
}
func (r *resourceRepo) GetByNodeAndType(nodeID uint, t models.ResourceType) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Where("node_id = ? AND type = ? AND is_active = true", nodeID, t).
		Order("quality_score DESC").Find(&resources).Error
	return resources, err
}
func (r *resourceRepo) Update(res *models.Resource) error { return r.db.Save(res).Error }
func (r *resourceRepo) Delete(id uint) error              { return r.db.Delete(&models.Resource{}, id).Error }
func (r *resourceRepo) BulkCreate(resources []models.Resource) error {
	return r.db.CreateInBatches(resources, 100).Error
}

// ── Project ───────────────────────────────────────────────────────────────────

type ProjectRepository interface {
	Create(project *models.Project) error
	GetByUser(userID uint) ([]models.Project, error)
	GetByIDAndUser(id, userID uint) (*models.Project, error)
	Update(project *models.Project) error
	Delete(id uint) error
	CreateTask(task *models.ProjectTask) error
	UpdateTask(task *models.ProjectTask) error
}

type projectRepo struct{ db *gorm.DB }

func NewProjectRepository(db *gorm.DB) ProjectRepository { return &projectRepo{db} }

func (r *projectRepo) Create(p *models.Project) error { return r.db.Create(p).Error }
func (r *projectRepo) GetByUser(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Preload("Tasks").Where("user_id = ?", userID).
		Order("created_at DESC").Find(&projects).Error
	return projects, err
}
func (r *projectRepo) GetByIDAndUser(id, userID uint) (*models.Project, error) {
	var p models.Project
	err := r.db.Preload("Tasks").Where("id = ? AND user_id = ?", id, userID).First(&p).Error
	return &p, err
}
func (r *projectRepo) Update(p *models.Project) error { return r.db.Save(p).Error }
func (r *projectRepo) Delete(id uint) error           { return r.db.Delete(&models.Project{}, id).Error }
func (r *projectRepo) CreateTask(t *models.ProjectTask) error { return r.db.Create(t).Error }
func (r *projectRepo) UpdateTask(t *models.ProjectTask) error { return r.db.Save(t).Error }
