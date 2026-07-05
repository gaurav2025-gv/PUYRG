package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"puyrg/backend/internal/middleware"
	"puyrg/backend/internal/models"
	"puyrg/backend/internal/repository"
	"puyrg/backend/internal/service"

	"gorm.io/gorm"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func requireUserID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
	}
	return id, ok
}

func pathParam(path, prefix, suffix string) (string, bool) {
	p := strings.TrimPrefix(path, prefix)
	p = strings.TrimSuffix(p, suffix)
	p = strings.Trim(p, "/")
	return p, p != ""
}

func parseUintParam(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	return uint(n), err
}

// ── Health ────────────────────────────────────────────────────────────────────

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "app": "PUYRG"})
}

// ── Auth Handlers ─────────────────────────────────────────────────────────────

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	user, err := h.auth.Register(body.Name, body.Email, body.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"user": user, "token": token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	user, err := h.auth.Login(body.Email, body.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user, "token": token})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	_ = userID
	// Return userID from context — full user fetch happens at service layer if needed
	writeJSON(w, http.StatusOK, map[string]any{"userId": userID})
}

// ── Dashboard Handler ─────────────────────────────────────────────────────────

type DashboardHandler struct {
	dashboard *service.DashboardService
}

func NewDashboardHandler(dashboard *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboard: dashboard}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	data, err := h.dashboard.GetForUser(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load dashboard")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

// ── Company Handlers ──────────────────────────────────────────────────────────

type CompanyHandler struct {
	companies  repository.CompanyRepository
	importance repository.CompanyNodeImportanceRepository
	readiness  *service.ReadinessService
}

func NewCompanyHandler(
	companies repository.CompanyRepository,
	importance repository.CompanyNodeImportanceRepository,
	readiness *service.ReadinessService,
) *CompanyHandler {
	return &CompanyHandler{
		companies:  companies,
		importance: importance,
		readiness:  readiness,
	}
}

func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := h.companies.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch companies")
		return
	}
	writeJSON(w, http.StatusOK, companies)
}

func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr, ok := pathParam(r.URL.Path, "/api/companies/", "")
	if !ok {
		writeError(w, http.StatusBadRequest, "company id required")
		return
	}
	id, err := parseUintParam(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid company id")
		return
	}
	company, err := h.companies.GetWithNodeImportance(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(w, http.StatusNotFound, "company not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, company)
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string              `json:"name"`
		Tier        models.CompanyTier  `json:"tier"`
		Description string              `json:"description"`
		WeightDSA   int                 `json:"weightDsa"`
		WeightCoreCS int                `json:"weightCoreCs"`
		WeightDev   int                 `json:"weightDevelopment"`
		WeightProj  int                 `json:"weightProjects"`
		WeightBeh   int                 `json:"weightBehavioral"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "company name is required")
		return
	}

	company := &models.Company{
		Name:             name,
		Slug:             slugify(name),
		Tier:             body.Tier,
		Description:      body.Description,
		IsCustom:         true,
		WeightDSA:        defaultInt(body.WeightDSA, 35),
		WeightCoreCS:     defaultInt(body.WeightCoreCS, 20),
		WeightDevelopment: defaultInt(body.WeightDev, 20),
		WeightProjects:   defaultInt(body.WeightProj, 15),
		WeightBehavioral: defaultInt(body.WeightBeh, 10),
	}
	if err := h.companies.Create(company); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create company")
		return
	}
	writeJSON(w, http.StatusCreated, company)
}

func (h *CompanyHandler) GetReadiness(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	// Extract company ID from path /api/companies/{id}/readiness
	path := strings.TrimSuffix(r.URL.Path, "/readiness")
	idStr, ok2 := pathParam(path, "/api/companies/", "")
	if !ok2 {
		writeError(w, http.StatusBadRequest, "company id required")
		return
	}
	id, err := parseUintParam(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid company id")
		return
	}
	snapshot, err := h.readiness.ComputeAndSave(userID, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

// ── Knowledge Node Handlers ───────────────────────────────────────────────────

type KnowledgeHandler struct {
	nodes repository.KnowledgeNodeRepository
}

func NewKnowledgeHandler(nodes repository.KnowledgeNodeRepository) *KnowledgeHandler {
	return &KnowledgeHandler{nodes: nodes}
}

func (h *KnowledgeHandler) ListDomains(w http.ResponseWriter, _ *http.Request) {
	domains, err := h.nodes.GetDomains()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch domains")
		return
	}
	writeJSON(w, http.StatusOK, domains)
}

func (h *KnowledgeHandler) ListAll(w http.ResponseWriter, _ *http.Request) {
	nodes, err := h.nodes.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch nodes")
		return
	}
	writeJSON(w, http.StatusOK, nodes)
}

func (h *KnowledgeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr, ok := pathParam(r.URL.Path, "/api/knowledge/nodes/", "")
	if !ok {
		writeError(w, http.StatusBadRequest, "node id required")
		return
	}
	id, err := parseUintParam(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid node id")
		return
	}
	node, err := h.nodes.GetWithChildren(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeError(w, http.StatusNotFound, "node not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, node)
}

func (h *KnowledgeHandler) GetChildren(w http.ResponseWriter, r *http.Request) {
	// /api/knowledge/nodes/{id}/children
	path := strings.TrimSuffix(r.URL.Path, "/children")
	idStr, ok := pathParam(path, "/api/knowledge/nodes/", "")
	if !ok {
		writeError(w, http.StatusBadRequest, "node id required")
		return
	}
	id, err := parseUintParam(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid node id")
		return
	}
	children, err := h.nodes.GetChildren(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, children)
}

func (h *KnowledgeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ParentID         *uint                  `json:"parentId"`
		Type             models.NodeType        `json:"type"`
		Name             string                 `json:"name"`
		Description      string                 `json:"description"`
		Difficulty       models.DifficultyLevel `json:"difficulty"`
		EstimatedMinutes int                    `json:"estimatedMinutes"`
		RevisionInterval int                    `json:"revisionIntervalDays"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	node := &models.KnowledgeNode{
		ParentID:             body.ParentID,
		Type:                 body.Type,
		Name:                 name,
		Slug:                 slugify(name),
		Description:          body.Description,
		Difficulty:           defaultDifficulty(body.Difficulty),
		EstimatedMinutes:     defaultInt(body.EstimatedMinutes, 60),
		RevisionIntervalDays: defaultInt(body.RevisionInterval, 7),
	}
	if err := h.nodes.Create(node); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create node")
		return
	}
	writeJSON(w, http.StatusCreated, node)
}

