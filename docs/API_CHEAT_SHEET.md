# DangDe! Activity API - Quick Reference Cheat Sheet

## Context
```javascript
window.DangdeAPI.context.activityId    // Activity ID number
window.DangdeAPI.context.assignmentId  // Assignment ID number
window.DangdeAPI.context.kidId         // Kid ID number
```

## Save Data
```javascript
await window.DangdeAPI.save({ key: 'score', value: 42 })
```

## Query Data
```javascript
const records = await window.DangdeAPI.query()              // Get all
const records = await window.DangdeAPI.query({ key: 'score' })  // Get one
```

## Update Data
```javascript
await window.DangdeAPI.update({ id: 'rec-id', value: 100 })
await window.DangdeAPI.update({ key: 'score', value: 100 })
```

## Delete Data
```javascript
await window.DangdeAPI.delete({ id: 'rec-id' })
await window.DangdeAPI.delete({ key: 'score' })
```

## Update Progress
```javascript
await window.DangdeAPI.setProgress({ progress: 50, status: 'in_progress' })
await window.DangdeAPI.setProgress({ progress: 100 })  // Auto-completes
```

---

## Common Pattern: Initialize Activity

```javascript
(async () => {
  // Check if already initialized
  const existing = await window.DangdeAPI.query({ key: 'initialized' })
  if (existing.length > 0) return
  
  // Initialize
  await window.DangdeAPI.save({
    key: 'game-state',
    value: { score: 0, level: 1 }
  })
})()
```

## Common Pattern: Load, Update, Save

```javascript
async function incrementScore(points) {
  // Load
  const records = await window.DangdeAPI.query({ key: 'game-state' })
  const state = records[0].value
  
  // Update
  state.score += points
  
  // Save
  await window.DangdeAPI.update({
    id: records[0].id,
    value: state
  })
  
  // Report progress
  await window.DangdeAPI.setProgress({ progress: state.score })
}
```

## Common Pattern: Error Handling

```javascript
try {
  await window.DangdeAPI.save({ key: 'score', value: 42 })
} catch (error) {
  console.error('Failed:', error.message)
}
```

---

## HTML/CSS/JS File Placeholders

### index.html
```html
<div class="game-shell">
  <p>Your activity UI goes here</p>
</div>
```

### styles.css
```css
.game-shell {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 2rem;
}
```

### app.js
```javascript
(async () => {
  // Your activity logic here
  const context = window.DangdeAPI.context
  console.log('Activity started:', context)
})()
```

---

## What Can Activities Do?
✅ DOM manipulation
✅ Canvas / WebGL
✅ Audio / Video
✅ Save/load data
✅ Track progress
✅ Handle clicks, input, touch

## What Can't Activities Do?
❌ Modify parent page
❌ Use localStorage
❌ Make CORS requests
❌ Navigate parent window
❌ Access cookies

---

## Tips
- Always `await` API calls
- Always wrap in `try-catch`
- Use descriptive key names: `game-state`, not `x`
- Call `setProgress()` regularly
- Test on mobile/tablet too
- Check browser console for errors

---

For full documentation: `docs/API_REFERENCE.md`
