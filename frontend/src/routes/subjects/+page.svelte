<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Subject } from '$lib/types';
	import { onMount } from 'svelte';

	let subjects: Subject[] = $state([]);
	let loading = $state(true);
	let showForm = $state(false);
	let saving = $state(false);
	let error = $state('');

	let newName = $state('');
	let newCode = $state('');
	let newIsLanguage = $state(false);
	let newIsCore = $state(true);

	onMount(async () => {
		const res = await api<Subject[]>('GET', '/subjects?limit=100');
		if (res.data) subjects = res.data;
		loading = false;
	});

	async function createSubject() {
		if (!newName.trim()) return;
		saving = true;
		error = '';
		const res = await api<{ id: string }>('POST', '/subjects', {
			name: newName.trim(),
			code: newCode.trim() || undefined,
			is_language: newIsLanguage,
			is_core: newIsCore
		});
		saving = false;
		if (res.error) {
			error = res.error.message;
			return;
		}
		subjects = [...subjects, {
			id: res.data!.id,
			school_id: '',
			name: newName.trim(),
			code: newCode.trim() || undefined,
			is_language: newIsLanguage,
			is_core: newIsCore
		}];
		newName = '';
		newCode = '';
		newIsLanguage = false;
		newIsCore = true;
		showForm = false;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Subjects</h1>
			<p class="text-sm text-slate-500 mt-1">{subjects.length} subjects</p>
		</div>
		<button onclick={() => showForm = !showForm} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			{showForm ? 'Cancel' : 'Add Subject'}
		</button>
	</div>

	{#if showForm}
		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
			<div class="grid grid-cols-1 sm:grid-cols-4 gap-3">
				<input bind:value={newName} placeholder="Subject name (e.g. Mathematics)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newCode} placeholder="Code (e.g. MATH)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<label class="flex items-center gap-2 text-sm">
					<input type="checkbox" bind:checked={newIsLanguage} class="rounded border-slate-300" />
					Language
				</label>
				<label class="flex items-center gap-2 text-sm">
					<input type="checkbox" bind:checked={newIsCore} class="rounded border-slate-300" />
					Core subject
				</label>
			</div>
			{#if error}
				<div class="text-sm text-danger-600">{error}</div>
			{/if}
			<button onclick={createSubject} disabled={saving || !newName.trim()} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{saving ? 'Saving...' : 'Create'}
			</button>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		<table class="w-full text-sm">
			<thead>
				<tr class="bg-slate-50 text-slate-600">
					<th class="text-left px-4 py-3 font-medium">Name</th>
					<th class="text-left px-4 py-3 font-medium">Code</th>
					<th class="text-left px-4 py-3 font-medium">Type</th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr><td colspan="3" class="px-4 py-8 text-center text-slate-400">Loading...</td></tr>
				{:else if subjects.length === 0}
					<tr><td colspan="3" class="px-4 py-8 text-center text-slate-400">No subjects yet. Add one above.</td></tr>
				{:else}
					{#each subjects as s (s.id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="px-4 py-3 font-medium">{s.name}</td>
							<td class="px-4 py-3 text-slate-500">{s.code || '—'}</td>
							<td class="px-4 py-3">
								{#if s.is_language}
									<span class="inline-block px-2 py-0.5 rounded-full text-xs bg-blue-100 text-blue-700">Language</span>
								{:else if s.is_core}
									<span class="inline-block px-2 py-0.5 rounded-full text-xs bg-green-100 text-green-700">Core</span>
								{:else}
									<span class="inline-block px-2 py-0.5 rounded-full text-xs bg-slate-100 text-slate-600">Elective</span>
								{/if}
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>
