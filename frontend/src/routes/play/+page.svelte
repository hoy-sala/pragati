<script lang="ts">
	import { apiUrl } from '$lib/api/client.svelte';
	import { goto } from '$app/navigation';
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

	const SUBJECT_COLORS = ['from-violet-500 to-purple-600', 'from-blue-500 to-cyan-500', 'from-emerald-500 to-teal-500', 'from-amber-500 to-orange-500', 'from-rose-500 to-pink-500', 'from-indigo-500 to-blue-500', 'from-fuchsia-500 to-purple-500'];
	const DIFFICULTY_COLORS: Record<string, string> = { easy: 'from-emerald-400 to-green-500', medium: 'from-amber-400 to-orange-500', hard: 'from-rose-400 to-red-500' };
	const DIFFICULTY_EMOJI: Record<string, string> = { easy: '🌟', medium: '🔥', hard: '💎' };
	const OPTION_COLORS = ['from-blue-500 to-indigo-600', 'from-emerald-500 to-teal-600', 'from-amber-500 to-orange-600', 'from-rose-500 to-pink-600'];
	const CONFETTI_COLORS = ['#fbbf24', '#f472b6', '#34d399', '#60a5fa', '#a78bfa', '#fb923c'];

	function shuffle<T>(arr: T[]): T[] {
		const a = [...arr];
		for (let i = a.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));
			[a[i], a[j]] = [a[j], a[i]];
		}
		return a;
	}

	function subjectColor(i: number) { return SUBJECT_COLORS[i % SUBJECT_COLORS.length]; }

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
	function playAgain() { playClick(); phase = 'difficulty'; }

	function timeColor() { return timeLeft > 10 ? '#34d399' : timeLeft > 5 ? '#fbbf24' : '#f43f5e'; }
	function timerWidth() { return `${(timeLeft / 15) * 100}%`; }
	function progressWidth() { return questions.length ? `${((currentIndex + 1) / questions.length) * 100}%` : '0%'; }
	function accuracy() { return questions.length ? Math.round((correctCount / questions.length) * 100) : 0; }
	function formatTime(ms: number) { const s = ms / 1000; return s < 60 ? `${s.toFixed(1)}s` : `${Math.floor(s / 60)}m ${Math.floor(s % 60)}s`; }
	function stars() { return accuracy() >= 80 ? '⭐⭐⭐' : accuracy() >= 50 ? '⭐⭐' : '⭐'; }
</script>

<svelte:head>
	<title>Quiz Arena</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
</svelte:head>

