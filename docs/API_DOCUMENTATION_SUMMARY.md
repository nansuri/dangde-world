# DangDe! World - Activity Development Documentation

## 📚 Overview

This document summarizes the activity development capabilities and API available to activity creators in the DangDe! World admin studio.

---

## 🎯 Quick Start

### For Activity Creators

1. **Log in** as an Admin in DangDe! World
2. **Navigate** to Admin Studio → Curriculum Management
3. **Choose** a template from the quick-start options, or create a new activity from scratch
4. **Use** the built-in API Reference (collapsible panel in the activity editor) to understand available functions
5. **Test** your activity by saving and launching it from the kids dashboard

### Files to Review

- **API Reference**: `docs/API_REFERENCE.md` - Complete API documentation with examples
- **This Document**: `docs/API_DOCUMENTATION_SUMMARY.md` - Quick overview and best practices

---

## 🎨 Quick-Start Templates

The admin studio provides 5 ready-to-use templates to get started quickly:

### 1. 🎮 Score Game
A simple game that tracks score and updates progress. Perfect for:
- Implementing point-based mechanics
- Learning basic save/query patterns
- Demonstrating progress tracking

### 2. ❓ Quiz
Multiple choice quiz with instant feedback. Perfect for:
- Educational assessments
- Learning if/else logic
- Displaying conditional UI

### 3. 🎨 Drawing
Simple canvas-based drawing app. Perfect for:
- Learning canvas APIs
- Understanding file/image storage
- Creative activities

### 4. 🧠 Memory Game
Classic flip-card memory game. Perfect for:
- Implementing complex game loops
- Learning event handling
- Building interactive games

### 5. 📝 [Create from Scratch]
Blank template with just the basics. Perfect for:
- Complete control and customization
- Implementing unique ideas

---

## 💻 Activity IDE Interface

### Sidebar (Left)
- **Workspace Panel**: Shows activity metadata
- **File Tabs**: Switch between HTML, CSS, and JavaScript
- **Metadata Form**: Edit title, description, category, difficulty, age group, icon
- **API Reference**: Collapsible panel with full API documentation
- **Save Button**: Submit your activity

### Editor (Right)
- **Tab Bar**: Visual tabs for each file with character counts
- **Line Numbers**: Gutter with line numbering
- **Code Editor**: Simple textarea-based editor for HTML, CSS, JS
- **Syntax**: No real-time syntax highlighting (yet)

### How to Use

1. **Enter metadata** in the left sidebar
2. **Click a file tab** (index.html, styles.css, app.js) to switch files
3. **Edit code** in the textarea
4. **Reference the API** using the collapsible API Reference panel
5. **Save** when ready
6. **Test** by running the activity as a kid

---

## 🔌 The DangdeAPI Bridge

### What is It?

Activities run in a **sandboxed iframe** and communicate with the parent app through `window.DangdeAPI`. This is your only way to interact with the parent app and persist data.

### Why Sandbox?

- **Security**: Activity code can't access parent's DOM, localStorage, or cookies
- **Isolation**: Activity styles and scripts don't pollute parent app
- **Safety**: Prevents malicious code from harming the platform

### What Can Activities Do?

✅ **Data Storage**: Save/load game state via `DangdeAPI.save()` and `query()`  
✅ **Progress Tracking**: Update assignment progress via `DangdeAPI.setProgress()`  
✅ **Context Access**: Read activity/kid/assignment info via `DangdeAPI.context`  
✅ **DOM Manipulation**: Full control over iframe's DOM  
✅ **Canvas/WebGL**: Drawing, animations, games inside the iframe  
✅ **Audio/Video**: Embed media inside the iframe  

### What Can't Activities Do?

❌ **No parent DOM access** - Can't modify the parent app's page  
❌ **No localStorage** - Must use DangdeAPI for persistence  
❌ **No arbitrary HTTP requests** - Can only fetch from `/api/*` endpoints  
❌ **No navigation** - Can't change the parent page's URL  
❌ **No same-origin cookies** - Each activity is isolated  

---

## 📖 Core API Methods

### Context
```javascript
const { activityId, assignmentId, kidId } = window.DangdeAPI.context
```

### Save Data
```javascript
await window.DangdeAPI.save({
  key: 'my-score',
  value: { score: 42, level: 3 }
})
```

### Query Data
```javascript
const records = await window.DangdeAPI.query({ key: 'my-score' })
```

### Update Data
```javascript
await window.DangdeAPI.update({
  key: 'my-score',
  value: { score: 100 }
})
```

### Delete Data
```javascript
await window.DangdeAPI.delete({ key: 'my-score' })
```

