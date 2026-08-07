<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { Save, X, BookOpen, Users, ClipboardCheck, Filter, Table } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import Select from '$lib/components/Select.svelte';
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

	let filteredClasses = $derived(
		selectedCategory && categories.find(c => c.id === selectedCategory)?.code === 'KREIS'
			? classes.filter(c => c.name === 'Class 10')
			: classes
	);
	let filteredSubjects = $derived(
		selectedCategory && categories.find(c => c.id === selectedCategory)?.code === 'KREIS'
			? subjects.filter(s => s.code && ['KAN', 'ENG', 'HIN', 'MAT', 'SCI', 'SOC'].includes(s.code))
			: subjects
	);

	let students = $state<MarkGridRow[]>([]);
	let maxMarks = $state(100);
	let gridVersion = $state(1);
	let loading = $state(false);
	let saving = $state(false);
	let statusMsg = $state('');
	let statusType = $state<'info' | 'error' | 'success'>('info');
	let tableEl = $state<HTMLDivElement>();
	let table = $state<Tabulator | null>(null);

	onMount(async () => {
		const sp = new URLSearchParams(window.location.search);
		selectedCategory = sp.get('category') ?? '';
		selectedClass = sp.get('class') ?? '';

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

		selectedSubject = sp.get('subject') ?? '';
		selectedAssessment = sp.get('assessment') ?? '';
	});

	let prevAssessment = $state('');
	$effect(() => {
		if (selectedAssessment && selectedAssessment !== prevAssessment) {
			prevAssessment = selectedAssessment;
			loadGrid();
		}
	});

	$effect(() => {
		if (tableEl && !table && students.length > 0) {
			initTable();
		}
	});

	let reqSeq = 0;
	async function loadAssessments() {
		const seq = ++reqSeq;
		const params = new URLSearchParams();
		if (selectedCategory) params.set('category_id', selectedCategory);
		if (selectedClass) params.set('class_id', selectedClass);
		if (selectedSubject) params.set('subject_id', selectedSubject);
		const res = await api<Assessment[]>('GET', '/assessments?' + params.toString());
		if (res.data && seq === reqSeq) assessments = res.data;
	}

	let prevSearch = $state('');
	let prevCat = $state('');
	let prevCls = $state('');
	let prevSub = $state('');
	$effect(() => {
		if (selectedCategory !== prevCat) {
			if (!filteredClasses.find(c => c.id === selectedClass)) selectedClass = '';
			if (!filteredSubjects.find(s => s.id === selectedSubject)) selectedSubject = '';
		}
		if ((selectedCategory !== prevCat || selectedClass !== prevCls || selectedSubject !== prevSub) && prevSearch !== '') {
			selectedAssessment = '';
			prevAssessment = '';
			students = [];
			if (table) { table.destroy(); table = null; }
			statusMsg = '';
		}
		prevCat = selectedCategory;
		prevCls = selectedClass;
		prevSub = selectedSubject;
	});

	$effect(() => {
		const qs = new URLSearchParams();
		if (selectedCategory) qs.set('category', selectedCategory);
		if (selectedClass) qs.set('class', selectedClass);
		if (selectedSubject) qs.set('subject', selectedSubject);
		if (selectedAssessment) qs.set('assessment', selectedAssessment);
		const newSearch = qs.toString() ? '/marks?' + qs.toString() : '/marks';
		if (newSearch !== prevSearch && prevSearch !== '') {
			prevSearch = newSearch;
			history.replaceState(null, '', newSearch);
		}
		if (prevSearch === '') prevSearch = newSearch;
		if (prevSearch !== '') loadAssessments();
	});

	async function loadGrid() {
		if (!selectedAssessment) return;
		loading = true;
		statusMsg = '';
		const res = await api<{ assessment: Assessment; students: MarkGridRow[] }>('GET', '/marks/grid?assessment_id=' + selectedAssessment);
		loading = false;
		if (res.error) {
			statusMsg = 'Error: ' + res.error.message;
			statusType = 'error';
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
			data: students.map((s, i) => ({
				...s,
				_idx: i + 1,
				_marks: s.is_absent ? '' : (s.marks_obtained >= 0 ? s.marks_obtained : ''),
				_absent: s.is_absent
			})),
			layout: 'fitColumns',
			height: 'calc(100vh - 320px)',
			width: '100%',
			selectable: false,
			clipboard: true,
			columns: [
				{ title: '#', field: '_idx', width: 50, hozAlign: 'center' },
				{ title: 'SATS', field: 'sats_number', width: 100 },
				{ title: 'Student Name', field: 'name', minWidth: 180, widthGrow: 3 },
				{ title: 'Father Name', field: 'father_name', minWidth: 140, widthGrow: 2 },
				{
					title: 'Marks / ' + maxMarks,
					field: '_marks',
					editor: 'input',
					editorParams: { elementAttributes: { inputmode: 'decimal' } },
					width: 120,
					hozAlign: 'center',
				cellEdited: (cell: any) => {
					const raw = String(cell.getValue() ?? '');
					const upper = raw.toUpperCase();
					if (upper === 'A' || upper === 'ABS') {
						cell.getRow().update({ _absent: true, _marks: '' });
					} else {
						cell.getRow().update({ _absent: false, marks_obtained: raw === '' ? 0 : parseFloat(raw) || 0 });
					}
					validateCell(cell);
					runValidation();
				},
					formatter: (cell: any) => {
						const data = cell.getRow().getData();
						if (data._absent) return 'ABS';
						const v = cell.getValue();
						return v === '' || v === null || v === undefined ? '' : String(v);
					}
				},
			],
			editorNavigation: (cell: any, direction: string) => {
				if (direction === 'left' || direction === 'right' || direction === 'prev' || direction === 'next') {
					return false;
				}
				return true;
			},
			clipboardPasteParser: (clipboard: string) => {
				const rows = clipboard.split('\n').filter(r => r.trim());
				return rows.map(row => row.split('\t'));
			},
			clipboardPasteAction: 'range',
		});
	}

	let validationErrors = $state<{ name: string; marks: number; max: number }[]>([]);

	function validateCell(cell: any) {
		const data = cell.getRow().getData();
		if (data._absent) {
			cell.getElement().style.backgroundColor = '';
			cell.getElement().style.color = '';
			return;
		}
		const raw = String(cell.getValue() ?? '');
		const val = parseFloat(raw);
		if (raw !== '' && !isNaN(val) && val > maxMarks) {
			cell.getElement().style.backgroundColor = '#fef2f2';
			cell.getElement().style.color = '#dc2626';
		} else if (raw !== '' && !isNaN(val) && val < 0) {
			cell.getElement().style.backgroundColor = '#fef2f2';
			cell.getElement().style.color = '#dc2626';
		} else {
			cell.getElement().style.backgroundColor = '';
			cell.getElement().style.color = '';
		}
	}

	function runValidation(): { name: string; marks: number; max: number }[] {
		if (!table) return [];
		const data = table.getData() as any[];
		const errs: { name: string; marks: number; max: number }[] = [];
		for (const d of data) {
			if (d._absent) continue;
			const raw = String(d._marks ?? '');
			if (raw === '') continue;
			const val = parseFloat(raw);
			if (isNaN(val)) continue;
			if (val > maxMarks || val < 0) {
				errs.push({ name: d.name, marks: val, max: maxMarks });
			}
		}
		validationErrors = errs;
		return errs;
	}

	async function saveMarks() {
		if (!table || !selectedAssessment) return;

		const errs = runValidation();
		if (errs.length > 0) {
			statusMsg = `Cannot save: ${errs.length} student${errs.length > 1 ? 's have' : ' has'} invalid marks (exceed max ${maxMarks} or negative). Highlighted in red.`;
			statusType = 'error';
			return;
		}

		saving = true;
		statusMsg = '';

		const data = table.getData() as any[];
		const marks: MarkInput[] = data.map((d: any) => ({
			student_id: d.student_id,
			marks_obtained: d._absent ? 0 : (parseFloat(d._marks) || 0),
			is_absent: !!d._absent,
			remarks: ''
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
			statusType = 'error';
		} else if (res.data) {
			const d = res.data as any;
			statusMsg = `Saved ${d.updated} mark${d.updated !== 1 ? 's' : ''}` + (d.errors?.length ? ` (${d.errors.length} error${d.errors.length !== 1 ? 's' : ''})` : '');
			statusType = d.errors?.length ? 'error' : 'success';
			gridVersion = (gridVersion || 0) + 1;
		}
	}

	function resetForm() {
		selectedAssessment = '';
		students = [];
		if (table) { table.destroy(); table = null; }
		statusMsg = '';
	}

	const statusStyles = {
		info: 'bg-blue-50 text-blue-700 border-blue-200',
		success: 'bg-green-50 text-green-700 border-green-200',
		error: 'bg-red-50 text-red-700 border-red-200',
	};
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-3">
			<div class="w-9 h-9 rounded-lg bg-primary-100 flex items-center justify-center">
				<Table size={18} class="text-primary-600" />
			</div>
			<div>
				<h1 class="text-xl font-bold text-slate-900">Marks Entry</h1>
				<p class="text-xs text-slate-500">Enter and manage student marks for assessments</p>
			</div>
		</div>
		<div class="flex gap-2">
			<Button onclick={saveMarks} disabled={saving || !selectedAssessment || !table} loading={saving} icon={Save}>
				Save Marks
			</Button>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-4 shadow-sm">
		<div class="flex items-center gap-2 mb-3">
			<Filter size={14} class="text-slate-400" />
			<span class="text-xs font-semibold text-slate-600 uppercase tracking-wider">Filters</span>
		</div>
		<div class="flex flex-wrap gap-3 items-end">
			<div class="w-44">
				<Select bind:value={selectedCategory} options={categories} label="Category" icon={ClipboardCheck} placeholder="All categories" />
			</div>
			<div class="w-44">
				<Select bind:value={selectedClass} options={filteredClasses} label="Class" icon={Users} placeholder="All classes" />
			</div>
			<div class="w-44">
				<Select bind:value={selectedSubject} options={filteredSubjects} label="Subject" icon={BookOpen} placeholder="All subjects" />
			</div>
			<div class="w-48">
				<Select bind:value={selectedAssessment} options={assessments} label="Assessment" icon={ClipboardCheck} placeholder="Select assessment" />
			</div>
			<Button onclick={resetForm} variant="ghost" class="px-3 py-2 text-sm rounded-lg border border-slate-300 hover:border-slate-400 bg-white text-slate-600 hover:text-slate-900 hover:bg-slate-50">Clear</Button>
		</div>
	</div>

	{#if statusMsg}
		<div class="flex items-center gap-2 text-sm px-4 py-2.5 rounded-lg border {statusStyles[statusType]}">
			<span>{statusMsg}</span>
		</div>
	{/if}

	{#if validationErrors.length > 0}
		<div class="bg-red-50 border border-red-200 rounded-lg p-4">
			<h3 class="text-sm font-semibold text-red-700 mb-2">Validation Errors ({validationErrors.length})</h3>
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-1.5">
				{#each validationErrors as e}
					<div class="text-xs text-red-600">
						<span class="font-medium">{e.name}</span>: {e.marks} &gt; max {e.max}
					</div>
				{/each}
			</div>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 shadow-sm w-full">
		{#if selectedAssessment}
			<div class="marks-grid-container">
				{#if loading}
					<div class="p-12 text-center text-sm text-slate-400">Loading marks grid...</div>
				{:else}
					<div bind:this={tableEl} class="marks-grid"></div>
				{/if}
			</div>
			<div class="px-4 py-2 bg-slate-50 border-t border-slate-100 text-xs text-slate-400 flex gap-4 items-center">
				<span>Type <kbd class="px-1 py-0.5 rounded bg-slate-200 text-slate-600 font-mono text-[10px]">a</kbd> to mark absent</span>
				<span class="w-px h-3 bg-slate-300"></span>
				<span>&uarr; &darr; Navigate rows</span>
				<span class="ml-auto">{students.length} student{students.length !== 1 ? 's' : ''}</span>
			</div>
		{:else}
			<div class="p-12 text-center">
				<Table size={32} class="mx-auto text-slate-300 mb-3" />
				<p class="text-sm text-slate-400">Select an assessment from the filters above to enter marks.</p>
			</div>
		{/if}
	</div>
</div>

<style>
	:global(.marks-grid-container) {
		width: 100%;
	}
	:global(.marks-grid-container .tabulator) {
		border: none;
		border-radius: 0;
	}
	:global(.marks-grid-container .tabulator .tabulator-header) {
		background: #f8fafc;
		border-bottom: 1px solid #e2e8f0;
	}
	:global(.marks-grid-container .tabulator .tabulator-header .tabulator-col) {
		background: #f8fafc;
	}
	:global(.marks-grid-container .tabulator .tabulator-row) {
		border-bottom: 1px solid #f1f5f9;
	}
	:global(.marks-grid-container .tabulator .tabulator-row .tabulator-cell) {
		padding: 8px 10px;
		border-right: 1px solid #f1f5f9;
	}
	:global(.marks-grid-container .tabulator .tabulator-row.tabulator-row-even) {
		background: #fafafa;
	}
	:global(.marks-grid-container .tabulator .tabulator-cell.tabulator-editing) {
		border: 2px solid #2563eb !important;
	}
	:global(.marks-grid-container .tabulator .tabulator-tableholder) {
		overflow: auto !important;
	}
</style>
