<script lang="ts">
	import { apiUrl } from '$lib/api/client.svelte';
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
	let scorePopups = $state<{ id: number; value: number; x: number; y: number }[]>([]);
	let popupId = $state(0);
	let showExitConfirm = $state(false);

	const SUBJECT_COLORS = ['from-violet-500 to-purple-600', 'from-blue-500 to-cyan-500', 'from-emerald-500 to-teal-500', 'from-amber-500 to-orange-500', 'from-rose-500 to-pink-500', 'from-indigo-500 to-blue-500', 'from-fuchsia-500 to-purple-500'];
	const DIFFICULTY_COLORS: Record<string, string> = { easy: 'from-emerald-400 to-green-500', medium: 'from-amber-400 to-orange-500', hard: 'from-rose-400 to-red-500' };

	function shuffle<T>(arr: T[]): T[] {
		const a = [...arr];
		for (let i = a.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));
			[a[i], a[j]] = [a[j], a[i]];
		}
		return a;
	}

	function subjectColor(i: number) { return SUBJECT_COLORS[i % SUBJECT_COLORS.length]; }

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
		loading = true;
		classes = (await api<PlayClass[]>('GET', '/play/classes')) ?? [];
		loading = false;
		phase = 'classes';
	}

	async function loadSubjects(cls: PlayClass) {
		selectedClass = cls;
		loading = true;
		subjects = (await api<PlaySubject[]>('GET', `/play/subjects?class_id=${cls.id}`)) ?? [];
		loading = false;
		phase = 'subjects';
	}

	async function loadTopics(sub: PlaySubject) {
		selectedSubject = sub;
		loading = true;
		topics = (await api<PlayTopic[]>('GET', `/play/topics?class_id=${selectedClass!.id}&subject_id=${sub.id}`)) ?? [];
		loading = false;
		phase = 'topics';
	}

	async function startQuiz(difficulty: string) {
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
		} else { streak = 0; }
		clearInterval(timerInterval);
	}

	function addPopup(value: number) {
		const id = ++popupId;
		scorePopups = [...scorePopups, { id, value, x: 50 + (Math.random() - 0.5) * 30, y: 40 }];
		setTimeout(() => { scorePopups = scorePopups.filter(p => p.id !== id); }, 1200);
	}

	function nextQuestion() {
		if (currentIndex >= questions.length - 1) { finishQuiz(); return; }
		currentIndex++; selectedKey = ''; answered = false; startTimer();
	}

	function finishQuiz() {
		clearInterval(timerInterval);
		phase = 'results';
	}

	function goBack() {
		const back: Partial<Record<Phase, Phase>> = {
			classes: 'welcome', subjects: 'classes', topics: 'subjects',
			difficulty: 'topics', quiz: 'welcome', results: 'welcome'
		};
		phase = back[phase] ?? 'welcome';
	}

	function confirmExit() { showExitConfirm = true; }
	function cancelExit() { showExitConfirm = false; }
	function confirmGoHome() { showExitConfirm = false; clearInterval(timerInterval); phase = 'welcome'; }

	function playAgain() { phase = 'difficulty'; }

	function timeColor() { return timeLeft > 10 ? 'bg-emerald-400' : timeLeft > 5 ? 'bg-amber-400' : 'bg-rose-500'; }
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

