import {
  AlertCircle,
  BookOpen,
  Brain,
  Building2,
  ChevronDown,
  ChevronRight,
  Flame,
  Layers,
  Loader2,
  RefreshCw,
  Sparkles,
  Target,
  TrendingUp,
  Trophy,
  Users,
  ZoomIn,
} from 'lucide-react'
import type { ReactNode } from 'react'
import { useState } from 'react'
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'
import './GoalProfile.css'

// ── Types ─────────────────────────────────────────────────────────────────────

export type GoalProfileData = {
  companyName: string
  tier: string
  overview: string
  careerTrack: string
  interviewFormat: string
  readinessWeights: { name: string; value: number; color: string }[]
  highestImpact: string[]
  requiredTopics: {
    name: string
    importance: 'Critical' | 'High' | 'Medium'
    minProbs: number
    notes: string
    solved?: number
  }[]
  dsaTopics: {
    name: string
    priority: 'Critical' | 'High' | 'Medium'
    solved?: number           // total solved across all subtopics
    subtopics: {
      name: string
      importance: 'Critical' | 'High' | 'Medium'
      easy: number
      medium: number
      hard: number
      total: number
      notes: string
      solved?: number         // solved count for this subtopic
    }[]
    drill?: TopicDrill
  }[]
  sections: { title: string; content: string; points: string[] }[]
  aiUsed: boolean
  savedAt: string
}

type DrillPattern = {
  name: string
  priority: 'Critical' | 'High' | 'Medium' | 'Low'
  easy: number
  medium: number
  hard: number
  total: number
  notes: string
}

type SubDrill = {
  categoryName: string
  topicName: string
  companyName: string
  totalProbs: number
  patterns: DrillPattern[]
}

type DrillCategory = {
  name: string
  priority: 'Critical' | 'High' | 'Medium' | 'Low'
  description: string
  patterns: DrillPattern[]
  subDrill?: SubDrill
}

type TopicDrill = {
  topicName: string
  companyName: string
  totalProbs: number
  categories: DrillCategory[]
  aiUsed: boolean
}

// ── Constants ─────────────────────────────────────────────────────────────────

const SECTION_ICONS: Record<string, ReactNode> = {
  'Interview Process': <Users size={16} />,
  'DSA Focus Areas': <Brain size={16} />,
  'System Design Expectations': <Layers size={16} />,
  'Behavioral & Culture': <Users size={16} />,
  'Preparation Timeline': <TrendingUp size={16} />,
  'Common Mistakes': <AlertCircle size={16} />,
}

const IMP_CLASS: Record<string, string> = {
  Critical: 'gp-imp-critical',
  High: 'gp-imp-high',
  Medium: 'gp-imp-medium',
  Low: 'gp-imp-low',
}

const PRIORITY_CLASS: Record<string, string> = {
  Critical: 'gp-priority-critical',
  High: 'gp-priority-high',
  Medium: 'gp-priority-medium',
  Low: 'gp-priority-low',
}

// ── Drill-down category accordion ─────────────────────────────────────────────

// ── DrillCategoryRow — with further expand ────────────────────────────────────

