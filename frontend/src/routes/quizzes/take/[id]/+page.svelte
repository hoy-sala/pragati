<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { QuizAttempt } from '$lib/types';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';

	const quizId = $page.params.id;

	let attempt = $state<QuizAttempt | null>(null);
	let questions = $state<any[]>([]);
	let responses = $state<Record<string, { selected_options?: string[]; text_answer: string }>>({});
	let currentIndex = $state(0);
	let loading = $state(true);
	let submitting = $state(false);
	let saving = $state(false);
	let timeLeft = $state<number | null>(null);
	let autoSaveMsg = $state('');

	let timerInterval: ReturnType<typeof setInterval> | undefined;

	onMount(async () => {
		const quizRes = await api<any>('GET', `/quizzes/${quizId}`);
		const dur = quizRes.data?.duration_min;
		if (dur) timeLeft = dur * 60;

		const attRes = await api<any>('POST', `/quizzes/${quizId}/attempts`);
		if (attRes.data) {
			const aid = attRes.data.id;
			attempt = { id: aid, status: attRes.data.status } as QuizAttempt;

			const detailRes = await api<any>('GET', `/quizzes/attempts/${aid}`);
			if (detailRes.data) {
				if (detailRes.data.questions) {
					questions = detailRes.data.questions;
					// parse options for mcq questions
					for (const q of questions) {
						if (typeof q.options === 'string') q.options = JSON.parse(q.options);
					}
				}
				if (detailRes.data.saved) {
					for (const r of detailRes.data.saved) {
						responses[r.question_id] = {
							selected_options: r.selected_options || [],
							text_answer: r.text_answer || '',
						};
					}
				}
			}
		}
		loading = false;

		timerInterval = setInterval(() => {
			if (timeLeft != null && timeLeft > 0) {
				timeLeft--;
				if (timeLeft % 30 === 0) autoSave();
			}
		}, 1000);
	});

	onDestroy(() => {
		if (timerInterval) clearInterval(timerInterval);
	});

	let currentQuestion = $derived(questions[currentIndex]);
	let currentOptions = $derived<any[]>(currentQuestion?.options ? JSON.parse(currentQuestion.options) : []);

	function selectOption(qid: string, optKey: string) {
		if (!responses[qid]) responses[qid] = { text_answer: '' };
		const current = responses[qid].selected_options || [];
		const idx = current.indexOf(optKey);
		if (idx >= 0) {
			current.splice(idx, 1);
		} else {
			current.push(optKey);
		}
		responses[qid] = { ...responses[qid], selected_options: [...current] };
	}

	function setText(qid: string, val: string) {
		if (!responses[qid]) responses[qid] = { text_answer: '' };
		responses[qid] = { ...responses[qid], text_answer: val };
	}

	async function autoSave() {
		if (!attempt) return;
		saving = true;
		const data = Object.entries(responses).map(([question_id, r]) => ({
			question_id,
			selected_options: r.selected_options || [],
			text_answer: r.text_answer || '',
		}));
		await api('PUT', `/quizzes/attempts/${attempt.id}/answers`, { responses: data });
		saving = false;
		autoSaveMsg = 'Saved at ' + new Date().toLocaleTimeString();
		setTimeout(() => autoSaveMsg = '', 3000);
	}

	async function submitQuiz() {
		if (!confirm('Submit your quiz? This cannot be undone.')) return;
		if (!attempt) return;
		submitting = true;
		await autoSave();
		const res = await api('POST', `/quizzes/attempts/${attempt.id}/submit`);
		submitting = false;
		if (res.data) {
			goto(`/quizzes/results/${attempt.id}`);
		}
	}

	function formatTime(sec: number): string {
		const m = Math.floor(sec / 60);
		const s = sec % 60;
		return `${m}:${s.toString().padStart(2, '0')}`;
	}

	let questionStatus = $derived(
		questions.map((q, i) => {
			const r = responses[q.id];
			const answered = r && (r.text_answer || (r.selected_options && r.selected_options.length > 0));
			return { index: i, answered: !!answered };
		})
	);
</script>

