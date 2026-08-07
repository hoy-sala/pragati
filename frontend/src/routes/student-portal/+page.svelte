<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { getAuthState } from '$lib/stores/auth.svelte';
	import { onMount } from 'svelte';
	import { User, Award, BookOpen, ClipboardCheck, TrendingUp, CheckCircle2, XCircle } from 'lucide-svelte';
	import type { AcademicYear, Student } from '$lib/types';

	const auth = getAuthState();
	const student = $derived(auth.currentUser as Student | null);

	let years = $state<AcademicYear[]>([]);
	let selectedYear = $state('');

	let report = $state<{
		student: { name: string; class: string; section?: string; roll_no: number; sats_number: string; gender?: string; date_of_birth?: string };
		academic_year: string;
		subjects: { subject_id: string; subject_code: string; subject_name: string; subject_type: string; assessments: { id: string; name: string; category: string; max: number; value: number; absent: boolean; has_mark: boolean }[]; total: number; max_max: number; percentage: number; grade: string; grade_label?: string }[];
		grand_total: number;
		grand_max: number;
		percentage: number;
		grade: string;
	} | null>(null);

	let quizHistory = $state<{
		quiz_title: string;
		percentage: number;
		passed: boolean;
		submitted_at: string;
	}[]>([]);

	let loading = $state(true);
	let activeTab = $state<'marks' | 'quizzes'>('marks');

	const gradeColors: Record<string, string> = {
		'A+': '#10b981', 'A': '#22c55e', 'B+': '#3b82f6', 'B': '#0ea5e9',
		'C+': '#f59e0b', 'C': '#f97316', 'D': '#ef4444', 'F': '#dc2626',
	};

	async function loadReport() {
		if (!selectedYear) { report = null; return; }
		const res = await api<typeof report>('GET', `/reports/student-me?academic_year_id=${selectedYear}`);
		if (res.data) report = res.data;
	}

	async function loadQuizHistory() {
		const res = await api<{ recent_attempts: typeof quizHistory }>('GET', '/dashboard/student');
		if (res.data) quizHistory = res.data.recent_attempts || [];
	}

	$effect(() => {
		if (selectedYear) loadReport();
	});

	onMount(async () => {
		const [yrRes, qRes] = await Promise.all([
			api<AcademicYear[]>('GET', '/academic-years'),
			api<{ recent_attempts: typeof quizHistory }>('GET', '/dashboard/student'),
		]);
		if (yrRes.data) {
			years = yrRes.data;
			const cur = yrRes.data.find(y => y.is_current);
			if (cur) { selectedYear = cur.id; await loadReport(); }
		}
		if (qRes.data) quizHistory = qRes.data.recent_attempts || [];
		loading = false;
	});
</script>

<svelte:head>
	<title>My Portal — Pragati</title>
</svelte:head>

