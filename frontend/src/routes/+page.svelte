<script lang="ts">
	import { getAuthState } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	const { isAuthenticated, isLoading } = getAuthState();

	$effect(() => {
		if (!isLoading && !isAuthenticated) {
			goto('/login');
		} else if (!isLoading && isAuthenticated) {
			goto('/dashboard');
		}
	});
</script>

{#if isLoading}
	<div class="flex h-screen items-center justify-center">
		<div class="text-slate-400 text-sm">Loading...</div>
	</div>
{/if}
