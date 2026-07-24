<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Class, AcademicYear } from '$lib/types';
	import { onMount } from 'svelte';

	interface HPCGridRow {
		student_id: string;
		sats_number: string;
		name: string;
		roll_no: number;
		entry_id: string;
		status: string;
		has_pdf: boolean;
	}

	let classes = $state<Class[]>([]);
	let years = $state<AcademicYear[]>([]);
	let grid = $state<HPCGridRow[]>([]);
	let selectedClass = $state('');
	let selectedYear = $state('');
	let selectedTerm = $state('Term1');
	let loading = $state(false);
	let statusMsg = $state('');
	let summary = $state({ total_students: 0, published_count: 0, draft_count: 0 });

	const terms = ['Term1', 'Term2'];

	onMount(async () => {
		const [cRes, yRes] = await Promise.all([
			api<Class[]>('GET', '/classes?limit=100'),
			api<AcademicYear[]>('GET', '/academic-years?limit=10')
		]);
		if (cRes.data) classes = cRes.data;
		if (yRes.data) {
			years = yRes.data;
			const current = yRes.data.find((y: { is_current: boolean }) => y.is_current);
			if (current) selectedYear = current.id;
		}
	});

	async function loadGrid() {
		if (!selectedClass || !selectedYear) return;
		loading = true;
		statusMsg = '';
		const params = `class_id=${selectedClass}&academic_year_id=${selectedYear}&term=${selectedTerm}`;
		const [gRes, rRes] = await Promise.all([
			api<HPCGridRow[]>('GET', '/hpc/grid?' + params),
			api('GET', '/hpc/reports/class?' + params)
		]);
		loading = false;
		if (gRes.data) grid = gRes.data;
		if (rRes.data) summary = rRes.data as any;
	}

	async function migrateFromMarks() {
		if (!selectedClass || !selectedYear) return;
		statusMsg = 'Migrating marks to HPC entries...';
		const res = await api<{ migrated: number }>('POST', '/hpc/migrate-from-marks?' +
			`class_id=${selectedClass}&academic_year_id=${selectedYear}&term=${selectedTerm}`);
		if (res.data) {
			statusMsg = `Migrated ${res.data.migrated} entries from marks.`;
			loadGrid();
		} else if (res.error) {
			statusMsg = 'Error: ' + res.error.message;
		}
	}

	async function publishAll() {
		const unpublished = grid.filter(r => r.entry_id && r.status === 'draft');
		if (unpublished.length === 0) { statusMsg = 'No draft entries to publish.'; return; }
		statusMsg = `Publishing ${unpublished.length} entries...`;
		for (const row of unpublished) {
			await api('POST', '/hpc/entries/publish', { entry_id: row.entry_id });
		}
		statusMsg = `Published ${unpublished.length} entries.`;
		loadGrid();
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h1 class="text-xl font-bold text-slate-900">Holistic Progress Card (HPC)</h1>
		<div class="flex gap-2">
			<button onclick={migrateFromMarks} disabled={!selectedClass}
				class="px-3 py-1.5 border border-slate-300 rounded-lg text-sm hover:bg-slate-50 transition-colors">
				Migrate from Marks
			</button>
			<button onclick={publishAll} disabled={!selectedClass}
				class="px-4 py-1.5 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				Publish All
			</button>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-4">
		<div class="flex flex-wrap gap-3 items-end">
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Class</label>
				<select bind:value={selectedClass} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					<option value="">Select Class</option>
					{#each classes as c}
						<option value={c.id}>{c.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Academic Year</label>
				<select bind:value={selectedYear} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					{#each years as y}
						<option value={y.id}>{y.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Term</label>
				<select bind:value={selectedTerm} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					{#each terms as t}
						<option value={t}>{t}</option>
					{/each}
				</select>
			</div>
			<button onclick={loadGrid}
				class="px-4 py-1.5 bg-slate-900 text-white rounded-lg text-sm font-medium hover:bg-slate-800 transition-colors">
				Load
			</button>
		</div>
	</div>

	{#if statusMsg}
		<div class="text-sm px-4 py-2 rounded-lg bg-slate-100 text-slate-700">{statusMsg}</div>
	{/if}

	{#if summary.total_students > 0}
		<div class="flex gap-4 text-sm">
			<span class="px-3 py-1 bg-slate-100 rounded-full">Total: <strong>{summary.total_students}</strong></span>
			<span class="px-3 py-1 bg-green-100 text-green-800 rounded-full">Published: <strong>{summary.published_count}</strong></span>
			<span class="px-3 py-1 bg-amber-100 text-amber-800 rounded-full">Draft: <strong>{summary.draft_count}</strong></span>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		<table class="w-full text-sm">
			<thead>
				<tr class="bg-slate-50 text-slate-600">
					<th class="text-left px-4 py-3 font-medium">#</th>
					<th class="text-left px-4 py-3 font-medium">SATS</th>
					<th class="text-left px-4 py-3 font-medium">Name</th>
					<th class="text-left px-4 py-3 font-medium">Status</th>
					<th class="text-left px-4 py-3 font-medium">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr><td colspan="5" class="px-4 py-8 text-center text-slate-400">Loading...</td></tr>
				{:else if grid.length === 0}
					<tr><td colspan="5" class="px-4 py-8 text-center text-slate-400">
						{selectedClass ? 'No HPC entries found. Click "Migrate from Marks" to create.' : 'Select a class above.'}
					</td></tr>
				{:else}
					{#each grid as row (row.student_id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="px-4 py-3">{row.roll_no}</td>
							<td class="px-4 py-3 font-mono text-xs">{row.sats_number}</td>
							<td class="px-4 py-3 font-medium">{row.name}</td>
							<td class="px-4 py-3">
								{#if !row.entry_id}
									<span class="text-xs text-slate-400">—</span>
								{:else if row.status === 'published'}
									<span class="px-2 py-0.5 text-xs bg-green-100 text-green-700 rounded-full">Published</span>
								{:else}
									<span class="px-2 py-0.5 text-xs bg-amber-100 text-amber-700 rounded-full">Draft</span>
								{/if}
							</td>
							<td class="px-4 py-3">
								<div class="flex gap-2">
									<a href="/hpc/entry/{row.student_id}?term={selectedTerm}&year={selectedYear}"
										class="text-xs text-primary-600 hover:text-primary-800 font-medium">
										{row.entry_id ? 'Edit' : 'Create'}
									</a>
									{#if row.has_pdf}
										<a href="/api/v1/hpc/entries/{row.entry_id}/pdf" target="_blank"
											class="text-xs text-slate-500 hover:text-slate-700">PDF</a>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>
