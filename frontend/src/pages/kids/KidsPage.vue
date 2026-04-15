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
        <GlassCard tag="section" class="hero-panel kids-hero-panel">
          <div>
            <PillBadge class="eyebrow">Hi {{ session?.name }}</PillBadge>
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
        </GlassCard>

        <GlassCard tag="section" class="kids-detail-panel">
          <template v-if="featuredAssignment">
            <div class="kids-detail-top">
              <span class="activity-icon kids-featured-icon">{{ featuredAssignment.activity.icon }}</span>
              <PillBadge class="status-chip">{{ featuredAssignment.status.replace('_', ' ') }}</PillBadge>
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
              <ActionButton @click="playActivity(featuredAssignment)">
                Play Featured
              </ActionButton>
              <ActionButton variant="ghost" @click="advance(featuredAssignment)">
                Quick Progress
              </ActionButton>
            </div>
          </template>

          <template v-else>
            <div class="activity-card empty-state-card kids-empty-panel">
              <div class="activity-card-top">
                <span class="activity-icon">{{ selectedCategory.icon }}</span>
                <PillBadge class="status-chip">No assignment</PillBadge>
              </div>
              <h3>No {{ selectedCategory.label.toLowerCase() }} activity yet</h3>
              <p>Ask your parent to assign a new activity in this category.</p>
            </div>
          </template>
        </GlassCard>

        <section class="card-grid kids-card-grid">
          <GlassCard
            v-for="assignment in gridAssignments"
            :key="assignment.id"
            tag="article"
            class="activity-card kids-activity-card"
          >
            <div class="activity-card-top">
              <span class="activity-icon">{{ assignment.activity.icon }}</span>
              <PillBadge class="status-chip">{{ assignment.status.replace('_', ' ') }}</PillBadge>
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
              <ActionButton @click="playActivity(assignment)">
                Play
              </ActionButton>
              <ActionButton variant="ghost" @click="setFeaturedAssignment(assignment.id)">
                Preview
              </ActionButton>
            </div>
          </GlassCard>

          <GlassCard v-if="featuredAssignment && gridAssignments.length === 0" class="activity-card empty-state-card kids-empty-panel">
             <div class="activity-card-top">
              <span class="activity-icon">✨</span>
            </div>
            <h3>More adventures await!</h3>
            <p>You are currently focusing on the featured activity above.</p>
          </GlassCard>
        </section>
      </div>
    </div>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '../../shared/ui/AppShell.vue'
import GlassCard from '../../shared/ui/GlassCard.vue'
import PillBadge from '../../shared/ui/PillBadge.vue'
import ActionButton from '../../shared/ui/ActionButton.vue'
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

<style scoped>
.page-workspace {
  height: 100%;
  min-height: 0;
  display: grid;
}

.kids-workspace {
  grid-template-rows: minmax(0, 1fr);
  overflow: hidden;
}

.kids-dashboard {
  display: grid;
  grid-template-rows: auto;
  gap: 1rem;
  min-height: 0;
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.kids-dashboard > section {
  min-height: 0;
  overflow: visible;
}

.kids-dashboard-split {
  grid-template-columns: minmax(320px, 0.95fr) minmax(0, 1.45fr);
  grid-template-rows: auto minmax(0, 1fr);
  grid-template-areas:
    "hero grid"
    "detail grid";
  align-items: stretch;
}

.kids-dashboard-split .kids-hero-panel {
  grid-area: hero;
  min-height: 220px;
}

.kids-dashboard-split .kids-detail-panel {
  grid-area: detail;
  align-self: stretch;
  min-height: 320px;
}

.kids-dashboard-split .kids-card-grid {
  grid-area: grid;
  height: 100%;
}

.hero-panel {
  padding: 1.5rem;
}

.kids-hero-panel {
  display: grid;
  gap: 1.25rem;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.kids-hero-panel > * {
  min-width: 0;
}

.kids-hero-badges {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.kids-hero-badges > * {
  min-width: 0;
}

.kids-hero-chip {
  background: rgba(255, 255, 255, 0.68);
  border: 1px solid var(--border);
  border-radius: 24px;
  padding: 1rem 1.1rem;
}

.kids-hero-chip span {
  display: block;
  color: var(--muted);
  word-break: break-word;
}

.kids-hero-chip strong {
  display: block;
  margin-top: 0.35rem;
  font-size: 2rem;
  word-break: break-word;
  overflow-wrap: break-word;
}

.kids-detail-panel {
  padding: 1.5rem;
  display: grid;
  gap: 1rem;
  min-width: 0;
}

.kids-detail-panel h3,
.kids-detail-panel p {
  min-width: 0;
  word-break: break-word;
  overflow-wrap: break-word;
}

.kids-detail-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.kids-featured-icon {
  font-size: 3rem;
}

.card-grid {
  display: grid;
  gap: 1rem;
}

.kids-card-grid {
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
}

.activity-card {
  padding: 1.25rem;
}

.kids-activity-card {
  min-height: 280px;
  display: grid;
  align-content: start;
  gap: 0.8rem;
  overflow: hidden;
}

.kids-activity-card h3,
.kids-activity-card p {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  word-break: break-word;
  overflow-wrap: break-word;
}

.kids-empty-panel {
  min-height: 220px;
}

.activity-card-top,
.progress-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  min-width: 0;
}

.activity-icon {
  font-size: 2rem;
}

.progress-track {
  height: 12px;
  flex: 1;
  border-radius: 999px;
  background: rgba(77, 142, 192, 0.12);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #00bcd4, #ffd84d);
  border-radius: inherit;
}

.card-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.empty-state-card {
  border-style: dashed;
  background: rgba(255, 255, 255, 0.55);
}

@media (max-width: 820px) {
  .hero-panel {
    border-radius: 22px;
  }

  .kids-hero-panel {
    grid-template-columns: 1fr;
  }

  .kids-hero-badges {
    grid-template-columns: 1fr;
  }

  .activity-card-top,
  .progress-row {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