<!-- confetti particles -->
{#each confettiPieces as cp (cp.id)}
	<div class="confetti-piece" style="left:{cp.x}%;background:{cp.color};animation-delay:{cp.delay}s;--rot:{cp.rot}deg"></div>
{/each}

<!-- score popups -->
{#each scorePopups as popup (popup.id)}
	<div class="score-popup" style="left:{popup.x}%;top:{popup.y}%">+{popup.value}⚡</div>
{/each}

<div class="play-root {screenShake ? 'shake' : ''}">
	<!-- exit confirm modal -->
	{#if showExitConfirm}
		<div class="fixed inset-0 z-[60] flex items-center justify-center bg-black/60 backdrop-blur-sm fade-in">
			<div class="glass rounded-3xl p-8 max-w-sm w-full mx-4 text-center space-y-5 fade-in scale-in">
				<div class="text-6xl bounce-in">🚪</div>
				<h3 class="text-2xl font-black text-white">Quit Quiz?</h3>
				<p class="text-white/60">Your progress will be lost!</p>
				<div class="space-y-3">
					<button onclick={confirmGoHome} class="btn-game w-full py-4 rounded-2xl bg-gradient-to-r from-rose-500 to-pink-600 text-white font-black text-lg shadow-xl">Yes, Quit</button>
					<button onclick={cancelExit} class="w-full py-4 rounded-2xl bg-white/10 text-white font-bold text-lg border border-white/20 hover:bg-white/20 active:scale-95 transition-all">Keep Playing</button>
				</div>
			</div>
		</div>
	{/if}

	{#if phase !== 'results'}
		<button onclick={confirmExit} class="fixed top-4 left-4 z-50 px-4 py-2 rounded-2xl bg-white/10 backdrop-blur border border-white/20 text-white/70 text-sm font-bold hover:bg-white/20 hover:text-white transition-all flex items-center gap-2 active:scale-95">
			🏠 Home
		</button>
	{/if}

	<!-- ═══ WELCOME ═══ -->
	{#if phase === 'welcome'}
		<div class="quiz-container fade-in">
			<div class="glass rounded-[2rem] p-8 md:p-12 text-center space-y-6 max-w-2xl mx-auto">
				<div class="text-8xl md:text-9xl bounce-in">🎮</div>
				<h1 class="text-5xl md:text-7xl font-black text-white tracking-tight game-title">Quiz Arena</h1>
				<p class="text-white/60 text-xl md:text-2xl font-medium">Test your knowledge & have fun!</p>
				<div class="space-y-4 max-w-md mx-auto">
					<input bind:value={playerName} placeholder="Enter your name..." maxlength={50}
						class="w-full px-6 py-4 rounded-2xl bg-white/10 border-2 border-white/20 text-white text-center text-xl placeholder-white/40 focus:outline-none focus:ring-4 focus:ring-purple-400/50 focus:border-purple-400/50 transition-all" />
					<button onclick={loadClasses} disabled={!playerName.trim()}
						class="btn-game w-full py-5 rounded-2xl bg-gradient-to-r from-amber-400 to-orange-500 text-white font-black text-xl shadow-xl shadow-orange-500/30 hover:shadow-orange-500/50 active:scale-95 transition-all disabled:opacity-30 disabled:cursor-not-allowed">
						🚀 Start Playing
					</button>
				</div>
			</div>
		</div>

	<!-- ═══ CLASSES ═══ -->
	{:else if phase === 'classes'}
		<div class="quiz-container fade-in">
			<div class="glass rounded-[2rem] p-6 md:p-8 space-y-5 max-w-3xl mx-auto">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-3 rounded-xl hover:bg-white/10 text-white/70 hover:text-white transition-colors text-xl">←</button>
					<div><h2 class="text-2xl md:text-3xl font-black text-white">Select Class</h2><p class="text-white/50">{playerName}'s turn</p></div>
				</div>
				{#if loading}
					<div class="text-center py-12 text-white/50 text-xl">Loading...</div>
				{:else if classes.length === 0}
					<div class="text-center py-12 text-white/50 text-xl">No classes with questions yet</div>
				{:else}
					<div class="grid grid-cols-2 md:grid-cols-3 gap-4">
						{#each classes as cls, i (cls.id)}
							<button onclick={() => loadSubjects(cls)}
								class="btn-game group p-5 md:p-6 rounded-2xl bg-gradient-to-br {subjectColor(i)} text-white text-left hover:scale-[1.04] active:scale-95 transition-transform shadow-xl">
								<div class="text-3xl mb-2">{['📚','📖','✏️','🎒','📝'][i % 5]}</div>
								<div class="text-xl md:text-2xl font-black">{cls.name}</div>
								<div class="text-white/70 text-sm mt-1">{cls.question_count} questions</div>
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>

	<!-- ═══ SUBJECTS ═══ -->
	{:else if phase === 'subjects'}
		<div class="quiz-container fade-in">
			<div class="glass rounded-[2rem] p-6 md:p-8 space-y-5 max-w-3xl mx-auto">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-3 rounded-xl hover:bg-white/10 text-white/70 hover:text-white transition-colors text-xl">←</button>
					<div><h2 class="text-2xl md:text-3xl font-black text-white">Select Subject</h2><p class="text-white/50">{selectedClass?.name}</p></div>
				</div>
				{#if loading}
					<div class="text-center py-12 text-white/50 text-xl">Loading...</div>
				{:else}
					<div class="grid grid-cols-2 md:grid-cols-3 gap-4">
						{#each subjects as sub, i (sub.id)}
							<button onclick={() => loadTopics(sub)}
								class="btn-game group p-5 rounded-2xl bg-gradient-to-br {subjectColor(i)} text-white text-left hover:scale-[1.04] active:scale-95 transition-transform shadow-xl">
								<div class="text-3xl mb-2">{['🔢','🔬','📖','🗣️','📝','🌍','🧠','🎨','🎵','⚽','💻'][i % 11]}</div>
								<div class="text-lg md:text-xl font-black">{sub.name}</div>
								<div class="text-white/70 text-sm mt-1">{sub.question_count} questions</div>
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>

	<!-- ═══ TOPICS ═══ -->
	{:else if phase === 'topics'}
		<div class="quiz-container fade-in">
			<div class="glass rounded-[2rem] p-6 md:p-8 space-y-5 max-w-2xl mx-auto">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-3 rounded-xl hover:bg-white/10 text-white/70 hover:text-white transition-colors text-xl">←</button>
					<div><h2 class="text-2xl md:text-3xl font-black text-white">Select Topic</h2><p class="text-white/50">{selectedSubject?.name}</p></div>
				</div>
				{#if loading}
					<div class="text-center py-12 text-white/50 text-xl">Loading...</div>
				{:else}
					<div class="flex flex-wrap gap-3">
						<button onclick={() => { playClick(); selectedTopic = ''; phase = 'difficulty'; }}
							class="btn-game px-6 py-3 rounded-2xl bg-gradient-to-r from-violet-500 to-purple-600 text-white font-bold text-lg hover:scale-[1.04] active:scale-95 transition-all shadow-lg">
							🎯 All Topics
						</button>
						{#each topics as topic (topic.name)}
							<button onclick={() => { playClick(); selectedTopic = topic.name; phase = 'difficulty'; }}
								class="btn-game px-6 py-3 rounded-2xl bg-white/10 text-white font-bold text-lg hover:bg-white/20 active:scale-95 transition-all border-2 border-white/20 hover:border-white/40">
								{topic.name}
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>

	<!-- ═══ DIFFICULTY ═══ -->
	{:else if phase === 'difficulty'}
		<div class="quiz-container fade-in">
			<div class="glass rounded-[2rem] p-6 md:p-8 space-y-5 max-w-lg mx-auto">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-3 rounded-xl hover:bg-white/10 text-white/70 hover:text-white transition-colors text-xl">←</button>
					<div>
						<h2 class="text-2xl md:text-3xl font-black text-white">Choose Difficulty</h2>
						<p class="text-white/50">{selectedClass?.name} → {selectedSubject?.name}{selectedTopic ? ` → ${selectedTopic}` : ''}</p>
					</div>
				</div>
				{#if loading}
					<div class="text-center py-12 text-white/50 text-xl">Loading...</div>
				{:else}
					<div class="space-y-4">
						{#each ['easy', 'medium', 'hard'] as diff (diff)}
							<button onclick={() => startQuiz(diff)}
								class="btn-game w-full py-5 md:py-6 rounded-2xl bg-gradient-to-r {DIFFICULTY_COLORS[diff]} text-white font-black text-xl md:text-2xl shadow-xl hover:scale-[1.03] active:scale-95 transition-transform">
								<span class="text-3xl">{DIFFICULTY_EMOJI[diff]}</span> {diff.charAt(0).toUpperCase() + diff.slice(1)}
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>

	<!-- ═══ QUIZ (landscape layout) ═══ -->
	{:else if phase === 'quiz' && questions.length > 0}
		{@const q = questions[currentIndex]}
		<div class="quiz-container fade-in">
			<div class="quiz-layout">
				<!-- LEFT: Score, Timer, Progress -->
				<div class="quiz-sidebar space-y-4">
					<!-- Score -->
					<div class="glass rounded-2xl p-4 text-center">
						<div class="text-5xl md:text-6xl font-black text-white game-title">{score}</div>
						<div class="text-white/50 text-xs font-bold uppercase tracking-widest">Points</div>
					</div>
					<!-- Streak -->
					{#if streak >= 2}
						<div class="glass rounded-2xl p-3 text-center animate-pulse">
							<div class="text-3xl font-black text-amber-400">🔥 {streak}x</div>
							<div class="text-white/50 text-xs font-bold">STREAK</div>
						</div>
					{/if}
					<!-- Timer ring -->
					<div class="glass rounded-2xl p-4 flex flex-col items-center gap-2">
						<div class="timer-ring-wrap">
							<svg class="timer-ring" viewBox="0 0 100 100">
								<circle cx="50" cy="50" r="42" fill="none" stroke="rgba(255,255,255,0.1)" stroke-width="8" />
								<circle cx="50" cy="50" r="42" fill="none" stroke={timeColor()} stroke-width="8" stroke-linecap="round"
									stroke-dasharray="264" stroke-dashoffset={264 - (timeLeft / 15) * 264}
									transform="rotate(-90 50 50)" class="timer-arc" />
							</svg>
							<div class="timer-text" style="color:{timeColor()}">{Math.ceil(timeLeft)}</div>
						</div>
					</div>
					<!-- Progress -->
					<div class="glass rounded-2xl p-3">
						<div class="flex justify-between text-white/60 text-xs font-bold mb-2">
							<span>Question</span>
							<span>{currentIndex + 1}/{questions.length}</span>
						</div>
						<div class="h-3 bg-white/10 rounded-full overflow-hidden">
							<div class="h-full bg-gradient-to-r from-purple-400 to-pink-400 rounded-full transition-all duration-500" style="width:{progressWidth()}"></div>
						</div>
						<div class="mt-2 text-center text-2xl">{correctCount > 0 ? '⭐'.repeat(Math.min(correctCount, 5)) : '❓'}</div>
					</div>
				</div>

				<!-- RIGHT: Question + Options -->
				<div class="quiz-main space-y-4">
					<div class="glass rounded-2xl p-5 md:p-6 {answerBounce ? 'answer-flash' : ''}">
						<p class="text-white text-xl md:text-2xl lg:text-3xl leading-relaxed font-bold">{q.question_text}</p>
					</div>

					<div class="quiz-options">
						{#each q.options as opt, oi (opt.key)}
							<button onclick={() => handleAnswer(opt.key)} disabled={answered}
								class="option-btn p-4 md:p-5 rounded-2xl text-left font-bold transition-all active:scale-[0.97] shadow-lg
								{answered && opt.key === selectedKey
									? (isCorrect ? 'option-correct' : 'option-wrong')
									: answered && opt.correct
										? 'option-correct'
										: 'bg-gradient-to-br ' + OPTION_COLORS[oi % 4] + ' text-white hover:scale-[1.02] border-2 border-white/10 hover:border-white/30'}">
								<div class="flex items-center gap-4">
									<span class="option-label">{opt.key}</span>
									<span class="text-base md:text-lg leading-snug">{opt.value}</span>
									{#if answered && opt.key === selectedKey}
										<span class="text-2xl ml-auto">{isCorrect ? '✅' : '❌'}</span>
									{:else if answered && opt.correct}
										<span class="text-2xl ml-auto">✅</span>
									{/if}
								</div>
							</button>
						{/each}
					</div>
				</div>
			</div>
		</div>

	<!-- ═══ RESULTS (selfie scorecard) ═══ -->
	{:else if phase === 'results'}
		<div class="quiz-container fade-in">
			<div class="scorecard rounded-[2rem] p-6 md:p-8 text-center space-y-5 max-w-lg mx-auto relative overflow-hidden">
				<div class="absolute inset-0 rounded-[2rem] scorecard-border pointer-events-none"></div>
				<div class="relative z-10">
					<p class="text-[10px] font-bold text-white/40 uppercase tracking-[0.25em]">Karnataka Residential Educational Institutions Society</p>
					<h2 class="text-base md:text-lg font-black text-white mt-1" style="text-shadow:0 1px 4px rgba(0,0,0,0.2)">MDRS (SC-32) BAHADDURGHATTA</h2>
				</div>
				<div class="text-4xl relative z-10 bounce-in">{stars()}</div>
				<div class="relative z-10">
					<div class="text-7xl md:text-8xl font-black text-white score-glow game-title">{score}</div>
					<p class="text-white/50 text-xs font-bold uppercase tracking-widest mt-1">Points</p>
				</div>
				<div class="relative z-10">
					<p class="text-xl md:text-2xl font-black text-white">{playerName}</p>
					<p class="text-white/40 text-sm">{selectedClass?.name} · {selectedSubject?.name}</p>
				</div>
				<div class="flex justify-center gap-6 relative z-10">
					<div class="text-center">
						<div class="text-2xl font-black text-emerald-400">{accuracy()}%</div>
						<div class="text-white/40 text-xs font-bold">ACCURACY</div>
					</div>
					<div class="w-px bg-white/10"></div>
					<div class="text-center">
						<div class="text-2xl font-black text-amber-400">🔥 {bestStreak}</div>
						<div class="text-white/40 text-xs font-bold">STREAK</div>
					</div>
					<div class="w-px bg-white/10"></div>
					<div class="text-center">
						<div class="text-2xl font-black text-blue-400">{correctCount}/{questions.length}</div>
						<div class="text-white/40 text-xs font-bold">CORRECT</div>
					</div>
				</div>
				<div class="relative z-10">
					<span class="inline-block px-4 py-1.5 rounded-full bg-white/10 text-white/60 text-xs font-bold uppercase tracking-wider">
						{selectedDifficulty} · {formatTime(Date.now() - startTime)}
					</span>
				</div>
				<div class="flex items-center gap-3 relative z-10">
					<div class="flex-1 h-px bg-white/10"></div>
					<span class="text-white/30 text-lg">📸</span>
					<div class="flex-1 h-px bg-white/10"></div>
				</div>
				<p class="text-white/30 text-xs font-bold relative z-10">Take a selfie with your score!</p>
				<div class="space-y-3 relative z-10 pt-1">
					<button onclick={playAgain} class="btn-game w-full py-4 rounded-2xl bg-gradient-to-r from-amber-400 to-orange-500 text-white font-black text-lg shadow-xl shadow-orange-500/30 active:scale-95 transition-all">
						🎮 Play Again
					</button>
					<a href="/" class="w-full py-4 rounded-2xl bg-white/10 text-white font-bold text-lg border border-white/20 hover:bg-white/20 active:scale-95 transition-all flex items-center justify-center gap-2">
						🏠 Back to Home
					</a>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	:global(body) { background: #1e1b4b !important; margin: 0; overflow: hidden; }
	.play-root { min-height: 100vh; min-height: 100dvh; background: linear-gradient(135deg, #312e81 0%, #7c3aed 30%, #c026d3 60%, #1e1b4b 100%); display: flex; align-items: center; justify-content: center; padding: 1rem; overflow-y: auto; -webkit-tap-highlight-color: transparent; position: relative; }

	.quiz-container { width: 100%; max-width: 100%; }

	.glass { background: rgba(255,255,255,0.08); backdrop-filter: blur(20px); border: 2px solid rgba(255,255,255,0.12); }

	.fade-in { animation: fadeSlideIn 0.35s ease-out; }
	@keyframes fadeSlideIn { from { opacity: 0; transform: translateY(16px) scale(0.98); } to { opacity: 1; transform: translateY(0) scale(1); } }
	.scale-in { animation: scaleIn 0.3s ease-out; }
	@keyframes scaleIn { from { opacity: 0; transform: scale(0.85); } to { opacity: 1; transform: scale(1); } }
	.bounce-in { animation: bounceIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1); }
	@keyframes bounceIn { from { transform: scale(0); } to { transform: scale(1); } }

	.game-title { background: linear-gradient(135deg, #fbbf24, #f472b6, #a78bfa); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; }

	.btn-game { position: relative; overflow: hidden; }
	.btn-game::after { content: ''; position: absolute; inset: 0; background: linear-gradient(to bottom, rgba(255,255,255,0.15), transparent); pointer-events: none; border-radius: inherit; }

	/* shake */
	.shake { animation: shakeIt 0.5s ease-out; }
	@keyframes shakeIt { 0%,100% { transform: translateX(0); } 20% { transform: translateX(-8px); } 40% { transform: translateX(8px); } 60% { transform: translateX(-4px); } 80% { transform: translateX(4px); } }

	/* score popups */
	.score-popup { position: fixed; z-index: 50; pointer-events: none; font-size: 2rem; font-weight: 900; color: #fbbf24; text-shadow: 0 0 20px rgba(251,191,36,0.6), 0 0 40px rgba(251,191,36,0.3); animation: scoreFloat 1.2s ease-out forwards; }
	@keyframes scoreFloat { 0% { opacity: 1; transform: translateY(0) scale(1); } 100% { opacity: 0; transform: translateY(-100px) scale(1.5); } }

	/* confetti */
	.confetti-piece { position: fixed; top: -20px; width: 12px; height: 12px; border-radius: 2px; z-index: 100; pointer-events: none; animation: confettiFall 2.5s ease-out forwards; }
	@keyframes confettiFall { 0% { opacity: 1; transform: translateY(0) rotate(0deg); } 100% { opacity: 0; transform: translateY(110vh) rotate(var(--rot, 720deg)); } }

	/* ═══ LAYOUT ═══ */
	.quiz-layout { display: flex; gap: 1rem; max-width: 1200px; margin: 0 auto; }
	.quiz-sidebar { width: 200px; flex-shrink: 0; display: flex; flex-direction: column; gap: 0.75rem; }
	.quiz-main { flex: 1; display: flex; flex-direction: column; gap: 0.75rem; min-width: 0; }
	.quiz-options { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; flex: 1; }

	/* option buttons */
	.option-btn { min-height: 80px; display: flex; align-items: center; }
	.option-label { width: 2.5rem; height: 2.5rem; border-radius: 12px; background: rgba(255,255,255,0.2); display: flex; align-items: center; justify-content: center; font-size: 1rem; font-weight: 900; flex-shrink: 0; }
	.option-correct { background: linear-gradient(135deg, #10b981, #059669) !important; color: white; border: 3px solid #34d399; animation: correctPop 0.3s ease-out; }
	.option-wrong { background: linear-gradient(135deg, #f43f5e, #e11d48) !important; color: white; border: 3px solid #fb7185; animation: wrongShake 0.4s ease-out; }
	@keyframes correctPop { 0% { transform: scale(1); } 50% { transform: scale(1.05); } 100% { transform: scale(1); } }
	@keyframes wrongShake { 0%,100% { transform: translateX(0); } 25% { transform: translateX(-6px); } 75% { transform: translateX(6px); } }
	.answer-flash { animation: answerFlash 0.3s ease-out; }
	@keyframes answerFlash { 0% { box-shadow: 0 0 0 0 rgba(255,255,255,0.3); } 100% { box-shadow: 0 0 0 0 transparent; } }

	/* timer ring */
	.timer-ring-wrap { position: relative; width: 100px; height: 100px; }
	.timer-ring { width: 100%; height: 100%; }
	.timer-arc { transition: stroke-dashoffset 0.1s linear; }
	.timer-text { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; font-size: 2rem; font-weight: 900; }

	/* scorecard */
	.scorecard { background: rgba(255,255,255,0.08); backdrop-filter: blur(20px); border: 2px solid rgba(255,255,255,0.15); box-shadow: 0 0 60px rgba(139,92,246,0.3), inset 0 1px 0 rgba(255,255,255,0.1); }
	.scorecard-border { border: 2px solid transparent; background: linear-gradient(135deg, rgba(251,191,36,0.2), rgba(139,92,246,0.2), rgba(236,72,153,0.2)) border-box; -webkit-mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0); -webkit-mask-composite: xor; mask-composite: exclude; }
	.score-glow { text-shadow: 0 0 40px rgba(251,191,36,0.5), 0 4px 12px rgba(0,0,0,0.4); }

	/* ═══ RESPONSIVE ═══ */
	@media (max-width: 768px) {
		.quiz-layout { flex-direction: column; }
		.quiz-sidebar { width: 100%; flex-direction: row; flex-wrap: wrap; gap: 0.5rem; }
		.quiz-sidebar > div { flex: 1; min-width: 100px; }
		.timer-ring-wrap { width: 60px; height: 60px; }
		.timer-text { font-size: 1.2rem; }
		.quiz-options { grid-template-columns: 1fr; }
		.option-btn { min-height: 60px; }
	}
</style>
