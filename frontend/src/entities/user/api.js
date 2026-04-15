import { http } from '../../shared/api/http.js'

export function loginUser(userId) {
  return http('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ userId }),
  })
}

export function listUsers(params = {}) {
  const query = new URLSearchParams(params).toString()
  return http(`/users${query ? `?${query}` : ''}`)
}