<div class="max-w-4xl mx-auto space-y-4">
	<div class="flex items-center justify-between bg-white rounded-xl border border-slate-200 p-4">
		<div class="flex items-center gap-4">
			<h1 class="text-lg font-semibold text-slate-900">{(questions as any)[0]?.quiz_title || 'Quiz'}</h1>
			{#if timeLeft != null}
				<span class="px-3 py-1 rounded-lg text-sm font-mono {timeLeft < 120 ? 'bg-danger-50 text-danger-700' : 'bg-slate-100 text-slate-700'}">
					{formatTime(timeLeft)}
				</span>
			{/if}
			{#if saving}
				<span class="text-xs text-slate-400">Saving...</span>
			{/if}
			{#if autoSaveMsg}
				<span class="text-xs text-green-600">{autoSaveMsg}</span>
			{/if}
		</div>
		<button onclick={submitQuiz} disabled={submitting}
			class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
			{submitting ? 'Submitting...' : 'Submit'}
		</button>
	</div>

	<div class="flex items-center justify-center gap-1.5 flex-wrap">
		{#each questionStatus as qs}
			<button onclick={() => currentIndex = qs.index}
				class="w-8 h-8 text-xs font-medium rounded-lg border transition-colors
					{qs.answered ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-slate-600 border-slate-300'}
					{qs.index === currentIndex ? 'ring-2 ring-primary-300' : ''}">
				{qs.index + 1}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
	{:else if currentQuestion}
		<div class="bg-white rounded-xl border border-slate-200 p-6 space-y-6">
			<div class="flex items-start justify-between gap-4">
				<p class="text-sm text-slate-900 leading-relaxed">{currentQuestion.question_text}</p>
				<span class="shrink-0 text-xs text-slate-400">{currentQuestion.marks ?? 1} mark{(currentQuestion.marks ?? 1) !== 1 ? 's' : ''}</span>
			</div>

			{#if currentQuestion.question_type === 'mcq'}
				<div class="space-y-2">
					{#each currentOptions as opt}
						<div class="flex items-center gap-3 p-3 rounded-lg border {responses[currentQuestion.id]?.selected_options?.includes(opt.key) ? 'border-primary-400 bg-primary-50' : 'border-slate-200 hover:border-slate-300'} cursor-pointer transition-colors"
							onclick={() => selectOption(currentQuestion.id, opt.key)}>
							<input type="radio" name={currentQuestion.id}
								checked={responses[currentQuestion.id]?.selected_options?.includes(opt.key)}
								class="shrink-0"
								onchange={() => selectOption(currentQuestion.id, opt.key)}>
							<span class="text-sm text-slate-700">{opt.value}</span>
						</div>
					{/each}
				</div>
			{:else if currentQuestion.question_type === 'true_false'}
				<div class="flex gap-4">
					<button onclick={() => setText(currentQuestion.id, 'true')}
						class="flex-1 p-3 rounded-lg border text-sm font-medium transition-colors
							{responses[currentQuestion.id]?.text_answer === 'true' ? 'border-primary-400 bg-primary-50 text-primary-700' : 'border-slate-200 hover:border-slate-300'}">
						True
					</button>
					<button onclick={() => setText(currentQuestion.id, 'false')}
						class="flex-1 p-3 rounded-lg border text-sm font-medium transition-colors
							{responses[currentQuestion.id]?.text_answer === 'false' ? 'border-primary-400 bg-primary-50 text-primary-700' : 'border-slate-200 hover:border-slate-300'}">
						False
					</button>
				</div>
			{:else if currentQuestion.question_type === 'fill_blank' || currentQuestion.question_type === 'short_answer'}
				<textarea value={responses[currentQuestion.id]?.text_answer ?? ''}
					oninput={(e) => setText(currentQuestion.id, (e.target as HTMLTextAreaElement).value)}
					rows={3} placeholder="Type your answer..."
					class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm resize-none"></textarea>
			{/if}
		</div>

		<div class="flex justify-between">
			<button onclick={() => currentIndex = Math.max(0, currentIndex - 1)}
				disabled={currentIndex === 0}
				class="px-4 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50 disabled:opacity-30 transition-colors">
				Previous
			</button>
			<button onclick={autoSave} class="px-4 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50 transition-colors">
				Save
			</button>
			<button onclick={() => currentIndex = Math.min(questions.length - 1, currentIndex + 1)}
				disabled={currentIndex === questions.length - 1}
				class="px-4 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50 disabled:opacity-30 transition-colors">
				Next
			</button>
		</div>
	{/if}
</div>
