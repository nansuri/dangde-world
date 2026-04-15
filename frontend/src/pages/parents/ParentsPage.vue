<template>
  <AppShell
    title="Parents Hub"
    subtitle="Track each child, assign activities, and monitor learning momentum."
    :items="navItems"
    @select="selectSection"
    @logout="logout"
  >
    <div class="page-workspace parents-workspace">
      <div class="admin-workspace">
        <section v-if="selectedSection === 'overview'" class="metrics-grid">
          <MetricCard label="Kids" :value="stats.kidsCount" />
          <MetricCard label="Assignments" :value="stats.assignmentsCount" />
          <MetricCard label="Completed" :value="stats.completedCount" />
          <MetricCard label="Average Progress" :value="`${stats.averageProgress}%`" />
        </section>

        <GlassCard v-if="selectedSection === 'overview'" tag="section" class="parents-card">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Family Snapshot</PillBadge>
              <h2>Your kids at a glance</h2>
            </div>
            <ActionButton variant="ghost" @click="selectedSection = 'kids'">
              Manage Family
            </ActionButton>
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
        </GlassCard>

        <GlassCard v-if="selectedSection === 'kids' && !isEditingKid" tag="section" class="parents-card">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Family Management</PillBadge>
              <h2>My kids</h2>
            </div>
            <ActionButton @click="startAddingKid">
              ➕ Add Kid
            </ActionButton>
          </div>
          <div class="assignment-list">
            <article v-for="kid in kids" :key="kid.id" class="assignment-row category-row">
              <div>
                <strong>{{ kid.avatar }} {{ kid.name }}</strong>
                <p>{{ kid.preferredLang === 'id' ? 'Bahasa Indonesia' : 'English' }}</p>
              </div>
              <div>
                <p class="muted">Activities</p>
                <strong>{{ countAssignmentsForKid(kid.id) }}</strong>
              </div>
              <div class="category-actions">
                <button class="icon-button" type="button" @click="editKid(kid)" title="Edit">
                  ✏️
                </button>
                <button class="icon-button delete" type="button" @click="deleteKidItem(kid)" title="Delete">
                  🗑️
                </button>
              </div>
            </article>
          </div>
        </GlassCard>

        <GlassCard v-if="selectedSection === 'kids' && isEditingKid" tag="section" class="parents-card">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Profile Editor</PillBadge>
              <h2>{{ kidForm.id ? 'Edit kid profile' : 'Add a new kid' }}</h2>
            </div>
          </div>
          <form class="activity-form" @submit.prevent="handleKidSubmit">
            <label>
              Kid Name
              <input v-model="kidForm.name" type="text" required placeholder="e.g., Nino" />
            </label>
            <label>
              Avatar
              <select v-model="kidForm.avatar">
                <option value="🐼">🐼 Panda</option>
                <option value="🦁">🦁 Lion</option>
                <option value="🦊">🦊 Fox</option>
                <option value="🐨">🐨 Koala</option>
                <option value="🐰">🐰 Rabbit</option>
                <option value="🐯">🐯 Tiger</option>
              </select>
            </label>
            <label>
              Preferred Language
              <select v-model="kidForm.preferredLang">
                <option value="en">English</option>
                <option value="id">Bahasa Indonesia</option>
              </select>
            </label>

            <div class="card-actions profile-actions">
              <ActionButton variant="ghost" @click="isEditingKid = false">Cancel</ActionButton>
              <ActionButton type="submit">{{ kidForm.id ? 'Update Kid' : 'Add Kid' }}</ActionButton>
            </div>
          </form>
        </GlassCard>

        <GlassCard v-if="selectedSection === 'assignments'" tag="section" class="parents-card">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Assign Learning</PillBadge>
              <h2>Create a new activity assignment</h2>
            </div>
          </div>
          <AssignmentComposer
            :kids="kids"
            :activities="activities"
            :parent-id="session.id"
            @assign="assignActivity"
          />
        </GlassCard>

        <GlassCard v-if="selectedSection === 'assignments'" tag="section" class="parents-card">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Current Progress</PillBadge>
              <h2>Assignments for your kids</h2>
            </div>
          </div>
          <div class="assignment-list">
            <article v-for="assignment in assignments" :key="assignment.id" class="assignment-row category-row">
              <div>
                <strong>{{ assignment.kid.name }}</strong>
                <p>{{ assignment.activity.title }} · {{ assignment.activity.language }}</p>
              </div>
              <div>
                <p class="muted">{{ assignment.status.replace('_', ' ') }}</p>
                <strong>{{ assignment.progress }}%</strong>
              </div>
              <div class="category-actions">
                <button class="icon-button delete" type="button" @click="deleteAssignmentItem(assignment)" title="Delete Assignment">
                  🗑️
                </button>
              </div>
            </article>
          </div>
        </GlassCard>

        <GlassCard v-if="selectedSection === 'reports'" tag="section" class="parents-card">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Reports</PillBadge>
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
        </GlassCard>
      </div>
    </div>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '../../shared/ui/AppShell.vue'
