<template>
  <div class="api-reference" :class="{ expanded: isExpanded }">
    <button class="api-reference-toggle" type="button" @click="isExpanded = !isExpanded">
      <span>📚</span>
      <span>API Reference</span>
      <span class="toggle-indicator">{{ isExpanded ? '▼' : '▶' }}</span>
    </button>

    <div v-if="isExpanded" class="api-reference-content">
      <p class="api-intro">
        Activities run in a sandboxed iframe and communicate with the parent app through the
        <code>window.DangdeAPI</code> bridge. All methods are async and must be awaited.
      </p>

      <section class="api-section">
        <h4>🔗 Context</h4>
        <p class="api-description">
          Access the current activity context with <code>window.DangdeAPI.context</code>:
        </p>
        <div class="api-code-block">
          <code>
            window.DangdeAPI.context.activityId<br />
            window.DangdeAPI.context.assignmentId<br />
            window.DangdeAPI.context.kidId
          </code>
        </div>
      </section>

      <section class="api-section">
        <h4>💾 save(payload)</h4>
        <p class="api-description">
          Save or upsert a key-value pair. If the key exists, it will be updated. Returns the saved
          record.
        </p>
        <div class="api-code-block">
          <code>
            await window.DangdeAPI.save(&#123;<br />
            &nbsp;&nbsp;key: 'my-score',<br />
            &nbsp;&nbsp;value: &#123; score: 42, level: 3 &#125;<br />
            &#125;)
          </code>
        </div>
        <p class="api-note"><strong>Required:</strong> <code>key</code> (string)</p>
        <p class="api-note"><strong>Value:</strong> Any JSON-serializable object or primitive</p>
      </section>

      <section class="api-section">
        <h4>🔍 query(payload)</h4>
        <p class="api-description">
          Query stored data for this activity. Optionally filter by key. Returns an array of records.
        </p>
        <div class="api-code-block">
          <code>
            // Get all stored data<br />
            const all = await window.DangdeAPI.query()<br />
            <br />
            // Get specific key<br />
            const data = await window.DangdeAPI.query(&#123; key: 'my-score' &#125;)
          </code>
        </div>
        <p class="api-note"><strong>Optional:</strong> <code>key</code> (string) to filter results</p>
      </section>

      <section class="api-section">
        <h4>✏️ update(payload)</h4>
        <p class="api-description">
          Update an existing record. Can update by <code>id</code> or <code>key</code>. Returns the
          updated record.
        </p>
        <div class="api-code-block">
          <code>
            // By ID<br />
            await window.DangdeAPI.update(&#123;<br />
            &nbsp;&nbsp;id: 'record-id',<br />
            &nbsp;&nbsp;key: 'my-score',<br />
            &nbsp;&nbsp;value: &#123; score: 50 &#125;<br />
            &#125;)<br />
            <br />
            // By Key<br />
            await window.DangdeAPI.update(&#123;<br />
            &nbsp;&nbsp;key: 'my-score',<br />
            &nbsp;&nbsp;value: &#123; score: 50 &#125;<br />
            &#125;)
          </code>
        </div>
        <p class="api-note"><strong>Required:</strong> <code>id</code> or <code>key</code></p>
      </section>

      <section class="api-section">
        <h4>🗑️ delete(payload)</h4>
        <p class="api-description">
          Delete a record or all records with a key. Returns count of deleted records.
        </p>
        <div class="api-code-block">
          <code>
            // By ID<br />
            await window.DangdeAPI.delete(&#123; id: 'record-id' &#125;)<br />
            <br />
            // By Key<br />
            await window.DangdeAPI.delete(&#123; key: 'my-score' &#125;)
          </code>
        </div>
        <p class="api-note"><strong>Required:</strong> <code>id</code> or <code>key</code></p>
      </section>

      <section class="api-section">
        <h4>📊 setProgress(payload)</h4>
        <p class="api-description">
          Update the assignment progress and status. Returns the updated assignment.
        </p>
        <div class="api-code-block">
          <code>
            await window.DangdeAPI.setProgress(&#123;<br />
            &nbsp;&nbsp;progress: 75,<br />
            &nbsp;&nbsp;status: 'in_progress'<br />
            &#125;)
          </code>
        </div>
        <p class="api-note"><strong>Optional:</strong> <code>progress</code> (0-100)</p>
        <p class="api-note">
          <strong>Optional:</strong> <code>status</code> ('in_progress', 'completed')
        </p>
        <p class="api-note">
          If progress reaches 100 and status is not set, it will default to 'completed'
        </p>
      </section>

      <section class="api-section">
        <h4>⚠️ Error Handling</h4>
        <p class="api-description">All API calls can throw errors. Wrap calls in try-catch:</p>
        <div class="api-code-block">
          <code>
            try &#123;<br />
            &nbsp;&nbsp;const result = await window.DangdeAPI.save(&#123; key: 'score', value: 10
            &#125;)<br />
            &#125; catch (error) &#123;<br />
            &nbsp;&nbsp;console.error('Failed to save:', error.message)<br />
            &#125;
          </code>
        </div>
      </section>

      <section class="api-section">
        <h4>📝 Complete Example</h4>
        <div class="api-code-block">
          <code>
            (async () => &#123;<br />
            &nbsp;&nbsp;const context = window.DangdeAPI.context<br />
            &nbsp;&nbsp;console.log('Activity ID:', context.activityId)<br />
            &nbsp;&nbsp;<br />
            &nbsp;&nbsp;// Save initial state<br />
            &nbsp;&nbsp;await window.DangdeAPI.save(&#123; key: 'game', value: &#123; score: 0,
            level: 1 &#125; &#125;)<br />
            &nbsp;&nbsp;<br />
            &nbsp;&nbsp;// Load state<br />
            &nbsp;&nbsp;const records = await window.DangdeAPI.query(&#123; key: 'game' &#125;)<br />
            &nbsp;&nbsp;console.log('Game state:', records[0].value)<br />
            &nbsp;&nbsp;<br />
            &nbsp;&nbsp;// Update progress<br />
            &nbsp;&nbsp;await window.DangdeAPI.setProgress(&#123; progress: 50, status:
            'in_progress' &#125;)<br />
            &#125;)()
          </code>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const isExpanded = ref(false)
</script>

<style scoped>
.api-reference {
  border-top: 1px solid var(--border);
  background: var(--surface);
}

.api-reference-toggle {
  width: 100%;
  padding: 0.75rem;
  background: none;
  border: none;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-weight: 600;
  color: var(--ink);
  text-align: left;
  cursor: pointer;
  transition: background-color 200ms ease;
}

.api-reference-toggle:hover {
  background: rgba(0, 0, 0, 0.04);
}

.toggle-indicator {
  margin-left: auto;
  font-size: 0.75rem;
}

.api-reference-content {
  padding: 1rem;
  max-height: 500px;
  overflow-y: auto;
  font-size: 0.875rem;
  line-height: 1.6;
}

.api-intro {
  margin: 0 0 1rem 0;
  padding: 0.75rem;
  background: rgba(55, 188, 255, 0.1);
  border-left: 3px solid var(--secondary);
  border-radius: 4px;
  color: var(--ink);
}

.api-section {
  margin-bottom: 1.5rem;
}

.api-section h4 {
  margin: 0 0 0.5rem 0;
  font-size: 0.95rem;
  color: var(--primary);
}

.api-description {
  margin: 0.5rem 0;
  color: var(--muted);
}

.api-code-block {
  background: rgba(23, 50, 77, 0.08);
  border: 1px solid var(--border);
  border-radius: 4px;
  padding: 0.75rem;
  margin: 0.75rem 0;
  overflow-x: auto;
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 0.8rem;
  line-height: 1.4;
}

.api-code-block code {
  color: var(--ink);
  white-space: pre-wrap;
  word-break: break-word;
}

code {
  background: rgba(23, 50, 77, 0.08);
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 0.85em;
}

.api-note {
  margin: 0.5rem 0 0 0;
  padding: 0.5rem;
  background: rgba(255, 215, 92, 0.1);
  border-left: 3px solid var(--primary);
  font-size: 0.8rem;
  color: var(--ink);
}

@media (max-height: 800px) {
  .api-reference-content {
    max-height: 300px;
  }
}
</style>
