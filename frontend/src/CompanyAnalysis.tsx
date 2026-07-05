import {
  AlertCircle,
  BookOpen,
  Brain,
  Building2,
  CheckCircle2,
  ChevronDown,
  ChevronRight,
  Flame,
  Layers,
  Loader2,
  Search,
  Sparkles,
  Target,
  TrendingUp,
  Trophy,
  Users,
} from 'lucide-react'
import type { FormEvent, ReactNode } from 'react'
import { useState } from 'react'
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'
import './CompanyAnalysis.css'

type ReadinessWeight = {
  name: string
  value: number
  color: string
}

type AnalysisTopic = {
  name: string
  importance: 'Critical' | 'High' | 'Medium'
  minProbs: number
  notes: string
}

type AnalysisSection = {
  title: string
  content: string
  points: string[]
}

type DSASubtopic = {
  name: string
  importance: 'Critical' | 'High' | 'Medium'
  easy: number
  medium: number
  hard: number
  total: number
  notes: string
}

type DSATopic = {
  name: string
  priority: 'Critical' | 'High' | 'Medium'
  subtopics: DSASubtopic[]
}

type CompanyAnalysis = {
  companyName: string
  tier: string
  overview: string
  readinessWeights: ReadinessWeight[]
  requiredTopics: AnalysisTopic[]
  dsaTopics: DSATopic[]
  careerTrack: string
  sections: AnalysisSection[]
  highestImpact: string[]
  interviewFormat: string
  aiUsed: boolean
}

const sectionIcons: Record<string, ReactNode> = {
  'Interview Process': <Users size={18} />,
  'DSA Focus Areas': <Brain size={18} />,
  'System Design Expectations': <Layers size={18} />,
  'Behavioral & Culture': <Users size={18} />,
  'Preparation Timeline': <TrendingUp size={18} />,
  'Common Mistakes': <AlertCircle size={18} />,
}

const importanceConfig = {
  Critical: { cls: 'importance-critical', label: 'Critical' },
  High: { cls: 'importance-high', label: 'High' },
  Medium: { cls: 'importance-medium', label: 'Medium' },
}

const SUGGESTED_COMPANIES = [
  'Rubrik', 'Google', 'Meta', 'Microsoft', 'Amazon',
  'Jane Street', 'Citadel Securities', 'Stripe', 'Uber', 'OpenAI',
]

