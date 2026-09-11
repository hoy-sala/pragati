export interface PageTab {
	href: string;
	label: string;
	/** If set, the tab is only shown to these roles. */
	roles?: string[];
}

export const MENTOR_TABS: PageTab[] = [
	{ href: '/mentors?tab=assignments', label: 'Assignments' },
	{ href: '/mentors?tab=roster', label: 'Roster' },
	{ href: '/mentors?tab=attendance', label: 'Attendance' },
	{ href: '/mentors?tab=logs', label: 'Logs' },
	{ href: '/mentors?tab=dashboard', label: 'Dashboard', roles: ['admin', 'principal'] },
];

export const HPC_TABS: PageTab[] = [
	{ href: '/hpc?tab=grid', label: 'Progress Grid' },
	{ href: '/hpc?tab=assess', label: 'LO Assessment' },
	{ href: '/hpc?tab=import', label: 'Import LOs', roles: ['admin'] },
	{ href: '/hpc?tab=config', label: 'Config', roles: ['admin'] },
];

export const REPORT_TABS: PageTab[] = [
	{ href: '/reports?view=marksheet', label: 'Mark Sheet' },
	{ href: '/reports?view=report', label: 'Report Card' },
	{ href: '/reports?view=mentors', label: 'Mentor-wise' },
];