// ── Attempt Handlers ──────────────────────────────────────────────────────────

type AttemptHandler struct {
	svc *service.AttemptService
	repo repository.AttemptRepository
}

func NewAttemptHandler(svc *service.AttemptService, repo repository.AttemptRepository) *AttemptHandler {
	return &AttemptHandler{svc: svc, repo: repo}
}

func (h *AttemptHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	attempts, err := h.repo.GetRecentByUser(userID, 50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, attempts)
}

func (h *AttemptHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	var body struct {
		Platform         string `json:"platform"`
		ProblemTitle     string `json:"problemTitle"`
		ProblemURL       string `json:"problemUrl"`
		Difficulty       string `json:"difficulty"`
		CFRating         int    `json:"cfRating"`
		Topic            string `json:"topic"`
		Subtopic         string `json:"subtopic"`
		Pattern          string `json:"pattern"`
		TimeTakenMinutes int    `json:"timeTakenMinutes"`
		Result           string `json:"result"`
		ConfidenceScore  int    `json:"confidenceScore"`
		RevisionNeeded   bool   `json:"revisionNeeded"`
		MistakeType      string `json:"mistakeType"`
		Notes            string `json:"notes"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	input := service.CreateAttemptInput{
		UserID:           userID,
		Platform:         models.Platform(defaultStr(body.Platform, "LeetCode")),
		ProblemTitle:     body.ProblemTitle,
		ProblemURL:       body.ProblemURL,
		Difficulty:       models.ProblemDifficulty(defaultStr(body.Difficulty, "Medium")),
		CFRating:         body.CFRating,
		Topic:            body.Topic,
		Subtopic:         body.Subtopic,
		Pattern:          body.Pattern,
		TimeTakenMinutes: body.TimeTakenMinutes,
		Result:           body.Result,
		ConfidenceScore:  body.ConfidenceScore,
		RevisionNeeded:   body.RevisionNeeded,
		MistakeType:      body.MistakeType,
		Notes:            body.Notes,
	}

	attempt, revision, err := h.svc.Create(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"attempt":  attempt,
		"revision": revision,
	})
}

// ── Revision Handlers ─────────────────────────────────────────────────────────

type RevisionHandler struct {
	svc  *service.RevisionService
	repo repository.RevisionRepository
}

func NewRevisionHandler(svc *service.RevisionService, repo repository.RevisionRepository) *RevisionHandler {
	return &RevisionHandler{svc: svc, repo: repo}
}

func (h *RevisionHandler) Today(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}
	revisions, err := h.repo.GetDueToday(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, revisions)
}

func (h *RevisionHandler) Complete(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	// /api/revisions/{id}/sessions
	path := strings.TrimSuffix(r.URL.Path, "/sessions")
	idStr, ok2 := pathParam(path, "/api/revisions/", "")
	if !ok2 {
		writeError(w, http.StatusBadRequest, "revision id required")
		return
	}
	id, err := parseUintParam(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid revision id")
		return
	}

	var body struct {
		Mode             string `json:"mode"`
		Result           string `json:"result"`
		TimeTakenMinutes int    `json:"timeTakenMinutes"`
		ConfidenceScore  int    `json:"confidenceScore"`
		AccuracyScore    int    `json:"accuracyScore"`
		NeededHint       bool   `json:"neededHint"`
		Notes            string `json:"notes"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	input := service.CompleteRevisionInput{
		UserID:             userID,
		RevisionScheduleID: id,
		Mode:               models.RevisionMode(body.Mode),
		Result:             body.Result,
		TimeTakenMinutes:   body.TimeTakenMinutes,
		ConfidenceScore:    body.ConfidenceScore,
		AccuracyScore:      body.AccuracyScore,
		NeededHint:         body.NeededHint,
		Notes:              body.Notes,
	}

	session, schedule, err := h.svc.Complete(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"session":  session,
		"schedule": schedule,
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func defaultStr(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func defaultInt(v, fallback int) int {
	if v <= 0 {
		return fallback
	}
	return v
}

func defaultDifficulty(d models.DifficultyLevel) models.DifficultyLevel {
	if d == 0 {
		return models.DifficultyMedium
	}
	return d
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
		} else if !prevDash {
			b.WriteRune('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
