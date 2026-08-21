<script lang="ts">
	import { login, staffLogin, studentLogin } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { GraduationCap, Phone, Key, Hash, Calendar, Shield, Gamepad2, Clock, LogIn, ArrowLeft } from 'lucide-svelte';

	type View = 'home' | 'login';
	let view = $state<View>('home');

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
		error = ''; loading = true;
		let result;
		try { result = await studentLogin(satsNumber, dateOfBirth); }
		catch { result = { error: 'Unable to reach server.' }; }
		loading = false;
		if ('error' in result) error = result.error; else goto('/dashboard');
	}

	async function handleStaffSubmit() {
		error = ''; loading = true;
		let result;
		try { result = await staffLogin(mobile, staffPassword); }
		catch { result = { error: 'Unable to reach server.' }; }
		loading = false;
		if ('error' in result) error = result.error; else goto('/dashboard');
	}

	async function handleAdminSubmit() {
		error = ''; loading = true;
		let result;
		try { result = await login(adminEmail, adminPassword); }
		catch { result = { error: 'Unable to reach server.' }; }
		loading = false;
		if ('error' in result) error = result.error; else goto('/dashboard');
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50 flex flex-col">
	{#if view === 'home'}
		<!-- Hero -->
		<div class="pt-12 pb-8 px-4 text-center fade-in">
			<div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-700 shadow-lg shadow-primary-200/60 mb-4">
				<GraduationCap size={32} class="text-white" />
			</div>
			<h1 class="text-4xl font-bold text-slate-800 font-kannada tracking-wide">ಪ್ರಗತಿ</h1>
			<div class="flex items-center justify-center gap-2 mt-1.5">
				<span class="text-sm font-semibold text-primary-600 uppercase tracking-[0.25em]">PRAGATI</span>
				<span class="w-1 h-1 rounded-full bg-primary-300"></span>
				<span class="text-xs text-slate-400">v1.0</span>
			</div>
			<p class="mt-3 text-sm text-slate-400 italic">Every Child Can Progress</p>
		</div>

		<!-- Feature Cards -->
		<div class="flex-1 flex items-start justify-center px-4 pb-12">
			<div class="w-full max-w-md space-y-4">
				<!-- Quiz Arena -->
				<a href="/play"
					class="group block rounded-2xl p-5 bg-gradient-to-br from-violet-500 via-purple-500 to-fuchsia-500 text-white shadow-lg shadow-purple-200/60 hover:shadow-purple-300/80 hover:scale-[1.02] active:scale-[0.98] transition-all">
					<div class="flex items-center gap-4">
						<div class="w-12 h-12 rounded-xl bg-white/20 flex items-center justify-center shrink-0">
							<Gamepad2 size={24} />
						</div>
						<div class="flex-1">
							<h2 class="text-lg font-bold">Quiz Arena</h2>
							<p class="text-white/70 text-sm mt-0.5">Play & learn — no login needed</p>
						</div>
						<span class="text-2xl opacity-60 group-hover:opacity-100 transition-opacity">→</span>
					</div>
				</a>

				<!-- Timetable -->
				<a href="/timetable"
					class="group block rounded-2xl p-5 bg-white border border-slate-200 shadow-sm hover:shadow-md hover:border-primary-200 hover:scale-[1.02] active:scale-[0.98] transition-all">
					<div class="flex items-center gap-4">
						<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center shrink-0 shadow-sm">
							<Clock size={24} class="text-white" />
						</div>
						<div class="flex-1">
							<h2 class="text-lg font-bold text-slate-800">Time Table</h2>
							<p class="text-slate-400 text-sm mt-0.5">View class-wise & subject-wise schedule</p>
						</div>
						<span class="text-slate-300 text-2xl group-hover:text-primary-400 transition-colors">→</span>
					</div>
				</a>

				<!-- Sign In -->
				<button onclick={() => view = 'login'}
					class="group w-full text-left rounded-2xl p-5 bg-white border border-slate-200 shadow-sm hover:shadow-md hover:border-primary-200 hover:scale-[1.02] active:scale-[0.98] transition-all">
					<div class="flex items-center gap-4">
						<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shrink-0 shadow-sm">
							<LogIn size={24} class="text-white" />
						</div>
						<div class="flex-1">
							<h2 class="text-lg font-bold text-slate-800">Sign In</h2>
							<p class="text-slate-400 text-sm mt-0.5">Students, teachers & administrators</p>
						</div>
						<span class="text-slate-300 text-2xl group-hover:text-primary-400 transition-colors">→</span>
					</div>
				</button>
			</div>
		</div>

	{:else}
		<!-- Login View -->
		<div class="flex-1 flex items-center justify-center px-4 py-8">
			<div class="w-full max-w-sm fade-in">
				<!-- Back button -->
				<button onclick={() => { view = 'home'; error = ''; }}
					class="flex items-center gap-1.5 text-sm text-slate-400 hover:text-slate-600 mb-6 transition-colors">
					<ArrowLeft size={16} /> Back
				</button>

				<!-- Header -->
				<div class="text-center mb-6">
					<div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 shadow-md shadow-primary-200/50 mb-3">
						<GraduationCap size={24} class="text-white" />
					</div>
					<h1 class="text-2xl font-bold text-slate-800">Sign In</h1>
					<p class="text-sm text-slate-400 mt-1">Choose your login type below</p>
				</div>

				<!-- Login Card -->
				<div class="bg-white rounded-2xl shadow-sm border border-slate-200 p-5 space-y-4">
					<!-- Tabs -->
					<div class="flex border border-slate-200 rounded-xl overflow-hidden">
						<button onclick={() => { activeTab = 'student'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors
								{activeTab === 'student' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">
							Student
						</button>
						<button onclick={() => { activeTab = 'staff'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors
								{activeTab === 'staff' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">
							Staff
						</button>
						<button onclick={() => { activeTab = 'admin'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors
								{activeTab === 'admin' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">
							Admin
						</button>
					</div>

					{#if activeTab === 'student'}
						<form onsubmit={handleStudentSubmit} class="space-y-3">
							<div class="relative">
								<Hash size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
								<input type="text" bind:value={satsNumber} maxlength="9" placeholder="SATS Number (9 digits)"
									class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-slate-50 focus:bg-white transition-colors" />
							</div>
							<div class="relative">
								<Calendar size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
								<input type="date" bind:value={dateOfBirth}
									class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-slate-50 focus:bg-white transition-colors" />
							</div>
							{#if error}
								<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>
							{/if}
							<button type="submit" disabled={loading || !satsNumber || !dateOfBirth}
								class="w-full py-2.5 px-4 bg-primary-600 text-white rounded-xl text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
								{loading ? 'Signing in...' : 'Sign in'}
							</button>
						</form>

					{:else if activeTab === 'staff'}
						<form onsubmit={handleStaffSubmit} class="space-y-3">
							<div class="relative">
								<Phone size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
								<input type="tel" bind:value={mobile} maxlength="10" placeholder="10-digit mobile number"
									class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-slate-50 focus:bg-white transition-colors" />
							</div>
							<div class="relative">
								<Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
								<input type="password" bind:value={staffPassword} placeholder="Password"
									class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-slate-50 focus:bg-white transition-colors" />
							</div>
							{#if error}
								<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>
							{/if}
							<button type="submit" disabled={loading || !mobile || !staffPassword}
								class="w-full py-2.5 px-4 bg-primary-600 text-white rounded-xl text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
								{loading ? 'Signing in...' : 'Sign in'}
							</button>
						</form>

					{:else}
						<form onsubmit={handleAdminSubmit} class="space-y-3">
							<div class="relative">
								<Shield size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
								<input type="email" bind:value={adminEmail} placeholder="Email address"
									class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-slate-50 focus:bg-white transition-colors" />
							</div>
							<div class="relative">
								<Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
								<input type="password" bind:value={adminPassword} placeholder="Password"
									class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-slate-50 focus:bg-white transition-colors" />
							</div>
							{#if error}
								<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>
							{/if}
							<button type="submit" disabled={loading || !adminEmail || !adminPassword}
								class="w-full py-2.5 px-4 bg-primary-600 text-white rounded-xl text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
								{loading ? 'Signing in...' : 'Sign in'}
							</button>
						</form>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.fade-in { animation: fadeUp 0.4s ease-out; }
	@keyframes fadeUp { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
