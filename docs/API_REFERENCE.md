# DangDe! World Activity Runtime API Reference

## Overview

Activities in DangDe! World run inside a sandboxed iframe and communicate with the parent application through the `window.DangdeAPI` bridge. This API provides a secure, asynchronous interface for activities to:

- Store and retrieve persistent data
- Update progress and status
- Access runtime context information

All API methods are **asynchronous** and must be awaited using `async/await` or `.then()` syntax.

---

## Context

Every activity has access to the current execution context:

```javascript
window.DangdeAPI.context
```

### Properties

| Property | Type | Description |
|----------|------|-------------|
| `activityId` | number | The ID of the current activity |
| `assignmentId` | number | The ID of the current assignment |
| `kidId` | number | The ID of the kid playing this activity |

### Example

```javascript
const { activityId, assignmentId, kidId } = window.DangdeAPI.context;
console.log(`Activity ${activityId} for kid ${kidId}`);
```

---

## API Methods

### save(payload)

Save or upsert a key-value pair. If a record with the same key already exists, it will be updated. Returns the saved record object.

**Parameters:**

```typescript
{
  key: string,          // Required: unique identifier for this data
  value: any            // Required: any JSON-serializable value or object
}
```

**Returns:**

```typescript
{
  id: string,
  key: string,
  value: any,
  activityId: number,
  assignmentId: number,
  kidId: number,
  createdAt: string,
  updatedAt: string
}
```

**Example:**

```javascript
// Save a simple value
await window.DangdeAPI.save({
  key: 'score',
  value: 42
});

// Save a complex object
await window.DangdeAPI.save({
  key: 'game-state',
  value: {
    level: 3,
    score: 150,
    items: ['sword', 'shield'],
    position: { x: 10, y: 20 }
  }
});

// Upsert (save will update if key exists)
const result = await window.DangdeAPI.save({
  key: 'score',
  value: 85  // This will replace 42
});
```

---

### query(payload?)

Query stored data for this activity. Optionally filter by a specific key. Returns an array of matching records.

**Parameters:**

```typescript
{
  key?: string    // Optional: filter by specific key
}
```

**Returns:**

```typescript
[
  {
    id: string,
    key: string,
    value: any,
    activityId: number,
    assignmentId: number,
    kidId: number,
    createdAt: string,
    updatedAt: string
  },
  // ... more records
]
```

**Example:**

```javascript
// Get all stored data for this activity
const allData = await window.DangdeAPI.query();

// Filter by specific key
const scoreRecords = await window.DangdeAPI.query({ key: 'score' });

// Access the first result
if (scoreRecords.length > 0) {
  const currentScore = scoreRecords[0].value;
  console.log('Current score:', currentScore);
}
```

---

### update(payload)

Update an existing record. Can update by record ID or key. Returns the updated record.

**Parameters:**

```typescript
{
  id?: string,          // Either id or key must be provided
  key?: string,         // Either id or key must be provided
  value: any            // Required: new value to set
}
```

**Returns:**

```typescript
{
  id: string,
  key: string,
  value: any,
  activityId: number,
  assignmentId: number,
  kidId: number,
  createdAt: string,
  updatedAt: string
}
```

**Example:**

```javascript
// Update by record ID
const updated = await window.DangdeAPI.update({
  id: 'some-record-id',
  value: { score: 200, completed: true }
});

// Update by key (updates first matching record)
const updated = await window.DangdeAPI.update({
  key: 'game-state',
  value: { level: 5, score: 500 }
});

// If you need to change the key too
const updated = await window.DangdeAPI.update({
  key: 'old-key',
  value: { /* new data */ }
  // Note: to change the key, use delete + save instead
});
```

---

### delete(payload)

Delete a record or all records with a specific key. Returns the count of deleted records.

**Parameters:**

```typescript
{
  id?: string,    // Either id or key must be provided
  key?: string    // Either id or key must be provided
}
```

**Returns:**

```typescript
{
  deleted: number  // Count of deleted records
}
```

**Example:**

```javascript
// Delete specific record by ID
const result = await window.DangdeAPI.delete({
  id: 'record-id-123'
});
console.log(`Deleted ${result.deleted} record(s)`);

// Delete all records with a key
const result = await window.DangdeAPI.delete({
  key: 'temporary-data'
});
console.log(`Deleted ${result.deleted} record(s)`);
```

---

### setProgress(payload)

Update the assignment progress and status. This is how activities report back to parents and admins how the kid is performing.

**Parameters:**

```typescript
{
  progress?: number,    // 0-100, optional
  status?: string       // 'in_progress' | 'completed', optional
}
```

**Returns:**

```typescript
{
  id: number,
  parentId: number,
  kidId: number,
  activityId: number,
  status: string,
  progress: number,
  createdAt: string,
  updatedAt: string
}
```

**Special Behavior:**

- If `progress` reaches 100 and `status` is not explicitly set, `status` will automatically be set to `'completed'`
- If `progress` is less than 100 and `status` is not set, it will be `'in_progress'`
- You can call this multiple times to incrementally update progress

**Example:**

```javascript
// Mark as started
await window.DangdeAPI.setProgress({
  progress: 10,
  status: 'in_progress'
});

// Update mid-activity
await window.DangdeAPI.setProgress({
  progress: 50
});

// Complete activity
await window.DangdeAPI.setProgress({
  progress: 100
  // status will auto-set to 'completed'
});

// Or explicitly set completed
await window.DangdeAPI.setProgress({
  progress: 100,
  status: 'completed'
});
```

---

## Error Handling

All API methods can throw errors. Common error scenarios:

- Network failures
- Invalid parameters
- Storage quota exceeded
- Authentication issues

**Best Practice:**

