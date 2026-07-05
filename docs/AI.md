# PUYRG AI System

## AI Purpose

AI should not be decorative. It should explain, prioritize, and personalize preparation.

AI should answer:

- What should I study today?
- Which weak topics block my target company?
- Which solved count is misleading?
- What revision is due?
- Which solved questions are likely to be forgotten?
- What project feature should I build next?
- How ready am I for a given company and role?

## AI Features

### Personalized Roadmap

Generates roadmap from:

- Career goals
- Companies
- Roles
- Current progress
- Weak patterns
- Deadline
- Available time

### Weak Topic Detection

Detects gaps such as:

- Many DP problems solved, but no Tree DP
- Graphs strong overall, but shortest path weak
- High solved count, low hard problem quality
- Repeated edge case mistakes

### Interview Simulator

Creates company-specific simulations:

```text
Google OA
- 2 Arrays
- 1 Graph
- 1 DP
Time: 90 minutes
```

### Resume Analyzer

Scores resume against target role and company.

### Code Review

Reviews solutions and project code for:

- Correctness
- Complexity
- Edge cases
- Clean code
- Production engineering practices

### Weekly Planner

Turns roadmap into weekly actions.

### Progress Prediction

Predicts readiness if the user follows the plan.

### Revision Intelligence

AI should help schedule and explain revisions using:

- Spaced repetition
- Overdue status
- Weak topic priority
- Problem difficulty
- Last confidence score
- Last accuracy score
- Memory decay estimate

Every solved question should complete at least 3 revisions before it is counted as mastered.

### Blind Revision

AI can start a blind re-solving session where the user sees only the question name and must solve it again. The system compares first solve time and revision solve time.

### Oral Revision

AI can ask concept checks without requiring code, for example:

- Explain Binary Search on Answer.
- How does DSU path compression work?
- When is BFS better than DFS?

Good oral answers can increase confidence and mastery for the related concept.

## AI/ML Career Track

AI/ML is a major career path, not a small optional feature.

Levels:

1. Mathematics
2. Python Ecosystem
3. Machine Learning
4. Deep Learning
5. NLP and LLMs
6. Computer Vision
7. MLOps
8. AI Engineering
9. Frameworks
10. Research

Example topics:

- Linear Algebra
- Probability
- Optimization
- NumPy
- Pandas
- Regression
- Classification
- Transformers
- RAG
- Vector Databases
- AI Agents
- MCP
- PyTorch
- Hugging Face
- RLHF
- Quantization
- Distillation

## Guardrails

AI should:

- Explain recommendations using data from progress and mappings.
- Avoid pretending offers are guaranteed.
- Prefer high-impact tasks over generic motivation.
- Surface uncertainty where data is incomplete.
