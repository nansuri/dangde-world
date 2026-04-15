<template>
  <div class="player-page" :class="{ 'player-page-split': isTabletLandscape }">
    <section class="player-topbar" :class="{ 'player-side-panel': isTabletLandscape }">
      <div>
        <p class="eyebrow">Activity Player</p>
        <h1>{{ assignment?.activity?.title || 'Loading activity...' }}</h1>
        <p class="muted">{{ assignment?.activity?.description }}</p>
      </div>
      <div v-if="assignment" class="player-summary-grid">
        <article class="metric-card player-metric-card">
          <p>Status</p>
          <strong>{{ assignment.status.replace('_', ' ') }}</strong>
        </article>
        <article class="metric-card player-metric-card">
          <p>Progress</p>
          <strong>{{ assignment.progress }}%</strong>
        </article>
      </div>
      <div v-if="assignment" class="player-note">
        <p class="eyebrow">Prompt</p>
        <p>{{ assignment.activity.prompt }}</p>
      </div>
      <div class="player-actions">
        <button class="ghost-button" type="button" @click="goBack">Back</button>
        <button class="primary-button" type="button" @click="markCompleted" :disabled="!assignment">
          Complete Activity
        </button>
      </div>
    </section>

    <section v-if="error" class="surface-card">
      <p class="error-text">{{ error }}</p>
    </section>

    <section v-else-if="assignment" class="player-frame-shell">
      <iframe
        ref="playerFrame"
        class="activity-player-frame"
        title="Activity Player"
        sandbox="allow-scripts"
        :srcdoc="srcdoc"
      />
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { readSession } from '../../features/auth/session.js'
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
const scriptCloseTag = '</scr' + 'ipt>'
const styleCloseTag = '</sty' + 'le>'
const { isTabletLandscape } = useViewportProfile()

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
    const response = await listAssignments({ kidId: session.id })
    assignment.value = response.items.find((item) => String(item.id) === route.params.assignmentId)
    if (!assignment.value) {
      error.value = 'Activity assignment not found.'
    }
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
  const response = await updateAssignment(assignment.value.id, { progress: 100, status: 'completed' })
  assignment.value = response.item
}

function goBack() {
  router.push('/kids')
}
</script>
