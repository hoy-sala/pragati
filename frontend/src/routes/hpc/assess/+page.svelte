<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Class, Subject } from '$lib/types';
	import { onMount } from 'svelte';
	import Select from '$lib/components/Select.svelte';

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
	const termOptions = terms.map(t => ({ id: t, name: t }));
	const proficiencyOptions = [
		{ id: '0', name: '—' },
		{ id: '1', name: '1-Beginning' },
		{ id: '2', name: '2-Developing' },
		{ id: '3', name: '3-Proficient' },
		{ id: '4', name: '4-Advanced' },
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
				<label for="assess-class" class="block text-xs font-medium text-slate-600 mb-1">Class</label>
				<Select id="assess-class" bind:value={selectedClass} options={classes} placeholder="Select" />
			</div>
			<div>
				<label for="assess-subject" class="block text-xs font-medium text-slate-600 mb-1">Subject</label>
				<Select id="assess-subject" bind:value={selectedSubject} options={subjects} placeholder="Select" />
			</div>
			<div>
				<label for="assess-term" class="block text-xs font-medium text-slate-600 mb-1">Term</label>
				<Select id="assess-term" bind:value={selectedTerm} options={termOptions} />
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
								<Select
									value={String(getCellValue(student.student_id, lo.id))}
									options={proficiencyOptions}
									size="sm"
									onselect={(v) => updateCell(student.student_id, lo.id, parseInt(v))}
								/>
							</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
</div>
