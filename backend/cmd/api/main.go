package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"puyrg/backend/internal/database"
	"puyrg/backend/internal/handler"
	"puyrg/backend/internal/middleware"
	"puyrg/backend/internal/models"
	"puyrg/backend/internal/repository"
	"puyrg/backend/internal/seed"
	"puyrg/backend/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to PostgreSQL
	cfg := database.ConfigFromEnv()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	if err := database.Ping(db); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("Connected to PostgreSQL ✓")

	// Auto-migrate all tables
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Seed initial data (idempotent)
	if err := seed.Run(db); err != nil {
		log.Printf("Seed warning: %v", err)
	}

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo       := repository.NewUserRepository(db)
	companyRepo    := repository.NewCompanyRepository(db)
	nodeRepo       := repository.NewKnowledgeNodeRepository(db)
	importanceRepo := repository.NewCompanyNodeImportanceRepository(db)
	attemptRepo    := repository.NewAttemptRepository(db)
	revisionRepo   := repository.NewRevisionRepository(db)
	masteryRepo    := repository.NewMasteryRepository(db)
	readinessRepo  := repository.NewReadinessRepository(db)
	projectRepo    := repository.NewProjectRepository(db)
	resourceRepo   := repository.NewResourceRepository(db)

	_ = projectRepo
	_ = resourceRepo

	// ── Services ──────────────────────────────────────────────────────────────
	authSvc      := service.NewAuthService(userRepo)
	attemptSvc   := service.NewAttemptService(attemptRepo, revisionRepo, masteryRepo, nodeRepo)
	revisionSvc  := service.NewRevisionService(revisionRepo)
	readinessSvc := service.NewReadinessService(attemptRepo, masteryRepo, companyRepo, importanceRepo, readinessRepo)
	dashboardSvc := service.NewDashboardService(attemptRepo, revisionRepo, masteryRepo, companyRepo, readinessRepo)

	// ── Handlers ──────────────────────────────────────────────────────────────
	authH      := handler.NewAuthHandler(authSvc)
	dashH      := handler.NewDashboardHandler(dashboardSvc)
	companyH   := handler.NewCompanyHandler(companyRepo, importanceRepo, readinessSvc)
	knowledgeH := handler.NewKnowledgeHandler(nodeRepo)
	attemptH   := handler.NewAttemptHandler(attemptSvc, attemptRepo)
	revisionH  := handler.NewRevisionHandler(revisionSvc, revisionRepo)

	// ── Router ────────────────────────────────────────────────────────────────
	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("POST /api/auth/register", authH.Register)
	mux.HandleFunc("POST /api/auth/login",    authH.Login)

	// Knowledge graph (public read)
	mux.HandleFunc("GET /api/knowledge/domains",       knowledgeH.ListDomains)
	mux.HandleFunc("GET /api/knowledge/nodes",         knowledgeH.ListAll)
	mux.HandleFunc("GET /api/knowledge/nodes/",        knowledgeH.GetByID)
	mux.HandleFunc("GET /api/companies",               companyH.List)
	mux.HandleFunc("GET /api/companies/",              companyH.GetByID)

	// Protected routes (require JWT)
	mux.Handle("GET /api/auth/me",                 middleware.RequireAuth(http.HandlerFunc(authH.Me)))
	mux.Handle("GET /api/dashboard",              middleware.RequireAuth(http.HandlerFunc(dashH.Get)))
	mux.Handle("POST /api/companies",             middleware.RequireAuth(http.HandlerFunc(companyH.Create)))
	mux.Handle("GET /api/companies/{id}/readiness", middleware.RequireAuth(http.HandlerFunc(companyH.GetReadiness)))
	mux.Handle("GET /api/attempts",               middleware.RequireAuth(http.HandlerFunc(attemptH.List)))
	mux.Handle("POST /api/attempts",              middleware.RequireAuth(http.HandlerFunc(attemptH.Create)))
	mux.Handle("GET /api/revisions/today",        middleware.RequireAuth(http.HandlerFunc(revisionH.Today)))
	mux.Handle("POST /api/revisions/{id}/sessions", middleware.RequireAuth(http.HandlerFunc(revisionH.Complete)))

	// AI endpoints (public — no auth required for analysis)
	mux.HandleFunc("POST /api/ai/company-analysis", companyAnalysisHandler)
	mux.HandleFunc("POST /api/ai/topic-drill",      topicDrillHandler)
	mux.HandleFunc("POST /api/ai/detect-topic",     detectTopicHandler)
	mux.HandleFunc("POST /api/ai/log-problem",      logProblemHandler)
	mux.Handle("POST /api/ai/companies/", middleware.RequireAuth(http.HandlerFunc(generateCompanyRoadmapHandler)))

	// CORS middleware
	corsHandler := corsMiddleware(mux)

	port := defaultString(os.Getenv("PORT"), "8080")
	log.Printf("PUYRG API running on :%s ✓\n", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// ── CORS ──────────────────────────────────────────────────────────────────────

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func defaultString(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func extractJSONObject(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		start := strings.Index(text, "\n")
		if start >= 0 {
			text = text[start+1:]
		}
		if idx := strings.LastIndex(text, "```"); idx >= 0 {
			text = text[:idx]
		}
		text = strings.TrimSpace(text)
	}
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		return text[start : end+1]
	}
	return ""
}

