import { http } from '../../shared/api/http.js'

export function listAssignments(params = {}) {
  const query = new URLSearchParams(params).toString()
  return http(`/assignments${query ? `?${query}` : ''}`)
}

export function createAssignment(payload) {
  return http('/assignments', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateAssignment(id, payload) {
  return http(`/assignments/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(payload),
  })
}

export function getParentStats(parentId) {
  return http(`/stats/parent/${parentId}`)
}

export function deleteAssignment(id) {
  return http(`/assignments/${id}`, {
    method: 'DELETE',
  })
}
