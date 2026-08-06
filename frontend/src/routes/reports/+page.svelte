<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { GraduationCap, Printer, Users, FileText } from 'lucide-svelte';
	import Select from '$lib/components/Select.svelte';
	import Button from '$lib/components/Button.svelte';
	import type { Class, AcademicYear } from '$lib/types';
	import { onMount } from 'svelte';

	type Tab = 'marksheet' | 'report';
	let activeTab = $state<Tab>('marksheet');

	let classes = $state<Class[]>([]);
	let years = $state<AcademicYear[]>([]);

	let selectedClass = $state('');
	let selectedYear = $state('');
	let selectedStudent = $state('');

	let loading = $state(false);
	let err = $state('');

	type MarkCell = { assessment_id: string; value: number; is_absent: boolean; has_mark: boolean };
	type SubjectAgg = { subject_id: string; subject_code: string; subject_name: string; subject_type: string; total: number; max_total: number; percentage: number; grade: string; grade_label?: string };
	type MarkSheetStudent = { student_id: string; sats_number: string; name: string; roll_no: number; marks: MarkCell[]; total: number; max_total: number; percentage: number; grade: string; rank: number; subjects: SubjectAgg[] };
	type MarkSheetAssessment = { id: string; name: string; subject_id: string; subject_code: string; subject_name: string; category_id: string; category_name: string; category_code: string; max_marks: number; date?: string; term: string; subject_type: string };
	type SubjectGroup = { subject_id: string; subject_code: string; subject_name: string; subject_type?: string; assessments: MarkSheetAssessment[] };
	type TermGroup = { term: string; subjects: SubjectGroup[] };

	let markSheetData = $state<{
		class_id: string; class_name: string; academic_year: string; term?: string;
		subjects: SubjectGroup[]; terms: TermGroup[]; assessments: MarkSheetAssessment[]; students: MarkSheetStudent[];
	} | null>(null);

	type ReportAssessment = { id: string; name: string; category: string; term: string; max: number; value: number; absent: boolean; has_mark: boolean };
	type ReportSubject = { subject_id: string; subject_code: string; subject_name: string; subject_type: string; assessments: ReportAssessment[]; total: number; max_max: number; percentage: number; grade: string; grade_label?: string };
	type ReportStudent = { name: string; class: string; section?: string; roll_no: number; sats_number: string; gender?: string; date_of_birth?: string };
	type StudentReport = { student: ReportStudent; academic_year: string; term?: string; subjects: ReportSubject[]; grand_total: number; grand_max: number; percentage: number; grade: string; attendance?: { present: number; total: number; percentage: number }; remarks?: string };

	let studentReport = $state<StudentReport | null>(null);
	let reportStudents = $state<{ id: string; name: string; roll_no: number }[]>([]);

	const termOptions = [
		{ id: '', name: 'All Terms' },
		{ id: 'Term 1', name: 'Term 1 (FA1, FA2, SA1)' },
		{ id: 'Term 2', name: 'Term 2 (FA3, FA4, SA2)' },
	];

	let selectedTerm = $state('');

	let studentOptions = $derived(
		reportStudents.map(s => ({ id: s.id, name: `${s.name} (Roll ${s.roll_no})` }))
	);

	function pctClass(pct: number): string {
		if (pct >= 90) return 'text-emerald-700 bg-emerald-50';
		if (pct >= 70) return 'text-green-700 bg-green-50';
		if (pct >= 50) return 'text-blue-700 bg-blue-50';
		if (pct >= 40) return 'text-sky-700 bg-sky-50';
		if (pct >= 30) return 'text-amber-700 bg-amber-50';
		return 'text-red-700 bg-red-50';
	}

	function gradeClass(grade: string): string {
		if (grade === 'A+') return 'text-emerald-700 bg-emerald-50';
		if (grade === 'A') return 'text-green-700 bg-green-50';
		if (grade === 'B+') return 'text-blue-700 bg-blue-50';
		if (grade === 'B') return 'text-sky-700 bg-sky-50';
		if (grade === 'C+') return 'text-amber-700 bg-amber-50';
		if (grade === 'C') return 'text-orange-700 bg-orange-50';
		if (grade === 'D') return 'text-red-600 bg-red-50';
		if (grade === 'F') return 'text-red-800 bg-red-100';
		return 'text-slate-600 bg-slate-100';
	}

	function subjectTypeLabel(t: string): string {
		return t === 'curricular' ? 'Curricular' : 'Co-curricular';
	}

	const curricularOrder = ['KAN', 'ENG', 'HIN', 'MAT', 'SCI', 'SOC'];

	onMount(async () => {
		const [cr, yr] = await Promise.all([
			api<Class[]>('GET', '/classes'),
			api<AcademicYear[]>('GET', '/academic-years'),
		]);
		if (cr.data) classes = cr.data;
		if (yr.data) {
			years = yr.data;
			const cur = yr.data.find(y => y.is_current);
			if (cur) selectedYear = cur.id;
		}
	});

	async function loadMarkSheet() {
		if (!selectedClass) { err = 'Select a class'; return; }
		err = ''; loading = true; markSheetData = null;
		try {
			const params = new URLSearchParams({ class_id: selectedClass });
			if (selectedYear) params.set('academic_year_id', selectedYear);
			if (selectedTerm) params.set('term', selectedTerm);
			const res = await api<typeof markSheetData>('GET', `/reports/mark-sheet?${params}`);
			if (res.data) markSheetData = res.data;
			else if (res.error) err = res.error.message;
		} catch (e) {
			err = 'Failed to load mark sheet';
		} finally {
			loading = false;
		}
	}

	async function loadReportStudents() {
		if (!selectedClass) return;
		const params = new URLSearchParams({ class_id: selectedClass });
		if (selectedYear) params.set('academic_year_id', selectedYear);
		const res = await api<{ id: string; name: string; roll_no: number; sats_number: string }[]>('GET', `/students?${params}`);
		reportStudents = res.data || [];
	}

	async function loadStudentReport() {
		if (!selectedStudent) { err = 'Select a student'; return; }
		err = ''; loading = true; studentReport = null;
		try {
			const params = new URLSearchParams({ student_id: selectedStudent });
			if (selectedYear) params.set('academic_year_id', selectedYear);
			if (selectedTerm) params.set('term', selectedTerm);
			const res = await api<StudentReport>('GET', `/reports/student?${params}`);
			if (res.data) studentReport = res.data;
			else if (res.error) err = res.error.message;
		} catch (e) {
			err = 'Failed to load report';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (selectedClass && activeTab === 'marksheet') loadMarkSheet();
	});

	$effect(() => {
		if (selectedTerm && activeTab === 'marksheet' && selectedClass) loadMarkSheet();
	});

	$effect(() => {
		if (selectedClass && activeTab === 'report') loadReportStudents();
	});

	function handlePrint() {
		window.print();
	}