<div class="play-root min-h-screen bg-gradient-to-br from-indigo-600 via-purple-600 to-violet-800 flex items-center justify-center p-4 overflow-hidden">
	{#if phase !== 'results'}
		<a href="/" class="fixed top-4 left-4 z-50 px-3 py-1.5 rounded-full bg-white/10 backdrop-blur border border-white/20 text-white/70 text-xs font-bold hover:bg-white/20 hover:text-white transition-all flex items-center gap-1.5">
			🏠 Home
		</a>
	{/if}

	{#each scorePopups as popup (popup.id)}
		<div class="score-popup" style="left: {popup.x}%; top: {popup.y}%">+{popup.value}</div>
	{/each}

	{#if showExitConfirm}
		<div class="fixed inset-0 z-[60] flex items-center justify-center bg-black/60 backdrop-blur-sm fade-in">
			<div class="glass rounded-2xl p-6 max-w-xs w-full mx-4 text-center space-y-4 fade-in">
				<div class="text-4xl">🚪</div>
				<h3 class="text-xl font-bold text-white">Quit Quiz?</h3>
				<p class="text-white/60 text-sm">Your current progress will be lost.</p>
				<div class="space-y-2">
					<button onclick={confirmGoHome}
						class="w-full py-2.5 rounded-xl bg-rose-500 text-white font-bold text-sm hover:bg-rose-600 active:scale-95 transition-all">
						Yes, Quit
					</button>
					<button onclick={cancelExit}
						class="w-full py-2.5 rounded-xl bg-white/10 text-white font-medium text-sm border border-white/20 hover:bg-white/20 active:scale-95 transition-all">
						Continue Playing
					</button>
				</div>
			</div>
		</div>
	{/if}

	<div class="w-full max-w-lg">
		{#if phase === 'welcome'}
			<div class="glass rounded-3xl p-8 text-center space-y-6 fade-in">
				<div class="text-6xl mb-2">🎮</div>
				<h1 class="text-4xl font-black text-white tracking-tight">Quiz Arena</h1>
				<p class="text-white/70 text-lg">Test your knowledge. Have fun!</p>
				<input bind:value={playerName} placeholder="Enter your name" maxlength={50}
					class="w-full px-5 py-3.5 rounded-xl bg-white/10 border border-white/20 text-white text-center text-lg placeholder-white/40 focus:outline-none focus:ring-2 focus:ring-white/50" />
				<button onclick={loadClasses} disabled={!playerName.trim()}
					class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-400 to-orange-500 text-white font-bold text-lg shadow-lg shadow-orange-500/30 hover:shadow-orange-500/50 active:scale-95 transition-all disabled:opacity-30 disabled:cursor-not-allowed">
					Start Playing
				</button>
			</div>

		{:else if phase === 'classes'}
			<div class="glass rounded-3xl p-6 space-y-4 fade-in">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-2 rounded-lg hover:bg-white/10 text-white/70 hover:text-white transition-colors text-lg">←</button>
					<div><h2 class="text-xl font-bold text-white">Select Class</h2><p class="text-white/50 text-sm">{playerName}'s turn</p></div>
				</div>
				{#if loading}
					<div class="text-center py-8 text-white/50">Loading classes...</div>
				{:else if classes.length === 0}
					<div class="text-center py-8 text-white/50">No classes with questions available</div>
				{:else}
					<div class="grid grid-cols-2 gap-3">
						{#each classes as cls, i (cls.id)}
							<button onclick={() => loadSubjects(cls)}
								class="card-btn group p-4 rounded-2xl bg-gradient-to-br {subjectColor(i)} text-white text-left hover:scale-[1.03] active:scale-95 transition-transform shadow-lg">
								<div class="text-lg font-bold">{cls.name}</div>
								<div class="text-white/70 text-xs mt-1">{cls.question_count} questions</div>
							</button>
						{/each}
					</div>
				{/if}
			</div>

		{:else if phase === 'subjects'}
			<div class="glass rounded-3xl p-6 space-y-4 fade-in">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-2 rounded-lg hover:bg-white/10 text-white/70 hover:text-white transition-colors text-lg">←</button>
					<div><h2 class="text-xl font-bold text-white">Select Subject</h2><p class="text-white/50 text-sm">{selectedClass?.name}</p></div>
				</div>
				{#if loading}
					<div class="text-center py-8 text-white/50">Loading subjects...</div>
				{:else}
					<div class="grid grid-cols-2 gap-3">
						{#each subjects as sub, i (sub.id)}
							<button onclick={() => loadTopics(sub)}
								class="card-btn group p-4 rounded-2xl bg-gradient-to-br {subjectColor(i)} text-white text-left hover:scale-[1.03] active:scale-95 transition-transform shadow-lg">
								<div class="text-lg font-bold">{sub.name}</div>
								<div class="text-white/70 text-xs mt-1">{sub.question_count} questions</div>
							</button>
						{/each}
					</div>
				{/if}
			</div>

		{:else if phase === 'topics'}
			<div class="glass rounded-3xl p-6 space-y-4 fade-in">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-2 rounded-lg hover:bg-white/10 text-white/70 hover:text-white transition-colors text-lg">←</button>
					<div><h2 class="text-xl font-bold text-white">Select Topic</h2><p class="text-white/50 text-sm">{selectedSubject?.name}</p></div>
				</div>
				{#if loading}
					<div class="text-center py-8 text-white/50">Loading topics...</div>
				{:else}
					<div class="flex flex-wrap gap-2">
						<button onclick={() => { selectedTopic = ''; phase = 'difficulty'; }}
							class="px-4 py-2 rounded-xl bg-white/10 text-white font-medium hover:bg-white/20 active:scale-95 transition-all border border-white/20">
							All Topics
						</button>
						{#each topics as topic (topic.name)}
							<button onclick={() => { selectedTopic = topic.name; phase = 'difficulty'; }}
								class="px-4 py-2 rounded-xl bg-white/10 text-white font-medium hover:bg-white/20 active:scale-95 transition-all border border-white/20">
								{topic.name}
							</button>
						{/each}
					</div>
				{/if}
			</div>

		{:else if phase === 'difficulty'}
			<div class="glass rounded-3xl p-6 space-y-4 fade-in">
				<div class="flex items-center gap-3">
					<button onclick={goBack} class="p-2 rounded-lg hover:bg-white/10 text-white/70 hover:text-white transition-colors text-lg">←</button>
					<div>
						<h2 class="text-xl font-bold text-white">Choose Difficulty</h2>
						<p class="text-white/50 text-sm">{selectedClass?.name} → {selectedSubject?.name}{selectedTopic ? ` → ${selectedTopic}` : ''}</p>
					</div>
				</div>
				{#if loading}
					<div class="text-center py-8 text-white/50">Loading questions...</div>
				{:else}
					<div class="space-y-3">
						{#each ['easy', 'medium', 'hard'] as diff (diff)}
							<button onclick={() => startQuiz(diff)}
								class="w-full py-4 rounded-2xl bg-gradient-to-r {DIFFICULTY_COLORS[diff]} text-white font-bold text-lg shadow-lg hover:scale-[1.02] active:scale-95 transition-transform">
								{diff === 'easy' ? '⭐' : diff === 'medium' ? '⭐⭐' : '⭐⭐⭐'} {diff.charAt(0).toUpperCase() + diff.slice(1)}
							</button>
						{/each}
					</div>
				{/if}
			</div>

		{:else if phase === 'quiz' && questions.length > 0}
			{@const q = questions[currentIndex]}
			<div class="space-y-3 fade-in">
				<div class="glass rounded-2xl p-4">
					<div class="flex items-center justify-between mb-2">
						<div class="flex items-center gap-2">
							<span class="text-white font-bold text-lg">⚡ {score}</span>
							{#if streak >= 2}
								<span class="px-2 py-0.5 rounded-full bg-amber-400/20 text-amber-300 text-xs font-bold animate-pulse">🔥 {streak}x</span>
							{/if}
						</div>
						<span class="text-white/60 text-xs">{currentIndex + 1}/{questions.length}</span>
					</div>
					<div class="h-1.5 bg-white/10 rounded-full overflow-hidden mb-2">
						<div class="h-full bg-white/40 rounded-full transition-all duration-300" style="width: {progressWidth()}"></div>
					</div>
					<div class="h-2 rounded-full overflow-hidden">
						<div class="h-full rounded-full transition-all duration-100 {timeColor()}" style="width: {timerWidth()}"></div>
					</div>
				</div>

				<div class="glass rounded-2xl p-6">
					<p class="text-white text-lg leading-relaxed font-medium">{q.question_text}</p>
				</div>

				<div class="space-y-2.5">
					{#each q.options as opt (opt.key)}
						<button onclick={() => handleAnswer(opt.key)} disabled={answered}
							class="w-full p-4 rounded-2xl text-left font-medium transition-all active:scale-[0.97] shadow-md
							{answered && opt.key === selectedKey
								? (isCorrect ? 'bg-emerald-500 text-white ring-4 ring-emerald-300' : 'bg-rose-500 text-white ring-4 ring-rose-300')
								: answered && opt.correct
									? 'bg-emerald-500/80 text-white'
									: 'bg-white/10 text-white hover:bg-white/20 border border-white/10'}">
							<div class="flex items-center gap-3">
								<span class="w-8 h-8 rounded-lg bg-white/10 flex items-center justify-center text-sm font-bold shrink-0
									{answered && opt.key === selectedKey ? 'bg-white/20' : ''}">
									{#if answered && opt.key === selectedKey}
										{isCorrect ? '✓' : '✗'}
									{:else}
										{opt.key}
									{/if}
								</span>
								<span class="text-sm leading-snug">{opt.value}</span>
							</div>
						</button>
					{/each}
				</div>

				{#if answered}
					<button onclick={nextQuestion}
						class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-400 to-orange-500 text-white font-bold text-base shadow-lg shadow-orange-500/30 hover:shadow-orange-500/50 active:scale-95 transition-all fade-in">
						{currentIndex >= questions.length - 1 ? '🏆 Finish' : 'Next →'}
					</button>
				{/if}
			</div>

		{:else if phase === 'results'}
			<div class="fade-in">
				<div class="scorecard rounded-3xl p-6 text-center space-y-4 relative overflow-hidden">
					<div class="absolute inset-0 rounded-3xl scorecard-border pointer-events-none"></div>

					<div class="relative z-10">
						<p class="text-[10px] font-bold text-white/40 uppercase tracking-[0.25em]">Karnataka Residential Educational Institutions Society</p>
						<h2 class="text-base font-black text-white mt-1" style="text-shadow:0 1px 4px rgba(0,0,0,0.2)">MDRS (SC-32) BAHADDURGHATTA</h2>
					</div>

					<div class="text-3xl relative z-10">{stars()}</div>

					<div class="relative z-10">
						<div class="text-6xl font-black text-white score-glow">{score}</div>
						<p class="text-white/50 text-xs font-bold uppercase tracking-widest mt-1">Points</p>
					</div>

					<div class="relative z-10">
						<p class="text-xl font-black text-white">{playerName}</p>
						<p class="text-white/40 text-xs">{selectedClass?.name} · {selectedSubject?.name}</p>
					</div>

					<div class="flex justify-center gap-4 relative z-10">
						<div class="text-center">
							<div class="text-xl font-bold text-emerald-400">{accuracy()}%</div>
							<div class="text-white/40 text-[10px] font-bold">ACCURACY</div>
						</div>
						<div class="w-px bg-white/10"></div>
						<div class="text-center">
							<div class="text-xl font-bold text-amber-400">🔥 {bestStreak}</div>
							<div class="text-white/40 text-[10px] font-bold">STREAK</div>
						</div>
						<div class="w-px bg-white/10"></div>
						<div class="text-center">
							<div class="text-xl font-bold text-blue-400">{correctCount}/{questions.length}</div>
							<div class="text-white/40 text-[10px] font-bold">CORRECT</div>
						</div>
					</div>

					<div class="relative z-10">
						<span class="inline-block px-3 py-1 rounded-full bg-white/10 text-white/60 text-[10px] font-bold uppercase tracking-wider">
							{selectedDifficulty} · {formatTime(Date.now() - startTime)}
						</span>
					</div>

					<div class="flex items-center gap-2 relative z-10">
						<div class="flex-1 h-px bg-white/10"></div>
						<span class="text-white/20 text-xs">📸</span>
						<div class="flex-1 h-px bg-white/10"></div>
					</div>

					<p class="text-white/30 text-[10px] font-bold relative z-10">Take a selfie with your score!</p>

					<div class="space-y-2 relative z-10 pt-1">
						<button onclick={playAgain}
							class="w-full py-3 rounded-2xl bg-gradient-to-r from-amber-400 to-orange-500 text-white font-bold text-sm shadow-lg shadow-orange-500/30 active:scale-95 transition-all">
							Play Again 🎮
						</button>
						<a href="/"
							class="w-full py-3 rounded-2xl bg-white/10 text-white font-medium text-sm border border-white/20 hover:bg-white/20 active:scale-95 transition-all flex items-center justify-center gap-2">
							🏠 Back to Home
						</a>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	:global(body) { background: #1e1b4b !important; }
	.play-root { -webkit-tap-highlight-color: transparent; }
	.glass { background: rgba(255,255,255,0.08); backdrop-filter: blur(20px); border: 1px solid rgba(255,255,255,0.12); }
	.fade-in { animation: fadeSlideIn 0.35s ease-out; }
	@keyframes fadeSlideIn { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: translateY(0); } }
	.score-popup {
		position: fixed; z-index: 50; pointer-events: none; font-size: 1.5rem; font-weight: 900;
		color: #fbbf24; text-shadow: 0 0 20px rgba(251,191,36,0.5); animation: scoreFloat 1.2s ease-out forwards;
	}
	@keyframes scoreFloat {
		0% { opacity: 1; transform: translateY(0) scale(1); }
		100% { opacity: 0; transform: translateY(-80px) scale(1.4); }
	}
	.scorecard {
		background: rgba(255,255,255,0.08);
		backdrop-filter: blur(20px);
		border: 2px solid rgba(255,255,255,0.15);
		box-shadow: 0 0 40px rgba(139,92,246,0.3), inset 0 1px 0 rgba(255,255,255,0.1);
	}
	.scorecard-border {
		border: 2px solid transparent;
		background: linear-gradient(135deg, rgba(251,191,36,0.2), rgba(139,92,246,0.2), rgba(236,72,153,0.2)) border-box;
		-webkit-mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0);
		-webkit-mask-composite: xor;
		mask-composite: exclude;
	}
	.score-glow {
		text-shadow: 0 0 30px rgba(251,191,36,0.4), 0 4px 8px rgba(0,0,0,0.3);
	}
</style>
