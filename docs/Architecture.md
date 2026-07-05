# PUYRG Architecture

## Architecture Philosophy

PUYRG should be data-driven. The frontend and backend should not hardcode roadmap logic such as "Graphs = 50 questions". Instead, requirements live in the database and are interpreted by services.

Correct model:

```text
Career Goal
  -> Role
  -> Company
  -> Knowledge Node
  -> Target Rule
  -> Progress
  -> Recommendation
```

## Planned Tech Stack

### Frontend

- React
- TypeScript
- Tailwind CSS
- Recharts

### Backend

- Go
- Gin
- GORM
- PostgreSQL
- Redis
- JWT

### Deployment

- Docker
- Docker Compose
- VPS, Render, or similar platform
- Later: Kubernetes basics

## Major Services

### Knowledge Graph Service

Manages domains, topics, subtopics, patterns, prerequisites, resources, and metadata.

### Progress Service

Stores Learning Ledger entries and aggregates progress by topic, company, role, and career mode.

### Revision Service

Tracks scheduled revisions, blind re-solving sessions, oral revision, memory decay, and mastery status. It ensures solved questions do not count as mastered until the required revision cycle is complete.

### Roadmap Service

Generates prioritized learning plans from goals, mappings, readiness weights, and progress.

### Readiness Service

Computes readiness scores for companies and roles.

### Recommendation Service

Finds highest-impact tasks by combining weakness, importance, required practice, revision due dates, and target deadlines.

### AI Service

Provides AI mentor, resume analysis, roadmap explanation, weak topic detection, and interview simulation.

### Admin Service

Allows editing topic taxonomy, mappings, resources, and target rules without code changes.

## High-Level Data Flow

```text
User logs problem attempt
  -> Learning Ledger stores attempt
  -> Revision schedule is created
  -> Progress aggregates update
  -> Topic mastery recalculates
  -> Readiness scores update
  -> Roadmap priorities refresh
  -> AI mentor explains next action
```

## No-Hardcode Rule

The following must be configurable:

- Company topic importance
- Role topic importance
- Minimum problems
- Difficulty split
- Readiness weights
- Revision intervals
- Mastery thresholds
- Memory decay rules
- Resource links
- Topic hierarchy
