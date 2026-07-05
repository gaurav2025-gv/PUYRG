package seed

import (
	"puyrg/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// seedCompanyNodeImportance seeds importance mappings for key topics per company.
func seedCompanyNodeImportance(db *gorm.DB) error {
	// Build maps
	companyMap := map[string]uint{}
	var companies []models.Company
	db.Find(&companies)
	for _, c := range companies {
		companyMap[c.Slug] = c.ID
	}

	nodeMap := map[string]uint{}
	var nodes []models.KnowledgeNode
	db.Find(&nodes)
	for _, n := range nodes {
		nodeMap[n.Slug] = n.ID
	}

	type mapping struct {
		company string
		node    string
		score   int   // 1-5
		freq    int   // 1-10
		label   string
		easy    int
		medium  int
		hard    int
	}

	mappings := []mapping{
		// ── Rubrik ────────────────────────────────────────────────────────────
		{"rubrik","arrays",5,9,"Critical",5,20,5},
		{"rubrik","strings",4,8,"Critical",5,15,5},
		{"rubrik","binary-search",5,9,"Critical",5,15,5},
		{"rubrik","graph-traversal",5,9,"Critical",5,20,5},
		{"rubrik","shortest-paths",4,8,"High",3,12,4},
		{"rubrik","topological-sort",4,8,"High",3,10,3},
		{"rubrik","dsu",4,7,"High",2,10,3},
		{"rubrik","mst",3,5,"Medium",0,5,2},
		{"rubrik","dp",5,9,"Critical",5,20,8},
		{"rubrik","segment-tree",4,7,"High",2,10,4},
		{"rubrik","fenwick-tree",3,6,"Medium",2,8,2},
		{"rubrik","binary-tree",5,9,"Critical",5,18,5},
		{"rubrik","bst",4,7,"High",3,10,3},
		{"rubrik","trie",4,7,"High",2,8,3},
		{"rubrik","heap",4,8,"High",3,10,3},
		{"rubrik","hashing",5,9,"Critical",5,15,3},
		{"rubrik","sliding-window",5,9,"Critical",5,15,3},
		{"rubrik","two-pointers",5,9,"Critical",5,15,3},
		{"rubrik","greedy",4,8,"High",3,12,4},
		{"rubrik","recursion-backtracking",4,8,"High",3,10,4},
		{"rubrik","linked-list",4,8,"High",3,12,3},
		{"rubrik","stack",4,7,"High",3,10,3},
		{"rubrik","monotonic-stack",4,7,"High",2,10,3},
		{"rubrik","operating-systems",5,9,"Critical",0,0,0},
		{"rubrik","dbms",5,9,"Critical",0,0,0},
		{"rubrik","computer-networks",4,8,"High",0,0,0},
		{"rubrik","concurrency",5,9,"Critical",0,0,0},
		{"rubrik","go-language",5,9,"Critical",0,0,0},
		{"rubrik","postgresql",5,9,"Critical",0,0,0},
		{"rubrik","redis",4,7,"High",0,0,0},
		{"rubrik","rest-api",5,9,"Critical",0,0,0},
		{"rubrik","lld",5,9,"Critical",0,0,0},
		{"rubrik","hld",4,7,"High",0,0,0},
		{"rubrik","distributed-systems",5,9,"Critical",0,0,0},

		// ── Google ────────────────────────────────────────────────────────────
		{"google","arrays",5,10,"Critical",5,25,8},
		{"google","strings",5,9,"Critical",5,20,8},
		{"google","binary-search",5,9,"Critical",5,20,5},
		{"google","graph-traversal",5,10,"Critical",5,25,10},
		{"google","shortest-paths",5,9,"Critical",3,15,6},
		{"google","topological-sort",5,9,"Critical",3,12,5},
		{"google","dsu",4,8,"High",2,12,4},
		{"google","dp",5,10,"Critical",5,25,10},
		{"google","segment-tree",4,7,"High",2,10,5},
		{"google","binary-tree",5,10,"Critical",5,20,8},
		{"google","trie",5,9,"Critical",3,12,5},
		{"google","heap",5,9,"Critical",3,15,5},
		{"google","hashing",5,10,"Critical",5,20,5},
		{"google","sliding-window",5,9,"Critical",5,15,5},
		{"google","two-pointers",5,9,"Critical",5,15,5},
		{"google","greedy",5,9,"Critical",3,15,5},
		{"google","recursion-backtracking",5,9,"Critical",3,12,5},
		{"google","bit-manipulation",4,7,"High",3,10,4},
		{"google","operating-systems",4,7,"High",0,0,0},
		{"google","dbms",4,7,"High",0,0,0},
		{"google","lld",4,7,"High",0,0,0},
		{"google","hld",5,9,"Critical",0,0,0},

		// ── Meta ──────────────────────────────────────────────────────────────
		{"meta","arrays",5,10,"Critical",5,25,8},
		{"meta","strings",5,9,"Critical",5,20,5},
		{"meta","binary-search",5,9,"Critical",5,20,5},
		{"meta","graph-traversal",5,9,"Critical",5,20,8},
		{"meta","dp",5,9,"Critical",5,20,8},
		{"meta","binary-tree",5,10,"Critical",5,20,8},
		{"meta","heap",4,8,"High",3,12,4},
		{"meta","hashing",5,9,"Critical",5,15,3},
		{"meta","sliding-window",5,9,"Critical",5,15,3},
		{"meta","two-pointers",5,9,"Critical",5,15,3},
		{"meta","greedy",4,8,"High",3,12,4},
		{"meta","recursion-backtracking",5,9,"Critical",3,12,5},

		// ── Microsoft ─────────────────────────────────────────────────────────
		{"microsoft","arrays",5,9,"Critical",5,20,5},
		{"microsoft","strings",5,9,"Critical",5,15,5},
		{"microsoft","binary-tree",5,9,"Critical",5,18,5},
		{"microsoft","graph-traversal",4,8,"High",3,15,5},
		{"microsoft","dp",4,8,"High",3,15,5},
		{"microsoft","hashing",4,8,"High",3,12,3},
		{"microsoft","linked-list",4,8,"High",3,12,3},
		{"microsoft","stack",4,7,"High",3,10,3},
		{"microsoft","operating-systems",4,7,"High",0,0,0},
		{"microsoft","dbms",4,7,"High",0,0,0},
		{"microsoft","lld",4,8,"High",0,0,0},
		{"microsoft","hld",3,6,"Medium",0,0,0},

		// ── Jane Street ───────────────────────────────────────────────────────
		{"jane-street","dp",5,10,"Critical",3,20,15},
		{"jane-street","dp-optimizations",5,9,"Critical",0,5,10},
		{"jane-street","graph-traversal",5,9,"Critical",3,15,10},
		{"jane-street","shortest-paths",5,9,"Critical",2,10,8},
		{"jane-street","number-theory",5,9,"Critical",0,5,8},
		{"jane-street","combinatorics",5,9,"Critical",0,5,8},
		{"jane-street","probability",5,9,"Critical",0,5,8},
		{"jane-street","segment-tree",5,9,"Critical",2,10,8},
		{"jane-street","network-flow",5,8,"Critical",0,5,8},
		{"jane-street","advanced-strings",4,7,"High",0,3,5},
		{"jane-street","fft-ntt",4,7,"High",0,2,5},
		{"jane-street","advanced-ds",4,7,"High",0,3,5},

		// ── Tower Research ────────────────────────────────────────────────────
		{"tower-research","dp",5,10,"Critical",3,20,12},
		{"tower-research","dp-optimizations",5,9,"Critical",0,5,8},
		{"tower-research","graph-traversal",5,9,"Critical",3,15,8},
		{"tower-research","segment-tree",5,9,"Critical",2,10,8},
		{"tower-research","number-theory",5,9,"Critical",0,5,8},
		{"tower-research","combinatorics",5,9,"Critical",0,5,8},
		{"tower-research","probability",5,9,"Critical",0,5,8},
		{"tower-research","advanced-ds",4,8,"High",0,3,5},

		// ── ICPC ──────────────────────────────────────────────────────────────
		{"icpc","dp",5,10,"Critical",3,20,15},
		{"icpc","dp-optimizations",5,10,"Critical",0,5,12},
		{"icpc","advanced-strings",5,9,"Critical",0,5,10},
		{"icpc","fft-ntt",5,9,"Critical",0,3,8},
		{"icpc","advanced-ds",5,9,"Critical",0,5,10},
		{"icpc","network-flow",5,9,"Critical",0,5,8},
		{"icpc","graph-matching",5,8,"Critical",0,3,8},
		{"icpc","number-theory",5,9,"Critical",0,5,8},
		{"icpc","geometry",4,8,"High",0,5,8},
		{"icpc","meet-in-middle",4,7,"High",0,3,5},
		{"icpc","constructive",4,8,"High",0,10,5},
		{"icpc","scc-bridges",4,7,"High",0,5,5},
		{"icpc","tree-algorithms",5,9,"Critical",0,8,8},
	}

	var records []models.CompanyNodeImportance
	for _, mm := range mappings {
		cID := companyMap[mm.company]
		nID := nodeMap[mm.node]
		if cID == 0 || nID == 0 {
			continue
		}
		total := mm.easy + mm.medium + mm.hard
		records = append(records, models.CompanyNodeImportance{
			CompanyID:          cID,
			NodeID:             nID,
			ImportanceScore:    mm.score,
			InterviewFrequency: mm.freq,
			ImportanceLabel:    mm.label,
			MinimumEasy:        mm.easy,
			MinimumMedium:      mm.medium,
			MinimumHard:        mm.hard,
			MinimumTotal:       total,
			RecommendedTotal:   total + total/3,
		})
	}

	return db.Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(records, 50).Error
}
