<script lang="ts">
	import { login, staffLogin, studentLogin } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { GraduationCap, Phone, Key, Hash, Calendar, Shield, Gamepad2, Clock, LogIn, ArrowLeft, Star, Zap, Trophy, BookOpen } from 'lucide-svelte';

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

<div class="min-h-screen flex flex-col overflow-hidden relative game-bg">
	<!-- Floating decorations -->
	<div class="deco deco-star1">⭐</div>
	<div class="deco deco-star2">🌟</div>
	<div class="deco deco-star3">✨</div>
	<div class="deco deco-cloud1">☁️</div>
	<div class="deco deco-cloud2">☁️</div>
	<div class="deco deco-coin1">🪙</div>
	<div class="deco deco-coin2">⭐</div>
	<div class="deco deco-heart1">💜</div>
	<div class="deco deco-zap">⚡</div>

	{#if view === 'home'}
		<!-- Hero -->
		<div class="relative z-10 pt-8 sm:pt-10 pb-3 px-4 text-center fade-in">
			<div class="bounce-in">
				<div class="inline-flex items-center justify-center w-20 h-20 rounded-[1.5rem] bg-gradient-to-br from-yellow-300 via-orange-400 to-pink-500 shadow-xl shadow-orange-400/40 mb-4 mascot">
					<span class="text-4xl">🏫</span>
				</div>
			</div>
			<h1 class="text-5xl sm:text-6xl font-black text-white tracking-tight drop-shadow-lg font-kannada" style="text-shadow: 3px 3px 0 rgba(0,0,0,0.15)">ಪ್ರಗತಿ</h1>
			<div class="flex items-center justify-center gap-2 mt-1">
				<span class="text-sm font-black text-white/80 uppercase tracking-widest">Pragati</span>
			</div>
			<p class="mt-2 text-sm text-white/70 font-medium">Learn. Play. Grow.</p>
		</div>

		<!-- Cards -->
		<div class="relative z-10 flex-1 flex items-start justify-center px-4 pb-6">
			<div class="w-full max-w-5xl">
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-4 items-stretch">

					<!-- Quiz Arena -->
					<a href="/play" class="game-card game-card-quiz group relative rounded-[1.5rem] p-5 text-white shadow-2xl hover:scale-[1.04] active:scale-[0.97] transition-all duration-200 overflow-hidden">
						<!-- Sparkle decorations -->
						<div class="absolute top-3 right-3 text-lg animate-spin-slow">⭐</div>
						<div class="absolute bottom-4 left-3 text-sm animate-bounce">🎯</div>
						<div class="absolute top-12 right-10 text-xs animate-pulse">🔥</div>

						<!-- Game icon -->
						<div class="w-16 h-16 rounded-2xl bg-white/20 backdrop-blur flex items-center justify-center mb-3 shadow-lg group-hover:rotate-12 transition-transform duration-300">
							<span class="text-3xl">🎮</span>
						</div>

						<h2 class="text-2xl font-black relative drop-shadow-sm">Quiz Arena</h2>
						<p class="text-white/70 text-sm mt-1.5 leading-relaxed relative font-medium">Battle through questions, earn points & climb the leaderboard!</p>

						<!-- Stats bar -->
						<div class="flex items-center gap-3 mt-4 relative">
							<div class="flex items-center gap-1 px-2.5 py-1 rounded-full bg-white/15 text-xs font-bold">
								<Star size={12} /> Points
							</div>
							<div class="flex items-center gap-1 px-2.5 py-1 rounded-full bg-white/15 text-xs font-bold">
								<Trophy size={12} /> Rank
							</div>
							<div class="flex items-center gap-1 px-2.5 py-1 rounded-full bg-white/15 text-xs font-bold">
								<Zap size={12} /> Streak
							</div>
						</div>

						<!-- Play button -->
						<div class="mt-4 relative">
							<div class="w-full py-2.5 rounded-xl bg-white/20 backdrop-blur text-center font-black text-sm tracking-wide group-hover:bg-white/30 transition-colors flex items-center justify-center gap-2">
								PLAY NOW <span class="text-lg">🚀</span>
							</div>
						</div>
					</a>

					<!-- Timetable -->
					<a href="/timetable" class="game-card game-card-time group relative rounded-[1.5rem] p-5 text-white shadow-2xl hover:scale-[1.04] active:scale-[0.97] transition-all duration-200 overflow-hidden">
						<div class="absolute top-3 right-3 text-lg animate-bounce" style="animation-delay:0.3s">📅</div>
						<div class="absolute bottom-4 left-3 text-sm animate-pulse">📚</div>
						<div class="absolute top-12 right-10 text-xs animate-spin-slow">🕐</div>

						<!-- Game icon -->
						<div class="w-16 h-16 rounded-2xl bg-white/20 backdrop-blur flex items-center justify-center mb-3 shadow-lg group-hover:-rotate-12 transition-transform duration-300">
							<span class="text-3xl">🗓️</span>
						</div>

						<h2 class="text-2xl font-black relative drop-shadow-sm">Timetable</h2>
						<p class="text-white/70 text-sm mt-1.5 leading-relaxed relative font-medium">See your weekly schedule — know what's next!</p>

						<!-- Mini schedule preview -->
						<div class="mt-4 space-y-1.5 relative">
							<div class="flex items-center gap-2">
								<div class="w-2 h-2 rounded-full bg-yellow-300"></div>
								<div class="h-2 flex-1 rounded-full bg-white/15 overflow-hidden"><div class="h-full bg-white/30 rounded-full" style="width:80%"></div></div>
								<span class="text-[10px] text-white/50 font-bold">8AM</span>
							</div>
							<div class="flex items-center gap-2">
								<div class="w-2 h-2 rounded-full bg-green-300"></div>
								<div class="h-2 flex-1 rounded-full bg-white/15 overflow-hidden"><div class="h-full bg-white/30 rounded-full" style="width:60%"></div></div>
								<span class="text-[10px] text-white/50 font-bold">9AM</span>
							</div>
							<div class="flex items-center gap-2">
								<div class="w-2 h-2 rounded-full bg-pink-300"></div>
								<div class="h-2 flex-1 rounded-full bg-white/15 overflow-hidden"><div class="h-full bg-white/30 rounded-full" style="width:90%"></div></div>
								<span class="text-[10px] text-white/50 font-bold">10AM</span>
							</div>
						</div>

						<div class="mt-4 relative">
							<div class="w-full py-2.5 rounded-xl bg-white/20 backdrop-blur text-center font-black text-sm tracking-wide group-hover:bg-white/30 transition-colors flex items-center justify-center gap-2">
								VIEW SCHEDULE <span class="text-lg">📋</span>
							</div>
						</div>
					</a>

					<!-- Sign In -->
					<button onclick={() => view = 'login'} class="game-card game-card-sign group relative rounded-[1.5rem] p-5 text-left bg-white shadow-2xl hover:scale-[1.04] active:scale-[0.97] transition-all duration-200 overflow-hidden border-4 border-white/50">
						<div class="absolute top-3 right-3 text-lg animate-bounce" style="animation-delay:0.5s">👋</div>
						<div class="absolute bottom-4 left-3 text-sm animate-pulse">🎓</div>
						<div class="absolute top-12 right-10 text-xs animate-spin-slow">🌟</div>

						<!-- Character avatars -->
						<div class="flex items-center gap-2 mb-3">
							<div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-xl shadow-lg shadow-blue-400/30 group-hover:animate-bounce">🧑‍🎓</div>
							<div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-amber-400 to-orange-500 flex items-center justify-center text-xl shadow-lg shadow-orange-400/30 group-hover:animate-bounce" style="animation-delay:0.1s">👩‍🏫</div>
							<div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center text-xl shadow-lg shadow-purple-400/30 group-hover:animate-bounce" style="animation-delay:0.2s">👨‍💼</div>
						</div>

						<h2 class="text-2xl font-black text-slate-800 relative">Sign In</h2>
						<p class="text-slate-400 text-sm mt-1.5 leading-relaxed relative font-medium">Jump back into your learning dashboard!</p>

						<!-- Role badges -->
						<div class="flex items-center gap-2 mt-4 relative">
							<span class="px-3 py-1 rounded-full bg-blue-100 text-blue-600 text-xs font-black">🧑‍🎓 Student</span>
							<span class="px-3 py-1 rounded-full bg-amber-100 text-amber-600 text-xs font-black">👩‍🏫 Teacher</span>
							<span class="px-3 py-1 rounded-full bg-purple-100 text-purple-600 text-xs font-black">👨‍💼 Admin</span>
						</div>

						<div class="mt-4 relative">
							<div class="w-full py-2.5 rounded-xl bg-gradient-to-r from-primary-500 to-primary-600 text-center font-black text-sm tracking-wide text-white group-hover:from-primary-600 group-hover:to-primary-700 transition-all flex items-center justify-center gap-2 shadow-lg shadow-primary-500/30">
								ENTER <span class="text-lg">🎯</span>
							</div>
						</div>
					</button>

				</div>
			</div>
		</div>

		<!-- Footer -->
		<div class="relative z-10 text-center pb-5 px-4">
			<div class="flex items-center justify-center gap-3 text-xs text-white/30 font-medium">
				<span class="px-2 py-0.5 rounded-full bg-white/5">📚 Classes 6–10</span>
				<span class="px-2 py-0.5 rounded-full bg-white/5">📝 CCE Pattern</span>
				<span class="px-2 py-0.5 rounded-full bg-white/5">🏆 HPC Reports</span>
			</div>
		</div>

	{:else}
		<!-- Login View -->
		<div class="relative z-10 flex-1 flex items-center justify-center px-4 py-8">
			<div class="w-full max-w-sm fade-in">
				<button onclick={() => { view = 'home'; error = ''; }}
					class="flex items-center gap-1.5 text-sm text-white/60 hover:text-white font-bold mb-6 transition-colors">
					<ArrowLeft size={16} /> Back
				</button>
				<div class="text-center mb-5">
					<div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-yellow-300 via-orange-400 to-pink-500 shadow-xl shadow-orange-400/30 mb-3 mascot">
						<span class="text-3xl">🏫</span>
					</div>
					<h1 class="text-3xl font-black text-white drop-shadow-lg">Welcome Back!</h1>
					<p class="text-sm text-white/60 mt-1.5 font-medium">Choose your character to enter</p>
				</div>
				<div class="bg-white/95 backdrop-blur-xl rounded-[1.5rem] shadow-2xl border-4 border-white/30 p-5 space-y-4">
					<div class="flex border-2 border-slate-200 rounded-xl overflow-hidden">
						<button onclick={() => { activeTab = 'student'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-black transition-all {activeTab === 'student' ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-lg shadow-blue-500/30' : 'bg-white text-slate-500 hover:bg-slate-50'}">🧑‍🎓 Student</button>
						<button onclick={() => { activeTab = 'staff'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-black transition-all {activeTab === 'staff' ? 'bg-gradient-to-r from-amber-500 to-orange-500 text-white shadow-lg shadow-amber-500/30' : 'bg-white text-slate-500 hover:bg-slate-50'}">👩‍🏫 Staff</button>
						<button onclick={() => { activeTab = 'admin'; error = ''; }}
							class="flex-1 py-2.5 text-sm font-black transition-all {activeTab === 'admin' ? 'bg-gradient-to-r from-purple-500 to-purple-600 text-white shadow-lg shadow-purple-500/30' : 'bg-white text-slate-500 hover:bg-slate-50'}">👨‍💼 Admin</button>
					</div>

					{#if activeTab === 'student'}
						<form onsubmit={handleStudentSubmit} class="space-y-3">
							<div class="relative"><Hash size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="text" bind:value={satsNumber} maxlength="9" placeholder="SATS Number (9 digits)" class="w-full pl-9 pr-3 py-2.5 rounded-xl border-2 border-slate-200 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							<div class="relative"><Calendar size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="date" bind:value={dateOfBirth} class="w-full pl-9 pr-3 py-2.5 rounded-xl border-2 border-slate-200 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							{#if error}<div class="text-sm font-bold text-red-600 bg-red-50 rounded-xl p-3 border-2 border-red-100">{error}</div>{/if}
							<button type="submit" disabled={loading || !satsNumber || !dateOfBirth} class="w-full py-3 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-xl text-sm font-black hover:from-blue-600 hover:to-blue-700 disabled:opacity-50 transition-all shadow-lg shadow-blue-500/30 flex items-center justify-center gap-2">{loading ? '⏳ Signing in...' : '🚀 Sign In'}</button>
						</form>
					{:else if activeTab === 'staff'}
						<form onsubmit={handleStaffSubmit} class="space-y-3">
							<div class="relative"><Phone size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="tel" bind:value={mobile} maxlength="10" placeholder="10-digit mobile number" class="w-full pl-9 pr-3 py-2.5 rounded-xl border-2 border-slate-200 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							<div class="relative"><Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="password" bind:value={staffPassword} placeholder="Password" class="w-full pl-9 pr-3 py-2.5 rounded-xl border-2 border-slate-200 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							{#if error}<div class="text-sm font-bold text-red-600 bg-red-50 rounded-xl p-3 border-2 border-red-100">{error}</div>{/if}
							<button type="submit" disabled={loading || !mobile || !staffPassword} class="w-full py-3 bg-gradient-to-r from-amber-500 to-orange-500 text-white rounded-xl text-sm font-black hover:from-amber-600 hover:to-orange-600 disabled:opacity-50 transition-all shadow-lg shadow-amber-500/30 flex items-center justify-center gap-2">{loading ? '⏳ Signing in...' : '🚀 Sign In'}</button>
						</form>
					{:else}
						<form onsubmit={handleAdminSubmit} class="space-y-3">
							<div class="relative"><Shield size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="email" bind:value={adminEmail} placeholder="Email address" class="w-full pl-9 pr-3 py-2.5 rounded-xl border-2 border-slate-200 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							<div class="relative"><Key size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" /><input type="password" bind:value={adminPassword} placeholder="Password" class="w-full pl-9 pr-3 py-2.5 rounded-xl border-2 border-slate-200 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 bg-slate-50 focus:bg-white transition-colors" /></div>
							{#if error}<div class="text-sm font-bold text-red-600 bg-red-50 rounded-xl p-3 border-2 border-red-100">{error}</div>{/if}
							<button type="submit" disabled={loading || !adminEmail || !adminPassword} class="w-full py-3 bg-gradient-to-r from-purple-500 to-purple-600 text-white rounded-xl text-sm font-black hover:from-purple-600 hover:to-purple-700 disabled:opacity-50 transition-all shadow-lg shadow-purple-500/30 flex items-center justify-center gap-2">{loading ? '⏳ Signing in...' : '🚀 Sign In'}</button>
						</form>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.game-bg {
		background: linear-gradient(160deg, #667eea 0%, #764ba2 25%, #f093fb 50%, #f5576c 75%, #fda085 100%);
		background-size: 400% 400%;
		animation: gradientShift 15s ease infinite;
	}
	@keyframes gradientShift {
		0% { background-position: 0% 50%; }
		50% { background-position: 100% 50%; }
		100% { background-position: 0% 50%; }
	}
	.game-card-quiz {
		background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
		border: 3px solid rgba(255,255,255,0.2);
	}
	.game-card-time {
		background: linear-gradient(135deg, #06b6d4 0%, #0ea5e9 50%, #3b82f6 100%);
		border: 3px solid rgba(255,255,255,0.2);
	}
	.game-card-sign {
		background: linear-gradient(135deg, #fefefe 0%, #f8fafc 100%);
		border: 3px solid rgba(99,102,241,0.2);
	}
	.mascot { animation: mascotFloat 3s ease-in-out infinite; }
	@keyframes mascotFloat {
		0%, 100% { transform: translateY(0) rotate(0deg); }
		25% { transform: translateY(-6px) rotate(-2deg); }
		75% { transform: translateY(3px) rotate(2deg); }
	}

	/* Floating decorations */
	.deco { position: fixed; z-index: 1; pointer-events: none; opacity: 0.5; }
	.deco-star1 { top: 8%; left: 8%; font-size: 1.5rem; animation: floatA 6s ease-in-out infinite; }
	.deco-star2 { top: 15%; right: 10%; font-size: 1.2rem; animation: floatB 7s ease-in-out infinite 1s; }
	.deco-star3 { bottom: 20%; left: 12%; font-size: 1rem; animation: floatA 5s ease-in-out infinite 2s; }
	.deco-cloud1 { top: 5%; right: 20%; font-size: 2rem; animation: driftRight 20s linear infinite; opacity: 0.3; }
	.deco-cloud2 { top: 25%; left: 5%; font-size: 1.5rem; animation: driftRight 25s linear infinite 5s; opacity: 0.25; }
	.deco-coin1 { bottom: 15%; right: 12%; font-size: 1.3rem; animation: floatB 4s ease-in-out infinite 0.5s; }
	.deco-coin2 { top: 40%; left: 5%; font-size: 1rem; animation: floatA 5.5s ease-in-out infinite 1.5s; }
	.deco-heart1 { top: 35%; right: 6%; font-size: 1.2rem; animation: floatB 6s ease-in-out infinite 2.5s; }
	.deco-zap { bottom: 30%; left: 8%; font-size: 1.1rem; animation: floatA 4.5s ease-in-out infinite 3s; }

	@keyframes floatA {
		0%, 100% { transform: translateY(0) translateX(0); }
		50% { transform: translateY(-15px) translateX(5px); }
	}
	@keyframes floatB {
		0%, 100% { transform: translateY(0) translateX(0); }
		50% { transform: translateY(10px) translateX(-8px); }
	}
	@keyframes driftRight {
		0% { transform: translateX(-100px); }
		100% { transform: translateX(calc(100vw + 100px)); }
	}

	:global(.animate-spin-slow) { animation: spin 4s linear infinite; }
	@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

	.bounce-in { animation: bounceIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55); }
	@keyframes bounceIn { 0% { transform: scale(0.3); opacity: 0; } 60% { transform: scale(1.1); } 100% { transform: scale(1); opacity: 1; } }
	.fade-in { animation: fadeUp 0.5s ease-out; }
	@keyframes fadeUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
</style>
