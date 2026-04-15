import { http } from '../../shared/api/http.js'

export function listActivities(params = {}) {
  const query = new URLSearchParams(params).toString()
  return http(`/activities${query ? `?${query}` : ''}`)
}

export function createActivity(payload) {
  return http('/activities', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateActivity(id, payload) {
  return http(`/activities/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}
