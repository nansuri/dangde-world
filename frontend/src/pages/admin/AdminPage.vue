<template>
  <AppShell
    title="Admin Studio"
    subtitle="Manage curriculum, categories, and learning activities from one colorful control room."
    :items="navItems"
    @select="selectSection"
    @logout="logout"
  >
    <div class="page-workspace admin-page-workspace">
      <div
        v-if="selectedSection === 'curriculum'"
        class="admin-workspace"
        :class="{ 'admin-workspace-management': selectedCurriculumPage === 'management' }"
      >
        <GlassCard tag="section" class="admin-subnav-shell">
          <div class="admin-subnav">
            <button
              v-for="item in curriculumSubPages"
              :key="item.key"
              class="nav-pill"
              type="button"
              :class="{ active: item.key === selectedCurriculumPage }"
              @click="selectedCurriculumPage = item.key"
            >
              <span>{{ item.icon }}</span>
              <span>{{ item.label }}</span>
            </button>
          </div>
        </GlassCard>

        <GlassCard
          v-if="selectedCurriculumPage === 'management'"
          tag="section"
          class="admin-management-shell"
        >
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Curriculum Management</PillBadge>
              <h2>{{ editingActivity ? 'Edit activity' : 'Create a new activity' }}</h2>
            </div>
            <ActionButton
              v-if="editingActivity"
              variant="ghost"
              @click="startCreating"
            >
              New Activity
            </ActionButton>
          </div>

          <ActivityTemplates
            v-if="!editingActivity"
            @select="applyTemplate"
          />

          <ActivityForm
            :categories="subCategories"
            :model-value="editingActivity"
            :submit-label="editingActivity ? 'Update Activity' : 'Add Activity'"
            @submit="handleActivitySubmit"
          />
        </GlassCard>

        <GlassCard v-if="selectedCurriculumPage === 'library'" tag="section">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Activity Library</PillBadge>
              <h2>Curriculum catalogue</h2>
            </div>
            <div class="admin-library-summary">
              <PillBadge class="status-chip">Page {{ currentPage }} / {{ totalPages }}</PillBadge>
              <span class="muted">{{ activities.length }} total activities</span>
            </div>
          </div>

          <div class="admin-table-shell">
            <div class="admin-table admin-table-head">
              <span>Title</span>
              <span>Category</span>
              <span>Difficulty</span>
              <span>Age Group</span>
              <span>Actions</span>
            </div>

            <article
              v-for="activity in paginatedActivities"
              :key="activity.id"
              class="admin-table admin-table-row"
              :class="{ active: editingActivity?.id === activity.id }"
            >
              <div>
                <strong>{{ activity.icon }} {{ activity.title }}</strong>
                <p class="muted">{{ activity.description }}</p>
              </div>
              <span>{{ activity.category.name }}</span>
              <span>{{ activity.difficulty }}</span>
              <span>{{ activity.ageGroup }}</span>
              <div class="admin-row-actions">
                <ActionButton variant="ghost" @click="editActivity(activity)">Edit</ActionButton>
              </div>
            </article>
          </div>

          <div class="admin-pagination">
            <ActionButton variant="ghost" :disabled="currentPage === 1" @click="previousPage">
              Previous
            </ActionButton>
            <div class="admin-page-indicators">
              <button
                v-for="page in visiblePages"
                :key="page"
                class="nav-pill"
                type="button"
                :class="{ active: page === currentPage }"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
            </div>
            <ActionButton
              variant="ghost"
              :disabled="currentPage === totalPages"
              @click="nextPage"
            >
              Next
            </ActionButton>
          </div>
        </GlassCard>
      </div>

      <template v-if="selectedSection === 'categories'">
        <div class="admin-workspace">
        <GlassCard v-if="!editingCategoryId" tag="section">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Categories</PillBadge>
              <h2>Learning structure</h2>
            </div>
            <ActionButton @click="startCreatingCategory">
              ➕ New Category
            </ActionButton>
          </div>
          <div class="assignment-list">
            <article
              v-for="category in rootCategories"
              :key="category.id"
              class="assignment-row category-row"
            >
              <div>
                <strong>{{ category.name }}</strong>
                <p>{{ category.type }}</p>
              </div>
              <div>
                <p class="muted">Sub-categories</p>
                <strong>{{ countChildren(category.id) }}</strong>
              </div>
              <div class="category-actions">
                <button class="icon-button" type="button" @click="editCategory(category)" title="Edit">
                  ✏️
                </button>
                <button class="icon-button delete" type="button" @click="deleteCategoryItem(category)" title="Delete">
                  🗑️
                </button>
              </div>
            </article>
          </div>
        </GlassCard>

        <GlassCard v-if="rootCategories.length > 0 && !editingCategoryId" tag="section">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Sub-Categories</PillBadge>
              <h2>Detailed paths</h2>
            </div>
            <ActionButton @click="startCreatingSubCategory">
              ➕ New Sub-Category
            </ActionButton>
          </div>
          <div class="assignment-list">
            <article
              v-for="category in subCategories"
              :key="category.id"
              class="assignment-row category-row"
            >
              <div>
                <strong>{{ category.name }}</strong>
                <p>{{ parentCategoryName(category.parentId) }}</p>
              </div>
              <div>
                <p class="muted">Activities</p>
                <strong>{{ countActivities(category.id) }}</strong>
              </div>
              <div class="category-actions">
                <button class="icon-button" type="button" @click="editCategory(category)" title="Edit">
                  ✏️
                </button>
                <button class="icon-button delete" type="button" @click="deleteCategoryItem(category)" title="Delete">
                  🗑️
                </button>
              </div>
            </article>
          </div>
        </GlassCard>

        <GlassCard v-if="editingCategoryId" tag="section">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Category Editor</PillBadge>
              <h2>{{ editingCategory ? 'Edit category' : 'Create category' }}</h2>
            </div>
          </div>
          <CategoryForm
            :categories="categories"
            :model-value="editingCategory"
            :submit-label="editingCategory ? 'Update Category' : 'Create Category'"
            @submit="handleCategorySubmit"
            @cancel="cancelEditingCategory"
          />
        </GlassCard>
        </div>
      </template>

      <template v-if="selectedSection === 'analytics'">
        <div class="admin-workspace">
        <section class="metrics-grid">
          <GlassCard tag="article" class="metric-card">
            <p>Activities</p>
            <strong>{{ activities.length }}</strong>
          </GlassCard>
          <GlassCard tag="article" class="metric-card">
            <p>Root Categories</p>
            <strong>{{ rootCategories.length }}</strong>
          </GlassCard>
          <GlassCard tag="article" class="metric-card">
            <p>Sub-Categories</p>
            <strong>{{ subCategories.length }}</strong>
          </GlassCard>
          <GlassCard tag="article" class="metric-card">
            <p>Code-Based</p>
            <strong>{{ codedActivities }}</strong>
          </GlassCard>
        </section>

        <GlassCard tag="section">
          <div class="section-head">
            <div>
              <PillBadge class="eyebrow">Coverage</PillBadge>
              <h2>Activities by category</h2>
            </div>
          </div>
          <div class="assignment-list">
            <article v-for="category in subCategories" :key="category.id" class="assignment-row">
              <div>
                <strong>{{ category.name }}</strong>
                <p>{{ parentCategoryName(category.parentId) }}</p>
              </div>
              <div>
                <p class="muted">Activities</p>
                <strong>{{ countActivities(category.id) }}</strong>
              </div>
            </article>
          </div>
        </GlassCard>
        </div>
      </template>
    </div>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '../../shared/ui/AppShell.vue'
