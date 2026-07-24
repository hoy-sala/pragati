<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { AvailableQuizItem } from '$lib/types';
	import { onMount } from 'svelte';

	let quizzes = $state<AvailableQuizItem[]>([]);
	let loading = $state(true);

	onMount(async () => {
		const res = await api<AvailableQuizItem[]>('GET', '/quizzes/available');
		if (res.data) quizzes = res.data;
		loading = false;
	});

	function statusBadge(status?: string): string {
		if (!status || status === 'graded') return 'bg-green-100 text-green-700';
		if (status === 'in_progress') return 'bg-blue-100 text-blue-700';
		return 'bg-slate-100 text-slate-600';
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-slate-900">Available Quizzes</h1>
		<p class="text-sm text-slate-500 mt-1">Take a quiz or continue an in-progress attempt.</p>
	</div>

	<div class="grid gap-4">
		{#if loading}
			<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
		{:else if quizzes.length === 0}
			<div class="p-8 text-center text-sm text-slate-400">No quizzes available right now.</div>
		{:else}
			{#each quizzes as q (q.id)}
				<div class="bg-white rounded-xl border border-slate-200 p-5 hover:border-primary-300 transition-colors">
					<div class="flex items-center justify-between gap-4">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2">
								<h3 class="text-base font-semibold text-slate-900">{q.title}</h3>
								{#if q.last_status}
									<span class="text-xs px-2 py-0.5 rounded-full {statusBadge(q.last_status)}">
										{q.last_status === 'in_progress' ? 'In Progress' : q.last_status === 'graded' ? 'Completed' : q.last_status}
									</span>
								{/if}
							</div>
							<p class="text-sm text-slate-500 mt-1 line-clamp-2">{q.description || ''}</p>
							<div class="flex flex-wrap gap-3 mt-2 text-xs text-slate-400">
								<span>{q.question_count} questions</span>
								<span>Pass: {q.pass_pct}%</span>
								<span>Attempts: {q.attempts_used}/{q.max_attempts}</span>
								{#if q.duration_min}<span>{q.duration_min} min</span>{/if}
								{#if q.last_score != null}<span>Last score: {q.last_score}%</span>{/if}
							</div>
						</div>
						<a href="/quizzes/take/{q.id}"
							class="shrink-0 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
							{q.last_status === 'in_progress' ? 'Continue' : 'Start'}
						</a>
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>
