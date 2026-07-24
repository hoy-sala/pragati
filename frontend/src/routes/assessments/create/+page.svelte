<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Class, Subject, AcademicYear, AssessmentCategory } from '$lib/types';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let categories = $state<AssessmentCategory[]>([]);
	let classes = $state<Class[]>([]);
	let subjects = $state<Subject[]>([]);
	let years = $state<AcademicYear[]>([]);

	let selectedCategory = $state('');
	let selectedSubject = $state('');
	let selectedClass = $state('');
	let section = $state('');
	let selectedYear = $state('');
	let name = $state('');
	let maxMarks = $state(100);
	let weightage = $state(100);
	let date = $state('');
	let chapters = $state('');
	let saving = $state(false);
	let error = $state('');
	let success = $state(false);

	onMount(async () => {
		const [catRes, classRes, subRes, yrRes] = await Promise.all([
			api<AssessmentCategory[]>('GET', '/assessment-categories'),
			api<Class[]>('GET', '/classes'),
			api<Subject[]>('GET', '/subjects'),
			api<AcademicYear[]>('GET', '/academic-years'),
		]);
		if (catRes.data) categories = catRes.data;
		if (classRes.data) classes = classRes.data;
		if (subRes.data) subjects = subRes.data;
		if (yrRes.data) {
			years = yrRes.data;
			const current = yrRes.data.find((y: { is_current: boolean }) => y.is_current);
			if (current) selectedYear = current.id;
		}
	});

	async function handleSubmit() {
		error = '';
		if (!selectedCategory || !selectedSubject || !selectedClass || !selectedYear) {
			error = 'Please fill all required fields.';
			return;
		}
		if (maxMarks <= 0) {
			error = 'Max marks must be positive.';
			return;
		}
		saving = true;
		const res = await api('POST', '/assessments', {
			category_id: selectedCategory,
			subject_id: selectedSubject,
			class_id: selectedClass,
			section_id: section,
			name: name,
			max_marks: maxMarks,
			weightage: weightage,
			date: date,
			chapters: chapters ? chapters.split(',').map(s => s.trim()) : [],
			academic_year_id: selectedYear
		});
		saving = false;
		if (res.error) {
			error = res.error.message;
		} else if (res.data) {
			success = true;
			setTimeout(() => goto('/assessments'), 1000);
		}
	}
</script>

<div class="max-w-2xl mx-auto space-y-6">
	<div class="flex items-center gap-4">
		<a href="/assessments" class="text-slate-400 hover:text-slate-600 transition-colors">&larr; Back</a>
		<h1 class="text-2xl font-bold text-slate-900">Create Assessment</h1>
	</div>

	{#if success}
		<div class="bg-green-50 border border-green-200 rounded-lg p-4 text-green-700 text-sm">Assessment created!</div>
	{/if}

	<form class="bg-white rounded-xl border border-slate-200 p-6 space-y-4" onsubmit={handleSubmit}>
		<div class="grid grid-cols-2 gap-4">
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Category *</label>
				<select bind:value={selectedCategory} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500">
					<option value="">Select</option>
					{#each categories as cat}
						<option value={cat.id}>{cat.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Subject *</label>
				<select bind:value={selectedSubject} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					{#each subjects as sub}
						<option value={sub.id}>{sub.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Class *</label>
				<select bind:value={selectedClass} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					{#each classes as cl}
						<option value={cl.id}>{cl.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Section</label>
				<input type="text" bind:value={section} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" placeholder="optional">
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Academic Year *</label>
				<select bind:value={selectedYear} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					{#each years as y}
						<option value={y.id}>{y.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Date</label>
				<input type="date" bind:value={date} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Max Marks</label>
				<input type="number" bind:value={maxMarks} min="1" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Weightage %</label>
				<input type="number" bind:value={weightage} min="0" max="100" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
		</div>

		<div>
			<label class="block text-sm font-medium text-slate-700 mb-1">Name</label>
			<input type="text" bind:value={name} placeholder="Term Exam, Unit Test, etc." class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
		</div>

		<div>
			<label class="block text-sm font-medium text-slate-700 mb-1">Chapters</label>
			<input type="text" bind:value={chapters} placeholder="separate with commas" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
		</div>

		{#if error}
			<div class="text-sm text-danger-600 bg-danger-50 rounded-lg p-3">{error}</div>
		{/if}

		<button type="submit" disabled={saving} class="w-full py-2 px-4 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
			{saving ? 'Creating...' : 'Create Assessment'}
		</button>
	</form>
</div>