### Update Progress
```javascript
await window.DangdeAPI.setProgress({
  progress: 75,
  status: 'in_progress'
})
```

**→ For full details, see `docs/API_REFERENCE.md`**

---

## 🛠️ Common Implementation Patterns

### Pattern 1: Initialize on Load

```javascript
(async () => {
  // Check if already initialized
  const existing = await window.DangdeAPI.query({ key: 'initialized' })
  
  if (existing.length === 0) {
    // First time
    await window.DangdeAPI.save({
      key: 'initialized',
      value: { startedAt: new Date().toISOString() }
    })
    
    // Initialize game state
    await window.DangdeAPI.save({
      key: 'game-state',
      value: { score: 0, level: 1, completed: false }
    })
  }
})()
```

### Pattern 2: Update and Save State

```javascript
async function updateGameState(updates) {
  try {
    // Load current state
    const records = await window.DangdeAPI.query({ key: 'game-state' })
    const current = records[0]?.value || {}
    
    // Merge updates
    const newState = { ...current, ...updates }
    
    // Save
    await window.DangdeAPI.update({
      id: records[0]?.id,
      value: newState
    })
    
    return newState
  } catch (error) {
    console.error('Failed to update state:', error)
  }
}
```

### Pattern 3: Error Handling

```javascript
async function safeApiCall(fn) {
  try {
    return await fn()
  } catch (error) {
    console.error('API Error:', error.message)
    // Gracefully degrade
    return null
  }
}

// Usage
const saved = await safeApiCall(() => 
  window.DangdeAPI.save({ key: 'score', value: 10 })
)
```

### Pattern 4: Track Progress

```javascript
async function setProgressPercent(percent) {
  const status = percent >= 100 ? 'completed' : 'in_progress'
  await window.DangdeAPI.setProgress({
    progress: Math.min(100, percent),
    status
  })
}
```

---

## 🎯 Best Practices

### Do's ✅

- **Do** wrap API calls in `try-catch`
- **Do** use meaningful key names: `game-state`, `player-score`, `level-progress`
- **Do** keep data sizes small (< 1MB per key)
- **Do** update progress regularly so parents see engagement
- **Do** validate data loaded from storage
- **Do** clean up temporary data when done

### Don'ts ❌

- **Don't** forget to `await` async calls
- **Don't** store large files/images in the database
- **Don't** assume `window.DangdeAPI` exists outside the iframe
- **Don't** make arbitrary CORS requests (use `/api/*` only)
- **Don't** store sensitive data unencrypted
- **Don't** ignore errors - always handle failures gracefully

---

## 📱 Responsive Design

Activities should be responsive and work on all screen sizes:

- **Desktop**: Full-size browser window
- **Tablet**: iPad in both portrait and landscape
- **Mobile**: Touch-friendly sizing

### Tips

- Use `max-width` for desktop
- Use viewport-relative units (`vh`, `vw`, `svh`)
- Make buttons large enough for touch (minimum 44x44px)
- Test in landscape orientation on tablets
- Use flexbox/grid for responsive layouts

---

## 🧪 Testing Your Activity

### Step 1: Create Activity
1. Log in as Admin
2. Go to Curriculum Management
3. Fill in metadata and code
4. Click "Add Activity"

### Step 2: Assign to Kid
1. Log in as Parent
2. Go to Assignments section
3. Create assignment with your activity
4. Assign to a kid

### Step 3: Test as Kid
1. Log out / Log in as Kid
2. Navigate to the activity
3. Launch the activity player
4. Test all interactions
5. Verify progress updates
6. Check parent can see progress

### Step 4: Debug
1. Open browser DevTools (F12)
2. Go to Console tab
3. Look for any errors
4. Check Network tab for failed API calls
5. Use `console.log()` to trace execution

---

## 🔍 API Reference Location

The complete API documentation is available in two places:

### 1. **In-App Reference** (While Creating Activity)
- Collapsible "📚 API Reference" panel in the activity editor
- Searchable and organized by method
- Shows parameters, returns, and examples
- Available while editing

### 2. **Offline Reference** (File-based)
- Location: `docs/API_REFERENCE.md`
- Complete with TypeScript definitions
- Troubleshooting section
- Best practices guide
- Complete examples

---

## 🚀 Advanced Topics

### Custom Data Types

Store any JSON-serializable data:

```javascript
// Numbers
await window.DangdeAPI.save({ key: 'score', value: 42 })

// Objects
await window.DangdeAPI.save({ 
  key: 'state', 
  value: { x: 10, y: 20, rotation: 45 } 
})

// Arrays
await window.DangdeAPI.save({ 
  key: 'items', 
  value: ['sword', 'shield', 'potion'] 
})

// Nested structures
await window.DangdeAPI.save({ 
  key: 'player', 
  value: {
    name: 'Alice',
    stats: { health: 100, mana: 50 },
    inventory: ['key', 'map']
  }
})
```

