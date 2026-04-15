<template>
  <button class="action-button" :class="variant" :type="type">
    <slot />
  </button>
</template>

<script setup>
defineProps({
  variant: {
    type: String,
    default: 'primary', // 'primary' | 'ghost'
  },
  type: {
    type: String,
    default: 'button',
  },
})
</script>

<style scoped>
.action-button {
  border: 0;
  border-radius: 18px;
  min-height: 58px;
  padding: 1rem 1.2rem;
  font-weight: 800;
  transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
  cursor: pointer;
  font: inherit;
  font-size: 16px;
  touch-action: manipulation;
}

.action-button.primary {
  background: linear-gradient(135deg, var(--primary), #ffa62e);
  color: white;
  box-shadow: 0 14px 30px rgba(255, 122, 24, 0.28);
}

.action-button.ghost {
  background: rgba(255, 255, 255, 0.7);
  color: var(--ink);
}

.action-button:hover,
.action-button:focus-visible {
  transform: translateY(-2px);
  outline: none;
}

.action-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}
</style>
