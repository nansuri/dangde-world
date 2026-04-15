<template>
  <AppShell
    title="Kids Zone"
    subtitle="Explore learning quests with bright colors, clear choices, and playful progress."
    :items="navItems"
    @select="selectCategory"
    @logout="logout"
  >
    <div class="page-workspace player-workspace">
      <div class="player-page" :class="{ 'player-page-split': isTabletLandscape }">
        <GlassCard tag="section" class="player-overview">
          <div class="player-head">
            <div>
              <PillBadge class="eyebrow">Activity Player</PillBadge>
              <h1>{{ assignment?.activity?.title || 'Loading activity...' }}</h1>
              <p class="muted">{{ assignment?.activity?.description || 'Please wait while we load your activity.' }}</p>
            </div>
            <PillBadge v-if="assignment" class="status-chip">{{ assignment.status.replace('_', ' ') }}</PillBadge>
          </div>

          <div v-if="assignment" class="player-summary-grid">
            <GlassCard tag="article" class="metric-card player-metric-card">
              <p>Progress</p>
              <strong>{{ assignment.progress }}%</strong>
            </GlassCard>
            <GlassCard tag="article" class="metric-card player-metric-card">
              <p>Language</p>
              <strong>{{ assignment.activity.language }}</strong>
            </GlassCard>
          </div>

          <div v-if="assignment" class="player-note">
            <PillBadge class="eyebrow">Prompt</PillBadge>
            <p>{{ assignment.activity.prompt }}</p>
          </div>

          <p v-if="actionError" class="error-text">{{ actionError }}</p>

          <div class="player-actions">
            <ActionButton variant="ghost" :disabled="isBusy" @click="goBack">Back</ActionButton>
            <ActionButton variant="ghost" :disabled="!assignment || isBusy" @click="reloadActivity">
              Reload
            </ActionButton>
            <ActionButton :disabled="!assignment || isBusy" @click="markCompleted">
              {{ isBusy ? 'Saving...' : 'Complete Activity' }}
            </ActionButton>
          </div>
        </GlassCard>

        <GlassCard v-if="error" tag="section" class="player-frame-shell">
          <p class="error-text">{{ error }}</p>
        </GlassCard>

        <GlassCard v-else-if="assignment" tag="section" class="player-frame-shell">
          <iframe
            ref="playerFrame"
            class="activity-player-frame"
            title="Activity Player"
            sandbox="allow-scripts allow-forms allow-modals allow-pointer-lock"
            :srcdoc="srcdoc"
          />
        </GlassCard>
      </div>
    </div>
  </AppShell>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppShell from '../../shared/ui/AppShell.vue'
import GlassCard from '../../shared/ui/GlassCard.vue'
import PillBadge from '../../shared/ui/PillBadge.vue'
import ActionButton from '../../shared/ui/ActionButton.vue'
import { clearSession, readSession } from '../../features/auth/session.js'
import { listAssignments, updateAssignment } from '../../entities/assignment/api.js'
import { useViewportProfile } from '../../shared/lib/useViewportProfile.js'
import {
  createActivityData,
  deleteActivityData,
  listActivityData,
  updateActivityData,
} from '../../entities/activity-data/api.js'

const route = useRoute()
const router = useRouter()
const playerFrame = ref(null)
const session = readSession()
const assignment = ref(null)
const error = ref('')
const actionError = ref('')
const isBusy = ref(false)
const scriptCloseTag = '</scr' + 'ipt>'
const styleCloseTag = '</sty' + 'le>'
const { isTabletLandscape } = useViewportProfile()
const selectedCategoryKey = ref('all')