import GlassCard from '../../shared/ui/GlassCard.vue'
import PillBadge from '../../shared/ui/PillBadge.vue'
import ActionButton from '../../shared/ui/ActionButton.vue'
import ActivityForm from '../../features/activity-management/ActivityForm.vue'
import ActivityTemplates from '../../features/activity-management/ActivityTemplates.vue'
import CategoryForm from '../../features/category-management/CategoryForm.vue'
import { clearSession } from '../../features/auth/session.js'
import { listCategories, createCategory, updateCategory, deleteCategory } from '../../entities/category/api.js'
import { createActivity, listActivities, updateActivity } from '../../entities/activity/api.js'

const router = useRouter()
const categories = ref([])
const activities = ref([])
const selectedSection = ref('curriculum')
const editingActivityId = ref(null)
const editingCategoryId = ref(null)
const currentPage = ref(1)
const selectedCurriculumPage = ref('management')
const pageSize = 4

const sectionItems = [
  { key: 'curriculum', label: 'Curriculum', icon: '🛠️' },
  { key: 'categories', label: 'Categories', icon: '🗂️' },
  { key: 'analytics', label: 'Analytics', icon: '📈' },
]
const curriculumSubPages = [
  { key: 'management', label: 'Management', icon: '✍️' },
  { key: 'library', label: 'Library', icon: '📚' },
]

