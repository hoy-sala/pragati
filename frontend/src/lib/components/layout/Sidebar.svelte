<script lang="ts">
	import { getAuthState, logout } from '$lib/stores/auth.svelte';
	import { page } from '$app/stores';
	import {
		LogOut, GraduationCap, LayoutDashboard, CalendarDays, Users, BookOpen,
		ClipboardCheck, ClipboardList, PlayCircle, Table, FileSpreadsheet,
		FileText, Settings, Heart, BarChart3, CheckSquare, Upload, HelpCircle, User, Award
	} from 'lucide-svelte';
	import type { ComponentType } from 'svelte';
	import type { User as UserType, Student } from '$lib/types';

	const auth = getAuthState();

	const roleLabels: Record<string, string> = {
		admin: 'Administrator', principal: 'Principal', teacher: 'Teacher',
		special_educator: 'Special Educator', student: 'Student', parent: 'Parent',
	};

	let displayName = $derived(
		auth.currentUser
			? (auth.currentUser as Student).first_name
				? `${(auth.currentUser as Student).first_name} ${(auth.currentUser as Student).last_name || ''}`.trim()
				: (auth.currentUser as UserType).name
			: ''
	);

	let effectiveRole = $derived(
		auth.currentUser
			? (auth.currentUser as Student).first_name
				? 'student'
				: (auth.currentUser as UserType).role
			: ''
	);

	let initials = $derived(
		displayName.split(' ').map(w => w[0]).filter(Boolean).slice(0, 2).join('').toUpperCase()
	);

	let displayTitle = $derived(roleLabels[effectiveRole] ?? effectiveRole.replace(/_/g, ' '));

	let roleBadgeColor = $derived(
		effectiveRole === 'admin' ? 'bg-purple-100 text-purple-700' :
		effectiveRole === 'principal' ? 'bg-blue-100 text-blue-700' :
		effectiveRole === 'teacher' ? 'bg-amber-100 text-amber-700' :
		effectiveRole === 'student' ? 'bg-emerald-100 text-emerald-700' :
		'bg-slate-100 text-slate-600'
	);

	type NavItem = {
		href: string;
		label: string;
		icon: ComponentType;
		roles: string[];
	};

	type NavSection = {
		label?: string;
		items: NavItem[];
	};

	const navSections: NavSection[] = [
		{
			items: [
				{ href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard, roles: ['admin', 'principal', 'teacher', 'special_educator', 'student', 'parent'] },
			]
		},
		{
			label: 'Academic',
			items: [
				{ href: '/timetable', label: 'Time Table', icon: CalendarDays, roles: ['admin', 'principal', 'teacher', 'special_educator', 'student', 'parent'] },
				{ href: '/classes', label: 'Classes', icon: Users, roles: ['admin', 'principal'] },
				{ href: '/subjects', label: 'Subjects', icon: BookOpen, roles: ['admin', 'principal'] },
				{ href: '/students', label: 'Students', icon: User, roles: ['admin', 'principal', 'teacher'] },
			]
		},
		{
			label: 'Assessment',
			items: [
				{ href: '/assessments', label: 'Assessments', icon: ClipboardCheck, roles: ['admin', 'principal', 'teacher'] },
				{ href: '/quizzes', label: 'Quizzes', icon: ClipboardList, roles: ['admin', 'principal', 'teacher'] },
				{ href: '/quizzes/available', label: 'Take Quiz', icon: PlayCircle, roles: ['admin', 'principal', 'teacher', 'student'] },
				{ href: '/marks', label: 'Marks Entry', icon: Table, roles: ['admin', 'principal', 'teacher'] },
				{ href: '/questions', label: 'Question Bank', icon: HelpCircle, roles: ['admin', 'principal', 'teacher'] },
			]
		},
		{
			label: 'Student Welfare',
			items: [
				{ href: '/mentors', label: 'Mentors', icon: Heart, roles: ['admin', 'principal', 'teacher', 'special_educator'] },
				{ href: '/hpc', label: 'HPC Cards', icon: FileSpreadsheet, roles: ['admin', 'principal', 'teacher'] },
				{ href: '/hpc/assess', label: 'LO Assessment', icon: CheckSquare, roles: ['admin', 'principal', 'teacher'] },
				{ href: '/hpc/lo-import', label: 'Import LOs', icon: Upload, roles: ['admin', 'principal'] },
			]
		},
		{
			label: 'Reports',
			items: [
				{ href: '/analytics', label: 'Analytics', icon: BarChart3, roles: ['admin', 'principal', 'teacher'] },
				{ href: '/reports', label: 'Reports', icon: FileText, roles: ['admin', 'principal', 'teacher', 'student', 'parent'] },
				{ href: '/certificates', label: 'Certificates', icon: Award, roles: ['admin'] },
			]
		},
		{
			label: 'System',
			items: [
				{ href: '/hpc/config', label: 'HPC Config', icon: Settings, roles: ['admin'] },
				{ href: '/settings', label: 'Settings', icon: Settings, roles: ['admin'] },
			]
		},
	];

	let visibleSections = $derived(
		navSections
			.map(section => ({
				label: section.label,
				items: section.items.filter(item => item.roles.includes(effectiveRole)),
			}))
			.filter(section => section.items.length > 0)
	);

	function isActive(href: string): boolean {
		const path = $page.url.pathname;
		if (href === '/dashboard') return path === '/dashboard';
		if (path === href) return true;
		if (path.startsWith(href + '/')) return true;
		return false;
	}

	function isSectionActive(items: NavItem[]): boolean {
		return items.some(item => isActive(item.href));
	}
