<template>
  <form class="activity-ide" @submit.prevent="submit">
    <section class="activity-ide-sidebar">
      <div class="activity-ide-toolbar">
        <div>
          <PillBadge class="eyebrow">Workspace</PillBadge>
          <h3>Activity Package</h3>
        </div>
        <ActionButton type="submit">{{ submitLabel }}</ActionButton>
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

      <div class="activity-ide-scrollable">
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
import ActionButton from '../../shared/ui/ActionButton.vue'
import PillBadge from '../../shared/ui/PillBadge.vue'
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

<style scoped>
.activity-ide {
  display: grid;
  grid-template-columns: minmax(280px, 0.9fr) minmax(0, 1.4fr);
  gap: 1rem;
  align-items: start;
  min-height: 0;
  height: 100%;
}

.activity-ide-sidebar,
.activity-editor-shell {
  border: 1px solid rgba(24, 78, 122, 0.1);
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.62);
}

.activity-ide-sidebar {
  padding: 1rem;
  display: grid;
  grid-template-rows: auto auto minmax(0, 1fr);
  gap: 1rem;
  min-height: 0;
  height: 100%;
  overflow: hidden;
}

.activity-ide-sidebar > :nth-child(n+3) {
  min-height: 0;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.activity-ide-scrollable {
  display: grid;
  gap: 1rem;
  min-height: 0;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.activity-ide-toolbar {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: start;
  gap: 0.9rem;
}

.activity-ide-toolbar > div {
  min-width: 0;
}

.activity-ide-toolbar h3 {
  margin: 0.3rem 0 0;
  font-size: 1.4rem;
  line-height: 1.15;
}

.activity-ide-tree {
  display: grid;
  gap: 0.5rem;
  padding: 0.25rem;
  border-radius: 20px;
  background: rgba(23, 50, 77, 0.06);
}

.activity-file-pill {
  border: 0;
  border-radius: 16px;
  min-height: 50px;
  padding: 0.8rem 0.95rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 0.7rem;
  background: transparent;
  color: var(--ink);
  cursor: pointer;
  font: inherit;
}

.activity-file-pill.active {
  background: rgba(255, 255, 255, 0.84);
  box-shadow: 0 10px 22px rgba(20, 54, 84, 0.08);
}

.activity-meta-panel {
  display: grid;
  gap: 1rem;
}

label {
  display: grid;
  gap: 0.45rem;
  font-weight: 700;
}

.activity-meta-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.activity-editor-shell {
  overflow: hidden;
  min-height: 0;
  display: grid;
  grid-template-rows: auto minmax(0, 1fr);
}

.activity-editor-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.85rem 1rem;
  border-bottom: 1px solid rgba(24, 78, 122, 0.08);
  background: linear-gradient(180deg, rgba(18, 31, 47, 0.96), rgba(27, 44, 65, 0.95));
  color: #eef6ff;
}

.activity-tab-strip {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.activity-editor-tab {
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px 14px 0 0;
  min-height: 44px;
  padding: 0.75rem 0.95rem;
  background: rgba(255, 255, 255, 0.06);
  color: inherit;
  display: flex;
  align-items: center;
  gap: 0.55rem;
  font-weight: 700;
  cursor: pointer;
  font: inherit;
}

.activity-editor-tab.active {
  background: rgba(255, 255, 255, 0.16);
  border-color: rgba(255, 255, 255, 0.18);
}

.activity-editor-stats {
  display: flex;
  gap: 0.9rem;
  color: rgba(238, 246, 255, 0.78);
  font-size: 0.9rem;
}

.activity-editor-body {
  display: grid;
  grid-template-columns: 58px minmax(0, 1fr);
  min-height: 0;
  height: 100%;
  background: linear-gradient(180deg, #0f1b2c, #152538);
}

.activity-editor-gutter {
  padding: 1rem 0.6rem;
  display: grid;
  align-content: start;
  gap: 0.15rem;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(214, 226, 242, 0.45);
  font-family: "SFMono-Regular", "Menlo", "Monaco", monospace;
  font-size: 0.9rem;
  text-align: right;
}

.activity-editor-pane {
  display: block;
}

.activity-code-editor {
  width: 100%;
  height: 100%;
  min-height: 0;
  border: 0;
  border-radius: 0;
  padding: 1rem 1.1rem;
  background: transparent;
  color: #edf5ff;
  font-family: "SFMono-Regular", "Menlo", "Monaco", monospace;
  line-height: 1.6;
  resize: none;
}

.activity-code-editor:focus {
  outline: none;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

@media (max-width: 1180px) {
  .activity-ide {
    grid-template-columns: 1fr;
    height: auto;
  }

  .activity-editor-topbar {
    flex-direction: column;
    align-items: flex-start;
  }

  .activity-ide-toolbar {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 820px) {
  .activity-meta-grid {
    grid-template-columns: 1fr;
  }

  .activity-editor-body {
    grid-template-columns: 44px minmax(0, 1fr);
    min-height: 520px;
  }

  .activity-code-editor {
    min-height: 520px;
  }
}
</style>
