<template>
  <form class="inline-form" @submit.prevent="submit">
    <label>
      Kid
      <select v-model="kidId">
        <option disabled value="">Choose kid</option>
        <option v-for="kid in kids" :key="kid.id" :value="String(kid.id)">{{ kid.name }}</option>
      </select>
    </label>

    <label>
      Activity
      <select v-model="activityId">
        <option disabled value="">Choose activity</option>
        <option v-for="activity in activities" :key="activity.id" :value="String(activity.id)">
          {{ activity.icon }} {{ activity.title }}
        </option>
      </select>
    </label>

    <button class="primary-button" type="submit">Assign</button>
  </form>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  kids: {
    type: Array,
    default: () => [],
  },
  activities: {
    type: Array,
    default: () => [],
  },
  parentId: Number,
})

const emit = defineEmits(['assign'])
const kidId = ref('')
const activityId = ref('')

function submit() {
  emit('assign', {
    parentId: props.parentId,
    kidId: Number(kidId.value),
    activityId: Number(activityId.value),
  })
  kidId.value = ''
  activityId.value = ''
}
</script>