</script>

<svelte:head>
	<title>Reports — Pragati School Management</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<div class="flex items-center justify-between no-print">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Reports</h1>
			<p class="text-sm text-slate-500 mt-0.5">Class mark sheets &amp; student report cards</p>
		</div>
		{#if (activeTab === 'marksheet' && markSheetData) || (activeTab === 'report' && studentReport)}
			<Button icon={Printer} onclick={handlePrint}>Print</Button>
		{/if}
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
		<div class="flex flex-wrap gap-3 items-end">
			<div class="w-44">
				<Select bind:value={selectedClass} options={classes.map(c => ({ id: c.id, name: c.name }))} placeholder="Select class" />
			</div>
			<div class="w-44">
				<Select bind:value={selectedYear} options={years.map(y => ({ id: y.id, name: y.name }))} placeholder="Academic year" />
			</div>
			<div class="w-48">
				<Select bind:value={selectedTerm} options={termOptions} placeholder="All Terms" />
			</div>

			<div class="flex bg-slate-100 rounded-lg p-1 ml-2">
				<button onclick={() => activeTab = 'marksheet'}
					class="px-4 py-1.5 text-sm font-medium rounded-md transition-colors {activeTab === 'marksheet' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
					<span class="flex items-center gap-1.5"><Users size={14} /> Mark Sheet</span>
				</button>
				<button onclick={() => activeTab = 'report'}
					class="px-4 py-1.5 text-sm font-medium rounded-md transition-colors {activeTab === 'report' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
					<span class="flex items-center gap-1.5"><FileText size={14} /> Report Card</span>
				</button>
			</div>

			{#if activeTab === 'report'}
				<div class="w-56">
					<Select bind:value={selectedStudent} options={studentOptions} placeholder="Select student" />
				</div>
				<Button onclick={loadStudentReport} disabled={!selectedStudent || loading}>
					{loading ? 'Loading...' : 'Generate'}
				</Button>
			{/if}
		</div>
		{#if err}
			<p class="text-sm text-danger-600 mt-2">{err}</p>
		{/if}
	</div>

	{#if loading}
		<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-slate-400 text-sm">Loading...</div>
	{:else if activeTab === 'marksheet' && markSheetData}
		{@const ms = markSheetData}
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden print-area">
			<div class="px-6 py-4 border-b border-slate-200 print-header">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-3">
						<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center">
							<GraduationCap size={20} class="text-white" />
						</div>
						<div>
							<h2 class="text-base font-bold text-slate-900">{ms.class_name} — Mark Sheet</h2>
							<p class="text-xs text-slate-500">Academic Year: {years.find(y => y.id === ms.academic_year)?.name || 'All'}</p>
						</div>
					</div>
					<div class="text-xs text-slate-400 hidden print-only">{new Date().toLocaleDateString()}</div>
				</div>
			</div>

			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200">
							<th rowspan="2" class="px-3 py-2 text-left font-semibold text-slate-600 border-r border-slate-200 sticky left-0 bg-slate-50 z-10 w-8">#</th>
							<th rowspan="2" class="px-3 py-2 text-left font-semibold text-slate-600 border-r border-slate-200 sticky left-8 bg-slate-50 z-10">Student</th>
							{#each ms.subjects as sg}
								<th colspan={sg.assessments.length} class="px-3 py-1.5 text-center font-semibold text-slate-700 border-r border-slate-200 border-b border-slate-100">
									<div class="text-xs">{sg.subject_code}</div>
									<div class="text-[10px] font-normal text-slate-400">{sg.subject_type === 'curricular' ? 'Curricular' : 'Co-curricular'}</div>
								</th>
							{/each}
							<th rowspan="2" class="px-3 py-2 text-center font-semibold text-slate-700 border-r border-slate-200 w-16">Total</th>
							<th rowspan="2" class="px-3 py-2 text-center font-semibold text-slate-700 border-r border-slate-200 w-14">%</th>
							<th rowspan="2" class="px-3 py-2 text-center font-semibold text-slate-700 border-r border-slate-200 w-14">Grade</th>
							<th rowspan="2" class="px-3 py-2 text-center font-semibold text-slate-700 w-12">Rank</th>
						</tr>
						<tr class="bg-slate-50 border-b border-slate-200">
							{#each ms.subjects as sg}
								{#each sg.assessments as a}
									<th class="px-2 py-1.5 text-center font-medium text-slate-500 border-r border-slate-200 text-[10px]">
										<div>{a.name}</div>
										<div class="text-[9px] text-slate-400">/{a.max_marks}</div>
									</th>
								{/each}
							{/each}
						</tr>
					</thead>
					<tbody>
						{#each ms.students as s, i}
							<tr class="border-b border-slate-100 hover:bg-slate-50/50">
								<td class="px-3 py-2 text-center text-slate-400 border-r border-slate-100 sticky left-0 bg-white z-10">{i + 1}</td>
								<td class="px-3 py-2 border-r border-slate-100 sticky left-8 bg-white z-10">
									<div class="font-medium text-slate-800 text-xs whitespace-nowrap">{s.name}</div>
									<div class="text-[10px] text-slate-400">Roll {s.roll_no}</div>
								</td>
								{#each s.marks as m}
									<td class="px-2 py-2 text-center border-r border-slate-100">
										{#if m.has_mark}
											<span class="text-xs font-medium {m.is_absent ? 'text-red-500' : 'text-slate-700'}">
												{m.is_absent ? 'AB' : m.value}
											</span>
										{:else}
											<span class="text-xs text-slate-300">—</span>
										{/if}
									</td>
								{/each}
								<td class="px-3 py-2 text-center font-semibold text-slate-800 border-r border-slate-100">{s.total}</td>
								<td class="px-3 py-2 text-center border-r border-slate-100">
									<span class="text-xs font-medium px-1.5 py-0.5 rounded {pctClass(s.percentage)}">{s.percentage.toFixed(1)}</span>
								</td>
								<td class="px-3 py-2 text-center border-r border-slate-100">
									<span class="text-xs font-bold px-2 py-0.5 rounded {gradeClass(s.grade)}">{s.grade}</span>
								</td>
								<td class="px-3 py-2 text-center">
									<span class="text-xs font-semibold text-slate-600">{s.rank}</span>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<div class="px-6 py-4 border-t border-slate-200 space-y-4">
			{#if ms.students[0]?.subjects}
			{@const curSubs = ms.students[0].subjects.filter(s => s.subject_type === 'curricular')}
			{@const coSubs = ms.students[0].subjects.filter(s => s.subject_type !== 'curricular')}
			{#if curSubs.length > 0}
				<div>
					<div class="text-[10px] font-semibold text-slate-400 uppercase tracking-wider mb-2">Curricular</div>
					<div class="flex flex-wrap gap-3">
						{#each curSubs as sg}
							<div class="border border-slate-200 rounded-lg p-3 min-w-40">
								<div class="flex items-center justify-between mb-2">
									<span class="text-xs font-semibold text-slate-700">{sg.subject_name}</span>
									<span class="text-[10px] text-slate-400">{sg.subject_code}</span>
								</div>
								<div class="text-lg font-bold text-slate-800">{sg.total}<span class="text-xs text-slate-400 font-normal">/{sg.max_total}</span></div>
								<div class="flex items-center justify-between mt-1">
									<span class="text-xs text-slate-500">{sg.percentage.toFixed(1)}%</span>
									<span class="text-xs font-bold px-1.5 py-0.5 rounded {gradeClass(sg.grade)}">{sg.grade}</span>
								</div>
							</div>
						{/each}
					</div>
				</div>
				{/if}
				{#if coSubs.length > 0}
				<div>
					<div class="text-[10px] font-semibold text-slate-400 uppercase tracking-wider mb-2">Co-curricular</div>
					<div class="flex flex-wrap gap-3">
						{#each coSubs as sg}
							<div class="border border-slate-200 rounded-lg p-3 min-w-44">
								<div class="flex items-center justify-between mb-2">
									<span class="text-xs font-semibold text-slate-700">{sg.subject_name}</span>
									<span class="text-[10px] text-slate-400">{sg.subject_code}</span>
								</div>
								<div class="text-lg font-bold text-slate-800">{sg.total}<span class="text-xs text-slate-400 font-normal">/{sg.max_total}</span></div>
								<div class="flex items-center justify-between mt-1">
									<span class="text-xs text-slate-500">{sg.percentage.toFixed(1)}%</span>
									<div class="flex flex-col items-end gap-0.5">
										<span class="text-xs font-bold px-1.5 py-0.5 rounded {gradeClass(sg.grade)}">{sg.grade}</span>
										{#if sg.grade_label}<span class="text-[10px] text-slate-500">{sg.grade_label}</span>{/if}
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/if}
			{/if}
		</div>
	</div>

	{:else if activeTab === 'report' && studentReport}
		{@const r = studentReport}
		<div class="bg-white rounded-xl border border-slate-200 print-area max-w-4xl mx-auto">
			<div class="px-8 py-6 border-b-2 border-primary-600">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-4">
						<div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shadow-md">
							<GraduationCap size={28} class="text-white" />
						</div>
						<div>
							<h1 class="text-xl font-bold text-slate-900 font-kannada tracking-wide">ಪ್ರಗತಿ</h1>
							<p class="text-xs text-slate-500">Morarji Desai Residential School, Bahaddurghatta</p>
						</div>
					</div>
					<div class="text-right">
						<h2 class="text-lg font-bold text-primary-700">Student Report Card</h2>
						<p class="text-xs text-slate-500">Year: {years.find(y => y.id === r.academic_year)?.name || '—'} {#if r.term} &middot; {r.term}{/if}</p>
					</div>
				</div>
			</div>

			<div class="px-8 py-4 bg-slate-50 border-b border-slate-200">
				<div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-sm">
					<div><span class="text-slate-500 text-xs">Name</span><p class="font-medium text-slate-800">{r.student.name}</p></div>
					<div><span class="text-slate-500 text-xs">Class & Section</span><p class="font-medium text-slate-800">{r.student.class} {r.student.section || ''}</p></div>
					<div><span class="text-slate-500 text-xs">Roll No</span><p class="font-medium text-slate-800">{r.student.roll_no}</p></div>
					<div><span class="text-slate-500 text-xs">SATS No</span><p class="font-medium text-slate-800">{r.student.sats_number}</p></div>
					<div><span class="text-slate-500 text-xs">Gender</span><p class="font-medium text-slate-800 capitalize">{r.student.gender || '—'}</p></div>
					<div><span class="text-slate-500 text-xs">Date of Birth</span><p class="font-medium text-slate-800">{r.student.date_of_birth || '—'}</p></div>
					{#if r.attendance}
						<div><span class="text-slate-500 text-xs">Attendance</span><p class="font-medium text-slate-800">{r.attendance.present}/{r.attendance.total} ({r.attendance.percentage.toFixed(1)}%)</p></div>
					{/if}
				</div>
			</div>

		<div class="px-8 py-6">
			{#if r.subjects}
			{@const curSubs = r.subjects.filter(s => s.subject_type === 'curricular')}
			{@const coSubs = r.subjects.filter(s => s.subject_type !== 'curricular')}

			{#if curSubs.length > 0}
			<h3 class="text-sm font-bold text-slate-700 uppercase tracking-wider mb-3">Scholastic Areas</h3>
			<table class="w-full text-sm mb-6">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200">
						<th class="px-3 py-2 text-left font-semibold text-slate-600">Subject</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-20">Total</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-20">Max</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-20">%</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-20">Grade</th>
					</tr>
				</thead>
				<tbody>
					{#each curSubs as sub}
						<tr class="border-b border-slate-100">
							<td class="px-3 py-2.5">
								<div class="font-medium text-slate-800">{sub.subject_name}</div>
								<div class="flex flex-wrap gap-1 mt-1">
									{#each sub.assessments as a}
										<span class="text-[10px] px-1.5 py-0.5 rounded {a.term === 'Term 1' ? 'bg-blue-50 text-blue-600' : 'bg-purple-50 text-purple-600'}" title="{a.name}">
											{a.name}: {#if a.has_mark}{a.absent ? 'AB' : a.value}{:else}—{/if}/{a.max}
										</span>
									{/each}
								</div>
							</td>
							<td class="px-3 py-2.5 text-center font-medium text-slate-800">{sub.total}</td>
							<td class="px-3 py-2.5 text-center text-slate-500">{sub.max_max}</td>
							<td class="px-3 py-2.5 text-center">
								<span class="text-xs font-medium px-1.5 py-0.5 rounded {pctClass(sub.percentage)}">{sub.percentage.toFixed(1)}</span>
							</td>
							<td class="px-3 py-2.5 text-center">
								<span class="text-xs font-bold px-2 py-0.5 rounded {gradeClass(sub.grade)}">{sub.grade}</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			{/if}

			{#if coSubs.length > 0}
			<h3 class="text-sm font-bold text-slate-700 uppercase tracking-wider mb-3">Co-Scholastic Areas</h3>
			<table class="w-full text-sm mb-6">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200">
						<th class="px-3 py-2 text-left font-semibold text-slate-600">Subject</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-20">Total</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-20">Max</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-20">%</th>
						<th class="px-3 py-2 text-center font-semibold text-slate-600 w-24">Grade</th>
					</tr>
				</thead>
				<tbody>
					{#each coSubs as sub}
						<tr class="border-b border-slate-100">
							<td class="px-3 py-2.5">
								<div class="font-medium text-slate-800">{sub.subject_name}</div>
								<div class="flex flex-wrap gap-1 mt-1">
									{#each sub.assessments as a}
										<span class="text-[10px] px-1.5 py-0.5 rounded {a.term === 'Term 1' ? 'bg-blue-50 text-blue-600' : 'bg-purple-50 text-purple-600'}" title="{a.name}">
											{a.name}: {#if a.has_mark}{a.absent ? 'AB' : a.value}{:else}—{/if}/{a.max}
										</span>
									{/each}
								</div>
							</td>
							<td class="px-3 py-2.5 text-center font-medium text-slate-800">{sub.total}</td>
							<td class="px-3 py-2.5 text-center text-slate-500">{sub.max_max}</td>
							<td class="px-3 py-2.5 text-center">
								<span class="text-xs font-medium px-1.5 py-0.5 rounded {pctClass(sub.percentage)}">{sub.percentage.toFixed(1)}</span>
							</td>
							<td class="px-3 py-2.5 text-center">
								<div class="flex flex-col items-center gap-0.5">
									<span class="text-xs font-bold px-2 py-0.5 rounded {gradeClass(sub.grade)}">{sub.grade}</span>
									{#if sub.grade_label}
										<span class="text-[10px] text-slate-500">{sub.grade_label}</span>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			{/if}

			{#if r.remarks}
			<div class="mt-2 border border-slate-200 rounded-lg p-4">
				<h4 class="text-xs font-bold text-slate-600 uppercase tracking-wider mb-1">Teacher Remarks</h4>
				<p class="text-sm text-slate-700">{r.remarks}</p>
			</div>
			{/if}
			{/if}
		</div>

		<div class="px-8 py-4 border-t border-slate-200 flex justify-between items-end">
			<div class="text-xs text-slate-400">Generated on {new Date().toLocaleDateString()} • Pragati v1.0</div>
			<div class="text-right">
				<div class="border-t border-slate-400 pt-1 text-xs text-slate-500 w-40 text-center">Class Teacher</div>
			</div>
		</div>
	</div>

	{:else if !loading}
		<div class="bg-white rounded-xl border border-slate-200 p-12 text-center">
			<FileText size={40} class="mx-auto text-slate-300 mb-3" />
			<p class="text-slate-500 text-sm">
				{#if activeTab === 'marksheet'}
					Select a class and academic year to view the mark sheet.
				{:else}
					Select a class and student to generate a report card.
				{/if}
			</p>
		</div>
	{/if}
</div>

<style>
	@media print {
		:global(body) { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
		:global(.no-print) { display: none !important; }
		.print-area { border: none !important; border-radius: 0 !important; box-shadow: none !important; max-width: 100% !important; margin: 0 !important; }
	}
</style>
