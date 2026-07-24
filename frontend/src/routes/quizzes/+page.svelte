<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { QuizListItem } from '$lib/types';
	import { onMount } from 'svelte';

	let quizzes = $state<QuizListItem[]>([]);
	let loading = $state(true);
	let deleting = $state<string | null>(null);

	onMount(loadQuizzes);

	async function loadQuizzes() {
		loading = true;
		const res = await api<QuizListItem[]>('GET', '/quizzes');
		if (res.data) quizzes = res.data;
		loading = false;
	}

	async function publish(id: string) {
		const res = await api('POST', `/quizzes/${id}/publish`);
		if (res.data) loadQuizzes();
	}

	async function remove(id: string) {
		if (!confirm('Delete this quiz?')) return;
		deleting = id;
		await api('DELETE', `/quizzes/${id}`);
		deleting = null;
		loadQuizzes();
	}

	function targetLabel(t: string): string {
		return t === 'student' ? 'Students' : 'Staff';
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Quizzes</h1>
			<p class="text-sm text-slate-500 mt-1">{quizzes.length} quizzes</p>
		</div>
		<a href="/quizzes/create" class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			Create Quiz
		</a>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 divide-y divide-slate-100">
		{#if loading}
			<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
		{:else if quizzes.length === 0}
			<div class="p-8 text-center text-sm text-slate-400">No quizzes yet.</div>
		{:else}
			{#each quizzes as q (q.id)}
				<div class="p-4 hover:bg-slate-50 transition-colors">
					<div class="flex items-center justify-between gap-4">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2">
								<h3 class="text-sm font-medium text-slate-900">{q.title}</h3>
								<span class="text-xs px-2 py-0.5 rounded-full {q.is_published ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}">
									{q.is_published ? 'Published' : 'Draft'}
								</span>
							</div>
							<p class="text-xs text-slate-500 mt-1 line-clamp-1">{q.description || 'No description'}</p>
							<div class="flex flex-wrap gap-3 mt-1.5 text-xs text-slate-400">
								<span>{targetLabel(q.target_type)}</span>
								<span>{q.question_count} questions</span>
								<span>{q.attempt_count} attempts</span>
								<span>Pass: {q.pass_pct}%</span>
								<span>by {q.created_by_name}</span>
							</div>
						</div>
						<div class="flex gap-1.5 shrink-0">
							{#if !q.is_published}
								<button onclick={() => publish(q.id)}
									class="px-3 py-1.5 text-xs font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700 transition-colors">
									Publish
								</button>
								<a href="/quizzes/{q.id}/edit"
									class="px-3 py-1.5 text-xs font-medium text-slate-700 border border-slate-300 rounded-lg hover:bg-slate-50 transition-colors">
									Edit
								</a>
								<button onclick={() => remove(q.id)}
									class="px-3 py-1.5 text-xs font-medium text-danger-600 border border-danger-300 rounded-lg hover:bg-danger-50 transition-colors">
									Delete
								</button>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>