```javascript
try {
  const result = await window.DangdeAPI.save({
    key: 'score',
    value: 42
  });
  console.log('Success:', result);
} catch (error) {
  console.error('Failed to save:', error.message);
  // Gracefully handle the error
  // - Show user message
  // - Retry with exponential backoff
  // - Fall back to local storage
}
```

---

## Common Patterns

### Initialize Activity State

```javascript
async function initializeActivity() {
  // Check if activity has been started before
  const existing = await window.DangdeAPI.query({ key: 'initialized' });
  
  if (existing.length > 0) {
    console.log('Activity already started');
    return;
  }
  
  // First time, initialize
  await window.DangdeAPI.save({
    key: 'initialized',
    value: { startedAt: new Date().toISOString() }
  });
  
  await window.DangdeAPI.save({
    key: 'game-state',
    value: { score: 0, level: 1, lives: 3 }
  });
  
  await window.DangdeAPI.setProgress({
    progress: 0,
    status: 'in_progress'
  });
}
```

### Load and Update Game State

```javascript
async function updateScore(points) {
  try {
    // Load current state
    const records = await window.DangdeAPI.query({ key: 'game-state' });
    const currentState = records[0].value;
    
    // Update state
    const newScore = currentState.score + points;
    const newState = { ...currentState, score: newScore };
    
    // Save
    await window.DangdeAPI.update({
      id: records[0].id,
      value: newState
    });
    
    // Update progress if applicable
    const progressPercentage = Math.min(100, newScore / 500 * 100);
    await window.DangdeAPI.setProgress({ progress: Math.round(progressPercentage) });
    
    return newState;
  } catch (error) {
    console.error('Error updating score:', error);
  }
}
```

### Auto-Save Game State

```javascript
async function setupAutoSave() {
  setInterval(async () => {
    try {
      const gameState = getCurrentGameState(); // Your function
      await window.DangdeAPI.save({
        key: 'auto-save',
        value: gameState
      });
      console.log('Auto-saved');
    } catch (error) {
      console.warn('Auto-save failed:', error);
    }
  }, 30000); // Every 30 seconds
}
```

### Clear Activity Data (Reset)

```javascript
async function resetActivity() {
  try {
    // Query all data
    const allData = await window.DangdeAPI.query();
    
    // Delete each record
    for (const record of allData) {
      await window.DangdeAPI.delete({ id: record.id });
    }
    
    console.log('Activity reset');
  } catch (error) {
    console.error('Reset failed:', error);
  }
}
```

---

## Activity Sandbox Limitations

Activities run in an iframe with the following sandbox restrictions:

- **No `allow-same-origin`** - Activities cannot access parent cookies or local storage
- **No `allow-top-navigation`** - Activities cannot navigate the parent window
- **No external API access** - Activities cannot make arbitrary CORS requests
- **CSS/Script isolated** - Activity styles and scripts don't affect parent app

**Allowed Features:**

- Internal DOM manipulation
- Local iframe DOM storage (persists for session)
- All DangdeAPI methods
- Standard browser APIs (fetch is limited, fetch to relative URLs like `/api/*` works)

---

## TypeScript Support

For TypeScript projects, here's the API interface:

```typescript
interface DangdeAPIContext {
  activityId: number;
  assignmentId: number;
  kidId: number;
}

interface ActivityDataRecord {
  id: string;
  key: string;
  value: any;
  activityId: number;
  assignmentId: number;
  kidId: number;
  createdAt: string;
  updatedAt: string;
}

interface AssignmentRecord {
  id: number;
  parentId: number;
  kidId: number;
  activityId: number;
  status: 'pending' | 'in_progress' | 'completed';
  progress: number;
  createdAt: string;
  updatedAt: string;
}

interface DangdeAPI {
  context: DangdeAPIContext;
  save(payload: { key: string; value: any }): Promise<ActivityDataRecord>;
  query(payload?: { key?: string }): Promise<ActivityDataRecord[]>;
  update(payload: { id?: string; key?: string; value: any }): Promise<ActivityDataRecord>;
  delete(payload: { id?: string; key?: string }): Promise<{ deleted: number }>;
  setProgress(payload?: { progress?: number; status?: string }): Promise<AssignmentRecord>;
}

declare global {
  interface Window {
    DangdeAPI: DangdeAPI;
  }
}
```

---

## Best Practices

1. **Always use try-catch** - Wrap API calls in error handling
2. **Avoid storing large data** - Keep data sizes reasonable for performance
3. **Use meaningful keys** - Use descriptive keys like `game-state`, `score`, `level-progress`
4. **Batch operations** - When possible, save complex objects instead of many small saves
5. **Update progress regularly** - Call `setProgress()` periodically so parents see activity
6. **Validate data** - Don't trust data from storage; validate before using
7. **Clean up** - Delete temporary data when no longer needed

---

## Troubleshooting

### "window.DangdeAPI is undefined"

The API is only available inside the activity iframe during activity runtime, not during admin preview or IDE editing. This is expected.

### "save requires a key"

The `save()` method requires a non-empty `key` parameter. Example:

```javascript
// ❌ Wrong
await window.DangdeAPI.save({ value: 42 });

// ✅ Correct
await window.DangdeAPI.save({ key: 'score', value: 42 });
```

### "update requires id or key"

The `update()` method needs either an `id` (from a previous query) or a `key` to identify which record to update.

### Data not persisting

Check browser console for errors. Ensure:
- You're using `await`
- The kid is logged in
- The activity is running in a real assignment (not preview)
- Storage isn't full

---

## Version History

- **Current Version**: 1.0.0
- **Date**: 2026-04-15
- **Status**: Stable

The API is part of the DangDe! World platform and will maintain backward compatibility for activities. Future versions will add features without breaking existing activity code.
