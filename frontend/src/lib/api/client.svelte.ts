import type { APIResponse } from '$lib/types';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:5050';

let accessToken: string | null = $state(null);
let refreshToken: string | null = $state(null);

export function setTokens(access: string, refresh: string) {
	accessToken = access;
	refreshToken = refresh;
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem('access_token', access);
		localStorage.setItem('refresh_token', refresh);
	}
}

export function loadTokens() {
	if (typeof localStorage !== 'undefined') {
		accessToken = localStorage.getItem('access_token');
		refreshToken = localStorage.getItem('refresh_token');
	}
}

export function clearTokens() {
	accessToken = null;
	refreshToken = null;
	if (typeof localStorage !== 'undefined') {
		localStorage.removeItem('access_token');
		localStorage.removeItem('refresh_token');
	}
}

async function refreshAccessToken(): Promise<boolean> {
	if (!refreshToken) return false;

	try {
		const res = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ refresh_token: refreshToken })
		});
		if (!res.ok) return false;
		const json: APIResponse<{ access_token: string; refresh_token: string }> = await res.json();
		if (json.data) {
			setTokens(json.data.access_token, json.data.refresh_token);
			return true;
		}
		return false;
	} catch {
		return false;
	}
}

export async function api<T = unknown>(
	method: string,
	path: string,
	body?: unknown
): Promise<APIResponse<T>> {
	const url = `${API_BASE}/api/v1${path}`;

	const headers: Record<string, string> = {
		'Content-Type': 'application/json'
	};

	if (accessToken) {
		headers['Authorization'] = `Bearer ${accessToken}`;
	}

	let res = await fetch(url, {
		method,
		headers,
		body: body ? JSON.stringify(body) : undefined
	});

	if (res.status === 401 && refreshToken) {
		const refreshed = await refreshAccessToken();
		if (refreshed) {
			headers['Authorization'] = `Bearer ${accessToken}`;
			res = await fetch(url, {
				method,
				headers,
				body: body ? JSON.stringify(body) : undefined
			});
		}
	}

	return res.json();
}

export function apiUrl(path: string): string {
	return `${API_BASE}/api/v1${path}`;
}
