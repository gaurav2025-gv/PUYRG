import {
  Activity,
  BookOpenCheck,
  Brain,
  Building2,
  CheckCircle2,
  Clock3,
  Flame,
  Loader2,
  Plus,
  RefreshCcw,
  Search,
  ShieldCheck,
  Sparkles,
  Target,
  Trash2,
  TrendingUp,
  X,
} from 'lucide-react'
import type { ReactNode } from 'react'
import { useMemo, useRef, useState } from 'react'
import { CompanyAnalysis } from './CompanyAnalysis'
import { GoalProfile } from './GoalProfile'
import type { GoalProfileData } from './GoalProfile'
import { HistoryPage } from './HistoryPage'
import { ProblemLogger } from './ProblemLogger'
import type { LoggedProblem } from './ProblemLogger'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import './App.css'

import { apiUrl } from './config'

// ── Revision types ────────────────────────────────────────────────────────────

type RevisionEntry = {
  id: string
  problemId: string
  problemTitle: string
  topic: string
  pattern: string
  difficulty: string
  qualityScore: number
  revisionNum: 1 | 2 | 3
  dueAt: string
  completedAt?: string
  loggedAt: string
}

// ── Storage helpers ───────────────────────────────────────────────────────────

const GOALS_KEY = 'puyrg_goal_profiles'
const PROBLEMS_KEY = 'puyrg_logged_problems'
const REVISIONS_KEY = 'puyrg_revisions'

function ls<T>(key: string): T[] {
  try { return JSON.parse(localStorage.getItem(key) ?? '[]') as T[] } catch { return [] }
}
function ss<T>(key: string, val: T[]) {
  localStorage.setItem(key, JSON.stringify(val))
}

// R1=3d, R2=12d, R3=45d
function revisionDueDate(loggedAt: string, num: 1 | 2 | 3): string {
  const gaps = { 1: 3, 2: 12, 3: 45 }
  const d = new Date(loggedAt)
  d.setDate(d.getDate() + gaps[num])
  return d.toISOString()
}

function isDue(dueAt: string) { return new Date(dueAt) <= new Date() }
function isOverdue(dueAt: string) {
  const d = new Date(dueAt); d.setDate(d.getDate() - 1); return d <= new Date()
}

// ── Mini components ───────────────────────────────────────────────────────────

function StatCard({ icon, label, value, detail, accent }: {
  icon: ReactNode; label: string; value: string; detail: string; accent?: string
}) {
  return (
    <section className="metric">
      <div className="metric-icon">{icon}</div>
      <div>
        <p>{label}</p>
        <strong style={accent ? { color: accent } : undefined}>{value}</strong>
        <span>{detail}</span>
      </div>
    </section>
  )
}

function ProgressBar({ value, color }: { value: number; color?: string }) {
  return (
    <div className="bar-track" aria-label={`${value}%`}>
      <span className="bar-fill" style={{ width: `${value}%`, background: color }} />
    </div>
  )
}

function EmptyState({ children }: { children: ReactNode }) {
  return <div className="empty-state">{children}</div>
}

type Page = 'dashboard' | 'company-analysis' | 'goal-profile' | 'history'

// ── App ───────────────────────────────────────────────────────────────────────

