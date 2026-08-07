<script lang="ts">
	import { getAuthState } from '$lib/stores/auth.svelte';
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';
	import type { User, Student } from '$lib/types';
	import { Sparkles, Trophy, Target, Flame, ArrowRight, PlayCircle, Award, CheckCircle2, BookOpen } from 'lucide-svelte';

	const auth = getAuthState();

	let isStudent = $derived(!!(auth.currentUser && (auth.currentUser as Student).first_name));

	let fullName = $derived(
		auth.currentUser
			? isStudent
				? `${(auth.currentUser as Student).first_name} ${(auth.currentUser as Student).last_name || ''}`.trim()
				: (auth.currentUser as User).name
			: ''
	);

	let firstName = $derived(fullName.split(' ')[0] || fullName);

	let loading = $state(true);

	let insights = $state<{
		student: { first_name: string; last_name: string; class_name: string };
		quizzes_taken: number;
		quizzes_passed: number;
		best_percentage: number;
		average_percentage: number;
		recent_attempts: {
			quiz_title: string;
			percentage: number;
			passed: boolean;
			submitted_at: string;
		}[];
	} | null>(null);

	let staffStats = $state<{
		total_students: number;
		total_teachers: number;
		total_classes: number;
		total_assessments: number;
		students_by_class: { class: string; count: number }[];
	} | null>(null);

	let staffDash = $state<{
		assignments: { class_name: string; subject_name: string }[];
		pending_assessments: { id: string; name: string; class_name: string; subject_name: string; max_marks: number; marks_count: number; student_count: number; due_date: string }[];
	} | null>(null);

	onMount(async () => {
		if (isStudent) {
			const res = await api<typeof insights>('GET', '/dashboard/student');
			if (res.data) insights = res.data;
		} else {
			const [statsRes, dashRes] = await Promise.all([
				api<typeof staffStats>('GET', '/dashboard/stats'),
				api<typeof staffDash>('GET', '/dashboard/staff'),
			]);
			if (statsRes.data) staffStats = statsRes.data;
			if (dashRes.data) staffDash = dashRes.data;
		}
		loading = false;
	});

	let bestScore = $derived(insights ? Math.round(insights.best_percentage) : 0);
	let avgScore = $derived(insights ? Math.round(insights.average_percentage) : 0);
	let scoreColor = $derived(bestScore >= 75 ? 'text-emerald-600' : bestScore >= 50 ? 'text-amber-600' : 'text-primary-600');

	let message = $derived(
		!insights || insights.quizzes_taken === 0
			? { title: 'Ready to start?', text: 'Take your first quiz and begin your learning journey!', icon: Sparkles }
			: bestScore >= 90
				? { title: 'Outstanding!', text: 'You are a quiz superstar. Keep shining bright!', icon: Award }
				: bestScore >= 70
					? { title: 'Great job!', text: 'You are doing wonderfully. A little more and you will be at the top!', icon: Trophy }
					: bestScore >= 50
						? { title: 'Good progress!', text: 'You are improving every day. Try again to score even higher!', icon: Target }
						: { title: 'Every expert was once a beginner', text: 'Practice makes perfect. Take a quiz and try again — you can do it!', icon: Flame }
	);

	let circleColor = $derived(bestScore >= 75 ? '#10b981' : bestScore >= 50 ? '#f59e0b' : '#3b82f6');
</script>

