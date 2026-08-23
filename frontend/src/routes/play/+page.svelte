<script lang="ts">
	import { apiUrl } from '$lib/api/client.svelte';
	import { goto } from '$app/navigation';
	import { Building2, GraduationCap, BookOpen, Layers, Sparkles } from 'lucide-svelte';
	import type { PlayClass, PlaySubject, PlayTopic, PlayQuestion } from '$lib/types';

	type Phase = 'welcome' | 'classes' | 'subjects' | 'topics' | 'difficulty' | 'quiz' | 'results';

	let phase = $state<Phase>('welcome');
	let playerName = $state('');
	let classes = $state<PlayClass[]>([]);
	let subjects = $state<PlaySubject[]>([]);
	let topics = $state<PlayTopic[]>([]);
	let questions = $state<PlayQuestion[]>([]);

	let selectedClass = $state<PlayClass | null>(null);
	let selectedSubject = $state<PlaySubject | null>(null);
	let selectedTopic = $state('');
	let selectedDifficulty = $state('');
	let loading = $state(false);

	let currentIndex = $state(0);
	let score = $state(0);
	let streak = $state(0);
	let bestStreak = $state(0);
	let correctCount = $state(0);
	let selectedKey = $state('');
	let answered = $state(false);
	let isCorrect = $state(false);
	let timeLeft = $state(15);
	let timerInterval = $state<ReturnType<typeof setInterval> | undefined>();
	let startTime = $state(0);
	let showExitConfirm = $state(false);
	let screenShake = $state(false);
	let scorePopups = $state<{ id: number; value: number; x: number; y: number }[]>([]);
	let popupId = $state(0);
	let confettiPieces = $state<{ id: number; x: number; color: string; rot: number; delay: number }[]>([]);
	let answerBounce = $state(false);

	const CONFETTI_COLORS = ['#FFC233', '#D8F3E3', '#0E7C71', '#FBDAD3', '#6B3FA0', '#FDE9C2'];

	function shuffle<T>(arr: T[]): T[] {
		const a = [...arr];
		for (let i = a.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));
			[a[i], a[j]] = [a[j], a[i]];
		}
		return a;
	}

	let audioCtx: AudioContext | null = null;
	function getAudio() {
		if (!audioCtx && typeof window !== 'undefined') {
			audioCtx = new (window.AudioContext || (window as any).webkitAudioContext)();
		}
		return audioCtx;
	}

	function playTone(freq: number, dur: number, type: OscillatorType = 'sine', vol = 0.12) {
		const ctx = getAudio();
		if (!ctx) return;
		if (ctx.state === 'suspended') ctx.resume();
		const osc = ctx.createOscillator();
		const gain = ctx.createGain();
		osc.connect(gain); gain.connect(ctx.destination);
		osc.type = type; osc.frequency.value = freq;
		gain.gain.setValueAtTime(vol, ctx.currentTime);
		gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + dur);
		osc.start(ctx.currentTime); osc.stop(ctx.currentTime + dur);
	}

	function playCorrect() { playTone(523, 0.1); setTimeout(() => playTone(659, 0.1), 80); setTimeout(() => playTone(784, 0.2), 160); }
	function playWrong() { playTone(200, 0.3, 'sawtooth', 0.08); }
	function playStreak() { playTone(600, 0.08); setTimeout(() => playTone(800, 0.08), 60); setTimeout(() => playTone(1000, 0.15), 120); }
	function playComplete() { [523, 659, 784, 1047].forEach((f, i) => setTimeout(() => playTone(f, 0.25), i * 120)); }
	function playClick() { playTone(440, 0.06, 'sine', 0.06); }

	function spawnConfetti() {
		const pieces = [];
		for (let i = 0; i < 40; i++) {
			pieces.push({ id: ++popupId, x: 10 + Math.random() * 80, color: CONFETTI_COLORS[i % CONFETTI_COLORS.length], rot: Math.random() * 360, delay: Math.random() * 0.5 });
		}
		confettiPieces = pieces;
		setTimeout(() => { confettiPieces = []; }, 3000);
	}

	function addPopup(value: number) {
		const id = ++popupId;
		scorePopups = [...scorePopups, { id, value, x: 50 + (Math.random() - 0.5) * 30, y: 40 }];
		setTimeout(() => { scorePopups = scorePopups.filter(p => p.id !== id); }, 1200);
	}

	async function api<T>(method: string, path: string, body?: unknown): Promise<T | null> {
		try {
			const opts: RequestInit = { method, headers: { 'Content-Type': 'application/json' } };
			if (body) opts.body = JSON.stringify(body);
			const res = await fetch(apiUrl(path), opts);
			const json = await res.json();
			return json.data ?? null;
		} catch { return null; }
	}

	async function loadClasses() {
		playClick();
		loading = true;
		classes = (await api<PlayClass[]>('GET', '/play/classes')) ?? [];
		loading = false;
		phase = 'classes';
	}

	async function loadSubjects(cls: PlayClass) {
		playClick();
		selectedClass = cls;
		loading = true;
		subjects = (await api<PlaySubject[]>('GET', `/play/subjects?class_id=${cls.id}`)) ?? [];
		loading = false;
		phase = 'subjects';
	}

	async function loadTopics(sub: PlaySubject) {
		playClick();
		selectedSubject = sub;
		loading = true;
		topics = (await api<PlayTopic[]>('GET', `/play/topics?class_id=${selectedClass!.id}&subject_id=${sub.id}`)) ?? [];
		loading = false;
		phase = 'topics';
	}

	async function startQuiz(difficulty: string) {
		playClick();
		selectedDifficulty = difficulty;
		loading = true;
		const topicParam = selectedTopic ? `&topic=${encodeURIComponent(selectedTopic)}` : '';
		const data = await api<PlayQuestion[]>('GET', `/play/quiz?class_id=${selectedClass!.id}&subject_id=${selectedSubject!.id}${topicParam}&difficulty=${difficulty}&limit=10`);
		loading = false;
		if (!data || data.length === 0) { alert('No questions found. Try different options.'); return; }
		questions = data.map(q => ({ ...q, options: shuffle(q.options) }));
		currentIndex = 0; score = 0; streak = 0; bestStreak = 0; correctCount = 0;
		selectedKey = ''; answered = false; startTime = Date.now();
		phase = 'quiz';
		startTimer();
	}

	function startTimer() {
		clearInterval(timerInterval);
		timeLeft = 15;
		timerInterval = setInterval(() => {
			timeLeft -= 0.1;
			if (timeLeft <= 0) { timeLeft = 0; clearInterval(timerInterval); handleAnswer(''); }
		}, 100);
	}

	function handleAnswer(key: string) {
		if (answered) return;
		answered = true; selectedKey = key;
		const correct = questions[currentIndex].options.find(o => o.correct);
		isCorrect = key !== '' && key === correct?.key;
		if (isCorrect) {
			correctCount++; streak++;
			if (streak > bestStreak) bestStreak = streak;
			const timePoints = Math.round(timeLeft * 10);
			const multiplier = Math.min(streak, 4);
			const pts = (100 + timePoints) * multiplier;
			score += pts;
			addPopup(pts);
			if (streak >= 3) playStreak(); else playCorrect();
			if (streak >= 5) spawnConfetti();
		} else {
			streak = 0;
			screenShake = true;
			playWrong();
			setTimeout(() => { screenShake = false; }, 500);
		}
		answerBounce = true;
		setTimeout(() => { answerBounce = false; }, 300);
		clearInterval(timerInterval);
		setTimeout(nextQuestion, 1200);
	}

	function nextQuestion() {
		if (currentIndex >= questions.length - 1) { finishQuiz(); return; }
		currentIndex++; selectedKey = ''; answered = false; startTimer();
	}

	function finishQuiz() {
		clearInterval(timerInterval);
		playComplete();
		spawnConfetti();
		phase = 'results';
	}

	function goBack() {
		playClick();
		const back: Partial<Record<Phase, Phase>> = {
			classes: 'welcome', subjects: 'classes', topics: 'subjects',
			difficulty: 'topics', quiz: 'welcome', results: 'welcome'
		};
		phase = back[phase] ?? 'welcome';
	}

	function confirmExit() { showExitConfirm = true; }
	function cancelExit() { showExitConfirm = false; }
	function confirmGoHome() { showExitConfirm = false; clearInterval(timerInterval); goto('/'); }
	function handleIdentityClick() {
		if (phase === 'welcome') goto('/');
		else confirmExit();
	}
	function playAgain() { playClick(); phase = 'difficulty'; }

	function timeColor() { return timeLeft > 10 ? '#0E7C71' : timeLeft > 5 ? '#B45309' : '#C2381B'; }
	function progressWidth() { return questions.length ? `${((currentIndex + 1) / questions.length) * 100}%` : '0%'; }
	function accuracy() { return questions.length ? Math.round((correctCount / questions.length) * 100) : 0; }
	function formatTime(ms: number) { const s = ms / 1000; return s < 60 ? `${s.toFixed(1)}s` : `${Math.floor(s / 60)}m ${Math.floor(s % 60)}s`; }
	function stars() { return accuracy() >= 80 ? '★★★' : accuracy() >= 50 ? '★★' : '★'; }
	function starsLabel() { return accuracy() >= 80 ? 'Excellent' : accuracy() >= 50 ? 'Good work' : 'Keep practising'; }
