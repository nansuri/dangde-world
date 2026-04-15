import { createRouter, createWebHistory } from 'vue-router'
import { readSession } from '../features/auth/session.js'
import LoginPage from '../pages/login/LoginPage.vue'
import KidsPage from '../pages/kids/KidsPage.vue'
import ActivityPlayerPage from '../pages/kids/ActivityPlayerPage.vue'
import ParentsPage from '../pages/parents/ParentsPage.vue'
import AdminPage from '../pages/admin/AdminPage.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: LoginPage },
  { path: '/kids', component: KidsPage, meta: { role: 'kid' } },
  { path: '/kids/activity/:assignmentId', component: ActivityPlayerPage, meta: { role: 'kid' } },
  { path: '/parents', component: ParentsPage, meta: { role: 'parent' } },
  { path: '/admin', component: AdminPage, meta: { role: 'admin' } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (!to.meta.role) {
    return true
  }

  const session = readSession()
  if (!session) {
    return '/login'
  }
  if (session.role !== to.meta.role) {
    return session.role === 'kid' ? '/kids' : session.role === 'parent' ? '/parents' : '/admin'
  }
  return true
})

export default router
