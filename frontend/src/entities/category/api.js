import { http } from '../../shared/api/http.js'

export function listCategories(params = {}) {
  const query = new URLSearchParams(params).toString()
  return http(`/categories${query ? `?${query}` : ''}`)
}

export function createCategory(data) {
  return http('/categories', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function updateCategory(id, data) {
  return http(`/categories/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export function deleteCategory(id) {
  return http(`/categories/${id}`, {
    method: 'DELETE',
  })
}
