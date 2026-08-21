<script lang="ts">
	import { login, staffLogin, studentLogin } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { GraduationCap, Phone, Key, Hash, Calendar, Shield, Gamepad2, Clock, LogIn, ArrowLeft, BookOpen, TrendingUp, Users, Award, Sparkles, ChevronRight } from 'lucide-svelte';

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
	<!-- Background -->
	<div class="hero-bg"></div>
	<div class="hero-orb hero-orb-1"></div>
	<div class="hero-orb hero-orb-2"></div>
	<div class="hero-orb hero-orb-3"></div>

	{#if view === 'home'}
		<!-- Hero Section -->
		<div class="relative z-10 pt-10 sm:pt-14 pb-6 px-4 text-center fade-in">
			<!-- Floating badges -->
			<div class="flex justify-center gap-3 mb-6">
				<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/10 backdrop-blur border border-white/15 text-white/70 text-xs font-medium">
					<TrendingUp size={12} /> NEP 2020 Aligned
				</span>
				<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/10 backdrop-blur border border-white/15 text-white/70 text-xs font-medium">
					<Sparkles size={12} /> Holistic Growth
				</span>
			</div>

			<!-- Logo -->
			<div class="logo-glow inline-flex items-center justify-center w-18 h-18 rounded-[1.25rem] bg-gradient-to-br from-white/20 to-white/5 backdrop-blur-xl border border-white/20 shadow-2xl mb-5">
				<GraduationCap size={36} class="text-white" />
			</div>

			<!-- Title -->
			<h1 class="text-5xl sm:text-6xl font-black text-white tracking-tight font-kannada">ಪ್ರಗತಿ</h1>
			<div class="flex items-center justify-center gap-2.5 mt-2">
				<span class="h-px w-8 bg-white/20"></span>
				<span class="text-sm font-bold text-white/60 uppercase tracking-[0.3em]">Pragati</span>
				<span class="h-px w-8 bg-white/20"></span>
			</div>

			<!-- Tagline -->
			<p class="mt-5 text-lg sm:text-xl text-white/80 font-light max-w-md mx-auto leading-relaxed">
				Empowering every <span class="text-white font-medium">student</span> and <span class="text-white font-medium">teacher</span> to reach their full potential
			</p>
			<p class="mt-2 text-sm text-white/40">A comprehensive school management platform for holistic development</p>
		</div>

		<!-- Feature Cards -->
		<div class="relative z-10 flex-1 flex items-start justify-center px-4 pb-8">
			<div class="w-full max-w-lg space-y-3.5">

				<!-- Quiz Arena -->
				<a href="/play"
					class="quiz-card group block rounded-2xl p-5 text-white shadow-2xl shadow-purple-500/20 hover:shadow-purple-500/40 hover:scale-[1.02] active:scale-[0.98] transition-all duration-300">
					<div class="flex items-start gap-4">
						<div class="w-14 h-14 rounded-2xl bg-white/15 backdrop-blur flex items-center justify-center shrink-0 border border-white/10 group-hover:scale-110 transition-transform">
							<Gamepad2 size={26} />
						</div>
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2">
								<h2 class="text-xl font-bold">Quiz Arena</h2>
								<span class="px-2 py-0.5 rounded-full bg-white/15 text-[10px] font-bold uppercase tracking-wider">Play</span>
							</div>
							<p class="text-white/60 text-sm mt-1 leading-relaxed">Challenge yourself with interactive quizzes across every subject. No login needed — just pick a topic and play!</p>
							<div class="flex items-center gap-4 mt-3">
								<span class="flex items-center gap-1 text-xs text-white/50"><BookOpen size={11} /> Any class</span>
								<span class="flex items-center gap-1 text-xs text-white/50"><TrendingUp size={11} /> Track scores</span>
								<span class="flex items-center gap-1 text-xs text-white/50"><Award size={11} /> Leaderboard</span>
							</div>
						</div>
						<div class="w-9 h-9 rounded-xl bg-white/10 flex items-center justify-center shrink-0 mt-1 group-hover:bg-white/20 transition-colors">
							<ChevronRight size={18} />
						</div>
					</div>
				</a>

				<!-- Timetable -->
				<a href="/timetable"
					class="group block rounded-2xl p-5 bg-white/10 backdrop-blur-xl border border-white/15 text-white shadow-xl hover:bg-white/15 hover:scale-[1.02] active:scale-[0.98] transition-all duration-300">
					<div class="flex items-start gap-4">
						<div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-sky-400/30 to-cyan-400/30 backdrop-blur flex items-center justify-center shrink-0 border border-white/10 group-hover:scale-110 transition-transform">
							<Clock size={26} />
						</div>
						<div class="flex-1 min-w-0">
							<h2 class="text-xl font-bold">Time Table</h2>
							<p class="text-white/50 text-sm mt-1 leading-relaxed">View the complete weekly schedule — class-wise and subject-wise. Plan your study routine effectively.</p>
							<div class="flex items-center gap-4 mt-3">
								<span class="flex items-center gap-1 text-xs text-white/50"><Users size={11} /> All classes</span>
								<span class="flex items-center gap-1 text-xs text-white/50"><BookOpen size={11} /> All subjects</span>
							</div>
						</div>
						<div class="w-9 h-9 rounded-xl bg-white/10 flex items-center justify-center shrink-0 mt-1 group-hover:bg-white/20 transition-colors">
							<ChevronRight size={18} />
						</div>
					</div>
				</a>

				<!-- Sign In -->
				<button onclick={() => view = 'login'}
					class="group w-full text-left rounded-2xl p-5 bg-white/95 backdrop-blur border border-slate-200/50 shadow-xl hover:shadow-2xl hover:scale-[1.02] active:scale-[0.98] transition-all duration-300">
					<div class="flex items-start gap-4">
						<div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shrink-0 shadow-lg shadow-primary-500/30 group-hover:scale-110 transition-transform">
							<LogIn size={26} class="text-white" />
						</div>
						<div class="flex-1 min-w-0">
							<h2 class="text-xl font-bold text-slate-800">Sign In</h2>
							<p class="text-slate-400 text-sm mt-1 leading-relaxed">Access your personalized dashboard — marks, assessments, HPC cards, mentor tracking & more.</p>
							<div class="flex items-center gap-3 mt-3">
								<span class="px-2.5 py-1 rounded-lg bg-blue-50 text-blue-600 text-xs font-medium">Student</span>
								<span class="px-2.5 py-1 rounded-lg bg-amber-50 text-amber-600 text-xs font-medium">Teacher</span>
								<span class="px-2.5 py-1 rounded-lg bg-purple-50 text-purple-600 text-xs font-medium">Admin</span>
							</div>
						</div>
						<div class="w-9 h-9 rounded-xl bg-slate-100 flex items-center justify-center shrink-0 mt-1 group-hover:bg-primary-50 group-hover:text-primary-600 text-slate-400 transition-colors">
							<ChevronRight size={18} />
						</div>
					</div>
				</button>
			</div>
		</div>

		<!-- Footer -->
		<div class="relative z-10 text-center pb-6 px-4">
			<div class="flex items-center justify-center gap-6 text-xs text-white/25">
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
							class="flex-1 py-2.5 text-sm font-medium transition-colors {activeTab === 'student' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">
							Student
						</button>
						<button onclick={() => { activeTab = 'staff'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors {activeTab === 'staff' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">
							Staff
						</button>
						<button onclick={() => { activeTab = 'admin'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-medium transition-colors {activeTab === 'admin' ? 'bg-primary-600 text-white' : 'bg-white text-slate-500 hover:bg-slate-50'}">
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
							{#if error}<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>{/if}
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
							{#if error}<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>{/if}
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
							{#if error}<div class="text-sm text-danger-600 bg-danger-50 rounded-xl p-3">{error}</div>{/if}
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
	.hero-bg {
		position: fixed; inset: 0; z-index: 0;
		background: linear-gradient(135deg, #0f172a 0%, #1e1b4b 30%, #312e81 50%, #1e3a5f 70%, #0f172a 100%);
	}
	.hero-orb {
		position: fixed; border-radius: 50%; z-index: 0; pointer-events: none; filter: blur(80px);
	}
	.hero-orb-1 {
		width: 400px; height: 400px; top: -100px; right: -100px;
		background: radial-gradient(circle, rgba(139,92,246,0.35) 0%, transparent 70%);
		animation: orbFloat 8s ease-in-out infinite;
	}
	.hero-orb-2 {
		width: 300px; height: 300px; bottom: 0; left: -80px;
		background: radial-gradient(circle, rgba(59,130,246,0.3) 0%, transparent 70%);
		animation: orbFloat 10s ease-in-out infinite reverse;
	}
	.hero-orb-3 {
		width: 200px; height: 200px; top: 40%; left: 60%;
		background: radial-gradient(circle, rgba(236,72,153,0.2) 0%, transparent 70%);
		animation: orbFloat 12s ease-in-out infinite 2s;
	}
	@keyframes orbFloat {
		0%, 100% { transform: translate(0, 0) scale(1); }
		33% { transform: translate(30px, -20px) scale(1.05); }
		66% { transform: translate(-20px, 15px) scale(0.95); }
	}
	.logo-glow {
		box-shadow: 0 0 40px rgba(139,92,246,0.3), 0 0 80px rgba(139,92,246,0.1);
	}
	.quiz-card {
		background: linear-gradient(135deg, #7c3aed 0%, #a855f7 40%, #c026d3 100%);
	}
	.fade-in { animation: fadeUp 0.5s ease-out; }
	@keyframes fadeUp { from { opacity: 0; transform: translateY(16px); } to { opacity: 1; transform: translateY(0); } }
</style>
