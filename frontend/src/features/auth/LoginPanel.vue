<template>
  <section class="login-card">
    <div class="login-copy">
      <p class="eyebrow">Play. Learn. Grow.</p>
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

      <button class="primary-button" :disabled="loading" type="submit">
        {{ loading ? 'Entering...' : 'Enter App' }}
      </button>
      <p v-if="error" class="error-text">{{ error }}</p>
    </form>
  </section>
</template>

<script setup>
import { ref } from 'vue'

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