// ── Company Roadmap (legacy AI endpoint) ──────────────────────────────────────

func generateCompanyRoadmapHandler(w http.ResponseWriter, r *http.Request) {
	companySlug := strings.TrimPrefix(r.URL.Path, "/api/ai/companies/")
	companySlug = strings.TrimSuffix(companySlug, "/")
	if companySlug == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "company slug required"})
		return
	}
	roadmap, err := generateRoadmapForCompany(companySlug)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, roadmap)
}

func generateRoadmapForCompany(companySlug string) (map[string]any, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	model := defaultString(os.Getenv("GEMINI_MODEL"), "gemini-2.5-flash")

	prompt := fmt.Sprintf(`Generate a company-specific interview preparation microtopic roadmap.
Company: %s
Return JSON only (no markdown): {"company": string, "roadmap": [{"topic": string, "subtopics": [{"name": string, "priority": "High"|"Medium"|"Low", "minProblems": number}]}]}`, companySlug)

	text, err := callGemini(key, model, "", prompt, 2000)
	if err != nil {
		return nil, err
	}
	text = extractJSONObject(text)

	var result map[string]any
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("parse roadmap JSON: %w", err)
	}
	return result, nil
}

// ── Topic Drill-Down ──────────────────────────────────────────────────────────

type TopicDrillRequest struct {
	TopicName   string `json:"topicName"`
	CompanyName string `json:"companyName"`
}

type DrillPattern struct {
	Name       string `json:"name"`
	Priority   string `json:"priority"`   // Critical / High / Medium / Low
	Easy       int    `json:"easy"`
	Medium     int    `json:"medium"`
	Hard       int    `json:"hard"`
	Total      int    `json:"total"`
	Notes      string `json:"notes"`
}

type DrillCategory struct {
	Name        string         `json:"name"`
	Priority    string         `json:"priority"`
	Description string         `json:"description"`
	Patterns    []DrillPattern `json:"patterns"`
}

type TopicDrillResponse struct {
	TopicName   string          `json:"topicName"`
	CompanyName string          `json:"companyName"`
	TotalProbs  int             `json:"totalProbs"`
	Categories  []DrillCategory `json:"categories"`
	AiUsed      bool            `json:"aiUsed"`
}

