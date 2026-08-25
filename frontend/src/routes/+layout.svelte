<script lang="ts">
	import '../app.css';
	import 'katex/dist/katex.min.css';
	import { onMount } from 'svelte';
	import { initAuth, getAuthState } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import Toast from '$lib/components/Toast.svelte';

	let { children } = $props();

	const auth = getAuthState();

	const publicRoutes = ['/login', '/timetable', '/play'];

	function isFullscreenRoute(path: string): boolean {
		return path.startsWith('/quizzes/take') || path.startsWith('/quizzes/results') || path.startsWith('/certificates/print') || path.startsWith('/play');
	}

	onMount(() => {
		initAuth();
	});

	$effect(() => {
		if (!auth.isLoading && !auth.isAuthenticated) {
			const path = $page.url.pathname;
			if (!publicRoutes.includes(path) && path !== '/') {
				goto('/login');
			}
		}
		if (!auth.isLoading && auth.isAuthenticated && $page.url.pathname === '/login') {
			goto('/dashboard');
		}
	});
</script>

{#if auth.isLoading}
	<div class="flex h-screen items-center justify-center">
		<div class="text-slate-400 text-sm">Loading...</div>
	</div>
{:else if auth.isAuthenticated}
	<div class="flex h-screen overflow-hidden">
		{#if !isFullscreenRoute($page.url.pathname)}
			<Sidebar />
		{/if}
		<main class="flex-1 overflow-y-auto p-6" class:p-0={isFullscreenRoute($page.url.pathname)}>
			{@render children()}
		</main>
	</div>
{:else}
	<main class="min-h-screen">
		{@render children()}
	</main>
{/if}
<Toast />
