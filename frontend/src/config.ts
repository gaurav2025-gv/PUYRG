// Backend API base URL
// In development: empty string (Vite proxy handles /api/*)
// In production: full Render URL
export const API_BASE =
  import.meta.env.VITE_API_URL ?? 'https://puyrg-ai.onrender.com'

export function apiUrl(path: string): string {
  return `${API_BASE}${path}`
}
