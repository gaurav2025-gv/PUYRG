package seed

import "puyrg/backend/internal/models"

// topicNodes returns all topic-level nodes, parented under domains.
// nodeMap: slug -> ID, populated from DB after domains are seeded.
func topicNodes(nodeMap map[string]uint) []models.KnowledgeNode {
	m := models.DifficultyMedium
	e := models.DifficultyEasy
	h := models.DifficultyHard

	dsaID := nodeMap["dsa"]
	mathID := nodeMap["mathematics"]
	csID := nodeMap["core-cs"]
	beID := nodeMap["backend"]
	sdID := nodeMap["system-design"]
	cpID := nodeMap["cp"]
	epID := nodeMap["engineering-practices"]
	dvID := nodeMap["devops"]
	isID := nodeMap["interview-skills"]
	aiID := nodeMap["ai-ml"]

	ptr := func(id uint) *uint { return &id }

	return []models.KnowledgeNode{
		// ── DSA ──────────────────────────────────────────────────────
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Arrays", Slug: "arrays", Difficulty: e, EstimatedMinutes: 300, RevisionIntervalDays: 7, SortOrder: 1},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Strings", Slug: "strings", Difficulty: e, EstimatedMinutes: 300, RevisionIntervalDays: 7, SortOrder: 2},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Hashing", Slug: "hashing", Difficulty: e, EstimatedMinutes: 180, RevisionIntervalDays: 7, SortOrder: 3},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Prefix Sum", Slug: "prefix-sum", Difficulty: e, EstimatedMinutes: 120, RevisionIntervalDays: 7, SortOrder: 4},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Two Pointers", Slug: "two-pointers", Difficulty: e, EstimatedMinutes: 150, RevisionIntervalDays: 7, SortOrder: 5},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Sliding Window", Slug: "sliding-window", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 7, SortOrder: 6},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Binary Search", Slug: "binary-search", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 7, SortOrder: 7},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Sorting", Slug: "sorting", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 10, SortOrder: 8},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Greedy", Slug: "greedy", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 7, SortOrder: 9},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Recursion & Backtracking", Slug: "recursion-backtracking", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 7, SortOrder: 10},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Bit Manipulation", Slug: "bit-manipulation", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 10, SortOrder: 11},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Matrix", Slug: "matrix", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 10, SortOrder: 12},

		// Linear DS
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Linked List", Slug: "linked-list", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 7, SortOrder: 13},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Stack", Slug: "stack", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 7, SortOrder: 14},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Queue & Deque", Slug: "queue-deque", Difficulty: m, EstimatedMinutes: 150, RevisionIntervalDays: 7, SortOrder: 15},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Monotonic Stack", Slug: "monotonic-stack", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 7, SortOrder: 16},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Heap / Priority Queue", Slug: "heap", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 7, SortOrder: 17},

		// Trees
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Binary Tree", Slug: "binary-tree", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 7, SortOrder: 18},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "BST", Slug: "bst", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 7, SortOrder: 19},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Trie", Slug: "trie", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 10, SortOrder: 20},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Tree Algorithms", Slug: "tree-algorithms", Difficulty: h, EstimatedMinutes: 480, RevisionIntervalDays: 10, SortOrder: 21},

		// Graphs
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Graph Traversal", Slug: "graph-traversal", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 7, SortOrder: 22},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Shortest Paths", Slug: "shortest-paths", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 7, SortOrder: 23},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Topological Sort", Slug: "topological-sort", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 7, SortOrder: 24},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "DSU / Union Find", Slug: "dsu", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 10, SortOrder: 25},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "MST", Slug: "mst", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 10, SortOrder: 26},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "SCC & Bridges", Slug: "scc-bridges", Difficulty: h, EstimatedMinutes: 300, RevisionIntervalDays: 14, SortOrder: 27},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Network Flow", Slug: "network-flow", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 14, SortOrder: 28},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Graph Matching", Slug: "graph-matching", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 14, SortOrder: 29},

		// Range Query DS
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Segment Tree", Slug: "segment-tree", Difficulty: h, EstimatedMinutes: 480, RevisionIntervalDays: 10, SortOrder: 30},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Fenwick Tree", Slug: "fenwick-tree", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 10, SortOrder: 31},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Sparse Table", Slug: "sparse-table", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 14, SortOrder: 32},
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Mo's Algorithm", Slug: "mos-algorithm", Difficulty: h, EstimatedMinutes: 300, RevisionIntervalDays: 14, SortOrder: 33},

		// DP
		{ParentID: ptr(dsaID), DomainID: ptr(dsaID), Type: models.NodeTypeTopic, Name: "Dynamic Programming", Slug: "dp", Difficulty: h, EstimatedMinutes: 600, RevisionIntervalDays: 7, SortOrder: 34},

		// ── Mathematics ──────────────────────────────────────────────
		{ParentID: ptr(mathID), DomainID: ptr(mathID), Type: models.NodeTypeTopic, Name: "Number Theory", Slug: "number-theory", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 10, SortOrder: 1},
		{ParentID: ptr(mathID), DomainID: ptr(mathID), Type: models.NodeTypeTopic, Name: "Combinatorics", Slug: "combinatorics", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 10, SortOrder: 2},
		{ParentID: ptr(mathID), DomainID: ptr(mathID), Type: models.NodeTypeTopic, Name: "Probability & Statistics", Slug: "probability", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 10, SortOrder: 3},
		{ParentID: ptr(mathID), DomainID: ptr(mathID), Type: models.NodeTypeTopic, Name: "Game Theory", Slug: "game-theory", Difficulty: h, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 4},
		{ParentID: ptr(mathID), DomainID: ptr(mathID), Type: models.NodeTypeTopic, Name: "Geometry", Slug: "geometry", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 14, SortOrder: 5},
		{ParentID: ptr(mathID), DomainID: ptr(mathID), Type: models.NodeTypeTopic, Name: "Linear Algebra", Slug: "linear-algebra", Difficulty: h, EstimatedMinutes: 300, RevisionIntervalDays: 14, SortOrder: 6},
		{ParentID: ptr(mathID), DomainID: ptr(mathID), Type: models.NodeTypeTopic, Name: "Discrete Mathematics", Slug: "discrete-math", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 7},

		// ── Core CS ───────────────────────────────────────────────────
		{ParentID: ptr(csID), DomainID: ptr(csID), Type: models.NodeTypeTopic, Name: "Operating Systems", Slug: "operating-systems", Difficulty: m, EstimatedMinutes: 480, RevisionIntervalDays: 7, SortOrder: 1},
		{ParentID: ptr(csID), DomainID: ptr(csID), Type: models.NodeTypeTopic, Name: "DBMS", Slug: "dbms", Difficulty: m, EstimatedMinutes: 480, RevisionIntervalDays: 7, SortOrder: 2},
		{ParentID: ptr(csID), DomainID: ptr(csID), Type: models.NodeTypeTopic, Name: "Computer Networks", Slug: "computer-networks", Difficulty: m, EstimatedMinutes: 360, RevisionIntervalDays: 7, SortOrder: 3},
		{ParentID: ptr(csID), DomainID: ptr(csID), Type: models.NodeTypeTopic, Name: "OOP", Slug: "oop", Difficulty: e, EstimatedMinutes: 240, RevisionIntervalDays: 7, SortOrder: 4},
		{ParentID: ptr(csID), DomainID: ptr(csID), Type: models.NodeTypeTopic, Name: "Concurrency", Slug: "concurrency", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 7, SortOrder: 5},
		{ParentID: ptr(csID), DomainID: ptr(csID), Type: models.NodeTypeTopic, Name: "Distributed Systems", Slug: "distributed-systems", Difficulty: h, EstimatedMinutes: 480, RevisionIntervalDays: 7, SortOrder: 6},
		{ParentID: ptr(csID), DomainID: ptr(csID), Type: models.NodeTypeTopic, Name: "Computer Architecture", Slug: "computer-architecture", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 7},

		// ── Backend Engineering ───────────────────────────────────────
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "Go Language", Slug: "go-language", Difficulty: m, EstimatedMinutes: 600, RevisionIntervalDays: 7, SortOrder: 1},
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "PostgreSQL", Slug: "postgresql", Difficulty: m, EstimatedMinutes: 480, RevisionIntervalDays: 7, SortOrder: 2},
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "Redis", Slug: "redis", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 10, SortOrder: 3},
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "REST API Design", Slug: "rest-api", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 7, SortOrder: 4},
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "gRPC", Slug: "grpc", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 14, SortOrder: 5},
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "Authentication & Security", Slug: "auth-security", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 10, SortOrder: 6},
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "Message Queues", Slug: "message-queues", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 7},
		{ParentID: ptr(beID), DomainID: ptr(beID), Type: models.NodeTypeTopic, Name: "Testing", Slug: "testing", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 8},

		// ── System Design ─────────────────────────────────────────────
		{ParentID: ptr(sdID), DomainID: ptr(sdID), Type: models.NodeTypeTopic, Name: "LLD", Slug: "lld", Difficulty: m, EstimatedMinutes: 480, RevisionIntervalDays: 7, SortOrder: 1},
		{ParentID: ptr(sdID), DomainID: ptr(sdID), Type: models.NodeTypeTopic, Name: "HLD", Slug: "hld", Difficulty: h, EstimatedMinutes: 480, RevisionIntervalDays: 7, SortOrder: 2},
		{ParentID: ptr(sdID), DomainID: ptr(sdID), Type: models.NodeTypeTopic, Name: "Design Patterns", Slug: "design-patterns", Difficulty: m, EstimatedMinutes: 360, RevisionIntervalDays: 10, SortOrder: 3},
		{ParentID: ptr(sdID), DomainID: ptr(sdID), Type: models.NodeTypeTopic, Name: "SOLID Principles", Slug: "solid", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 10, SortOrder: 4},
		{ParentID: ptr(sdID), DomainID: ptr(sdID), Type: models.NodeTypeTopic, Name: "Scalability & Caching", Slug: "scalability-caching", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 10, SortOrder: 5},
		{ParentID: ptr(sdID), DomainID: ptr(sdID), Type: models.NodeTypeTopic, Name: "Database Design", Slug: "database-design", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 10, SortOrder: 6},

		// ── CP ────────────────────────────────────────────────────────
		{ParentID: ptr(cpID), DomainID: ptr(cpID), Type: models.NodeTypeTopic, Name: "Advanced String Algorithms", Slug: "advanced-strings", Difficulty: h, EstimatedMinutes: 480, RevisionIntervalDays: 14, SortOrder: 1},
		{ParentID: ptr(cpID), DomainID: ptr(cpID), Type: models.NodeTypeTopic, Name: "FFT / NTT", Slug: "fft-ntt", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 21, SortOrder: 2},
		{ParentID: ptr(cpID), DomainID: ptr(cpID), Type: models.NodeTypeTopic, Name: "Advanced DP Optimizations", Slug: "dp-optimizations", Difficulty: h, EstimatedMinutes: 480, RevisionIntervalDays: 14, SortOrder: 3},
		{ParentID: ptr(cpID), DomainID: ptr(cpID), Type: models.NodeTypeTopic, Name: "Advanced Data Structures", Slug: "advanced-ds", Difficulty: h, EstimatedMinutes: 600, RevisionIntervalDays: 14, SortOrder: 4},
		{ParentID: ptr(cpID), DomainID: ptr(cpID), Type: models.NodeTypeTopic, Name: "Constructive Algorithms", Slug: "constructive", Difficulty: h, EstimatedMinutes: 300, RevisionIntervalDays: 14, SortOrder: 5},
		{ParentID: ptr(cpID), DomainID: ptr(cpID), Type: models.NodeTypeTopic, Name: "Meet in the Middle", Slug: "meet-in-middle", Difficulty: h, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 6},

		// ── Engineering Practices ─────────────────────────────────────
		{ParentID: ptr(epID), DomainID: ptr(epID), Type: models.NodeTypeTopic, Name: "Clean Code", Slug: "clean-code", Difficulty: e, EstimatedMinutes: 180, RevisionIntervalDays: 14, SortOrder: 1},
		{ParentID: ptr(epID), DomainID: ptr(epID), Type: models.NodeTypeTopic, Name: "Code Review", Slug: "code-review", Difficulty: e, EstimatedMinutes: 120, RevisionIntervalDays: 14, SortOrder: 2},
		{ParentID: ptr(epID), DomainID: ptr(epID), Type: models.NodeTypeTopic, Name: "Debugging & Profiling", Slug: "debugging-profiling", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 14, SortOrder: 3},
		{ParentID: ptr(epID), DomainID: ptr(epID), Type: models.NodeTypeTopic, Name: "Observability & Logging", Slug: "observability", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 14, SortOrder: 4},

		// ── DevOps ────────────────────────────────────────────────────
		{ParentID: ptr(dvID), DomainID: ptr(dvID), Type: models.NodeTypeTopic, Name: "Docker", Slug: "docker", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 1},
		{ParentID: ptr(dvID), DomainID: ptr(dvID), Type: models.NodeTypeTopic, Name: "Kubernetes", Slug: "kubernetes", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 21, SortOrder: 2},
		{ParentID: ptr(dvID), DomainID: ptr(dvID), Type: models.NodeTypeTopic, Name: "CI/CD", Slug: "ci-cd", Difficulty: m, EstimatedMinutes: 180, RevisionIntervalDays: 14, SortOrder: 3},
		{ParentID: ptr(dvID), DomainID: ptr(dvID), Type: models.NodeTypeTopic, Name: "Linux & Shell", Slug: "linux-shell", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 4},

		// ── Interview Skills ──────────────────────────────────────────
		{ParentID: ptr(isID), DomainID: ptr(isID), Type: models.NodeTypeTopic, Name: "Resume Building", Slug: "resume", Difficulty: e, EstimatedMinutes: 120, RevisionIntervalDays: 30, SortOrder: 1},
		{ParentID: ptr(isID), DomainID: ptr(isID), Type: models.NodeTypeTopic, Name: "Behavioral Interviews", Slug: "behavioral", Difficulty: m, EstimatedMinutes: 240, RevisionIntervalDays: 14, SortOrder: 2},
		{ParentID: ptr(isID), DomainID: ptr(isID), Type: models.NodeTypeTopic, Name: "Mock Interviews", Slug: "mock-interviews", Difficulty: m, EstimatedMinutes: 300, RevisionIntervalDays: 7, SortOrder: 3},
		{ParentID: ptr(isID), DomainID: ptr(isID), Type: models.NodeTypeTopic, Name: "Communication", Slug: "communication", Difficulty: e, EstimatedMinutes: 120, RevisionIntervalDays: 30, SortOrder: 4},

		// ── AI / ML ───────────────────────────────────────────────────
		{ParentID: ptr(aiID), DomainID: ptr(aiID), Type: models.NodeTypeTopic, Name: "Machine Learning Basics", Slug: "ml-basics", Difficulty: m, EstimatedMinutes: 480, RevisionIntervalDays: 14, SortOrder: 1},
		{ParentID: ptr(aiID), DomainID: ptr(aiID), Type: models.NodeTypeTopic, Name: "Deep Learning", Slug: "deep-learning", Difficulty: h, EstimatedMinutes: 600, RevisionIntervalDays: 14, SortOrder: 2},
		{ParentID: ptr(aiID), DomainID: ptr(aiID), Type: models.NodeTypeTopic, Name: "LLMs & Transformers", Slug: "llms", Difficulty: h, EstimatedMinutes: 480, RevisionIntervalDays: 14, SortOrder: 3},
		{ParentID: ptr(aiID), DomainID: ptr(aiID), Type: models.NodeTypeTopic, Name: "AI Engineering", Slug: "ai-engineering", Difficulty: h, EstimatedMinutes: 360, RevisionIntervalDays: 14, SortOrder: 4},
	}
}
