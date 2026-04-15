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
        <section class="surface-card admin-subnav-shell">
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
        </section>

        <section
          v-if="selectedCurriculumPage === 'management'"
          class="surface-card admin-management-shell"
        >
          <div class="section-head">
            <div>
              <p class="eyebrow">Curriculum Management</p>
              <h2>{{ editingActivity ? 'Edit activity' : 'Create a new activity' }}</h2>
            </div>
            <button
              v-if="editingActivity"
              class="ghost-button"
              type="button"
              @click="startCreating"
            >
              New Activity
            </button>
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
        </section>

        <section v-if="selectedCurriculumPage === 'library'" class="surface-card">
          <div class="section-head">
            <div>
              <p class="eyebrow">Activity Library</p>
              <h2>Curriculum catalogue</h2>
            </div>
            <div class="admin-library-summary">
              <span class="status-chip">Page {{ currentPage }} / {{ totalPages }}</span>
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
                <button class="ghost-button" type="button" @click="editActivity(activity)">Edit</button>
              </div>
            </article>
          </div>

          <div class="admin-pagination">
            <button class="ghost-button" type="button" :disabled="currentPage === 1" @click="previousPage">
              Previous
            </button>
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
            <button
              class="ghost-button"
              type="button"
              :disabled="currentPage === totalPages"
              @click="nextPage"
            >
              Next
            </button>
          </div>
        </section>
      </div>

      <template v-if="selectedSection === 'categories'">
        <div class="admin-workspace">
        <section v-if="!editingCategoryId" class="surface-card">
          <div class="section-head">
            <div>
              <p class="eyebrow">Categories</p>
              <h2>Learning structure</h2>
            </div>
            <button class="primary-button" type="button" @click="startCreatingCategory">
              ➕ New Category
            </button>
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
        </section>

        <section v-if="rootCategories.length > 0 && !editingCategoryId" class="surface-card">
          <div class="section-head">
            <div>
              <p class="eyebrow">Sub-Categories</p>
              <h2>Detailed paths</h2>
            </div>
            <button class="primary-button" type="button" @click="startCreatingSubCategory">
              ➕ New Sub-Category
            </button>
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
        </section>

        <section v-if="editingCategoryId" class="surface-card">
          <div class="section-head">
            <div>
              <p class="eyebrow">Category Editor</p>
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
        </section>
        </div>
      </template>

      <template v-if="selectedSection === 'analytics'">
        <div class="admin-workspace">
        <section class="metrics-grid">
          <article class="metric-card">
            <p>Activities</p>
            <strong>{{ activities.length }}</strong>
          </article>
          <article class="metric-card">
            <p>Root Categories</p>
            <strong>{{ rootCategories.length }}</strong>
          </article>
          <article class="metric-card">
            <p>Sub-Categories</p>
            <strong>{{ subCategories.length }}</strong>
          </article>
          <article class="metric-card">
            <p>Code-Based</p>
            <strong>{{ codedActivities }}</strong>
          </article>
        </section>

        <section class="surface-card">
          <div class="section-head">
            <div>
              <p class="eyebrow">Coverage</p>
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
        </section>
        </div>
      </template>
    </div>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '../../shared/ui/AppShell.vue'
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