import GlassCard from '../../shared/ui/GlassCard.vue'
import PillBadge from '../../shared/ui/PillBadge.vue'
import ActionButton from '../../shared/ui/ActionButton.vue'
import MetricCard from '../../widgets/dashboard/MetricCard.vue'
import AssignmentComposer from '../../features/assignment/AssignmentComposer.vue'
import { readSession, clearSession } from '../../features/auth/session.js'
import { createUser, updateUser, deleteUser, listUsers } from '../../entities/user/api.js'
import { listActivities } from '../../entities/activity/api.js'
import { createAssignment, getParentStats, listAssignments, deleteAssignment } from '../../entities/assignment/api.js'

const router = useRouter()
const session = readSession()
const kids = ref([])
const activities = ref([])
const assignments = ref([])
const selectedSection = ref('overview')
const isEditingKid = ref(false)
const kidForm = ref({
  id: null,
  name: '',
  avatar: '🐼',
  preferredLang: 'en'
})

const stats = reactive({
  kidsCount: 0,
  assignmentsCount: 0,
  completedCount: 0,
  averageProgress: 0,
})

const sectionItems = [
  { key: 'overview', label: 'Overview', icon: '📊' },
  { key: 'kids', label: 'Kids', icon: '👶' },
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
  isEditingKid.value = false
}

function editKid(kid) {
  kidForm.value = { ...kid }
  isEditingKid.value = true
  selectedSection.value = 'kids'
}

function startAddingKid() {
  kidForm.value = { id: null, name: '', avatar: '🐼', preferredLang: 'en' }
  isEditingKid.value = true
}

async function handleKidSubmit() {
  const payload = {
    ...kidForm.value,
    role: 'kid',
    parentId: session.id
  }

  try {
    if (payload.id) {
      await updateUser(payload.id, payload)
    } else {
      await createUser(payload)
    }
    await loadData()
    isEditingKid.value = false
  } catch (error) {
    console.error('Error saving kid:', error)
    alert('Failed to save kid: ' + error.message)
  }
}

async function deleteKidItem(kid) {
  if (!confirm(`Delete kid "${kid.name}"? This will also remove their assignments and progress.`)) {
    return
  }
  try {
    await deleteUser(kid.id)
    await loadData()
  } catch (error) {
    console.error('Error deleting kid:', error)
    alert('Failed to delete kid: ' + error.message)
  }
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

async function deleteAssignmentItem(assignment) {
  if (!confirm(`Delete assignment for ${assignment.kid.name}? This cannot be undone.`)) {
    return
  }
  try {
    await deleteAssignment(assignment.id)
    await loadData()
  } catch (error) {
    console.error('Error deleting assignment:', error)
    alert('Failed to delete assignment: ' + error.message)
  }
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

.parents-workspace {
  grid-template-rows: minmax(0, 1fr);
  overflow: hidden;
}

.admin-workspace {
  display: grid;
  grid-template-rows: auto;
  gap: 1rem;
  min-height: 0;
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.admin-workspace > section {
  min-height: 0;
  overflow: visible;
}

.parents-card {
  padding: 1.2rem 1.25rem;
}

.metrics-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.1rem;
}

.section-head h2 {
  margin: 0.35rem 0 0;
  line-height: 1.2;
}

.assignment-list {
  display: grid;
  gap: 0.8rem;
}

.assignment-row {
  padding: 0.95rem 1rem;
  border: 1px solid rgba(24, 78, 122, 0.1);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.56);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.assignment-row strong {
  display: block;
  margin-bottom: 0.2rem;
}

.assignment-row p {
  margin: 0;
}

.assignment-row .muted {
  font-size: 0.88rem;
}

.category-row {
  flex-wrap: wrap;
}

.category-actions {
  display: flex;
  gap: 0.5rem;
  margin-left: auto;
}

.icon-button {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(24, 78, 122, 0.12);
  font-size: 1.1rem;
  cursor: pointer;
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  padding: 0;
  transition: all 200ms;
  border-radius: 12px;
}

.icon-button:hover {
  background: rgba(255, 255, 255, 0.95);
  transform: translateY(-1px);
  box-shadow: 0 8px 16px rgba(24, 78, 122, 0.12);
}

.icon-button.delete:hover {
  background: rgba(244, 67, 54, 0.1);
}

.activity-form {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

label {
  display: grid;
  gap: 0.45rem;
  font-weight: 700;
}

label input,
label select {
  min-height: 46px;
  border-radius: 12px;
  border: 1px solid rgba(24, 78, 122, 0.14);
  padding: 0.6rem 0.75rem;
  font: inherit;
  background: rgba(255, 255, 255, 0.88);
}

label input:focus,
label select:focus {
  outline: none;
  border-color: rgba(22, 183, 214, 0.5);
  box-shadow: 0 0 0 3px rgba(22, 183, 214, 0.14);
}

.card-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.profile-actions {
  grid-column: 1 / -1;
  margin-top: 0.65rem;
  justify-content: flex-end;
}

@media (max-width: 820px) {
  .metrics-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .parents-card {
    padding: 1rem;
  }

  .section-head {
    align-items: flex-start;
    flex-direction: column;
  }

  .section-head h2 {
    margin-top: 0.25rem;
  }

  .assignment-row {
    padding: 0.85rem;
  }

  .activity-form {
    grid-template-columns: 1fr;
  }

  .profile-actions {
    justify-content: stretch;
  }
}
</style>
