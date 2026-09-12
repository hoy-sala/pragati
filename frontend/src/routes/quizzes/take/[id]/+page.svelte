<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { QuizAttempt } from '$lib/types';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount, onDestroy } from 'svelte';
	import { X, ChevronLeft, ChevronRight, Check } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import Modal from '$lib/components/Modal.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import MathText from '$lib/components/MathText.svelte';

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
	let showSubmit = $state(false);

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
	let currentOptions = $derived<any[]>(typeof currentQuestion?.options === 'string' ? JSON.parse(currentQuestion.options) : (currentQuestion?.options || []));

	function selectOption(qid: string, optKey: string) {
		if (!responses[qid]) responses[qid] = { text_answer: '' };
		const current = responses[qid].selected_options || [];
		const selected = current.includes(optKey) ? [] : [optKey];
		responses[qid] = { ...responses[qid], selected_options: selected };
	}

	function setText(qid: string, val: string) {
		if (!responses[qid]) responses[qid] = { text_answer: '' };
		responses[qid] = { ...responses[qid], text_answer: val };
	}

	async function autoSave(quiet = false) {
		if (!attempt) return;
		saving = true;
		const data = Object.entries(responses).map(([question_id, r]) => ({
			question_id,
			selected_options: r.selected_options || [],
			text_answer: r.text_answer || '',
		}));
		await api('PUT', `/quizzes/attempts/${attempt.id}/answers`, { responses: data });
		saving = false;
		if (!quiet) {
			autoSaveMsg = 'Saved at ' + new Date().toLocaleTimeString();
			setTimeout(() => autoSaveMsg = '', 3000);
		}
	}

	function gotoQuestion(i: number) {
		const next = Math.max(0, Math.min(questions.length - 1, i));
		if (next !== currentIndex) {
			autoSave(true);
			currentIndex = next;
		}
	}

	async function submitQuiz() {
		if (!attempt) return;
		submitting = true;
		await autoSave(true);
		const res = await api('POST', `/quizzes/attempts/${attempt.id}/submit`);
		submitting = false;
		showSubmit = false;
		if (res.data) {
			goto(`/quizzes/results/${attempt.id}`);
		}
	}

	function exitQuiz() {
		autoSave(true);
		if (window.history.length > 1) window.history.back();
		else goto('/quizzes/available');
	}

	function formatTime(sec: number): string {
		const m = Math.floor(sec / 60);
		const s = sec % 60;
		return `${m}:${s.toString().padStart(2, '0')}`;
	}

	function isAnswered(q: any): boolean {
		const r = responses[q.question_id];
		return !!(r && (r.text_answer || (r.selected_options && r.selected_options.length > 0)));
	}

	let questionStatus = $derived(
		questions.map((q, i) => ({ index: i, answered: isAnswered(q) }))
	);
	let answeredCount = $derived(questionStatus.filter(q => q.answered).length);
	let progressPct = $derived(questions.length ? Math.round(((currentIndex + 1) / questions.length) * 100) : 0);
	let isLast = $derived(questions.length > 0 && currentIndex === questions.length - 1);
</script>

