<script lang="ts">
	import { getAuthState, logout } from '$lib/stores/auth.svelte';
	import { page } from '$app/stores';

	const { currentUser } = getAuthState();

	const navItems = [
		{ href: '/dashboard', label: 'Dashboard', icon: 'LayoutDashboard', roles: ['admin', 'principal', 'teacher', 'special_educator', 'student', 'parent'] },
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
		navItems.filter(item => item.roles.includes(currentUser?.role ?? ''))
	);

	function isActive(href: string): boolean {
		return $page.url.pathname.startsWith(href);
	}
</script>

<aside class="w-64 bg-white border-r border-slate-200 flex flex-col h-full no-print">
	<div class="p-4 border-b border-slate-200">
		<h1 class="text-lg font-bold text-primary-600 font-kannada tracking-wider">ಪ್ರಗತಿ</h1>
		<p class="text-xs text-slate-500">PRAGATI</p>
	</div>

	<nav class="flex-1 overflow-y-auto p-2 space-y-1">
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

	<div class="p-4 border-t border-slate-200">
		<div class="text-sm text-slate-700">{currentUser?.name}</div>
		<div class="text-xs text-slate-400">{currentUser?.role}</div>
		<button
			onclick={logout}
			class="mt-2 text-xs text-danger-600 hover:text-danger-700 transition-colors"
		>
			Sign out
		</button>
	</div>
</aside>
