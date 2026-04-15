<template>
  <GlassCard tag="section" class="login-card">
    <div class="login-copy">
      <PillBadge class="eyebrow">Play. Learn. Grow.</PillBadge>
      <h1>Welcome to DangDe! World</h1>
      <p>
        A colorful learning playground for kids, parents, and admins. Choose a demo user
        to enter the correct experience.
      </p>
    </div>

    <form class="login-form" @submit.prevent="submit">
      <label>
        Demo User
        <select v-model="selectedUserId">
          <option disabled value="">Select a role</option>
          <option v-for="user in users" :key="user.id" :value="String(user.id)">
            {{ user.avatar }} {{ user.name }} · {{ user.role }}
          </option>
        </select>
      </label>

      <ActionButton :disabled="loading" type="submit">
        {{ loading ? 'Entering...' : 'Enter App' }}
      </ActionButton>
      <p v-if="error" class="error-text">{{ error }}</p>
    </form>
  </GlassCard>
</template>

<script setup>
import { ref } from 'vue'
import GlassCard from '../../shared/ui/GlassCard.vue'
import PillBadge from '../../shared/ui/PillBadge.vue'
import ActionButton from '../../shared/ui/ActionButton.vue'

const props = defineProps({
  users: {
    type: Array,
    default: () => [],
  },
  loading: Boolean,
  error: String,
})

const emit = defineEmits(['submit'])
const selectedUserId = ref('')

function submit() {
  emit('submit', Number(selectedUserId.value))
}
</script>

<style scoped>
.login-card {
  padding: 2rem;
  display: grid;
  gap: 2rem;
  align-content: center;
}

.login-copy h1 {
  margin: 0.25rem 0 0.5rem;
  font-size: clamp(1.8rem, 2.5vw, 2.8rem);
  word-break: break-word;
  overflow-wrap: break-word;
}

.login-copy p {
  color: var(--muted);
  word-break: break-word;
  overflow-wrap: break-word;
}

.login-form {
  display: grid;
  gap: 1rem;
}

label {
  display: grid;
  gap: 0.45rem;
  font-weight: 700;
}

.error-text {
  color: #bf1650;
  margin: 0;
}

@media (max-width: 820px) {
  .login-card {
    padding: 2rem;
  }
}
</style>
