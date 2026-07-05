# PUYRG API Design

This is an initial REST API sketch. It may evolve before backend implementation.

## Auth

```text
POST /api/auth/register
POST /api/auth/login
POST /api/auth/logout
GET  /api/auth/me
```

## Goals

```text
GET  /api/goals
POST /api/users/me/goals
GET  /api/users/me/goals
```

## Companies and Roles

```text
GET /api/companies
GET /api/companies/:id
GET /api/roles
GET /api/roles/:id
```

## Knowledge Graph

```text
GET  /api/knowledge/domains
GET  /api/knowledge/nodes
GET  /api/knowledge/nodes/:id
GET  /api/knowledge/nodes/:id/children
GET  /api/knowledge/nodes/:id/prerequisites
GET  /api/knowledge/nodes/:id/resources
```

## Learning Ledger

```text
GET  /api/ledger/attempts
POST /api/ledger/attempts
GET  /api/ledger/stats
GET  /api/ledger/heatmap
GET  /api/ledger/mistakes
```

Example attempt request:

```json
{
  "platform": "LeetCode",
  "problemTitle": "Number of Connected Components in an Undirected Graph",
  "problemUrl": "https://leetcode.com/problems/number-of-connected-components-in-an-undirected-graph/",
  "difficulty": "Medium",
  "nodeId": "connected-components",
  "timeTakenMinutes": 28,
  "result": "Accepted",
  "confidenceScore": 8,
  "revisionNeeded": true,
  "mistakeType": "Logic",
  "notes": "Forgot to initialize visited for all nodes."
}
```

## Revision

```text
GET  /api/revisions/today
GET  /api/revisions/overdue
GET  /api/revisions/stats
POST /api/revisions/:scheduleId/sessions
POST /api/revisions/:scheduleId/blind-start
POST /api/revisions/:scheduleId/oral-start
PATCH /api/revisions/schedules/:scheduleId
```

Example revision session request:

```json
{
  "mode": "blind_resolve",
  "result": "Solved",
  "timeTakenMinutes": 26,
  "confidenceScore": 8,
  "accuracyScore": 85,
  "neededHint": false,
  "notes": "Solved faster than first attempt after remembering BFS layering."
}
```

## Progress

```text
GET /api/progress/summary
GET /api/progress/topics
GET /api/progress/nodes/:id
GET /api/progress/timeline
```

## Readiness

```text
GET /api/readiness
GET /api/readiness/companies/:companyId
GET /api/readiness/roles/:roleId
```

## Roadmaps

```text
GET  /api/roadmaps/current
POST /api/roadmaps/generate
GET  /api/roadmaps/:id
POST /api/roadmaps/:id/items/:itemId/complete
```

## AI

```text
POST /api/ai/mentor
POST /api/ai/recommendations
POST /api/ai/interview-simulation
POST /api/ai/resume-analyze
POST /api/ai/code-review
GET  /api/ai/insights
```

## Admin

```text
POST   /api/admin/knowledge/nodes
PATCH  /api/admin/knowledge/nodes/:id
DELETE /api/admin/knowledge/nodes/:id
POST   /api/admin/company-mappings
PATCH  /api/admin/company-mappings/:id
POST   /api/admin/resources
PATCH  /api/admin/resources/:id
```