func topicDrillHandler(w http.ResponseWriter, r *http.Request) {
	var req TopicDrillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	req.TopicName = strings.TrimSpace(req.TopicName)
	req.CompanyName = strings.TrimSpace(req.CompanyName)
	if req.TopicName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "topicName is required"})
		return
	}

	result, err := generateTopicDrill(req.TopicName, req.CompanyName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func generateTopicDrill(topicName, companyName string) (*TopicDrillResponse, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	model := defaultString(os.Getenv("GEMINI_MODEL"), "gemini-2.5-flash")

	companyContext := ""
	if companyName != "" {
		companyContext = fmt.Sprintf("The candidate is preparing specifically for %s.", companyName)
	}

	systemPrompt := `You are PUYRG AI — an expert competitive programming and interview preparation system.
You have deep knowledge of every DSA topic, its subcategories, and all problem patterns.
Respond ONLY with valid JSON. No markdown, no explanation outside JSON.`

	userPrompt := fmt.Sprintf(`Generate a COMPLETE, EXHAUSTIVE drill-down for the DSA topic: "%s"
%s

Requirements:
- Break the topic into ALL its major categories/subcategories (minimum 8-15 categories)
- Each category must have ALL its atomic problem patterns (minimum 5-15 patterns each)
- Problem counts must be realistic for interview preparation
- Priority based on interview frequency: Critical (appears in >60%% of interviews), High (30-60%%), Medium (10-30%%), Low (<10%%)
- total = easy + medium + hard (always enforce)
- Be EXHAUSTIVE — a student should never need to look elsewhere for this topic

Example for "Graphs":
Categories: Graph Representation, DFS, BFS, Cycle Detection, Topological Sort,
Shortest Paths, MST, DSU, SCC & Bridges, Euler Path, Binary Lifting & LCA,
Tree DP, Heavy-Light Decomposition, Centroid Decomposition, Network Flow,
Graph Matching, Functional Graph, Graph DP, Advanced Graph Algorithms

Example for "DFS":
Patterns: Basic DFS Traversal, DFS on Grid/Matrix, Connected Components via DFS,
DFS with Visited Coloring, DFS Discovery/Finish Time, DFS Tree Edges Classification,
DFS on DAG, Cycle Detection Undirected (DFS), Cycle Detection Directed (DFS),
Bipartite Check via DFS, DFS + Backtracking, DFS + Memoization,
DFS Flood Fill, Path Finding DFS, DFS on State Space,
Iterative DFS, DFS on Implicit Graph, Multi-component DFS,
DFS in Functional Graphs, Euler Tour via DFS

Return this EXACT JSON:
{
  "topicName": string,
  "companyName": string,
  "totalProbs": number (sum of all pattern totals),
  "categories": [
    {
      "name": string,
      "priority": "Critical"|"High"|"Medium"|"Low",
      "description": string (1 sentence),
      "patterns": [
        {
          "name": string,
          "priority": "Critical"|"High"|"Medium"|"Low",
          "easy": number,
          "medium": number,
          "hard": number,
          "total": number,
          "notes": string (why important for %s, 1 sentence)
        }
      ]
    }
  ],
  "aiUsed": true
}`, topicName, companyContext, companyName)

	text, err := callGemini(key, model, systemPrompt, userPrompt, 32000)
	if err != nil {
		return nil, err
	}

	text = extractJSONObject(text)
	if text == "" {
		return nil, fmt.Errorf("Gemini returned empty response")
	}

	var result TopicDrillResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("parse drill JSON: %w — raw: %.300s", err, text)
	}

	// Recompute totalProbs from actual data
	total := 0
	for _, cat := range result.Categories {
		for _, p := range cat.Patterns {
			total += p.Total
		}
	}
	result.TotalProbs = total
	result.AiUsed = true
	return &result, nil
}

// ── Company Analysis ──────────────────────────────────────────────────────────

type CompanyAnalysisRequest struct {
	CompanyName string `json:"companyName"`
}

type ReadinessWeight struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

type AnalysisTopic struct {
	Name       string `json:"name"`
	Importance string `json:"importance"`
	MinProbs   int    `json:"minProbs"`
	Notes      string `json:"notes"`
}

type DSASubtopic struct {
	Name       string `json:"name"`
	Importance string `json:"importance"`
	Easy       int    `json:"easy"`
	Medium     int    `json:"medium"`
	Hard       int    `json:"hard"`
	Total      int    `json:"total"`
	Notes      string `json:"notes"`
}

type DSATopic struct {
	Name      string       `json:"name"`
	Priority  string       `json:"priority"`
	Subtopics []DSASubtopic `json:"subtopics"`
}

type CompanyAnalysisSection struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Points  []string `json:"points"`
}

type CompanyAnalysisResponse struct {
	CompanyName      string                   `json:"companyName"`
	Tier             string                   `json:"tier"`
	Overview         string                   `json:"overview"`
	ReadinessWeights []ReadinessWeight        `json:"readinessWeights"`
	RequiredTopics   []AnalysisTopic          `json:"requiredTopics"`
	DSATopics        []DSATopic               `json:"dsaTopics"`
	CareerTrack      string                   `json:"careerTrack"`
	Sections         []CompanyAnalysisSection `json:"sections"`
	HighestImpact    []string                 `json:"highestImpact"`
	InterviewFormat  string                   `json:"interviewFormat"`
	AiUsed           bool                     `json:"aiUsed"`
}

func companyAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	var req CompanyAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}
	req.CompanyName = strings.TrimSpace(req.CompanyName)
	if req.CompanyName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "companyName is required"})
		return
	}
	analysis, err := generateCompanyAnalysis(req.CompanyName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, analysis)
}

func fullDocsContext() string {
	return `
=== PUYRG COMPANY PROFILES ===
TIER 1: Rubrik — DSA 30, Core CS 25, Backend 25, Projects 10, Behavioral 10
TIER 2: Google — DSA 45, Core CS 20, Development 15, Projects 10, Behavioral 10
TIER 2: Meta, Microsoft, Amazon, Apple, Atlassian, Snowflake, Databricks, Stripe, Uber
TIER 3 (HFT): Jane Street, HRT, Citadel, Tower Research, Optiver, IMC — DSA/CP 55, Math 20
TIER 4 (CP): ICPC, Codeforces, AtCoder

=== TOPIC TAXONOMY ===
DSA: Arrays, Strings, Binary Search, Two Pointers, Sliding Window, Hashing, Prefix Sum,
     Sorting, Greedy, Recursion/Backtracking, Bit Manipulation, Linked List, Stack,
     Queue/Deque, Monotonic Stack, Heap, Binary Tree, BST, Trie, Tree Algorithms,
     Graph Traversal, Shortest Paths, Topological Sort, DSU, MST, SCC/Bridges,
     Network Flow, Graph Matching, Segment Tree, Fenwick Tree, Sparse Table, Mo's, DP
Mathematics: Number Theory, Combinatorics, Probability, Game Theory, Geometry
Core CS: OS, DBMS, Computer Networks, OOP, Concurrency, Distributed Systems
Backend: Go, PostgreSQL, Redis, REST, gRPC, Auth, Message Queues, Testing
System Design: LLD, HLD, Design Patterns, SOLID, Scalability/Caching, Database Design
CP: Advanced Strings, FFT/NTT, DP Optimizations, Advanced DS, Constructive, Meet in Middle

=== MASTERY LEVELS ===
0: Never Studied | 1: Concept Known | 2: Basic Problems | 3: Medium Comfortable
4: Hard Comfortable | 5: Interview Ready | 6: CP Ready | 7: Can Teach Others

=== REVISION SYSTEM ===
R1: 3 days | R2: 12 days | R3: 45 days → Mastered (accuracy > 80%, no hints)
`
}

