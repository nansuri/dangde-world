<template>
  <button
    class="template-button"
    type="button"
    :title="template.description"
    @click="emit('select', template)"
  >
    <span class="template-icon">{{ template.icon }}</span>
    <span class="template-name">{{ template.name }}</span>
  </button>
</template>

<script setup>
defineProps({
  template: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['select'])
</script>

<style scoped>
.template-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: white;
  border: 2px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: all 200ms;
  text-align: center;
  color: var(--ink);
}

.template-button:hover {
  border-color: var(--primary);
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.template-button:focus-visible {
  outline: 2px solid var(--secondary);
  outline-offset: 2px;
}

.template-icon {
  font-size: 2rem;
  display: block;
}

.template-name {
  font-weight: 600;
  font-size: 0.85rem;
}
</style>
