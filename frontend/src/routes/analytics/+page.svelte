<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';
	import { BarChart3, Users, BookOpen, ClipboardCheck, TrendingUp, Award, AlertCircle } from 'lucide-svelte';
	import Select from '$lib/components/Select.svelte';
	import SearchFilter from '$lib/components/SearchFilter.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import type { Class, Assessment } from '$lib/types';

	type Tab = 'overview' | 'performance' | 'progress';
	let activeTab = $state<Tab>('overview');

	let classes = $state<Class[]>([]);
	let selectedClass = $state('');

	// Overview data
	let overview = $state<{
		total_students: number; total_teachers: number; total_classes: number; total_assessments: number;
		students_by_class: { class: string; count: number }[];
	} | null>(null);

	// Performance data (derived from mark sheet)
	let perfData = $state<{
		students: { name: string; roll_no: number; percentage: number; grade: string }[];
		subjects: { subject_code: string; subject_name: string; average: number; grade: string }[];
		grade_dist: { grade: string; count: number; cls: string }[];
		class_name: string;
	} | null>(null);

	// Progress data
	let progressAssessments = $state<(Assessment & { progress: number })[]>([]);
	let progressSearch = $state('');
	let progressPage = $state(1);
	const progressPageSize = 15;

	let filteredProgress = $derived(() => {
		if (!progressSearch.trim()) return progressAssessments;
		const q = progressSearch.toLowerCase();
		return progressAssessments.filter(a => a.name?.toLowerCase().includes(q) || a.subject_name?.toLowerCase().includes(q));
	});

	let paginatedProgress = $derived(() => {
		const filtered = filteredProgress();
		const start = (progressPage - 1) * progressPageSize;
		return filtered.slice(start, start + progressPageSize);
	});

	let totalProgress = $derived(filteredProgress().length);

	function onProgressPageChange(p: number) { progressPage = p; }
	function resetProgressPage() { progressPage = 1; }

	let loading = $state(true);
	let perfLoading = $state(false);

	const gradeOrder = ['A+', 'A', 'B+', 'B', 'C+', 'C', 'D', 'F'];
	const gradeColors: Record<string, string> = {
		'A+': '#10b981', 'A': '#22c55e', 'B+': '#3b82f6', 'B': '#0ea5e9',
		'C+': '#f59e0b', 'C': '#f97316', 'D': '#ef4444', 'F': '#dc2626',
	};

	const tabItems: { id: Tab; label: string; icon: typeof BarChart3 }[] = [
		{ id: 'overview', label: 'Overview', icon: BarChart3 },
		{ id: 'performance', label: 'Class Performance', icon: TrendingUp },
		{ id: 'progress', label: 'Marks Progress', icon: ClipboardCheck },
	];

	async function loadOverview() {
		const res = await api<typeof overview>('GET', '/dashboard/stats');
		if (res.data) overview = res.data;
	}

	async function loadPerformance() {
		if (!selectedClass) return;
		perfLoading = true;
		const res = await api<{
			class_name: string;
			students: { name: string; roll_no: number; percentage: number; grade: string; subjects: { subject_code: string; subject_name: string; percentage: number; grade: string }[] }[];
		}>('GET', `/reports/mark-sheet?class_id=${selectedClass}`);
		perfLoading = false;
		if (!res.data) return;

		const d = res.data;
		const subjectsMap: Record<string, { code: string; name: string; total: number; count: number }> = {};
		const gradeCounts: Record<string, number> = {};

		for (const s of d.students) {
			gradeCounts[s.grade] = (gradeCounts[s.grade] || 0) + 1;
			for (const sub of s.subjects) {
				if (!subjectsMap[sub.subject_code]) {
					subjectsMap[sub.subject_code] = { code: sub.subject_code, name: sub.subject_name, total: 0, count: 0 };
				}
				subjectsMap[sub.subject_code].total += sub.percentage;
				subjectsMap[sub.subject_code].count++;
			}
		}

		const subjects = Object.values(subjectsMap).map(sm => {
			const avg = sm.count > 0 ? sm.total / sm.count : 0;
			return { subject_code: sm.code, subject_name: sm.name, average: avg, grade: gradeFromAvg(avg) };
		});

		const grade_dist = gradeOrder.map(g => ({
			grade: g, count: gradeCounts[g] || 0, cls: gradeColors[g],
		}));

		perfData = {
			students: d.students.map(s => ({ name: s.name, roll_no: s.roll_no, percentage: s.percentage, grade: s.grade })),
			subjects,
			grade_dist,
			class_name: d.class_name,
		};
	}

	function gradeFromAvg(avg: number): string {
		if (avg >= 90) return 'A+';
		if (avg >= 70) return 'A';
		if (avg >= 50) return 'B+';
		if (avg >= 40) return 'B';
		if (avg >= 30) return 'C+';
		if (avg < 30) return 'C';
		return 'C';
	}

	async function loadProgress() {
		const params = new URLSearchParams({ limit: '100' });
		if (selectedClass) params.set('class_id', selectedClass);
		const res = await api<Assessment[]>('GET', `/assessments?${params}`);
		if (res.data) {
			progressAssessments = res.data.map(a => ({
				...a,
				progress: a.student_count ? Math.round(((a.marks_count ?? 0) / a.student_count) * 100) : 0,
			}));
		}
	}

	onMount(async () => {
		const [cr] = await Promise.all([api<Class[]>('GET', '/classes')]);
		if (cr.data) classes = cr.data;
		await loadOverview();
		loading = false;
	});

	$effect(() => {
		if (selectedClass && activeTab === 'performance') loadPerformance();
		if (selectedClass && activeTab === 'progress') loadProgress();
	});

	$effect(() => {
		if (activeTab === 'progress' && progressAssessments.length === 0) loadProgress();
	});

	const maxStudents = $derived(overview ? Math.max(...overview.students_by_class.map(s => s.count), 1) : 1);
	const topStudents = $derived(perfData ? [...perfData.students].sort((a, b) => b.percentage - a.percentage).slice(0, 10) : []);
	const maxGrade = $derived(perfData ? Math.max(...perfData.grade_dist.map(g => g.count), 1) : 1);
