<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Question, Subject } from '$lib/types';
	import { typeColors, typeLabels } from '$lib/utils/questionUtils';
	import { onMount } from 'svelte';
	import Select from '$lib/components/Select.svelte';
	import SearchFilter from '$lib/components/SearchFilter.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import MathText from '$lib/components/MathText.svelte';

	let allQuestions = $state<Question[]>([]);
	let subjects = $state<Subject[]>([]);
	let loading = $state(true);
	let filterSubject = $state('');
	let filterType = $state('');
	let search = $state('');
	let page = $state(1);
	const pageSize = 20;

	let filteredQuestions = $derived(allQuestions.filter(q => {
		if (filterSubject && q.subject_id !== filterSubject) return false;
		if (filterType && q.question_type !== filterType) return false;
		if (!search.trim()) return true;
		return q.question_text?.toLowerCase().includes(search.toLowerCase());
	}));

	let paginatedQuestions = $derived(filteredQuestions.slice((page - 1) * pageSize, page * pageSize));
	let totalQuestions = $derived(filteredQuestions.length);

	function onPageChange(p: number) { page = p; }
	function resetPage() { page = 1; }

	onMount(async () => {
		const [subRes] = await Promise.all([api<Subject[]>('GET', '/subjects')]);
		if (subRes.data) subjects = subRes.data;
		loadQuestions();
	});

	async function loadQuestions() {
		loading = true;
		const res = await api<Question[]>('GET', '/questions?limit=500');
		if (res.data) allQuestions = res.data;
		loading = false;
	}

	function subjectName(id: string): string {
		return subjects.find(s => s.id === id)?.name || id.slice(0, 8);
	}

	const typeOptions = [
		{ id: 'mcq', name: 'MCQ' },
		{ id: 'true_false', name: 'True/False' },
		{ id: 'fill_blank', name: 'Fill Blank' },
		{ id: 'short_answer', name: 'Short Answer' },
	];
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Question Bank</h1>
			<p class="text-sm text-slate-500 mt-1">{totalQuestions} questions</p>
		</div>
		<div class="flex gap-2">
			<a href="/questions/import" class="px-4 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50 transition-colors">
				Import
			</a>
			<a href="/questions/create" class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
				New Question
			</a>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-4">
		<div class="flex flex-wrap gap-3">
			<div class="flex-1 min-w-48">
				<SearchFilter bind:value={search} placeholder="Search questions..." onInput={resetPage} />
			</div>
			<div class="w-40">
				<Select bind:value={filterSubject} options={subjects} placeholder="All Subjects" onselect={resetPage} />
			</div>
			<div class="w-36">
				<Select bind:value={filterType} options={typeOptions} placeholder="All Types" onselect={resetPage} />
			</div>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 divide-y divide-slate-100">
		{#if loading}
			<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
		{:else if totalQuestions === 0}
			<div class="p-8 text-center text-sm text-slate-400">No questions found.</div>
		{:else}
			{#each paginatedQuestions as q (q.id)}
				<div class="p-4 hover:bg-slate-50 transition-colors">
					<div class="flex items-start justify-between gap-4">
						<div class="flex-1 min-w-0">
							<p class="text-sm text-slate-900 leading-relaxed"><MathText text={q.question_text} /></p>
							<div class="flex flex-wrap gap-2 mt-2">
								<span class="text-xs px-2 py-0.5 rounded-full {typeColors[q.question_type] || 'bg-slate-100 text-slate-600'}">
									{typeLabels[q.question_type] || q.question_type}
								</span>
								<span class="text-xs text-slate-500">{q.difficulty}</span>
								<span class="text-xs text-slate-500">{subjectName(q.subject_id)}</span>
								<span class="text-xs text-slate-400">{q.marks} mark{q.marks !== 1 ? 's' : ''}</span>
							</div>
						</div>
						<div class="text-xs text-slate-400 whitespace-nowrap">
							{q.chapters?.length ? q.chapters.join(', ') : ''}
						</div>
					</div>
				</div>
			{/each}
		{/if}
		<Pagination total={totalQuestions} pageSize={pageSize} page={page} onChange={onPageChange} />
	</div>
</div>
