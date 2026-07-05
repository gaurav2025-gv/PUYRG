package seed

import "puyrg/backend/internal/models"

func domainNodes() []models.KnowledgeNode {
	d := models.DifficultyBeginner
	return []models.KnowledgeNode{
		{Type: models.NodeTypeDomain, Name: "DSA", Slug: "dsa", Description: "Data Structures and Algorithms — interview and CP core.", Difficulty: d, SortOrder: 1},
		{Type: models.NodeTypeDomain, Name: "Mathematics", Slug: "mathematics", Description: "Number theory, combinatorics, probability, geometry.", Difficulty: d, SortOrder: 2},
		{Type: models.NodeTypeDomain, Name: "Core CS", Slug: "core-cs", Description: "OS, DBMS, CN, OOP, Computer Architecture.", Difficulty: d, SortOrder: 3},
		{Type: models.NodeTypeDomain, Name: "Backend Engineering", Slug: "backend", Description: "Go, PostgreSQL, Redis, Docker, REST, gRPC.", Difficulty: d, SortOrder: 4},
		{Type: models.NodeTypeDomain, Name: "System Design", Slug: "system-design", Description: "LLD, HLD, scalability, distributed systems.", Difficulty: d, SortOrder: 5},
		{Type: models.NodeTypeDomain, Name: "Competitive Programming", Slug: "cp", Description: "Advanced algorithms for ICPC, Codeforces, AtCoder.", Difficulty: d, SortOrder: 6},
		{Type: models.NodeTypeDomain, Name: "Engineering Practices", Slug: "engineering-practices", Description: "Clean code, testing, profiling, observability.", Difficulty: d, SortOrder: 7},
		{Type: models.NodeTypeDomain, Name: "DevOps", Slug: "devops", Description: "Docker, Kubernetes, CI/CD, monitoring.", Difficulty: d, SortOrder: 8},
		{Type: models.NodeTypeDomain, Name: "Projects", Slug: "projects", Description: "Mini labs, major projects, production features.", Difficulty: d, SortOrder: 9},
		{Type: models.NodeTypeDomain, Name: "Interview Skills", Slug: "interview-skills", Description: "Resume, behavioral, mock interviews, communication.", Difficulty: d, SortOrder: 10},
		{Type: models.NodeTypeDomain, Name: "AI / ML", Slug: "ai-ml", Description: "Machine learning, deep learning, LLMs, AI engineering.", Difficulty: d, SortOrder: 11},
	}
}