### Conditional Logic

```javascript
const records = await window.DangdeAPI.query({ key: 'game-state' })
const gameState = records[0]?.value

if (gameState.level >= 5) {
  // Unlock special content
} else if (gameState.score > 1000) {
  // Show bonus challenge
}
```

### Auto-save

```javascript
setInterval(async () => {
  const gameState = getCurrentGameState()
  try {
    await window.DangdeAPI.save({
      key: 'auto-save',
      value: gameState
    })
  } catch (error) {
    console.warn('Auto-save failed:', error)
  }
}, 30000) // Every 30 seconds
```

---

## 📞 Support & Resources

### If You Get Stuck

1. **Check the API Reference** - Most questions are answered there
2. **Look at Templates** - All 5 templates demonstrate different patterns
3. **Check Console** - DevTools console shows errors and hints
4. **Review Examples** - `docs/API_REFERENCE.md` has many code examples

### Common Issues

| Issue | Solution |
|-------|----------|
| "window.DangdeAPI is undefined" | API only exists inside activity iframe during runtime |
| Data not saving | Check browser console for errors; verify internet connection |
| Activity not showing | Check that metadata is filled in and activity is assigned |
| Styles not applying | Verify CSS is in the styles.css file (not index.html) |
| Script not running | Check for JavaScript errors in console; verify syntax |

---

## 📊 What's Next?

### Planned Improvements

- ✨ **Live Preview**: See changes in real-time while editing
- ✨ **Syntax Highlighting**: Better code editor with color coding
- ✨ **Linting**: Real-time error detection as you type
- ✨ **Component Library**: Pre-built reusable UI components
- ✨ **Testing Tools**: Built-in testing framework for activities
- ✨ **Version Control**: Save and restore activity versions

### Future API Additions

- 🔮 **Media Upload**: Upload images/audio directly from activity
- 🔮 **Analytics**: Track detailed user behavior
- 🔮 **Multiplayer**: Activities with multiple kids
- 🔮 **Real-time Sync**: Live data sync between devices

---

## 📝 Version Information

- **Documentation Version**: 1.0.0
- **API Version**: 1.0.0
- **Last Updated**: 2026-04-15
- **Status**: Stable and ready for production

---

## 🎓 Learning Resources

### For Beginners
1. Start with a template
2. Modify the template to learn
3. Read the API Reference
4. Create your first original activity

### For Intermediate Developers
1. Study the Score Game template
2. Understand save/query/update patterns
3. Implement custom game logic
4. Add progress tracking

### For Advanced Developers
1. Build complex interactive experiences
2. Optimize performance and bundle size
3. Create reusable patterns
4. Push the boundaries of what's possible

---

## 💡 Tips & Tricks

### Tip 1: Use Good Key Names
```javascript
// ✅ Good
await window.DangdeAPI.save({ key: 'player-score', value: 42 })

// ❌ Bad
await window.DangdeAPI.save({ key: 'x', value: 42 })
```

### Tip 2: Batch Operations
```javascript
// ❌ Slow: Many individual saves
for (let i = 0; i < 100; i++) {
  await window.DangdeAPI.save({ key: `item-${i}`, value: data[i] })
}

// ✅ Fast: One save with array
await window.DangdeAPI.save({ 
  key: 'items', 
  value: data 
})
```

### Tip 3: Check Before Update
```javascript
// Always check if data exists before updating
const existing = await window.DangdeAPI.query({ key: 'score' })
if (existing.length > 0) {
  await window.DangdeAPI.update({
    id: existing[0].id,
    value: newValue
  })
} else {
  await window.DangdeAPI.save({
    key: 'score',
    value: newValue
  })
}
```

### Tip 4: Graceful Degradation
```javascript
// If API fails, still show UI
try {
  await window.DangdeAPI.save({ key: 'score', value: 100 })
} catch (error) {
  console.warn('Could not save to server')
  // Continue anyway - UI still works
}
```

---

## 📄 Files in This Documentation Suite

```
docs/
├── API_REFERENCE.md                    # Complete API reference
├── API_DOCUMENTATION_SUMMARY.md        # This file
└── plan/
    ├── 2026-04-15-initial-webapp-plan.md
    ├── 2026-04-15-ipad-safari-compatibility.md
    └── 2026-04-15-custom-activity-runtime.md
```

Happy creating! 🎉
