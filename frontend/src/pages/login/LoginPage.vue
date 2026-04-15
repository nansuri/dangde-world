<template>
  <div class="login-page">
    <GlassCard class="login-hero">
      <PillBadge>Kahoot-inspired learning UI</PillBadge>
      <h2>Learning feels like a game, and every role gets its own experience.</h2>
    </GlassCard>

    <LoginPanel :users="users" :loading="loading" :error="error" @submit="handleLogin" />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import GlassCard from '../../shared/ui/GlassCard.vue'
import PillBadge from '../../shared/ui/PillBadge.vue'
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

<style scoped>
.login-page {
  min-height: 100vh;
  min-height: 100svh;
  min-height: 100dvh;
  display: grid;
  grid-template-columns: 1.1fr 1fr;
  gap: 2rem;
  padding:
    max(1.25rem, var(--safe-top))
    max(1.25rem, var(--safe-right))
    max(1.25rem, var(--safe-bottom))
    max(1.25rem, var(--safe-left));
  align-items: stretch;
}

.login-hero {
  padding: 3rem;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.14), rgba(255, 255, 255, 0.42)),
    linear-gradient(135deg, #ff8b3d 0%, #ffd85d 50%, #71dbff 100%);
  color: #15314b;
}

.login-hero h2 {
  font-size: clamp(2.5rem, 4vw, 4.75rem);
  line-height: 0.95;
  margin: 1rem 0 0;
  max-width: 10ch;
  word-break: break-word;
  overflow-wrap: break-word;
}

@media (max-width: 1180px) {
  .login-page {
    grid-template-columns: 1fr;
    gap: 1.25rem;
  }

  .login-hero {
    min-height: 260px;
    justify-content: center;
  }
}

@media (max-width: 820px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .login-hero {
    min-height: 220px;
    padding: 2rem;
  }
}
</style>