function DrillCategoryRow({
  cat,
  topicName,
  companyName,
  onSubDrillSave,
}: {
  cat: DrillCategory & { subDrill?: SubDrill }
  topicName: string
  companyName: string
  onSubDrillSave: (categoryName: string, subDrill: SubDrill) => void
}) {
  const [open, setOpen] = useState(cat.priority === 'Critical' || cat.priority === 'High')
  const [subLoading, setSubLoading] = useState(false)
  const [subError, setSubError] = useState('')
  const [subExpanded, setSubExpanded] = useState(false)

  const catTotal = cat.patterns.reduce((s, p) => s + p.total, 0)
  const catEasy = cat.patterns.reduce((s, p) => s + p.easy, 0)
  const catMedium = cat.patterns.reduce((s, p) => s + p.medium, 0)
  const catHard = cat.patterns.reduce((s, p) => s + p.hard, 0)

  async function expandFurther() {
    if (cat.subDrill) {
      setSubExpanded((v) => !v)
      return
    }
    setSubLoading(true)
    setSubError('')
    try {
      const resp = await fetch('/api/ai/topic-drill', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          topicName: `${topicName} → ${cat.name}`,
          companyName,
        }),
      })
      if (!resp.ok) {
        const b = (await resp.json()) as { error?: string }
        throw new Error(b.error ?? 'Failed')
      }
      const data = (await resp.json()) as TopicDrill
      // Flatten all patterns from all returned categories into one SubDrill
      const allPatterns = data.categories.flatMap((c) => c.patterns)
      const subDrill: SubDrill = {
        categoryName: cat.name,
        topicName,
        companyName,
        totalProbs: allPatterns.reduce((s, p) => s + p.total, 0),
        patterns: allPatterns,
      }
      onSubDrillSave(cat.name, subDrill)
      setSubExpanded(true)
    } catch (err) {
      setSubError(err instanceof Error ? err.message : 'Failed')
    } finally {
      setSubLoading(false)
    }
  }

  return (
    <div className="drill-cat-card">
      <button
        type="button"
        className="drill-cat-header"
        onClick={() => setOpen((v) => !v)}
        aria-expanded={open}
      >
        <div className="drill-cat-left">
          <span className={`gp-priority ${PRIORITY_CLASS[cat.priority] ?? ''}`}>
            {cat.priority}
          </span>
          <strong>{cat.name}</strong>
          {cat.description && <span className="drill-cat-desc">{cat.description}</span>}
          {cat.subDrill && (
            <span className="drill-badge">
              <ZoomIn size={10} /> {cat.subDrill.patterns.length} deep patterns
            </span>
          )}
        </div>
        <div className="drill-cat-right">
          <span className="gp-pill gp-pill-easy">{catEasy}E</span>
          <span className="gp-pill gp-pill-medium">{catMedium}M</span>
          <span className="gp-pill gp-pill-hard">{catHard}H</span>
          <span className="drill-cat-total">{catTotal}</span>
          <span className="drill-cat-count">{cat.patterns.length} patterns</span>
          {open ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
        </div>
      </button>

      {open && (
        <div className="drill-patterns-table">
          <div className="drill-pt-row drill-pt-head">
            <span>Pattern</span>
            <span>Priority</span>
            <span>Easy</span>
            <span>Med</span>
            <span>Hard</span>
            <span>Total</span>
          </div>
          {cat.patterns.map((p) => (
            <div key={p.name} className="drill-pt-row">
              <div className="drill-pt-name">
                <span>{p.name}</span>
                {p.notes && <small>{p.notes}</small>}
              </div>
              <span className={`gp-imp-badge ${IMP_CLASS[p.priority] ?? ''}`}>{p.priority}</span>
              <span className="gp-n gp-n-easy">{p.easy}</span>
              <span className="gp-n gp-n-medium">{p.medium}</span>
              <span className="gp-n gp-n-hard">{p.hard}</span>
              <span className="gp-n gp-n-total">{p.total}</span>
            </div>
          ))}
          <div className="drill-pt-row drill-pt-totals">
            <span>Subtotal</span><span />
            <span className="gp-n gp-n-easy">{catEasy}</span>
            <span className="gp-n gp-n-medium">{catMedium}</span>
            <span className="gp-n gp-n-hard">{catHard}</span>
            <span className="gp-n gp-n-total">{catTotal}</span>
          </div>

          {/* Further expand button */}
          <div className="drill-expand-row">
            {subError && <span className="drill-err"><AlertCircle size={12} /> {subError}</span>}
            <button
              type="button"
              className={`drill-expand-btn ${cat.subDrill ? 'drill-expand-btn-done' : ''}`}
              onClick={(e) => { e.stopPropagation(); void expandFurther() }}
              disabled={subLoading}
            >
              {subLoading ? (
                <><Loader2 size={13} className="spin-icon" /> Expanding {cat.name}...</>
              ) : cat.subDrill ? (
                subExpanded
                  ? <><ChevronDown size={13} /> Hide deep patterns ({cat.subDrill.patterns.length})</>
                  : <><ZoomIn size={13} /> Show deep patterns ({cat.subDrill.patterns.length})</>
              ) : (
                <><ZoomIn size={13} /> Expand {cat.name} further — all atomic patterns</>
              )}
            </button>
          </div>

          {/* Sub-drill patterns table */}
          {cat.subDrill && subExpanded && (
            <div className="sub-drill-view">
              <div className="drill-pt-row drill-pt-head">
                <span>Atomic Pattern</span>
                <span>Priority</span>
                <span>Easy</span>
                <span>Med</span>
                <span>Hard</span>
                <span>Total</span>
              </div>
              {cat.subDrill.patterns.map((p) => (
                <div key={p.name} className="drill-pt-row">
                  <div className="drill-pt-name">
                    <span>{p.name}</span>
                    {p.notes && <small>{p.notes}</small>}
                  </div>
                  <span className={`gp-imp-badge ${IMP_CLASS[p.priority] ?? ''}`}>{p.priority}</span>
                  <span className="gp-n gp-n-easy">{p.easy}</span>
                  <span className="gp-n gp-n-medium">{p.medium}</span>
                  <span className="gp-n gp-n-hard">{p.hard}</span>
                  <span className="gp-n gp-n-total">{p.total}</span>
                </div>
              ))}
              <div className="drill-pt-row drill-pt-totals">
                <span>Deep Total</span><span />
                <span className="gp-n gp-n-easy">{cat.subDrill.patterns.reduce((s, p) => s + p.easy, 0)}</span>
                <span className="gp-n gp-n-medium">{cat.subDrill.patterns.reduce((s, p) => s + p.medium, 0)}</span>
                <span className="gp-n gp-n-hard">{cat.subDrill.patterns.reduce((s, p) => s + p.hard, 0)}</span>
                <span className="gp-n gp-n-total">{cat.subDrill.totalProbs}</span>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

// ── DrillView — full topic deep breakdown ─────────────────────────────────────

function DrillView({
  drill,
  onSubDrillSave,
}: {
  drill: TopicDrill
  onSubDrillSave: (categoryName: string, subDrill: SubDrill) => void
}) {
  const critCount = drill.categories.filter((c) => c.priority === 'Critical').length
  const highCount = drill.categories.filter((c) => c.priority === 'High').length
  const totalPatterns = drill.categories.reduce((s, c) => s + c.patterns.length, 0)

  return (
    <div className="drill-view">
      <div className="drill-stats-bar">
        <div className="drill-stat">
          <strong>{drill.categories.length}</strong><span>categories</span>
        </div>
        <div className="drill-stat">
          <strong>{totalPatterns}</strong><span>patterns</span>
        </div>
        <div className="drill-stat drill-stat-critical">
          <strong>{critCount}</strong><span>critical</span>
        </div>
        <div className="drill-stat drill-stat-high">
          <strong>{highCount}</strong><span>high priority</span>
        </div>
        <div className="drill-stat drill-stat-total">
          <strong>{drill.totalProbs}</strong><span>total problems</span>
        </div>
      </div>

      <div className="drill-cats">
        {drill.categories.map((cat) => (
          <DrillCategoryRow
            key={cat.name}
            cat={cat}
            topicName={drill.topicName}
            companyName={drill.companyName}
            onSubDrillSave={onSubDrillSave}
          />
        ))}
      </div>

      <div className="drill-grand-total">
        <Sparkles size={13} />
        <span>
          Complete <strong>{drill.topicName}</strong> breakdown ·{' '}
          {drill.totalProbs} problems across {totalPatterns} patterns
        </span>
      </div>
    </div>
  )
}

// ── TopicAccordion with drill-down ────────────────────────────────────────────

type TopicAccordionProps = {
  topics: GoalProfileData['dsaTopics']
  companyName: string
  onDrillSave: (topicName: string, drill: TopicDrill) => void
}

function TopicAccordion({ topics, companyName, onDrillSave }: TopicAccordionProps) {
  const [open, setOpen] = useState<Set<string>>(
    () => new Set(topics.filter((t) => t.priority === 'Critical').map((t) => t.name)),
  )
  const [drillLoading, setDrillLoading] = useState<Set<string>>(new Set())
  const [drillError, setDrillError] = useState<Record<string, string>>({})
  const [drillExpanded, setDrillExpanded] = useState<Set<string>>(new Set())

  function toggleTopic(name: string) {
    setOpen((prev) => {
      const next = new Set(prev)
      next.has(name) ? next.delete(name) : next.add(name)
      return next
    })
  }

  async function expandTopic(topicName: string) {
    // If already has drill data, just toggle show/hide
    const topic = topics.find((t) => t.name === topicName)
    if (topic?.drill) {
      setDrillExpanded((prev) => {
        const next = new Set(prev)
        next.has(topicName) ? next.delete(topicName) : next.add(topicName)
        return next
      })
      return
    }

    setDrillLoading((prev) => new Set(prev).add(topicName))
    setDrillError((prev) => ({ ...prev, [topicName]: '' }))

    try {
        const controller = new AbortController()
        const timeout = setTimeout(() => controller.abort(), 175000)
        const resp = await fetch('/api/ai/topic-drill', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ topicName, companyName }),
          signal: controller.signal,
        })
        clearTimeout(timeout)
      if (!resp.ok) {
        const body = (await resp.json()) as { error?: string }
        throw new Error(body.error ?? 'Drill failed')
      }
      const drill = (await resp.json()) as TopicDrill
      onDrillSave(topicName, drill)
      setDrillExpanded((prev) => new Set(prev).add(topicName))
    } catch (err) {
      setDrillError((prev) => ({
        ...prev,
        [topicName]: err instanceof Error ? err.message : 'Failed to expand',
      }))
    } finally {
      setDrillLoading((prev) => {
        const next = new Set(prev)
        next.delete(topicName)
        return next
      })
    }
  }

  const grandTotal = topics.reduce((s, t) => {
    if (t.drill) {
      return s + t.drill.totalProbs
    }
    return s + t.subtopics.reduce((ss, st) => ss + st.total, 0)
  }, 0)

  const grandEasy = topics.reduce((s, t) => s + t.subtopics.reduce((ss, st) => ss + st.easy, 0), 0)
  const grandMedium = topics.reduce((s, t) => s + t.subtopics.reduce((ss, st) => ss + st.medium, 0), 0)
  const grandHard = topics.reduce((s, t) => s + t.subtopics.reduce((ss, st) => ss + st.hard, 0), 0)

  return (
    <div className="gp-accordion-wrap">
      <div className="gp-block-header">
        <Brain size={18} />
        <span>DSA Topics & Problem Targets</span>
        <span className="gp-badge">{topics.length} topics · {grandTotal.toLocaleString()} problems</span>
      </div>

      <div className="gp-topic-list">
        {topics.map((topic) => {
          const isOpen = open.has(topic.name)
          const isDrillLoading = drillLoading.has(topic.name)
          const isDrillExpanded = drillExpanded.has(topic.name)
          const hasDrill = !!topic.drill
          const err = drillError[topic.name]

          const tTotal = topic.subtopics.reduce((s, st) => s + st.total, 0)
          const tEasy = topic.subtopics.reduce((s, st) => s + st.easy, 0)
          const tMedium = topic.subtopics.reduce((s, st) => s + st.medium, 0)
          const tHard = topic.subtopics.reduce((s, st) => s + st.hard, 0)

          return (
            <div key={topic.name} className="gp-topic-card">
              {/* Topic header */}
              <button
                type="button"
                className="gp-topic-header"
                onClick={() => toggleTopic(topic.name)}
                aria-expanded={isOpen}
              >
                <div className="gp-topic-left">
                  <span className={`gp-priority ${PRIORITY_CLASS[topic.priority] ?? ''}`}>
                    {topic.priority}
                  </span>
                  <strong>{topic.name}</strong>
                  <span className="gp-sub-count">{topic.subtopics.length} patterns</span>
                  {hasDrill && (
                    <span className="drill-badge">
                      <ZoomIn size={10} /> {topic.drill!.categories.length} categories · {topic.drill!.totalProbs} probs
                    </span>
                  )}
                  {(topic.solved ?? 0) > 0 && (
                    <span className="gp-solved-badge">
                      ✓ {topic.solved} solved
                    </span>
                  )}
                </div>
                <div className="gp-topic-right">
                  {/* Progress bar */}
                  {tTotal > 0 && (
                    <div className="gp-topic-progress-wrap" title={`${topic.solved ?? 0} / ${tTotal} solved`}>
                      <div
                        className="gp-topic-progress-fill"
                        style={{ width: `${Math.min(100, ((topic.solved ?? 0) / tTotal) * 100)}%` }}
                      />
                    </div>
                  )}
                  <span className="gp-topic-progress-label">
                    {topic.solved ?? 0}<span className="gp-topic-progress-sep">/</span>{tTotal}
                  </span>
                  <span className="gp-pill gp-pill-easy">{tEasy}E</span>
                  <span className="gp-pill gp-pill-medium">{tMedium}M</span>
                  <span className="gp-pill gp-pill-hard">{tHard}H</span>
                  {isOpen ? <ChevronDown size={15} /> : <ChevronRight size={15} />}
                </div>
              </button>

              {/* Subtopics table (overview level) */}
              {isOpen && (
                <div className="gp-subtopics">
                  <div className="gp-st-row gp-st-head">
                    <span>Pattern</span>
                    <span>Importance</span>
                    <span>Solved</span>
                    <span>Easy</span>
                    <span>Med</span>
                    <span>Hard</span>
                    <span>Total</span>
                  </div>
                  {topic.subtopics.map((st) => {
                    const solvedCount = st.solved ?? 0
                    const pct = st.total > 0 ? Math.min(100, (solvedCount / st.total) * 100) : 0
                    const isDone = solvedCount >= st.total && st.total > 0
                    return (
                    <div key={st.name} className={`gp-st-row ${isDone ? 'gp-st-row-done' : ''}`}>
                      <div className="gp-st-name">
                        <span>{isDone ? '✓ ' : ''}{st.name}</span>
                        {st.notes && <small>{st.notes}</small>}
                        {st.total > 0 && (
                          <div className="gp-st-progress-bar">
                            <div className="gp-st-progress-fill" style={{ width: `${pct}%` }} />
                          </div>
                        )}
                      </div>
                      <span className={`gp-imp-badge ${IMP_CLASS[st.importance] ?? ''}`}>
                        {st.importance}
                      </span>
                      <span className={`gp-n ${solvedCount > 0 ? 'gp-n-solved' : 'gp-n-zero'}`}>
                        {solvedCount}
                      </span>
                      <span className="gp-n gp-n-easy">{st.easy}</span>
                      <span className="gp-n gp-n-medium">{st.medium}</span>
                      <span className="gp-n gp-n-hard">{st.hard}</span>
                      <span className="gp-n gp-n-total">{st.total}</span>
                    </div>
                    )
                  })}
                  <div className="gp-st-row gp-st-totals">
                    <span>Subtotal</span>
                    <span />
                    <span className="gp-n gp-n-solved">{topic.subtopics.reduce((s, st) => s + (st.solved ?? 0), 0)}</span>
                    <span className="gp-n gp-n-easy">{tEasy}</span>
                    <span className="gp-n gp-n-medium">{tMedium}</span>
                    <span className="gp-n gp-n-hard">{tHard}</span>
                    <span className="gp-n gp-n-total">{tTotal}</span>
                  </div>

                  {/* Expand button */}
                  <div className="drill-expand-row">
                    {err && <span className="drill-err"><AlertCircle size={12} /> {err}</span>}
                    <button
                      type="button"
                      className={`drill-expand-btn ${hasDrill ? 'drill-expand-btn-done' : ''}`}
                      onClick={(e) => { e.stopPropagation(); void expandTopic(topic.name) }}
                      disabled={isDrillLoading}
                    >
                      {isDrillLoading ? (
                        <><Loader2 size={13} className="spin-icon" /> Generating deep breakdown...</>
                      ) : hasDrill ? (
                        isDrillExpanded
                          ? <><ChevronDown size={13} /> Hide deep breakdown</>
                          : <><ZoomIn size={13} /> Show deep breakdown ({topic.drill!.categories.length} categories)</>
                      ) : (
                        <><ZoomIn size={13} /> Expand — get all {topic.name} patterns & subtopics</>
                      )}
                    </button>
                    {hasDrill && !isDrillLoading && (
                      <button
                        type="button"
                        className="drill-refresh-btn"
                        onClick={(e) => {
                          e.stopPropagation()
                          // Force regenerate
                          onDrillSave(topic.name, null as unknown as TopicDrill)
                          void expandTopic(topic.name)
                        }}
                        title="Regenerate"
                      >
                        <RefreshCw size={12} />
                      </button>
                    )}
                  </div>

                  {/* Drill-down view */}
                  {hasDrill && isDrillExpanded && topic.drill && (
                    <DrillView
                      drill={topic.drill}
                      onSubDrillSave={(categoryName, subDrill) => {
                        const updatedDrill: TopicDrill = {
                          ...topic.drill!,
                          categories: topic.drill!.categories.map((c) =>
                            c.name === categoryName ? { ...c, subDrill } : c,
                          ),
                        }
                        onDrillSave(topic.name, updatedDrill)
                      }}
                    />
                  )}
                </div>
              )}
            </div>
          )
        })}
      </div>

      {/* Grand total bar */}
      <div className="gp-grand-total">
        <div className="gp-gt-item"><span className="gp-gt-dot gp-gt-easy" /><span>Easy</span><strong>{grandEasy}</strong></div>
        <div className="gp-gt-item"><span className="gp-gt-dot gp-gt-medium" /><span>Medium</span><strong>{grandMedium}</strong></div>
        <div className="gp-gt-item"><span className="gp-gt-dot gp-gt-hard" /><span>Hard</span><strong>{grandHard}</strong></div>
        <div className="gp-gt-item gp-gt-highlight"><span>Total (overview)</span><strong>{grandTotal.toLocaleString()}</strong></div>
      </div>
    </div>
  )
}

// ── SectionCard ───────────────────────────────────────────────────────────────

function SectionCard({ title, content, points }: { title: string; content: string; points: string[] }) {
  const icon = SECTION_ICONS[title] ?? <BookOpen size={16} />
  return (
    <article className="gp-section-card">
      <div className="gp-section-head">
        <span className="gp-section-icon">{icon}</span>
        <h4>{title}</h4>
      </div>
      {content && <p className="gp-section-body">{content}</p>}
      {points.length > 0 && (
        <ul className="gp-section-points">
          {points.map((pt, i) => (
            <li key={i}><ChevronRight size={13} /><span>{pt}</span></li>
          ))}
        </ul>
      )}
    </article>
  )
}

// ── Main GoalProfile component ────────────────────────────────────────────────

type Props = {
  data: GoalProfileData
  onBack: () => void
  onUpdate: (updated: GoalProfileData) => void
}

export function GoalProfile({ data, onBack, onUpdate }: Props) {
  const totalProblems = data.dsaTopics.reduce((s, t) => {
    if (t.drill) return s + t.drill.totalProbs
    return s + t.subtopics.reduce((ss, st) => ss + st.total, 0)
  }, 0)

  const tierColor =
    data.tier.includes('1') ? '#7c5cff' :
    data.tier.includes('2') ? '#72e3ff' :
    data.tier.includes('3') ? '#ff8db5' :
    data.tier.includes('4') ? '#ffb36b' : '#9c83ff'

  function handleDrillSave(topicName: string, drill: TopicDrill) {
    const updated: GoalProfileData = {
      ...data,
      dsaTopics: data.dsaTopics.map((t) =>
        t.name === topicName ? { ...t, drill } : t,
      ),
    }
    onUpdate(updated)
  }

  return (
    <div className="gp-root">

      {/* Page header */}
      <div className="gp-page-header">
        <button type="button" className="gp-back-btn" onClick={onBack}>
          <ChevronRight size={16} style={{ transform: 'rotate(180deg)' }} />
          All Goals
        </button>

        <div className="gp-hero">
          <div className="gp-hero-icon"><Building2 size={32} /></div>
          <div className="gp-hero-text">
            <h1>{data.companyName}</h1>
            <div className="gp-hero-meta">
              <span className="gp-tier-badge" style={{ borderColor: tierColor, color: tierColor }}>
                {data.tier}
              </span>
              <span className="gp-track-badge">
                <Trophy size={11} /> {data.careerTrack}
              </span>
              {data.aiUsed && (
                <span className="gp-ai-badge"><Sparkles size={11} /> AI Generated</span>
              )}
            </div>
          </div>
          <div className="gp-hero-stats">
            <div className="gp-hero-stat">
              <strong>{totalProblems.toLocaleString()}</strong>
              <span>total problems</span>
            </div>
            <div className="gp-hero-stat">
              <strong>{data.dsaTopics.length}</strong>
              <span>DSA topics</span>
            </div>
            <div className="gp-hero-stat gp-hero-stat-solved">
              <strong>{data.dsaTopics.reduce((s, t) => s + (t.solved ?? 0), 0)}</strong>
              <span>solved</span>
            </div>
            <div className="gp-hero-stat">
              <strong>{data.dsaTopics.filter(t => t.drill).length}</strong>
              <span>drilled deep</span>
            </div>
          </div>
        </div>

        <p className="gp-overview">{data.overview}</p>
      </div>

      {/* Top grid: weights + impact */}
      <div className="gp-top-grid">
        <div className="gp-panel">
          <div className="gp-block-header"><Target size={16} /><span>Readiness Weights</span></div>
          <div className="gp-weights-row">
            <ResponsiveContainer width={180} height={180}>
              <PieChart>
                <Pie data={data.readinessWeights} dataKey="value" innerRadius={48} outerRadius={80} paddingAngle={3}>
                  {data.readinessWeights.map((w) => (
                    <Cell key={w.name} fill={w.color} />
                  ))}
                </Pie>
                <Tooltip formatter={(v) => `${v}%`} />
              </PieChart>
            </ResponsiveContainer>
            <div className="gp-weight-legend">
              {data.readinessWeights.map((w) => (
                <div key={w.name} className="gp-weight-row">
                  <span className="gp-weight-dot" style={{ background: w.color }} />
                  <span className="gp-weight-label">{w.name}</span>
                  <div className="gp-weight-bar-wrap">
                    <div className="gp-weight-bar" style={{ width: `${w.value}%`, background: w.color }} />
                  </div>
                  <strong>{w.value}%</strong>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="gp-panel">
          <div className="gp-block-header"><Flame size={16} /><span>Highest Impact Tasks</span></div>
          <ol className="gp-impact-list">
            {data.highestImpact.map((task, i) => (
              <li key={i}>
                <span className="gp-impact-num">{i + 1}</span>
                <span>{task}</span>
              </li>
            ))}
          </ol>
        </div>

        <div className="gp-panel gp-panel-full">
          <div className="gp-block-header"><Users size={16} /><span>Interview Format</span></div>
          <p className="gp-format-text">{data.interviewFormat}</p>
        </div>
      </div>

      {/* Required topics grid */}
      <div className="gp-required-wrap">
        <div className="gp-block-header">
          <Brain size={18} /><span>Required Topics</span>
          <span className="gp-badge">{data.requiredTopics.length} topics</span>
          <div className="gp-imp-legend">
            <span className="gp-imp-dot gp-imp-dot-critical" /> Critical: {data.requiredTopics.filter(t => t.importance === 'Critical').length}
            <span className="gp-imp-dot gp-imp-dot-high" /> High: {data.requiredTopics.filter(t => t.importance === 'High').length}
            <span className="gp-imp-dot gp-imp-dot-medium" /> Medium: {data.requiredTopics.filter(t => t.importance === 'Medium').length}
          </div>
        </div>
        <div className="gp-req-grid">
          {data.requiredTopics.map((t) => (
            <article key={t.name} className={`gp-req-card gp-req-${t.importance.toLowerCase()}`}>
              <div className="gp-req-top">
                <strong>{t.name}</strong>
                <span className={`gp-imp-badge ${IMP_CLASS[t.importance] ?? ''}`}>{t.importance}</span>
              </div>
              <div className="gp-req-bottom">
                <span className="gp-req-probs">{t.minProbs} problems</span>
                {t.notes && <span className="gp-req-notes">{t.notes}</span>}
              </div>
            </article>
          ))}
        </div>
      </div>

      {/* DSA Topics accordion with drill-down */}
      {data.dsaTopics.length > 0 && (
        <TopicAccordion
          topics={data.dsaTopics}
          companyName={data.companyName}
          onDrillSave={handleDrillSave}
        />
      )}

      {/* Sections */}
      {data.sections.length > 0 && (
        <div className="gp-sections-wrap">
          <div className="gp-block-header"><BookOpen size={18} /><span>Detailed Breakdown</span></div>
          <div className="gp-sections-grid">
            {data.sections.map((s) => (
              <SectionCard key={s.title} title={s.title} content={s.content} points={s.points} />
            ))}
          </div>
        </div>
      )}

      {/* Footer */}
      <div className="gp-footer">
        <Sparkles size={13} />
        Generated by PUYRG AI · {new Date(data.savedAt).toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' })}
      </div>
    </div>
  )
}
