export interface PageTab {
	href: string;
	label: string;
	/** If set, the tab is only shown to these roles. */
	roles?: string[];
}

export const MENTOR_TABS: PageTab[] = [
	{ href: '/mentors', label: 'Assignments' },
	{ href: '/mentors/roster', label: 'Roster' },
	{ href: '/mentors/attendance', label: 'Attendance' },
	{ href: '/mentors/logs', label: 'Logs' },
	{ href: '/mentors/dashboard', label: 'Dashboard', roles: ['admin', 'principal'] },
];

export const HPC_TABS: PageTab[] = [
	{ href: '/hpc', label: 'Progress Grid' },
	{ href: '/hpc/assess', label: 'LO Assessment' },
	{ href: '/hpc/lo-import', label: 'Import LOs', roles: ['admin'] },
	{ href: '/hpc/config', label: 'Config', roles: ['admin'] },
];
