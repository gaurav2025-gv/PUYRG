import {
  AlertCircle,
  BookOpen,
  Brain,
  ChevronDown,
  ChevronRight,
  Clock,
  Filter,
  Flame,
  Search,
  Sparkles,
  TrendingUp,
  Trophy,
} from 'lucide-react'
import { useMemo, useState } from 'react'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import type { LoggedProblem } from './ProblemLogger'
import './HistoryPage.css'

// ── Helpers ───────────────────────────────────────────────────────────────────

const DIFF_COLOR: Record<string, string> = {
  Easy: '#67e8b9',
  Medium: '#ffb36b',
  Hard: '#ff8db5',
}

function fmtDate(iso: string) {
  const d = new Date(iso)
  return d.toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' })
}

function fmtRelative(iso: string) {
  const diff = Date.now() - new Date(iso).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  const days = Math.floor(hrs / 24)
  if (days < 7) return `${days}d ago`
  return fmtDate(iso)
}

// ── Stats row ─────────────────────────────────────────────────────────────────

function StatsRow({ problems }: { problems: LoggedProblem[] }) {
  const total = problems.length
  const totalQuality = problems.reduce((s, p) => s + p.qualityScore, 0)
  const easy = problems.filter(p => p.difficulty === 'Easy').length
  const medium = problems.filter(p => p.difficulty === 'Medium').length
  const hard = problems.filter(p => p.difficulty === 'Hard').length
  const avgQuality = total > 0 ? (totalQuality / total).toFixed(1) : '0'
  const topicCount = new Set(problems.map(p => p.topic)).size

  return (
    <div className="hp-stats-row">
      <div className="hp-stat">
        <strong>{total}</strong>
        <span>Total Solved</span>
      </div>
      <div className="hp-stat">
        <strong style={{ color: '#67e8b9' }}>{easy}</strong>
        <span>Easy</span>
      </div>
      <div className="hp-stat">
        <strong style={{ color: '#ffb36b' }}>{medium}</strong>
        <span>Medium</span>
      </div>
      <div className="hp-stat">
        <strong style={{ color: '#ff8db5' }}>{hard}</strong>
        <span>Hard</span>
      </div>
      <div className="hp-stat">
        <strong style={{ color: 'var(--orange)' }}>{totalQuality}</strong>
        <span>Quality Pts</span>
      </div>
      <div className="hp-stat">
        <strong style={{ color: 'var(--purple-2)' }}>{avgQuality}</strong>
        <span>Avg Quality</span>
      </div>
      <div className="hp-stat">
        <strong style={{ color: 'var(--cyan)' }}>{topicCount}</strong>
        <span>Topics</span>
      </div>
    </div>
  )
}

// ── Topic heatmap ─────────────────────────────────────────────────────────────

