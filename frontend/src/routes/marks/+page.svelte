<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Assessment, AssessmentCategory, Class, Subject, AcademicYear, MarkGridRow, MarkInput } from '$lib/types';
	import { onMount } from 'svelte';
	import { TabulatorFull as Tabulator } from 'tabulator-tables';
	import 'tabulator-tables/dist/css/tabulator.min.css';

	let categories = $state<AssessmentCategory[]>([]);
	let classes = $state<Class[]>([]);
	let subjects = $state<Subject[]>([]);
	let years = $state<AcademicYear[]>([]);
	let assessments = $state<Assessment[]>([]);

	let selectedCategory = $state('');
	let selectedClass = $state('');
	let selectedSubject = $state('');
	let selectedAssessment = $state('');

	let students = $state<MarkGridRow[]>([]);
	let maxMarks = $state(100);
	let gridVersion = $state(1);
	let loading = $state(false);
	let saving = $state(false);
	let statusMsg = $state('');
	let tableEl = $state<HTMLDivElement>();
	let table = $state<Tabulator | null>(null);

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
			const current = yrRes.data.find(y => y.is_current);
		}
	});

	$effect(() => {
		if (selectedCategory || selectedClass || selectedSubject) {
			loadAssessments();
		}
	});

	async function loadAssessments() {
		const params = new URLSearchParams();
		if (selectedCategory) params.set('category_id', selectedCategory);
		if (selectedClass) params.set('class_id', selectedClass);
		if (selectedSubject) params.set('subject_id', selectedSubject);
		const res = await api<Assessment[]>('GET', '/assessments?' + params.toString());
		if (res.data) assessments = res.data;
	}

	async function loadGrid() {
		if (!selectedAssessment) return;
		loading = true;
		statusMsg = '';
		const res = await api<{ assessment: Assessment; students: MarkGridRow[] }>('GET', '/marks/grid?assessment_id=' + selectedAssessment);
		loading = false;
		if (res.error) {
			statusMsg = 'Error: ' + res.error.message;
			return;
		}
		if (res.data) {
			students = res.data.students;
			maxMarks = res.data.assessment.max_marks;
			gridVersion = res.data.assessment.version;
			initTable();
		}
	}

	function initTable() {
		if (table) { table.destroy(); table = null; }
		if (!tableEl) return;

		table = new Tabulator(tableEl, {
			data: students.map(s => ({
				...s,
				_marks: s.is_absent ? '' : (s.marks_obtained >= 0 ? s.marks_obtained : ''),
				_absent: s.is_absent
			})),
			layout: 'fitColumns',
			height: 'calc(100vh - 280px)',
			selectable: false,
			clipboard: true,
			columns: [
				{ title: '#', field: 'roll_no', width: 60, hozAlign: 'center', frozen: true },
				{ title: 'SATS', field: 'sats_number', width: 100, frozen: true },
				{ title: 'Name', field: 'name', width: 200, frozen: true },
				{
					title: 'Marks / ' + maxMarks,
					field: '_marks',
					editor: 'input',
					editorParams: { elementAttributes: { inputmode: 'decimal' } },
					width: 120,
					hozAlign: 'center',
					cellEdited: (cell: any) => {
						const val = cell.getValue();
						const row = cell.getRow().getData();
						cell.getRow().update({ is_absent: false, marks_obtained: val === '' ? 0 : parseFloat(val) || 0 });
					},
					formatter: (cell: any) => {
						const data = cell.getRow().getData();
						if (data._absent) return 'ABS';
						const v = cell.getValue();
						return v === '' || v === null || v === undefined ? '' : String(v);
					}
				},
				{
					title: 'Absent',
					field: '_absent',
					editor: 'tickCross',
					width: 80,
					hozAlign: 'center',
					cellEdited: (cell: any) => {
						const val = cell.getValue();
						cell.getRow().update({ is_absent: val, marks_obtained: 0 });
						if (val) cell.getRow().getCell('_marks').setValue('');
					}
				},
				{
					title: 'Remarks',
					field: 'remarks',
					editor: 'input',
					width: 200
				}
			],
			keybindings: {
				navUp: true,
				navDown: true,
				navLeft: true,
				navRight: true,
				navPrev: true,
				navNext: true,
				undo: true,
				redo: true,
				clipboardCopy: true,
				clipboardPaste: true,
			},
			clipboardPasteParser: (clipboard: string) => {
				const rows = clipboard.split('\n').filter(r => r.trim());
				return rows.map(row => row.split('\t'));
			},
			clipboardPasteAction: 'range',
		});
	}

	async function saveMarks() {
		if (!table || !selectedAssessment) return;
		saving = true;
		statusMsg = '';

		const data = table.getData() as any[];
		const marks: MarkInput[] = data.map((d: any) => ({
			student_id: d.student_id,
			marks_obtained: d._absent ? 0 : (parseFloat(d._marks) || 0),
			is_absent: !!d._absent,
			remarks: d.remarks || ''
		}));

		const res = await api('PUT', '/marks/batch', {
			assessment_id: selectedAssessment,
			version: gridVersion,
			marks
		});

		saving = false;
		if (res.error) {
			statusMsg = res.error.code === 'VERSION_CONFLICT'
				? 'Conflict: marks were updated by another user. Please refresh.'
				: 'Error: ' + res.error.message;
		} else if (res.data) {
			const d = res.data as any;
			statusMsg = `Saved ${d.updated} marks` + (d.errors?.length ? ` (${d.errors.length} errors)` : '');
			gridVersion = (gridVersion || 0) + 1;
		}
	}

	function resetForm() {
		selectedAssessment = '';
		students = [];
		if (table) { table.destroy(); table = null; }
		statusMsg = '';
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h1 class="text-xl font-bold text-slate-900">Marks Entry</h1>
		<div class="flex gap-2">
			<button onclick={saveMarks} disabled={saving || !selectedAssessment || !table}
				class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{saving ? 'Saving...' : 'Save Marks'}
			</button>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-4">
		<div class="flex flex-wrap gap-3 items-end">
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Category</label>
				<select bind:value={selectedCategory} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					<option value="">All</option>
					{#each categories as c}
						<option value={c.id}>{c.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Class</label>
				<select bind:value={selectedClass} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					<option value="">All</option>
					{#each classes as c}
						<option value={c.id}>{c.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Subject</label>
				<select bind:value={selectedSubject} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
					<option value="">All</option>
					{#each subjects as s}
						<option value={s.id}>{s.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Assessment</label>
				<select bind:value={selectedAssessment} onchange={loadGrid} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm min-w-[160px]">
					<option value="">Select</option>
					{#each assessments as a}
						<option value={a.id}>{a.name || a.category_id}</option>
					{/each}
				</select>
			</div>
			<button onclick={resetForm} class="px-3 py-1.5 text-sm text-slate-500 hover:text-slate-700 transition-colors">Clear</button>
		</div>
	</div>

	{#if statusMsg}
		<div class="text-sm px-4 py-2 rounded-lg bg-slate-100 text-slate-700">{statusMsg}</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		{#if loading}
			<div class="p-8 text-center text-sm text-slate-400">Loading...</div>
		{:else if selectedAssessment}
			<div bind:this={tableEl} class="marks-grid"></div>
		{:else}
			<div class="p-8 text-center text-sm text-slate-400">
				Select an assessment from the filters above to enter marks.
			</div>
		{/if}
	</div>

	<div class="text-xs text-slate-400 flex gap-4">
		<span>&larr; &rarr; &uarr; &darr; Navigate</span>
		<span>Tab / Enter: Next cell</span>
		<span>Ctrl+C / Ctrl+V: Copy/Paste</span>
	</div>
</div>

<style>
	:global(.marks-grid .tabulator) {
		border: none;
		border-radius: 0;
	}
	:global(.marks-grid .tabulator .tabulator-header) {
		background: #f8fafc;
		border-bottom: 1px solid #e2e8f0;
	}
	:global(.marks-grid .tabulator .tabulator-row) {
		border-bottom: 1px solid #f1f5f9;
	}
	:global(.marks-grid .tabulator .tabulator-row .tabulator-cell) {
		padding: 6px 8px;
		border-right: 1px solid #f1f5f9;
	}
	:global(.marks-grid .tabulator .tabulator-row.tabulator-row-even) {
		background: #fafafa;
	}
	:global(.marks-grid .tabulator .tabulator-cell.tabulator-editing) {
		border: 2px solid #2563eb;
	}
</style>
