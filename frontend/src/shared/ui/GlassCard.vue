<template>
  <component :is="tag" class="glass-card">
    <slot />
  </component>
</template>

<script setup>
defineProps({
  tag: {
    type: String,
    default: 'div',
  },
})
</script>

<style scoped>
.glass-card {
  background: var(--surface);
  backdrop-filter: blur(18px);
  border: 1px solid var(--border);
  border-radius: 28px;
  box-shadow: var(--shadow);
  min-width: 0;
}

@media (max-width: 820px) {
  .glass-card {
    border-radius: 22px;
  }
}
</style>
