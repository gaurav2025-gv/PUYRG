import {
  AlertCircle,
  CheckCircle2,
  ChevronRight,
  Clipboard,
  Image,
  Loader2,
  Sparkles,
  X,
  Zap,
} from 'lucide-react'
import type { FormEvent } from 'react'
import { useRef, useState } from 'react'
import './ProblemLogger.css'

// ── Types ─────────────────────────────────────────────────────────────────────

export type LoggedProblem = {
  id: string
  problemTitle: string
  platform: string
  difficulty: 'Easy' | 'Medium' | 'Hard'
  cfRating: number
  topic: string
  subtopic: string
  pattern: string
  qualityScore: number
  qualityReason: string
  summary: string
  confidence: number
  notes: string
  loggedAt: string
}

type UploadedImage = {
  name: string
  mimeType: string
  dataUrl: string
}

type Props = {
  onProblemLogged: (problem: LoggedProblem) => void
  recentProblems: LoggedProblem[]
}

const DIFFICULTY_COLOR: Record<string, string> = {
  Easy: '#67e8b9',
  Medium: '#ffb36b',
  Hard: '#ff8db5',
}

const QUALITY_LABEL = (score: number) => {
  if (score >= 9) return 'Elite'
  if (score >= 7) return 'Hard'
  if (score >= 5) return 'Medium'
  if (score >= 3) return 'Standard'
  return 'Basic'
}