<div class="max-w-5xl mx-auto space-y-6">
	<!-- Profile Header -->
	<div class="bg-gradient-to-br from-primary-600 via-primary-500 to-sky-400 rounded-2xl p-6 text-white relative overflow-hidden">
		<div class="absolute -top-10 -right-10 w-48 h-48 bg-white/10 rounded-full blur-2xl"></div>
		<div class="relative flex items-center gap-4">
			<div class="w-16 h-16 rounded-2xl bg-white/20 flex items-center justify-center">
				<User size={32} class="text-white" />
			</div>
			<div>
				<h1 class="text-2xl font-bold">{student?.first_name} {student?.last_name || ''}</h1>
				<div class="flex items-center gap-3 mt-1 text-primary-100 text-sm">
					<span>{report?.student.class || ''}</span>
					<span>Roll {report?.student.roll_no || ''}</span>
					<span>SATS: {report?.student.sats_number || ''}</span>
				</div>
			</div>
			{#if report}
				<div class="ml-auto text-right hidden sm:block">
					<div class="text-3xl font-bold">{report.percentage.toFixed(1)}%</div>
					<div class="text-primary-100 text-sm">Overall</div>
				</div>
			{/if}
		</div>
	</div>

	<!-- Tabs -->
	<div class="flex gap-1 bg-slate-100 rounded-lg p-1 no-print">
		<button onclick={() => activeTab = 'marks'}
			class="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab === 'marks' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
			<BookOpen size={16} /> My Marks
		</button>
		<button onclick={() => activeTab = 'quizzes'}
			class="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab === 'quizzes' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
			<ClipboardCheck size={16} /> Quiz History
		</button>
	</div>

	{#if loading}
		<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Loading...</div>

	{:else if activeTab === 'marks'}
		<!-- Marks Tab -->
		<div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
			<div class="w-56">
				<select bind:value={selectedYear} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">Select academic year</option>
					{#each years as y}
						<option value={y.id}>{y.name}</option>
					{/each}
				</select>
			</div>
		</div>

		{#if !report}
			<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">
				{#if !selectedYear}
					Select an academic year to view your marks
				{:else}
					No marks found for this year
				{/if}
			</div>
		{:else}
			<!-- Overall Summary -->
			<div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
				<div class="bg-white rounded-xl border border-slate-200 p-4 text-center">
					<div class="text-2xl font-bold text-slate-900">{report.grand_total}</div>
					<div class="text-xs text-slate-500 mt-1">Total Marks</div>
				</div>
				<div class="bg-white rounded-xl border border-slate-200 p-4 text-center">
					<div class="text-2xl font-bold text-slate-900">{report.grand_max}</div>
					<div class="text-xs text-slate-500 mt-1">Max Marks</div>
				</div>
				<div class="bg-white rounded-xl border border-slate-200 p-4 text-center">
					<div class="text-2xl font-bold" style="color: {gradeColors[report.grade]}">{report.percentage.toFixed(1)}%</div>
					<div class="text-xs text-slate-500 mt-1">Percentage</div>
				</div>
				<div class="bg-white rounded-xl border border-slate-200 p-4 text-center">
					<div class="text-2xl font-bold" style="color: {gradeColors[report.grade]}">{report.grade}</div>
					<div class="text-xs text-slate-500 mt-1">Grade</div>
				</div>
			</div>

			<!-- Subject-wise Report -->
			<div class="bg-white rounded-xl border border-slate-200 mt-6">
				<div class="px-6 py-4 border-b border-slate-200">
					<h2 class="text-base font-semibold text-slate-900">Subject-wise Performance</h2>
				</div>
				<div class="divide-y divide-slate-100">
					{#each report.subjects as sub}
						<div class="px-6 py-4">
							<div class="flex items-center justify-between mb-2">
								<div class="flex items-center gap-3">
									<span class="text-sm font-medium text-slate-800">{sub.subject_name}</span>
									<span class="text-xs text-slate-400">({sub.subject_code})</span>
									{#if sub.subject_type === 'curricular'}
										<span class="text-[10px] px-1.5 py-0.5 rounded bg-blue-50 text-blue-600">Curricular</span>
									{:else}
										<span class="text-[10px] px-1.5 py-0.5 rounded bg-purple-50 text-purple-600">Co-curricular</span>
									{/if}
								</div>
								<div class="flex items-center gap-3">
									<span class="text-sm text-slate-600">{sub.total}/{sub.max_max}</span>
									<span class="text-sm font-bold" style="color: {gradeColors[sub.grade]}">{sub.percentage.toFixed(1)}%</span>
									<span class="text-sm font-bold px-2 py-0.5 rounded" style="color: {gradeColors[sub.grade]}; background: {gradeColors[sub.grade]}15">{sub.grade}</span>
								</div>
							</div>
							<div class="flex flex-wrap gap-1.5">
								{#each sub.assessments as a}
									<span class="text-[11px] px-2 py-1 rounded border {a.absent ? 'bg-red-50 border-red-200 text-red-600' : a.has_mark ? 'bg-emerald-50 border-emerald-200 text-emerald-700' : 'bg-slate-50 border-slate-200 text-slate-400'}">
										{a.category}: {#if a.absent}ABS{:else}{a.has_mark ? a.value : '—'}{/if}/{a.max}
									</span>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

	{:else if activeTab === 'quizzes'}
		<!-- Quiz History Tab -->
		{#if quizHistory.length === 0}
			<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">
				<ClipboardCheck size={40} class="mx-auto text-slate-300 mb-3" />
				<p>No quiz attempts yet</p>
				<a href="/quizzes/available" class="text-sm text-primary-600 hover:text-primary-700 mt-2 inline-block">Take a quiz</a>
			</div>
		{:else}
			<div class="bg-white rounded-xl border border-slate-200">
				<div class="px-6 py-4 border-b border-slate-200">
					<h2 class="text-base font-semibold text-slate-900">Recent Quiz Attempts</h2>
				</div>
				<div class="divide-y divide-slate-100">
					{#each quizHistory as q}
						<div class="px-6 py-3 flex items-center gap-3">
							<div class="w-9 h-9 rounded-lg shrink-0 flex items-center justify-center {q.passed ? 'bg-emerald-50 text-emerald-600' : 'bg-rose-50 text-rose-500'}">
								{#if q.passed}<CheckCircle2 size={18} />{:else}<XCircle size={18} />{/if}
							</div>
							<div class="min-w-0 flex-1">
								<div class="text-sm font-medium text-slate-800 truncate">{q.quiz_title}</div>
								<div class="text-xs text-slate-400">{q.submitted_at}</div>
							</div>
							<div class="text-right shrink-0">
								<div class="text-sm font-bold {q.passed ? 'text-emerald-600' : 'text-rose-500'}">{Math.round(q.percentage)}%</div>
								<div class="text-[11px] text-slate-400">{q.passed ? 'Passed' : 'Try again'}</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>
