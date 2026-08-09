<script lang="ts">
	import { getAuthState, logout } from '$lib/stores/auth.svelte';
	import { page } from '$app/stores';
	import { LogOut, GraduationCap } from 'lucide-svelte';
	import type { User, Student } from '$lib/types';
	import Search from './Search.svelte';

	const auth = getAuthState();

	let displayName = $derived(
		auth.currentUser
			? (auth.currentUser as Student).first_name
				? `${(auth.currentUser as Student).first_name} ${(auth.currentUser as Student).last_name || ''}`
				: (auth.currentUser as User).name
			: ''
	);

	let displayRole = $derived(
		auth.currentUser
			? (auth.currentUser as Student).first_name
				? 'student'
				: (auth.currentUser as User).role
			: ''
	);

	let effectiveRole = $derived(
		auth.currentUser
			? (auth.currentUser as Student).first_name
				? 'student'
				: (auth.currentUser as User).role
			: ''
	);

	let initials = $derived(
		displayName
			.split(' ')
			.map(w => w[0])
			.filter(Boolean)
			.slice(0, 2)
			.join('')
			.toUpperCase()
	);

	const roleLabels: Record<string, string> = {
		admin: 'Administrator',
		principal: 'Principal',
		teacher: 'Teacher',
		special_educator: 'Special Educator',
		student: 'Student',
		parent: 'Parent',
	};

	let displayTitle = $derived(roleLabels[displayRole] ?? displayRole.replace(/_/g, ' '));

	let roleBadgeColor = $derived(
		displayRole === 'admin' ? 'bg-purple-100 text-purple-700' :
		displayRole === 'principal' ? 'bg-blue-100 text-blue-700' :
		displayRole === 'teacher' ? 'bg-amber-100 text-amber-700' :
		displayRole === 'student' ? 'bg-emerald-100 text-emerald-700' :
		'bg-slate-100 text-slate-600'
	);

	const navItems = [
		{ href: '/dashboard', label: 'Dashboard', icon: 'LayoutDashboard', roles: ['admin', 'principal', 'teacher', 'special_educator', 'student', 'parent'] },
		{ href: '/timetable', label: 'Time Table', icon: 'CalendarDays', roles: ['admin', 'principal', 'teacher', 'special_educator', 'student', 'parent'] },
		{ href: '/mentors', label: 'Mentors', icon: 'Heart', roles: ['admin', 'principal', 'teacher', 'special_educator'] },
		{ href: '/students', label: 'Students', icon: 'Users', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/classes', label: 'Classes', icon: 'GraduationCap', roles: ['admin', 'principal'] },
		{ href: '/subjects', label: 'Subjects', icon: 'BookOpen', roles: ['admin', 'principal'] },
		{ href: '/questions', label: 'Question Bank', icon: 'HelpCircle', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/assessments', label: 'Assessments', icon: 'ClipboardCheck', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/quizzes', label: 'Quizzes', icon: 'ClipboardList', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/quizzes/available', label: 'Take Quiz', icon: 'PlayCircle', roles: ['admin', 'principal', 'teacher', 'student'] },
		{ href: '/marks', label: 'Marks Entry', icon: 'Table', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/hpc', label: 'HPC Cards', icon: 'FileSpreadsheet', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/hpc/assess', label: 'LO Assessment', icon: 'CheckSquare', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/hpc/lo-import', label: 'Import LOs', icon: 'Upload', roles: ['admin', 'principal'] },
		{ href: '/hpc/config', label: 'HPC Config', icon: 'Settings', roles: ['admin'] },
		{ href: '/analytics', label: 'Analytics', icon: 'BarChart3', roles: ['admin', 'principal', 'teacher'] },
		{ href: '/reports', label: 'Reports', icon: 'FileText', roles: ['admin', 'principal', 'teacher', 'student', 'parent'] },
		{ href: '/settings', label: 'Settings', icon: 'Settings', roles: ['admin'] },
	];

	let visibleItems = $derived(
		navItems.filter(item => item.roles.includes(effectiveRole))
	);

	function isActive(href: string): boolean {
		return $page.url.pathname.startsWith(href);
	}
</script>

<aside class="w-64 bg-white border-r border-slate-200 flex flex-col h-full no-print">
	<div class="px-4 py-5 border-b border-slate-200 bg-gradient-to-r from-primary-50/60 to-white">
		<div class="flex items-center gap-3">
			<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shadow-sm shadow-primary-200">
				<GraduationCap size={20} class="text-white" />
			</div>
			<div>
				<h1 class="text-base font-bold text-slate-800 font-kannada tracking-wide leading-tight">ಪ್ರಗತಿ</h1>
				<div class="flex items-center gap-1.5 mt-0.5">
					<span class="text-[10px] font-semibold text-primary-600 uppercase tracking-[0.2em]">PRAGATI</span>
					<span class="w-1 h-1 rounded-full bg-primary-300"></span>
					<span class="text-[10px] text-slate-400">v1.0</span>
				</div>
			</div>
		</div>
		<div class="mt-2.5 pt-2.5 border-t border-slate-200/60">
			<p class="text-[11px] text-slate-400 italic leading-relaxed">Every Child Can Progress</p>
		</div>
	</div>

	<Search />

	<nav class="flex-1 overflow-y-auto p-2 space-y-1 scrollbar-none">
		{#each visibleItems as item (item.href)}
			<a
				href={item.href}
				class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-slate-600 hover:bg-slate-100 hover:text-slate-900 transition-colors"
				class:bg-primary-50={isActive(item.href)}
				class:text-primary-700={isActive(item.href)}
				class:font-medium={isActive(item.href)}
			>
				{item.label}
			</a>
		{/each}
	</nav>

	<div class="border-t border-slate-200 p-3">
		<div class="flex items-center gap-3 px-2 py-2 rounded-lg bg-slate-50">
			<div class="w-9 h-9 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 text-xs font-semibold shrink-0">
				{initials}
			</div>
			<div class="min-w-0 flex-1">
				<div class="text-sm font-medium text-slate-800 truncate">{displayName}</div>
				<div class="flex items-center gap-1.5 mt-0.5">
					<span class="inline-block w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
					<span class="text-xs text-slate-500">{displayTitle}</span>
				</div>
			</div>
			<button
				onclick={logout}
				title="Sign out"
				class="p-1.5 rounded-md text-slate-400 hover:text-danger-600 hover:bg-danger-50 transition-colors shrink-0"
			>
				<LogOut size={16} />
			</button>
		</div>
	</div>
</aside>