</script>

<svelte:head>
	<title>Quiz Arena — Pragati</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
</svelte:head>

{#each confettiPieces as cp (cp.id)}
	<div class="confetti-piece" style="left:{cp.x}%;background:{cp.color};animation-delay:{cp.delay}s;--rot:{cp.rot}deg"></div>
{/each}

{#each scorePopups as popup (popup.id)}
	<div class="score-popup" style="left:{popup.x}%;top:{popup.y}%">+{popup.value}</div>
{/each}

<div class="page {screenShake ? 'shake' : ''}">
	<header class="quiz-header">
		<button onclick={handleIdentityClick} class="qh-identity" aria-label="Go home">
			<span class="qh-logo"><Building2 size={16} /></span>
			<span class="qh-wordmark qh-wordmark-kannada">ಪ್ರಗತಿ</span>
			<span class="qh-sub">MDRS (SC-32) Bahaddurghatta</span>
		</button>
		{#if phase !== 'results'}
			<button onclick={confirmExit} class="qh-home">← Home</button>
		{/if}
	</header>

	{#if showExitConfirm}
		<div class="modal-scrim" role="presentation" onclick={cancelExit}></div>
		<div class="modal" role="dialog" aria-modal="true" aria-label="Quit quiz">
			<button class="modal-close" onclick={cancelExit} aria-label="Close">✕</button>
			<div class="modal-emoji" aria-hidden="true"><Layers size={28} /></div>
			<h3 class="modal-title">Quit quiz?</h3>
			<p class="modal-sub">Your progress will be lost.</p>
			<div class="modal-actions">
				<button onclick={confirmGoHome} class="btn-primary">Yes, quit</button>
				<button onclick={cancelExit} class="btn-ghost">Keep playing</button>
			</div>
		</div>
	{/if}

	<!-- ═══ WELCOME ═══ -->
	{#if phase === 'welcome'}
		<div class="center fade-in">
			<div class="q-card welcome-card">
				<div class="welcome-icon" aria-hidden="true"><Sparkles size={28} /></div>
				<h1 class="welcome-title">Quizzes</h1>
				<p class="welcome-sub">What's your name, superstar?</p>
				<form class="welcome-form" onsubmit={(e) => { e.preventDefault(); if (playerName.trim()) loadClasses(); }}>
					<label class="field">
						<span class="field-label">Your name</span>
						<input bind:value={playerName} placeholder="Enter your name…" maxlength={50} class="input" required />
					</label>
					<button type="submit" disabled={!playerName.trim()} class="btn-primary btn-block">
						Start playing →
					</button>
					<p class="hint">Just your name — ready to play?</p>
				</form>
			</div>
		</div>

	<!-- ═══ CLASSES ═══ -->
	{:else if phase === 'classes'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">Select class</h2>
					<p class="section-sub">{playerName ? `${playerName}’s turn` : 'Choose your class'}</p>
				</div>
			</div>
			{#if loading}
				<div class="empty">Loading…</div>
			{:else if classes.length === 0}
				<div class="empty">No classes with questions yet.</div>
			{:else}
				<div class="pick-grid">
					{#each classes as cls (cls.id)}
						<button onclick={() => loadSubjects(cls)} class="pick-card">
							<span class="pick-icon"><GraduationCap size={18} /></span>
							<span class="pick-name">{cls.name}</span>
							<span class="pick-meta">{cls.question_count} questions</span>
						</button>
					{/each}
				</div>
			{/if}
		</div>

	<!-- ═══ SUBJECTS ═══ -->
	{:else if phase === 'subjects'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">Select subject</h2>
					<p class="section-sub">{selectedClass?.name}</p>
				</div>
			</div>
			{#if loading}
				<div class="empty">Loading…</div>
			{:else}
				<div class="pick-grid">
					{#each subjects as sub (sub.id)}
						<button onclick={() => loadTopics(sub)} class="pick-card">
							<span class="pick-icon"><BookOpen size={18} /></span>
							<span class="pick-name">{sub.name}</span>
							<span class="pick-meta">{sub.question_count} questions</span>
						</button>
					{/each}
				</div>
			{/if}
		</div>

	<!-- ═══ TOPICS ═══ -->
	{:else if phase === 'topics'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">Select topic</h2>
					<p class="section-sub">{selectedSubject?.name}</p>
				</div>
			</div>
			{#if loading}
				<div class="empty">Loading…</div>
			{:else}
				<div class="topics">
					<button onclick={() => { playClick(); selectedTopic = ''; phase = 'difficulty'; }} class="topic topic-all">All topics</button>
					{#each topics as topic (topic.name)}
						<button onclick={() => { playClick(); selectedTopic = topic.name; phase = 'difficulty'; }} class="topic">{topic.name}</button>
					{/each}
				</div>
			{/if}
		</div>

	<!-- ═══ DIFFICULTY ═══ -->
	{:else if phase === 'difficulty'}
		<div class="center fade-in">
			<div class="q-card diff-card">
				<div class="section-head" style="margin:0">
					<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
					<div>
						<h2 class="section-title">Choose difficulty</h2>
						<p class="section-sub">{selectedClass?.name} → {selectedSubject?.name}{selectedTopic ? ` → ${selectedTopic}` : ''}</p>
					</div>
				</div>
				{#if loading}
					<div class="empty">Loading…</div>
				{:else}
					<div class="diff-stack">
						<button onclick={() => startQuiz('easy')} class="diff-btn diff-easy">
							<span class="diff-left">Easy</span>
							<span class="diff-meta">
								<span class="dots"><span class="dot filled"></span><span class="dot"></span><span class="dot"></span></span>
								Gentle start
							</span>
						</button>
						<button onclick={() => startQuiz('medium')} class="diff-btn diff-medium">
							<span class="diff-left">Medium</span>
							<span class="diff-meta">
								<span class="dots"><span class="dot filled"></span><span class="dot filled"></span><span class="dot"></span></span>
								Level up
							</span>
						</button>
						<button onclick={() => startQuiz('hard')} class="diff-btn diff-hard">
							<span class="diff-left">Hard</span>
							<span class="diff-meta">
								<span class="dots"><span class="dot filled"></span><span class="dot filled"></span><span class="dot filled"></span></span>
								Expert only
							</span>
						</button>
					</div>
				{/if}
			</div>
		</div>

	<!-- ═══ QUIZ ═══ -->
	{:else if phase === 'quiz' && questions.length > 0}
		{@const q = questions[currentIndex]}
		<div class="quiz-shell fade-in">
			<!-- main card -->
			<div class="quiz-main">
				<div class="q-meta">
					<div class="q-tags">
						<span class="tag">{selectedClass?.name}</span>
						<span class="tag">{selectedSubject?.name}</span>
						{#if selectedTopic}<span class="tag">{selectedTopic}</span>{/if}
						<span class="tag tag-difficulty tag-{selectedDifficulty}">{selectedDifficulty}</span>
					</div>
					<span class="counter">{currentIndex + 1} / {questions.length}</span>
				</div>

				<div class="q-card q-question {answerBounce ? 'answer-pop' : ''}">
					<fieldset class="q-fieldset">
						<legend class="q-legend">{q.question_text}</legend>
						<div class="options">
							{#each q.options as opt, oi (opt.key)}
								{@const isSelected = answered && opt.key === selectedKey}
								{@const isCorrectOpt = !!opt.correct}
								<label class="opt
									{answered && isSelected ? (isCorrect ? 'opt-correct' : 'opt-incorrect') : ''}
									{answered && !isSelected && isCorrectOpt ? 'opt-correct' : ''}
									{!answered && selectedKey === opt.key ? 'opt-selected' : ''}"
								>
									<input type="radio" name="q{currentIndex}" value={opt.key} disabled={answered} checked={selectedKey === opt.key} onchange={() => handleAnswer(opt.key)} class="opt-input" />
									<span class="opt-row">
										<span class="opt-radio"><span class="opt-dot"></span></span>
										<span class="opt-text">{opt.value}</span>
										{#if answered && isSelected}
											<span class="opt-state">{isCorrect ? '✓' : '✕'}</span>
										{:else if answered && isCorrectOpt}
											<span class="opt-state">✓</span>
										{/if}
									</span>
								</label>
							{/each}
						</div>
					</fieldset>
					{#if answered}
						<div class="feedback {isCorrect ? '' : 'feedback-bad'}">
							<div class="feedback-head">
								<span class="feedback-badge">{isCorrect ? '✓' : '✕'}</span>
								<span class="feedback-title">{isCorrect ? 'Correct!' : 'Not quite'}</span>
								{#if streak >= 3 && isCorrect}<span class="streak-chip">{streak} streak</span>{/if}
							</div>
							<p class="feedback-text">
								{#if isCorrect}
									Nice work — you earned points with a {streak}× multiplier.
								{:else}
									The correct answer is highlighted above.
								{/if}
							</p>
						</div>
					{/if}
				</div>
			</div>

			<!-- side panel -->
			<aside class="quiz-side">
				<div class="side-card">
					<div class="side-row">
						<span class="side-label">Score</span>
						<span class="side-value mono">{score.toLocaleString()}</span>
					</div>
					<div class="progress">
						<div class="progress-bar" style="width:{progressWidth()}"></div>
					</div>
					<div class="side-grid">
						<div class="stat">
							<span class="stat-v mono">{accuracy()}%</span>
							<span class="stat-k">Accuracy</span>
						</div>
						<div class="stat">
							<span class="stat-v mono">{correctCount}/{questions.length}</span>
							<span class="stat-k">Correct</span>
						</div>
						<div class="stat">
							<span class="stat-v mono">{bestStreak}</span>
							<span class="stat-k">Best streak</span>
						</div>
					</div>
				</div>

				<div class="side-card">
					<div class="timer-head">
						<span class="side-label">Time</span>
						<span class="timer-num mono" style="color:{timeColor()}">{Math.ceil(timeLeft)}s</span>
					</div>
					<div class="timer-track">
						<div class="timer-fill" style="width:{(timeLeft/15)*100}%; background:{timeColor()}"></div>
					</div>
					{#if streak >= 2}
						<div class="streak-box">{streak}× streak</div>
					{/if}
				</div>

				<div class="side-card side-muted">
					<p class="side-hint">Tip: quick answers earn a bit extra. Get 3 in a row for a streak!</p>
				</div>
			</aside>
		</div>

	<!-- ═══ RESULTS ═══ -->
	{:else if phase === 'results'}
		<div class="center fade-in">
			<div class="q-card score-card">
				<p class="score-kicker">MDRS (SC-32) Bahaddurghatta · KREIS</p>
				<div class="score-stars" aria-hidden="true">{stars()}</div>
				<div class="score-stars-label">{starsLabel()}</div>
				<div class="score-num mono">{score.toLocaleString()}</div>
				<p class="score-label">Points</p>
				<p class="score-name">{playerName}</p>
				<p class="score-meta">{selectedClass?.name} · {selectedSubject?.name}{selectedTopic ? ` · ${selectedTopic}` : ''} · {selectedDifficulty}</p>
				<div class="score-stats">
					<div class="s-stat"><span class="s-v mono">{accuracy()}%</span><span class="s-k">Accuracy</span></div>
					<div class="s-sep"></div>
					<div class="s-stat"><span class="s-v mono">{bestStreak}</span><span class="s-k">Best streak</span></div>
					<div class="s-sep"></div>
					<div class="s-stat"><span class="s-v mono">{correctCount}/{questions.length}</span><span class="s-k">Correct</span></div>
				</div>
				<p class="score-time">{formatTime(Date.now() - startTime)} · {questions.length} questions</p>
				<div class="score-actions">
					<button onclick={playAgain} class="btn-primary btn-block">Play again →</button>
					<a href="/" class="btn-ghost btn-block" style="text-align:center; text-decoration:none; display:flex; justify-content:center;">Back to home</a>
				</div>
				<p class="hint">Take a screenshot to share your score.</p>
			</div>
		</div>
	{/if}
</div>

<style>
	.page {
		max-width: 1040px;
		margin: 0 auto;
		padding: 1rem clamp(1rem, 4vw, 1.5rem) 3rem;
	}
	.quiz-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
		margin: 0.75rem 0 1.5rem;
	}
	.qh-identity { display: flex; align-items: center; gap: 0.6rem; color: var(--ink); background: transparent; border: 0; padding: 0; cursor: pointer; font: inherit; text-align: left; }
	.qh-sub { font-size: 0.78rem; color: var(--ink-soft); font-weight: 600; margin-left: 0.35rem; }
	.qh-logo {
		width: 38px; height: 38px; border-radius: 10px; border: 2.5px solid var(--ink);
		background: var(--paper); box-shadow: 3px 3px 0 var(--ink);
		display: grid; place-items: center; font-size: 1.15rem; flex: none;
	}
	.qh-wordmark { font-family: var(--font-display); font-weight: 800; letter-spacing: -0.02em; font-size: 1.15rem; }
	.qh-wordmark-kannada { font-family: 'Anek Kannada', var(--font-display), system-ui, sans-serif; font-weight: 700; }
	.qh-dot { color: var(--coral); }
	.qh-home {
		font-family: var(--font-body); font-weight: 700; font-size: 0.9rem;
		border: 2.5px solid var(--ink); background: var(--paper); border-radius: 10px;
		padding: 0.4rem 0.75rem; cursor: pointer; color: var(--ink);
	}
	.qh-home:hover { background: var(--cream-deep); }

	/* modal */
	.modal-scrim { position: fixed; inset: 0; background: rgba(31,26,46,0.45); backdrop-filter: blur(2px); z-index: 60; }
	.modal {
		position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%);
		z-index: 61; background: var(--paper); border: 2.5px solid var(--ink);
		box-shadow: 6px 6px 0 var(--ink); border-radius: 18px;
		padding: 1.35rem; width: min(420px, calc(100vw - 2rem));
		text-align: center;
	}
	.modal-close {
		position: absolute; top: 0.6rem; right: 0.6rem;
		width: 36px; height: 36px; border-radius: 10px; border: 2.5px solid var(--ink);
		background: var(--paper); cursor: pointer; display: grid; place-items: center; color: var(--ink);
	}
	.modal-emoji { font-size: 2.5rem; }
	.modal-title { font-family: var(--font-display); font-size: 1.35rem; font-weight: 800; margin: 0.4rem 0 0; }
	.modal-sub { color: var(--ink-soft); margin: 0.25rem 0 0; }
	.modal-actions { display: flex; flex-direction: column; gap: 0.6rem; margin-top: 1rem; }

	/* shared */
	.fade-in { animation: fadeIn 0.2s ease-out; }
	@keyframes fadeIn { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: translateY(0); } }
	.center { display: flex; justify-content: center; }
	.stack { max-width: 720px; margin: 0 auto; }
	.q-card {
		background: var(--paper);
		border: 2.5px solid var(--ink);
		box-shadow: 6px 6px 0 var(--ink);
		border-radius: 18px;
		padding: clamp(1.15rem, 1rem + 0.8vw, 1.5rem);
	}
	.welcome-card { max-width: 480px; width: 100%; text-align: center; }
	.welcome-icon { display: grid; place-items: center; color: var(--ink-soft); margin: 0 auto; }
	.welcome-title { font-family: var(--font-display); font-size: clamp(1.6rem, 1.3rem + 1vw, 2.1rem); font-weight: 800; margin: 0.5rem 0 0; line-height: 1.08; letter-spacing: -0.02em; }
	.welcome-sub { color: var(--ink-soft); font-size: 0.94rem; line-height: 1.6; margin: 0.5rem 0 0; }
	.welcome-form { margin-top: 1.15rem; display: flex; flex-direction: column; gap: 0.9rem; text-align: left; }

	.field { display: flex; flex-direction: column; gap: 0.35rem; }
	.field-label { font-size: 0.78rem; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase; color: var(--ink); }
	.input {
		width: 100%; min-height: 46px; padding: 0.6rem 0.85rem;
		border-radius: 12px; border: 2.5px solid var(--ink); background: var(--cream);
		color: var(--ink); font-size: 1rem; font-weight: 500; outline: none;
	}
	.input:focus { background: var(--paper); border-color: var(--plum); box-shadow: 0 0 0 3px rgba(107,63,160,0.18); }
	.hint { color: var(--ink-soft); font-size: 0.85rem; margin: 0; text-align: center; }

	.section-head { display: flex; gap: 0.75rem; align-items: center; margin-bottom: 1rem; }
	.back-btn {
		width: 44px; height: 44px; border-radius: 12px; border: 2.5px solid var(--ink);
		background: var(--paper); color: var(--ink); display: grid; place-items: center;
		font-size: 1.15rem; cursor: pointer; flex: none;
	}
	.back-btn:hover { background: var(--cream-deep); }
	.section-title { font-family: var(--font-display); font-size: 1.32rem; font-weight: 800; margin: 0; line-height: 1.1; letter-spacing: -0.015em; }
	.section-sub { color: var(--ink-soft); font-size: 0.88rem; margin: 0.15rem 0 0; }
	.empty { text-align: center; padding: 2rem; color: var(--ink-soft); font-weight: 600; }

	.pick-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 0.8rem; }
	.pick-card {
		text-align: left; background: var(--paper); border: 2px solid var(--ink);
		box-shadow: 3px 3px 0 var(--ink); border-radius: 14px; padding: 1rem;
		cursor: pointer; display: flex; flex-direction: column; align-items: flex-start; gap: 0.45rem;
		transition: box-shadow 120ms, transform 120ms, background 120ms;
		min-width: 0; overflow: hidden;
	}
	.pick-card:hover { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--ink); background: var(--cream); }
	.pick-card:active { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
	.pick-icon { width: 36px; height: 36px; border-radius: 10px; border: 1.5px solid var(--ink); background: var(--cream); display: grid; place-items: center; color: var(--ink); flex: none; }
	.pick-name { font-family: var(--font-display); font-weight: 700; font-size: 1.02rem; line-height: 1.2; }
	.pick-meta { color: var(--ink-soft); font-size: 0.76rem; font-weight: 600; font-family: var(--font-mono); }

	.topics { display: flex; flex-wrap: wrap; gap: 0.6rem; }
	.topic {
		min-height: 40px; padding: 0.5rem 0.9rem; border-radius: 999px;
		border: 1.5px solid var(--ink); background: var(--paper); color: var(--ink);
		font-weight: 600; font-size: 0.88rem; cursor: pointer;
	}
	.topic:hover { background: var(--cream); }
	.topic-all { background: var(--amber); border-width: 2px; }

	.diff-card { max-width: 520px; width: 100%; display: flex; flex-direction: column; gap: 1rem; }
	.diff-stack { display: flex; flex-direction: column; gap: 0.75rem; }
	.diff-btn {
		display: flex; justify-content: space-between; align-items: center; gap: 0.75rem;
		min-height: 56px; padding: 0.75rem 1rem; border-radius: 12px;
		border: 2px solid var(--ink); box-shadow: 3px 3px 0 var(--ink);
		cursor: pointer; transition: box-shadow 120ms, transform 120ms;
		color: var(--ink);
	}
	.diff-btn:hover { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--ink); }
	.diff-btn:active { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
	.diff-easy { background: var(--mint); }
	.diff-medium { background: var(--amber-tint); }
	.diff-hard { background: var(--coral-tint); }
	.diff-left { display: flex; align-items: center; gap: 0.5rem; font-family: var(--font-display); font-weight: 700; font-size: 1.02rem; }
	.diff-emoji { font-size: 1.35rem; }
	.diff-meta { display: flex; align-items: center; gap: 0.6rem; font-size: 0.82rem; font-weight: 700; color: var(--ink-soft); }
	.dots { display: inline-flex; gap: 3px; }
	.dot { width: 8px; height: 8px; border-radius: 50%; border: 1.5px solid var(--ink); }
	.dot.filled { background: var(--ink); }

	/* quiz shell */
	.quiz-shell { display: grid; grid-template-columns: minmax(0, 1fr) 320px; gap: 1.25rem; align-items: start; }
	@media (max-width: 860px) { .quiz-shell { grid-template-columns: 1fr; } }
	.quiz-main { min-width: 0; }
	.q-meta { display: flex; justify-content: space-between; align-items: center; gap: 0.75rem; flex-wrap: wrap; margin-bottom: 1rem; }
	.q-tags { display: flex; flex-wrap: wrap; gap: 0.4rem; }
	.tag {
		font-size: 0.78rem; font-weight: 700; border: 2px solid var(--ink);
		background: var(--paper); border-radius: 999px; padding: 0.22rem 0.55rem;
	}
	.tag-difficulty { text-transform: capitalize; }
	.tag-easy { background: var(--mint); } .tag-medium { background: var(--amber-tint); } .tag-hard { background: var(--coral-tint); }
	.counter { font-family: var(--font-mono); font-weight: 700; color: var(--ink-soft); font-size: 0.85rem; }
	.q-fieldset { border: 0; margin: 0; padding: 0; min-inline-size: 0; }
	.q-legend { font-family: var(--font-body); font-weight: 700; font-size: clamp(1rem, 0.96rem + 0.4vw, 1.18rem); line-height: 1.45; color: var(--ink); margin: 0 0 0.9rem; width: 100%; overflow-wrap: anywhere; word-break: break-word; }
	.options { display: flex; flex-direction: column; gap: 0.55rem; min-width: 0; }
	.opt {
		display: block; border: 2px solid var(--ink); background: var(--paper);
		border-radius: 12px; cursor: pointer; position: relative;
	}
	.opt:hover { background: var(--cream-deep); }
	.opt-input { position: absolute; opacity: 0; width: 1px; height: 1px; }
	.opt-row { display: flex; align-items: center; gap: 0.7rem; min-height: 44px; padding: 0.75rem 0.9rem; }
	.opt-radio {
		width: 22px; height: 22px; border-radius: 50%; border: 2px solid var(--ink);
		background: var(--paper); display: grid; place-items: center; flex: none;
	}
	.opt-dot { width: 10px; height: 10px; border-radius: 50%; background: var(--ink); transform: scale(0); transition: transform 120ms; }
	.opt-input:checked + .opt-row .opt-dot { transform: scale(1); }
	.opt-text { flex: 1; min-width: 0; font-size: 0.96rem; line-height: 1.4; font-weight: 500; overflow-wrap: anywhere; word-break: break-word; }
	.opt-state { width: 26px; height: 26px; display: grid; place-items: center; font-weight: 800; flex: none; }
	.opt-selected { border-color: var(--plum); background: var(--amber-tint); }
	.opt-correct { border-color: var(--teal); background: var(--mint); }
	.opt-correct .opt-radio { border-color: var(--teal); }
	.opt-correct .opt-dot { background: var(--teal); }
	.opt-correct .opt-state { color: var(--teal); }
	.opt-incorrect { border-color: var(--coral); background: var(--coral-tint); }
	.opt-incorrect .opt-radio { border-color: var(--coral); }
	.opt-incorrect .opt-dot { background: var(--coral); }
	.opt-incorrect .opt-state { color: var(--coral); }
	.opt:focus-within { outline: 3px solid var(--plum); outline-offset: 3px; }
	.feedback {
		margin-top: 1rem; border: 2.5px solid var(--teal); background: var(--mint);
		border-radius: 14px; padding: 1rem; display: flex; flex-direction: column; gap: 0.6rem;
	}
	.feedback-bad { border-color: var(--coral); background: var(--coral-tint); }
	.feedback-head { display: flex; align-items: center; gap: 0.6rem; flex-wrap: wrap; }
	.feedback-badge {
		width: 32px; height: 32px; border-radius: 50%; background: var(--teal); color: var(--paper);
		display: grid; place-items: center; font-weight: 800; flex: none;
	}
	.feedback-bad .feedback-badge { background: var(--coral); }
	.feedback-title { font-family: var(--font-display); font-weight: 800; font-size: 1.15rem; }
	.streak-chip {
		margin-left: auto; font-size: 0.82rem; font-weight: 700;
		border: 2px solid var(--ink); background: var(--amber-tint); border-radius: 999px;
		padding: 0.2rem 0.55rem;
	}
	.feedback-text { margin: 0; line-height: 1.5; color: var(--ink); font-weight: 500; }

	.answer-pop { animation: pop 0.22s ease-out; }
	@keyframes pop { 0% { transform: scale(1); } 50% { transform: scale(1.015); } 100% { transform: scale(1); } }

	/* side */
	.quiz-side { display: flex; flex-direction: column; gap: 0.7rem; min-width: 0; }
	.side-card {
		background: var(--paper); border: 2px solid var(--ink);
		box-shadow: 3px 3px 0 var(--ink); border-radius: 14px; padding: 0.9rem;
		display: flex; flex-direction: column; gap: 0.65rem;
	}
	.side-muted { background: var(--cream-deep); box-shadow: 2px 2px 0 var(--ink); }
	.side-row { display: flex; justify-content: space-between; align-items: baseline; gap: 0.5rem; }
	.side-label { font-size: 0.75rem; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; color: var(--ink-soft); }
	.side-value { font-size: 1.6rem; font-weight: 800; }
	.mono { font-family: var(--font-mono); }
	.progress { height: 10px; background: var(--cream); border: 2px solid var(--ink); border-radius: 999px; overflow: hidden; }
	.progress-bar { height: 100%; background: var(--ink); border-radius: 999px; transition: width 0.3s; }
	.side-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0.5rem; text-align: center; }
	.stat-v { font-weight: 800; font-size: 1rem; display: block; }
	.stat-k { font-size: 0.7rem; font-weight: 700; letter-spacing: 0.06em; text-transform: uppercase; color: var(--ink-soft); }
	.timer-head { display: flex; justify-content: space-between; align-items: center; }
	.timer-num { font-size: 1.35rem; font-weight: 800; }
	.timer-track { height: 10px; background: var(--cream); border: 2px solid var(--ink); border-radius: 999px; overflow: hidden; }
	.timer-fill { height: 100%; border-radius: 999px; transition: width 0.1s linear; }
	.streak-box {
		text-align: center; font-weight: 800; padding: 0.45rem;
		background: var(--amber-tint); border: 2px solid var(--ink); border-radius: 999px;
	}
	.side-hint { margin: 0; font-size: 0.88rem; line-height: 1.45; color: var(--ink-soft); }

	/* results */
	.score-card { max-width: 520px; width: 100%; text-align: center; display: flex; flex-direction: column; gap: 0.75rem; align-items: center; }
	.score-kicker { font-family: var(--font-mono); font-size: 0.68rem; font-weight: 700; letter-spacing: 0.12em; text-transform: uppercase; color: var(--ink-soft); margin: 0; }
	.score-stars { font-size: 1.6rem; }
	.score-num { font-size: clamp(2.75rem, 2rem + 3vw, 4rem); font-weight: 800; line-height: 1; letter-spacing: -0.03em; }
	.score-label { font-size: 0.75rem; font-weight: 700; letter-spacing: 0.1em; text-transform: uppercase; color: var(--ink-soft); margin: -0.4rem 0 0; }
	.score-name { font-family: var(--font-display); font-size: 1.35rem; font-weight: 800; margin: 0; }
	.score-meta { color: var(--ink-soft); font-size: 0.9rem; margin: 0; }
	.score-stats { display: flex; align-items: center; gap: 1rem; justify-content: center; width: 100%; }
	.s-stat { text-align: center; flex: 1; }
	.s-v { font-weight: 800; font-size: 1.15rem; display: block; }
	.s-k { font-size: 0.68rem; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; color: var(--ink-soft); }
	.s-sep { width: 1px; height: 36px; background: var(--ink); opacity: 0.15; }
	.score-time { font-family: var(--font-mono); font-size: 0.82rem; color: var(--ink-soft); margin: 0; }
	.score-actions { width: 100%; display: flex; flex-direction: column; gap: 0.6rem; margin-top: 0.25rem; }

	/* buttons */
	.btn-primary {
		min-height: 44px; border-radius: 12px; border: 2.5px solid var(--ink);
		background: var(--amber); color: var(--ink); box-shadow: 3px 3px 0 var(--ink);
		font-family: var(--font-display); font-weight: 700; font-size: 0.98rem;
		cursor: pointer; padding: 0.55rem 1.1rem; display: inline-flex; align-items: center; justify-content: center; gap: 0.35rem;
		transition: transform 120ms, box-shadow 120ms;
	}
	.btn-primary:hover { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--ink); }
	.btn-primary:active { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
	.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; transform: none; box-shadow: 3px 3px 0 var(--ink); }
	.btn-block { width: 100%; }
	.btn-ghost {
		min-height: 48px; border-radius: 14px; border: 2.5px solid var(--ink);
		background: var(--paper); color: var(--ink); font-weight: 700; font-size: 1rem;
		cursor: pointer; padding: 0.6rem 1.2rem; font-family: var(--font-body);
	}
	.btn-ghost:hover { background: var(--cream-deep); }

	/* confetti / popups */
	.confetti-piece { position: fixed; top: -20px; width: 12px; height: 12px; border-radius: 2px; z-index: 100; pointer-events: none; animation: confettiFall 2.4s ease-out forwards; }
	@keyframes confettiFall { 0% { opacity: 1; transform: translateY(0) rotate(0deg); } 100% { opacity: 0; transform: translateY(110vh) rotate(var(--rot, 720deg)); } }
	.score-popup {
		position: fixed; z-index: 50; pointer-events: none; font-family: var(--font-mono);
		font-size: 1.6rem; font-weight: 800; color: var(--ink);
		background: var(--amber); border: 2px solid var(--ink); border-radius: 999px;
		padding: 0.15rem 0.55rem; box-shadow: 3px 3px 0 var(--ink);
		animation: scoreFloat 1.1s ease-out forwards;
	}
	@keyframes scoreFloat { 0% { opacity: 1; transform: translateY(0) scale(1); } 100% { opacity: 0; transform: translateY(-90px) scale(1.05); } }
	.shake { animation: shakeIt 0.45s ease-out; }
	@keyframes shakeIt { 0%,100% { transform: translateX(0); } 20% { transform: translateX(-6px); } 40% { transform: translateX(6px); } 60% { transform: translateX(-4px); } 80% { transform: translateX(4px); } }

	@media (prefers-reduced-motion: reduce) { .fade-in, .answer-pop, .shake, .confetti-piece, .score-popup { animation: none; transition: none; } }
</style>
