<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Question, Subject } from '$lib/types';
	import { typeColors, typeLabels } from '$lib/utils/questionUtils';
	import { onMount } from 'svelte';

	let questions = $state<Question[]>([]);
	let subjects = $state<Subject[]>([]);
	let loading = $state(true);
	let filterSubject = $state('');
	let filterType = $state('');
	let search = $state('');

	onMount(async () => {
		const [subRes] = await Promise.all([api<Subject[]>('GET', '/subjects')]);
		if (subRes.data) subjects = subRes.data;
		loadQuestions();
	});

	async function loadQuestions() {
		loading = true;
		const params = new URLSearchParams();
		if (filterSubject) params.set('subject_id', filterSubject);
		if (filterType) params.set('type', filterType);
		if (search) params.set('search', search);

		const res = await api<Question[]>('GET', '/questions?' + params.toString());
		if (res.data) questions = res.data;
		loading = false;
	}

	function subjectName(id: string): string {
		return subjects.find(s => s.id === id)?.name || id.slice(0, 8);
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Question Bank</h1>
			<p class="text-sm text-slate-500 mt-1">{questions.length} questions</p>
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
			<select bind:value={filterSubject} onchange={loadQuestions} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
				<option value="">All Subjects</option>
				{#each subjects as s}
					<option value={s.id}>{s.name}</option>
				{/each}
			</select>
			<select bind:value={filterType} onchange={loadQuestions} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
				<option value="">All Types</option>
				<option value="mcq">MCQ</option>
				<option value="true_false">True/False</option>
				<option value="fill_blank">Fill Blank</option>
				<option value="short_answer">Short Answer</option>
			</select>
			<input type="text" bind:value={search} placeholder="Search questions..." oninput={loadQuestions}
				class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm flex-1 min-w-[200px]">
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 divide-y divide-slate-100">
		{#if loading}
			<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
		{:else if questions.length === 0}
			<div class="p-8 text-center text-sm text-slate-400">No questions yet. Create or import some.</div>
		{:else}
			{#each questions as q (q.id)}
				<div class="p-4 hover:bg-slate-50 transition-colors">
					<div class="flex items-start justify-between gap-4">
						<div class="flex-1 min-w-0">
							<p class="text-sm text-slate-900 leading-relaxed">{q.question_text}</p>
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
	</div>
</div>
