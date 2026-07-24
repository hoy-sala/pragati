import { api, setTokens, clearTokens, loadTokens } from '$lib/api/client.svelte';
import type { User } from '$lib/types';

let currentUser: User | null = $state(null);
let isAuthenticated = $derived(currentUser !== null);
let isLoading = $state(true);

export function getAuthState() {
	return {
		get currentUser() { return currentUser; },
		get isAuthenticated() { return isAuthenticated; },
		get isLoading() { return isLoading; }
	};
}

export async function initAuth() {
	loadTokens();
	const token = typeof localStorage !== 'undefined' ? localStorage.getItem('access_token') : null;
	if (!token) {
		isLoading = false;
		return;
	}

	try {
		const res = await api<User>('GET', '/auth/me');
		if (res.data) {
			currentUser = res.data;
		} else {
			clearTokens();
		}
	} catch {
		clearTokens();
	}
	isLoading = false;
}

export async function login(email: string, password: string): Promise<{ user: User } | { error: string }> {
	const res = await api<{ user: User; access_token: string; refresh_token: string; expires_in: number }>(
		'POST', '/auth/login', { email, password }
	);

	if (res.error) {
		return { error: res.error.message };
	}

	if (res.data) {
		setTokens(res.data.access_token, res.data.refresh_token);
		currentUser = res.data.user;
		return { user: res.data.user };
	}

	return { error: 'unexpected error' };
}

export async function logout() {
	await api('POST', '/auth/logout');
	clearTokens();
	currentUser = null;
}
