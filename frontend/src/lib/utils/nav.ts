import type { ComponentType } from 'svelte';
import {
	GraduationCap, CalendarDays, House,
	ClipboardCheck, ClipboardList, Table, FileSpreadsheet,
	FileText, Settings, Heart, HelpCircle, User, Award
} from 'lucide-svelte';
import type { User as UserType, Student } from '$lib/types';
import { effectiveRole } from '$lib/utils/roles';

export type NavItem = {
	href: string;
	label: string;
	icon: ComponentType;
	roles: string[];
};

export type NavSection = {
	label?: string;
	items: NavItem[];
};

export const NAV_SECTIONS: NavSection[] = [
	{
		items: [
			{ href: '/home', label: 'Home', icon: House, roles: ['admin', 'principal', 'teacher', 'special_educator'] },
		]
	},
	{
		label: 'Academic',
		items: [
			{ href: '/timetable', label: 'Time Table', icon: CalendarDays, roles: ['admin', 'principal', 'teacher', 'special_educator', 'student', 'parent'] },
			{ href: '/students', label: 'Students', icon: User, roles: ['admin', 'principal', 'teacher'] },
		]
	},
	{
		label: 'Assessment',
		items: [
			{ href: '/assessments', label: 'Assessments', icon: ClipboardCheck, roles: ['admin', 'principal', 'teacher'] },
			{ href: '/marks', label: 'Marks Entry', icon: Table, roles: ['admin', 'principal', 'teacher'] },
			{ href: '/questions', label: 'Question Bank', icon: HelpCircle, roles: ['admin', 'principal', 'teacher'] },
			{ href: '/quizzes', label: 'Quizzes', icon: ClipboardList, roles: ['admin', 'principal', 'teacher'] },
		]
	},
	{
		label: 'Student Welfare',
		items: [
			{ href: '/mentors', label: 'Mentors', icon: Heart, roles: ['admin', 'principal', 'teacher', 'special_educator'] },
			{ href: '/hpc', label: 'HPC Cards', icon: FileSpreadsheet, roles: ['admin', 'principal', 'teacher'] },
		]
	},
	{
		label: 'Reports',
		items: [
			{ href: '/reports', label: 'Reports', icon: FileText, roles: ['admin', 'principal', 'teacher', 'student', 'parent'] },
			{ href: '/certificates', label: 'Certificates', icon: Award, roles: ['admin'] },
		]
	},
	{
		label: 'System',
		items: [
			{ href: '/settings', label: 'Settings', icon: Settings, roles: ['admin'] },
		]
	},
];

export const ROLE_LABELS: Record<string, string> = {
	admin: 'Administrator', principal: 'Principal', teacher: 'Teacher',
	special_educator: 'Special Educator', student: 'Student', parent: 'Parent',
};

export function allHrefs(sections: NavSection[] = NAV_SECTIONS): string[] {
	return sections.flatMap(s => s.items.map(i => i.href));
}

/** Longest-prefix wins so /mentors doesn't stay active on /mentors/roster.
 *  Hrefs with a query (?tab=) match exactly against pathname + search. */
export function isNavActive(href: string, pathname: string, search = '', hrefs: string[] = allHrefs()): boolean {
	if (href.includes('?')) return pathname + search === href;
	if (pathname === href) return true;
	if (!pathname.startsWith(href + '/')) return false;
	return !hrefs.some(h => {
		if (h.includes('?') || h === href || h.length <= href.length) return false;
		return pathname === h || pathname.startsWith(h + '/');
	});
}

export function visibleSections(role: string, sections: NavSection[] = NAV_SECTIONS): NavSection[] {
	return sections
		.map(section => ({
			label: section.label,
			items: section.items.filter(item => item.roles.includes(role)),
		}))
		.filter(section => section.items.length > 0);
}

export function userDisplayName(user: UserType | Student | null | undefined): string {
	if (!user) return '';
	if ((user as Student).first_name) {
		return `${(user as Student).first_name} ${(user as Student).last_name || ''}`.trim();
	}
	return (user as UserType).name;
}

export function userInitials(name: string): string {
	return name.split(' ').map(w => w[0]).filter(Boolean).slice(0, 2).join('').toUpperCase();
}

export function roleTitle(role: string): string {
	return ROLE_LABELS[role] ?? role.replace(/_/g, ' ');
}

/** Bottom-tab priority for the mobile app bar. First 4 visible to the role win. */
const BOTTOM_TAB_ORDER = [
	'/home', '/marks', '/students', '/reports', '/mentors', '/assessments',
	'/timetable', '/quizzes', '/questions', '/hpc', '/certificates', '/settings',
];

export function bottomTabs(role: string): NavItem[] {
	const visible = visibleSections(role).flatMap(s => s.items);
	const byHref = new Map(visible.map(i => [i.href, i]));
	const tabs: NavItem[] = [];
	for (const href of BOTTOM_TAB_ORDER) {
		const item = byHref.get(href);
		if (item) tabs.push(item);
		if (tabs.length === 4) break;
	}
	return tabs;
}

export { effectiveRole };
