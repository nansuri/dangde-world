<template>
  <div class="app-shell">
    <AppSidebar
      :title="title"
      :subtitle="subtitle"
      :items="items"
      @logout="$emit('logout')"
      @select="(item) => $emit('select', item)"
    />

    <main class="app-content">
      <div class="app-content-inner">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup>
import AppSidebar from './AppSidebar.vue'

defineProps({
  title: String,
  subtitle: String,
  items: {
    type: Array,
    default: () => [],
  },
})

defineEmits(['logout', 'select'])
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  min-height: 100svh;
  min-height: 100dvh;
  height: 100vh;
  height: 100svh;
  height: 100dvh;
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 1.5rem;
  padding:
    max(1rem, var(--safe-top))
    max(1rem, var(--safe-right))
    max(1rem, var(--safe-bottom))
    max(1rem, var(--safe-left));
  align-items: start;
}

.app-content {
  align-content: start;
  min-width: 0;
  min-height: 0;
  height: calc(100vh - max(2rem, var(--safe-top) + var(--safe-bottom)));
  height: calc(100svh - max(2rem, var(--safe-top) + var(--safe-bottom)));
  height: calc(100dvh - max(2rem, var(--safe-top) + var(--safe-bottom)));
  overflow-y: auto;
  overflow-x: hidden;
  -webkit-overflow-scrolling: touch;
  overscroll-behavior: contain;
  padding-right: 0.125rem;
}

.app-content-inner {
  height: 100%;
  min-height: 0;
}

@media (max-width: 1180px) {
  .app-shell {
    grid-template-columns: 1fr;
    gap: 1rem;
    height: auto;
    min-height: 100vh;
    min-height: 100svh;
    min-height: 100dvh;
  }

  .app-content {
    height: auto;
    min-height: 0;
    overflow: visible;
    padding-right: 0;
  }
}
</style>
