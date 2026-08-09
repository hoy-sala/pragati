<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { QuizListItem } from '$lib/types';
	import { onMount } from 'svelte';
	import SearchFilter from '$lib/components/SearchFilter.svelte';
	import Pagination from '$lib/components/Pagination.svelte';

	let allQuizzes = $state<QuizListItem[]>([]);
	let loading = $state(true);
	let deleting = $state<string | null>(null);
	let search = $state('');
	let filterStatus = $state('');
	let page = $state(1);
	const pageSize = 15;

	let filteredQuizzes = $derived.by(() => {
		let result = allQuizzes;
		if (filterStatus === 'published') result = result.filter(q => q.is_published);
		if (filterStatus === 'draft') result = result.filter(q => !q.is_published);
		if (search.trim()) {
			const q = search.toLowerCase();
			result = result.filter(item => item.title?.toLowerCase().includes(q) || item.description?.toLowerCase().includes(q));
		}
		return result;
	});

	let paginatedQuizzes = $derived.by(() => {
		const start = (page - 1) * pageSize;
		return filteredQuizzes.slice(start, start + pageSize);
	});

	let totalQuizzes = $derived(filteredQuizzes.length);

	function onPageChange(p: number) { page = p; }
	function resetPage() { page = 1; }

	onMount(loadQuizzes);

	async function loadQuizzes() {
		loading = true;
		const res = await api<QuizListItem[]>('GET', '/quizzes');
		if (res.data) allQuizzes = res.data;
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
			<p class="text-sm text-slate-500 mt-1">{totalQuizzes} quizzes</p>
		</div>
		<a href="/quizzes/create" class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			Create Quiz
		</a>
	</div>

	<div class="bg-white rounded-xl border border-slate-200">
		<div class="p-4 border-b border-slate-200 flex flex-wrap gap-3">
			<div class="flex-1 min-w-48">
				<SearchFilter bind:value={search} placeholder="Search quizzes..." onInput={resetPage} />
			</div>
			<div class="w-36">
				<select bind:value={filterStatus} onchange={resetPage} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">All Status</option>
					<option value="published">Published</option>
					<option value="draft">Draft</option>
				</select>
			</div>
		</div>

		{#if loading}
			<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
		{:else if totalQuizzes === 0}
			<div class="p-8 text-center text-sm text-slate-400">No quizzes found.</div>
		{:else}
			<div class="divide-y divide-slate-100">
				{#each paginatedQuizzes as q (q.id)}
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
			</div>
			<Pagination {total} pageSize={pageSize} {page} onChange={onPageChange} />
		{/if}
	</div>
</div>
