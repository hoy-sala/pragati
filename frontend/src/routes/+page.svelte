<script lang="ts">
	import { getAuthState } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	const auth = getAuthState();

	$effect(() => {
		if (!auth.isLoading && !auth.isAuthenticated) {
			goto('/login');
		} else if (!auth.isLoading && auth.isAuthenticated) {
			goto('/dashboard');
		}
	});
</script>

{#if auth.isLoading}
	<div class="flex h-screen items-center justify-center">
		<div class="text-slate-400 text-sm">Loading...</div>
	</div>
{/if}
