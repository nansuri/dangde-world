<template>
  <AppShell
    title="Parents Hub"
    subtitle="Track each child, assign activities, and monitor learning momentum."
    :items="navItems"
    @select="selectSection"
    @logout="logout"
  >
    <div class="page-workspace parents-workspace">
      <section v-if="selectedSection === 'overview'" class="metrics-grid">
        <MetricCard label="Kids" :value="stats.kidsCount" />
        <MetricCard label="Assignments" :value="stats.assignmentsCount" />
        <MetricCard label="Completed" :value="stats.completedCount" />
        <MetricCard label="Average Progress" :value="`${stats.averageProgress}%`" />
      </section>

      <section v-if="selectedSection === 'overview'" class="surface-card">
        <div class="section-head">
          <div>
            <p class="eyebrow">Family Snapshot</p>
            <h2>Your kids at a glance</h2>
          </div>
        </div>
        <div class="assignment-list">
          <article v-for="kid in kids" :key="kid.id" class="assignment-row">
            <div>
              <strong>{{ kid.avatar }} {{ kid.name }}</strong>
              <p>{{ kid.preferredLang === 'id' ? 'Bahasa Indonesia' : 'English' }} learner</p>
            </div>
            <div>
              <p class="muted">Assigned activities</p>
              <strong>{{ countAssignmentsForKid(kid.id) }}</strong>
            </div>
          </article>
        </div>
      </section>

      <section v-if="selectedSection === 'assignments'" class="surface-card">
        <div class="section-head">
          <div>
            <p class="eyebrow">Assign Learning</p>
            <h2>Create a new activity assignment</h2>
          </div>
        </div>
        <AssignmentComposer
          :kids="kids"
          :activities="activities"
          :parent-id="session.id"
          @assign="assignActivity"
        />
      </section>

      <section v-if="selectedSection === 'assignments'" class="surface-card">
        <div class="section-head">
          <div>
            <p class="eyebrow">Current Progress</p>
            <h2>Assignments for your kids</h2>
          </div>
        </div>
        <div class="assignment-list">
          <article v-for="assignment in assignments" :key="assignment.id" class="assignment-row">
            <div>
              <strong>{{ assignment.kid.name }}</strong>
              <p>{{ assignment.activity.title }} · {{ assignment.activity.language }}</p>
            </div>
            <div>
              <p class="muted">{{ assignment.status.replace('_', ' ') }}</p>
              <strong>{{ assignment.progress }}%</strong>
            </div>
          </article>
        </div>
      </section>

      <section v-if="selectedSection === 'reports'" class="surface-card">
        <div class="section-head">
          <div>
            <p class="eyebrow">Reports</p>
            <h2>Progress by child</h2>
          </div>
        </div>
        <div class="assignment-list">
          <article v-for="kid in kids" :key="kid.id" class="assignment-row">
            <div>
              <strong>{{ kid.name }}</strong>
              <p>{{ summarizeKidStatus(kid.id) }}</p>
            </div>
            <div>
              <p class="muted">Average progress</p>
              <strong>{{ averageKidProgress(kid.id) }}%</strong>
            </div>
          </article>
        </div>
      </section>
    </div>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '../../shared/ui/AppShell.vue'
import MetricCard from '../../widgets/dashboard/MetricCard.vue'
import AssignmentComposer from '../../features/assignment/AssignmentComposer.vue'
import { readSession, clearSession } from '../../features/auth/session.js'
import { listUsers } from '../../entities/user/api.js'
import { listActivities } from '../../entities/activity/api.js'
import { createAssignment, getParentStats, listAssignments } from '../../entities/assignment/api.js'

const router = useRouter()
const session = readSession()
const kids = ref([])
const activities = ref([])
const assignments = ref([])
const selectedSection = ref('overview')
const stats = reactive({
  kidsCount: 0,
  assignmentsCount: 0,
  completedCount: 0,
  averageProgress: 0,
})

const sectionItems = [
  { key: 'overview', label: 'Overview', icon: '📊' },
  { key: 'assignments', label: 'Assignments', icon: '🧩' },
  { key: 'reports', label: 'Reports', icon: '🌈' },
]

const navItems = computed(() =>
  sectionItems.map((item) => ({
    ...item,
    active: item.key === selectedSection.value,
  })),
)

onMounted(loadData)

async function loadData() {
  const [kidsResponse, activitiesResponse, assignmentsResponse, statsResponse] = await Promise.all([
    listUsers({ role: 'kid', parentId: session.id }),
    listActivities(),
    listAssignments({ parentId: session.id }),
    getParentStats(session.id),
  ])

  kids.value = kidsResponse.items
  activities.value = activitiesResponse.items
  assignments.value = assignmentsResponse.items
  Object.assign(stats, statsResponse)
}

async function assignActivity(payload) {
  await createAssignment(payload)
  await loadData()
  selectedSection.value = 'assignments'
}

function selectSection(item) {
  selectedSection.value = item.key
}

function countAssignmentsForKid(kidId) {
  return assignments.value.filter((assignment) => assignment.kid.id === kidId).length
}

function averageKidProgress(kidId) {
  const kidAssignments = assignments.value.filter((assignment) => assignment.kid.id === kidId)
  if (!kidAssignments.length) {
    return 0
  }
  const total = kidAssignments.reduce((sum, assignment) => sum + assignment.progress, 0)
  return Math.round(total / kidAssignments.length)
}

function summarizeKidStatus(kidId) {
  const kidAssignments = assignments.value.filter((assignment) => assignment.kid.id === kidId)
  if (!kidAssignments.length) {
    return 'No assignments yet.'
  }

  const completed = kidAssignments.filter((assignment) => assignment.status === 'completed').length
  const inProgress = kidAssignments.filter((assignment) => assignment.status === 'in_progress').length
  return `${completed} completed, ${inProgress} in progress, ${kidAssignments.length} total`
}

function logout() {
  clearSession()
  router.push('/login')
}
</script>