<div class="max-w-2xl mx-auto px-4 pb-4 md:px-0 min-h-full flex flex-col">
	<!-- App header -->
	<div class="sticky top-0 z-10 -mx-4 px-4 pt-3 pb-2 bg-slate-50/95 backdrop-blur md:mx-0 md:px-0">
		<div class="flex items-center gap-2">
			<button
				onclick={exitQuiz}
				aria-label="Exit quiz"
				class="w-10 h-10 rounded-full bg-white border border-slate-200 flex items-center justify-center text-slate-600 active:scale-95 transition-transform shrink-0"
			>
				<X size={20} />
			</button>
			<div class="flex-1 min-w-0 text-center">
				<div class="text-sm font-bold text-slate-900">Question {questions.length ? currentIndex + 1 : '–'} of {questions.length}</div>
				<div class="text-xs text-slate-500 truncate">{(questions as any)[0]?.quiz_title || 'Quiz'}</div>
			</div>
			{#if timeLeft != null}
				<span class="px-3 py-1.5 rounded-full text-sm font-mono font-semibold shrink-0 {timeLeft < 120 ? 'bg-danger-50 text-danger-700' : 'bg-white border border-slate-200 text-slate-700'}">
					{formatTime(timeLeft)}
				</span>
			{/if}
		</div>
		<div class="mt-2 h-1.5 rounded-full bg-slate-200 overflow-hidden">
			<div class="h-full bg-primary-600 rounded-full transition-all" style="width: {progressPct}%"></div>
		</div>
		<div class="mt-1.5 flex items-center justify-between text-[11px] text-slate-400">
			<span>{answeredCount}/{questions.length} answered</span>
			{#if saving}
				<span>Saving...</span>
			{:else if autoSaveMsg}
				<span class="text-green-600">{autoSaveMsg}</span>
			{/if}
		</div>
	</div>

	<!-- Question navigator -->
	{#if questions.length > 1}
		<div class="flex gap-2 overflow-x-auto py-2 scrollbar-none" role="tablist" aria-label="Questions">
			{#each questionStatus as qs}
				<button role="tab" aria-selected={qs.index === currentIndex} aria-label="Question {qs.index + 1}"
					onclick={() => gotoQuestion(qs.index)}
					class="min-w-10 min-h-10 w-10 h-10 text-sm font-semibold rounded-full border shrink-0 transition-colors active:scale-95
						{qs.answered ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-slate-600 border-slate-300'}
						{qs.index === currentIndex ? 'ring-2 ring-primary-300 ring-offset-1' : ''}">
					{qs.index + 1}
				</button>
			{/each}
		</div>
	{/if}

	{#if loading}
		<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
	{:else if !currentQuestion}
		<EmptyState title="No questions found in this quiz." />
	{:else}
		<!-- Question card -->
		<div class="flex-1 bg-white rounded-2xl border border-slate-200 p-5 space-y-5 shadow-sm">
			<div class="flex items-start justify-between gap-3">
				<p class="text-base text-slate-900 leading-relaxed font-medium"><MathText text={currentQuestion.question_text} /></p>
				<span class="shrink-0 text-xs text-slate-400 bg-slate-100 rounded-full px-2 py-0.5">{currentQuestion.marks ?? 1} mark{(currentQuestion.marks ?? 1) !== 1 ? 's' : ''}</span>
			</div>

			{#if currentQuestion.question_type === 'mcq'}
				<div class="space-y-3" role="radiogroup" aria-label="Answer options">
					{#each currentOptions as opt, i}
						{@const selected = responses[currentQuestion.question_id]?.selected_options?.includes(opt.key)}
						<label class="flex items-center gap-3 p-4 min-h-[60px] rounded-2xl border-2 cursor-pointer transition-colors active:scale-[0.99]
							{selected ? 'border-primary-500 bg-primary-50' : 'border-slate-200 bg-white active:border-slate-300'}">
							<input type="radio" name={currentQuestion.question_id}
								checked={selected}
								class="sr-only"
								onchange={() => selectOption(currentQuestion.question_id, opt.key)} />
							<span class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold shrink-0 transition-colors
								{selected ? 'bg-primary-600 text-white' : 'bg-slate-100 text-slate-500'}">
								{selected ? '✓' : String.fromCharCode(65 + i)}
							</span>
							<span class="text-[15px] text-slate-800 flex-1"><MathText text={opt.value} /></span>
						</label>
					{/each}
				</div>
			{:else if currentQuestion.question_type === 'true_false'}
				<div class="flex gap-3">
					<button onclick={() => setText(currentQuestion.question_id, 'true')}
						class="flex-1 min-h-[64px] rounded-2xl border-2 text-base font-semibold transition-colors active:scale-[0.99]
							{responses[currentQuestion.question_id]?.text_answer === 'true' ? 'border-primary-500 bg-primary-50 text-primary-700' : 'border-slate-200 text-slate-600'}">
						True
					</button>
					<button onclick={() => setText(currentQuestion.question_id, 'false')}
						class="flex-1 min-h-[64px] rounded-2xl border-2 text-base font-semibold transition-colors active:scale-[0.99]
							{responses[currentQuestion.question_id]?.text_answer === 'false' ? 'border-primary-500 bg-primary-50 text-primary-700' : 'border-slate-200 text-slate-600'}">
						False
					</button>
				</div>
			{:else if currentQuestion.question_type === 'fill_blank' || currentQuestion.question_type === 'short_answer'}
				<textarea value={responses[currentQuestion.question_id]?.text_answer ?? ''}
					oninput={(e) => setText(currentQuestion.question_id, (e.target as HTMLTextAreaElement).value)}
					rows={4} placeholder="Type your answer..."
					class="w-full px-4 py-3 rounded-2xl border-2 border-slate-200 text-base resize-none focus:outline-none focus:border-primary-400"></textarea>
			{/if}
		</div>

		<!-- Bottom action bar -->
		<div class="sticky bottom-0 -mx-4 px-4 pt-2 bg-gradient-to-t from-slate-50 via-slate-50 to-transparent md:mx-0 md:px-0" style="padding-bottom: calc(0.75rem + env(safe-area-inset-bottom));">
			<div class="flex gap-3">
				<Button
					variant="secondary"
					icon={ChevronLeft}
					size="lg"
					onclick={() => gotoQuestion(currentIndex - 1)}
					disabled={currentIndex === 0}
				>
					Prev
				</Button>
				{#if isLast}
					<Button
						icon={Check}
						size="lg"
						class="flex-1"
						onclick={() => (showSubmit = true)}
						disabled={submitting}
					>
						Submit Quiz
					</Button>
				{:else}
					<Button
						size="lg"
						class="flex-1"
						onclick={() => gotoQuestion(currentIndex + 1)}
					>
						Next <ChevronRight size={18} />
					</Button>
				{/if}
			</div>
		</div>
	{/if}
</div>

<Modal bind:open={showSubmit} title="Submit quiz?">
	<p class="text-sm text-slate-600">
		You answered <strong>{answeredCount} of {questions.length}</strong> questions.
		{#if answeredCount < questions.length}
			<span class="text-amber-600 font-medium"> {questions.length - answeredCount} left unanswered.</span>
		{/if}
		This cannot be undone.
	</p>
	{#snippet footer()}
		<Button variant="ghost" onclick={() => (showSubmit = false)}>Keep answering</Button>
		<Button onclick={submitQuiz} loading={submitting} disabled={submitting}>
			{submitting ? 'Submitting...' : 'Submit'}
		</Button>
	{/snippet}
</Modal>