func generateCompanyAnalysis(companyName string) (*CompanyAnalysisResponse, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set — add it to backend/.env")
	}
	model := defaultString(os.Getenv("GEMINI_MODEL"), "gemini-2.5-flash")

	systemPrompt := fmt.Sprintf(`You are PUYRG AI — an expert interview preparation system with deep knowledge of software engineering interviews.

CONTEXT:
%s

CRITICAL RULES:
1. Respond ONLY with valid JSON. Zero markdown, zero explanation outside JSON.
2. Be EXHAUSTIVE. This is a complete interview OS, not a quick summary.
3. Every number must be realistic and specific to %s's actual hiring bar.
4. total = easy + medium + hard (always enforce this).
5. readinessWeights must sum to exactly 100.`, fullDocsContext(), companyName)

	userPrompt := fmt.Sprintf(`Generate a COMPLETE, EXHAUSTIVE interview preparation profile for "%s".

REQUIREMENTS:
- dsaTopics: Include ALL relevant DSA topics (minimum 15-20 topics for Tier 1/2, 20-25 for HFT/CP)
- Each DSA topic must have ALL its important patterns as subtopics (minimum 8-15 patterns per topic)
- Problem counts must reflect REAL preparation needed: easy 3-15, medium 10-30, hard 3-15 per pattern
- Total problems across all topics should be 3000-6000 for Tier 1/2 companies
- requiredTopics: Include ALL topics (DSA + Core CS + Backend + System Design) — minimum 20 topics
- Be SPECIFIC to %s's known interview patterns, not generic

Return this EXACT JSON (no deviations):
{
  "companyName": string,
  "tier": string,
  "overview": string (3-4 sentences, very specific to company),
  "careerTrack": string,
  "interviewFormat": string (detailed, mention all rounds),
  "readinessWeights": [
    {"name": string, "value": number, "color": string (hex)}
  ],
  "highestImpact": [string] (top 10 highest impact preparation tasks, very specific),
  "requiredTopics": [
    {
      "name": string,
      "importance": "Critical"|"High"|"Medium",
      "minProbs": number,
      "notes": string (1 sentence, company-specific)
    }
  ],
  "dsaTopics": [
    {
      "name": string (ONE specific atomic topic — e.g. "Arrays" NOT "Arrays & Sequences", "1D DP" NOT "Dynamic Programming", "DFS" NOT "Graph Traversal"),
      "priority": "Critical"|"High"|"Medium",
      "subtopics": [
        {
          "name": string (specific pattern),
          "importance": "Critical"|"High"|"Medium",
          "easy": number,
          "medium": number,
          "hard": number,
          "total": number,
          "notes": string (why this pattern matters for %s)
        }
      ]
    }
  ],
  "sections": [
    {
      "title": string,
      "content": string,
      "points": [string]
    }
  ]
}

CRITICAL RULE FOR dsaTopics:
- Each entry in dsaTopics MUST be ONE atomic topic — NOT a group
- WRONG: {"name": "Dynamic Programming"} — this is a group, NOT allowed
- WRONG: {"name": "Graph Algorithms"} — this is a group, NOT allowed  
- WRONG: {"name": "Trees"} — this is a group, NOT allowed
- RIGHT: {"name": "1D DP"}, {"name": "Tree DP"}, {"name": "Bitmask DP"} — separate entries
- RIGHT: {"name": "DFS"}, {"name": "BFS"}, {"name": "Dijkstra"} — separate entries
- RIGHT: {"name": "Binary Tree"}, {"name": "BST"}, {"name": "Trie"} — separate entries
- MINIMUM 70 separate entries in dsaTopics array
- MAXIMUM subtopics per topic: 8 (keep concise — user can drill down for more)
- Keep subtopic notes SHORT (5 words max)

EACH of these MUST be its OWN separate dsaTopics entry:
Time Complexity, Recursion Basics,
Arrays, Matrix, Prefix Sum, Difference Array, Sliding Window, Two Pointers,
Binary Search, Ternary Search, Sorting, Greedy,
String Basics, String Hashing, KMP Algorithm, Z Algorithm, Rabin-Karp, Suffix Array, Suffix Automaton,
Linked List, Stack, Queue, Deque, Monotonic Stack, Monotonic Queue,
Heap / Priority Queue, Ordered Set / PBDS,
Bit Manipulation, Bitmasking, Recursion, Backtracking,
Binary Tree, BST, N-ary Tree, Trie, Binary Lifting, LCA,
Tree DP, Rerooting DP, Centroid Decomposition, Heavy-Light Decomposition, DSU on Tree, Euler Tour,
Graph Representation, DFS, BFS, Multi-Source BFS, 0-1 BFS,
Cycle Detection, Bipartite Check, Topological Sort,
Dijkstra, Bellman-Ford, Floyd-Warshall, SPFA,
Kruskal MST, Prim MST, DSU / Union Find,
SCC (Kosaraju/Tarjan), Bridges & Articulation Points, Euler Path,
Functional Graphs, Graph DP,
Fenwick Tree, Segment Tree, Segment Tree Lazy, Sparse Table, Sqrt Decomposition,
1D DP, 2D DP, Knapsack DP, LIS DP, LCS DP, Grid DP,
Tree DP (standalone), Bitmask DP, Digit DP, Interval DP, Profile DP,
Probability DP, Game Theory DP, DP on DAG,
CHT Optimization, D&C DP Optimization, Knuth Optimization,
Number Theory, Combinatorics, Probability & Statistics, Game Theory, Matrix Exponentiation,
Computational Geometry, Convex Hull,
Max Flow, Min Cost Flow, Bipartite Matching, General Matching,
FFT / NTT, Treap / Splay, Link-Cut Tree, Li Chao Tree, Wavelet Tree,
Meet in the Middle, Coordinate Compression, Sweep Line,
Mo's Algorithm, Small-to-Large Merging, Divide & Conquer Optimization, Convex Hull Trick

For Tier 1/2 (product companies) ALSO add as separate entries:
Operating Systems, DBMS, Computer Networks, OOP Concepts, Concurrency,
Go Language, PostgreSQL, Redis, REST API Design, LLD, HLD, Distributed Systems

Only skip topics clearly irrelevant to %s's tier/focus.
Total problems across all dsaTopics should be 10,000-20,000.
ADVANCED ALGORITHMS: FFT/NTT, Advanced Data Structures (Treap/Splay/LCT/Wavelet), Offline Algorithms, Randomized Algorithms
CP TECHNIQUES: Meet in the Middle, Coordinate Compression, Sweep Line, Mo's Algorithm, Small-to-Large Merging, Divide & Conquer Optimization, Convex Hull Trick
CORE CS (for Tier 1/2): Operating Systems, DBMS, Computer Networks, OOP, Concurrency, Distributed Systems
BACKEND (for Tier 1): Go Language, PostgreSQL, Redis, REST API Design, gRPC, Authentication & Security, System Design LLD, System Design HLD

Only include topics relevant to %s's tier and focus. Skip HFT-specific topics for product companies, skip backend for pure CP companies.
Total problems across all topics should be 10000-20000 for Tier 1/2 companies doing full preparation.`,
		companyName, companyName, companyName)

	text, err := callGemini(key, model, systemPrompt, userPrompt, 32000)
	if err != nil {
		return nil, err
	}

	text = extractJSONObject(text)
	if text == "" {
		return nil, fmt.Errorf("Gemini returned empty response")
	}

	var analysis CompanyAnalysisResponse
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		return nil, fmt.Errorf("parse analysis JSON: %w — raw: %.300s", err, text)
	}
	analysis.AiUsed = true
	return &analysis, nil
}

