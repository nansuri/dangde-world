<template>
  <span class="pill-badge">
    <slot />
  </span>
</template>

<style scoped>
.pill-badge {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  max-width: 100%;
  padding: 0.35rem 0.7rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.65);
  color: #234a6e;
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  word-break: break-word;
  overflow-wrap: break-word;
}
</style>
