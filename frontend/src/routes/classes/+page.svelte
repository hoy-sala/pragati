<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Class, AcademicYear } from '$lib/types';
	import { onMount } from 'svelte';

	let classes: Class[] = $state([]);
	let academicYears: AcademicYear[] = $state([]);
	let loading = $state(true);
	let showForm = $state(false);
	let saving = $state(false);
	let error = $state('');

	let newName = $state('');
	let newCode = $state('');
	let newSortOrder = $state(0);
	let newAcademicYearId = $state('');

	onMount(async () => {
		const [classRes, yearRes] = await Promise.all([
			api<Class[]>('GET', '/classes?limit=50'),
			api<AcademicYear[]>('GET', '/academic-years?limit=50')
		]);
		if (classRes.data) classes = classRes.data;
		if (yearRes.data) {
			academicYears = yearRes.data;
			const current = yearRes.data.find(y => y.is_current);
			if (current) newAcademicYearId = current.id;
		}
		loading = false;
	});

	async function createClass() {
		if (!newName.trim()) return;
		saving = true;
		error = '';
		const res = await api<{ id: string }>('POST', '/classes', {
			name: newName.trim(),
			code: newCode.trim() || undefined,
			sort_order: newSortOrder,
			academic_year_id: newAcademicYearId || undefined
		});
		saving = false;
		if (res.error) {
			error = res.error.message;
			return;
		}
		classes = [...classes, {
			id: res.data!.id,
			school_id: '',
			name: newName.trim(),
			code: newCode.trim() || undefined,
			sort_order: newSortOrder,
			created_at: new Date().toISOString()
		}];
		newName = '';
		newCode = '';
		newSortOrder = classes.length;
		showForm = false;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Classes</h1>
			<p class="text-sm text-slate-500 mt-1">{classes.length} classes</p>
		</div>
		<button onclick={() => showForm = !showForm} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			{showForm ? 'Cancel' : 'Add Class'}
		</button>
	</div>

	{#if showForm}
		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
			<div class="grid grid-cols-1 sm:grid-cols-4 gap-3">
				<input bind:value={newName} placeholder="Class name (e.g. Class 6)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newCode} placeholder="Code (e.g. VI)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newSortOrder} type="number" placeholder="Sort order" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<select bind:value={newAcademicYearId} class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500">
					<option value="">Select academic year</option>
					{#each academicYears as y (y.id)}
						<option value={y.id}>{y.name}{y.is_current ? ' (current)' : ''}</option>
					{/each}
				</select>
			</div>
			{#if error}
				<div class="text-sm text-danger-600">{error}</div>
			{/if}
			<button onclick={createClass} disabled={saving || !newName.trim()} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{saving ? 'Saving...' : 'Create'}
			</button>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		<table class="w-full text-sm">
			<thead>
				<tr class="bg-slate-50 text-slate-600">
					<th class="text-left px-4 py-3 font-medium">#</th>
					<th class="text-left px-4 py-3 font-medium">Name</th>
					<th class="text-left px-4 py-3 font-medium">Code</th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr><td colspan="3" class="px-4 py-8 text-center text-slate-400">Loading...</td></tr>
				{:else if classes.length === 0}
					<tr><td colspan="3" class="px-4 py-8 text-center text-slate-400">No classes yet. Add one above.</td></tr>
				{:else}
					{#each classes as c, i (c.id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="px-4 py-3 text-slate-400">{c.sort_order || i + 1}</td>
							<td class="px-4 py-3 font-medium">{c.name}</td>
							<td class="px-4 py-3 text-slate-500">{c.code || '—'}</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>
