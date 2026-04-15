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

export function createUser(payload) {
  return http('/users', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateUser(id, payload) {
  return http(`/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export function deleteUser(id) {
  return http(`/users/${id}`, {
    method: 'DELETE',
  })
}