function DSATopicsSection({ topics }: { topics: DSATopic[] }) {
  const [openTopics, setOpenTopics] = useState<Set<string>>(
    () => new Set(topics.filter((t) => t.priority === 'Critical').map((t) => t.name))
  )

  function toggle(name: string) {
    setOpenTopics((prev) => {
      const next = new Set(prev)
      if (next.has(name)) next.delete(name)
      else next.add(name)
      return next
    })
  }

  const totalQuestions = topics.reduce(
    (sum, t) => sum + t.subtopics.reduce((s, st) => s + st.total, 0),
    0,
  )

  return (
    <div className="ca-dsa-section">
      <div className="ca-section-title">
        <Brain size={18} />
        DSA Topics Breakdown
        <span className="ca-topic-count">{topics.length} topics · {totalQuestions} total problems</span>
      </div>

      <div className="ca-dsa-list">
        {topics.map((topic) => {
          const isOpen = openTopics.has(topic.name)
          const topicTotal = topic.subtopics.reduce((s, st) => s + st.total, 0)
          const topicEasy = topic.subtopics.reduce((s, st) => s + st.easy, 0)
          const topicMedium = topic.subtopics.reduce((s, st) => s + st.medium, 0)
          const topicHard = topic.subtopics.reduce((s, st) => s + st.hard, 0)
          return (
            <div key={topic.name} className="ca-dsa-topic">
              {/* Topic header — clickable */}
              <button
                type="button"
                className="ca-dsa-topic-header"
                onClick={() => toggle(topic.name)}
                aria-expanded={isOpen}
              >
                <div className="ca-dsa-topic-left">
                  <span className={`ca-dsa-priority ca-dsa-priority-${topic.priority.toLowerCase()}`}>
                    {topic.priority}
                  </span>
                  <strong>{topic.name}</strong>
                  <span className="ca-dsa-sub-count">{topic.subtopics.length} subtopics</span>
                </div>
                <div className="ca-dsa-topic-right">
                  <span className="ca-dsa-pill ca-dsa-pill-easy">{topicEasy}E</span>
                  <span className="ca-dsa-pill ca-dsa-pill-medium">{topicMedium}M</span>
                  <span className="ca-dsa-pill ca-dsa-pill-hard">{topicHard}H</span>
                  <span className="ca-dsa-total">{topicTotal} total</span>
                  {isOpen ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                </div>
              </button>

              {/* Subtopics */}
              {isOpen && (
                <div className="ca-dsa-subtopics">
                  {/* Table header */}
                  <div className="ca-dsa-subtopic-row ca-dsa-subtopic-header">
                    <span>Subtopic / Pattern</span>
                    <span>Importance</span>
                    <span>Easy</span>
                    <span>Medium</span>
                    <span>Hard</span>
                    <span>Total</span>
                  </div>
                  {topic.subtopics.map((sub) => (
                    <div key={sub.name} className="ca-dsa-subtopic-row">
                      <div className="ca-dsa-subtopic-name">
                        <span>{sub.name}</span>
                        {sub.notes && <small>{sub.notes}</small>}
                      </div>
                      <span className={`ca-dsa-imp ca-dsa-imp-${sub.importance.toLowerCase()}`}>
                        {sub.importance}
                      </span>
                      <span className="ca-dsa-num ca-dsa-num-easy">{sub.easy}</span>
                      <span className="ca-dsa-num ca-dsa-num-medium">{sub.medium}</span>
                      <span className="ca-dsa-num ca-dsa-num-hard">{sub.hard}</span>
                      <span className="ca-dsa-num ca-dsa-num-total">{sub.total}</span>
                    </div>
                  ))}
                  {/* Topic totals row */}
                  <div className="ca-dsa-subtopic-row ca-dsa-subtopic-totals">
                    <span>Topic Total</span>
                    <span />
                    <span className="ca-dsa-num ca-dsa-num-easy">{topicEasy}</span>
                    <span className="ca-dsa-num ca-dsa-num-medium">{topicMedium}</span>
                    <span className="ca-dsa-num ca-dsa-num-hard">{topicHard}</span>
                    <span className="ca-dsa-num ca-dsa-num-total">{topicTotal}</span>
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>

      {/* Grand total summary */}
      <div className="ca-dsa-summary">
        <div className="ca-dsa-summary-item">
          <span className="ca-dsa-summary-dot ca-dsa-pill-easy" />
          <span>Easy</span>
          <strong>{topics.reduce((s, t) => s + t.subtopics.reduce((ss, st) => ss + st.easy, 0), 0)}</strong>
        </div>
        <div className="ca-dsa-summary-item">
          <span className="ca-dsa-summary-dot ca-dsa-pill-medium" />
          <span>Medium</span>
          <strong>{topics.reduce((s, t) => s + t.subtopics.reduce((ss, st) => ss + st.medium, 0), 0)}</strong>
        </div>
        <div className="ca-dsa-summary-item">
          <span className="ca-dsa-summary-dot ca-dsa-pill-hard" />
          <span>Hard</span>
          <strong>{topics.reduce((s, t) => s + t.subtopics.reduce((ss, st) => ss + st.hard, 0), 0)}</strong>
        </div>
        <div className="ca-dsa-summary-item ca-dsa-summary-total">
          <span>Total Problems</span>
          <strong>{totalQuestions}</strong>
        </div>
      </div>
    </div>
  )
}

function SectionCard({ section }: { section: AnalysisSection }) {
  const icon = sectionIcons[section.title] ?? <BookOpen size={18} />
  return (
    <article className="ca-section-card">
      <div className="ca-section-header">
        <span className="ca-section-icon">{icon}</span>
        <h3>{section.title}</h3>
      </div>
      {section.content ? <p className="ca-section-content">{section.content}</p> : null}
      {section.points?.length > 0 && (
        <ul className="ca-section-points">
          {section.points.map((point, idx) => (
            <li key={idx}>
              <ChevronRight size={14} />
              <span>{point}</span>
            </li>
          ))}
        </ul>
      )}
    </article>
  )
}

function TopicBadge({ topic }: { topic: AnalysisTopic }) {
  const cfg = importanceConfig[topic.importance] ?? importanceConfig.Medium
  return (
    <article className={`ca-topic-badge ${cfg.cls}`}>
      <div className="ca-topic-top">
        <strong>{topic.name}</strong>
        <span className="ca-topic-imp">{cfg.label}</span>
      </div>
      <div className="ca-topic-bottom">
        <span>{topic.minProbs} problems</span>
        {topic.notes ? <span className="ca-topic-notes">{topic.notes}</span> : null}
      </div>
    </article>
  )
}

export function CompanyAnalysis() {
  const [query, setQuery] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [analysis, setAnalysis] = useState<CompanyAnalysis | null>(null)

  async function runAnalysis(name: string) {
    const trimmed = name.trim()
    if (!trimmed) return
    setLoading(true)
    setError('')
    setAnalysis(null)

    try {
      const response = await fetch('/api/ai/company-analysis', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ companyName: trimmed }),
      })
      if (!response.ok) {
        const body = (await response.json()) as { error?: string }
        throw new Error(body.error ?? 'Analysis failed')
      }
      const data = (await response.json()) as CompanyAnalysis
      setAnalysis(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    void runAnalysis(query)
  }

  function handleSuggestion(company: string) {
    setQuery(company)
    void runAnalysis(company)
  }

  return (
    <div className="ca-root">
      {/* Search hero */}
      <section className="ca-hero">
        <div className="ca-hero-badge">
          <Sparkles size={14} />
          AI-Powered Analysis
        </div>
        <h1 className="ca-hero-title">Company Deep Dive</h1>
        <p className="ca-hero-sub">
          Enter any company name. PUYRG AI will analyze their interview style, required topics,
          readiness weights, and give you the exact preparation roadmap.
        </p>

        <form className="ca-search-form" onSubmit={handleSubmit}>
          <div className="ca-search-box">
            <Search size={20} />
            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="e.g. Google, Jane Street, Rubrik, OpenAI..."
              aria-label="Company name"
              autoFocus
            />
            <button type="submit" className="ca-search-btn" disabled={loading || !query.trim()}>
              {loading ? <Loader2 size={18} className="ca-spin" /> : <><Brain size={16} /> Analyze</>}
            </button>
          </div>
        </form>

        {/* Suggestion chips */}
        {!analysis && !loading && (
          <div className="ca-suggestions">
            {SUGGESTED_COMPANIES.map((company) => (
              <button
                key={company}
                type="button"
                className="ca-chip"
                onClick={() => handleSuggestion(company)}
              >
                {company}
              </button>
            ))}
          </div>
        )}
      </section>

      {/* Error */}
      {error && (
        <div className="ca-error">
          <AlertCircle size={16} />
          {error}
        </div>
      )}

      {/* Loading skeleton */}
      {loading && (
        <div className="ca-loading">
          <Loader2 size={36} className="ca-spin ca-spin-big" />
          <p>PUYRG AI is analyzing <strong>{query}</strong>...</p>
          <span>Checking company tier, interview patterns, topic requirements, and roadmap</span>
        </div>
      )}

      {/* Analysis results */}
      {analysis && !loading && (
        <div className="ca-results">
          {/* Header */}
          <div className="ca-results-header">
            <div className="ca-company-badge">
              <Building2 size={28} />
            </div>
            <div>
              <h2 className="ca-company-name">{analysis.companyName}</h2>
              <div className="ca-meta-row">
                <span className="ca-tier-badge">{analysis.tier}</span>
                <span className="ca-track-badge">
                  <Trophy size={12} />
                  {analysis.careerTrack}
                </span>
              </div>
            </div>
            {analysis.aiUsed && (
              <div className="ca-ai-badge">
                <Sparkles size={13} />
                OpenAI
              </div>
            )}
          </div>

          <p className="ca-overview">{analysis.overview}</p>

          {/* Readiness weights + Interview format */}
          <div className="ca-top-grid">
            <div className="ca-panel">
              <div className="ca-panel-label">
                <Target size={16} /> Readiness Weights
              </div>
              <div className="ca-weights-inner">
                <ResponsiveContainer width="100%" height={200}>
                  <PieChart>
                    <Pie
                      data={analysis.readinessWeights}
                      dataKey="value"
                      innerRadius={52}
                      outerRadius={84}
                      paddingAngle={3}
                    >
                      {analysis.readinessWeights.map((entry) => (
                        <Cell key={entry.name} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip formatter={(val) => `${val}%`} />
                  </PieChart>
                </ResponsiveContainer>
                <div className="ca-weight-legend">
                  {analysis.readinessWeights.map((w) => (
                    <div key={w.name} className="ca-weight-item">
                      <span className="ca-weight-dot" style={{ background: w.color }} />
                      <span className="ca-weight-name">{w.name}</span>
                      <strong>{w.value}%</strong>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="ca-panel">
              <div className="ca-panel-label">
                <Flame size={16} /> Highest Impact Tasks
              </div>
              <ol className="ca-impact-list">
                {analysis.highestImpact.map((task, idx) => (
                  <li key={idx}>
                    <span className="ca-impact-num">{idx + 1}</span>
                    <span>{task}</span>
                  </li>
                ))}
              </ol>
            </div>

            <div className="ca-panel ca-panel-full">
              <div className="ca-panel-label">
                <Users size={16} /> Interview Format
              </div>
              <p className="ca-interview-format">{analysis.interviewFormat}</p>
            </div>
          </div>

          {/* Required topics */}
          <div className="ca-topics-section">
            <div className="ca-section-title">
              <Brain size={18} />
              Required Topics
              <span className="ca-topic-count">{analysis.requiredTopics.length} topics</span>
            </div>
            <div className="ca-topics-filter">
              <span className="ca-filter-label">
                <span className="dot-critical" /> Critical: {analysis.requiredTopics.filter(t => t.importance === 'Critical').length}
              </span>
              <span className="ca-filter-label">
                <span className="dot-high" /> High: {analysis.requiredTopics.filter(t => t.importance === 'High').length}
              </span>
              <span className="ca-filter-label">
                <span className="dot-medium" /> Medium: {analysis.requiredTopics.filter(t => t.importance === 'Medium').length}
              </span>
            </div>
            <div className="ca-topics-grid">
              {analysis.requiredTopics.map((topic) => (
                <TopicBadge key={topic.name} topic={topic} />
              ))}
            </div>
          </div>

          {/* DSA Topics Breakdown */}
          {analysis.dsaTopics?.length > 0 && (
            <DSATopicsSection topics={analysis.dsaTopics} />
          )}

          {/* Sections */}
          <div className="ca-sections-grid">
            {analysis.sections.map((section) => (
              <SectionCard key={section.title} section={section} />
            ))}
          </div>

          {/* Re-search */}
          <div className="ca-re-search">
            <button
              type="button"
              className="ca-re-btn"
              onClick={() => {
                setAnalysis(null)
                setQuery('')
              }}
            >
              <Search size={16} /> Analyze another company
            </button>
            <div className="ca-suggestions ca-suggestions-bottom">
              {SUGGESTED_COMPANIES.filter((c) => c !== analysis.companyName).slice(0, 6).map((company) => (
                <button
                  key={company}
                  type="button"
                  className="ca-chip"
                  onClick={() => handleSuggestion(company)}
                >
                  {company}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Empty idle state */}
      {!analysis && !loading && !error && (
        <div className="ca-idle">
          <div className="ca-idle-icon">
            <Building2 size={32} />
          </div>
          <p>Search for any company above to get a full PUYRG AI analysis.</p>
          <ul className="ca-idle-list">
            <li><CheckCircle2 size={14} /> Interview format and rounds</li>
            <li><CheckCircle2 size={14} /> DSA topics with subtopics and question counts</li>
            <li><CheckCircle2 size={14} /> Required topics with minimum problem targets</li>
            <li><CheckCircle2 size={14} /> Readiness weight breakdown</li>
            <li><CheckCircle2 size={14} /> Highest impact preparation tasks</li>
            <li><CheckCircle2 size={14} /> Common mistakes and how to avoid them</li>
          </ul>
        </div>
      )}
    </div>
  )
}