// callGemini sends a request to Gemini API and returns the text response.
// systemPrompt can be empty string if not needed.
func callGemini(apiKey, model, systemPrompt, userPrompt string, maxTokens int) (string, error) {
	return callGeminiWithImage(apiKey, model, systemPrompt, userPrompt, "", "", maxTokens)
}

func callGeminiWithImage(apiKey, model, systemPrompt, userPrompt, imageMimeType, imageData string, maxTokens int) (string, error) {
	imageData = strings.TrimSpace(imageData)
	imageMimeType = strings.TrimSpace(imageMimeType)
	if idx := strings.Index(imageData, ","); strings.HasPrefix(imageData, "data:") && idx >= 0 {
		imageData = imageData[idx+1:]
	}

	type inlineData struct {
		MimeType string `json:"mime_type"`
		Data     string `json:"data"`
	}
	type part struct {
		Text       string      `json:"text,omitempty"`
		InlineData *inlineData `json:"inline_data,omitempty"`
	}
	type content struct {
		Role  string `json:"role"`
		Parts []part `json:"parts"`
	}
	type genConfig struct {
		MaxOutputTokens int     `json:"maxOutputTokens"`
		Temperature     float64 `json:"temperature"`
		ResponseMimeType string  `json:"responseMimeType,omitempty"`
	}
	type systemInstruction struct {
		Parts []part `json:"parts"`
	}

	type requestBody struct {
		Contents           []content          `json:"contents"`
		GenerationConfig   genConfig          `json:"generationConfig"`
		SystemInstruction  *systemInstruction `json:"systemInstruction,omitempty"`
	}

	parts := []part{{Text: userPrompt}}
	if imageData != "" {
		if imageMimeType == "" {
			imageMimeType = "image/png"
		}
		parts = append(parts, part{InlineData: &inlineData{MimeType: imageMimeType, Data: imageData}})
	}

	reqBody := requestBody{
		Contents: []content{
			{Role: "user", Parts: parts},
		},
		GenerationConfig: genConfig{
			MaxOutputTokens: maxTokens,
			Temperature:     0.3,
			ResponseMimeType: "application/json",
		},
	}

	if systemPrompt != "" {
		reqBody.SystemInstruction = &systemInstruction{
			Parts: []part{{Text: systemPrompt}},
		}
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model, apiKey,
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 180 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	var decoded map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", fmt.Errorf("decode Gemini response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("Gemini API error %d: %v", resp.StatusCode, decoded["error"])
	}

	// Parse: candidates[0].content.parts[0].text
	candidates, ok := decoded["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("Gemini returned no candidates")
	}
	candidate, ok := candidates[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid candidate format")
	}
	contentObj, ok := candidate["content"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid content format")
	}
	responseParts, ok := contentObj["parts"].([]any)
	if !ok || len(responseParts) == 0 {
		return "", fmt.Errorf("no parts in response")
	}
	firstPart, ok := responseParts[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid part format")
	}
	text, ok := firstPart["text"].(string)
	if !ok {
		return "", fmt.Errorf("no text in response part")
	}
	return text, nil
}

