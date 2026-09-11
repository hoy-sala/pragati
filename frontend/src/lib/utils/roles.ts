import type { Student, User } from '$lib/types';

/** Mirrors Sidebar's role derivation: students are identified by first_name, staff by role. */
export function effectiveRole(user: User | Student | null | undefined): string {
	if (!user) return '';
	if ((user as Student).first_name) return 'student';
	return (user as User).role ?? '';
}

export function hasRole(user: User | Student | null | undefined, ...roles: string[]): boolean {
	return roles.includes(effectiveRole(user));
}
