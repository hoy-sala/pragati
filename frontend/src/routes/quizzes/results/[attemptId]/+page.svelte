<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { QuizResultData } from '$lib/types';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	const attemptId = $page.params.attemptId;

	let result = $state<QuizResultData | null>(null);
	let loading = $state(true);

	onMount(async () => {
		const res = await api<QuizResultData>('GET', `/quizzes/attempts/${attemptId}/result`);
		if (res.data) result = res.data;
		loading = false;
	});

	let passed = $derived(result?.attempt.passed ?? false);
	let scorePct = $derived(result ? Math.round((result.total_awarded / result.total_marks) * 100) : 0);
</script>

<div class="max-w-3xl mx-auto space-y-6">
	{#if loading}
		<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
	{:else if !result}
		<div class="p-8 text-center text-sm text-slate-400">Result not found.</div>
	{:else}
		<div class="bg-white rounded-xl border border-slate-200 p-6 text-center space-y-3">
			<div class="inline-flex items-center justify-center w-20 h-20 rounded-full {passed ? 'bg-green-100' : 'bg-danger-100'}">
				<span class="text-3xl font-bold {passed ? 'text-green-700' : 'text-danger-700'}">{scorePct}%</span>
			</div>
			<h1 class="text-2xl font-bold text-slate-900">{result.quiz.title}</h1>
			<div class="flex justify-center gap-6 text-sm text-slate-500">
				<span>Score: {result.total_awarded}/{result.total_marks}</span>
				<span>Pass: {result.quiz.pass_pct}%</span>
				<span>{passed ? 'Passed' : 'Failed'}</span>
			</div>
		</div>

		<div class="space-y-3">
			<h2 class="text-lg font-semibold text-slate-900">Responses</h2>
			{#each result.responses as r (r.id)}
				<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
					<div class="flex items-start justify-between gap-4">
						<p class="text-sm text-slate-900 leading-relaxed">{r.question_text}</p>
						<span class="shrink-0 text-xs text-slate-400">{r.marks_awarded}/{r.marks_total}</span>
					</div>

					{#if r.question_type === 'mcq'}
						<div class="space-y-1">
							{#each r.options || [] as opt}
								<div class="text-sm px-3 py-1.5 rounded-lg
									{opt.correct ? 'bg-green-50 text-green-800' : ''}
									{(r.selected_options || []).includes(opt.key) && !opt.correct ? 'bg-danger-50 text-danger-800' : 'text-slate-600'}">
									{opt.value}
									{#if opt.correct}<span class="text-xs ml-1">(correct)</span>{/if}
									{#if (r.selected_options || []).includes(opt.key) && !opt.correct}<span class="text-xs ml-1">(your answer)</span>{/if}
								</div>
							{/each}
						</div>
					{:else if r.question_type === 'true_false'}
						<div class="flex gap-3 text-sm">
							<span class="px-3 py-1 rounded-lg {r.text_answer === 'true' ? (r.is_correct ? 'bg-green-50 text-green-700' : 'bg-danger-50 text-danger-700') : ''}">
								True
							</span>
							<span class="px-3 py-1 rounded-lg {r.text_answer === 'false' ? (r.is_correct ? 'bg-green-50 text-green-700' : 'bg-danger-50 text-danger-700') : ''}">
								False
							</span>
							<span class="text-slate-400">| Correct: {r.correct_answer}</span>
						</div>
					{:else}
						<div class="space-y-1 text-sm">
							<p><span class="text-slate-500">Your answer:</span> {r.text_answer || '(no answer)'}</p>
							<p><span class="text-slate-500">Correct answer:</span> {r.correct_answer}</p>
						</div>
					{/if}
				</div>
			{/each}
		</div>

		<div class="flex gap-3">
			<a href="/quizzes/available" class="px-4 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50">
				Back to Quizzes
			</a>
		</div>
	{/if}
</div>