{#if isStudent}
	<div class="max-w-5xl mx-auto space-y-6">
		<div class="bg-gradient-to-br from-primary-600 via-primary-500 to-sky-400 rounded-2xl p-6 sm:p-8 text-white relative overflow-hidden">
			<div class="absolute -top-10 -right-10 w-48 h-48 bg-white/10 rounded-full blur-2xl"></div>
			<div class="absolute bottom-0 right-16 w-32 h-32 bg-white/10 rounded-full"></div>
			<div class="relative space-y-3">
				<div class="flex items-center gap-2 text-primary-100 text-sm font-medium">
					<span class="w-2 h-2 rounded-full bg-amber-300"></span>
					{insights?.student.class_name || 'Student'}
				</div>
				<h1 class="text-2xl sm:text-3xl font-bold">
					Welcome, {fullName}!
				</h1>
				<p class="text-primary-50 text-sm sm:text-base leading-relaxed">
					{message.title} {message.text}
				</p>
				<a href="/quizzes/available"
					class="inline-flex items-center gap-2 mt-2 px-5 py-2.5 bg-white text-primary-700 rounded-xl text-sm font-semibold shadow hover:shadow-lg hover:-translate-y-0.5 transition-all">
					<PlayCircle size={18} />
					Take a Quiz
					<ArrowRight size={16} />
				</a>
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
			<div class="bg-white rounded-2xl border border-slate-200 p-5 flex items-center gap-4">
				<div class="w-12 h-12 rounded-xl bg-primary-50 flex items-center justify-center">
					<BookOpen size={22} class="text-primary-600" />
				</div>
				<div>
					<div class="text-2xl font-bold text-slate-900">{loading ? '--' : insights?.quizzes_taken ?? 0}</div>
					<div class="text-sm text-slate-500">Quizzes taken</div>
				</div>
			</div>
			<div class="bg-white rounded-2xl border border-slate-200 p-5 flex items-center gap-4">
				<div class="w-12 h-12 rounded-xl bg-emerald-50 flex items-center justify-center">
					<Trophy size={22} class="text-emerald-600" />
				</div>
				<div>
					<div class="text-2xl font-bold text-slate-900">{loading ? '--' : insights?.quizzes_passed ?? 0}</div>
					<div class="text-sm text-slate-500">Quizzes passed</div>
				</div>
			</div>
			<div class="bg-white rounded-2xl border border-slate-200 p-5 flex items-center gap-4">
				<div class="w-12 h-12 rounded-xl bg-amber-50 flex items-center justify-center">
					<Target size={22} class="text-amber-600" />
				</div>
				<div>
					<div class="text-2xl font-bold {scoreColor}">{loading ? '--' : bestScore}%</div>
					<div class="text-sm text-slate-500">Best score</div>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-5 gap-4">
			<div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200 p-6">
				<h2 class="text-base font-semibold text-slate-900 mb-4">Your progress</h2>
				<div class="flex flex-col items-center gap-3">
					<div class="relative w-36 h-36">
						<svg viewBox="0 0 120 120" class="w-full h-full -rotate-90">
							<circle cx="60" cy="60" r="52" fill="none" stroke="#e2e8f0" stroke-width="12" />
							<circle cx="60" cy="60" r="52" fill="none" stroke={circleColor} stroke-width="12" stroke-linecap="round"
								stroke-dasharray="326.7" stroke-dashoffset={loading || !insights ? 326.7 : 326.7 * (1 - avgScore / 100)} />
						</svg>
						<div class="absolute inset-0 flex flex-col items-center justify-center">
							<span class="text-3xl font-bold text-slate-900">{loading ? '--' : avgScore}%</span>
							<span class="text-xs text-slate-500">avg score</span>
						</div>
					</div>
					<div class="flex items-center gap-2 text-sm text-slate-500">
						<Flame size={16} class="text-orange-500" />
						Keep learning — every quiz makes you smarter!
					</div>
				</div>
			</div>

			<div class="lg:col-span-3 bg-white rounded-2xl border border-slate-200 p-6">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-base font-semibold text-slate-900">Recent attempts</h2>
					<a href="/quizzes/available" class="text-sm text-primary-600 hover:text-primary-700 font-medium inline-flex items-center gap-1">
						View quizzes <ArrowRight size={14} />
					</a>
				</div>
				{#if loading}
					<div class="p-6 text-center text-sm text-slate-400">Loading...</div>
				{:else if !insights || insights.recent_attempts.length === 0}
					<div class="p-6 text-center space-y-2">
						<Sparkles size={28} class="mx-auto text-primary-300" />
						<p class="text-sm text-slate-500">No quizzes yet. Take your first quiz to see your results here!</p>
						<a href="/quizzes/available" class="inline-flex items-center gap-1.5 mt-1 px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
							<PlayCircle size={16} /> Start now
						</a>
					</div>
				{:else}
					<div class="space-y-3">
						{#each insights.recent_attempts as a}
							<div class="flex items-center gap-3 p-3 rounded-xl border border-slate-100 hover:bg-slate-50 transition-colors">
								<div class="w-9 h-9 rounded-lg shrink-0 flex items-center justify-center
									{a.passed ? 'bg-emerald-50 text-emerald-600' : 'bg-rose-50 text-rose-500'}">
									{#if a.passed}
										<CheckCircle2 size={18} />
									{:else}
										<Award size={18} />
									{/if}
								</div>
								<div class="min-w-0 flex-1">
									<div class="text-sm font-medium text-slate-800 truncate">{a.quiz_title}</div>
									<div class="text-xs text-slate-400">{a.submitted_at}</div>
								</div>
								<div class="text-right shrink-0">
									<div class="text-sm font-bold {a.passed ? 'text-emerald-600' : 'text-rose-500'}">{Math.round(a.percentage)}%</div>
									<div class="text-[11px] text-slate-400">{a.passed ? 'Passed' : 'Try again'}</div>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	</div>
{:else}
	<div class="space-y-6">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-bold text-slate-900">Dashboard</h1>
				<p class="text-sm text-slate-500 mt-1">Welcome, {fullName}</p>
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="bg-white rounded-xl border border-slate-200 p-4">
				<div class="text-2xl font-bold text-primary-600">{loading ? '--' : staffStats?.total_students ?? 0}</div>
				<div class="text-sm text-slate-500 mt-1">Students</div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-4">
				<div class="text-2xl font-bold text-primary-600">{loading ? '--' : staffStats?.total_teachers ?? 0}</div>
				<div class="text-sm text-slate-500 mt-1">Teachers</div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-4">
				<div class="text-2xl font-bold text-primary-600">{loading ? '--' : staffStats?.total_classes ?? 0}</div>
				<div class="text-sm text-slate-500 mt-1">Classes</div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-4">
				<div class="text-2xl font-bold text-primary-600">{loading ? '--' : staffStats?.total_assessments ?? 0}</div>
				<div class="text-sm text-slate-500 mt-1">Assessments</div>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- My Assignments -->
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<h2 class="text-base font-semibold text-slate-900 mb-3">My Classes & Subjects</h2>
				{#if loading}
					<div class="text-sm text-slate-400">Loading...</div>
				{:else if staffDash?.assignments && staffDash.assignments.length > 0}
					<div class="flex flex-wrap gap-2">
						{#each staffDash.assignments as a}
							<span class="text-xs px-2.5 py-1.5 rounded-lg bg-slate-100 text-slate-700">
								{a.class_name} · {a.subject_name}
							</span>
						{/each}
					</div>
				{:else}
					<div class="text-sm text-slate-400">No class assignments yet</div>
				{/if}
			</div>

			<!-- Pending Marks Entry -->
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<h2 class="text-base font-semibold text-slate-900 mb-3">Pending Marks Entry</h2>
				{#if loading}
					<div class="text-sm text-slate-400">Loading...</div>
				{:else if staffDash?.pending_assessments && staffDash.pending_assessments.length > 0}
					<div class="space-y-2">
						{#each staffDash.pending_assessments as p}
							<a href="/marks?assessment_id={p.id}" class="flex items-center gap-3 p-2.5 rounded-lg border border-slate-100 hover:bg-slate-50 transition-colors">
								<div class="flex-1 min-w-0">
									<div class="text-sm font-medium text-slate-800 truncate">{p.name || 'Untitled'}</div>
									<div class="text-xs text-slate-400">{p.class_name} · {p.subject_name}</div>
								</div>
								<div class="text-right shrink-0">
									<div class="text-xs font-medium {p.marks_count > 0 ? 'text-amber-600' : 'text-red-500'}">{p.marks_count}/{p.student_count}</div>
									<div class="text-[10px] text-slate-400">marks</div>
								</div>
							</a>
						{/each}
					</div>
				{:else}
					<div class="text-sm text-slate-400">All caught up! No pending marks entry.</div>
				{/if}
			</div>
		</div>

		{#if staffStats?.students_by_class && staffStats.students_by_class.length > 0}
			<div class="bg-white rounded-xl border border-slate-200 p-4">
				<h2 class="text-lg font-semibold text-slate-900 mb-3">Students by Class</h2>
				<div class="space-y-2">
					{#each staffStats.students_by_class as item}
						<div class="flex items-center gap-3">
							<span class="text-sm font-medium text-slate-700 w-20">{item.class}</span>
							<div class="flex-1 h-5 bg-slate-100 rounded-full overflow-hidden">
								<div class="h-full bg-primary-500 rounded-full transition-all" style="width: {Math.max(4, (item.count / Math.max(...staffStats.students_by_class.map(s => s.count))) * 100)}%"></div>
							</div>
							<span class="text-sm text-slate-600 w-8 text-right">{item.count}</span>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
{/if}
