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
          <div class="kids-hero-content">
            <PillBadge class="eyebrow">Hi {{ session?.name }}</PillBadge>
            <h2>{{ selectedCategory.label }} adventures</h2>
            <p class="kids-hero-desc">Ready to play? Choose a game to start your learning quest!</p>
          </div>

          <div class="kids-hero-badges">
            <div class="kids-hero-chip">
              <span class="chip-label">Assignments</span>
              <strong class="chip-value">{{ filteredAssignments.length }}</strong>
            </div>
            <div class="kids-hero-chip">
              <span class="chip-label">Done</span>
              <strong class="chip-value">{{ completedAssignments }}</strong>
            </div>
          </div>
        </GlassCard>

        <GlassCard tag="section" class="kids-detail-panel">
          <template v-if="featuredAssignment">
            <div class="kids-detail-top">
              <span class="activity-icon kids-featured-icon">{{ featuredAssignment.activity.icon }}</span>
              <div class="kids-detail-meta">
                <PillBadge class="status-chip">{{ featuredAssignment.status.replace('_', ' ') }}</PillBadge>
                <span class="difficulty-tag">{{ featuredAssignment.activity.difficulty }}</span>
              </div>
            </div>
            <h3 class="kids-featured-title">{{ featuredAssignment.activity.title }}</h3>
            <p class="kids-featured-desc">{{ featuredAssignment.activity.description }}</p>
            
            <div class="kids-progress-shell">
              <div class="progress-row">
                <div class="progress-track">
                  <div class="progress-fill" :style="{ width: `${featuredAssignment.progress}%` }"></div>
                </div>
                <strong>{{ featuredAssignment.progress }}%</strong>
              </div>
            </div>

            <div class="card-actions kids-featured-actions">
              <ActionButton class="play-button" @click="playActivity(featuredAssignment)">
                Start Adventure
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

          <GlassCard v-if="featuredAssignment && gridAssignments.length === 0" class="activity-card empty-state-card kids-mini-empty">
             <div class="activity-card-top">
              <span class="activity-icon small-icon">✨</span>
              <p class="mini-desc">More adventures await!</p>
            </div>
            <p class="muted small-text">You are currently focusing on the featured activity above.</p>
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
  min-height: auto;
}

.kids-dashboard-split .kids-detail-panel {
  grid-area: detail;
  align-self: stretch;
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
  gap: 1.5rem;
  grid-template-columns: 1fr auto;
  align-items: center;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.85), rgba(255, 248, 220, 0.75));
}

.kids-hero-content h2 {
  font-size: 2.2rem;
  color: #ff7a18;
  text-shadow: 0 2px 4px rgba(255, 122, 24, 0.1);
}

.kids-hero-desc {
  font-size: 1.1rem;
  color: var(--muted);
  max-width: 40ch;
}

.kids-hero-badges {
  display: flex;
  gap: 1rem;
}

.kids-hero-chip {
  background: white;
  border: 2px solid rgba(255, 122, 24, 0.1);
  border-radius: 20px;
  padding: 0.8rem 1.2rem;
  text-align: center;
  min-width: 100px;
  box-shadow: 0 8px 20px rgba(255, 122, 24, 0.05);
}

.chip-label {
  display: block;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--muted);
  font-weight: 700;
}

.chip-value {
  display: block;
  font-size: 1.8rem;
  color: #ff7a18;
}

.kids-detail-panel {
  padding: 1.8rem;
  display: grid;
  gap: 1.25rem;
  min-width: 0;
  border: 2px solid rgba(22, 183, 214, 0.15);
}

.kids-detail-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.kids-featured-icon {
  font-size: 4rem;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1));
}

.kids-detail-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.5rem;
}

.difficulty-tag {
  font-size: 0.75rem;
  font-weight: 800;
  color: var(--secondary);
  text-transform: uppercase;
  background: rgba(22, 183, 214, 0.1);
  padding: 0.2rem 0.6rem;
  border-radius: 6px;
}

.kids-featured-title {
  font-size: 1.8rem;
  margin: 0;
}

.kids-featured-desc {
  font-size: 1.1rem;
  margin: 0;
  line-height: 1.4;
}

.kids-progress-shell {
  background: rgba(255, 255, 255, 0.5);
  padding: 1rem;
  border-radius: 18px;
  border: 1px solid rgba(24, 78, 122, 0.05);
}

.play-button {
  flex: 1.5;
  font-size: 1.1rem !important;
  box-shadow: 0 10px 25px rgba(255, 122, 24, 0.3) !important;
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

.kids-activity-card h3 {
  font-size: 1.3rem;
  margin: 0;
}

.kids-activity-card p {
  font-size: 0.95rem;
  color: var(--muted);
}

.kids-empty-panel {
  min-height: 220px;
}

.kids-mini-empty {
  min-height: auto;
  padding: 1rem 1.25rem;
  border-style: solid;
  border-width: 1px;
  border-color: rgba(24, 78, 122, 0.1);
  background: rgba(255, 255, 255, 0.4);
}

.small-icon {
  font-size: 1.4rem;
}

.mini-desc {
  margin: 0;
  font-weight: 800;
  font-size: 1rem;
  color: var(--ink);
}

.small-text {
  font-size: 0.85rem;
  margin-top: 0.4rem;
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

@media (max-width: 1180px) {
  .kids-hero-panel {
    grid-template-columns: 1fr;
    text-align: center;
  }
  
  .kids-hero-content {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .kids-hero-badges {
    justify-content: center;
  }
}

@media (max-width: 820px) {
  .hero-panel {
    border-radius: 22px;
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
