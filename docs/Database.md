# PUYRG Database Design

This document defines the first-pass relational model. It will evolve before implementation.

## Core Identity Tables

### users

Stores account and profile information.

Suggested fields:

- id
- name
- email
- password_hash
- current_level
- target_deadline
- created_at
- updated_at

### career_goals

Examples: Software Engineer, Backend Engineer, AI Engineer, Competitive Programmer, Quant Developer.

### roles

Examples: SWE, Backend Engineer, AI Engineer, ML Engineer, SRE, DevOps Engineer, Quant Developer.

### companies

Stores company profiles such as Rubrik, Google, Microsoft, Jane Street, and custom companies.

## Knowledge Graph Tables

### domains

Examples: DSA, Core CS, Backend, System Design, AI/ML, DevOps.

### topics

Examples: Graphs, Dynamic Programming, Operating Systems, PostgreSQL.

### subtopics

Examples: DFS, Tree DP, Transactions, Scheduling.

### patterns

Examples: Connected Components, Binary Search on Answer, Lazy Propagation.

### knowledge_nodes

Generalized graph node table for flexible hierarchy.

Suggested fields:

- id
- parent_id
- domain_id
- node_type
- name
- slug
- description
- difficulty
- estimated_minutes
- revision_interval_days
- is_active

### prerequisites

Maps one knowledge node to prerequisite nodes.

## Mapping Tables

### company_node_importance

Maps company to knowledge node.

Fields:

- company_id
- node_id
- importance_score
- interview_frequency
- minimum_easy
- minimum_medium
- minimum_hard
- minimum_total
- recommended_total

### role_node_importance

Maps role to knowledge node.

### career_goal_node_importance

Maps career goal to knowledge node.

### readiness_weights

Defines readiness dimensions by company, role, or career goal.

## Learning Ledger Tables

### problems

Canonical problem catalog.

Fields:

- id
- platform
- external_id
- title
- url
- difficulty
- rating
- default_quality_weight

### problem_node_tags

Maps problems to topic, subtopic, and pattern nodes.

### attempts

Stores user problem-solving attempts.

Fields:

- id
- user_id
- problem_id
- node_id
- result
- time_taken_minutes
- confidence_score
- revision_needed
- mistake_type
- notes
- quality_weight
- attempted_at

### revision_schedules

Stores the planned revision lifecycle for a solved question.

Fields:

- id
- user_id
- problem_id
- attempt_id
- current_revision_number
- required_revisions
- next_revision_at
- memory_estimate
- status
- created_at
- updated_at

### revision_sessions

Stores each revision attempt.

Fields:

- id
- user_id
- problem_id
- attempt_id
- revision_number
- mode
- result
- time_taken_minutes
- confidence_score
- accuracy_score
- needed_hint
- notes
- revised_at

Modes:

- normal
- blind_resolve
- oral_concept

### mastery_records

Stores whether a problem or knowledge node is mastered by a user.

Fields:

- id
- user_id
- entity_type
- entity_id
- mastered_at
- mastery_score
- last_accuracy_score
- source

### progress_snapshots

Stores daily aggregate progress for trend charts.

## Resource Tables

### resources

Fields:

- id
- node_id
- type
- title
- url
- source
- estimated_minutes
- quality_score

## Roadmap Tables

### roadmaps

Stores generated or saved roadmaps.

### roadmap_items

Stores tasks in a roadmap: problems, revisions, projects, reading, mocks.

## Project Tables

### projects

Tracks mini labs, mini projects, major projects, and production features.

### project_tasks

Tracks project milestones such as JWT, Redis, Kafka, Docker, tests, deployment.

## AI Tables

### ai_analyses

Stores AI-generated insights and reports.

### ai_recommendations

Stores recommended next actions.

## Achievement Tables

### achievements

Defines achievement badges and milestones.

### user_achievements

Tracks earned achievements.

## Indexing Notes

Likely indexes:

- attempts(user_id, attempted_at)
- attempts(user_id, node_id)
- revision_schedules(user_id, next_revision_at)
- revision_sessions(user_id, revised_at)
- knowledge_nodes(parent_id)
- knowledge_nodes(slug)
- company_node_importance(company_id, node_id)
- role_node_importance(role_id, node_id)
- resources(node_id)
