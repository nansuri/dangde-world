<template>
  <article class="metric-card">
    <p>{{ label }}</p>
    <strong>{{ value }}</strong>
  </article>
</template>

<script setup>
defineProps({
  label: String,
  value: [String, Number],
})
</script>