function TopicHeatmap({ problems }: { problems: LoggedProblem[] }) {
  const topicData = useMemo(() => {
    const map: Record<string, { solved: number; quality: number }> = {}
    for (const p of problems) {
      if (!map[p.topic]) map[p.topic] = { solved: 0, quality: 0 }
      map[p.topic].solved++
      map[p.topic].quality += p.qualityScore
    }
    return Object.entries(map)
      .map(([name, { solved, quality }]) => ({ name, solved, quality }))
      .sort((a, b) => b.solved - a.solved)
      .slice(0, 12)
  }, [problems])

  if (topicData.length === 0) return null

  return (
    <div className="hp-panel">
      <div className="hp-panel-title">
        <Brain size={16} /> Topic Coverage Heatmap
      </div>
      <ResponsiveContainer width="100%" height={200}>
        <BarChart data={topicData} margin={{ top: 4, right: 8, left: -20, bottom: 40 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="rgba(177,162,219,0.1)" />
          <XAxis
            dataKey="name"
            angle={-35}
            textAnchor="end"
            tick={{ fill: 'var(--muted)', fontSize: 11 }}
            tickLine={false}
            axisLine={false}
          />
          <YAxis tick={{ fill: 'var(--muted)', fontSize: 11 }} tickLine={false} axisLine={false} />
          <Tooltip
            contentStyle={{ background: 'var(--panel-strong)', border: '1px solid var(--line)', borderRadius: 12, fontSize: 12 }}
            formatter={(val, name) => [val, name === 'solved' ? 'Problems' : 'Quality']}
          />
          <Bar dataKey="solved" radius={[4, 4, 0, 0]}>
            {topicData.map((_, i) => (
              <Cell key={i} fill={`hsl(${250 + i * 12}, 70%, ${55 + i * 2}%)`} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}

// ── Monthly trend ─────────────────────────────────────────────────────────────

function MonthlyTrend({ problems }: { problems: LoggedProblem[] }) {
  const trendData = useMemo(() => {
    const map: Record<string, { solved: number; quality: number }> = {}
    for (const p of problems) {
      const key = new Date(p.loggedAt).toLocaleDateString('en-IN', { month: 'short', year: '2-digit' })
      if (!map[key]) map[key] = { solved: 0, quality: 0 }
      map[key].solved++
      map[key].quality += p.qualityScore
    }
    return Object.entries(map)
      .map(([month, { solved, quality }]) => ({ month, solved, quality }))
      .slice(-8)
  }, [problems])

  if (trendData.length === 0) return null

  return (
    <div className="hp-panel">
      <div className="hp-panel-title">
        <TrendingUp size={16} /> Monthly Progress
      </div>
      <ResponsiveContainer width="100%" height={160}>
        <BarChart data={trendData} margin={{ top: 4, right: 8, left: -20, bottom: 4 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="rgba(177,162,219,0.1)" />
          <XAxis dataKey="month" tick={{ fill: 'var(--muted)', fontSize: 11 }} tickLine={false} axisLine={false} />
          <YAxis tick={{ fill: 'var(--muted)', fontSize: 11 }} tickLine={false} axisLine={false} />
          <Tooltip
            contentStyle={{ background: 'var(--panel-strong)', border: '1px solid var(--line)', borderRadius: 12, fontSize: 12 }}
          />
          <Bar dataKey="solved" fill="var(--purple)" radius={[4, 4, 0, 0]} name="Problems" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}

// ── Problem card ──────────────────────────────────────────────────────────────

function ProblemCard({ problem, expanded, onToggle }: {
  problem: LoggedProblem
  expanded: boolean
  onToggle: () => void
}) {
  const diffColor = DIFF_COLOR[problem.difficulty] ?? '#9c83ff'

  return (
    <article className="hp-problem-card">
      <button type="button" className="hp-problem-header" onClick={onToggle}>
        <div className="hp-problem-left">
          <span className="hp-diff-badge" style={{ background: diffColor + '22', color: diffColor }}>
            {problem.difficulty[0]}
          </span>
          <div className="hp-problem-info">
            <span className="hp-problem-title">{problem.problemTitle}</span>
            <span className="hp-problem-meta">
              {problem.platform} · {problem.topic} › {problem.pattern}
            </span>
          </div>
        </div>
        <div className="hp-problem-right">
          <span className="hp-quality-badge">+{problem.qualityScore}pts</span>
          <span className="hp-problem-date">{fmtRelative(problem.loggedAt)}</span>
          {expanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
        </div>
      </button>

      {expanded && (
        <div className="hp-problem-detail">
          <div className="hp-detail-row">
            <span className="hp-detail-label">Summary</span>
            <span>{problem.summary}</span>
          </div>
          <div className="hp-detail-row">
            <span className="hp-detail-label">Pattern</span>
            <span className="hp-detail-tag">{problem.pattern}</span>
          </div>
          <div className="hp-detail-row">
            <span className="hp-detail-label">Quality</span>
            <div className="hp-quality-track">
              <div className="hp-quality-fill" style={{ width: `${problem.qualityScore * 10}%` }} />
            </div>
            <span className="hp-quality-num">{problem.qualityScore}/10</span>
          </div>
          <div className="hp-detail-row">
            <span className="hp-detail-label">Quality Reason</span>
            <span className="hp-muted">{problem.qualityReason}</span>
          </div>
          {problem.notes && (
            <div className="hp-detail-row hp-notes-row">
              <span className="hp-detail-label">Notes</span>
              <p className="hp-notes-text">{problem.notes}</p>
            </div>
          )}
          <div className="hp-detail-row">
            <span className="hp-detail-label">Logged</span>
            <span className="hp-muted">{fmtDate(problem.loggedAt)}</span>
          </div>
        </div>
      )}
    </article>
  )
}

// ── Main HistoryPage ──────────────────────────────────────────────────────────

type Props = {
  problems: LoggedProblem[]
}

export function HistoryPage({ problems }: Props) {
  const [search, setSearch] = useState('')
  const [filterDiff, setFilterDiff] = useState<string>('All')
  const [filterTopic, setFilterTopic] = useState<string>('All')
  const [filterPlatform, setFilterPlatform] = useState<string>('All')
  const [sortBy, setSortBy] = useState<'date' | 'quality' | 'topic'>('date')
  const [expanded, setExpanded] = useState<string | null>(null)
  const [showFilters, setShowFilters] = useState(false)

  const allTopics = useMemo(() => ['All', ...Array.from(new Set(problems.map(p => p.topic))).sort()], [problems])
  const allPlatforms = useMemo(() => ['All', ...Array.from(new Set(problems.map(p => p.platform))).sort()], [problems])

  const filtered = useMemo(() => {
    let list = [...problems]
    if (search.trim()) {
      const q = search.toLowerCase()
      list = list.filter(p =>
        p.problemTitle.toLowerCase().includes(q) ||
        p.topic.toLowerCase().includes(q) ||
        p.pattern.toLowerCase().includes(q) ||
        p.notes?.toLowerCase().includes(q)
      )
    }
    if (filterDiff !== 'All') list = list.filter(p => p.difficulty === filterDiff)
    if (filterTopic !== 'All') list = list.filter(p => p.topic === filterTopic)
    if (filterPlatform !== 'All') list = list.filter(p => p.platform === filterPlatform)

    if (sortBy === 'date') list.sort((a, b) => new Date(b.loggedAt).getTime() - new Date(a.loggedAt).getTime())
    else if (sortBy === 'quality') list.sort((a, b) => b.qualityScore - a.qualityScore)
    else list.sort((a, b) => a.topic.localeCompare(b.topic))

    return list
  }, [problems, search, filterDiff, filterTopic, filterPlatform, sortBy])

  const totalQuality = problems.reduce((s, p) => s + p.qualityScore, 0)

  return (
    <div className="hp-root">
      {/* Header */}
      <div className="hp-header">
        <div className="hp-header-left">
          <div className="hp-eyebrow"><BookOpen size={14} /> Learning Ledger</div>
          <h2 className="hp-title">Problem History</h2>
          <p className="hp-subtitle">
            {problems.length} problems logged · {totalQuality} quality points · {new Set(problems.map(p => p.topic)).size} topics covered
          </p>
        </div>
        <div className="hp-header-badges">
          <span className="hp-badge hp-badge-purple">
            <Sparkles size={11} /> AI Tracked
          </span>
          <span className="hp-badge hp-badge-green">
            <Trophy size={11} /> {problems.filter(p => p.qualityScore >= 5).length} Hard+
          </span>
        </div>
      </div>

      {/* Stats */}
      {problems.length > 0 && <StatsRow problems={problems} />}

      {/* Charts */}
      {problems.length >= 3 && (
        <div className="hp-charts-grid">
          <TopicHeatmap problems={problems} />
          <MonthlyTrend problems={problems} />
        </div>
      )}

      {/* Filters bar */}
      <div className="hp-filters-bar">
        <div className="hp-search-wrap">
          <Search size={14} />
          <input
            className="hp-search"
            value={search}
            onChange={e => setSearch(e.target.value)}
            placeholder="Search problems, topics, patterns, notes..."
          />
        </div>

        <button
          type="button"
          className={`hp-filter-toggle ${showFilters ? 'hp-filter-toggle-active' : ''}`}
          onClick={() => setShowFilters(v => !v)}
        >
          <Filter size={13} /> Filters
        </button>

        <div className="hp-sort-wrap">
          <Clock size={13} />
          <select className="hp-select" value={sortBy} onChange={e => setSortBy(e.target.value as typeof sortBy)}>
            <option value="date">Latest First</option>
            <option value="quality">Highest Quality</option>
            <option value="topic">By Topic</option>
          </select>
        </div>
      </div>

      {showFilters && (
        <div className="hp-filter-chips">
          <div className="hp-filter-group">
            <span className="hp-filter-label">Difficulty</span>
            {['All', 'Easy', 'Medium', 'Hard'].map(d => (
              <button
                key={d}
                type="button"
                className={`hp-chip ${filterDiff === d ? 'hp-chip-active' : ''}`}
                onClick={() => setFilterDiff(d)}
              >
                {d}
              </button>
            ))}
          </div>
          <div className="hp-filter-group">
            <span className="hp-filter-label">Platform</span>
            {allPlatforms.map(p => (
              <button
                key={p}
                type="button"
                className={`hp-chip ${filterPlatform === p ? 'hp-chip-active' : ''}`}
                onClick={() => setFilterPlatform(p)}
              >
                {p}
              </button>
            ))}
          </div>
          <div className="hp-filter-group">
            <span className="hp-filter-label">Topic</span>
            {allTopics.map(t => (
              <button
                key={t}
                type="button"
                className={`hp-chip ${filterTopic === t ? 'hp-chip-active' : ''}`}
                onClick={() => setFilterTopic(t)}
              >
                {t}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Results count */}
      <div className="hp-results-count">
        {filtered.length === problems.length
          ? `${problems.length} problems`
          : `${filtered.length} of ${problems.length} problems`}
        {(search || filterDiff !== 'All' || filterTopic !== 'All' || filterPlatform !== 'All') && (
          <button
            type="button"
            className="hp-clear-filters"
            onClick={() => { setSearch(''); setFilterDiff('All'); setFilterTopic('All'); setFilterPlatform('All') }}
          >
            Clear filters
          </button>
        )}
      </div>

      {/* Problem list */}
      {problems.length === 0 ? (
        <div className="hp-empty">
          <AlertCircle size={32} />
          <p>No problems logged yet.</p>
          <span>Go to Dashboard → Log Solved Problem to start tracking.</span>
        </div>
      ) : filtered.length === 0 ? (
        <div className="hp-empty">
          <Search size={28} />
          <p>No problems match your filters.</p>
        </div>
      ) : (
        <div className="hp-problem-list">
          {filtered.map(p => (
            <ProblemCard
              key={p.id}
              problem={p}
              expanded={expanded === p.id}
              onToggle={() => setExpanded(expanded === p.id ? null : p.id)}
            />
          ))}
        </div>
      )}

      {/* Quality summary */}
      {problems.length > 0 && (
        <div className="hp-quality-summary">
          <Flame size={14} />
          <span>
            Total Quality Score: <strong>{totalQuality} points</strong> across {problems.length} problems.
            Average: <strong>{(totalQuality / problems.length).toFixed(1)}/10</strong>
          </span>
        </div>
      )}
    </div>
  )
}
