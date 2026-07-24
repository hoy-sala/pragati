<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { initAuth, getAuthState } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';

	const { isAuthenticated, isLoading } = getAuthState();

	const publicRoutes = ['/login', '/timetable'];

	function isFullscreenRoute(path: string): boolean {
		return path.startsWith('/quizzes/take') || path.startsWith('/quizzes/results');
	}

	onMount(() => {
		initAuth();
	});

	$effect(() => {
		if (!isLoading && !isAuthenticated) {
			const path = $page.url.pathname;
			if (!publicRoutes.includes(path) && path !== '/') {
				goto('/login');
			}
		}
		if (!isLoading && isAuthenticated && $page.url.pathname === '/login') {
			goto('/dashboard');
		}
	});
</script>

{#if isLoading}
	<div class="flex h-screen items-center justify-center">
		<div class="text-slate-400 text-sm">Loading...</div>
	</div>
{:else if isAuthenticated}
	<div class="flex h-screen overflow-hidden">
		{#if !isFullscreenRoute($page.url.pathname)}
			<Sidebar />
		{/if}
		<main class="flex-1 overflow-y-auto p-6" class:p-0={isFullscreenRoute($page.url.pathname)}>
			<slot />
		</main>
	</div>
{:else}
	<main class="min-h-screen">
		<slot />
	</main>
{/if}
