<script lang="ts">
	import { login } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit() {
		error = '';
		loading = true;
		const result = await login(email, password);
		loading = false;
		if ('error' in result) {
			error = result.error;
		} else {
			goto('/dashboard');
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-slate-100 px-4">
	<div class="w-full max-w-sm">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-bold text-primary-600 font-kannada tracking-wider">ಪ್ರಗತಿ</h1>
			<p class="text-sm text-slate-500 mt-1">PRAGATI</p>
			<p class="text-xs text-slate-400 mt-2">Every Child Can Progress</p>
		</div>

		<form onsubmit={handleSubmit} class="bg-white rounded-xl shadow-sm border border-slate-200 p-6 space-y-4">
			<div>
				<label for="email" class="block text-sm font-medium text-slate-700 mb-1">Email</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
					placeholder="you@school.edu"
				/>
			</div>

			<div>
				<label for="password" class="block text-sm font-medium text-slate-700 mb-1">Password</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
				/>
			</div>

			{#if error}
				<div class="text-sm text-danger-600 bg-danger-50 rounded-lg p-3">{error}</div>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="w-full py-2 px-4 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
			>
				{loading ? 'Signing in...' : 'Sign in'}
			</button>
		</form>

		<div class="text-center mt-6">
			<a href="/timetable" class="text-xs text-slate-400 hover:text-primary-600 transition-colors underline underline-offset-2">
				View Master School Timetable 2026-27
			</a>
		</div>
	</div>
</div>
