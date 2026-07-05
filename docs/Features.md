# PUYRG Features

## Core Modules

### 1. Goal and Career Modes

Users can select one or more career goals:

- Software Engineer
- Backend Engineer
- AI Engineer
- ML Engineer
- Competitive Programmer
- Research Engineer
- Quant Developer
- DevOps Engineer
- SRE
- Security Engineer

Selected goals change topic priorities, minimum practice targets, revision strategy, and readiness scoring.

### 2. Company Profiles

Users can prepare for companies such as:

- Rubrik
- Google
- Meta
- Microsoft
- Amazon
- Apple
- Atlassian
- Adobe
- Salesforce
- Snowflake
- Databricks
- Uber
- Airbnb
- LinkedIn
- Stripe
- Palantir
- Jane Street
- HRT
- Citadel Securities
- Optiver
- IMC
- Tower Research
- Jump Trading
- DRW
- Graviton
- Quadeye

Company readiness is calculated from weighted skill areas and topic-level requirements.

### 3. Master Knowledge Graph

The knowledge graph is the foundation:

```text
Domain -> Topic -> Subtopic -> Pattern -> Problem/Resource/Progress
```

It supports DSA, Core CS, Backend, System Design, AI/ML, CP, Projects, Soft Skills, and more.

### 4. Learning Ledger

Learning Ledger stores every meaningful learning action, especially problem attempts.

Each question entry includes:

- Platform
- Problem name, ID, and URL
- Difficulty
- Topic, subtopic, and pattern
- Time taken
- Result
- Confidence
- Revision needed
- Mistake type
- Notes
- Quality weight

This enables insight such as:

- "You solved 120 DP problems, but Tree DP is still weak."
- "42% of mistakes are observation-based."
- "Rubrik readiness improved from 74.0% to 74.6%."

### 5. Revision Engine

PUYRG must track retention, not just progress. Every solved question should require at least 3 revisions before it can be considered mastered.

```text
First Solve
  -> Revision 1 after 2-3 days
  -> Revision 2 after 10-15 days
  -> Revision 3 after 45-60 days
  -> Mastered
```

Question cards should show:

- Solved status
- Revision 1 status
- Revision 2 status
- Revision 3 status
- Mastery percentage
- Last revision date
- Next revision date
- Memory estimate

Dashboard metrics:

- Solved questions
- Questions needing revision
- Mastered questions
- Overdue revisions
- Revision accuracy
- Forgotten questions

Smart revision priority:

1. Overdue questions
2. Weak topics
3. Hard problems
4. Recently forgotten questions

Revision modes:

- Normal revision
- Blind re-solving mode
- AI oral/concept revision

Blind re-solving mode shows only the question name and starts a timer. The system compares first solve time and revision solve time to measure improvement.

AI oral revision can ask questions like:

- Explain Binary Search on Answer.
- How does DSU work?
- When would you choose BFS over DFS?

Mastery conditions:

- Solved
- Revision 1 complete
- Revision 2 complete
- Revision 3 complete
- Last revision accuracy greater than 80%

Company readiness should weight mastered questions more than merely solved questions.

### 6. Readiness Score

Base readiness dimensions:

```text
DSA                 35
Core CS             20
Development         20
Projects            15
Interview Skills    10
Total              100
```

Company and role profiles override these weights.

### 7. Roadmap Engine

The roadmap engine combines:

- Selected goals
- Selected companies
- Selected roles
- Current progress
- Weak topics
- Company requirements
- Revision dates
- Overdue revisions
- Mastered question count
- Available time

It produces daily, weekly, and long-term plans.

### 8. Analytics

Analytics include:

- Daily progress
- Weekly consistency
- Topic heatmap
- Weakest topics
- Strongest topics
- Quality-adjusted solved score
- Average problem time
- Acceptance rate
- Mistake-type distribution
- Company readiness trend
- Revision accuracy
- Mastered vs solved ratio
- Memory decay trend

### 9. Resources Layer

Each node can include:

- Best article
- Best video
- Practice list
- Revision notes
- Project usage
- Estimated learning time

### 10. AI Mentor

AI features include:

- Personalized roadmap
- "What should I study today?"
- Weak topic detection
- Interview simulator
- Resume analyzer
- Code review
- Motivation and weekly planning
- Readiness prediction
- Spaced revision planning
- Oral revision

### 11. Admin Panel

Admins can manage:

- Topics
- Subtopics
- Patterns
- Company mappings
- Role mappings
- Minimum practice targets
- Resources
- Difficulty splits
- Interview frequencies
- Revision intervals
- Readiness weights

## Non-Goals

- PUYRG should not become a generic todo app.
- PUYRG should not hardcode company targets into frontend logic.
- PUYRG should not reward raw solved count without topic coverage and question quality.
