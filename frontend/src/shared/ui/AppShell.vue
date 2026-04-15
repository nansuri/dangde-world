<template>
  <div class="app-shell">
    <aside class="app-sidebar">
      <div class="sidebar-head">
        <p class="eyebrow">DangDe! World</p>
        <div class="sidebar-title-wrap">
          <h1>{{ title }}</h1>
          <button class="ghost-button sidebar-logout-mobile" type="button" @click="$emit('logout')">
            Logout
          </button>
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

      <button class="ghost-button sidebar-logout-desktop" type="button" @click="$emit('logout')">
        Logout
      </button>
    </aside>

    <main class="app-content">
      <div class="app-content-inner">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup>
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
