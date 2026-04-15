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
			Description: "A 3-stage literacy quest: build words, identify sounds, and complete mini spelling checks.",
			Language:    "English + Bahasa Indonesia",
			Difficulty:  "medium",
			AgeGroup:    "4-7",
			CategoryID:  4,
			Prompt:      "Complete stage challenges by building words and matching missing letters. Reach all stars to finish.",
			Icon:        "🔤",
			HTMLCode: `<div class='activity-root'>
  <h1>Alphabet Echo Adventure</h1>
  <p id='stage-label'></p>
  <p id='instruction'></p>
  <div id='target' class='target-word'></div>
  <div id='letter-grid' class='letter-grid'></div>
  <p id='attempt' class='attempt-word'>Your word: </p>
  <div class='actions'>
    <button id='check-btn'>Check</button>
    <button id='reset-btn'>Reset</button>
    <button id='next-btn' disabled>Next Stage</button>
  </div>
  <div class='stats'>
    <span id='score'>Stars: 0</span>
    <span id='progress'>Progress: 0%</span>
  </div>
  <p id='feedback'></p>
</div>`,
			CSSCode: `.activity-root{font-family:'Avenir Next',sans-serif;padding:20px;min-height:100vh;background:linear-gradient(180deg,#fff6d8,#dff8ff);color:#17324d}
h1{margin:0 0 10px}
#instruction{margin:0 0 12px}
.target-word{font-size:34px;font-weight:900;letter-spacing:8px;padding:12px 14px;border-radius:18px;background:#fff;margin:10px 0}
.letter-grid{display:grid;grid-template-columns:repeat(5,minmax(0,1fr));gap:10px;margin:12px 0}
.letter-grid button,.actions button{min-height:58px;border:0;border-radius:18px;font-size:22px;font-weight:800;background:#fff;box-shadow:0 10px 20px rgba(0,0,0,.08)}
.actions{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:10px;margin-top:8px}
.attempt-word{font-size:22px;font-weight:800}
.stats{display:flex;gap:14px;margin-top:12px;font-weight:800}
#feedback{margin-top:12px;font-size:18px;font-weight:800}`,
			JSCode: `const stages=[{word:'BACA',hint:'Build the word BACA from the letters.',pool:['B','A','C','A','D','E','I']},{word:'BUKU',hint:'Build the word BUKU.',pool:['B','U','K','U','O','A','R']},{word:'MAIN',hint:'Build the word MAIN.',pool:['M','A','I','N','S','T','L']}];
let stageIndex=0;let picked=[];let score=0;let cleared=0;
const stageLabel=document.getElementById('stage-label');const instruction=document.getElementById('instruction');const target=document.getElementById('target');const grid=document.getElementById('letter-grid');const attempt=document.getElementById('attempt');const feedback=document.getElementById('feedback');const scoreEl=document.getElementById('score');const progressEl=document.getElementById('progress');const nextBtn=document.getElementById('next-btn');
function progressForCurrent(){return Math.min(100,Math.round((cleared/stages.length)*100));}
function renderStage(){const stage=stages[stageIndex];stageLabel.textContent='Stage '+(stageIndex+1)+' / '+stages.length;instruction.textContent=stage.hint;target.textContent=stage.word.split('').map(()=>'_').join(' ');picked=[];attempt.textContent='Your word: ';feedback.textContent='';nextBtn.disabled=true;grid.innerHTML='';stage.pool.forEach((letter)=>{const b=document.createElement('button');b.type='button';b.textContent=letter;b.addEventListener('click',()=>{picked.push(letter);attempt.textContent='Your word: '+picked.join('');if(picked.length===stage.word.length){Array.from(grid.querySelectorAll('button')).forEach((x)=>x.disabled=true);}});grid.appendChild(b);});}
async function saveState(extra){await window.DangdeAPI.save({key:'alphabet-echo-state',value:Object.assign({stageIndex,score,cleared,lastWord:picked.join(''),progress:progressForCurrent()},extra||{})});await window.DangdeAPI.setProgress({progress:progressForCurrent(),status:progressForCurrent()>=100?'completed':'in_progress'});}
document.getElementById('check-btn').addEventListener('click',async()=>{const stage=stages[stageIndex];const current=picked.join('');if(!current){feedback.textContent='Pick letters first.';return;}if(current===stage.word){feedback.textContent='Great! Correct word.';score+=2;cleared=Math.max(cleared,stageIndex+1);nextBtn.disabled=stageIndex===stages.length-1?true:false;}else{feedback.textContent='Not quite. Try again.';score=Math.max(0,score-1);}scoreEl.textContent='Stars: '+score;progressEl.textContent='Progress: '+progressForCurrent()+'%';await saveState({lastCheckCorrect:current===stage.word});if(stageIndex===stages.length-1&&current===stage.word){feedback.textContent='Amazing! All stages complete.';await window.DangdeAPI.setProgress({progress:100,status:'completed'});}});
document.getElementById('reset-btn').addEventListener('click',()=>{renderStage();});
nextBtn.addEventListener('click',async()=>{if(stageIndex<stages.length-1){stageIndex+=1;renderStage();await saveState({movedStage:true});}});
(async()=>{const rows=await window.DangdeAPI.query({key:'alphabet-echo-state'});const latest=rows[0];if(latest&&latest.value){stageIndex=Math.min(stages.length-1,latest.value.stageIndex||0);score=latest.value.score||0;cleared=latest.value.cleared||0;}scoreEl.textContent='Stars: '+score;progressEl.textContent='Progress: '+progressForCurrent()+'%';renderStage();})();`,
		},
		{
			ID:          2,
			Title:       "Count the Stars",
			Description: "Solve a progressive 3-level math mission with mixed addition and subtraction.",
			Language:    "English",
			Difficulty:  "medium",
			AgeGroup:    "5-8",
			CategoryID:  3,
			Prompt:      "Answer each level's questions. Reach the target score to unlock the next level.",
			Icon:        "➕",
			HTMLCode: `<div class='activity-root'>
  <h1>Count the Stars Challenge</h1>
  <p id='level'></p>
  <div class='problem-card'>
    <p id='problem'></p>
    <p id='progress-text'></p>
  </div>
  <div id='answer-row' class='answer-row'></div>
  <div class='stats'>
    <span id='score'>Score: 0</span>
    <span id='streak'>Streak: 0</span>
  </div>
  <p id='result'></p>
  <div class='actions'>
    <button id='next-question'>Next Question</button>
  </div>
</div>`,
			CSSCode: `.activity-root{font-family:'Avenir Next',sans-serif;padding:24px;min-height:100vh;background:linear-gradient(180deg,#eef9ff,#fff4ce)}
.problem-card{background:#fff;border-radius:22px;padding:18px;box-shadow:0 12px 24px rgba(0,0,0,.08)}
#problem{font-size:36px;font-weight:900;margin:0}
#progress-text{margin:8px 0 0;font-weight:700}
.answer-row{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:14px;margin-top:18px}
.answer-row button,.actions button{min-height:74px;border:0;border-radius:20px;font-size:28px;font-weight:800;background:#fff;box-shadow:0 10px 20px rgba(0,0,0,.08)}
.stats{display:flex;gap:16px;margin-top:14px;font-weight:800}
#result{margin-top:12px;font-size:22px;font-weight:800}`,
			JSCode: `const levels=[{name:'Level 1',target:3,pool:[[2,'+',3,5],[1,'+',4,5],[6,'-',2,4],[3,'+',2,5]]},{name:'Level 2',target:4,pool:[[7,'+',5,12],[9,'-',3,6],[8,'+',4,12],[10,'-',4,6]]},{name:'Level 3',target:5,pool:[[12,'+',9,21],[14,'-',6,8],[15,'+',7,22],[20,'-',11,9]]}];
let levelIndex=0;let solved=0;let score=0;let streak=0;let current=null;let asked=0;
const levelEl=document.getElementById('level');const problemEl=document.getElementById('problem');const answersEl=document.getElementById('answer-row');const resultEl=document.getElementById('result');const scoreEl=document.getElementById('score');const streakEl=document.getElementById('streak');const progressText=document.getElementById('progress-text');
function rnd(n){return Math.floor(Math.random()*n);}
function makeOptions(answer){const s=new Set([answer]);while(s.size<3){const delta=rnd(5)+1;const val=Math.max(0,answer+(rnd(2)?delta:-delta));s.add(val);}return Array.from(s).sort(()=>Math.random()-0.5);}
function chooseQuestion(){const pool=levels[levelIndex].pool;current=pool[rnd(pool.length)];asked+=1;problemEl.textContent=current[0]+' '+current[1]+' '+current[2]+' = ?';progressText.textContent='Need '+levels[levelIndex].target+' correct to clear this level.';const options=makeOptions(current[3]);answersEl.innerHTML='';options.forEach((opt)=>{const b=document.createElement('button');b.type='button';b.textContent=String(opt);b.addEventListener('click',()=>submitAnswer(opt));answersEl.appendChild(b);});}
async function sync(){const p=Math.min(100,Math.round(((levelIndex)+(solved/levels[levelIndex].target))/levels.length*100));await window.DangdeAPI.save({key:'count-stars-state',value:{levelIndex,solved,score,streak,asked,progress:p}});await window.DangdeAPI.setProgress({progress:p,status:p>=100?'completed':'in_progress'});scoreEl.textContent='Score: '+score;streakEl.textContent='Streak: '+streak;}
async function submitAnswer(choice){const ok=choice===current[3];if(ok){resultEl.textContent='Correct!';solved+=1;score+=2;streak+=1;}else{resultEl.textContent='Oops, keep trying!';score=Math.max(0,score-1);streak=0;}if(solved>=levels[levelIndex].target){if(levelIndex<levels.length-1){levelIndex+=1;solved=0;resultEl.textContent='Great! '+levels[levelIndex-1].name+' cleared. Welcome to '+levels[levelIndex].name+'.';}else{resultEl.textContent='Excellent! All levels complete.';await window.DangdeAPI.setProgress({progress:100,status:'completed'});}}await sync();renderLevel();}
function renderLevel(){levelEl.textContent=levels[levelIndex].name;chooseQuestion();}
document.getElementById('next-question').addEventListener('click',()=>{renderLevel();});
(async()=>{const rows=await window.DangdeAPI.query({key:'count-stars-state'});const latest=rows[0];if(latest&&latest.value){levelIndex=Math.min(levels.length-1,latest.value.levelIndex||0);solved=latest.value.solved||0;score=latest.value.score||0;streak=latest.value.streak||0;asked=latest.value.asked||0;}scoreEl.textContent='Score: '+score;streakEl.textContent='Streak: '+streak;renderLevel();await sync();})();`,
		},
		{
			ID:          3,
			Title:       "Arabic Letter Path",
			Description: "Learn Arabic letters through recognition, sound matching, and final mixed review.",
			Language:    "Arabic",
			Difficulty:  "hard",
			AgeGroup:    "5-8",
			CategoryID:  5,
			Prompt:      "Finish 3 stages: identify letters, choose sounds, then pass a mixed quiz.",
			Icon:        "📖",
			HTMLCode: `<div class='activity-root'>
  <h1>Arabic Letter Path</h1>
  <p id='stage'></p>
  <div id='card' class='card'></div>
  <p id='question'></p>
  <div id='options' class='options'></div>
  <div class='stats'>
    <span id='score'>Points: 0</span>
    <span id='progress'>Progress: 0%</span>
  </div>
  <p id='message'></p>
</div>`,
			CSSCode: `.activity-root{font-family:'Avenir Next',sans-serif;padding:24px;min-height:100vh;background:linear-gradient(180deg,#fff7ef,#eef8ff);text-align:center}
.card{display:grid;place-items:center;height:220px;border-radius:28px;background:#fff;font-size:110px;font-weight:800;margin:18px 0;box-shadow:0 12px 24px rgba(0,0,0,.08)}
.options{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px;max-width:520px;margin:0 auto}
.options button{min-height:64px;border:0;border-radius:18px;background:#fff;font-size:24px;font-weight:800;box-shadow:0 10px 20px rgba(0,0,0,.08)}
.stats{display:flex;justify-content:center;gap:14px;margin-top:16px;font-weight:800}
#message{margin-top:14px;font-size:20px;font-weight:800}`,
			JSCode: `const stages=[{name:'Stage 1: Letter Shape',items:[{letter:'ب',answer:'ba',opts:['ba','ta','da']},{letter:'ت',answer:'ta',opts:['ba','ta','sa']},{letter:'د',answer:'da',opts:['da','ra','za']}]},{name:'Stage 2: Sound Match',items:[{letter:'ر',answer:'ra',opts:['ra','ba','ha']},{letter:'س',answer:'sa',opts:['sa','ta','ya']},{letter:'م',answer:'ma',opts:['ma','na','la']}]},{name:'Stage 3: Mixed Review',items:[{letter:'ن',answer:'na',opts:['na','ma','sa']},{letter:'ي',answer:'ya',opts:['ya','ra','da']},{letter:'ل',answer:'la',opts:['la','ba','na']}]}];
let stageIndex=0;let itemIndex=0;let points=0;let completed=0;
const stageEl=document.getElementById('stage');const card=document.getElementById('card');const question=document.getElementById('question');const options=document.getElementById('options');const msg=document.getElementById('message');const scoreEl=document.getElementById('score');const progressEl=document.getElementById('progress');
function currentItem(){return stages[stageIndex].items[itemIndex];}
function progress(){const total=stages.reduce((sum,s)=>sum+s.items.length,0);return Math.min(100,Math.round((completed/total)*100));}
async function sync(extra){const p=progress();await window.DangdeAPI.save({key:'arabic-letter-path-state',value:Object.assign({stageIndex,itemIndex,points,completed,progress:p},extra||{})});await window.DangdeAPI.setProgress({progress:p,status:p>=100?'completed':'in_progress'});scoreEl.textContent='Points: '+points;progressEl.textContent='Progress: '+p+'%';}
function render(){const item=currentItem();stageEl.textContent=stages[stageIndex].name;card.textContent=item.letter;question.textContent='Choose the correct sound for this letter.';options.innerHTML='';item.opts.forEach((opt)=>{const b=document.createElement('button');b.type='button';b.textContent=opt;b.addEventListener('click',()=>answer(opt));options.appendChild(b);});}
async function answer(choice){const item=currentItem();if(choice===item.answer){points+=2;completed+=1;msg.textContent='Correct! '+item.letter+' sounds like '+item.answer+'.';itemIndex+=1;if(itemIndex>=stages[stageIndex].items.length){stageIndex+=1;itemIndex=0;if(stageIndex>=stages.length){msg.textContent='Excellent! You completed all Arabic stages.';await window.DangdeAPI.setProgress({progress:100,status:'completed'});stageIndex=stages.length-1;itemIndex=stages[stageIndex].items.length-1;}}}else{points=Math.max(0,points-1);msg.textContent='Not quite. Try again.';}await sync({lastChoice:choice});if(progress()<100){render();}}
(async()=>{const rows=await window.DangdeAPI.query({key:'arabic-letter-path-state'});const latest=rows[0];if(latest&&latest.value){stageIndex=Math.min(stages.length-1,latest.value.stageIndex||0);itemIndex=Math.min(stages[stageIndex].items.length-1,latest.value.itemIndex||0);points=latest.value.points||0;completed=latest.value.completed||0;}await sync();render();})();`,
		},
		{
			ID:          4,
			Title:       "Tell the Time",
			Description: "Master time reading through 3 levels: hour, half-hour, and quarter-time challenges.",
			Language:    "English",
			Difficulty:  "hard",
			AgeGroup:    "5-8",
			CategoryID:  6,
			Prompt:      "Solve the clock challenges in each level. Clear all rounds to complete the mission.",
			Icon:        "🕒",
			HTMLCode: `<div class='activity-root'>
  <h1>Tell the Time Quest</h1>
  <p id='level'></p>
  <div class='clock-card'>
    <div id='clock-time' class='clock-time'>3:00</div>
  </div>
  <p id='question'></p>
  <div id='answers' class='answer-row'></div>
  <div class='stats'>
    <span id='score'>Score: 0</span>
    <span id='progress'>Progress: 0%</span>
  </div>
  <p id='result'></p>
  <button id='next-round'>Next Round</button>
</div>`,
			CSSCode: `.activity-root{font-family:'Avenir Next',sans-serif;padding:24px;min-height:100vh;background:linear-gradient(180deg,#f4f9ff,#fff2d5)}
.clock-card{display:grid;place-items:center;height:210px;border-radius:30px;background:#fff;box-shadow:0 14px 26px rgba(0,0,0,.08)}
.clock-time{font-size:70px;font-weight:900}
.answer-row{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px;margin-top:16px}
.answer-row button,#next-round{min-height:68px;border:0;border-radius:20px;background:#fff;font-size:24px;font-weight:800;box-shadow:0 10px 20px rgba(0,0,0,.08)}
.stats{display:flex;gap:14px;margin-top:14px;font-weight:800}
#result{margin-top:12px;font-size:22px;font-weight:800}
#next-round{margin-top:10px;width:100%}`,
			JSCode: `const levels=[{name:'Level 1: Full Hour',target:3,pool:['1:00','3:00','6:00','9:00']},{name:'Level 2: Half Hour',target:3,pool:['1:30','4:30','7:30','10:30']},{name:'Level 3: Quarter Time',target:4,pool:['2:15','5:45','8:15','11:45']}];
let levelIndex=0;let score=0;let solved=0;let round=0;let current='3:00';
const levelEl=document.getElementById('level');const clockEl=document.getElementById('clock-time');const qEl=document.getElementById('question');const answers=document.getElementById('answers');const scoreEl=document.getElementById('score');const progressEl=document.getElementById('progress');const resultEl=document.getElementById('result');
function pick(arr){return arr[Math.floor(Math.random()*arr.length)];}
function optionsFor(answer,pool){const set=new Set([answer]);while(set.size<4){set.add(pick(pool));}return Array.from(set).sort(()=>Math.random()-0.5);}
function progress(){const base=levelIndex/levels.length;const part=(solved/levels[levelIndex].target)/levels.length;return Math.min(100,Math.round((base+part)*100));}
async function sync(extra){const p=progress();await window.DangdeAPI.save({key:'tell-time-state',value:Object.assign({levelIndex,score,solved,round,current,progress:p},extra||{})});await window.DangdeAPI.setProgress({progress:p,status:p>=100?'completed':'in_progress'});scoreEl.textContent='Score: '+score;progressEl.textContent='Progress: '+p+'%';}
function renderRound(){const lv=levels[levelIndex];levelEl.textContent=lv.name;current=pick(lv.pool);clockEl.textContent=current;qEl.textContent='What time is shown?';answers.innerHTML='';optionsFor(current,lv.pool).forEach((opt)=>{const b=document.createElement('button');b.type='button';b.textContent=opt;b.addEventListener('click',()=>answer(opt));answers.appendChild(b);});}
async function answer(choice){const ok=choice===current;if(ok){score+=2;solved+=1;resultEl.textContent='Correct!';}else{score=Math.max(0,score-1);resultEl.textContent='Not yet. Look carefully at the minute part.';}if(solved>=levels[levelIndex].target){if(levelIndex<levels.length-1){resultEl.textContent='Great! '+levels[levelIndex].name+' cleared.';levelIndex+=1;solved=0;}else{resultEl.textContent='Fantastic! You mastered all levels.';await window.DangdeAPI.setProgress({progress:100,status:'completed'});}}await sync({lastChoice:choice,correct:ok});if(progress()<100){renderRound();}}
document.getElementById('next-round').addEventListener('click',()=>{round+=1;renderRound();});
(async()=>{const rows=await window.DangdeAPI.query({key:'tell-time-state'});const latest=rows[0];if(latest&&latest.value){levelIndex=Math.min(levels.length-1,latest.value.levelIndex||0);score=latest.value.score||0;solved=latest.value.solved||0;round=latest.value.round||0;current=latest.value.current||current;}await sync();renderRound();})();`,
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