</script>

<aside class="w-60 bg-white border-r border-slate-200 flex flex-col h-full no-print">
	<div class="px-4 py-4 border-b border-slate-200">
		<div class="flex items-center gap-3">
			<div class="w-9 h-9 rounded-lg bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shrink-0">
				<GraduationCap size={18} class="text-white" />
			</div>
			<div class="min-w-0">
				<h1 class="text-sm font-bold text-slate-800 font-kannada tracking-wide leading-tight">ಪ್ರಗತಿ</h1>
				<p class="text-[10px] text-slate-400 leading-tight">PRAGATI v1.0</p>
			</div>
		</div>
	</div>

	<nav class="flex-1 overflow-y-auto px-2 py-3 space-y-5 scrollbar-none">
		{#each visibleSections as section}
			<div>
				{#if section.label}
					<div class="px-2 mb-1.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider">
						{section.label}
					</div>
				{/if}
				<div class="space-y-0.5">
					{#each section.items as item (item.href)}
						<a
							href={item.href}
							class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-md text-[13px] transition-colors {isActive(item.href)
								? 'bg-primary-50 text-primary-700 font-medium'
								: 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
						>
							<item.icon size={16} class="{isActive(item.href) ? 'text-primary-600' : 'text-slate-400'}" />
							<span class="truncate">{item.label}</span>
						</a>
					{/each}
				</div>
			</div>
		{/each}
	</nav>

	<div class="border-t border-slate-200 p-2.5">
		<div class="flex items-center gap-2.5 px-2 py-1.5">
			<div class="w-8 h-8 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 text-xs font-semibold shrink-0">
				{initials}
			</div>
			<div class="min-w-0 flex-1">
				<div class="text-sm font-medium text-slate-800 truncate leading-tight">{displayName}</div>
				<div class="flex items-center gap-1.5 mt-0.5">
					<span class="inline-block w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
					<span class="text-[11px] text-slate-500">{displayTitle}</span>
				</div>
			</div>
			<button
				onclick={logout}
				title="Sign out"
				class="p-1.5 rounded-md text-slate-400 hover:text-danger-600 hover:bg-danger-50 transition-colors shrink-0"
			>
				<LogOut size={15} />
			</button>
		</div>
	</div>
</aside>