// ensure models import is used
var _ = models.NodeTypeDomain

// ── Topic Detection ───────────────────────────────────────────────────────────

type DetectTopicRequest struct {
	Text        string `json:"text"`
	CompanyName string `json:"companyName"`
}

type DetectedSubtopic struct {
	Name       string `json:"name"`
	Importance string `json:"importance"`
	Easy       int    `json:"easy"`
	Medium     int    `json:"medium"`
	Hard       int    `json:"hard"`
	Total      int    `json:"total"`
	Notes      string `json:"notes"`
}

type DetectedTopic struct {
	TopicName    string             `json:"topicName"`
	IsNew        bool               `json:"isNew"`
	Priority     string             `json:"priority"`
	DetectedFrom string             `json:"detectedFrom"`
	Subtopics    []DetectedSubtopic `json:"subtopics"`
}

type DetectTopicResponse struct {
	Topics  []DetectedTopic `json:"topics"`
	Summary string          `json:"summary"`
}

func detectTopicHandler(w http.ResponseWriter, r *http.Request) {
	var req DetectTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text is required"})
		return
	}
	result, err := detectTopicsFromText(req.Text, req.CompanyName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func detectTopicsFromText(inputText, companyName string) (*DetectTopicResponse, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	model := defaultString(os.Getenv("GEMINI_MODEL"), "gemini-2.5-flash")

	companyCtx := ""
	companyNote := ""
	if companyName != "" {
		companyCtx = fmt.Sprintf("The candidate is preparing for %s.", companyName)
		companyNote = fmt.Sprintf(" for %s", companyName)
	}

	systemPrompt := `You are PUYRG AI — an expert at identifying DSA topics, patterns and algorithms from problem descriptions, code snippets, competitive programming questions, or any technical text.
Respond ONLY with valid JSON. No markdown fences, no explanation outside JSON.`

	userPrompt := fmt.Sprintf(`Analyze the following text and identify ALL DSA topics and patterns present.
%s

Text to analyze:
---
%s
---

Instructions:
1. Identify every DSA topic/pattern mentioned or implied in the text
2. Group patterns under their parent topic (e.g. "Dijkstra" goes under "Shortest Paths")
3. If the text mentions a new/advanced pattern not commonly in standard lists, set isNew=true
4. Set realistic practice targets (easy/medium/hard) for each detected pattern
5. detectedFrom should be a short quote from the input that triggered the detection

Return EXACT JSON:
{
  "summary": string (1-2 sentences describing what topics were found),
  "topics": [
    {
      "topicName": string (parent topic e.g. "Graphs", "Dynamic Programming", "Arrays"),
      "isNew": boolean,
      "priority": "Critical"|"High"|"Medium"|"Low",
      "detectedFrom": string,
      "subtopics": [
        {
          "name": string (specific pattern),
          "importance": "Critical"|"High"|"Medium",
          "easy": number,
          "medium": number,
          "hard": number,
          "total": number,
          "notes": string (why this matters%s)
        }
      ]
    }
  ]
}`, companyCtx, inputText, companyNote)

	rawText, err := callGemini(key, model, systemPrompt, userPrompt, 8000)
	if err != nil {
		return nil, err
	}
	rawText = extractJSONObject(rawText)
	if rawText == "" {
		return nil, fmt.Errorf("Gemini returned empty response")
	}
	var result DetectTopicResponse
	if err := json.Unmarshal([]byte(rawText), &result); err != nil {
		return nil, fmt.Errorf("parse detect JSON: %w", err)
	}
	return &result, nil
}

// ── Log Problem ───────────────────────────────────────────────────────────────

type LogProblemRequest struct {
	Text          string `json:"text"`          // pasted question text or OCR'd screenshot text
	ImageData     string `json:"imageData"`     // base64 image data or data URL
	ImageMimeType string `json:"imageMimeType"` // image/png, image/jpeg, etc.
}

type LogProblemResponse struct {
	ProblemTitle  string  `json:"problemTitle"`
	Platform      string  `json:"platform"`      // LeetCode / Codeforces / CSES / AtCoder / GFG / Other
	Difficulty    string  `json:"difficulty"`     // Easy / Medium / Hard
	CFRating      int     `json:"cfRating"`       // 0 if not CF
	Topic         string  `json:"topic"`          // parent topic e.g. "Graphs"
	Subtopic      string  `json:"subtopic"`       // e.g. "DFS"
	Pattern       string  `json:"pattern"`        // e.g. "DFS on Grid"
	QualityScore  int     `json:"qualityScore"`   // 1-10 points
	QualityReason string  `json:"qualityReason"`  // why this quality score
	Summary       string  `json:"summary"`        // 1-sentence summary of what was detected
	Confidence    float64 `json:"confidence"`     // 0-1 AI confidence
}

func logProblemHandler(w http.ResponseWriter, r *http.Request) {
	var req LogProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	req.Text = strings.TrimSpace(req.Text)
	req.ImageData = strings.TrimSpace(req.ImageData)
	req.ImageMimeType = strings.TrimSpace(req.ImageMimeType)
	if req.Text == "" && req.ImageData == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "text or image is required"})
		return
	}
	result, err := analyzeProblem(req.Text, req.ImageMimeType, req.ImageData)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func analyzeProblem(inputText, imageMimeType, imageData string) (*LogProblemResponse, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	model := defaultString(os.Getenv("GEMINI_MODEL"), "gemini-2.5-flash")

	systemPrompt := `You are PUYRG AI. Extract DSA problem metadata.
RULES: Respond ONLY with a single valid JSON object. No markdown. No explanation. No extra text before or after the JSON.`

	inputText = strings.TrimSpace(inputText)
	if inputText == "" {
		inputText = "The problem is in the attached image. Read the screenshot, identify the problem, then extract metadata."
	}

	userPrompt := fmt.Sprintf(`Extract metadata from this DSA problem and return ONLY a JSON object.

Problem text: %s

Quality scoring: LeetCode Easy=2, Medium=3, Hard=5. CF rating 800-1200=2, 1300-1600=4, 1700-1900=6, 2000+=8. ICPC=9. Unknown=3.

Return ONLY this JSON with no extra text:
{"problemTitle":"<title max 8 words>","platform":"LeetCode","difficulty":"Medium","cfRating":0,"topic":"Graphs","subtopic":"DFS","pattern":"DFS on Grid","qualityScore":3,"qualityReason":"LC Medium","summary":"Find connected components","confidence":0.95}`, inputText)

	text, err := callGeminiWithImage(key, model, systemPrompt, userPrompt, imageMimeType, imageData, 4096)
	if err != nil {
		return nil, err
	}
	text = extractJSONObject(text)
	if text == "" {
		return nil, fmt.Errorf("Gemini did not return valid JSON")
	}
	var result LogProblemResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("parse log JSON: %w - raw: %.200s", err, text)
	}
	normalizeLogProblem(&result)
	return &result, nil
}

func normalizeLogProblem(result *LogProblemResponse) {
	result.ProblemTitle = defaultString(strings.TrimSpace(result.ProblemTitle), "Detected Problem")
	result.Platform = defaultString(strings.TrimSpace(result.Platform), "Other")
	result.Difficulty = defaultString(strings.TrimSpace(result.Difficulty), "Medium")
	result.Topic = defaultString(strings.TrimSpace(result.Topic), "DSA")
	result.Subtopic = defaultString(strings.TrimSpace(result.Subtopic), result.Topic)
	result.Pattern = defaultString(strings.TrimSpace(result.Pattern), result.Subtopic)
	result.QualityReason = defaultString(strings.TrimSpace(result.QualityReason), "Estimated from problem difficulty")
	result.Summary = defaultString(strings.TrimSpace(result.Summary), "Problem metadata detected.")
	if result.QualityScore < 1 {
		result.QualityScore = 3
	}
	if result.QualityScore > 10 {
		result.QualityScore = 10
	}
	if result.Confidence <= 0 {
		result.Confidence = 0.6
	}
	if result.Confidence > 1 {
		result.Confidence = 1
	}
}