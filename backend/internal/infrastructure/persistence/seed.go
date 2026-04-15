package persistence

import (
	"dangde-world/backend/internal/domain/activity"
	"dangde-world/backend/internal/domain/assignment"
	"dangde-world/backend/internal/domain/category"
	"dangde-world/backend/internal/domain/user"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	// Truncate tables to start fresh
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&userModel{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&categoryModel{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&activityModel{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&assignmentModel{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&activityDataModel{})

	users := NewUserRepository(db)
	categories := NewCategoryRepository(db)
	activities := NewActivityRepository(db)
	assignments := NewAssignmentRepository(db)

	if err := users.CreateBatch([]user.User{
		{ID: 1, Name: "Alya Admin", Role: user.RoleAdmin, Avatar: "🛠️", PreferredLang: "en"},
		{ID: 2, Name: "Prita Parent", Role: user.RoleParent, Avatar: "🌟", PreferredLang: "id"},
		{ID: 3, Name: "Nino", Role: user.RoleKid, Avatar: "🦁", PreferredLang: "en", ParentID: uintPtr(2)},
		{ID: 4, Name: "Sasa", Role: user.RoleKid, Avatar: "🐼", PreferredLang: "id", ParentID: uintPtr(2)},
	}); err != nil {
		return err
	}

	seedCategories := []category.Category{
		{ID: 1, Name: "Learning", Slug: "learning", Type: "learning"},
		{ID: 2, Name: "Playing", Slug: "playing", Type: "playing"},
		{ID: 3, Name: "Math", Slug: "math", Type: "learning", ParentID: uintPtr(1)},
		{ID: 4, Name: "Alphabets", Slug: "alphabets", Type: "learning", ParentID: uintPtr(1)},
		{ID: 5, Name: "Arabic Reading", Slug: "arabic-reading", Type: "learning", ParentID: uintPtr(1)},
		{ID: 6, Name: "Clock Reading", Slug: "clock-reading", Type: "learning", ParentID: uintPtr(1)},
	}
	if err := categories.CreateBatch(seedCategories); err != nil {
		return err
	}

	seedActivities := []activity.Activity{
		{
			ID:          1,
			Title:       "Alphabet Echo",
			Description: "Learn letters and say simple English and Bahasa Indonesia sounds out loud.",
			Language:    "English + Bahasa Indonesia",
			Difficulty:  "easy",
			AgeGroup:    "4-6",
			CategoryID:  4,
			Prompt:      "Tap a letter, hear it, and repeat it. Try saying 'Ba Ca = Baca'.",
			Icon:        "🔤",
			HTMLCode:    "<div class='activity-root'><h1>Alphabet Echo</h1><p id='word'>Tap a letter to build a word.</p><div class='letter-grid'><button data-letter='B'>B</button><button data-letter='A'>A</button><button data-letter='C'>C</button><button data-letter='A'>A</button></div><div class='actions'><button id='save-word'>Save Word</button><button id='load-word'>Load Word</button></div></div>",
			CSSCode:     ".activity-root{font-family:'Avenir Next',sans-serif;padding:24px;min-height:100vh;background:linear-gradient(180deg,#fff6d8,#dff8ff);color:#17324d}.letter-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:16px;margin:24px 0}.letter-grid button,#save-word,#load-word{min-height:72px;border:0;border-radius:24px;font-size:28px;font-weight:800;background:#fff;box-shadow:0 16px 30px rgba(0,0,0,.08)}.actions{display:grid;grid-template-columns:1fr 1fr;gap:16px}",
			JSCode:      "const letters=[];const word=document.getElementById('word');document.querySelectorAll('[data-letter]').forEach((button)=>{button.addEventListener('click',()=>{letters.push(button.dataset.letter);word.textContent=letters.join('')||'Tap a letter to build a word.';});});document.getElementById('save-word').addEventListener('click',async()=>{await window.DangdeAPI.save({key:'latest-word',value:{letters,word:letters.join('')}});await window.DangdeAPI.setProgress({progress:80,status:'in_progress'});word.textContent=`Saved: ${letters.join('')}`;});document.getElementById('load-word').addEventListener('click',async()=>{const rows=await window.DangdeAPI.query({key:'latest-word'});const latest=rows[0];if(!latest){word.textContent='No saved word yet.';return;}letters.splice(0,letters.length,...(latest.value.letters||[]));word.textContent=`Loaded: ${latest.value.word||''}`;});",
		},
		{
			ID:          2,
			Title:       "Count the Stars",
			Description: "Practice simple addition with bright visual counters.",
			Language:    "English",
			Difficulty:  "easy",
			AgeGroup:    "4-7",
			CategoryID:  3,
			Prompt:      "How many stars do you have when 2 join 3?",
			Icon:        "➕",
			HTMLCode:    "<div class='activity-root'><h1>Count the Stars</h1><p>Use the buttons to answer 2 + 3.</p><div class='answer-row'><button data-answer='4'>4</button><button data-answer='5'>5</button><button data-answer='6'>6</button></div><p id='result'></p></div>",
			CSSCode:     ".activity-root{font-family:'Avenir Next',sans-serif;padding:24px;min-height:100vh;background:linear-gradient(180deg,#eef9ff,#fff4ce)}.answer-row{display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-top:24px}.answer-row button{min-height:88px;border:0;border-radius:24px;font-size:30px;font-weight:800;background:#fff;box-shadow:0 16px 30px rgba(0,0,0,.08)}#result{margin-top:24px;font-size:24px;font-weight:700}",
			JSCode:      "const result=document.getElementById('result');document.querySelectorAll('[data-answer]').forEach((button)=>{button.addEventListener('click',async()=>{const correct=button.dataset.answer==='5';result.textContent=correct?'Correct!':'Try again!';await window.DangdeAPI.save({key:'last-answer',value:{answer:button.dataset.answer,correct}});await window.DangdeAPI.setProgress({progress:correct?100:40,status:correct?'completed':'in_progress'});});});",
		},
		{
			ID:          3,
			Title:       "Arabic Letter Path",
			Description: "Follow the shapes of basic Arabic characters and hear the sounds.",
			Language:    "Arabic",
			Difficulty:  "medium",
			AgeGroup:    "5-8",
			CategoryID:  5,
			Prompt:      "Trace the character and match it to the right sound.",
			Icon:        "📖",
			HTMLCode:    "<div class='activity-root'><h1>Arabic Letter Path</h1><p>Learn the shape and sound of the letter ب</p><div class='card'>ب</div><button id='remember'>I remember this</button><p id='message'></p></div>",
			CSSCode:     ".activity-root{font-family:'Avenir Next',sans-serif;padding:24px;min-height:100vh;background:linear-gradient(180deg,#fff7ef,#eef8ff);text-align:center}.card{display:grid;place-items:center;height:240px;border-radius:32px;background:#fff;font-size:120px;font-weight:800;margin:24px 0;box-shadow:0 16px 30px rgba(0,0,0,.08)}button{min-height:72px;padding:0 24px;border:0;border-radius:24px;background:#16b7d6;color:#fff;font-size:24px;font-weight:800}",
			JSCode:      "document.getElementById('remember').addEventListener('click',async()=>{await window.DangdeAPI.save({key:'arabic-letter',value:{letter:'ب',remembered:true}});await window.DangdeAPI.setProgress({progress:100,status:'completed'});document.getElementById('message').textContent='Great work!';});",
		},
		{
			ID:          4,
			Title:       "Tell the Time",
			Description: "Understand hours and half-hours by reading colorful clocks.",
			Language:    "English",
			Difficulty:  "medium",
			AgeGroup:    "5-8",
			CategoryID:  6,
			Prompt:      "Move the hands and say what time the clock shows.",
			Icon:        "🕒",
			HTMLCode:    "<div class='activity-root'><h1>Tell the Time</h1><p>What time is shown?</p><div class='clock-card'>3:00</div><div class='answer-row'><button data-answer='3:00'>3:00</button><button data-answer='4:00'>4:00</button></div><p id='result'></p></div>",
			CSSCode:     ".activity-root{font-family:'Avenir Next',sans-serif;padding:24px;min-height:100vh;background:linear-gradient(180deg,#f4f9ff,#fff2d5)}.clock-card{display:grid;place-items:center;height:240px;border-radius:32px;background:#fff;font-size:72px;font-weight:800;box-shadow:0 16px 30px rgba(0,0,0,.08)}.answer-row{display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-top:24px}.answer-row button{min-height:80px;border:0;border-radius:24px;background:#fff;font-size:26px;font-weight:800;box-shadow:0 16px 30px rgba(0,0,0,.08)}#result{margin-top:20px;font-size:24px;font-weight:700}",
			JSCode:      "const result=document.getElementById('result');document.querySelectorAll('[data-answer]').forEach((button)=>{button.addEventListener('click',async()=>{const correct=button.dataset.answer==='3:00';result.textContent=correct?'Yes, it is 3:00!':'Not yet, look at the short hand.';await window.DangdeAPI.save({key:'clock-answer',value:{answer:button.dataset.answer,correct}});await window.DangdeAPI.setProgress({progress:correct?100:50,status:correct?'completed':'in_progress'});});});",
		},
	}
	if err := activities.CreateBatch(seedActivities); err != nil {
		return err
	}

	seedAssignments := []assignment.Assignment{
		{ID: 1, ParentID: 2, KidID: 3, ActivityID: 1, Status: "in_progress", Progress: 60},
		{ID: 2, ParentID: 2, KidID: 3, ActivityID: 2, Status: "assigned", Progress: 10},
		{ID: 3, ParentID: 2, KidID: 4, ActivityID: 4, Status: "completed", Progress: 100},
	}
	return assignments.CreateBatch(seedAssignments)
}

func uintPtr(value uint) *uint {
	return &value
}

func seedCategoriesToModels(items []category.Category) []categoryModel {
	models := make([]categoryModel, 0, len(items))
	for _, item := range items {
		models = append(models, toCategoryModel(item))
	}
	return models
}

func seedActivitiesToModels(items []activity.Activity) []activityModel {
	models := make([]activityModel, 0, len(items))
	for _, item := range items {
		models = append(models, toActivityModel(item))
	}
	return models
}

func seedAssignmentsToModels(items []assignment.Assignment) []assignmentModel {
	models := make([]assignmentModel, 0, len(items))
	for _, item := range items {
		models = append(models, toAssignmentModel(item))
	}
	return models
}
