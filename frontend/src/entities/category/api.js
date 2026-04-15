import { http } from '../../shared/api/http.js'

export function listCategories(params = {}) {
  const query = new URLSearchParams(params).toString()
  return http(`/categories${query ? `?${query}` : ''}`)
}
