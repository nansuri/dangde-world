<template>
  <div class="template-selector">
    <button
      v-for="template in templates"
      :key="template.id"
      class="template-button"
      type="button"
      :title="template.description"
      @click="selectTemplate(template)"
    >
      <span class="template-icon">{{ template.icon }}</span>
      <span class="template-name">{{ template.name }}</span>
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  onSelect: Function,
})

const templates = [
  {
    id: 'score-game',
    name: 'Score Game',
    icon: '🎮',
    description: 'Simple game that tracks score and progress',
    htmlCode: `<div class="game-container">
  <div class="score-display">
    <p>Score: <span id="score">0</span></p>
  </div>
  <div class="game-area">
    <button id="actionBtn" class="action-button">Click Me!</button>
  </div>
  <p class="instructions">Click the button to score points</p>
</div>`,
    cssCode: `.game-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 2rem;
  padding: 2rem;
  font-family: sans-serif;
}

.score-display {
  font-size: 2rem;
  font-weight: bold;
  color: #ff7a18;
}

.game-area {
  display: flex;
  gap: 1rem;
}

.action-button {
  padding: 1rem 2rem;
  font-size: 1.2rem;
  background: #37bcff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 200ms;
}

.action-button:hover {
  background: #1e9fd4;
  transform: scale(1.05);
}

.action-button:active {
  transform: scale(0.95);
}

.instructions {
  color: #57748f;
  font-size: 0.9rem;
}`,
    jsCode: `(async () => {
  let score = 0;
  const scoreDisplay = document.getElementById('score');
  const actionBtn = document.getElementById('actionBtn');

  // Load saved score if exists
  try {
    const records = await window.DangdeAPI.query({ key: 'game-score' });
    if (records.length > 0) {
      score = records[0].value;
      scoreDisplay.textContent = score;
    }
  } catch (error) {
    console.error('Error loading score:', error);
  }

  actionBtn.addEventListener('click', async () => {
    score += 10;
    scoreDisplay.textContent = score;

    try {
      // Save score
      await window.DangdeAPI.save({
        key: 'game-score',
        value: score
      });

      // Update progress (100 points = 100%)
      const progress = Math.min(100, score);
      await window.DangdeAPI.setProgress({ progress });
    } catch (error) {
      console.error('Error saving:', error);
    }
  });
})();`,
  },
  {
    id: 'quiz',
    name: 'Quiz',
    icon: '❓',
    description: 'Multiple choice quiz with scoring',
    htmlCode: `<div class="quiz-container">
  <div class="quiz-header">
    <h2>Quiz Time!</h2>
    <p>Question <span id="current">1</span> of 3</p>
  </div>
  
  <div class="quiz-body">
    <p id="question" class="question">What is 2 + 2?</p>
    <div class="options">
      <button class="option-btn" data-correct="true">4</button>
      <button class="option-btn">3</button>
      <button class="option-btn">5</button>
      <button class="option-btn">6</button>
    </div>
  </div>
  
  <div class="quiz-footer">
    <p id="feedback" class="feedback"></p>
    <button id="nextBtn" class="next-button" disabled>Next →</button>
  </div>
</div>`,
    cssCode: `.quiz-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 2rem;
  background: linear-gradient(135deg, #ffd75c 0%, #37bcff 100%);
}

.quiz-header {
  text-align: center;
  color: white;
  margin-bottom: 2rem;
}

.quiz-header h2 {
  margin: 0;
  font-size: 2rem;
}

.quiz-header p {
  margin: 0.5rem 0 0 0;
  opacity: 0.9;
}

.quiz-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2rem;
}

.question {
  font-size: 1.5rem;
  color: white;
  text-align: center;
  margin: 0;
}

.options {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.option-btn {
  padding: 1rem;
  background: white;
  border: 2px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: bold;
  transition: all 200ms;
}

.option-btn:hover {
  transform: translateY(-2px);
}

.option-btn.selected {
  border-color: #17324d;
}

.option-btn.correct {
  background: #4caf50;
  color: white;
}

.option-btn.incorrect {
  background: #f44336;
  color: white;
}

.quiz-footer {
  text-align: center;
  color: white;
}

.feedback {
  margin: 0 0 1rem 0;
  font-weight: bold;
  min-height: 1.5rem;
}

.next-button {
  padding: 0.75rem 2rem;
  background: white;
  color: #37bcff;
  border: none;
  border-radius: 8px;
  font-weight: bold;
  cursor: pointer;
  transition: all 200ms;
}

.next-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.next-button:not(:disabled):hover {
  transform: scale(1.05);
}`,
    jsCode: `(async () => {
  const questions = [
    { q: 'What is 2 + 2?', options: ['4', '3', '5', '6'], correct: 0 },
    { q: 'What is 3 + 5?', options: ['8', '7', '9', '6'], correct: 0 },
    { q: 'What is 10 - 3?', options: ['7', '6', '8', '5'], correct: 0 },
  ];

  let currentQuestion = 0;
  let score = 0;
  let answered = false;

  const questionEl = document.getElementById('question');
  const currentEl = document.getElementById('current');
  const optionBtns = document.querySelectorAll('.option-btn');
  const feedbackEl = document.getElementById('feedback');
  const nextBtn = document.getElementById('nextBtn');

  function renderQuestion() {
    answered = false;
    nextBtn.disabled = true;
    feedbackEl.textContent = '';
    
    const q = questions[currentQuestion];
    questionEl.textContent = q.q;
    currentEl.textContent = currentQuestion + 1;

    optionBtns.forEach((btn, i) => {
      btn.textContent = q.options[i];
      btn.classList.remove('selected', 'correct', 'incorrect');
      btn.disabled = false;
      btn.onclick = () => selectOption(i, q.correct);
    });
  }

  function selectOption(selectedIdx, correctIdx) {
    answered = true;
    optionBtns.forEach((btn, i) => {
      btn.disabled = true;
      if (i === correctIdx) {
        btn.classList.add('correct');
      } else if (i === selectedIdx) {
        btn.classList.add('incorrect');
      }
    });

    if (selectedIdx === correctIdx) {
      score++;
      feedbackEl.textContent = '✓ Correct!';
      feedbackEl.style.color = '#4caf50';
    } else {
      feedbackEl.textContent = '✗ Wrong! The correct answer is ' + questions[currentQuestion].options[correctIdx];
      feedbackEl.style.color = '#f44336';
    }

    nextBtn.disabled = false;
  }

  nextBtn.onclick = async () => {
    currentQuestion++;
    if (currentQuestion < questions.length) {
      renderQuestion();
    } else {
      // Quiz complete
      const percent = Math.round((score / questions.length) * 100);
      questionEl.textContent = \`Quiz Complete! You scored \${score}/\${questions.length} (\${percent}%)\`;
      nextBtn.style.display = 'none';

      try {
        await window.DangdeAPI.save({
          key: 'quiz-results',
          value: { score, percent }
        });
        await window.DangdeAPI.setProgress({ progress: percent });
      } catch (error) {
        console.error('Error saving results:', error);
      }
    }
  };

  renderQuestion();
})();`,
  },
  {
    id: 'drawing',
    name: 'Drawing',
    icon: '🎨',
    description: 'Simple drawing canvas',
    htmlCode: `<div class="drawing-container">
  <div class="drawing-toolbar">
    <button id="clearBtn" class="tool-btn">🗑️ Clear</button>
    <label class="color-picker">
      Color: <input type="color" id="colorPicker" value="#ff7a18" />
    </label>
    <label class="size-picker">
      Size: <input type="range" id="sizePicker" min="1" max="50" value="5" />
    </label>
  </div>
  <canvas id="drawingCanvas"></canvas>
  <button id="saveBtn" class="save-btn">💾 Save Drawing</button>
</div>`,
    cssCode: `.drawing-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 1rem;
  gap: 1rem;
  background: linear-gradient(135deg, #ffd75c 0%, #37bcff 100%);
}

.drawing-toolbar {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  align-items: center;
  background: rgba(255, 255, 255, 0.9);
  padding: 1rem;
  border-radius: 8px;
}

.tool-btn, .save-btn {
  padding: 0.5rem 1rem;
  background: #37bcff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: all 200ms;
}

.tool-btn:hover, .save-btn:hover {
  background: #1e9fd4;
}

.color-picker, .size-picker {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: bold;
}

.color-picker input, .size-picker input {
  cursor: pointer;
}

#drawingCanvas {
  flex: 1;
  background: white;
  border-radius: 8px;
  border: 2px solid rgba(255, 255, 255, 0.5);
  cursor: crosshair;
}

.save-btn {
  align-self: center;
  padding: 1rem 2rem;
  font-size: 1rem;
}`,
    jsCode: `(async () => {
  const canvas = document.getElementById('drawingCanvas');
  const ctx = canvas.getContext('2d');
  const colorPicker = document.getElementById('colorPicker');
  const sizePicker = document.getElementById('sizePicker');
  const clearBtn = document.getElementById('clearBtn');
  const saveBtn = document.getElementById('saveBtn');

  // Set canvas size
  function resizeCanvas() {
    canvas.width = canvas.offsetWidth;
    canvas.height = canvas.offsetHeight;
  }
  resizeCanvas();
  window.addEventListener('resize', resizeCanvas);

  let isDrawing = false;

  canvas.addEventListener('mousedown', (e) => {
    isDrawing = true;
    const rect = canvas.getBoundingClientRect();
    ctx.beginPath();
    ctx.moveTo(e.clientX - rect.left, e.clientY - rect.top);
  });

  canvas.addEventListener('mousemove', (e) => {
    if (!isDrawing) return;
    const rect = canvas.getBoundingClientRect();
    ctx.lineTo(e.clientX - rect.left, e.clientY - rect.top);
    ctx.strokeStyle = colorPicker.value;
    ctx.lineWidth = sizePicker.value;
    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
    ctx.stroke();
  });

  canvas.addEventListener('mouseup', () => isDrawing = false);
  canvas.addEventListener('mouseleave', () => isDrawing = false);

  clearBtn.onclick = () => {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
  };

  saveBtn.onclick = async () => {
    try {
      const imageData = canvas.toDataURL();
      await window.DangdeAPI.save({
        key: 'drawing',
        value: imageData
      });
      await window.DangdeAPI.setProgress({ progress: 100 });
      alert('Drawing saved!');
    } catch (error) {
      console.error('Error saving drawing:', error);
    }
  };
})();`,
  },
  {
    id: 'memory-game',
    name: 'Memory Game',
    icon: '🧠',
    description: 'Flip card memory matching game',
    htmlCode: `<div class="memory-container">
  <div class="memory-header">
    <h2>Memory Game</h2>
    <p>Matches: <span id="matches">0</span>/6</p>
  </div>
  <div class="memory-grid" id="grid"></div>
  <button id="restartBtn" class="restart-btn">🔄 Restart</button>
</div>`,
    cssCode: `.memory-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  padding: 2rem;
  gap: 2rem;
  background: linear-gradient(135deg, #ffd75c 0%, #37bcff 100%);
}

.memory-header {
  text-align: center;
  color: white;
}

.memory-header h2 {
  margin: 0;
  font-size: 2rem;
}

.memory-header p {
  margin: 0.5rem 0 0 0;
  font-size: 1.2rem;
}

.memory-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  width: 100%;
  max-width: 400px;
}

.memory-card {
  aspect-ratio: 1;
  background: white;
  border: 2px solid rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  cursor: pointer;
  font-size: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 200ms;
  user-select: none;
}

.memory-card:hover:not(.flipped):not(.matched) {
  transform: scale(1.05);
}

.memory-card.flipped {
  background: #37bcff;
}

.memory-card.matched {
  background: #4caf50;
  cursor: default;
}

.restart-btn {
  padding: 1rem 2rem;
  background: white;
  color: #37bcff;
  border: none;
  border-radius: 8px;
  font-weight: bold;
  font-size: 1rem;
  cursor: pointer;
  transition: all 200ms;
}

.restart-btn:hover {
  transform: scale(1.05);
}`,
    jsCode: `(async () => {
  const emojis = ['🎈', '🎸', '🎲', '🎭', '🎪', '🎯'];
  const grid = document.getElementById('grid');
  const matchesEl = document.getElementById('matches');
  const restartBtn = document.getElementById('restartBtn');
  
  let cards = [];
  let flipped = [];
  let matched = 0;

  function init() {
    cards = [...emojis, ...emojis].sort(() => Math.random() - 0.5);
    flipped = [];
    matched = 0;
    matchesEl.textContent = matched;
    grid.innerHTML = '';

    cards.forEach((emoji, i) => {
      const card = document.createElement('div');
      card.className = 'memory-card';
      card.textContent = '?';
      card.dataset.index = i;
      card.onclick = () => flipCard(card, i);
      grid.appendChild(card);
    });
  }

  function flipCard(card, i) {
    if (flipped.length === 2 || flipped.includes(i) || card.classList.contains('matched')) return;

    card.textContent = cards[i];
    card.classList.add('flipped');
    flipped.push(i);

    if (flipped.length === 2) {
      setTimeout(checkMatch, 600);
    }
  }

  function checkMatch() {
    const [i1, i2] = flipped;
    if (cards[i1] === cards[i2]) {
      document.querySelectorAll('[data-index]')[i1].classList.add('matched');
      document.querySelectorAll('[data-index]')[i2].classList.add('matched');
      matched++;
      matchesEl.textContent = matched;

      if (matched === emojis.length) {
        setTimeout(() => alert('You won!'), 500);
        saveProgress();
      }
    } else {
      document.querySelectorAll('[data-index]')[i1].textContent = '?';
      document.querySelectorAll('[data-index]')[i2].textContent = '?';
      document.querySelectorAll('[data-index]')[i1].classList.remove('flipped');
      document.querySelectorAll('[data-index]')[i2].classList.remove('flipped');
    }
    flipped = [];
  }

  async function saveProgress() {
    try {
      const progress = (matched / emojis.length) * 100;
      await window.DangdeAPI.save({
        key: 'game-progress',
        value: { matched, total: emojis.length }
      });
      await window.DangdeAPI.setProgress({ progress: 100 });
    } catch (error) {
      console.error('Error saving progress:', error);
    }
  }

  restartBtn.onclick = init;

  init();
})();`,
  },
];

function selectTemplate(template) {
  if (props.onSelect) {
    props.onSelect(template);
  }
}
</script>

<style scoped>
.template-selector {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 1rem;
  padding: 1rem;
  background: var(--surface);
  border-radius: 8px;
}

.template-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: white;
  border: 2px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: all 200ms;
  text-align: center;
}

.template-button:hover {
  border-color: var(--primary);
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.template-icon {
  font-size: 2rem;
}

.template-name {
  font-weight: 600;
  color: var(--ink);
  font-size: 0.9rem;
}
</style>
