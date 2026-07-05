package service

import (
	"errors"
	"fmt"
	"math"
	"puyrg/backend/internal/models"
	"puyrg/backend/internal/repository"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ── Auth Service ──────────────────────────────────────────────────────────────

type AuthService struct {
	users repository.UserRepository
}

func NewAuthService(users repository.UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	_, err := s.users.GetByEmail(email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := s.users.Create(user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.users.GetByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid email or password")
	}
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

// ── Attempt Service ───────────────────────────────────────────────────────────

type AttemptService struct {
	attempts  repository.AttemptRepository
	revisions repository.RevisionRepository
	mastery   repository.MasteryRepository
	nodes     repository.KnowledgeNodeRepository
}

func NewAttemptService(
	attempts repository.AttemptRepository,
	revisions repository.RevisionRepository,
	mastery repository.MasteryRepository,
	nodes repository.KnowledgeNodeRepository,
) *AttemptService {
	return &AttemptService{
		attempts:  attempts,
		revisions: revisions,
		mastery:   mastery,
		nodes:     nodes,
	}
}

type CreateAttemptInput struct {
	UserID           uint
	Platform         models.Platform
	ProblemTitle     string
	ProblemURL       string
	Difficulty       models.ProblemDifficulty
	CFRating         int
	Topic            string
	Subtopic         string
	Pattern          string
	TimeTakenMinutes int
	Result           string
	ConfidenceScore  int
	RevisionNeeded   bool
	MistakeType      string
	Notes            string
}

func (s *AttemptService) Create(input CreateAttemptInput) (*models.Attempt, *models.RevisionSchedule, error) {
	if strings.TrimSpace(input.ProblemTitle) == "" {
		return nil, nil, errors.New("problem title is required")
	}
	if strings.TrimSpace(input.Topic) == "" {
		return nil, nil, errors.New("topic is required")
	}

	confidence := clamp(input.ConfidenceScore, 1, 10)
	if confidence == 0 {
		confidence = 7
	}

	now := time.Now().UTC()
	attempt := &models.Attempt{
		UserID:           input.UserID,
		Platform:         input.Platform,
		ProblemTitle:     strings.TrimSpace(input.ProblemTitle),
		ProblemURL:       input.ProblemURL,
		Difficulty:       input.Difficulty,
		CFRating:         input.CFRating,
		Topic:            strings.TrimSpace(input.Topic),
		Subtopic:         strings.TrimSpace(input.Subtopic),
		Pattern:          strings.TrimSpace(input.Pattern),
		TimeTakenMinutes: input.TimeTakenMinutes,
		Result:           defaultStr(input.Result, "Accepted"),
		ConfidenceScore:  confidence,
		RevisionNeeded:   input.RevisionNeeded,
		MistakeType:      input.MistakeType,
		Notes:            input.Notes,
		QualityWeight:    qualityWeight(string(input.Platform), string(input.Difficulty), input.CFRating),
		AttemptedAt:      now,
	}

	// Try to link to knowledge node
	if input.Pattern != "" {
		node, err := s.nodes.GetBySlug(slugify(input.Topic + "-" + input.Pattern))
		if err == nil {
			attempt.NodeID = &node.ID
		}
	}

	if err := s.attempts.Create(attempt); err != nil {
		return nil, nil, fmt.Errorf("create attempt: %w", err)
	}

	// Create revision schedule
	var revision *models.RevisionSchedule
	if input.RevisionNeeded {
		revision = &models.RevisionSchedule{
			UserID:             input.UserID,
			AttemptID:          attempt.ID,
			ProblemTitle:       attempt.ProblemTitle,
			Pattern:            attempt.Pattern,
			CurrentRevisionNum: 1,
			RequiredRevisions:  3,
			NextRevisionAt:     now.Add(72 * time.Hour), // R1 after 3 days
			MemoryEstimate:     memoryEstimate(confidence, now.Add(72*time.Hour)),
			Status:             models.RevisionStatusScheduled,
		}
		if err := s.revisions.Create(revision); err != nil {
			return attempt, nil, fmt.Errorf("create revision: %w", err)
		}
	}

	return attempt, revision, nil
}

// ── Revision Service ──────────────────────────────────────────────────────────

type RevisionService struct {
	revisions repository.RevisionRepository
}

func NewRevisionService(revisions repository.RevisionRepository) *RevisionService {
	return &RevisionService{revisions: revisions}
}

type CompleteRevisionInput struct {
	UserID             uint
	RevisionScheduleID uint
	Mode               models.RevisionMode
	Result             string
	TimeTakenMinutes   int
	ConfidenceScore    int
	AccuracyScore      int
	NeededHint         bool
	Notes              string
}

func (s *RevisionService) Complete(input CompleteRevisionInput) (*models.RevisionSession, *models.RevisionSchedule, error) {
	schedule, err := s.revisions.GetByID(input.RevisionScheduleID)
	if err != nil {
		return nil, nil, errors.New("revision schedule not found")
	}
	if schedule.UserID != input.UserID {
		return nil, nil, errors.New("unauthorized")
	}
	if schedule.Status == models.RevisionStatusMastered {
		return nil, nil, errors.New("this problem is already mastered")
	}

	now := time.Now().UTC()
	confidence := clamp(input.ConfidenceScore, 1, 10)
	accuracy := clamp(input.AccuracyScore, 0, 100)

	session := &models.RevisionSession{
		UserID:             input.UserID,
		RevisionScheduleID: schedule.ID,
		AttemptID:          schedule.AttemptID,
		RevisionNumber:     schedule.CurrentRevisionNum,
		Mode:               defaultRevisionMode(input.Mode),
		Result:             defaultStr(input.Result, "Solved"),
		TimeTakenMinutes:   input.TimeTakenMinutes,
		ConfidenceScore:    confidence,
		AccuracyScore:      accuracy,
		NeededHint:         input.NeededHint,
		Notes:              input.Notes,
		RevisedAt:          now,
	}

	if err := s.revisions.CreateSession(session); err != nil {
		return nil, nil, fmt.Errorf("create session: %w", err)
	}

	// Update schedule
	schedule.LastRevisionAt = &now
	schedule.LastRevisionAccuracy = accuracy

	isMastered := schedule.CurrentRevisionNum >= schedule.RequiredRevisions &&
		accuracy >= 80 &&
		!input.NeededHint

	if isMastered {
		schedule.Status = models.RevisionStatusMastered
		schedule.MemoryEstimate = 100
	} else {
		if schedule.CurrentRevisionNum < schedule.RequiredRevisions {
			schedule.CurrentRevisionNum++
		}
		schedule.Status = models.RevisionStatusScheduled
		schedule.NextRevisionAt = nextRevisionDate(now, schedule.CurrentRevisionNum)
		schedule.MemoryEstimate = memoryEstimate(confidence, schedule.NextRevisionAt)
	}

	if err := s.revisions.Update(schedule); err != nil {
		return nil, nil, fmt.Errorf("update schedule: %w", err)
	}

	return session, schedule, nil
}

// ── Dashboard Service ─────────────────────────────────────────────────────────

type DashboardService struct {
	attempts   repository.AttemptRepository
	revisions  repository.RevisionRepository
	mastery    repository.MasteryRepository
	companies  repository.CompanyRepository
	readiness  repository.ReadinessRepository
}

func NewDashboardService(
	attempts repository.AttemptRepository,
	revisions repository.RevisionRepository,
	mastery repository.MasteryRepository,
	companies repository.CompanyRepository,
	readiness repository.ReadinessRepository,
) *DashboardService {
	return &DashboardService{
		attempts:  attempts,
		revisions: revisions,
		mastery:   mastery,
		companies: companies,
		readiness: readiness,
	}
}

type DashboardData struct {
	SolvedQuestions   int64
	MasteredQuestions int64
	NeedRevision      int64
	OverdueRevision   int64
	QualityScore      int
	EasySolved        int64
	MediumSolved      int64
	HardSolved        int64
	RecentAttempts    []models.Attempt
	TodayRevisions    []models.RevisionSchedule
	TopicStats        []repository.TopicStat
	MonthlyTrend      []repository.MonthlyTrend
	CompanyReadiness  []CompanyReadinessSummary
}

type CompanyReadinessSummary struct {
	Company   models.Company
	Score     float64
	Snapshot  *models.ReadinessSnapshot
}

func (s *DashboardService) GetForUser(userID uint) (*DashboardData, error) {
	data := &DashboardData{}
	var err error

	data.SolvedQuestions, err = s.attempts.CountByUser(userID)
	if err != nil {
		return nil, err
	}
	data.MasteredQuestions, err = s.revisions.CountMastered(userID)
	if err != nil {
		return nil, err
	}
	data.NeedRevision, err = s.revisions.CountPending(userID)
	if err != nil {
		return nil, err
	}
	data.QualityScore, err = s.attempts.GetQualityScoreForUser(userID)
	if err != nil {
		return nil, err
	}
	data.EasySolved, _ = s.attempts.CountByUserAndDifficulty(userID, models.DiffEasy)
	data.MediumSolved, _ = s.attempts.CountByUserAndDifficulty(userID, models.DiffMedium)
	data.HardSolved, _ = s.attempts.CountByUserAndDifficulty(userID, models.DiffHard)
	data.RecentAttempts, _ = s.attempts.GetRecentByUser(userID, 5)
	data.TodayRevisions, _ = s.revisions.GetDueToday(userID)
	data.TopicStats, _ = s.attempts.GetTopicStatsForUser(userID)
	data.MonthlyTrend, _ = s.attempts.GetMonthlyTrendForUser(userID)

	// Count overdue
	overdue, _ := s.revisions.GetOverdue(userID)
	data.OverdueRevision = int64(len(overdue))

	// Company readiness
	snapshots, _ := s.readiness.GetAllLatestForUser(userID)
	snapshotMap := map[uint]*models.ReadinessSnapshot{}
	for i := range snapshots {
		snapshotMap[snapshots[i].CompanyID] = &snapshots[i]
	}

	companies, _ := s.companies.GetAll()
	for _, company := range companies {
		summary := CompanyReadinessSummary{Company: company}
		if snap, ok := snapshotMap[company.ID]; ok {
			summary.Score = snap.OverallScore
			summary.Snapshot = snap
		}
		data.CompanyReadiness = append(data.CompanyReadiness, summary)
	}

	return data, nil
}

// ── Readiness Service ─────────────────────────────────────────────────────────

type ReadinessService struct {
	attempts   repository.AttemptRepository
	mastery    repository.MasteryRepository
	companies  repository.CompanyRepository
	importance repository.CompanyNodeImportanceRepository
	readiness  repository.ReadinessRepository
}

func NewReadinessService(
	attempts repository.AttemptRepository,
	mastery repository.MasteryRepository,
	companies repository.CompanyRepository,
	importance repository.CompanyNodeImportanceRepository,
	readiness repository.ReadinessRepository,
) *ReadinessService {
	return &ReadinessService{
		attempts:   attempts,
		mastery:    mastery,
		companies:  companies,
		importance: importance,
		readiness:  readiness,
	}
}

// ComputeAndSave computes readiness score for a user+company and saves snapshot
func (s *ReadinessService) ComputeAndSave(userID, companyID uint) (*models.ReadinessSnapshot, error) {
	company, err := s.companies.GetByID(companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	mappings, err := s.importance.GetByCompany(companyID)
	if err != nil {
		return nil, err
	}

	topicStats, err := s.attempts.GetTopicStatsForUser(userID)
	if err != nil {
		return nil, err
	}

	// Build topic quality map
	topicQuality := map[string]int{}
	for _, stat := range topicStats {
		topicQuality[strings.ToLower(stat.Topic)] = stat.QualityScore
	}

	// Calculate DSA score from mappings
	var totalRequired, totalAchieved float64
	for _, mapping := range mappings {
		required := float64(mapping.MinimumTotal)
		if required == 0 {
			required = 10
		}
		nodeName := strings.ToLower(mapping.Node.Name)
		achieved := math.Min(float64(topicQuality[nodeName]), required)
		totalRequired += required
		totalAchieved += achieved
	}

	dsaScore := 0.0
	if totalRequired > 0 {
		dsaScore = math.Min(100, (totalAchieved/totalRequired)*100)
	}

	// Weighted overall score using company weights
	// For now use DSA score as primary with company weights applied
	overallScore := (dsaScore * float64(company.WeightDSA) / 100)

	snapshot := &models.ReadinessSnapshot{
		UserID:       userID,
		CompanyID:    companyID,
		OverallScore: math.Round(overallScore*10) / 10,
		DSAScore:     math.Round(dsaScore*10) / 10,
		SnapshotDate: time.Now().UTC(),
	}

	if err := s.readiness.Save(snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func qualityWeight(platform, difficulty string, cfRating int) int {
	platform = strings.ToLower(platform)
	difficulty = strings.ToLower(strings.TrimSpace(difficulty))

	if cfRating > 0 {
		switch {
		case cfRating >= 2400:
			return 10
		case cfRating >= 2100:
			return 8
		case cfRating >= 1900:
			return 6
		case cfRating >= 1600:
			return 4
		case cfRating >= 1200:
			return 2
		default:
			return 1
		}
	}

	switch difficulty {
	case "hard":
		return 4
	case "medium":
		return 2
	case "easy":
		return 1
	default:
		if strings.Contains(platform, "icpc") {
			return 8
		}
		return 2
	}
}

func memoryEstimate(confidence int, dueAt time.Time) int {
	base := 50 + confidence*5
	if time.Now().UTC().After(dueAt) {
		base -= 20
	}
	return clamp(base, 10, 100)
}

func nextRevisionDate(now time.Time, revNum int) time.Time {
	switch revNum {
	case 2:
		return now.Add(12 * 24 * time.Hour)
	case 3:
		return now.Add(45 * 24 * time.Hour)
	default:
		return now.Add(3 * 24 * time.Hour)
	}
}

func defaultStr(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func defaultRevisionMode(m models.RevisionMode) models.RevisionMode {
	if m == "" {
		return models.RevisionModeNormal
	}
	return m
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevDash = false
		} else if !prevDash {
			b.WriteRune('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
