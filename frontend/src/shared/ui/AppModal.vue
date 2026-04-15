<template>
  <Teleport to="body">
    <div v-if="open" class="app-modal-overlay" @click.self="emit('close')">
      <section class="app-modal-panel" :style="panelStyle" role="dialog" aria-modal="true">
        <header class="app-modal-header">
          <div>
            <p v-if="subtitle" class="app-modal-subtitle">{{ subtitle }}</p>
            <h3 class="app-modal-title">{{ title }}</h3>
          </div>
          <button class="app-modal-close" type="button" @click="emit('close')" aria-label="Close dialog">
            ✕
          </button>
        </header>
        <div class="app-modal-body">
          <slot />
        </div>
        <footer v-if="$slots.footer" class="app-modal-footer">
          <slot name="footer" />
        </footer>
      </section>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  open: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: '',
  },
  subtitle: {
    type: String,
    default: '',
  },
  maxWidth: {
    type: String,
    default: '680px',
  },
})

const emit = defineEmits(['close'])

const panelStyle = computed(() => ({
  maxWidth: props.maxWidth,
}))

function onKeydown(event) {
  if (event.key === 'Escape' && props.open) {
    emit('close')
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
.app-modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(16, 28, 44, 0.46);
  backdrop-filter: blur(6px);
  display: grid;
  place-items: center;
  padding: 1rem;
}

.app-modal-panel {
  width: min(100%, var(--modal-max, 680px));
  max-height: min(86vh, 920px);
  overflow: auto;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(24, 78, 122, 0.12);
  box-shadow: 0 30px 60px rgba(12, 28, 48, 0.24);
}

.app-modal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1rem 0.5rem;
}

.app-modal-subtitle {
  margin: 0;
  font-size: 0.8rem;
  color: var(--muted);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.app-modal-title {
  margin: 0.25rem 0 0;
}

.app-modal-close {
  border: 1px solid rgba(24, 78, 122, 0.16);
  background: rgba(255, 255, 255, 0.82);
  border-radius: 12px;
  min-width: 38px;
  min-height: 38px;
  cursor: pointer;
  font: inherit;
  font-weight: 800;
}

.app-modal-body {
  padding: 0.5rem 1rem 1rem;
}

.app-modal-footer {
  padding: 0 1rem 1rem;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}
</style>