</script>

<svelte:head>
	<title>Analytics — Pragati School Management</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<div class="flex items-center gap-3">
		<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-700 flex items-center justify-center">
			<BarChart3 size={20} class="text-white" />
		</div>
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Analytics</h1>
			<p class="text-sm text-slate-500">School performance insights</p>
		</div>
	</div>

	<div class="flex gap-1 bg-slate-100 rounded-lg p-1 no-print">
		{#each tabItems as item}
			<button onclick={() => activeTab = item.id}
				class="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab === item.id ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
				<item.icon size={16} />
				{item.label}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Loading...</div>

	{:else if activeTab === 'overview' && overview}
		<!-- Overview -->
		<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<div class="flex items-center gap-3 mb-3">
					<div class="w-9 h-9 rounded-lg bg-blue-50 flex items-center justify-center"><Users size={18} class="text-blue-600" /></div>
					<span class="text-xs text-slate-500 font-medium">Students</span>
				</div>
				<div class="text-2xl font-bold text-slate-900">{overview.total_students}</div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<div class="flex items-center gap-3 mb-3">
					<div class="w-9 h-9 rounded-lg bg-emerald-50 flex items-center justify-center"><Users size={18} class="text-emerald-600" /></div>
					<span class="text-xs text-slate-500 font-medium">Teachers</span>
				</div>
				<div class="text-2xl font-bold text-slate-900">{overview.total_teachers}</div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<div class="flex items-center gap-3 mb-3">
					<div class="w-9 h-9 rounded-lg bg-purple-50 flex items-center justify-center"><BookOpen size={18} class="text-purple-600" /></div>
					<span class="text-xs text-slate-500 font-medium">Classes</span>
				</div>
				<div class="text-2xl font-bold text-slate-900">{overview.total_classes}</div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<div class="flex items-center gap-3 mb-3">
					<div class="w-9 h-9 rounded-lg bg-amber-50 flex items-center justify-center"><ClipboardCheck size={18} class="text-amber-600" /></div>
					<span class="text-xs text-slate-500 font-medium">Assessments</span>
				</div>
				<div class="text-2xl font-bold text-slate-900">{overview.total_assessments}</div>
			</div>
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-6">
			<h2 class="text-base font-semibold text-slate-900 mb-4">Students by Class</h2>
			<div class="space-y-3">
				{#each overview.students_by_class as sc}
					<div class="flex items-center gap-3">
						<span class="text-sm text-slate-600 w-24 shrink-0">{sc.class}</span>
						<div class="flex-1 h-6 bg-slate-100 rounded-full overflow-hidden">
							<div class="h-full bg-gradient-to-r from-primary-500 to-primary-600 rounded-full transition-all" style="width: {(sc.count / maxStudents) * 100}%"></div>
						</div>
						<span class="text-sm font-medium text-slate-700 w-8 text-right">{sc.count}</span>
					</div>
				{/each}
			</div>
		</div>

	{:else if activeTab === 'performance'}
		<!-- Class Performance -->
		<div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
			<div class="w-56">
				<Select bind:value={selectedClass} options={[{ id: '', name: 'Select a class' }, ...classes.map(c => ({ id: c.id, name: c.name }))]} placeholder="Select a class" />
			</div>
		</div>

		{#if !selectedClass}
			<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Select a class to view performance analytics</div>
		{:else if perfLoading}
			<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Loading...</div>
		{:else if perfData}
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<!-- Grade Distribution -->
				<div class="bg-white rounded-xl border border-slate-200 p-6">
					<h2 class="text-base font-semibold text-slate-900 mb-4">Grade Distribution — {perfData.class_name}</h2>
					<div class="space-y-2">
						{#each perfData.grade_dist as gd}
							{#if gd.count > 0}
								<div class="flex items-center gap-3">
									<span class="text-xs font-bold w-8" style="color: {gd.cls}">{gd.grade}</span>
									<div class="flex-1 h-5 bg-slate-100 rounded overflow-hidden">
										<div class="h-full rounded transition-all" style="width: {(gd.count / maxGrade) * 100}%; background: {gd.cls}"></div>
									</div>
									<span class="text-xs text-slate-600 w-8 text-right">{gd.count}</span>
								</div>
							{/if}
						{/each}
					</div>
				</div>

				<!-- Subject Averages -->
				<div class="bg-white rounded-xl border border-slate-200 p-6">
					<h2 class="text-base font-semibold text-slate-900 mb-4">Subject Averages</h2>
					<div class="space-y-3">
						{#each perfData.subjects as sub}
							<div>
								<div class="flex items-center justify-between mb-1">
									<span class="text-sm text-slate-700">{sub.subject_name} <span class="text-xs text-slate-400">({sub.subject_code})</span></span>
									<span class="text-xs font-medium px-1.5 py-0.5 rounded" style="color: {gradeColors[sub.grade]}">{sub.grade} · {sub.average.toFixed(1)}%</span>
								</div>
								<div class="h-2 bg-slate-100 rounded-full overflow-hidden">
									<div class="h-full rounded-full transition-all" style="width: {sub.average}%; background: {gradeColors[sub.grade]}"></div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			</div>

			<!-- Top Performers -->
			<div class="bg-white rounded-xl border border-slate-200 p-6 mt-6">
				<h2 class="text-base font-semibold text-slate-900 mb-4"><Award size={16} class="inline mr-1 text-amber-500" /> Top Performers</h2>
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3">
					{#each topStudents as s, i}
						<div class="border border-slate-200 rounded-lg p-3 text-center">
							<div class="text-lg font-bold" style="color: {gradeColors[s.grade]}">#{i + 1}</div>
							<div class="text-sm font-medium text-slate-800 truncate">{s.name}</div>
							<div class="text-xs text-slate-500">Roll {s.roll_no}</div>
							<div class="text-sm font-bold mt-1" style="color: {gradeColors[s.grade]}">{s.percentage.toFixed(1)}%</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

	{:else if activeTab === 'progress'}
		<!-- Marks Progress -->
		<div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
			<div class="flex flex-wrap gap-3">
				<div class="w-44">
					<Select bind:value={selectedClass} options={[{ id: '', name: 'All Classes' }, ...classes.map(c => ({ id: c.id, name: c.name }))]} placeholder="All Classes" />
				</div>
				<div class="flex-1 min-w-48">
					<SearchFilter bind:value={progressSearch} placeholder="Search assessments..." onInput={resetProgressPage} />
				</div>
			</div>
		</div>

		{#if totalProgress === 0}
			<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">No assessments found</div>
		{:else}
			<div class="bg-white rounded-xl border border-slate-200">
				<div class="px-6 py-4 border-b border-slate-200">
					<h2 class="text-base font-semibold text-slate-900">Marks Entry Progress</h2>
				</div>
				<div class="p-4 space-y-3">
					{#each paginatedProgress() as a}
						<div class="flex items-center gap-3">
							<div class="flex-1 min-w-0">
								<div class="flex items-center justify-between mb-1">
									<span class="text-sm text-slate-700 truncate">{a.name || 'Untitled'} <span class="text-xs text-slate-400">({a.class_name || a.class_id})</span></span>
									<span class="text-xs text-slate-500 shrink-0 ml-2">{a.marks_count ?? 0}/{a.student_count ?? 0}</span>
								</div>
								<div class="h-2 bg-slate-100 rounded-full overflow-hidden">
									<div class="h-full rounded-full transition-all {a.progress >= 80 ? 'bg-emerald-500' : a.progress >= 40 ? 'bg-amber-500' : 'bg-red-400'}" style="width: {a.progress}%"></div>
								</div>
							</div>
							{#if !a.is_published}
								<span class="text-[10px] px-1.5 py-0.5 rounded bg-amber-50 text-amber-600 shrink-0">Draft</span>
							{:else}
								<span class="text-[10px] px-1.5 py-0.5 rounded bg-emerald-50 text-emerald-600 shrink-0">Published</span>
							{/if}
						</div>
					{/each}
				</div>
			<Pagination total={totalProgress} pageSize={progressPageSize} page={progressPage} onChange={onProgressPageChange} />
		</div>

		<!-- Summary stats -->
			<div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-6">
				<div class="bg-white rounded-xl border border-slate-200 p-5 text-center">
					<div class="text-2xl font-bold text-emerald-600">{progressAssessments.filter(a => a.progress >= 80).length}</div>
					<div class="text-xs text-slate-500 mt-1">Complete (≥80%)</div>
				</div>
				<div class="bg-white rounded-xl border border-slate-200 p-5 text-center">
					<div class="text-2xl font-bold text-amber-600">{progressAssessments.filter(a => a.progress >= 40 && a.progress < 80).length}</div>
					<div class="text-xs text-slate-500 mt-1">In Progress</div>
				</div>
				<div class="bg-white rounded-xl border border-slate-200 p-5 text-center">
					<div class="text-2xl font-bold text-red-500">{progressAssessments.filter(a => a.progress < 40).length}</div>
					<div class="text-xs text-slate-500 mt-1">Pending (&lt;40%)</div>
				</div>
			</div>
		{/if}
	{/if}
</div>
