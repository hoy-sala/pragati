<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Class, Subject } from '$lib/types';
	import { onMount } from 'svelte';

	let classes = $state<Class[]>([]);
	let subjects = $state<Subject[]>([]);
	let selectedClass = $state('');
	let selectedSubject = $state('');
	let loText = $state('');
	let importing = $state(false);
	let statusMsg = $state('');

	const sampleLO = `code,description,domain,expected_level,sort_order
LA-6.1.1,Reads and comprehends a variety of texts,cognitive,3,1
LA-6.1.2,Identifies main ideas and supporting details,cognitive,3,2
LA-6.2.1,Writes clear and coherent paragraphs,psychomotor,3,1`;

	onMount(async () => {
		const [cRes, sRes] = await Promise.all([
			api<Class[]>('GET', '/classes?limit=100'),
			api<Subject[]>('GET', '/subjects?limit=50')
		]);
		if (cRes.data) classes = cRes.data;
		if (sRes.data) subjects = sRes.data;
	});

	async function importLO() {
		if (!selectedClass || !selectedSubject || !loText.trim()) return;
		importing = true;
		statusMsg = 'Importing...';

		const lines = loText.trim().split('\n');
		const header = lines[0].toLowerCase().split(',');
		const codeIdx = header.indexOf('code');
		const descIdx = header.indexOf('description');
		const domainIdx = header.indexOf('domain');
		const levelIdx = header.indexOf('expected_level');
		const sortIdx = header.indexOf('sort_order');

		if (codeIdx === -1 || descIdx === -1) {
			statusMsg = 'CSV must include "code" and "description" columns.';
			importing = false;
			return;
		}

		const outcomes = [];
		for (let i = 1; i < lines.length; i++) {
			const cols = lines[i].split(',');
			if (cols.length < 2) continue;
			outcomes.push({
				code: cols[codeIdx]?.trim() || '',
				description: cols[descIdx]?.trim() || '',
				domain: domainIdx >= 0 ? (cols[domainIdx]?.trim() || 'cognitive') : 'cognitive',
				expected_level: levelIdx >= 0 ? parseInt(cols[levelIdx]?.trim()) || 3 : 3,
				sort_order: sortIdx >= 0 ? parseInt(cols[sortIdx]?.trim()) || 0 : i,
			});
		}

		if (outcomes.length === 0) {
			statusMsg = 'No valid outcomes found in input.';
			importing = false;
			return;
		}

		const res = await api('POST', '/hpc/learning-outcomes/import', {
			subject_id: selectedSubject,
			class_id: selectedClass,
			outcomes,
		});

		importing = false;
		if (res.data) {
			statusMsg = `Imported ${(res.data as any).imported} learning outcomes.`;
		} else if (res.error) {
			statusMsg = 'Error: ' + res.error.message;
		}
	}
</script>

<div class="max-w-3xl mx-auto space-y-4">
	<h1 class="text-xl font-bold text-slate-900">Import Learning Outcomes</h1>

	<div class="bg-white rounded-xl border border-slate-200 p-4">
		<div class="flex flex-wrap gap-3 items-end">
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Class</label>
				<select bind:value={selectedClass} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					{#each classes as c}
						<option value={c.id}>{c.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Subject</label>
				<select bind:value={selectedSubject} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					{#each subjects as s}
						<option value={s.id}>{s.name}</option>
					{/each}
				</select>
			</div>
			<button onclick={importLO} disabled={importing || !selectedSubject || !loText.trim()}
				class="px-4 py-1.5 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{importing ? 'Importing...' : 'Import'}
			</button>
		</div>
	</div>

	{#if statusMsg}
		<div class="text-sm px-4 py-2 rounded-lg bg-slate-100 text-slate-700">{statusMsg}</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-2">
		<label class="block text-sm font-medium text-slate-700">CSV Input</label>
		<p class="text-xs text-slate-500">Paste CSV with columns: code, description, domain, expected_level, sort_order</p>
		<textarea bind:value={loText} rows="10"
			class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono resize-none"
			placeholder={sampleLO}></textarea>
		<details class="text-xs text-slate-500">
			<summary class="cursor-pointer hover:text-slate-700">Show sample</summary>
			<pre class="mt-1 p-2 bg-slate-50 rounded text-xs">{sampleLO}</pre>
		</details>
	</div>
</div>
