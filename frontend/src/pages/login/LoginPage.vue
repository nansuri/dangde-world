<template>
  <div class="login-page">
    <div class="login-hero">
      <div class="hero-badge">Kahoot-inspired learning UI</div>
      <h2>Learning feels like a game, and every role gets its own experience.</h2>
    </div>

    <LoginPanel :users="users" :loading="loading" :error="error" @submit="handleLogin" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import LoginPanel from '../../features/auth/LoginPanel.vue'
import { loginUser, listUsers } from '../../entities/user/api.js'
import { saveSession } from '../../features/auth/session.js'

const router = useRouter()
const users = ref([])
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  const response = await listUsers()
  users.value = response.items
})

async function handleLogin(userId) {
  if (!userId) {
    error.value = 'Select a demo user first.'
    return
  }

  loading.value = true
  error.value = ''
  try {
    const response = await loginUser(userId)
    saveSession(response.user)
    router.push(response.user.role === 'kid' ? '/kids' : response.user.role === 'parent' ? '/parents' : '/admin')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>
