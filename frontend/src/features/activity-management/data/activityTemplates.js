export const ACTIVITY_TEMPLATES = [
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
      await window.DangdeAPI.save({ key: 'game-score', value: score });
      await window.DangdeAPI.setProgress({ progress: Math.min(100, score) });
    } catch (error) {
      console.error('Error saving score:', error);
    }
  });
})();`,
  },
  {
    id: 'quiz',
    name: 'Quiz',
    icon: '❓',
    description: 'Multiple choice quiz with instant feedback',
    htmlCode: `<div class="quiz-container">
  <div class="question-panel">
    <h2 id="question">What is 2 + 2?</h2>
    <div class="options">
      <button class="option-btn" data-option="0">3</button>
      <button class="option-btn" data-option="1">4</button>
      <button class="option-btn" data-option="2">5</button>
    </div>
  </div>
  <div class="score-panel">
    <p>Score: <span id="score">0</span>/<span id="total">0</span></p>
  </div>
</div>`,
    cssCode: `.quiz-container {
  display: grid;
  grid-template-rows: 1fr auto;
  height: 100%;
  padding: 2rem;
  gap: 2rem;
}

.question-panel {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2rem;
}

.question-panel h2 {
  font-size: 1.8rem;
  color: var(--ink);
}

.options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 1rem;
}

.option-btn {
  padding: 1rem;
  font-size: 1rem;
  background: #37bcff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 200ms;
}

.option-btn:hover {
  background: #1e9fd4;
}

.option-btn.correct {
  background: #4caf50;
}

.option-btn.incorrect {
  background: #f44336;
}

.score-panel {
  text-align: center;
  padding: 1rem;
  background: rgba(0,0,0,0.05);
  border-radius: 8px;
}`,
    jsCode: `(async () => {
  const options = document.querySelectorAll('.option-btn');
  let score = 0;
  let total = 0;
  const scoreDisplay = document.getElementById('score');
  const totalDisplay = document.getElementById('total');

  options.forEach(btn => {
    btn.addEventListener('click', async () => {
      total++;
      const isCorrect = btn.dataset.option === '1';
      if (isCorrect) score++;
      
      scoreDisplay.textContent = score;
      totalDisplay.textContent = total;
      
      try {
        await window.DangdeAPI.save({
          key: 'quiz-score',
          value: { score, total }
        });
      } catch (error) {
        console.error('Error:', error);
      }
    });
  });
})();`,
  },
  {
    id: 'memory-game',
    name: 'Memory Game',
    icon: '🧠',
    description: 'Matching pairs memory game',
    htmlCode: `<div class="memory-container">
  <h2>Find the Pairs!</h2>
  <div class="memory-grid" id="grid"></div>
  <button id="restartBtn" class="restart-button">Restart</button>
</div>`,
    cssCode: `.memory-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem;
  gap: 2rem;
}

.memory-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  max-width: 400px;
}

.memory-card {
  aspect-ratio: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 2rem;
  transition: all 200ms;
  color: transparent;
}

.memory-card.flipped {
  color: inherit;
  background: white;
}

.memory-card:hover {
  transform: scale(1.05);
}

.restart-button {
  padding: 0.75rem 2rem;
  font-size: 1rem;
  background: #ff7a18;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}`,
    jsCode: `(async () => {
  const emojis = ['🍎', '🍌', '🍇', '🍓', '🍎', '🍌', '🍇', '🍓'];
  const grid = document.getElementById('grid');
  const restartBtn = document.getElementById('restartBtn');
  let flipped = [];
  let matched = 0;

  function shuffle(arr) {
    return [...arr].sort(() => Math.random() - 0.5);
  }

  function init() {
    grid.innerHTML = '';
    flipped = [];
    matched = 0;
    const shuffled = shuffle(emojis);
    
    shuffled.forEach((emoji, i) => {
      const card = document.createElement('button');
      card.className = 'memory-card';
      card.textContent = '?';
      card.dataset.emoji = emoji;
      card.dataset.index = i;
      card.addEventListener('click', handleCardClick);
      grid.appendChild(card);
    });
  }

  function handleCardClick(e) {
    if (flipped.length === 2) return;
    const card = e.target;
    if (card.classList.contains('flipped')) return;
    
    card.classList.add('flipped');
    card.textContent = card.dataset.emoji;
    flipped.push(card);

    if (flipped.length === 2) {
      checkMatch();
    }
  }

  function checkMatch() {
    const [card1, card2] = flipped;
    const match = card1.dataset.emoji === card2.dataset.emoji;

    if (match) {
      matched++;
      if (matched === emojis.length / 2) {
        setTimeout(() => alert('You won!'), 500);
        saveProgress();
      }
    } else {
      setTimeout(() => {
        card1.classList.remove('flipped');
        card2.classList.remove('flipped');
        card1.textContent = '?';
        card2.textContent = '?';
        flipped = [];
      }, 600);
    }
    flipped = [];
  }

  async function saveProgress() {
    try {
      const progress = (matched / (emojis.length / 2)) * 100;
      await window.DangdeAPI.save({
        key: 'game-progress',
        value: { matched, total: emojis.length / 2 }
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
