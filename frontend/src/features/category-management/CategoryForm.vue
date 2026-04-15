<template>
  <form class="category-form" @submit.prevent="submit">
    <div class="form-group">
      <label>
        Category Name
        <input v-model="form.name" type="text" required placeholder="e.g., Mathematics" />
      </label>
    </div>

    <div class="form-group">
      <label>
        Slug
        <input v-model="form.slug" type="text" required placeholder="e.g., math" />
      </label>
      <p class="form-hint">URL-friendly name (lowercase, no spaces)</p>
    </div>

    <div class="form-row">
      <label>
        Type
        <select v-model="form.type" required>
          <option value="learning">Learning</option>
          <option value="playing">Playing</option>
        </select>
      </label>

      <label>
        Parent Category
        <select v-model="form.parentId">
          <option :value="null">None (Root)</option>
          <option v-for="cat in availableParents" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
      </label>
    </div>

    <div class="form-actions">
      <button type="button" class="ghost-button" @click="cancel">Cancel</button>
      <button type="submit" class="primary-button">{{ submitLabel }}</button>
    </div>
  </form>
</template>

<script setup>
import { computed, reactive, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: null,
  },
  categories: {
    type: Array,
    default: () => [],
  },
  submitLabel: {
    type: String,
    default: 'Save',
  },
})

const emit = defineEmits(['submit', 'cancel'])

const form = reactive({
  name: '',
  slug: '',
  type: 'learning',
  parentId: null,
})

const availableParents = computed(() => {
  // Only show root categories (parentId is null)
  return props.categories.filter((cat) => !cat.parentId && cat.id !== props.modelValue?.id)
})

watch(
  () => props.modelValue,
  (value) => {
    if (value) {
      Object.assign(form, value)
    } else {
      form.name = ''
      form.slug = ''
      form.type = 'learning'
      form.parentId = null
    }
  },
  { immediate: true },
)

function submit() {
  emit('submit', {
    ...form,
    parentId: form.parentId ? Number(form.parentId) : null,
  })
}

function cancel() {
  emit('cancel')
}
</script>

<style scoped>
.category-form {
  display: grid;
  gap: 1.5rem;
  padding: 1rem;
  background: var(--surface);
  border-radius: 8px;
}

.form-group,
.form-row {
  display: grid;
  gap: 0.5rem;
}

.form-row {
  grid-template-columns: 1fr 1fr;
}

label {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  font-weight: 600;
  color: var(--ink);
}

input,
select {
  padding: 0.75rem;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: white;
  font-size: 1rem;
}

input:focus,
select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(255, 122, 24, 0.1);
}

.form-hint {
  margin: 0;
  font-size: 0.85rem;
  color: var(--muted);
  font-weight: 400;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-weight: 600;
  cursor: pointer;
  transition: all 200ms;
}

.primary-button {
  background: var(--primary);
  color: white;
}

.primary-button:hover {
  background: var(--primary-dark);
}

.ghost-button {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--ink);
}

.ghost-button:hover {
  background: rgba(0, 0, 0, 0.04);
}
</style>