export default function App() {
  const [activePage, setActivePage] = useState<Page>('dashboard')
  const [activeGoal, setActiveGoal] = useState<GoalProfileData | null>(null)
  const [goals, setGoals] = useState<GoalProfileData[]>(() => ls<GoalProfileData>(GOALS_KEY))
  const [problems, setProblems] = useState<LoggedProblem[]>(() => ls<LoggedProblem>(PROBLEMS_KEY))
  const [revisions, setRevisions] = useState<RevisionEntry[]>(() => ls<RevisionEntry>(REVISIONS_KEY))
  const [goalInput, setGoalInput] = useState('')
  const [goalLoading, setGoalLoading] = useState(false)
  const [goalError, setGoalError] = useState('')
  const goalInputRef = useRef<HTMLInputElement>(null)

  // ── Derived stats ─────────────────────────────────────────────────────────

  const totalQuality = useMemo(() => problems.reduce((s, p) => s + p.qualityScore, 0), [problems])
  const topicsSet = useMemo(() => new Set(problems.map(p => p.topic)), [problems])
  const patternsSet = useMemo(() => new Set(problems.map(p => p.pattern)), [problems])

  const dueRevisions = useMemo(
    () => revisions.filter(r => !r.completedAt && isDue(r.dueAt)),
    [revisions]
  )
  const overdueRevisions = useMemo(
    () => dueRevisions.filter(r => isOverdue(r.dueAt)),
    [dueRevisions]
  )

  // Topic distribution for heatmap
  const topicChartData = useMemo(() => {
    const map: Record<string, number> = {}
    for (const p of problems) {
      map[p.topic] = (map[p.topic] ?? 0) + 1
    }
    return Object.entries(map)
      .map(([name, solved]) => ({ name, solved }))
      .sort((a, b) => b.solved - a.solved)
      .slice(0, 10)
  }, [problems])

  // Goal readiness weights (from first goal or default)
  const readinessMix = useMemo(() => {
    if (goals.length === 0) return []
    return goals[0].readinessWeights ?? []
  }, [goals])

  // ── Revision logic ────────────────────────────────────────────────────────

  function scheduleRevisions(problem: LoggedProblem) {
    const newRevs: RevisionEntry[] = [1, 2, 3].map((num) => ({
      id: `${problem.id}-r${num}`,
      problemId: problem.id,
      problemTitle: problem.problemTitle,
      topic: problem.topic,
      pattern: problem.pattern,
      difficulty: problem.difficulty,
      qualityScore: problem.qualityScore,
      revisionNum: num as 1 | 2 | 3,
      dueAt: revisionDueDate(problem.loggedAt, num as 1 | 2 | 3),
      loggedAt: problem.loggedAt,
    }))
    const updated = [...revisions, ...newRevs]
    setRevisions(updated)
    ss(REVISIONS_KEY, updated)
  }

  function completeRevision(rev: RevisionEntry) {
    const updated = revisions.map(r =>
      r.id === rev.id ? { ...r, completedAt: new Date().toISOString() } : r
    )
    setRevisions(updated)
    ss(REVISIONS_KEY, updated)
  }

  // ── Problem logging ───────────────────────────────────────────────────────

  function handleProblemLogged(problem: LoggedProblem) {
    const updatedProblems = [problem, ...problems]
    setProblems(updatedProblems)
    ss(PROBLEMS_KEY, updatedProblems)

    // Schedule R1/R2/R3
    scheduleRevisions(problem)

    // Auto-update goal profiles
    const tL = problem.topic.toLowerCase()
    const sL = problem.subtopic.toLowerCase()
    const pL = problem.pattern.toLowerCase()
    const kw = (s: string) => s.toLowerCase().replace(/[^a-z0-9 ]/g, '').split(/\s+/).filter(w => w.length > 2)
    const tKw = kw(tL); const sKw = kw(sL); const pKw = kw(pL)
    function fz(a: string, b: string): boolean {
      if (a === b || a.includes(b) || b.includes(a)) return true
      return kw(a).some(k => kw(b).includes(k))
    }

    const updatedGoals = goals.map(goal => {
      let changed = false
      const newTopics = goal.dsaTopics.map(t => {
        const tn = t.name.toLowerCase()
        const topicHit = fz(tn, tL) || fz(tn, sL) || tKw.some(k => tn.includes(k)) || sKw.some(k => tn.includes(k))
        if (!topicHit) return t
        let topicChanged = false
        const newSubs = t.subtopics.map(st => {
          const sn = st.name.toLowerCase()
          const hit = fz(sn, pL) || fz(sn, sL) || pKw.some(k => sn.includes(k)) || sKw.some(k => sn.includes(k)) || kw(sn).some(k => pKw.includes(k) || sKw.includes(k))
          if (!hit) return st
          topicChanged = true; changed = true
          return { ...st, solved: (st.solved ?? 0) + 1 }
        })
        if (!topicChanged && t.subtopics.length > 0) {
          topicChanged = true; changed = true
          const first = { ...t.subtopics[0], solved: (t.subtopics[0].solved ?? 0) + 1 }
          const merged = [first, ...t.subtopics.slice(1)]
          return { ...t, subtopics: merged, solved: merged.reduce((s, st) => s + (st.solved ?? 0), 0) }
        }
        return { ...t, subtopics: newSubs, solved: newSubs.reduce((s, st) => s + (st.solved ?? 0), 0) }
      })
      if (!changed) return goal
      return { ...goal, dsaTopics: newTopics }
    })
    setGoals(updatedGoals)
    ss(GOALS_KEY, updatedGoals)
    if (activeGoal) {
      const fresh = updatedGoals.find(g => g.companyName === activeGoal.companyName)
      if (fresh) setActiveGoal(fresh)
    }
  }

  // ── Goal profile logic ────────────────────────────────────────────────────

  async function generateGoalProfile(name: string) {
    const trimmed = name.trim()
    if (!trimmed) return
    const existing = goals.find(g => g.companyName.toLowerCase() === trimmed.toLowerCase())
    if (existing) { setActiveGoal(existing); setActivePage('goal-profile'); setGoalInput(''); return }
    setGoalLoading(true); setGoalError('')
    try {
      const ctrl = new AbortController()
      const t = setTimeout(() => ctrl.abort(), 175000)
      const resp = await fetch(apiUrl('/api/ai/company-analysis'), {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ companyName: trimmed }), signal: ctrl.signal,
      })
      clearTimeout(t)
      if (!resp.ok) { const b = await resp.json() as { error?: string }; throw new Error(b.error ?? 'Failed') }
      const data = await resp.json() as Omit<GoalProfileData, 'savedAt'>
      const profile: GoalProfileData = { ...data, savedAt: new Date().toISOString() }
      const updated = [profile, ...goals.filter(g => g.companyName.toLowerCase() !== trimmed.toLowerCase())]
      setGoals(updated); ss(GOALS_KEY, updated)
      setActiveGoal(profile); setActivePage('goal-profile'); setGoalInput('')
    } catch (err) { setGoalError(err instanceof Error ? err.message : 'Something went wrong') }
    finally { setGoalLoading(false) }
  }

  function deleteGoal(name: string, e: React.MouseEvent) {
    e.stopPropagation()
    const updated = goals.filter(g => g.companyName !== name)
    setGoals(updated); ss(GOALS_KEY, updated)
    if (activeGoal?.companyName === name) { setActivePage('dashboard'); setActiveGoal(null) }
  }

  // ── Render ────────────────────────────────────────────────────────────────

  return (
    <main className="app-shell">
      {/* Sidebar */}
      <aside className="sidebar" aria-label="Primary navigation">
        <div className="brand">
          <div className="brand-mark">P</div>
          <div><strong>PUYRG</strong><span>Practice Until You Reach Goal</span></div>
        </div>

        <nav>
          <a className={activePage === 'dashboard' ? 'active' : ''} href="#" onClick={e => { e.preventDefault(); setActivePage('dashboard') }}>
            <Activity size={18} /> Dashboard
          </a>
          <a className={activePage === 'history' ? 'active' : ''} href="#" onClick={e => { e.preventDefault(); setActivePage('history') }}>
            <BookOpenCheck size={18} /> Problem History
            {problems.length > 0 && <span className="nav-count">{problems.length}</span>}
          </a>
          <a className={activePage === 'company-analysis' ? 'active' : ''} href="#" onClick={e => { e.preventDefault(); setActivePage('company-analysis') }}>
            <Building2 size={18} /> Company Analysis
          </a>
        </nav>

        {/* My Goals */}
        <div className="goals-section">
          <div className="goals-header"><Sparkles size={13} /><span>MY GOALS</span></div>
          <form className="goals-add-form" onSubmit={e => { e.preventDefault(); void generateGoalProfile(goalInput) }}>
            <input ref={goalInputRef} value={goalInput} onChange={e => setGoalInput(e.target.value)}
              placeholder="Add company / exam..." className="goals-add-input" disabled={goalLoading} />
            <button type="submit" className="goals-add-btn" disabled={goalLoading || !goalInput.trim()}>
              {goalLoading ? <Loader2 size={14} className="spin-icon" /> : <Plus size={14} />}
            </button>
          </form>
          {goalError && <div className="goals-error"><X size={12} /> {goalError}</div>}
          {goalLoading && <div className="goals-loading"><Loader2 size={13} className="spin-icon" /><span>Generating {goalInput}...</span></div>}
          <div className="goals-list">
            {goals.map(g => (
              <button key={g.companyName} type="button"
                className={`goals-item ${activeGoal?.companyName === g.companyName && activePage === 'goal-profile' ? 'goals-item-active' : ''}`}
                onClick={() => { setActiveGoal(g); setActivePage('goal-profile') }}>
                <Building2 size={13} /><span>{g.companyName}</span>
                <button type="button" className="goals-del" onClick={e => deleteGoal(g.companyName, e)}><Trash2 size={11} /></button>
              </button>
            ))}
            {goals.length === 0 && !goalLoading && <span className="goals-empty">No goals yet. Add one above.</span>}
          </div>
        </div>
      </aside>

      {/* Workspace */}
      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">Interview Operating System</p>
            <h1>Track. Revise. Master.</h1>
          </div>
          <div className="search-box">
            <Search size={18} />
            <input placeholder="Search topic, company, problem" />
          </div>
        </header>

        {activePage === 'company-analysis' ? <CompanyAnalysis /> :
         activePage === 'history' ? <HistoryPage problems={problems} /> :
         activePage === 'goal-profile' && activeGoal ? (
          <GoalProfile data={activeGoal} onBack={() => { setActivePage('dashboard'); setActiveGoal(null) }}
            onUpdate={updated => {
              setActiveGoal(updated)
              const u = goals.map(g => g.companyName === updated.companyName ? updated : g)
              setGoals(u); ss(GOALS_KEY, u)
            }} />
        ) : (
          <>
            {/* Stats */}
            <section className="metrics-grid">
              <StatCard icon={<BookOpenCheck size={20} />} label="Solved problems" value={String(problems.length)}
                detail={`${totalQuality} quality points`} />
              <StatCard icon={<ShieldCheck size={20} />} label="Hard problems"
                value={String(problems.filter(p => p.difficulty === 'Hard').length)}
                detail={`${problems.filter(p => p.difficulty === 'Medium').length} med · ${problems.filter(p => p.difficulty === 'Easy').length} easy`}
                accent="#ff8db5" />
              <StatCard icon={<RefreshCcw size={20} />} label="Revisions due"
                value={String(dueRevisions.length)}
                detail={`${overdueRevisions.length} overdue · ${revisions.filter(r => r.completedAt).length} done`}
                accent={dueRevisions.length > 0 ? '#ffb36b' : undefined} />
              <StatCard icon={<Target size={20} />} label="Topics covered"
                value={String(topicsSet.size)}
                detail={`${patternsSet.size} patterns practiced`} />
            </section>

            {/* Main content grid */}
            <section className="content-grid">

              {/* Log problem */}
              <section className="panel">
                <ProblemLogger onProblemLogged={handleProblemLogged} recentProblems={problems} />
              </section>

              {/* Revision Engine */}
              <section className="panel wide">
                <div className="panel-heading">
                  <div><p className="eyebrow">R1 · R2 · R3 Spaced Repetition</p><h2>Revision Queue</h2></div>
                  <RefreshCcw size={20} />
                </div>
                {dueRevisions.length === 0 ? (
                  <EmptyState>
                    {revisions.length === 0
                      ? 'Log a solved problem to auto-schedule R1 (3 days), R2 (12 days), R3 (45 days).'
                      : `All caught up! Next revision in ${(() => {
                          const pending = revisions.filter(r => !r.completedAt).sort((a, b) => new Date(a.dueAt).getTime() - new Date(b.dueAt).getTime())
                          if (!pending.length) return '—'
                          const diff = Math.ceil((new Date(pending[0].dueAt).getTime() - Date.now()) / 86400000)
                          return `${diff} day${diff !== 1 ? 's' : ''}: ${pending[0].problemTitle}`
                        })()}`
                    }
                  </EmptyState>
                ) : (
                  <div className="revision-list">
                    {dueRevisions.slice(0, 8).map(rev => (
                      <article key={rev.id} className="revision-item">
                        <div>
                          <strong>{rev.problemTitle}</strong>
                          <span>{rev.topic} · {rev.pattern}</span>
                        </div>
                        <div className="revision-meta">
                          <span className={`priority ${isOverdue(rev.dueAt) ? 'overdue' : 'high'}`}>
                            {isOverdue(rev.dueAt) ? 'Overdue' : 'Due'}
                          </span>
                          <span>R{rev.revisionNum}</span>
                          <span>{rev.difficulty}</span>
                          <span>+{rev.qualityScore}pts</span>
                        </div>
                        <button className="inline-action" type="button" onClick={() => completeRevision(rev)}>
                          <CheckCircle2 size={16} /> Mark revised
                        </button>
                      </article>
                    ))}
                    {dueRevisions.length > 8 && (
                      <p style={{ color: 'var(--muted)', fontSize: 12, textAlign: 'center', margin: '8px 0 0' }}>
                        +{dueRevisions.length - 8} more due — go to Problem History for full list
                      </p>
                    )}
                  </div>
                )}
              </section>

              {/* Goal profiles readiness */}
              <section className="panel wide">
                <div className="panel-heading">
                  <div><p className="eyebrow">My Goals</p><h2>Company Readiness</h2></div>
                  <Building2 size={20} />
                </div>
                {goals.length === 0 ? (
                  <EmptyState>Add a company goal from the sidebar to see readiness breakdown.</EmptyState>
                ) : (
                  <div className="company-list">
                    {goals.map(g => {
                      const totalRequired = g.dsaTopics.reduce((s, t) => s + t.subtopics.reduce((ss2, st) => ss2 + st.total, 0), 0)
                      const totalSolved = g.dsaTopics.reduce((s, t) => s + (t.solved ?? 0), 0)
                      const pct = totalRequired > 0 ? Math.min(100, Math.round((totalSolved / totalRequired) * 100)) : 0
                      return (
                        <article key={g.companyName} className="company-row" style={{ cursor: 'pointer' }}
                          onClick={() => { setActiveGoal(g); setActivePage('goal-profile') }}>
                          <div><strong>{g.companyName}</strong><span>{g.tier}</span></div>
                          <div className="company-focus">
                            {g.readinessWeights.slice(0, 3).map(w => <span key={w.name}>{w.name}</span>)}
                          </div>
                          <div className="readiness">
                            <strong>{pct}%</strong>
                            <small>{totalSolved} / {totalRequired} problems</small>
                            <ProgressBar value={pct} />
                          </div>
                          <span className="inline-action">View Profile →</span>
                        </article>
                      )
                    })}
                  </div>
                )}
              </section>

              {/* Topic heatmap */}
              <section className="panel wide">
                <div className="panel-heading">
                  <div><p className="eyebrow">Practice distribution</p><h2>Topic Heatmap</h2></div>
                  <Activity size={20} />
                </div>
                {topicChartData.length === 0 ? (
                  <EmptyState>Log problems to see your topic distribution.</EmptyState>
                ) : (
                  <ResponsiveContainer width="100%" height={220}>
                    <BarChart data={topicChartData} margin={{ left: -20, right: 10, bottom: 30 }}>
                      <CartesianGrid strokeDasharray="3 3" stroke="rgba(177,162,219,0.1)" />
                      <XAxis dataKey="name" angle={-30} textAnchor="end" tick={{ fill: 'var(--muted)', fontSize: 11 }} tickLine={false} axisLine={false} />
                      <YAxis tick={{ fill: 'var(--muted)', fontSize: 11 }} tickLine={false} axisLine={false} />
                      <Tooltip contentStyle={{ background: 'var(--panel-strong)', border: '1px solid var(--line)', borderRadius: 10, fontSize: 12 }} />
                      <Bar dataKey="solved" radius={[4,4,0,0]} name="Problems">
                        {topicChartData.map((_, i) => <Cell key={i} fill={`hsl(${250 + i * 15}, 70%, 60%)`} />)}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                )}
              </section>

              {/* Readiness mix pie */}
              {readinessMix.length > 0 && (
                <section className="panel chart-panel">
                  <div className="panel-heading">
                    <div><p className="eyebrow">{goals[0]?.companyName} weights</p><h2>Readiness Mix</h2></div>
                    <TrendingUp size={20} />
                  </div>
                  <ResponsiveContainer width="100%" height={200}>
                    <PieChart>
                      <Pie data={readinessMix} dataKey="value" innerRadius={54} outerRadius={82} paddingAngle={3}>
                        {readinessMix.map(w => <Cell key={w.name} fill={w.color} />)}
                      </Pie>
                      <Tooltip formatter={v => `${v}%`} />
                    </PieChart>
                  </ResponsiveContainer>
                  <div className="legend">
                    {readinessMix.map(w => (
                      <span key={w.name}><i style={{ background: w.color }} />{w.name} {w.value}%</span>
                    ))}
                  </div>
                </section>
              )}

              {/* AI Mentor */}
              <section className="panel">
                <div className="panel-heading">
                  <div><p className="eyebrow">AI mentor</p><h2>Next best move</h2></div>
                  <Brain size={20} />
                </div>
                <p className="mentor-text">
                  {problems.length === 0
                    ? 'Start by logging one solved problem. PUYRG will auto-schedule R1/R2/R3 revisions and track your progress across all company profiles.'
                    : dueRevisions.length > 0
                      ? `You have ${dueRevisions.length} revision${dueRevisions.length > 1 ? 's' : ''} due (${overdueRevisions.length} overdue). Complete these before adding new problems — mastered problems count more than solved ones.`
                      : `Great work! ${problems.length} problems logged, ${topicsSet.size} topics covered. Focus on your weakest topics next.`}
                </p>
                <button className="primary-button" type="button" onClick={() => setActivePage('history')}>
                  <Flame size={18} /> View full history
                </button>
              </section>

              {/* Recent problems */}
              <section className="panel">
                <div className="panel-heading">
                  <div><p className="eyebrow">Recent activity</p><h2>Last solved</h2></div>
                  <BookOpenCheck size={20} />
                </div>
                <div className="ledger-list">
                  {problems.length === 0 ? (
                    <EmptyState>No problems logged yet.</EmptyState>
                  ) : (
                    problems.slice(0, 6).map(p => (
                      <article key={p.id}>
                        <div>
                          <strong>{p.problemTitle}</strong>
                          <span>{p.platform} · {p.topic} · {p.pattern}</span>
                        </div>
                        <small style={{ color: p.difficulty === 'Hard' ? '#ff8db5' : p.difficulty === 'Medium' ? '#ffb36b' : '#67e8b9' }}>
                          {p.difficulty} · +{p.qualityScore}pts
                        </small>
                      </article>
                    ))
                  )}
                </div>
              </section>

              {/* Revision stats */}
              <section className="panel">
                <div className="panel-heading">
                  <div><p className="eyebrow">Retention</p><h2>Revision Stats</h2></div>
                  <Clock3 size={20} />
                </div>
                {revisions.length === 0 ? (
                  <EmptyState>Revisions appear after logging your first problem.</EmptyState>
                ) : (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
                    {([1,2,3] as const).map(n => {
                      const total = revisions.filter(r => r.revisionNum === n).length
                      const done = revisions.filter(r => r.revisionNum === n && r.completedAt).length
                      const pct = total > 0 ? Math.round((done / total) * 100) : 0
                      return (
                        <div key={n} style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
                          <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 13 }}>
                            <span style={{ color: 'var(--text)', fontWeight: 750 }}>
                              R{n} ({n === 1 ? '3 days' : n === 2 ? '12 days' : '45 days'})
                            </span>
                            <span style={{ color: 'var(--muted)' }}>{done}/{total} · {pct}%</span>
                          </div>
                          <ProgressBar value={pct} color={n === 1 ? '#7c5cff' : n === 2 ? '#72e3ff' : '#67e8b9'} />
                        </div>
                      )
                    })}
                    <p style={{ color: 'var(--muted)', fontSize: 12, margin: '4px 0 0' }}>
                      Mastered: {revisions.filter(r => r.revisionNum === 3 && r.completedAt).length} problems completed all 3 revisions
                    </p>
                  </div>
                )}
              </section>

            </section>
          </>
        )}
      </section>
    </main>
  )
}
