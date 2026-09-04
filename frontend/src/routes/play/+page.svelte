<script lang="ts">
	import { apiUrl } from '$lib/api/client.svelte';
	import { goto } from '$app/navigation';
	import { Building2, GraduationCap, BookOpen, Layers, Sparkles, Hash, Zap, Check, Globe, Newspaper, Clock, Atom, Users, Trophy, Cpu } from 'lucide-svelte';
	import type { PlayClass, PlaySubject, PlayTopic, PlayQuestion } from '$lib/types';
	import { ELEMENTS, CATEGORY_LABELS, elementClass, elementPeriod, PERIODIC_DIFFICULTIES } from '$lib/data/elements';
	import MathText from '$lib/components/MathText.svelte';

	type Phase = 'welcome' | 'mode' | 'classes' | 'subjects' | 'topics' | 'tables' | 'tables-ready' | 'gk' | 'coming-soon' | 'periodic' | 'periodic-ready' | 'quiz' | 'results' | 'team-create' | 'team-created' | 'team-list' | 'team-play';
	type GenQ = PlayQuestion & { uid: number; isRepeat?: boolean };
	type QuizEntry = 'subject' | 'gk' | 'ca';

	let phase = $state<Phase>('welcome');
	let playerName = $state('');	let classes = $state<PlayClass[]>([]);
	let subjects = $state<PlaySubject[]>([]);
	let topics = $state<PlayTopic[]>([]);
	let questions = $state<GenQ[]>([]);

	let selectedClass = $state<PlayClass | null>(null);
	let selectedSubject = $state<PlaySubject | null>(null);
	let selectedTopic = $state('');
	let selectedDifficulty = $state('');
	let loading = $state(false);
	let quizEntry = $state<QuizEntry>('subject');
	let comingSoonLabel = $state('');

	// ── Team quiz create (public) — Custom Quiz / Millionaire ──
	let teamTitle = $state('');
	let teamCount = $state(4);
	let perTeam = $state(5);
	let teamChapters = $state<string[]>([]);
	let allSubjectsForTeam = $state<PlaySubject[]>([]);
	let subjectTopicsMap = $state<Record<string, PlayTopic[]>>({});
	let teamCreating = $state(false);
	let teamError = $state('');
	let createdTeamQuiz = $state<any>(null);
	let customTeamNames = $state<string[]>(['A','B','C','D']);
	let timerChoice = $state(30);
	// ── Team play (round-robin, 10 pts, Millionaire) ──
	let teamPlayQuestions = $state<Record<string, any[]>>({});
	let teamPlayOrder = $state<{team:string, q:any}[]>([]);
	let teamPlayIndex = $state(0);
	let teamScores = $state<Record<string, number>>({});
	let teamPlayRevealed = $state(false);
	let teamPlayAnswered = $state(false);
	let teamPlaySelectedKey = $state('');
	let teamPlayIsCorrect = $state(false);
	let teamQuizzes = $state<any[]>([]);
	let teamQuizzesLoading = $state(false);
	let teamTimerSec = $state(30);
	let lifelineUsed = $state<Record<string, boolean>>({});
	let hiddenOptKeys = $state<string[]>([]);
	let teamFinished = $state(false);

	let currentIndex = $state(0);
	let score = $state(0);
	let streak = $state(0);
	let bestStreak = $state(0);
	let correctCount = $state(0);
	let selectedKey = $state('');
	let answered = $state(false);
	let isCorrect = $state(false);
	let timeLeft = $state(30);
	let timerInterval = $state<ReturnType<typeof setInterval> | undefined>();
	let startTime = $state(0);
	let showExitConfirm = $state(false);
	let screenShake = $state(false);
	let scorePopups = $state<{ id: number; value: number; x: number; y: number }[]>([]);
	let popupId = $state(0);
	let confettiPieces = $state<{ id: number; x: number; color: string; rot: number; delay: number }[]>([]);
	let answerBounce = $state(false);

	// ── Times tables ──
	let tablesActive = $state(false);
	let rushMode = $state(false);
	let selectedTable = $state(0); // 0 = mixed
	let rushLeft = $state(60);
	// ── Periodic table ──
	let periodicActive = $state(false);
	let selectedPeriodicMax = $state(0); // 30 | 60 | 118, 0 = none
	let qUid = $state(0);
	let uniqueTotal = $state(0);
	let uniqueWrong = $state<number[]>([]);
	let justMastered = $state(false);
	let mastered = $state<number[]>([]);
	let bestRush = $state<Record<string, number>>({});

	if (typeof window !== 'undefined') {
		try { mastered = JSON.parse(localStorage.getItem('pragati:tables:mastered') || '[]'); } catch { mastered = []; }
		try { bestRush = JSON.parse(localStorage.getItem('pragati:tables:best') || '{}'); } catch { bestRush = {}; }
		playerName = localStorage.getItem('pragati:player:name') || '';
		if (playerName) phase = 'mode';
	}
	function savePlayerName() {
		try { localStorage.setItem('pragati:player:name', playerName.trim()); } catch { /* ignore */ }
	}
	function saveTables() {
		try {
			localStorage.setItem('pragati:tables:mastered', JSON.stringify(mastered));
			localStorage.setItem('pragati:tables:best', JSON.stringify(bestRush));
		} catch { /* ignore */ }
	}

	function makeTableQ(a: number, b: number): GenQ {
		const ans = a * b;
		const cands = shuffle([ans + a, ans - a, ans + b, ans - b, ans + 10, ans - 10, a * (b + 1), b > 1 ? a * (b - 1) : -1, a + b]);
		const distractors: number[] = [];
		for (const c of cands) {
			if (distractors.length >= 3) break;
			if (c > 0 && c !== ans && !distractors.includes(c)) distractors.push(c);
		}
		let filler = ans + 1;
		while (distractors.length < 3) { if (filler !== ans && !distractors.includes(filler)) distractors.push(filler); filler++; }
		const options = shuffle([ans, ...distractors]).map((v, i) => ({ key: 'ABCD'[i], value: String(v), correct: v === ans }));
		return { question_text: `What is ${a} × ${b}?`, options, uid: ++qUid, isRepeat: false };
	}

	function randPair(): [number, number] {
		if (selectedTable === 0) return [2 + Math.floor(Math.random() * 9), 1 + Math.floor(Math.random() * 10)];
		return [selectedTable, 1 + Math.floor(Math.random() * 10)];
	}

	function genTablesQuestions(n = 10): GenQ[] {
		const out: GenQ[] = [];
		if (selectedTable === 0) {
			const used = new Set<string>();
			while (out.length < n) {
				const [a, b] = randPair();
				const k = `${a}x${b}`;
				if (!used.has(k)) { used.add(k); out.push(makeTableQ(a, b)); }
			}
		} else {
			for (const b of shuffle([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]).slice(0, n)) out.push(makeTableQ(selectedTable, b));
		}
		return out;
	}

	function startTablesQuiz(mode: 'practice' | 'rush') {
		playClick();
		tablesActive = true; periodicActive = false;
		rushMode = mode === 'rush';
		justMastered = false;
		selectedDifficulty = mode;
		uniqueTotal = 10; uniqueWrong = [];
		questions = genTablesQuestions(10);
		currentIndex = 0; score = 0; streak = 0; bestStreak = 0; correctCount = 0;
		selectedKey = ''; answered = false; startTime = Date.now();
		phase = 'quiz';
		if (rushMode) startRushTimer(); else startTimer();
	}

	let lastTickSecond = 99;
	function startRushTimer() {
		clearInterval(timerInterval);
		rushLeft = 60; lastTickSecond = 99;
		timerInterval = setInterval(() => {
			rushLeft -= 0.1;
			const s = Math.ceil(rushLeft);
			if (s !== lastTickSecond) { lastTickSecond = s; if (s <= 5 && s > 0) playTick(); }
			if (rushLeft <= 0) { rushLeft = 0; clearInterval(timerInterval); finishQuiz(); }
		}, 100);
	}

	function tableLabel() { return selectedTable === 0 ? 'Mixed tables' : `Table of ${selectedTable}`; }
	function periodicLabel() {
		if (!periodicActive) return '';
		const d = PERIODIC_DIFFICULTIES.find(x => x.id === selectedPeriodicMax);
		return d ? `Periodic · ${d.label} (${d.range})` : 'Periodic Table';
	}
	function displayAccuracy() {
		if ((tablesActive || periodicActive) && !rushMode) return Math.round(((uniqueTotal - uniqueWrong.length) / Math.max(uniqueTotal, 1)) * 100);
		return accuracy();
	}
	function displayCorrect() {
		if ((tablesActive || periodicActive) && !rushMode) return `${uniqueTotal - uniqueWrong.length}/${uniqueTotal}`;
		return `${correctCount}/${questions.length}`;
	}
	function rushColor() { return rushLeft > 20 ? '#0E7C71' : rushLeft > 10 ? '#B45309' : '#C2381B'; }

	// ── Periodic table ──
	function makeElementQ(max: number): GenQ {
		const pool = ELEMENTS.filter(e => e.z <= max);
		const el = pool[Math.floor(Math.random() * pool.length)];
		const tier: 'easy' | 'medium' | 'hard' = max <= 30 ? 'easy' : max <= 60 ? 'medium' : 'hard';
		type QKind = 'sym' | 'z' | 'nameBySym' | 'nameByZ' | 'class' | 'family' | 'period';
		let kinds: QKind[];
		if (tier === 'easy') kinds = ['sym', 'z', 'nameBySym'];
		else if (tier === 'medium') kinds = ['sym', 'z', 'nameBySym', 'nameByZ', 'class'];
		else kinds = ['sym', 'z', 'nameBySym', 'nameByZ', 'class', 'family', 'period'];
		const kind = kinds[Math.floor(Math.random() * kinds.length)];
		let qtext = '';
		let correct = '';
		let distractors: string[] = [];
		if (kind === 'sym') {
			qtext = `What is the symbol of ${el.name}?`;
			correct = el.sym;
			distractors = shuffle(pool.filter(x => x.sym !== correct).map(x => x.sym)).slice(0, 3);
		} else if (kind === 'z') {
			qtext = `What is the atomic number of ${el.name}?`;
			correct = String(el.z);
			const nums = [String(el.z + 1), String(Math.max(1, el.z - 1)), String(el.z + 2), String(el.z + 5), String(el.z + 10)];
			distractors = shuffle(nums.filter(n => n !== correct)).slice(0, 3);
		} else if (kind === 'nameBySym') {
			qtext = `Which element has the symbol ${el.sym}?`;
			correct = el.name;
			distractors = shuffle(pool.filter(x => x.name !== correct).map(x => x.name)).slice(0, 3);
		} else if (kind === 'nameByZ') {
			qtext = `Which element has atomic number ${el.z}?`;
			correct = el.name;
			distractors = shuffle(pool.filter(x => x.name !== correct).map(x => x.name)).slice(0, 3);
		} else if (kind === 'class') {
			qtext = `${el.name} is a…`;
			correct = elementClass(el.cat);
			distractors = shuffle(['Metal', 'Non-metal', 'Metalloid'].filter(c => c !== correct));
		} else if (kind === 'family') {
			qtext = `${el.name} belongs to which family?`;
			correct = CATEGORY_LABELS[el.cat];
			const cats = [...new Set(pool.filter(x => x.cat !== el.cat && x.cat !== 'unknown').map(x => CATEGORY_LABELS[x.cat]))];
			distractors = shuffle(cats).slice(0, 3);
		} else {
			qtext = `Which period is ${el.name} in?`;
			const p = elementPeriod(el.z);
			correct = String(p);
			distractors = shuffle(['1', '2', '3', '4', '5', '6', '7'].filter(x => x !== correct)).slice(0, 3);
		}
		const opts = shuffle([correct, ...distractors].slice(0, 4)).map((v, i) => ({ key: 'ABCD'[i], value: v, correct: v === correct }));
		return { question_text: qtext, options: opts, uid: ++qUid, isRepeat: false } as GenQ;
	}

	function genPeriodicQuestions(max: number, n = 10): GenQ[] {
		const out: GenQ[] = [];
		const used = new Set<string>();
		let attempts = 0;
		while (out.length < n && attempts < n * 10) {
			const q = makeElementQ(max);
			if (!used.has(q.question_text)) { used.add(q.question_text); out.push(q); }
			attempts++;
		}
		while (out.length < n) out.push(makeElementQ(max));
		return out;
	}

	function startPeriodicQuiz(mode: 'practice' | 'rush') {
		playClick();
		periodicActive = true; tablesActive = false; rushMode = mode === 'rush';
		justMastered = false; selectedDifficulty = mode;
		uniqueTotal = 10; uniqueWrong = [];
		questions = genPeriodicQuestions(selectedPeriodicMax, 10);
		currentIndex = 0; score = 0; streak = 0; bestStreak = 0; correctCount = 0;
		selectedKey = ''; answered = false; startTime = Date.now();
		phase = 'quiz';
		if (rushMode) startRushTimer(); else startTimer();
	}

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
	let masterGain: GainNode | null = null;
	let soundMuted = $state(false);
	if (typeof window !== 'undefined') {
		try { soundMuted = localStorage.getItem('pragati:sound:muted') === '1'; } catch { soundMuted = false; }
	}
	function getAudio() {
		if (!audioCtx && typeof window !== 'undefined') {
			audioCtx = new (window.AudioContext || (window as any).webkitAudioContext)();
			masterGain = audioCtx.createGain();
			masterGain.gain.value = soundMuted ? 0 : 0.9;
			masterGain.connect(audioCtx.destination);
		}
		if (masterGain) masterGain.gain.value = soundMuted ? 0 : 0.9;
		return audioCtx;
	}
	function toggleMute() {
		soundMuted = !soundMuted;
		try { localStorage.setItem('pragati:sound:muted', soundMuted ? '1' : '0'); } catch { /* ignore */ }
		if (masterGain && audioCtx) masterGain.gain.setValueAtTime(soundMuted ? 0 : 0.9, audioCtx.currentTime);
		if (!soundMuted) playClick();
	}

	function playTone(freq: number, dur: number, type: OscillatorType = 'sine', vol = 0.1, delayMs = 0, slideTo?: number) {
		if (soundMuted) return;
		const ctx = getAudio();
		if (!ctx || !masterGain) return;
		if (ctx.state === 'suspended') ctx.resume();
		const t0 = ctx.currentTime + delayMs / 1000;
		const osc = ctx.createOscillator();
		const gain = ctx.createGain();
		osc.connect(gain); gain.connect(masterGain);
		osc.type = type;
		osc.frequency.setValueAtTime(freq, t0);
		if (slideTo) osc.frequency.exponentialRampToValueAtTime(slideTo, t0 + dur);
		// Soft attack + smooth decay (no clicks / harsh edges)
		gain.gain.setValueAtTime(0.0001, t0);
		gain.gain.exponentialRampToValueAtTime(Math.max(vol, 0.001), t0 + 0.015);
		gain.gain.exponentialRampToValueAtTime(0.0001, t0 + dur);
		osc.start(t0); osc.stop(t0 + dur + 0.05);
	}

	// Polished game sounds — warm chimes, never harsh
	function playClick() { playTone(660, 0.07, 'triangle', 0.05); }
	function playCorrect() {
		// Warm rising chime: E5 → B5 with soft harmonic
		playTone(659, 0.14, 'triangle', 0.09);
		playTone(987, 0.22, 'triangle', 0.08, 90);
		playTone(1318, 0.18, 'sine', 0.04, 90);
	}
	function playWrong() {
		// Gentle descending "uh-oh": two soft sine steps, no buzz
		playTone(233, 0.16, 'sine', 0.08);
		playTone(174, 0.24, 'sine', 0.08, 130);
	}
	function playStreak() {
		// Sparkle run for streaks
		playTone(784, 0.09, 'triangle', 0.07);
		playTone(987, 0.09, 'triangle', 0.07, 70);
		playTone(1174, 0.2, 'triangle', 0.08, 140);
	}
	function playComplete() {
		// Little fanfare: C–E–G–C with warm overlap
		playTone(523, 0.18, 'triangle', 0.08);
		playTone(659, 0.18, 'triangle', 0.08, 130);
		playTone(784, 0.18, 'triangle', 0.08, 260);
		playTone(1046, 0.34, 'triangle', 0.09, 390);
		playTone(1318, 0.3, 'sine', 0.04, 390);
	}
	function playTick() { playTone(1180, 0.05, 'sine', 0.045); }
	function playReveal() {
		// Soft shimmer when answer is revealed
		playTone(520, 0.12, 'sine', 0.05, 0, 780);
		playTone(780, 0.14, 'sine', 0.04, 80);
	}
	function playTimeUp() {
		playTone(392, 0.16, 'triangle', 0.07);
		playTone(311, 0.26, 'triangle', 0.07, 140);
	}

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
		tablesActive = false; periodicActive = false;
		rushMode = false;
		quizEntry = 'subject';
		loading = true;
		classes = (await api<PlayClass[]>('GET', '/play/classes')) ?? [];
		loading = false;
		if (classes.length === 0) { comingSoonLabel = 'Class quizzes'; phase = 'coming-soon'; return; }
		phase = 'classes';
	}

	async function loadSubjects(cls: PlayClass) {
		playClick();
		selectedClass = cls;
		loading = true;
		const all = (await api<PlaySubject[]>('GET', `/play/subjects?class_id=${cls.id}`)) ?? [];
		subjects = all.filter(s => !/general\s*knowledge|current\s*affairs/i.test(s.name));
		loading = false;
		if (subjects.length === 0) { comingSoonLabel = `${cls.name} quizzes`; phase = 'coming-soon'; return; }
		phase = 'subjects';
	}

	async function loadTopics(sub: PlaySubject) {
		playClick();
		selectedSubject = sub;
		loading = true;
		const classParam = selectedClass ? `class_id=${selectedClass.id}&` : '';
		topics = (await api<PlayTopic[]>('GET', `/play/topics?${classParam}subject_id=${sub.id}`)) ?? [];
		loading = false;
		if (topics.length === 0) { comingSoonLabel = `${sub.name} quizzes`; phase = 'coming-soon'; return; }
		phase = 'topics';
	}

	async function loadGK() {
		playClick();
		tablesActive = false; periodicActive = false; rushMode = false;
		selectedClass = null;
		loading = true;
		const subs = (await api<PlaySubject[]>('GET', '/play/subjects')) ?? [];
		const gk = subs.find(s => /general\s*knowledge/i.test(s.name));
		if (!gk) { loading = false; comingSoonLabel = 'General Knowledge'; phase = 'coming-soon'; return; }
		selectedSubject = gk;
		topics = (await api<PlayTopic[]>('GET', `/play/topics?subject_id=${gk.id}`)) ?? [];
		loading = false;
		phase = 'gk';
	}

	async function loadCA() {
		playClick();
		tablesActive = false; periodicActive = false; rushMode = false;
		selectedClass = null;
		loading = true;
		const subs = (await api<PlaySubject[]>('GET', '/play/subjects')) ?? [];
		const ca = subs.find(s => /current\s*affairs/i.test(s.name));
		if (!ca) { loading = false; comingSoonLabel = 'Current Affairs'; phase = 'coming-soon'; return; }
		selectedSubject = ca;
		topics = (await api<PlayTopic[]>('GET', `/play/topics?subject_id=${ca.id}`)) ?? [];
		loading = false;
		if (topics.length === 0) { comingSoonLabel = 'Current Affairs'; phase = 'coming-soon'; return; }
		quizEntry = 'ca';
		selectedTopic = '';
		phase = 'topics';
	}

	async function loadComputerAwareness() {
		playClick();
		tablesActive = false; periodicActive = false; rushMode = false;
		selectedClass = null;
		loading = true;
		const subs = (await api<PlaySubject[]>('GET', '/play/subjects')) ?? [];
		const caw = subs.find(s => /computer\s*awareness/i.test(s.name));
		if (!caw) { loading = false; comingSoonLabel = 'Computer Awareness'; phase = 'coming-soon'; return; }
		selectedSubject = caw;
		topics = (await api<PlayTopic[]>('GET', `/play/topics?subject_id=${caw.id}`)) ?? [];
		loading = false;
		if (topics.length === 0) { comingSoonLabel = 'Computer Awareness'; phase = 'coming-soon'; return; }
		quizEntry = 'ca';
		selectedTopic = '';
		phase = 'topics';
	}

	async function openTeamCreate() {
		playClick();
		teamError = '';
		syncCustomNames();
		phase = 'team-create';
		loading = true;
		const subs = (await api<PlaySubject[]>('GET', '/play/subjects')) ?? [];
		// Show all except Current Affairs — including Kannada (may have 0 topics) and Social (32 topics)
		allSubjectsForTeam = subs.filter(s => !/current\s*affairs/i.test(s.name));
		for (const s of allSubjectsForTeam) {
			const t = (await api<PlayTopic[]>('GET', `/play/topics?subject_id=${s.id}`)) ?? [];
			subjectTopicsMap[s.id] = t;
		}
		loading = false;
	}
	async function openTeamList() {
		playClick();
		phase = 'team-list';
		teamQuizzesLoading = true;
		try {
			const res = await fetch(apiUrl('/team-quizzes/'));
			const json = await res.json();
			teamQuizzes = (json as any)?.data ?? json ?? [];
			if (!Array.isArray(teamQuizzes)) teamQuizzes = [];
		} catch { teamQuizzes = []; }
		teamQuizzesLoading = false;
	}
	function syncCustomNames() {
		const want = teamCount;
		const cur = [...customTeamNames];
		while (cur.length < want) cur.push(String.fromCharCode(65 + cur.length));
		customTeamNames = cur.slice(0, want);
	}
	function teamTimerLabel() { return teamTimerSec === 0 ? 'No timer' : `${teamTimerSec}s`; }
	async function openTeamQuiz(id: string) {
		playClick();
		loading = true;
		try {
			const res = await fetch(apiUrl(`/team-quizzes/${id}`));
			const json = await res.json();
			const data: any = (json as any)?.data ?? json;
			teamPlayQuestions = data.questions_by_team ?? {};
			teamPlayOrder = [];
			teamScores = {};
			lifelineUsed = {};
			hiddenOptKeys = [];
			teamTimerSec = data.timer_sec ?? 30;
			const teamNames = Object.keys(teamPlayQuestions).sort();
			for (const n of teamNames) { teamScores[n] = 0; lifelineUsed[n] = false; }
			const maxPer = Math.max(...teamNames.map(n => teamPlayQuestions[n].length));
			for (let r = 0; r < maxPer; r++) {
				for (const n of teamNames) {
					const q = teamPlayQuestions[n][r];
					if (q) teamPlayOrder.push({ team: n, q: { ...q, options: shuffle(q.options ?? []) } });
				}
			}
			teamPlayIndex = 0;
			teamPlayRevealed = false;
			teamPlayAnswered = false;
			teamPlaySelectedKey = '';
			teamFinished = false;
			createdTeamQuiz = data;
			phase = 'team-play';
			if (teamTimerSec > 0) startTeamTimer(); else { clearInterval(timerInterval); timeLeft = 0; }
		} catch { teamError = 'Failed to load'; }
		loading = false;
	}
	function toggleTeamChapter(ch: string) {
		playClick();
		if (teamChapters.includes(ch)) teamChapters = teamChapters.filter(c => c !== ch);
		else teamChapters = [...teamChapters, ch];
	}
	async function createTeamQuiz() {
		if (teamChapters.length === 0) { teamError = 'Pick at least one chapter.'; return; }
		syncCustomNames();
		const namesToSend = customTeamNames.slice(0, teamCount).map((n,i)=> (n||'').trim() || String.fromCharCode(65+i));
		const titleToSend = teamTitle.trim() || `Custom Quiz — ${teamChapters.slice(0,2).join(', ')}${teamChapters.length>2?'…':''}`;
		teamCreating = true; teamError = '';
		try {
			const res = await fetch(apiUrl('/team-quizzes/'), {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ title: titleToSend, description: '', teams: teamCount, per_team: perTeam, chapters: teamChapters, timer_sec: timerChoice, team_names: namesToSend })
			});
			const json = await res.json();
			if (!res.ok) { teamError = (json as any)?.message || (json as any)?.error || JSON.stringify(json) || 'Failed to create'; teamCreating = false; return; }
			teamCreating = false;
			const data: any = (json as any)?.data ?? json;
			createdTeamQuiz = data;
			// Immediately start round-robin play
			const qbt: Record<string, any[]> = data.questions_by_team ?? data.questions ?? {};
			teamPlayQuestions = qbt;
			teamPlayOrder = [];
			teamScores = {};
			lifelineUsed = {};
			hiddenOptKeys = [];
			teamTimerSec = data.timer_sec ?? timerChoice ?? 30;
			const teamNames = Object.keys(qbt).sort();
			for (const n of teamNames) { teamScores[n] = 0; lifelineUsed[n] = false; }
			const maxPer = Math.max(...teamNames.map(n => qbt[n].length));
			for (let r = 0; r < maxPer; r++) {
				for (const n of teamNames) {
					const q = qbt[n][r];
					if (q) teamPlayOrder.push({ team: n, q: { ...q, options: shuffle(q.options ?? []) } });
				}
			}
			teamPlayIndex = 0;
			teamPlayRevealed = false;
			teamPlayAnswered = false;
			teamPlaySelectedKey = '';
			teamFinished = false;
			phase = 'team-play';
			if (teamTimerSec > 0) { clearInterval(timerInterval); timeLeft = teamTimerSec; } else { clearInterval(timerInterval); timeLeft = 0; }
		} catch (e) { teamError = 'Network error'; teamCreating = false; }
	}
	function startTeamTimer() {
		if (teamTimerSec === 0) { clearInterval(timerInterval); timeLeft = 0; return; }
		playClick();
		clearInterval(timerInterval);
		timeLeft = teamTimerSec; lastTickSecond = 99;
		timerInterval = setInterval(() => {
			timeLeft -= 0.1;
			const s = Math.ceil(timeLeft);
			if (s !== lastTickSecond) { lastTickSecond = s; if (s <= 5 && s > 0) playTick(); }
			if (timeLeft <= 0) { timeLeft = 0; clearInterval(timerInterval); playTimeUp(); }
		}, 100);
	}
	function handleTeamAnswer(key: string) {
		if (teamPlayAnswered) return;
		teamPlayAnswered = true;
		teamPlaySelectedKey = key;
		const cur = teamPlayOrder[teamPlayIndex];
		const correct = cur.q.options.find((o: any) => o.correct)?.key;
		teamPlayIsCorrect = key === correct;
		if (teamPlayIsCorrect) {
			teamScores[cur.team] = (teamScores[cur.team] ?? 0) + 10;
			playCorrect();
			if ((teamScores[cur.team] ?? 0) % 30 === 0) playStreak();
			spawnConfetti();
			addPopup(10);
		} else {
			playWrong();
			screenShake = true;
			setTimeout(() => { screenShake = false; }, 500);
		}
		clearInterval(timerInterval);
	}
	function useFiftyFifty() {
		const cur = teamPlayOrder[teamPlayIndex];
		if (!cur || teamPlayAnswered || lifelineUsed[cur.team]) return;
		const wrongs = cur.q.options.filter((o: any) => !o.correct).map((o: any) => o.key);
		hiddenOptKeys = shuffle(wrongs).slice(0, 2);
		lifelineUsed[cur.team] = true;
		playClick();
	}
	function nextTeamQuestion() {
		if (teamPlayIndex >= teamPlayOrder.length - 1) { clearInterval(timerInterval); playComplete(); spawnConfetti(); teamFinished = true; phase = 'team-created'; return; }
		teamPlayIndex++;
		teamPlayRevealed = false;
		teamPlayAnswered = false;
		teamPlaySelectedKey = '';
		hiddenOptKeys = [];
		if (teamTimerSec > 0) { clearInterval(timerInterval); timeLeft = teamTimerSec; } else { clearInterval(timerInterval); timeLeft = 0; }
	}

	async function startQuiz(difficulty: string) {
		playClick();
		tablesActive = false; periodicActive = false; rushMode = false;
		selectedDifficulty = difficulty;
		loading = true;
		const classParam = selectedClass ? `class_id=${selectedClass.id}&` : '';
		const topicParam = selectedTopic ? `&topic=${encodeURIComponent(selectedTopic)}` : '';
		const data = await api<PlayQuestion[]>('GET', `/play/quiz?${classParam}subject_id=${selectedSubject!.id}${topicParam}&difficulty=${difficulty}&limit=10`);
		loading = false;
		if (!data || data.length === 0) {
			comingSoonLabel = selectedTopic ? `${selectedTopic} (${difficulty})` : `${selectedSubject?.name} (${difficulty})`;
			phase = 'coming-soon';
			return;
		}
		questions = data.map(q => ({ ...q, options: shuffle(q.options) }));
		currentIndex = 0; score = 0; streak = 0; bestStreak = 0; correctCount = 0;
		selectedKey = ''; answered = false; startTime = Date.now();
		phase = 'quiz';
		startTimer();
	}

	function startTimer() {
		clearInterval(timerInterval);
		timeLeft = 30; lastTickSecond = 99;
		timerInterval = setInterval(() => {
			timeLeft -= 0.1;
			const s = Math.ceil(timeLeft);
			if (s !== lastTickSecond) { lastTickSecond = s; if (s <= 5 && s > 0) playTick(); }
			if (timeLeft <= 0) { timeLeft = 0; clearInterval(timerInterval); playTimeUp(); handleAnswer(''); }
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
			const pts = rushMode ? 100 * Math.min(streak, 5) : (100 + Math.round(timeLeft * 10)) * Math.min(streak, 4);
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
		if (!rushMode) clearInterval(timerInterval);
		// Tables / Periodic practice: re-queue wrong questions a few positions later (active recall)
		if ((tablesActive || periodicActive) && !rushMode) {
			const q = questions[currentIndex];
			if (!q.isRepeat) {
				if (!isCorrect && !uniqueWrong.includes(q.uid)) uniqueWrong = [...uniqueWrong, q.uid];
			}
			if (!isCorrect && !q.isRepeat && questions.length < 15) {
				const clone: GenQ = { ...q, uid: ++qUid, isRepeat: true };
				const at = Math.min(currentIndex + 3, questions.length);
				questions = [...questions.slice(0, at), clone, ...questions.slice(at)];
			}
		}
		setTimeout(nextQuestion, rushMode ? 700 : 1200);
	}

	function nextQuestion() {
		if (phase !== 'quiz') return;
		if (rushMode) {
			currentIndex++; selectedKey = ''; answered = false;
			if (periodicActive) {
				questions = [...questions, makeElementQ(selectedPeriodicMax)];
			} else {
				const [a, b] = randPair();
				questions = [...questions, makeTableQ(a, b)];
			}
			return;
		}
		if (tablesActive || periodicActive) {
			if (currentIndex >= questions.length - 1) { finishQuiz(); return; }
			currentIndex++; selectedKey = ''; answered = false; startTimer();
			return;
		}
		if (currentIndex >= questions.length - 1) { finishQuiz(); return; }
		currentIndex++; selectedKey = ''; answered = false; startTimer();
	}

	function finishQuiz() {
		clearInterval(timerInterval);
		if (tablesActive && !rushMode && selectedTable !== 0) {
			const acc = (uniqueTotal - uniqueWrong.length) / Math.max(uniqueTotal, 1);
			if (acc >= 0.9 && !mastered.includes(selectedTable)) {
				mastered = [...mastered, selectedTable];
				justMastered = true;
			}
		}
		if (periodicActive && rushMode) {
			const key = `p${selectedPeriodicMax}`;
			if (score > (bestRush[key] || 0)) bestRush = { ...bestRush, [key]: score };
		} else if (tablesActive && rushMode) {
			const key = selectedTable === 0 ? 'mixed' : String(selectedTable);
			if (score > (bestRush[key] || 0)) bestRush = { ...bestRush, [key]: score };
		}
		if (tablesActive || periodicActive) saveTables();
		playComplete();
		spawnConfetti();
		phase = 'results';
	}

	function goBack() {
		playClick();
		if (phase === 'topics') {
			phase = quizEntry === 'subject' ? 'subjects' : quizEntry === 'gk' ? 'gk' : 'mode';
			return;
		}
		const back: Partial<Record<Phase, Phase>> = {
			mode: 'welcome', classes: 'mode', subjects: 'classes', tables: 'mode', 'tables-ready': 'tables',
			periodic: 'gk', 'periodic-ready': 'periodic',
			gk: 'mode', 'coming-soon': 'mode',
			'team-create': 'mode', 'team-created': 'mode', 'team-list': 'mode', 'team-play': 'team-create',
			quiz: 'welcome', results: 'welcome'
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
	function playAgain() { playClick(); phase = periodicActive ? 'periodic-ready' : tablesActive ? 'tables-ready' : 'topics'; }

	function timeColor() { return timeLeft > 20 ? '#0E7C71' : timeLeft > 10 ? '#B45309' : '#C2381B'; }
	function progressWidth() { return questions.length ? `${((currentIndex + 1) / questions.length) * 100}%` : '0%'; }
	function accuracy() { return questions.length ? Math.round((correctCount / questions.length) * 100) : 0; }
	function formatTime(ms: number) { const s = ms / 1000; return s < 60 ? `${s.toFixed(1)}s` : `${Math.floor(s / 60)}m ${Math.floor(s % 60)}s`; }
	function stars() { return displayAccuracy() >= 80 ? '★★★' : displayAccuracy() >= 50 ? '★★' : '★'; }
	function starsLabel() { return displayAccuracy() >= 80 ? 'Excellent' : displayAccuracy() >= 50 ? 'Good work' : 'Keep practising'; }
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
		<div style="display:flex;gap:0.5rem;align-items:center">
			<button onclick={toggleMute} class="qh-home" aria-label={soundMuted ? 'Unmute sounds' : 'Mute sounds'} title={soundMuted ? 'Unmute sounds' : 'Mute sounds'}>{soundMuted ? '🔇' : '🔊'}</button>
			{#if phase !== 'results'}
				<button onclick={confirmExit} class="qh-home">← Home</button>
			{/if}
		</div>
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
				<form class="welcome-form" onsubmit={(e) => { e.preventDefault(); if (playerName.trim()) { savePlayerName(); tablesActive = false; rushMode = false; phase = 'mode'; } }}>
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

	<!-- ═══ MODE ═══ -->
	{:else if phase === 'mode'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">What will you play?</h2>
					<p class="section-sub">{playerName}'s turn</p>
				</div>
			</div>
			<div class="mode-grid">
				<button onclick={loadClasses} class="pick-card mode-card">
					<span class="pick-icon"><BookOpen size={18} /></span>
					<span class="pick-name">Subject Quizzes</span>
					<span class="pick-meta">Maths, Science &amp; more</span>
				</button>
				<button onclick={loadGK} class="pick-card mode-card">
					<span class="pick-icon"><Globe size={18} /></span>
					<span class="pick-name">General Knowledge</span>
					<span class="pick-meta">1000 Freedom · 500 Karnataka · 1528 Qs</span>
				</button>
				<button onclick={loadCA} class="pick-card mode-card">
					<span class="pick-icon"><Newspaper size={18} /></span>
					<span class="pick-name">Current Affairs</span>
					<span class="pick-meta">What's happening around the world</span>
				</button>
				<button onclick={loadComputerAwareness} class="pick-card mode-card">
					<span class="pick-icon"><Cpu size={18} /></span>
					<span class="pick-name">Computer Awareness</span>
					<span class="pick-meta">Hardware · Software · Internet</span>
				</button>
				<button onclick={openTeamCreate} class="pick-card mode-card">
					<span class="pick-icon"><Trophy size={18} /></span>
					<span class="pick-name">Custom Quiz</span>
					<span class="pick-meta">Teams · 10 pts · Your timer</span>
				</button>
			</div>
		</div>

	<!-- ═══ CUSTOM QUIZ (team, immediate) ═══ -->
	{:else if phase === 'team-create'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">Custom Quiz</h2>
					<p class="section-sub">Teams · 10 pts · Your timer · Starts immediately</p>
				</div>
			</div>
			<div class="q-card" style="display:flex;flex-direction:column;gap:1rem">
				<div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:0.8rem">
					<label class="field">
						<span class="field-label">Number of teams</span>
						<select bind:value={teamCount} onchange={() => syncCustomNames()} class="input">
							{#each [2,3,4,5,6,7,8] as n}<option value={n}>{n} teams</option>{/each}
						</select>
					</label>
					<label class="field">
						<span class="field-label">Questions per team</span>
						<select bind:value={perTeam} class="input">
							{#each [5,10,15,20] as n}<option value={n}>{n} (multiples of 5)</option>{/each}
						</select>
						<span class="hint" style="text-align:left">Total {teamCount * perTeam} · Round-robin</span>
					</label>
					<label class="field">
						<span class="field-label">Timer</span>
						<select bind:value={timerChoice} class="input">
							<option value={0}>No timer</option>
							<option value={30}>30s</option>
							<option value={45}>45s</option>
							<option value={60}>60s</option>
						</select>
						<span class="hint" style="text-align:left">Expiry never reveals answer</span>
					</label>
				</div>
				<div style="border:2px dashed var(--ink);border-radius:12px;padding:0.85rem;background:var(--paper)">
					<p class="field-label" style="margin:0 0 0.6rem">Team names — editable</p>
					<div style="display:grid;grid-template-columns:1fr 1fr;gap:0.6rem">
						{#each Array.from({length: teamCount}, (_, i) => i) as i}
							<label class="field">
								<span class="field-label">Team {String.fromCharCode(65+i)}</span>
								<input bind:value={customTeamNames[i]} placeholder="Team {String.fromCharCode(65+i)}" maxlength={20} class="input" oninput={() => syncCustomNames()} />
							</label>
						{/each}
					</div>
				</div>
				{#if loading}
					<div class="empty">Loading chapters…</div>
				{:else}
					<div style="border:2px solid var(--ink);border-radius:12px;padding:0.85rem;background:var(--cream)">
						<p class="field-label" style="margin:0 0 0.6rem">Chapters — tick at least one (grouped by subject)</p>
						{#each allSubjectsForTeam as subj}
							{@const chapCount = (subjectTopicsMap[subj.id] ?? []).length}
							<div style="margin:0.7rem 0 0.35rem;font-weight:700;font-family:var(--font-display);font-size:0.95rem;opacity:{chapCount===0?0.55:1}">{subj.name} ({chapCount}) {#if chapCount===0}<span class="tag" style="margin-left:0.4rem;font-size:0.7rem">No chapters yet</span>{/if}</div>
							{#if chapCount > 0}
							<div class="topics" style="margin:0">
								{#each (subjectTopicsMap[subj.id] ?? []) as t}
									{@const full = t.name}
									<button onclick={() => toggleTeamChapter(full)} class="topic {teamChapters.includes(full) ? 'topic-selected' : ''}">{full.replace('Karnataka:','').replace('Freedom:','')}</button>
								{/each}
							</div>
							{:else}
							<p class="hint" style="text-align:left;margin:0 0 0.3rem">No chapters yet — add questions with [chapter:...] to enable.</p>
							{/if}
						{/each}
						<p class="hint" style="margin-top:0.7rem;text-align:left">{teamChapters.length} chapter(s) selected</p>
					</div>
				{/if}
				{#if teamError}<div class="feedback feedback-bad" style="margin:0"><p class="feedback-text">{teamError}</p></div>{/if}
				<button onclick={createTeamQuiz} disabled={teamCreating || teamChapters.length===0} class="btn-primary btn-block">{teamCreating ? 'Starting…' : 'Start Custom Quiz →'}</button>
				<p class="hint">Starts immediately — round-robin · No repeats · Difficulty balanced · 10 pts</p>
			</div>
		</div>

	<!-- ═══ TEAM CREATED / FINISHED ═══ -->
	{:else if phase === 'team-created'}
		{@const finishOrder = Object.keys(teamScores).sort((a,b)=>(teamScores[b]??0)-(teamScores[a]??0))}
		{@const topScore = finishOrder.length ? (teamScores[finishOrder[0]] ?? 0) : 0}
		{@const winners = finishOrder.filter(t => (teamScores[t] ?? 0) === topScore)}
		<div class="center fade-in">
			<div class="q-card" style="max-width:560px;width:100%;text-align:center">
				<div class="welcome-icon" aria-hidden="true"><Trophy size={28} /></div>
				{#if teamFinished && finishOrder.length}
					<h2 class="welcome-title" style="font-size:1.5rem">🏆 {winners.length > 1 ? `Tie: ${winners.join(' & ')}!` : `${winners[0]} Wins!`}</h2>
					<p class="welcome-sub" style="margin-top:0.5rem">Final scores · 10 pts per correct</p>
					<div style="display:flex;flex-direction:column;gap:0.5rem;margin:1rem 0;text-align:left">
						{#each finishOrder as t, i}
							<div class="side-row" style="background:{i===0?'var(--amber-tint)':'transparent'};border-radius:8px;padding:0.4rem 0.6rem">
								<span class="side-label">{i+1}. {t} {#if i===0}🏆{/if}</span>
								<span class="side-value mono">{teamScores[t] ?? 0} pts</span>
							</div>
						{/each}
					</div>
				{:else}
					<h2 class="welcome-title" style="font-size:1.5rem">Team Quiz Created!</h2>
					<p class="welcome-sub" style="margin-top:0.5rem">{createdTeamQuiz?.title ?? teamTitle} · {createdTeamQuiz?.teams ?? teamCount} teams · {createdTeamQuiz?.per_team ?? perTeam} Qs each · 10 pts per correct</p>
					<div style="background:var(--cream);border:2px solid var(--ink);border-radius:12px;padding:0.8rem;margin:1rem 0;text-align:left">
						<p class="field-label" style="margin:0 0 0.3rem">Chapters</p>
						<p style="font-size:0.9rem;line-height:1.5">{(createdTeamQuiz?.chapters ?? teamChapters).join(', ')}</p>
						<p class="hint" style="text-align:left;margin-top:0.5rem">Total {(createdTeamQuiz?.teams ?? teamCount) * (createdTeamQuiz?.per_team ?? perTeam)} questions · No repeats · Difficulty balanced · Round-robin</p>
					</div>
				{/if}
				<div style="display:flex;flex-direction:column;gap:0.6rem">
					<button onclick={() => { playClick(); phase = 'mode'; }} class="btn-primary btn-block">Back to Quizzes →</button>
					<button onclick={() => { playClick(); teamFinished = false; phase = 'team-create'; }} class="btn-ghost btn-block">Create Another</button>
				</div>
				<p class="hint" style="margin-top:0.8rem">Round-robin on projector · 10 pts · Expiry never reveals answer.</p>
			</div>
		</div>

	<!-- ═══ TEAM LIST ═══ -->
	{:else if phase === 'team-list'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">Team Quizzes</h2>
					<p class="section-sub">Public · Round-robin · 10 pts · Tap to play</p>
				</div>
			</div>
			{#if teamQuizzesLoading}
				<div class="empty">Loading…</div>
			{:else if teamQuizzes.length === 0}
				<div class="empty">No team quizzes yet.</div>
				<button onclick={openTeamCreate} class="btn-primary btn-block" style="margin-top:1rem">Create Team Quiz →</button>
			{:else}
				<div class="mode-grid">
					{#each teamQuizzes as tq}
						<button onclick={() => openTeamQuiz(tq.id)} class="pick-card mode-card">
							<span class="pick-icon"><Trophy size={18} /></span>
							<span class="pick-name">{tq.title}</span>
							<span class="pick-meta">{tq.teams} teams · {tq.per_team} each · {(tq.chapters ?? []).slice(0,2).join(', ')}{(tq.chapters?.length>2?'…':'')}</span>
						</button>
					{/each}
				</div>
			{/if}
		</div>

	<!-- ═══ TEAM PLAY (round-robin, 10 pts) ═══ -->
	{:else if phase === 'team-play' && teamPlayOrder.length > 0}
		{@const cur = teamPlayOrder[teamPlayIndex]}
		{@const teamLabels = Object.keys(teamScores).sort()}
		{@const curLifelineDone = !!lifelineUsed[cur.team]}
		<div class="quiz-shell fade-in">
			<div class="quiz-main">
				<div class="q-meta">
					<div class="q-tags">
						<span class="tag">Team {cur.team}</span>
						<span class="tag">{createdTeamQuiz?.title ?? 'Custom Quiz'}</span>
						<span class="tag tag-difficulty">{teamTimerSec === 0 ? 'No timer' : `${teamTimerSec}s`}</span>
						<span class="tag">Q {teamPlayIndex + 1}/{teamPlayOrder.length}</span>
					</div>
					<span class="counter">{cur.team}'s turn · 10 pts</span>
				</div>
				<div class="q-card q-question">
					<div style="display:flex;justify-content:space-between;align-items:center;gap:0.6rem;margin-bottom:0.8rem">
						<button onclick={useFiftyFifty} disabled={curLifelineDone || teamPlayAnswered} class="btn-ghost" style="min-height:36px;padding:0.3rem 0.8rem;border-radius:999px;opacity:{curLifelineDone ? 0.5 : 1}" title="50:50 — remove two wrong options (once per team)">50:50 {#if curLifelineDone}✓{:else}· {cur.team}{/if}</button>
						<span class="counter">{teamScores[cur.team] ?? 0} pts</span>
					</div>
					<fieldset class="q-fieldset">
						<legend class="q-legend"><MathText text={cur.q.question_text} /></legend>
						<div class="options">
							{#each cur.q.options as opt}
								{@const isSelected = teamPlayAnswered && opt.key === teamPlaySelectedKey}
								{@const isCorrectOpt = !!opt.correct}
								{@const showCorrect = (teamPlayAnswered || teamPlayRevealed) && !isSelected && isCorrectOpt}
								{@const hidden = hiddenOptKeys.includes(opt.key)}
								{#if !hidden}
								<label class="opt {teamPlayAnswered && isSelected ? (teamPlayIsCorrect ? 'opt-correct' : 'opt-incorrect') : ''} {showCorrect ? 'opt-correct' : ''} {!teamPlayAnswered && !teamPlayRevealed && teamPlaySelectedKey === opt.key ? 'opt-selected' : ''}">
									<input type="radio" name="tq{teamPlayIndex}" value={opt.key} disabled={teamPlayAnswered} checked={teamPlaySelectedKey === opt.key} onchange={() => handleTeamAnswer(opt.key)} class="opt-input" />
									<span class="opt-row">
										<span class="opt-radio"><span class="opt-dot"></span></span>
										<span class="opt-text"><MathText text={opt.value} /></span>
										{#if teamPlayAnswered && isSelected}<span class="opt-state">{teamPlayIsCorrect ? '✓' : '✕'}</span>{:else if showCorrect}<span class="opt-state">✓</span>{/if}
									</span>
								</label>
								{/if}
							{/each}
						</div>
					</fieldset>
					{#if teamPlayAnswered}
						<div class="feedback {teamPlayIsCorrect ? '' : 'feedback-bad'}">
							<div class="feedback-head">
								<span class="feedback-badge">{teamPlayIsCorrect ? '✓' : '✕'}</span>
								<span class="feedback-title">{teamPlayIsCorrect ? 'Correct! +10' : 'Not quite'}</span>
								<span class="streak-chip">{cur.team}</span>
							</div>
						</div>
					{:else if teamPlayRevealed}
						<div class="feedback">
							<div class="feedback-head">
								<span class="feedback-badge">✓</span>
								<span class="feedback-title">Answer revealed</span>
								<span class="streak-chip">{cur.team}</span>
							</div>
							<p class="feedback-text">Correct answer highlighted above. Tap a team’s answer to award 10 points.</p>
						</div>
					{/if}
					<div style="display:flex;gap:0.6rem;margin-top:1rem">
						{#if !teamPlayAnswered && !teamPlayRevealed}
							<button onclick={() => { teamPlayRevealed = true; clearInterval(timerInterval); playReveal(); }} class="btn-ghost" style="flex:1">Reveal Answer</button>
						{:else if teamPlayRevealed && !teamPlayAnswered}
							<button onclick={() => { teamPlayRevealed = false; if (teamTimerSec>0) startTeamTimer(); }} class="btn-ghost" style="flex:1">Hide Answer</button>
						{/if}
						{#if teamPlayAnswered || teamPlayRevealed}
							<button onclick={nextTeamQuestion} class="btn-primary" style="flex:1">{teamPlayIndex >= teamPlayOrder.length - 1 ? 'Finish →' : 'Next →'}</button>
						{/if}
					</div>
				</div>
			</div>
			<aside class="quiz-side">
				<div class="side-card">
					<div class="side-row"><span class="side-label">Current Team</span><span class="side-value mono">{cur.team}</span></div>
					{#if teamTimerSec > 0}
					<div class="timer-head">
						<span class="side-label">Time</span>
						<div style="display:flex;align-items:center;gap:0.5rem">
							<span class="timer-num mono" style="color:{timeColor()}">{Math.ceil(timeLeft)}s</span>
							<button onclick={startTeamTimer} class="btn-ghost" style="min-height:32px;padding:0.25rem 0.5rem;border-radius:8px" title="Start {teamTimerLabel()} timer after reading">
								<Clock size={16} />
							</button>
						</div>
					</div>
					<div class="timer-track"><div class="timer-fill" style="width:{(timeLeft/Math.max(teamTimerSec,1))*100}%;background:{timeColor()}"></div></div>
					<p class="hint" style="font-size:0.75rem;margin:0">Host: click <Clock size={12} style="display:inline;vertical-align:middle" /> after reading Q & options</p>
					{:else}
					<p class="hint" style="font-size:0.8rem;margin:0">No timer — host controls pace. Click an answer to award 10.</p>
					{/if}
				</div>
				<div class="side-card">
					<p class="side-label" style="margin-bottom:0.5rem">Scores (10 pts each)</p>
					{#each [...teamLabels].sort((a,b)=>(teamScores[b]??0)-(teamScores[a]??0)) as t, i}
						<div class="side-row"><span class="side-label">{i+1}. {t} {lifelineUsed[t] ? '· 50:50 ✓' : ''}</span><span class="side-value mono">{teamScores[t] ?? 0}</span></div>
					{/each}
				</div>
				<div class="side-card side-muted"><p class="side-hint">Host reads aloud. Timer expiry does not reveal answer. Correct = +10 with sound & confetti.</p></div>
			</aside>
		</div>

	<!-- ═══ GK HUB ═══ -->
	{:else if phase === 'gk'}
		{@const ktopics = topics.filter(t => t.name.startsWith('Karnataka:'))}
		{@const ftopics = topics.filter(t => t.name.startsWith('Freedom:'))}
		{@const otherTopics = topics.filter(t => !t.name.startsWith('Karnataka:') && !t.name.startsWith('Freedom:'))}
		{@const kLabels: Record<string,string> = {'Karnataka:Districts':'Districts & Divisions','Karnataka:Geography':'Geography & Nature','Karnataka:History':'History & Kingdoms','Karnataka:Culture':'Culture & Heritage','Karnataka:Language':'Language & Literature','Karnataka:Wildlife':'Wildlife & Reserves','Karnataka:Symbols':'Government & Symbols','Karnataka:Legends':'Legends & Today'}}
		{@const fLabels: Record<string,string> = {'Freedom:1857':'1857 & Early Resistance','Freedom:1885-1905':'INC & Swadeshi (1885-1905)','Freedom:1905-1919':'Revolutionary & Home Rule','Freedom:Gandhian-I':'Gandhian I: NCM','Freedom:Gandhian-II':'Gandhian II: CDM & 1935','Freedom:Revolutionary-INA':'Revolutionaries & INA','Freedom:Leadership':'Leaders & Personalities','Freedom:Endgame':'Quit India & Partition'}}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">General Knowledge</h2>
					<p class="section-sub">Pick a challenge, {playerName} · {topics.length} topics</p>
				</div>
			</div>
			<div class="mode-grid">
				<button onclick={() => { playClick(); tablesActive = true; periodicActive = false; rushMode = false; phase = 'tables'; }} class="pick-card mode-card">
					<span class="pick-icon"><Hash size={18} /></span>
					<span class="pick-name">Times Tables</span>
					<span class="pick-meta">Learn · Practice · Rush</span>
				</button>
				<button onclick={() => { playClick(); periodicActive = false; tablesActive = false; selectedPeriodicMax = 0; phase = 'periodic'; }} class="pick-card mode-card">
					<span class="pick-icon"><Atom size={18} /></span>
					<span class="pick-name">Periodic Table</span>
					<span class="pick-meta">Elements · Easy to Hard</span>
				</button>
			</div>
			{#if loading}
				<div class="empty">Loading…</div>
			{:else}
				{#if ktopics.length}
					<div style="margin-top:1.1rem">
						<p class="tchip-label" style="margin:0 0 0.5rem">Karnataka — Namma Nadu · 500 questions</p>
						<div class="mode-grid">
							{#each ktopics as topic (topic.name)}
								<button onclick={() => { playClick(); selectedTopic = topic.name; quizEntry = 'gk'; phase = 'topics'; }} class="pick-card mode-card">
									<span class="pick-icon"><Globe size={18} /></span>
									<span class="pick-name">{kLabels[topic.name] ?? topic.name.replace('Karnataka:','')}</span>
									<span class="pick-meta">Karnataka · Quiz · Easy to Hard</span>
								</button>
							{/each}
						</div>
					</div>
				{/if}
				{#if ftopics.length}
					<div style="margin-top:1.1rem">
						<p class="tchip-label" style="margin:0 0 0.5rem">Indian Freedom Movement · 1000 questions</p>
						<div class="mode-grid">
							{#each ftopics as topic (topic.name)}
								<button onclick={() => { playClick(); selectedTopic = topic.name; quizEntry = 'gk'; phase = 'topics'; }} class="pick-card mode-card">
									<span class="pick-icon"><Globe size={18} /></span>
									<span class="pick-name">{fLabels[topic.name] ?? topic.name.replace('Freedom:','')}</span>
									<span class="pick-meta">Freedom · Quiz · Easy to Hard</span>
								</button>
							{/each}
						</div>
					</div>
				{/if}
				{#if otherTopics.length}
					<div style="margin-top:1.1rem">
						<p class="tchip-label" style="margin:0 0 0.5rem">India & More</p>
						<div class="mode-grid">
							{#each otherTopics as topic (topic.name)}
								<button onclick={() => { playClick(); selectedTopic = topic.name; quizEntry = 'gk'; phase = 'topics'; }} class="pick-card mode-card">
									<span class="pick-icon"><Globe size={18} /></span>
									<span class="pick-name">{topic.name}</span>
									<span class="pick-meta">Quiz · Easy to Hard</span>
								</button>
							{/each}
						</div>
					</div>
				{/if}
			{/if}
		</div>

	<!-- ═══ COMING SOON ═══ -->
	{:else if phase === 'coming-soon'}
		<div class="center fade-in">
			<div class="q-card welcome-card">
				<div class="welcome-icon" aria-hidden="true"><Clock size={28} /></div>
				<h1 class="welcome-title">Coming soon!</h1>
				<p class="welcome-sub">{comingSoonLabel} are being prepared right now. Check back soon — new questions are on the way!</p>
				<div class="welcome-form">
					<button onclick={() => { playClick(); phase = 'mode'; }} class="btn-primary btn-block">← Choose another quiz</button>
				</div>
			</div>
		</div>

	<!-- ═══ TABLES PICK ═══ -->
	{:else if phase === 'tables'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">Pick a table</h2>
					<p class="section-sub">Master 2 to 10 first, then take the challenge</p>
				</div>
			</div>
			<div class="tchip-grid">
				{#each [2, 3, 4, 5, 6, 7, 8, 9, 10] as n (n)}
					<button onclick={() => { playClick(); selectedTable = n; phase = 'tables-ready'; }} class="tchip {mastered.includes(n) ? 'tchip-done' : ''}">
						{n}
						{#if mastered.includes(n)}<span class="tchip-check"><Check size={12} strokeWidth={3} /></span>{/if}
					</button>
				{/each}
			</div>
			<p class="tchip-label">Challenge · 11 to 20</p>
			<div class="tchip-grid">
				{#each [11, 12, 13, 14, 15, 16, 17, 18, 19, 20] as n (n)}
					<button onclick={() => { playClick(); selectedTable = n; phase = 'tables-ready'; }} class="tchip {mastered.includes(n) ? 'tchip-done' : ''}">
						{n}
						{#if mastered.includes(n)}<span class="tchip-check"><Check size={12} strokeWidth={3} /></span>{/if}
					</button>
				{/each}
			</div>
			<button onclick={() => { playClick(); selectedTable = 0; phase = 'tables-ready'; }} class="btn-ghost btn-block" style="margin-top:1rem">Mixed — random from 2 to 10</button>
			<p class="hint">Get 9 or more right in Practice to master a table ✓</p>
		</div>

	<!-- ═══ PERIODIC PICK ═══ -->
	{:else if phase === 'periodic'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">Periodic Table</h2>
					<p class="section-sub">30 easy · 60 medium · 118 hard</p>
				</div>
			</div>
			<div class="diff-stack">
				{#each PERIODIC_DIFFICULTIES as d (d.id)}
					<button onclick={() => { playClick(); selectedPeriodicMax = d.id; phase = 'periodic-ready'; }} class="diff-btn {d.id === 30 ? 'diff-easy' : d.id === 60 ? 'diff-medium' : 'diff-hard'}">
						<span class="diff-left"><Atom size={16} /> {d.label} — {d.range}</span>
						<span class="diff-meta">{d.desc}</span>
					</button>
				{/each}
			</div>
		</div>

	<!-- ═══ PERIODIC READY / LEARN ═══ -->
	{:else if phase === 'periodic-ready'}
		<div class="center fade-in">
			<div class="q-card" style="max-width:560px;width:100%">
				<div class="section-head" style="margin:0 0 0.9rem">
					<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
					<div>
						<h2 class="section-title">Periodic Table — {PERIODIC_DIFFICULTIES.find(x => x.id === selectedPeriodicMax)?.label ?? ''}</h2>
						<p class="section-sub">{PERIODIC_DIFFICULTIES.find(x => x.id === selectedPeriodicMax)?.range ?? ''} · {selectedPeriodicMax === 30 ? 'Symbols & numbers' : selectedPeriodicMax === 60 ? 'Plus metal vs non-metal' : 'All families & periods'}</p>
					</div>
				</div>
				<div class="p-info">
					<div class="p-info-row"><span class="p-info-k">Pool</span><span class="p-info-v">{PERIODIC_DIFFICULTIES.find(x => x.id === selectedPeriodicMax)?.range ?? ''}</span></div>
					<div class="p-info-row"><span class="p-info-k">Focus</span><span class="p-info-v">{selectedPeriodicMax === 30 ? 'Symbols & atomic numbers' : selectedPeriodicMax === 60 ? 'Plus Metal / Non-metal / Metalloid' : 'Families & periods — high-school level'}</span></div>
					<div class="p-info-row"><span class="p-info-k">Example</span><span class="p-info-v p-info-ex">{selectedPeriodicMax === 30 ? 'What is the symbol of Oxygen? → O' : selectedPeriodicMax === 60 ? 'Sodium is a… → Metal' : 'Which family is Iodine in? → Halogen'}</span></div>
				</div>
				<div class="welcome-form">
					<button onclick={() => startPeriodicQuiz('practice')} class="btn-primary btn-block">Practice — 10 questions →</button>
					<button onclick={() => startPeriodicQuiz('rush')} class="btn-ghost btn-block"><Zap size={16} /> Rush — 60 seconds</button>
					<p class="hint">Practice shuffles the pool. Rush adds new questions non-stop for 60 seconds.</p>
				</div>
			</div>
		</div>

	<!-- ═══ TABLES READY / LEARN ═══ -->
	{:else if phase === 'tables-ready'}
		<div class="center fade-in">
			<div class="q-card" style="max-width:520px;width:100%">
				<div class="section-head" style="margin:0 0 0.9rem">
					<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
					<div>
						<h2 class="section-title">{tableLabel()}</h2>
						<p class="section-sub">{selectedTable === 0 ? 'Any question from 2×1 to 10×10' : `See the pattern, then race it`}</p>
					</div>
				</div>
				{#if selectedTable !== 0}
					<div class="tt-grid">
						{#each [1, 2, 3, 4, 5, 6, 7, 8, 9, 10] as i (i)}
							<div class="tt-row"><span class="tt-expr">{selectedTable} × {i}</span><span class="tt-eq">=</span><span class="tt-val">{selectedTable * i}</span></div>
						{/each}
					</div>
				{/if}
				<div class="welcome-form">
					<button onclick={() => startTablesQuiz('practice')} class="btn-primary btn-block">Practice — 10 questions →</button>
					<button onclick={() => startTablesQuiz('rush')} class="btn-ghost btn-block"><Zap size={16} /> Rush — 60 seconds</button>
					<p class="hint">Learn the pattern first — then beat the clock!</p>
				</div>
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

	<!-- ═══ TOPICS + DIFFICULTY (merged start screen) ═══ -->
	{:else if phase === 'topics'}
		<div class="stack fade-in">
			<div class="section-head">
				<button onclick={goBack} class="back-btn" aria-label="Back">←</button>
				<div>
					<h2 class="section-title">{selectedTopic || 'Select topic'}</h2>
					<p class="section-sub">{selectedClass ? `${selectedClass.name} · ` : ''}{selectedSubject?.name}</p>
				</div>
			</div>
			{#if loading}
				<div class="empty">Loading…</div>
			{:else}
				{@const topicFamily = selectedTopic.includes(':') ? selectedTopic.split(':')[0] : ''}
				{@const visibleTopics = topicFamily && (topicFamily === 'Karnataka' || topicFamily === 'Freedom') ? topics.filter(t => t.name.startsWith(topicFamily + ':')) : topics}
				{#if !(quizEntry === 'gk' && selectedTopic)}
					<div class="topics">
						<button onclick={() => { playClick(); selectedTopic = ''; }} class="topic {selectedTopic === '' ? 'topic-selected' : ''}">All topics</button>
						{#each visibleTopics as topic (topic.name)}
							<button onclick={() => { playClick(); selectedTopic = topic.name; }} class="topic {selectedTopic === topic.name ? 'topic-selected' : ''}">{topic.name.replace('Karnataka:','').replace('Freedom:','')}</button>
						{/each}
					</div>
				{:else}
					<p class="section-sub" style="margin:0 0 0.6rem; color:var(--ink-soft)">{selectedTopic.replace('Karnataka:','').replace('Freedom:','')} · choose difficulty to start</p>
				{/if}
				<div class="diff-stack" style="margin-top:1.25rem">
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

	<!-- ═══ QUIZ ═══ -->
	{:else if phase === 'quiz' && questions.length > 0}
		{@const q = questions[currentIndex]}
		<div class="quiz-shell fade-in">
			<!-- main card -->
			<div class="quiz-main">
				<div class="q-meta">
					<div class="q-tags">
						{#if tablesActive}
							<span class="tag">Times tables</span>
							<span class="tag">{tableLabel()}</span>
							<span class="tag tag-difficulty">{rushMode ? 'Rush · 60s' : 'Practice'}</span>
						{:else if periodicActive}
							<span class="tag">Periodic Table</span>
							<span class="tag">{periodicLabel()}</span>
							<span class="tag tag-difficulty">{rushMode ? 'Rush · 60s' : 'Practice'}</span>
						{:else}
							{#if selectedClass}<span class="tag">{selectedClass.name}</span>{/if}
							<span class="tag">{selectedSubject?.name}</span>
							{#if selectedTopic}<span class="tag">{selectedTopic}</span>{/if}
							<span class="tag tag-difficulty tag-{selectedDifficulty}">{selectedDifficulty}</span>
						{/if}
					</div>
					<span class="counter">{rushMode ? `${questions.length} asked` : `${currentIndex + 1} / ${questions.length}`}</span>
				</div>

				<div class="q-card q-question {answerBounce ? 'answer-pop' : ''}">
					<fieldset class="q-fieldset">
						<legend class="q-legend"><MathText text={q.question_text} /></legend>
						<div class="options">
							{#each q.options as opt, oi (`${currentIndex}-${opt.key}`)}
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
										<span class="opt-text"><MathText text={opt.value} /></span>
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
							<span class="stat-v mono">{displayAccuracy()}%</span>
							<span class="stat-k">Accuracy</span>
						</div>
						<div class="stat">
							<span class="stat-v mono">{displayCorrect()}</span>
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
						<span class="side-label">{rushMode ? 'Time left' : 'Time'}</span>
						<span class="timer-num mono" style="color:{rushMode ? rushColor() : timeColor()}">{Math.ceil(rushMode ? rushLeft : timeLeft)}s</span>
					</div>
					<div class="timer-track">
						<div class="timer-fill" style="width:{(rushMode ? rushLeft / 60 : timeLeft / 30) * 100}%; background:{rushMode ? rushColor() : timeColor()}"></div>
					</div>
					{#if streak >= 2}
						<div class="streak-box">{streak}× streak</div>
					{/if}
				</div>

				<div class="side-card side-muted">
					<p class="side-hint">{rushMode ? 'Speed round! Every correct answer counts — go, go, go!' : 'Tip: quick answers earn a bit extra. Get 3 in a row for a streak!'}</p>
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
				<p class="score-meta">
					{#if tablesActive}{tableLabel()} · {rushMode ? 'Rush' : 'Practice'}{:else if periodicActive}{periodicLabel()} · {rushMode ? 'Rush' : 'Practice'}{:else}{#if selectedClass}{selectedClass.name} · {/if}{selectedSubject?.name}{selectedTopic ? ` · ${selectedTopic}` : ''} · {selectedDifficulty}{/if}
				</p>
				{#if justMastered}
					<span class="master-badge"><Check size={14} strokeWidth={3} /> Table mastered!</span>
				{/if}
				<div class="score-stats">
					<div class="s-stat"><span class="s-v mono">{displayAccuracy()}%</span><span class="s-k">Accuracy</span></div>
					<div class="s-sep"></div>
					<div class="s-stat"><span class="s-v mono">{bestStreak}</span><span class="s-k">Best streak</span></div>
					<div class="s-sep"></div>
					<div class="s-stat"><span class="s-v mono">{displayCorrect()}</span><span class="s-k">Correct</span></div>
				</div>
				{#if (tablesActive || periodicActive) && rushMode}
					<p class="score-time">Personal best: {Math.max(score, bestRush[(periodicActive ? `p${selectedPeriodicMax}` : selectedTable === 0 ? 'mixed' : String(selectedTable))] || 0).toLocaleString()}</p>
				{:else}
					<p class="score-time">{formatTime(Date.now() - startTime)} · {(tablesActive || periodicActive) ? `${uniqueTotal} questions` : `${questions.length} questions`}</p>
				{/if}
				{#if (tablesActive || periodicActive) && !rushMode && uniqueWrong.length > 0}
					<p class="score-time">Practise these: {uniqueWrong.slice(0, 4).map(uid => { const q = questions.find(qq => qq.uid === uid); return q ? q.question_text.replace('What is ', '').replace('?', '') : ''; }).filter(Boolean).join(' · ')}</p>
				{/if}
				<div class="score-actions">
					<button onclick={playAgain} class="btn-primary btn-block">Play again →</button>
					<button onclick={() => { playClick(); phase = periodicActive ? 'periodic' : tablesActive ? 'tables' : 'mode'; }} class="btn-ghost btn-block">Choose {periodicActive ? 'level' : tablesActive ? 'table' : 'quiz'}</button>
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
	.topic-selected { background: var(--ink); color: var(--paper); border-width: 2px; }
	.topic-selected:hover { background: var(--ink); }
	.topic-all { background: var(--amber); border-width: 2px; }

	/* mode select */
	.mode-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0.8rem; }
	@media (max-width: 560px) { .mode-grid { grid-template-columns: 1fr; } }
	.mode-card { flex-direction: column; align-items: flex-start; gap: 0.4rem; padding: 1.15rem; }
	.mode-card .pick-name { font-size: 1.1rem; }

	/* tables pick chips */
	.tchip-grid { display: flex; flex-wrap: wrap; justify-content: center; gap: 0.55rem; }
	.tchip {
		position: relative; width: 56px; min-height: 52px; border-radius: 12px;
		border: 2px solid var(--ink); background: var(--paper); color: var(--ink);
		font-family: var(--font-display); font-weight: 800; font-size: 1.15rem;
		cursor: pointer; display: grid; place-items: center;
		box-shadow: 3px 3px 0 var(--ink); transition: box-shadow 120ms, transform 120ms, background 120ms;
	}
	.tchip:hover { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--ink); background: var(--cream); }
	.tchip:active { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
	.tchip-done { background: var(--mint); }
	.tchip-check {
		position: absolute; top: -7px; right: -7px; width: 18px; height: 18px;
		border-radius: 50%; border: 2px solid var(--ink); background: var(--teal); color: var(--paper);
		display: grid; place-items: center;
	}
	.tchip-label {
		font-family: var(--font-mono); font-size: 0.68rem; font-weight: 700;
		letter-spacing: 0.08em; text-transform: uppercase; color: var(--ink-soft);
		margin: 0.9rem 0 0.45rem;
	}

	/* tables learn grid */
	.tt-grid {
		display: grid; grid-template-columns: 1fr 1fr; gap: 0.4rem 1.2rem;
		background: var(--cream); border: 2px solid var(--ink); border-radius: 14px;
		padding: 0.9rem 1.1rem;
	}
	.tt-row { display: flex; align-items: baseline; gap: 0.5rem; font-size: 1.02rem; padding: 0.15rem 0; }
	.tt-expr { font-weight: 600; color: var(--ink-soft); min-width: 4.2rem; }
	.tt-eq { color: var(--ink-soft); }
	.tt-val { font-family: var(--font-display); font-weight: 800; font-size: 1.15rem; color: var(--ink); }

	/* periodic preview */
	.p-info { background: var(--cream); border: 2px solid var(--ink); border-radius: 12px; padding: 0.9rem 1rem; display: flex; flex-direction: column; gap: 0.5rem; }
	.p-info-row { display: flex; justify-content: space-between; align-items: baseline; gap: 1rem; font-size: 0.92rem; }
	.p-info-k { font-size: 0.68rem; font-weight: 700; letter-spacing: 0.08em; text-transform: uppercase; color: var(--ink-soft); }
	.p-info-v { font-weight: 600; color: var(--ink); text-align: right; }
	.p-info-ex { font-family: var(--font-mono); font-size: 0.85rem; color: var(--plum); font-weight: 600; }

	/* mastery badge */
	.master-badge {
		display: inline-flex; align-items: center; gap: 0.35rem;
		font-size: 0.85rem; font-weight: 700;
		border: 2px solid var(--ink); background: var(--amber); border-radius: 999px;
		padding: 0.28rem 0.75rem; color: var(--ink);
	}

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
	.opt-text { flex: 1; min-width: 0; font-size: 1.02rem; line-height: 1.5; font-weight: 700; color: var(--ink); overflow-wrap: anywhere; word-break: break-word; letter-spacing: -0.01em; }
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
