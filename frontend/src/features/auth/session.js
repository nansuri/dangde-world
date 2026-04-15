const STORAGE_KEY = 'dangde-world-session'

export function readSession() {
  const raw = localStorage.getItem(STORAGE_KEY)
  return raw ? JSON.parse(raw) : null
}

export function saveSession(user) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(user))
}

export function clearSession() {
  localStorage.removeItem(STORAGE_KEY)
}