const navItems = computed(() =>
  sectionItems.map((item) => ({
    ...item,
    active: item.key === selectedSection.value,
  })),
)

const subCategories = computed(() => categories.value.filter((item) => item.parentId))
const rootCategories = computed(() => categories.value.filter((item) => !item.parentId))
const codedActivities = computed(
  () => activities.value.filter((item) => item.htmlCode || item.cssCode || item.jsCode).length,
)
const editingActivity = computed(
  () => activities.value.find((item) => item.id === editingActivityId.value) || null,
)
const editingCategory = computed(
  () => categories.value.find((item) => item.id === editingCategoryId.value) || null,
)
const totalPages = computed(() => Math.max(1, Math.ceil(activities.value.length / pageSize)))
const paginatedActivities = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return activities.value.slice(start, start + pageSize)
})
const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 1)
  const end = Math.min(totalPages.value, start + 2)
  for (let page = start; page <= end; page += 1) {
    pages.push(page)
  }
  return pages
})

onMounted(loadData)

async function loadData() {
  const [categoriesResponse, activitiesResponse] = await Promise.all([
    listCategories(),
    listActivities(),
  ])

  categories.value = categoriesResponse.items
  activities.value = activitiesResponse.items
}

async function createNewActivity(payload) {
  await createActivity(payload)
  await loadData()
  startCreating()
}

async function saveActivity(id, payload) {
  await updateActivity(id, payload)
  await loadData()
}

async function handleActivitySubmit(payload) {
  if (editingActivity.value) {
    await saveActivity(editingActivity.value.id, payload)
    return
  }
  await createNewActivity(payload)
}

function applyTemplate(template) {
  editingActivityId.value = null
  selectedCurriculumPage.value = 'management'
  // Set a small delay to allow the form to render before setting values
  setTimeout(() => {
    const form = {
      title: template.name,
      description: template.description,
      prompt: template.description,
      language: 'English',
      difficulty: 'easy',
      ageGroup: '4-8',
      categoryId: subCategories.value[0]?.id || 1,
      icon: template.icon,
      htmlCode: template.htmlCode,
      cssCode: template.cssCode,
      jsCode: template.jsCode,
    }
    editingActivityId.value = 'template-' + template.id
    Object.assign(editingActivity.value || {}, form)
  }, 0)
}

function selectSection(item) {
  selectedSection.value = item.key
}

function editActivity(activity) {
  editingActivityId.value = activity.id
  selectedCurriculumPage.value = 'management'
}

function startCreating() {
  editingActivityId.value = null
  selectedCurriculumPage.value = 'management'
}

function startCreatingCategory() {
  editingCategoryId.value = null
}

function startCreatingSubCategory() {
  editingCategoryId.value = null
}

async function handleCategorySubmit(payload) {
  try {
    if (editingCategory.value) {
      // Update
      await updateCategory(editingCategory.value.id, payload)
    } else {
      // Create
      await createCategory(payload)
    }
    await loadData()
    editingCategoryId.value = null
  } catch (error) {
    console.error('Error saving category:', error)
    alert('Failed to save category: ' + error.message)
  }
}

function editCategory(category) {
  editingCategoryId.value = category.id
}

async function deleteCategoryItem(category) {
  if (!confirm(`Delete category "${category.name}"? This cannot be undone.`)) {
    return
  }
  try {
    await deleteCategory(category.id)
    await loadData()
  } catch (error) {
    console.error('Error deleting category:', error)
    alert('Failed to delete category: ' + error.message)
  }
}

