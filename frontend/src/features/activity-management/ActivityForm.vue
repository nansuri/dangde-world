<template>
  <form class="activity-ide" @submit.prevent="submit">
    <section class="activity-ide-sidebar">
      <div class="activity-ide-toolbar">
        <div>
          <p class="eyebrow">Workspace</p>
          <h3>Activity Package</h3>
        </div>
        <button class="primary-button" type="submit">{{ submitLabel }}</button>
      </div>

      <div class="activity-ide-tree">
        <button
          v-for="file in editorTabs"
          :key="file.key"
          class="activity-file-pill"
          type="button"
          :class="{ active: activeTab === file.key }"
          @click="activeTab = file.key"
        >
          <span>{{ file.icon }}</span>
          <span>{{ file.label }}</span>
        </button>
      </div>

      <APIReference />

      <div class="activity-meta-panel">
        <label>
          Title
          <input v-model="form.title" required />
        </label>
        <label>
          Description
          <textarea v-model="form.description" rows="3" required />
        </label>
        <label>
          Prompt
          <textarea v-model="form.prompt" rows="3" required />
        </label>
        <div class="activity-meta-grid">
          <label>
            Language
            <input v-model="form.language" required />
          </label>
          <label>
            Difficulty
            <select v-model="form.difficulty" required>
              <option value="easy">easy</option>
              <option value="medium">medium</option>
              <option value="hard">hard</option>
            </select>
          </label>
          <label>
            Age Group
            <input v-model="form.ageGroup" required />
          </label>
          <label>
            Category
            <select v-model="form.categoryId" required>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </label>
          <label>
            Icon
            <input v-model="form.icon" placeholder="🎨" />
          </label>
        </div>
      </div>
    </section>

    <section class="activity-editor-shell">
      <div class="activity-editor-topbar">
        <div class="activity-tab-strip">
          <button
            v-for="file in editorTabs"
            :key="file.key"
            class="activity-editor-tab"
            type="button"
            :class="{ active: activeTab === file.key }"
            @click="activeTab = file.key"
          >
            <span>{{ file.icon }}</span>
            <span>{{ file.fileName }}</span>
          </button>
        </div>
        <div class="activity-editor-stats">
          <span>{{ activeFile.label }}</span>
          <span>{{ currentCodeLength }} chars</span>
        </div>
      </div>

      <div class="activity-editor-body">
        <div class="activity-editor-gutter">
          <span v-for="line in visibleLineNumbers" :key="line">{{ line }}</span>
        </div>
        <label class="activity-editor-pane">
          <span class="sr-only">{{ activeFile.label }}</span>
          <textarea
            v-if="activeTab === 'html'"
            v-model="form.htmlCode"
            class="activity-code-editor"
            rows="18"
            placeholder="<div class='game-shell'>Hello activity</div>"
            spellcheck="false"
          />
          <textarea
            v-else-if="activeTab === 'css'"
            v-model="form.cssCode"
            class="activity-code-editor"
            rows="18"
            placeholder=".game-shell { display: grid; }"
            spellcheck="false"
          />
          <textarea
            v-else
            v-model="form.jsCode"
            class="activity-code-editor"
            rows="18"
            placeholder="window.DangdeAPI.save({ key: 'score', value: { score: 10 } })"
            spellcheck="false"
          />
        </label>
      </div>
    </section>
  </form>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import APIReference from './APIReference.vue'

const props = defineProps({
  categories: {
    type: Array,
    default: () => [],
  },
  modelValue: {
    type: Object,
    default: null,
  },
  submitLabel: {
    type: String,
    default: 'Save Activity',
  },
})

const emit = defineEmits(['submit'])
const activeTab = ref('html')
const editorTabs = [
  { key: 'html', label: 'Markup', fileName: 'index.html', icon: '🌐' },
  { key: 'css', label: 'Styles', fileName: 'styles.css', icon: '🎨' },
  { key: 'js', label: 'Logic', fileName: 'app.js', icon: '⚙️' },
]

const createDefaultForm = () => ({
  title: '',
  description: '',
  language: 'English',
  difficulty: 'easy',
  ageGroup: '4-6',
  categoryId: props.categories[0]?.id || 1,
  prompt: '',
  icon: '✨',
  htmlCode: '',
  cssCode: '',
  jsCode: '',
})

const form = reactive(createDefaultForm())

const activeFile = computed(() => editorTabs.find((item) => item.key === activeTab.value) || editorTabs[0])
const currentCode = computed(() => {
  if (activeTab.value === 'html') {
    return form.htmlCode
  }
  if (activeTab.value === 'css') {
    return form.cssCode
  }
  return form.jsCode
})
const currentCodeLength = computed(() => currentCode.value.length)
const visibleLineNumbers = computed(() => {
  const lines = Math.max(18, currentCode.value.split('\n').length)
  return Array.from({ length: lines }, (_, index) => index + 1)
})

watch(
  () => props.modelValue,
  (value) => {
    Object.assign(form, value || createDefaultForm())
  },
  { immediate: true },
)

watch(
  () => props.categories,
  (categories) => {
    if (!form.categoryId && categories.length) {
      form.categoryId = categories[0].id
    }
  },
  { immediate: true },
)

function submit() {
  emit('submit', { ...form, categoryId: Number(form.categoryId) })
  Object.assign(form, createDefaultForm())
  activeTab.value = 'html'
}
</script>
