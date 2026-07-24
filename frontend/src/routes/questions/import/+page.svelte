<script lang="ts">
	import { api, apiUrl } from '$lib/api/client.svelte';
	import type { Subject } from '$lib/types';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let subjects = $state<Subject[]>([]);
	let subjectId = $state('');
	let importMode = $state('gift');
	let giftText = $state('');
	let csvFile = $state<File | null>(null);
	let importing = $state(false);
	let result = $state<{ imported: number; errors: { line: number; message: string }[] } | null>(null);
	let error = $state('');

	const openB = '{';
	const closeB = '}';
	const giftTF = `${openB}TRUE${closeB}`;
	const sampleGIFT = `// MCQ example
Who wrote the Indian National Anthem? {
    ~Mahatma Gandhi
    ~Jawaharlal Nehru
    =Rabindranath Tagore
    ~Bankim Chandra
}

// True/False
The Earth is flat. {FALSE}

// Fill in blank
The capital of India is {=New Delhi ~Mumbai ~Kolkata}.`;

	onMount(async () => {
		const res = await api<Subject[]>('GET', '/subjects');
		if (res.data) subjects = res.data;
	});

	async function handleImport() {
		error = '';
		result = null;
		if (!subjectId) { error = 'Select a subject.'; return; }
		importing = true;

		if (importMode === 'gift') {
			if (!giftText.trim()) { error = 'Paste GIFT format text.'; importing = false; return; }
			const params = new URLSearchParams({ subject_id: subjectId });
			const res = await fetch(apiUrl('/questions/import/gift?' + params.toString()), {
				method: 'POST',
				headers: { 'Authorization': 'Bearer ' + localStorage.getItem('access_token'), 'Content-Type': 'text/plain' },
				body: giftText
			});
			const json = await res.json();
			if (json.data) result = json.data;
			if (json.error) error = json.error.message;
		} else {
			if (!csvFile) { error = 'Select a CSV file.'; importing = false; return; }
			const formData = new FormData();
			formData.append('file', csvFile);
			const params = new URLSearchParams({ subject_id: subjectId });
			const res = await fetch(apiUrl('/questions/import/csv?' + params.toString()), {
				method: 'POST',
				headers: { 'Authorization': 'Bearer ' + localStorage.getItem('access_token') },
				body: formData
			});
			const json = await res.json();
			if (json.data) result = json.data;
			if (json.error) error = json.error.message;
		}

		importing = false;
	}
</script>

<div class="max-w-3xl mx-auto space-y-6">
	<div class="flex items-center gap-4">
		<a href="/questions" class="text-slate-400 hover:text-slate-600">&larr; Back</a>
		<h1 class="text-xl font-bold text-slate-900">Import Questions</h1>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-6 space-y-4">
		<div class="flex gap-4 items-end">
			<div class="flex-1">
				<label class="block text-sm font-medium text-slate-700 mb-1">Subject *</label>
				<select bind:value={subjectId} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					{#each subjects as s}
						<option value={s.id}>{s.name}</option>
					{/each}
				</select>
			</div>
			<div class="flex gap-1 bg-slate-100 rounded-lg p-1">
				<button onclick={() => importMode = 'gift'}
					class="px-3 py-1.5 text-sm rounded-md transition-colors {importMode === 'gift' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500'}">
					GIFT Format
				</button>
				<button onclick={() => importMode = 'csv'}
					class="px-3 py-1.5 text-sm rounded-md transition-colors {importMode === 'csv' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500'}">
					CSV File
				</button>
			</div>
		</div>

		{#if importMode === 'gift'}
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">GIFT Format Text</label>
				<textarea bind:value={giftText} rows="12" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm font-mono resize-none" placeholder={sampleGIFT}></textarea>
			</div>
			<details class="text-xs text-slate-500">
				<summary class="cursor-pointer hover:text-slate-700">GIFT format reference</summary>
				<pre class="mt-2 p-3 bg-slate-50 rounded-lg text-xs leading-relaxed">{sampleGIFT}</pre>
				<p class="mt-2">Rules: <code>~</code> wrong answer, <code>=</code> correct answer, <code>{'{TRUE}'}</code> or <code>{'{FALSE}'}</code> for T/F, <code>{'{"{=correct ~wrong}"}'}</code> for fill blank.</p>
			</details>
		{:else}
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">CSV File</label>
				<input type="file" accept=".csv" onchange={(e: Event) => { const el = e.target as HTMLInputElement; csvFile = el.files?.[0] || null; }}
					class="w-full text-sm file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-medium file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100">
			</div>
			<details class="text-xs text-slate-500">
				<summary class="cursor-pointer hover:text-slate-700">CSV format reference</summary>
				<p class="mt-2">Columns: <code>type, question, a, b, c, d, answer, marks, difficulty, chapter, tags</code></p>
				<p><code>type</code>: mcq, true_false, fill_blank, short_answer</p>
				<p><code>answer</code>: for MCQ use option key (A/B/C/D), for T/F use TRUE/FALSE, for others use text</p>
				<p><code>difficulty</code>: easy, medium, hard (default: medium)</p>
			</details>
		{/if}

		{#if error}
			<div class="text-sm text-danger-600 bg-danger-50 rounded-lg p-3">{error}</div>
		{/if}

		{#if result}
			<div class="text-sm bg-green-50 border border-green-200 rounded-lg p-4">
				<p class="text-green-700 font-medium">Imported {result.imported} question{result.imported !== 1 ? 's' : ''}</p>
				{#if result.errors && result.errors.length > 0}
					<ul class="mt-2 space-y-1">
						{#each result.errors as e}
							<li class="text-amber-700">Line {e.line}: {e.message}</li>
						{/each}
					</ul>
				{/if}
			</div>
		{/if}

		<div class="flex gap-3">
			<button onclick={handleImport} disabled={importing || !subjectId}
				class="px-6 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{importing ? 'Importing...' : 'Import'}
			</button>
			<button onclick={() => goto('/questions')} class="px-6 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm hover:bg-slate-50 transition-colors">
				Cancel
			</button>
		</div>
	</div>
</div>
