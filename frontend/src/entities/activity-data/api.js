import { http } from '../../shared/api/http.js'

export function listActivityData(params = {}) {
  const query = new URLSearchParams(params).toString()
  return http(`/activity-data${query ? `?${query}` : ''}`)
}

export function createActivityData(payload) {
  return http('/activity-data', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateActivityData(id, payload) {
  return http(`/activity-data/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export function deleteActivityData(id) {
  return http(`/activity-data/${id}`, {
    method: 'DELETE',
  })
}
