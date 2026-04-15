<template>
  <div class="app-root">
    <RouterView />
  </div>
</template>

<script setup>
import { RouterView } from 'vue-router'
import { useSafariViewport } from './providers/useSafariViewport.js'

useSafariViewport()
</script>
