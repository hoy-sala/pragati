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

<div class="min-h-screen flex flex-col overflow-hidden relative">
	<div class="hero-bg"></div>
	<div class="hero-orb hero-orb-1"></div>
	<div class="hero-orb hero-orb-2"></div>
	<div class="hero-orb hero-orb-3"></div>

	{#if view === 'home'}
		<!-- Hero -->
		<div class="relative z-10 pt-8 sm:pt-10 pb-4 px-4 text-center fade-in">
			<div class="flex justify-center gap-2.5 mb-5">
				<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/10 backdrop-blur border border-white/15 text-white/60 text-xs font-medium">
					<span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
					NEP 2020
				</span>
				<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/10 backdrop-blur border border-white/15 text-white/60 text-xs font-medium">
					<span class="w-1.5 h-1.5 rounded-full bg-amber-400 animate-pulse"></span>
					Holistic Development
				</span>
			</div>

			<div class="logo-glow inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-white/20 to-white/5 backdrop-blur-xl border border-white/20 shadow-2xl mb-4">
				<GraduationCap size={32} class="text-white" />
			</div>
			<h1 class="text-4xl sm:text-5xl font-black text-white tracking-tight font-kannada">ಪ್ರಗತಿ</h1>
			<div class="flex items-center justify-center gap-2.5 mt-1.5">
				<span class="h-px w-6 bg-white/20"></span>
				<span class="text-xs font-bold text-white/50 uppercase tracking-[0.3em]">Pragati</span>
				<span class="h-px w-6 bg-white/20"></span>
			</div>
			<p class="mt-3 text-base sm:text-lg text-white/60 max-w-sm mx-auto">Empowering every student and teacher to reach their full potential</p>
		</div>

		<!-- Cards -->
		<div class="relative z-10 flex-1 flex items-start justify-center px-4 pb-6">
			<div class="w-full max-w-5xl">
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">

					<!-- Quiz Arena Card -->
					<a href="/play" class="card-quiz group relative rounded-2xl p-5 text-white shadow-2xl shadow-purple-500/20 hover:shadow-purple-500/40 hover:scale-[1.03] active:scale-[0.98] transition-all duration-300 overflow-hidden">
						<!-- Decorative visual elements -->
						<div class="absolute -top-6 -right-6 w-24 h-24 rounded-full bg-white/5"></div>
						<div class="absolute -bottom-8 -left-8 w-32 h-32 rounded-full bg-white/5"></div>
						<div class="absolute top-3 right-3 flex gap-1">
							<span class="w-2 h-2 rounded-full bg-yellow-300/60 animate-bounce" style="animation-delay:0s"></span>
							<span class="w-1.5 h-1.5 rounded-full bg-pink-300/50 animate-bounce" style="animation-delay:0.15s"></span>
							<span class="w-1 h-1 rounded-full bg-cyan-300/40 animate-bounce" style="animation-delay:0.3s"></span>
						</div>
						<!-- Visual: Growth chart icon -->
						<div class="mb-4 flex items-end gap-1.5">
							<div class="w-3 h-5 rounded-sm bg-white/15 bar-bounce" style="animation-delay:0s"></div>
							<div class="w-3 h-8 rounded-sm bg-white/20 bar-bounce" style="animation-delay:0.1s"></div>
							<div class="w-3 h-11 rounded-sm bg-white/25 bar-bounce" style="animation-delay:0.2s"></div>
							<div class="w-3 h-14 rounded-sm bg-white/30 bar-bounce" style="animation-delay:0.3s"></div>
							<div class="ml-1 w-5 h-5 rounded-full bg-white/20 flex items-center justify-center">
								<Gamepad2 size={12} class="text-white" />
							</div>
						</div>
						<h2 class="text-xl font-bold relative">Quiz Arena</h2>
						<p class="text-white/50 text-sm mt-1.5 leading-relaxed relative">Test knowledge across every subject with interactive quizzes and climb the leaderboard.</p>
						<div class="flex items-center gap-3 mt-4 relative">
							<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-white/10 text-[11px] text-white/60 font-medium">Any class</span>
							<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-white/10 text-[11px] text-white/60 font-medium">Scores</span>
							<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-white/10 text-[11px] text-white/60 font-medium">Rank</span>
						</div>
						<div class="absolute bottom-4 right-4 w-8 h-8 rounded-lg bg-white/10 flex items-center justify-center group-hover:bg-white/20 transition-colors">
							<span class="text-white/60 group-hover:text-white transition-colors">→</span>
						</div>
					</a>

					<!-- Timetable Card -->
					<a href="/timetable" class="card-time group relative rounded-2xl p-5 text-white shadow-xl hover:shadow-2xl hover:scale-[1.03] active:scale-[0.98] transition-all duration-300 overflow-hidden border border-white/10">
						<div class="absolute -top-6 -right-6 w-24 h-24 rounded-full bg-white/5"></div>
						<div class="absolute -bottom-8 -left-8 w-32 h-32 rounded-full bg-white/5"></div>
						<!-- Visual: Timetable grid -->
						<div class="mb-4 grid grid-cols-5 gap-1 w-fit">
							{#each Array(15) as _, i}
								<div class="w-3.5 h-3 rounded-sm bar-bounce {i % 3 === 0 ? 'bg-white/25' : 'bg-white/10'}" style="animation-delay:{i * 0.03}s"></div>
							{/each}
							<div class="col-span-5 mt-0.5">
								<div class="h-1.5 rounded-full bg-white/10 overflow-hidden w-full">
									<div class="h-full rounded-full bg-sky-300/40" style="width:65%"></div>
								</div>
							</div>
						</div>
						<h2 class="text-xl font-bold relative">Time Table</h2>
						<p class="text-white/50 text-sm mt-1.5 leading-relaxed relative">View the complete weekly schedule — class-wise and subject-wise.</p>
						<div class="flex items-center gap-3 mt-4 relative">
							<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-white/10 text-[11px] text-white/60 font-medium">All classes</span>
							<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-white/10 text-[11px] text-white/60 font-medium">All subjects</span>
						</div>
						<div class="absolute bottom-4 right-4 w-8 h-8 rounded-lg bg-white/10 flex items-center justify-center group-hover:bg-white/20 transition-colors">
							<span class="text-white/60 group-hover:text-white transition-colors">→</span>
						</div>
					</a>

					<!-- Sign In Card -->
					<button onclick={() => view = 'login'} class="card-sign group relative rounded-2xl p-5 text-left bg-white/95 backdrop-blur border border-white/30 shadow-xl hover:shadow-2xl hover:scale-[1.03] active:scale-[0.98] transition-all duration-300 overflow-hidden">
						<div class="absolute -top-6 -right-6 w-24 h-24 rounded-full bg-primary-50"></div>
						<div class="absolute -bottom-8 -left-8 w-32 h-32 rounded-full bg-primary-50"></div>
						<!-- Visual: User avatars / progress ring -->
						<div class="mb-4 flex items-center gap-2">
							<div class="flex -space-x-2">
								<div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 border-2 border-white flex items-center justify-center text-[10px] font-bold text-white shadow-sm">S</div>
								<div class="w-8 h-8 rounded-full bg-gradient-to-br from-amber-400 to-orange-500 border-2 border-white flex items-center justify-center text-[10px] font-bold text-white shadow-sm">T</div>
								<div class="w-8 h-8 rounded-full bg-gradient-to-br from-purple-400 to-purple-600 border-2 border-white flex items-center justify-center text-[10px] font-bold text-white shadow-sm">A</div>
							</div>
							<div class="w-10 h-10 rounded-full border-[3px] border-primary-200 flex items-center justify-center relative">
								<svg class="absolute inset-0 w-full h-full -rotate-90" viewBox="0 0 36 36">
									<circle cx="18" cy="18" r="15" fill="none" stroke="currentColor" stroke-width="3" class="text-primary-100" />
									<circle cx="18" cy="18" r="15" fill="none" stroke="currentColor" stroke-width="3" stroke-dasharray="70 100" stroke-linecap="round" class="text-primary-500 progress-ring" />
								</svg>
								<span class="text-[10px] font-bold text-primary-600 relative z-10">✓</span>
							</div>
						</div>
						<h2 class="text-xl font-bold text-slate-800 relative">Sign In</h2>
						<p class="text-slate-400 text-sm mt-1.5 leading-relaxed relative">Access your personalized dashboard — marks, assessments, reports & more.</p>
						<div class="flex items-center gap-2 mt-4 relative">
							<span class="px-2 py-0.5 rounded-md bg-blue-50 text-blue-600 text-[11px] font-medium">Student</span>
							<span class="px-2 py-0.5 rounded-md bg-amber-50 text-amber-600 text-[11px] font-medium">Teacher</span>
							<span class="px-2 py-0.5 rounded-md bg-purple-50 text-purple-600 text-[11px] font-medium">Admin</span>
						</div>
						<div class="absolute bottom-4 right-4 w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center group-hover:bg-primary-50 text-slate-400 group-hover:text-primary-500 transition-colors">
							<span>→</span>
						</div>
					</button>

				</div>
			</div>
		</div>

		<!-- Footer -->
		<div class="relative z-10 text-center pb-5 px-4">
			<div class="flex items-center justify-center gap-5 text-xs text-white/20">
				<span>Classes 6–10</span>
				<span class="w-1 h-1 rounded-full bg-white/15"></span>
				<span>CCE Pattern</span>
				<span class="w-1 h-1 rounded-full bg-white/15"></span>
				<span>HPC Reports</span>
			</div>
		</div>

	{:else}
		<!-- Login View -->
		<div class="relative z-10 flex-1 flex items-center justify-center px-4 py-8">
			<div class="w-full max-w-sm fade-in">
				<button onclick={() => { view = 'home'; error = ''; }}
					class="flex items-center gap-1.5 text-sm text-white/50 hover:text-white/80 mb-6 transition-colors">
					<ArrowLeft size={16} /> Back to home
				</button>
				<div class="text-center mb-6">
					<div class="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-white/10 backdrop-blur-xl border border-white/15 shadow-xl mb-4">
						<GraduationCap size={28} class="text-white" />
					</div>
					<h1 class="text-3xl font-bold text-white">Welcome Back</h1>
					<p class="text-sm text-white/50 mt-2">Sign in to access your dashboard</p>
				</div>
				<div class="bg-white/95 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/20 p-5 space-y-4">
					<div class="flex border border-slate-200 rounded-xl overflow-hidden">
						<button onclick={() => { activeTab = 'student'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors {activeTab === 'student' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">Student</button>
						<button onclick={() => { activeTab = 'staff'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors {activeTab === 'staff' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">Staff</button>
						<button onclick={() => { activeTab = 'admin'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors {activeTab === 'admin' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">Admin</button>
					</div>

					{#if activeTab === 'student'}
						<form onsubmit={handleStudentSubmit} class="space-y-3">
							<div class="relative"><Hash size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="text" bind:value={satsNumber} maxlength="9" placeholder="SATS Number (9 digits)" class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							<div class="relative"><Calendar size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="date" bind:value={dateOfBirth} class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							{#if error}<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>{/if}
							<button type="submit" disabled={loading || !satsNumber || !dateOfBirth} class="w-full py-2.5 bg-primary-600 text-white rounded-xl text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">{loading ? 'Signing in...' : 'Sign in'}</button>
						</form>
					{:else if activeTab === 'staff'}
						<form onsubmit={handleStaffSubmit} class="space-y-3">
							<div class="relative"><Phone size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="tel" bind:value={mobile} maxlength="10" placeholder="10-digit mobile number" class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							<div class="relative"><Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="password" bind:value={staffPassword} placeholder="Password" class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							{#if error}<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>{/if}
							<button type="submit" disabled={loading || !mobile || !staffPassword} class="w-full py-2.5 bg-primary-600 text-white rounded-xl text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">{loading ? 'Signing in...' : 'Sign in'}</button>
						</form>
					{:else}
						<form onsubmit={handleAdminSubmit} class="space-y-3">
							<div class="relative"><Shield size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="email" bind:value={adminEmail} placeholder="Email address" class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							<div class="relative"><Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="password" bind:value={adminPassword} placeholder="Password" class="w-full pl-9 pr-3 py-2.5 rounded-xl border border-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							{#if error}<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>{/if}
							<button type="submit" disabled={loading || !adminEmail || !adminPassword} class="w-full py-2.5 bg-primary-600 text-white rounded-xl text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">{loading ? 'Signing in...' : 'Sign in'}</button>
						</form>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.hero-bg { position:fixed; inset:0; z-index:0; background:linear-gradient(135deg,#0f172a 0%,#1e1b4b 30%,#312e81 50%,#1e3a5f 70%,#0f172a 100%); }
	.hero-orb { position:fixed; border-radius:50%; z-index:0; pointer-events:none; filter:blur(80px); }
	.hero-orb-1 { width:400px; height:400px; top:-100px; right:-100px; background:radial-gradient(circle,rgba(139,92,246,0.35) 0%,transparent 70%); animation:orbFloat 8s ease-in-out infinite; }
	.hero-orb-2 { width:300px; height:300px; bottom:0; left:-80px; background:radial-gradient(circle,rgba(59,130,246,0.3) 0%,transparent 70%); animation:orbFloat 10s ease-in-out infinite reverse; }
	.hero-orb-3 { width:200px; height:200px; top:40%; left:60%; background:radial-gradient(circle,rgba(236,72,153,0.2) 0%,transparent 70%); animation:orbFloat 12s ease-in-out infinite 2s; }
	@keyframes orbFloat { 0%,100%{transform:translate(0,0) scale(1)} 33%{transform:translate(30px,-20px) scale(1.05)} 66%{transform:translate(-20px,15px) scale(0.95)} }
	.logo-glow { box-shadow:0 0 40px rgba(139,92,246,0.3),0 0 80px rgba(139,92,246,0.1); }
	.card-quiz { background:linear-gradient(135deg,#7c3aed 0%,#a855f7 40%,#c026d3 100%); }
	.card-time { background:linear-gradient(135deg,#0ea5e9 0%,#06b6d4 50%,#14b8a6 100%); }
	.bar-bounce { animation:barGrow 0.6s ease-out both; }
	@keyframes barGrow { from{transform:scaleY(0);transform-origin:bottom} to{transform:scaleY(1);transform-origin:bottom} }
	.progress-ring { animation:ringDraw 1.5s ease-out both; }
	@keyframes ringDraw { from{stroke-dashoffset:100} to{stroke-dashoffset:30} }
	.fade-in { animation:fadeUp 0.5s ease-out; }
	@keyframes fadeUp { from{opacity:0;transform:translateY(16px)} to{opacity:1;transform:translateY(0)} }
</style>