function cancelEditingCategory() {
  editingCategoryId.value = null
}

function goToPage(page) {
  currentPage.value = page
}

function previousPage() {
  currentPage.value = Math.max(1, currentPage.value - 1)
}

function nextPage() {
  currentPage.value = Math.min(totalPages.value, currentPage.value + 1)
}

function countChildren(parentId) {
  return categories.value.filter((item) => item.parentId === parentId).length
}

function countActivities(categoryId) {
  return activities.value.filter((item) => item.categoryId === categoryId).length
}

function parentCategoryName(parentId) {
  return categories.value.find((item) => item.id === parentId)?.name || 'No parent'
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

.admin-page-workspace {
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

.admin-workspace-management {
  grid-template-rows: auto minmax(0, 1fr);
}

.admin-management-shell {
  min-height: 0;
  display: grid;
  grid-template-rows: auto minmax(0, 1fr);
  overflow: hidden;
  padding: 1.5rem;
}

.admin-subnav-shell {
  padding: 1rem;
}

.admin-subnav {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.admin-subnav .nav-pill {
  width: auto;
  min-width: 180px;
  justify-content: center;
  border: 0;
  border-radius: 18px;
  min-height: 58px;
  padding: 1rem 1.2rem;
  font-weight: 800;
  transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
  background: rgba(255, 255, 255, 0.58);
  display: flex;
  align-items: center;
  gap: 0.7rem;
  cursor: pointer;
}

.admin-subnav .nav-pill.active {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.85);
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.admin-library-summary {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.admin-table-shell {
  display: grid;
  gap: 0.75rem;
}

.admin-table {
  display: grid;
  grid-template-columns: minmax(0, 2.2fr) minmax(120px, 1fr) 110px 110px 120px;
  gap: 1rem;
  align-items: center;
}

.admin-table-head {
  padding: 0 0.5rem 0.5rem;
  color: var(--muted);
  font-size: 0.9rem;
  font-weight: 700;
  border-bottom: 1px solid rgba(24, 78, 122, 0.1);
}

.admin-table-row {
  padding: 1rem;
  border: 1px solid rgba(24, 78, 122, 0.08);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.58);
}

.admin-table-row.active {
  border-color: rgba(22, 183, 214, 0.35);
  box-shadow: 0 14px 30px rgba(22, 183, 214, 0.12);
}

.admin-row-actions {
  display: flex;
  justify-content: flex-end;
}

.admin-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-top: 1.5rem;
}

.admin-page-indicators {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.admin-page-indicators .nav-pill {
  min-width: 58px;
  justify-content: center;
  border: 0;
  border-radius: 18px;
  min-height: 58px;
  padding: 1rem 1.2rem;
  font-weight: 800;
  background: rgba(255, 255, 255, 0.58);
  cursor: pointer;
}

.admin-page-indicators .nav-pill.active {
  background: rgba(255, 255, 255, 0.85);
}

.assignment-list {
  display: grid;
  gap: 1rem;
}

.assignment-row {
  padding: 1rem 0;
  border-bottom: 1px solid rgba(24, 78, 122, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.assignment-row:last-child {
  border-bottom: 0;
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
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0.5rem;
  transition: all 200ms;
  border-radius: 4px;
}

.icon-button:hover {
  background: rgba(0, 0, 0, 0.08);
  transform: scale(1.1);
}

.icon-button.delete:hover {
  background: rgba(244, 67, 54, 0.1);
}

.metrics-grid {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.metric-card {
  padding: 1.2rem;
}

.metric-card strong {
  display: block;
  margin-top: 0.4rem;
  font-size: 2rem;
  word-break: break-word;
  overflow-wrap: break-word;
}

@media (max-width: 1180px) {
  .admin-activity-card {
    grid-template-columns: 1fr;
  }

  .admin-table {
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }

  .admin-table-head {
    display: none;
  }

  .admin-row-actions {
    justify-content: flex-start;
  }

  .admin-pagination {
    flex-direction: column;
    align-items: stretch;
  }

  .admin-subnav {
    display: grid;
  }

  .admin-subnav .nav-pill {
    width: 100%;
    min-width: 0;
  }
}

@media (max-width: 820px) {
  .metrics-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .section-head {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
