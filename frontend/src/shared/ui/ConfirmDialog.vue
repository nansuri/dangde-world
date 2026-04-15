<template>
  <AppModal :open="open" :title="title" :subtitle="subtitle" max-width="520px" @close="$emit('cancel')">
    <p class="confirm-message">{{ message }}</p>
    <div class="confirm-actions">
      <ActionButton variant="ghost" type="button" @click="$emit('cancel')">
        {{ cancelLabel }}
      </ActionButton>
      <ActionButton type="button" @click="$emit('confirm')">
        {{ confirmLabel }}
      </ActionButton>
    </div>
  </AppModal>
</template>

<script setup>
import AppModal from './AppModal.vue'
import ActionButton from './ActionButton.vue'

defineProps({
  open: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: 'Please confirm',
  },
  subtitle: {
    type: String,
    default: 'Confirmation',
  },
  message: {
    type: String,
    default: 'Are you sure you want to continue?',
  },
  confirmLabel: {
    type: String,
    default: 'Confirm',
  },
  cancelLabel: {
    type: String,
    default: 'Cancel',
  },
})

defineEmits(['confirm', 'cancel'])
</script>

<style scoped>
.confirm-message {
  margin: 0 0 1rem;
  line-height: 1.45;
}

.confirm-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
}
</style>
