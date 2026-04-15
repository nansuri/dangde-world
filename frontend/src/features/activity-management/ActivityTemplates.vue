<template>
  <div class="template-selector">
    <TemplateButton
      v-for="template in templates"
      :key="template.id"
      :template="template"
      @select="handleSelect"
    />
  </div>
</template>

<script setup>
import { ACTIVITY_TEMPLATES } from './data/activityTemplates'
import TemplateButton from './components/TemplateButton.vue'

defineProps({
  onSelect: Function,
})

const emit = defineEmits(['select'])
const templates = ACTIVITY_TEMPLATES

function handleSelect(template) {
  emit('select', template)
}
</script>

<style scoped>
.template-selector {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1rem;
  padding: 1rem;
}
</style>