const categoryItems = [
  { key: 'all', label: 'My Activities', icon: '🧩', matches: [] },
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

const bridgeContext = computed(() => ({
  activityId: assignment.value?.activityId,
  assignmentId: assignment.value?.id,
  kidId: session?.id,
}))

const srcdoc = computed(() => {
  if (!assignment.value?.activity) {
    return ''
  }

  const activity = assignment.value.activity
  const contextJSON = JSON.stringify(bridgeContext.value)
  const bridgeCode = `
    const pendingRequests = new Map();
    let requestCounter = 0;
    const context = ${contextJSON};

    function sendRequest(action, payload) {
      return new Promise((resolve, reject) => {
        const requestId = 'req-' + (++requestCounter);
        pendingRequests.set(requestId, { resolve, reject });
        parent.postMessage({ source: 'dangde-runtime', action, payload, requestId }, '*');
      });
    }

    window.addEventListener('message', (event) => {
      const data = event.data || {};
      if (data.source !== 'dangde-host' || !data.requestId) {
        return;
      }
      const pending = pendingRequests.get(data.requestId);
      if (!pending) {
        return;
      }
      pendingRequests.delete(data.requestId);
      if (data.error) {
        pending.reject(new Error(data.error));
        return;
      }
      pending.resolve(data.result);
    });

    window.DangdeAPI = {
      context,
      save: (payload) => sendRequest('save', payload),
      query: (payload = {}) => sendRequest('query', payload),
      update: (payload) => sendRequest('update', payload),
      delete: (payload) => sendRequest('delete', payload),
      setProgress: (payload) => sendRequest('setProgress', payload),
    };
  `.replaceAll(scriptCloseTag, '<\\/script>')

  const activityJS = (activity.jsCode || '').replaceAll(scriptCloseTag, '<\\/script>')
  const activityCSS = (activity.cssCode || '').replaceAll(styleCloseTag, '<\\/style>')

  return [
    '<!doctype html>',
    '<html><head><meta charset="UTF-8" />',
    '<meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover" />',
    '<style>html,body{margin:0;min-height:100%;}body{min-height:100vh;min-height:100svh;min-height:100dvh;}',
    activityCSS,
    '</style></head><body>',
    activity.htmlCode || '<div style="padding:24px;font-family:sans-serif">No activity code found.</div>',
    '<script>',
    bridgeCode,
    '<\/script><script>',
    activityJS,
    '<\/script></body></html>',
  ].join('')
})

onMounted(async () => {
  await loadAssignment()
  window.addEventListener('message', handleRuntimeMessage)
})

onUnmounted(() => {
  window.removeEventListener('message', handleRuntimeMessage)
})

async function loadAssignment() {
  try {
    actionError.value = ''
    const response = await listAssignments({ kidId: session.id })
    assignment.value = response.items.find((item) => String(item.id) === route.params.assignmentId)
    if (!assignment.value) {
      error.value = 'Activity assignment not found.'
      return
    }
    error.value = ''
    syncCategoryFromAssignment()
  } catch (err) {
    error.value = err.message
  }
}

async function handleRuntimeMessage(event) {
  if (event.source !== playerFrame.value?.contentWindow) {
    return
  }

  const data = event.data || {}
  if (data.source !== 'dangde-runtime' || !data.requestId) {
    return
  }

  try {
    const result = await runRuntimeAction(data.action, data.payload || {})
    event.source.postMessage({ source: 'dangde-host', requestId: data.requestId, result }, '*')
  } catch (err) {
    event.source.postMessage(
      { source: 'dangde-host', requestId: data.requestId, error: err.message || 'Runtime request failed' },
      '*',
    )
  }
}

async function runRuntimeAction(action, payload) {
  const context = bridgeContext.value

  switch (action) {
    case 'save':
      return upsertRuntimeValue(payload.key, payload.value)
    case 'query':
      return queryRuntimeValues(payload.key)
    case 'update':
      return updateRuntimeValue(payload)
    case 'delete':
      return deleteRuntimeValue(payload)
    case 'setProgress':
      return setActivityProgress(payload)
    default:
      throw new Error(`Unsupported runtime action: ${action}`)
  }

  async function queryRuntimeValues(key) {
    const response = await listActivityData({
      activityId: context.activityId,
      assignmentId: context.assignmentId,
      kidId: context.kidId,
      ...(key ? { key } : {}),
    })
    return response.items.map(parseRuntimeItem)
  }

  async function upsertRuntimeValue(key, value) {
    if (!key) {
      throw new Error('save requires a key')
    }

    const existing = await queryRuntimeValues(key)
    const serializedValue = JSON.stringify(value)

    if (existing.length) {
      const response = await updateActivityData(existing[0].id, { key, value: serializedValue })
      return parseRuntimeItem(response.item)
    }

    const response = await createActivityData({
      activityId: context.activityId,
      assignmentId: context.assignmentId,
      kidId: context.kidId,
      key,
      value: serializedValue,
    })
    return parseRuntimeItem(response.item)
  }

  async function updateRuntimeValue(runtimePayload) {
    const serializedValue = JSON.stringify(runtimePayload.value)

    if (runtimePayload.id) {
      const response = await updateActivityData(runtimePayload.id, {
        key: runtimePayload.key,
        value: serializedValue,
      })
      return parseRuntimeItem(response.item)
    }

    if (runtimePayload.key) {
      return upsertRuntimeValue(runtimePayload.key, runtimePayload.value)
    }

    throw new Error('update requires id or key')
  }

  async function deleteRuntimeValue(runtimePayload) {
    if (runtimePayload.id) {
      await deleteActivityData(runtimePayload.id)
      return { deleted: 1 }
    }

    if (!runtimePayload.key) {
      throw new Error('delete requires id or key')
    }

    const existing = await queryRuntimeValues(runtimePayload.key)
    await Promise.all(existing.map((item) => deleteActivityData(item.id)))
    return { deleted: existing.length }
  }

  async function setActivityProgress(runtimePayload) {
    const progress = Number(runtimePayload.progress || 0)
    const status = runtimePayload.status || (progress >= 100 ? 'completed' : 'in_progress')
    const response = await updateAssignment(context.assignmentId, { progress, status })
    assignment.value = response.item
    return response.item
  }
}

function parseRuntimeItem(item) {
  try {
    return { ...item, value: JSON.parse(item.value) }
  } catch {
    return item
  }
}

async function markCompleted() {
  if (!assignment.value) {
    return
  }
  try {
    isBusy.value = true
    actionError.value = ''
    const response = await updateAssignment(assignment.value.id, { progress: 100, status: 'completed' })
    assignment.value = response.item
  } catch (err) {
    actionError.value = err.message || 'Unable to complete activity right now.'
  } finally {
    isBusy.value = false
  }
}

function goBack() {
  router.push('/kids').catch(() => {})
}

function syncCategoryFromAssignment() {
  const slug = (assignment.value?.activity?.category?.slug || '').toLowerCase()
  if (!slug) {
    selectedCategoryKey.value = 'all'
    return
  }
  const matched = categoryItems.find((item) => item.matches.some((match) => slug.includes(match)))
  selectedCategoryKey.value = matched?.key || 'all'
}

function selectCategory(item) {
  selectedCategoryKey.value = item.key
  if (!assignment.value) {
    return
  }
  if (item.key === 'all') {
    return
  }

  const slug = (assignment.value.activity.category?.slug || '').toLowerCase()
  const isCurrentCategory = item.matches.some((match) => slug.includes(match))
  if (isCurrentCategory) {
    return
  }

  goBack()
}

function logout() {
  clearSession()
  router.push('/login').catch(() => {})
}

function reloadActivity() {
  // Resetting srcdoc by nudging assignment reference re-renders the iframe safely.
  if (assignment.value) {
    assignment.value = { ...assignment.value }
  }
}
</script>

<style scoped>
.page-workspace {
  height: 100%;
  min-height: 0;
  display: grid;
}

.player-workspace {
  grid-template-rows: minmax(0, 1fr);
  overflow: hidden;
}

.player-page {
  display: grid;
  grid-template-rows: auto minmax(72vh, 1fr);
  gap: 1rem;
  min-height: 0;
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.player-page-split {
  grid-template-columns: minmax(280px, 0.62fr) minmax(0, 2.2fr);
  grid-template-rows: minmax(0, 1fr);
  align-items: stretch;
}

.player-overview {
  padding: 0.75rem;
  display: grid;
  gap: 0.55rem;
  min-width: 0;
  align-content: start;
  max-height: 34vh;
  overflow-y: auto;
}

.player-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  min-width: 0;
}

.player-overview h1 {
  margin: 0.2rem 0 0.25rem;
  font-size: clamp(1.05rem, 1.8vw, 1.35rem);
  line-height: 1.2;
  word-break: break-word;
  overflow-wrap: break-word;
}

.player-overview .muted {
  margin: 0;
  font-size: 0.9rem;
  line-height: 1.35;
}

.player-summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.45rem;
  width: 100%;
}

