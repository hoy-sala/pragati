<script lang="ts">
	import { login, staffLogin, studentLogin } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { GraduationCap, Phone, Key, Hash, Calendar, Shield } from 'lucide-svelte';

	type Tab = 'student' | 'staff' | 'admin';
	let activeTab = $state<Tab>('student');

	let satsNumber = $state('');
	let dateOfBirth = $state('');

	let mobile = $state('');
	let staffPassword = $state('');

	let adminEmail = $state('');
	let adminPassword = $state('');

	let error = $state('');
	let loading = $state(false);

	async function handleStudentSubmit() {
		error = '';
		loading = true;
		let result;
		try {
			result = await studentLogin(satsNumber, dateOfBirth);
		} catch (e) {
			result = { error: 'Unable to reach server. Check your connection.' };
		}
		loading = false;
		if ('error' in result) {
			error = result.error;
		} else {
			goto('/dashboard');
		}
	}

	async function handleStaffSubmit() {
		error = '';
		loading = true;
		let result;
		try {
			result = await staffLogin(mobile, staffPassword);
		} catch (e) {
			result = { error: 'Unable to reach server. Check your connection.' };
		}
		loading = false;
		if ('error' in result) {
			error = result.error;
		} else {
			goto('/dashboard');
		}
	}

	async function handleAdminSubmit() {
		error = '';
		loading = true;
		let result;
		try {
			result = await login(adminEmail, adminPassword);
		} catch (e) {
			result = { error: 'Unable to reach server. Check your connection.' };
		}
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
			<div class="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-700 shadow-lg shadow-primary-200 mb-4">
				<GraduationCap size={28} class="text-white" />
			</div>
			<h1 class="text-3xl font-bold text-slate-800 font-kannada tracking-wide">ಪ್ರಗತಿ</h1>
			<div class="flex items-center justify-center gap-2 mt-1.5">
				<span class="text-sm font-semibold text-primary-600 uppercase tracking-[0.25em]">PRAGATI</span>
				<span class="w-1 h-1 rounded-full bg-primary-300"></span>
				<span class="text-xs text-slate-400">v1.0</span>
			</div>
			<div class="mt-3 pt-3 border-t border-slate-200 max-w-[200px] mx-auto">
				<p class="text-xs text-slate-400 italic">Every Child Can Progress</p>
			</div>
		</div>

		<div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6 space-y-4">
			<div class="flex border border-slate-200 rounded-lg overflow-hidden">
				<button
					onclick={() => { activeTab = 'student'; error = ''; }}
					class="flex-1 py-2 text-sm font-medium transition-colors
						{activeTab === 'student' ? 'bg-primary-600 text-white' : 'bg-white text-slate-600 hover:bg-slate-50'}"
				>
					Student
				</button>
				<button
					onclick={() => { activeTab = 'staff'; error = ''; }}
					class="flex-1 py-2 text-sm font-medium transition-colors
						{activeTab === 'staff' ? 'bg-primary-600 text-white' : 'bg-white text-slate-600 hover:bg-slate-50'}"
				>
					Staff
				</button>
				<button
					onclick={() => { activeTab = 'admin'; error = ''; }}
					class="flex-1 py-2 text-sm font-medium transition-colors
						{activeTab === 'admin' ? 'bg-primary-600 text-white' : 'bg-white text-slate-600 hover:bg-slate-50'}"
				>
					Admin
				</button>
			</div>

			{#if activeTab === 'student'}
				<form onsubmit={handleStudentSubmit} class="space-y-4">
					<div class="relative">
						<Hash size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							type="text"
							bind:value={satsNumber}
							maxlength="9"
							placeholder="SATS Number (9 digits)"
							class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
						/>
					</div>

					<div class="relative">
						<Calendar size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							type="date"
							bind:value={dateOfBirth}
							placeholder="Date of Birth"
							class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
						/>
					</div>

					{#if error}
						<div class="text-sm text-danger-600 bg-danger-50 rounded-lg p-3">{error}</div>
					{/if}

					<button
						type="submit"
						disabled={loading || !satsNumber || !dateOfBirth}
						class="w-full py-2 px-4 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
					>
						{loading ? 'Signing in...' : 'Sign in'}
					</button>
				</form>
			{:else if activeTab === 'staff'}
				<form onsubmit={handleStaffSubmit} class="space-y-4">
					<div class="relative">
						<Phone size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							type="tel"
							bind:value={mobile}
							maxlength="10"
							placeholder="10-digit mobile number"
							class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
						/>
					</div>

					<div class="relative">
						<Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							type="password"
							bind:value={staffPassword}
							placeholder="Password"
							class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
						/>
					</div>

					{#if error}
						<div class="text-sm text-danger-600 bg-danger-50 rounded-lg p-3">{error}</div>
					{/if}

					<button
						type="submit"
						disabled={loading || !mobile || !staffPassword}
						class="w-full py-2 px-4 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
					>
						{loading ? 'Signing in...' : 'Sign in'}
					</button>
				</form>
			{:else}
				<form onsubmit={handleAdminSubmit} class="space-y-4">
					<div class="relative">
						<Shield size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							type="email"
							bind:value={adminEmail}
							placeholder="Email address"
							class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
						/>
					</div>

					<div class="relative">
						<Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							type="password"
							bind:value={adminPassword}
							placeholder="Password"
							class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
						/>
					</div>

					{#if error}
						<div class="text-sm text-danger-600 bg-danger-50 rounded-lg p-3">{error}</div>
					{/if}

					<button
						type="submit"
						disabled={loading || !adminEmail || !adminPassword}
						class="w-full py-2 px-4 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors"
					>
						{loading ? 'Signing in...' : 'Sign in'}
					</button>
				</form>
			{/if}
		</div>

		<div class="text-center mt-6 flex items-center justify-center gap-4">
			<a href="/timetable" class="text-xs text-slate-400 hover:text-primary-600 transition-colors underline underline-offset-2">
				View Master School Timetable 2026-27
			</a>
			<span class="text-slate-300">·</span>
			<a href="/play" class="text-xs font-bold text-primary-500 hover:text-primary-700 transition-colors underline underline-offset-2">
			 🎮 Play Quiz Arena
			</a>
		</div>
	</div>
</div>
