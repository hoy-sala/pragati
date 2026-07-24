<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Class, Subject } from '$lib/types';
	import { onMount } from 'svelte';

	let classes = $state<Class[]>([]);
	let subjects = $state<Subject[]>([]);
	let selectedClass = $state('');
	let selectedSubject = $state('');
	let selectedTerm = $state('Term1');

	let loColumns = $state<any[]>([]);
	let students = $state<any[]>([]);
	let gridData = $state<any[]>([]);
	let loading = $state(false);
	let saving = $state(false);
	let statusMsg = $state('');

	const terms = ['Term1', 'Term2'];
	const proficiencyOptions = [
		{ value: 0, label: '—' },
		{ value: 1, label: '1-Beginning' },
		{ value: 2, label: '2-Developing' },
		{ value: 3, label: '3-Proficient' },
		{ value: 4, label: '4-Advanced' },
	];

	onMount(async () => {
		const [cRes, sRes] = await Promise.all([
			api<Class[]>('GET', '/classes?limit=100'),
			api<Subject[]>('GET', '/subjects?limit=50')
		]);
		if (cRes.data) classes = cRes.data;
		if (sRes.data) subjects = sRes.data;
	});

	async function loadGrid() {
		if (!selectedClass || !selectedSubject) return;
		loading = true;
		statusMsg = '';
		const params = `class_id=${selectedClass}&subject_id=${selectedSubject}&term=${selectedTerm}`;
		const res = await api<any>('GET', '/hpc/assessments?' + params);
		loading = false;
		if (res.data) {
			loColumns = res.data.columns || [];
			students = res.data.students || [];
			gridData = res.data.grid || [];
		} else if (res.error) {
			statusMsg = 'Error: ' + res.error.message;
		}
	}

	function getCellValue(studentId: string, loId: string): number {
		const row = gridData.find((r: any) => r.student.student_id === studentId);
		if (row && row.cells[loId]) return row.cells[loId].level || 0;
		return 0;
	}

	async function updateCell(studentId: string, loId: string, level: number) {
		const assessment = {
			learning_outcome_id: loId,
			proficiency_level: level
		};
		await api('POST', '/hpc/assess', {
			student_id: studentId,
			subject_id: selectedSubject,
			term: selectedTerm,
			assessments: [assessment]
		});
	}

	async function saveAll() {
		saving = true;
		statusMsg = 'Saving all assessments...';
		let total = 0;
		for (const row of gridData) {
			const assessments = loColumns
				.filter((lo: any) => row.cells[lo.id] && row.cells[lo.id].level > 0)
				.map((lo: any) => ({
					learning_outcome_id: lo.id,
					proficiency_level: row.cells[lo.id].level
				}));
			if (assessments.length === 0) continue;
			await api('POST', '/hpc/assess', {
				student_id: row.student.student_id,
				subject_id: selectedSubject,
				term: selectedTerm,
				assessments
			});
			total += assessments.length;
		}
		saving = false;
		statusMsg = `Saved ${total} assessments.`;
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h1 class="text-xl font-bold text-slate-900">Learning Outcome Assessment Grid</h1>
		<button onclick={saveAll} disabled={saving || loading || !selectedSubject}
			class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
			{saving ? 'Saving...' : 'Save All'}
		</button>
	</div>

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

	<div class="bg-white rounded-xl border border-slate-200 overflow-x-auto">
		{#if loading}
			<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
		{:else if loColumns.length === 0}
			<div class="p-8 text-center text-sm text-slate-400">
				{selectedSubject ? 'No learning outcomes configured for this subject. Import them first.' : 'Select a subject and class to load.'}
			</div>
		{:else}
			<table class="w-full text-sm">
				<thead>
					<tr class="bg-slate-50">
						<th class="sticky left-0 bg-slate-50 z-10 px-3 py-2 text-left font-medium text-xs" style="min-width:180px;">Student</th>
						{#each loColumns as lo}
							<th class="px-2 py-2 text-center font-medium text-xs" style="min-width:80px; max-width:100px;"
								title={lo.description}>
								{lo.code}
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each students as student}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="sticky left-0 bg-white z-10 px-3 py-1.5 font-medium text-xs">{student.name}</td>
							{#each loColumns as lo}
								<td class="px-1 py-1.5 text-center">
									<select class="text-xs px-1 py-0.5 border border-slate-300 rounded"
										value={getCellValue(student.student_id, lo.id)}
										onchange={(e) => { const el = e.target as HTMLSelectElement; updateCell(student.student_id, lo.id, parseInt(el.value)); }}>
										{#each proficiencyOptions as opt}
											<option value={opt.value}>{opt.label}</option>
										{/each}
									</select>
								</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
</div>
