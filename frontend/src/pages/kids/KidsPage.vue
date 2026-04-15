<template>
  <AppShell
    title="Kids Zone"
    subtitle="Explore learning quests with bright colors, clear choices, and playful progress."
    :items="navItems"
    @select="selectCategory"
    @logout="logout"
  >
    <div class="page-workspace kids-workspace">
      <div class="kids-dashboard" :class="{ 'kids-dashboard-split': isTabletLandscape }">
        <section class="hero-panel kids-hero-panel">
          <div>
            <p class="eyebrow">Hi {{ session?.name }}</p>
            <h2>{{ selectedCategory.label }} adventures</h2>
            <p>Choose a game from the left, then continue the matching activity here.</p>
          </div>

          <div class="kids-hero-badges">
            <div class="kids-hero-chip">
              <span>Assignments</span>
              <strong>{{ filteredAssignments.length }}</strong>
            </div>
            <div class="kids-hero-chip">
              <span>Done</span>
              <strong>{{ completedAssignments }}</strong>
            </div>
          </div>
        </section>

        <section class="kids-detail-panel" :class="{ 'surface-card': true }">
          <template v-if="featuredAssignment">
            <div class="kids-detail-top">
              <span class="activity-icon kids-featured-icon">{{ featuredAssignment.activity.icon }}</span>
              <span class="status-chip">{{ featuredAssignment.status.replace('_', ' ') }}</span>
            </div>
            <h3>{{ featuredAssignment.activity.title }}</h3>
            <p>{{ featuredAssignment.activity.description }}</p>
            <p class="muted">{{ featuredAssignment.activity.prompt }}</p>
            <div class="progress-row">
              <div class="progress-track">
                <div class="progress-fill" :style="{ width: `${featuredAssignment.progress}%` }"></div>
              </div>
              <strong>{{ featuredAssignment.progress }}%</strong>
            </div>
            <div class="card-actions">
              <button class="primary-button" type="button" @click="playActivity(featuredAssignment)">
                Play Featured
              </button>
              <button class="ghost-button" type="button" @click="advance(featuredAssignment)">
                Quick Progress
              </button>
            </div>
          </template>

          <template v-else>
            <div class="activity-card empty-state-card kids-empty-panel">
              <div class="activity-card-top">
                <span class="activity-icon">{{ selectedCategory.icon }}</span>
                <span class="status-chip">No assignment</span>
              </div>
              <h3>No {{ selectedCategory.label.toLowerCase() }} activity yet</h3>
              <p>Ask your parent to assign a new activity in this category.</p>
            </div>
          </template>
        </section>

        <section class="card-grid kids-card-grid">
          <article
            v-for="assignment in gridAssignments"
            :key="assignment.id"
            class="activity-card kids-activity-card"
          >
            <div class="activity-card-top">
              <span class="activity-icon">{{ assignment.activity.icon }}</span>
              <span class="status-chip">{{ assignment.status.replace('_', ' ') }}</span>
            </div>
            <h3>{{ assignment.activity.title }}</h3>
            <p>{{ assignment.activity.description }}</p>
            <div class="progress-row">
              <div class="progress-track">
                <div class="progress-fill" :style="{ width: `${assignment.progress}%` }"></div>
              </div>
              <strong>{{ assignment.progress }}%</strong>
            </div>
            <div class="card-actions">
              <button class="primary-button" type="button" @click="playActivity(assignment)">
                Play
              </button>
              <button class="ghost-button" type="button" @click="setFeaturedAssignment(assignment.id)">
                Preview
              </button>
            </div>
          </article>

          <div v-if="featuredAssignment && gridAssignments.length === 0" class="activity-card empty-state-card kids-empty-panel">
             <div class="activity-card-top">
              <span class="activity-icon">✨</span>
            </div>
            <h3>More adventures await!</h3>
            <p>You are currently focusing on the featured activity above.</p>
          </div>
        </section>
      </div>
    </div>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '../../shared/ui/AppShell.vue'
import { readSession, clearSession } from '../../features/auth/session.js'
import { listAssignments, updateAssignment } from '../../entities/assignment/api.js'
import { useViewportProfile } from '../../shared/lib/useViewportProfile.js'

const router = useRouter()
const session = ref(readSession())
const assignments = ref([])
const selectedCategoryKey = ref('math')
const selectedAssignmentId = ref(null)
const { isTabletLandscape } = useViewportProfile()

const categoryItems = [
  { key: 'math', label: 'Math', icon: '➕', matches: ['math'] },
  { key: 'alphabets', label: 'Alphabets', icon: '🔤', matches: ['alphabets'] },
  { key: 'arabic', label: 'Arabic', icon: '📖', matches: ['arabic'] },
  { key: 'clock', label: 'Clock', icon: '🕒', matches: ['clock'] },
]

const navItems = computed(() =>
  categoryItems.map((item) => ({
    ...item,
    active: item.key === selectedCategoryKey.value,
  })),
)

const selectedCategory = computed(
  () => categoryItems.find((item) => item.key === selectedCategoryKey.value) || categoryItems[0],
)

const filteredAssignments = computed(() =>
  assignments.value.filter((assignment) => {
    const slug = assignment.activity.category?.slug || ''
    return selectedCategory.value.matches.some((match) => slug.includes(match))
  }),
)

const featuredAssignment = computed(
  () =>
    filteredAssignments.value.find((assignment) => assignment.id === selectedAssignmentId.value) ||
    filteredAssignments.value[0] ||
    null,
)

const gridAssignments = computed(() => {
  if (!featuredAssignment.value) return filteredAssignments.value
  return filteredAssignments.value.filter((a) => a.id !== featuredAssignment.value.id)
})

const completedAssignments = computed(
  () => filteredAssignments.value.filter((assignment) => assignment.status === 'completed').length,
)

onMounted(loadAssignments)

watch(filteredAssignments, (items) => {
  if (!items.length) {
    selectedAssignmentId.value = null
    return
  }

  const exists = items.some((assignment) => assignment.id === selectedAssignmentId.value)
  if (!exists) {
    selectedAssignmentId.value = items[0].id
  }
})

async function loadAssignments() {
  const response = await listAssignments({ kidId: session.value.id })
  assignments.value = response.items
  if (!selectedAssignmentId.value && response.items.length) {
    selectedAssignmentId.value = response.items[0].id
  }
}

async function advance(item) {
  const progress = Math.min(item.progress + 20, 100)
  const status = progress === 100 ? 'completed' : 'in_progress'
  await updateAssignment(item.id, { progress, status })
  await loadAssignments()
}

function playActivity(item) {
  router.push(`/kids/activity/${item.id}`)
}

function selectCategory(item) {
  selectedCategoryKey.value = item.key
  selectedAssignmentId.value = filteredAssignments.value[0]?.id || null
}

function setFeaturedAssignment(assignmentId) {
  selectedAssignmentId.value = assignmentId
}

function logout() {
  clearSession()
  router.push('/login')
}
</script>