export function ProblemLogger({ onProblemLogged, recentProblems }: Props) {
  const [text, setText] = useState('')
  const [notes, setNotes] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [preview, setPreview] = useState<LoggedProblem | null>(null)
  const [saved, setSaved] = useState(false)
  const [uploadedImage, setUploadedImage] = useState<UploadedImage | null>(null)
  const fileRef = useRef<HTMLInputElement>(null)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    const trimmed = text.trim()
    if (!trimmed && !uploadedImage) return
    setLoading(true)
    setError('')
    setPreview(null)
    setSaved(false)

    try {
      const resp = await fetch('/api/ai/log-problem', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          text: trimmed,
          imageData: uploadedImage?.dataUrl,
          imageMimeType: uploadedImage?.mimeType,
        }),
      })
      if (!resp.ok) {
        const b = (await resp.json()) as { error?: string }
        throw new Error(b.error ?? 'Analysis failed')
      }
      const data = (await resp.json()) as Omit<LoggedProblem, 'id' | 'loggedAt'>
      const problem: LoggedProblem = {
        ...data,
        notes: notes.trim(),
        id: `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
        loggedAt: new Date().toISOString(),
      }
      setPreview(problem)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  function confirmSave() {
    if (!preview) return
    onProblemLogged(preview)
    setSaved(true)
    setPreview(null)
    setText('')
    setNotes('')
    setUploadedImage(null)
    setTimeout(() => setSaved(false), 3000)
  }

  function handleImageUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = () => {
      setUploadedImage({
        name: file.name,
        mimeType: file.type || 'image/png',
        dataUrl: String(reader.result ?? ''),
      })
      setError('')
    }
    reader.onerror = () => setError('Could not read image file')
    reader.readAsDataURL(file)
    if (fileRef.current) fileRef.current.value = ''
  }

  const totalSolved = recentProblems.length
  const totalQuality = recentProblems.reduce((s, p) => s + p.qualityScore, 0)

  return (
    <div className="pl-root">
      <div className="pl-header">
        <div className="pl-title">
          <Zap size={18} />
          <span>Log Solved Problem</span>
        </div>
        <div className="pl-stats">
          <span className="pl-stat-chip">{totalSolved} logged</span>
          <span className="pl-stat-chip pl-stat-quality">{totalQuality} quality pts</span>
        </div>
      </div>

      <form className="pl-form" onSubmit={handleSubmit}>
        <div className="pl-textarea-wrap">
          <textarea
            className="pl-textarea"
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder={`Paste question text, problem name, or describe what you solved...

Examples:
• "Given an array, find the maximum sum subarray" 
• "LC 200 Number of Islands - DFS solution"
• "CF 1234A - solved in 25 min, easy greedy"
• "CSES - Dijkstra's Shortest Path"
`}
            rows={5}
            disabled={loading}
          />
          <div className="pl-textarea-actions">
            <button
              type="button"
              className="pl-icon-btn"
              title="Upload screenshot"
              onClick={() => fileRef.current?.click()}
              disabled={loading}
            >
              <Image size={15} />
            </button>
            <button
              type="button"
              className="pl-icon-btn"
              title="Paste from clipboard"
              onClick={async () => {
                try {
                  const t = await navigator.clipboard.readText()
                  setText((prev) => prev + (prev ? '\n' : '') + t)
                } catch { /* ignore */ }
              }}
              disabled={loading}
            >
              <Clipboard size={15} />
            </button>
          </div>
        </div>

        <input
          ref={fileRef}
          type="file"
          accept="image/*"
          className="pl-hidden-input"
          onChange={handleImageUpload}
        />

        {/* Image preview */}
        {uploadedImage && (
          <div className="pl-img-preview">
            <img src={uploadedImage.dataUrl} alt={uploadedImage.name} className="pl-img-thumb" />
            <span>{uploadedImage.name}</span>
            <button type="button" className="pl-img-remove" onClick={() => setUploadedImage(null)}>
              <X size={12} />
            </button>
          </div>
        )}

        {/* Notes */}
        <input
          className="pl-notes-input"
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          placeholder="Optional note (approach, mistakes, key insight...)"
          disabled={loading}
        />

        <button
          type="submit"
          className="pl-submit-btn"
          disabled={loading || (!text.trim() && !uploadedImage)}
        >
          {loading
            ? <><Loader2 size={15} className="spin-icon" /> Analyzing...</>
            : <><Sparkles size={15} /> Analyze & Log</>}
        </button>
      </form>

      {error && (
        <div className="pl-error">
          <AlertCircle size={14} /> {error}
        </div>
      )}

      {saved && (
        <div className="pl-success">
          <CheckCircle2 size={14} /> Problem logged! All matching profiles updated.
        </div>
      )}

      {/* Preview card — confirm before save */}
      {preview && (
        <div className="pl-preview-card">
          <div className="pl-preview-header">
            <span className="pl-preview-title">AI Detected</span>
            <span className="pl-confidence">{Math.round(preview.confidence * 100)}% confident</span>
            <button
              type="button"
              className="pl-dismiss"
              onClick={() => setPreview(null)}
            >
              <X size={14} />
            </button>
          </div>

          <h4 className="pl-preview-problem">{preview.problemTitle}</h4>
          <p className="pl-preview-summary">{preview.summary}</p>

          <div className="pl-preview-meta">
            <span className="pl-meta-chip pl-meta-platform">{preview.platform}</span>
            <span
              className="pl-meta-chip pl-meta-diff"
              style={{ color: DIFFICULTY_COLOR[preview.difficulty] }}
            >
              {preview.difficulty}
              {preview.cfRating > 0 && ` (CF ${preview.cfRating})`}
            </span>
            <span className="pl-meta-chip pl-meta-topic">
              {preview.topic} › {preview.subtopic}
            </span>
            <span className="pl-meta-chip pl-meta-pattern">{preview.pattern}</span>
          </div>

          <div className="pl-quality-bar">
            <div className="pl-quality-info">
              <Sparkles size={12} />
              <span>Quality Score: <strong>{preview.qualityScore}/10</strong> — {QUALITY_LABEL(preview.qualityScore)}</span>
            </div>
            <div className="pl-quality-track">
              <div
                className="pl-quality-fill"
                style={{ width: `${preview.qualityScore * 10}%` }}
              />
            </div>
            <small>{preview.qualityReason}</small>
          </div>

          <div className="pl-preview-actions">
            <button type="button" className="pl-confirm-btn" onClick={confirmSave}>
              <CheckCircle2 size={14} /> Confirm & Save
            </button>
            <button type="button" className="pl-cancel-btn" onClick={() => setPreview(null)}>
              Cancel
            </button>
          </div>
        </div>
      )}

      {/* Recent problems */}
      {recentProblems.length > 0 && (
        <div className="pl-recent">
          <div className="pl-recent-title">Recent</div>
          <div className="pl-recent-list">
            {recentProblems.slice(0, 5).map((p) => (
              <div key={p.id} className="pl-recent-item">
                <div className="pl-recent-left">
                  <span
                    className="pl-recent-diff"
                    style={{ background: DIFFICULTY_COLOR[p.difficulty] + '22', color: DIFFICULTY_COLOR[p.difficulty] }}
                  >
                    {p.difficulty[0]}
                  </span>
                  <div>
                    <span className="pl-recent-name">{p.problemTitle}</span>
                    <span className="pl-recent-meta">{p.topic} › {p.pattern}</span>
                  </div>
                </div>
                <div className="pl-recent-right">
                  <span className="pl-recent-quality">+{p.qualityScore}pts</span>
                  <ChevronRight size={12} />
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
