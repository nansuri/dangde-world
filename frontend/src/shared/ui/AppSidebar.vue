<template>
  <aside class="app-sidebar">
    <div class="sidebar-head">
      <PillBadge class="eyebrow">DangDe! World</PillBadge>
      <div class="sidebar-title-wrap">
        <h1>{{ title }}</h1>
        <ActionButton
          class="sidebar-logout-mobile"
          variant="ghost"
          @click="$emit('logout')"
        >
          Logout
        </ActionButton>
      </div>
      <p class="sidebar-copy">{{ subtitle }}</p>
    </div>

    <nav class="nav-list">
      <button
        v-for="item in items"
        :key="item.label"
        class="nav-pill"
        type="button"
        :class="{ active: item.active }"
        @click="$emit('select', item)"
      >
        <span>{{ item.icon }}</span>
        <span>{{ item.label }}</span>
      </button>
    </nav>

    <ActionButton
      class="sidebar-logout-desktop"
      variant="ghost"
      @click="$emit('logout')"
    >
      Logout
    </ActionButton>
  </aside>
</template>

<script setup>
import PillBadge from './PillBadge.vue'
import ActionButton from './ActionButton.vue'

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
.app-sidebar {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.74), rgba(255, 244, 214, 0.82));
  border: 1px solid var(--border);
  border-radius: 28px;
  box-shadow: var(--shadow);
  padding: 1.5rem;
  position: sticky;
  top: max(1rem, var(--safe-top));
  align-self: start;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 1.5rem;
  height: calc(100vh - max(2rem, var(--safe-top) + var(--safe-bottom)));
  height: calc(100svh - max(2rem, var(--safe-top) + var(--safe-bottom)));
  height: calc(100dvh - max(2rem, var(--safe-top) + var(--safe-bottom)));
  overflow: hidden;
}

.sidebar-head {
  display: grid;
  gap: 0.65rem;
  min-width: 0;
}

.sidebar-head > * {
  min-width: 0;
}

.sidebar-title-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.sidebar-title-wrap h1 {
  margin: 0;
  min-width: 0;
  word-break: break-word;
  overflow-wrap: break-word;
}

.sidebar-logout-mobile {
  display: none;
}

.sidebar-copy {
  color: var(--muted);
  word-break: break-word;
  overflow-wrap: break-word;
}

.nav-list {
  display: grid;
  gap: 1rem;
}

.nav-pill {
  border: 0;
  border-radius: 18px;
  min-height: 58px;
  padding: 1rem 1.2rem;
  font-weight: 800;
  transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
  background: rgba(255, 255, 255, 0.58);
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 0.7rem;
  width: 100%;
  min-width: 0;
  cursor: pointer;
  font: inherit;
}

.nav-pill span:last-child {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nav-pill.active,
.nav-pill:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.85);
}

@media (max-width: 1180px) {
  .app-sidebar {
    min-height: auto;
    height: auto;
    padding: 1rem;
    gap: 1rem;
    position: sticky;
    top: max(0.75rem, var(--safe-top));
    z-index: 20;
    overflow: visible;
  }

  .sidebar-title-wrap {
    align-items: flex-start;
  }

  .sidebar-logout-mobile {
    display: inline-flex;
    flex-shrink: 0;
    min-height: 44px;
    padding: 0.5rem 1rem;
  }

  .sidebar-logout-desktop {
    display: none;
  }

  .nav-list {
    grid-auto-flow: column;
    grid-auto-columns: minmax(140px, 1fr);
    overflow-x: auto;
    overflow-y: hidden;
    padding-bottom: 0.25rem;
    scrollbar-width: none;
    overscroll-behavior-x: contain;
    -webkit-overflow-scrolling: touch;
  }

  .nav-list::-webkit-scrollbar {
    display: none;
  }

  .nav-pill {
    min-width: 140px;
    justify-content: center;
  }
}

@media (max-width: 820px) {
  .sidebar-title-wrap {
    align-items: flex-start;
    flex-direction: column;
  }

  .nav-list {
    grid-auto-columns: minmax(160px, 1fr);
  }
}
</style>