.player-metric-card {
  min-height: 72px;
  border-radius: 16px;
  padding: 0.6rem 0.7rem;
}

.player-metric-card p {
  margin: 0;
  font-size: 0.78rem;
}

.player-metric-card strong {
  display: block;
  margin-top: 0.1rem;
  font-size: 1.05rem;
  line-height: 1.2;
  word-break: break-word;
}

.player-note {
  background: rgba(255, 255, 255, 0.62);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 0.55rem 0.65rem;
}

.player-note p {
  margin: 0.3rem 0 0;
  font-size: 0.88rem;
  line-height: 1.35;
}

.player-actions {
  display: flex;
  gap: 0.45rem;
  flex-wrap: wrap;
}

.player-actions :deep(.action-button) {
  min-height: 44px;
  padding: 0.62rem 0.85rem;
  font-size: 0.9rem;
  border-radius: 14px;
}

.player-frame-shell {
  padding: 0.75rem;
  min-height: 0;
  height: 100%;
  overflow: hidden;
  display: grid;
}

.activity-player-frame {
  width: 100%;
  height: 100%;
  border: 0;
  border-radius: 24px;
  background: #fff;
}

.error-text {
  color: #bf1650;
  margin: 0;
}

@media (max-width: 1180px) {
  .player-page-split {
    grid-template-columns: 1fr;
    grid-template-rows: auto minmax(72vh, 1fr);
  }

  .player-overview {
    max-height: 30vh;
  }
}

@media (max-width: 820px) {
  .player-page {
    grid-template-rows: auto minmax(68vh, 1fr);
  }

  .player-overview {
    max-height: 32vh;
  }

  .player-head {
    align-items: flex-start;
    flex-direction: column;
  }

  .player-summary-grid {
    grid-template-columns: 1fr;
  }

  .player-frame-shell {
    padding: 0.5rem;
  }
}
</style>
