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
	<div class="starfield"></div>
	<!-- Floating decorations -->
	<div class="deco deco-moon">🌙</div>
	<div class="deco deco-planet">🪐</div>
	<div class="deco deco-rocket">🚀</div>
	<div class="deco deco-ufo">🛸</div>
	<div class="deco deco-star1">⭐</div>
	<div class="deco deco-star2">✨</div>
	<div class="deco deco-star3">🌟</div>
	<div class="deco deco-star4">✦</div>
	<div class="deco deco-star5">✧</div>
	<!-- Shooting stars / comets -->
	<div class="shooting-star ss1"></div>
	<div class="shooting-star ss2"></div>
	<div class="shooting-star ss3"></div>
	<div class="comet comet1"></div>
	<div class="comet comet2"></div>

	{#if view === 'home'}
		<!-- Hero -->
		<div class="relative z-10 pt-8 sm:pt-10 pb-3 px-4 text-center fade-in">
			<div class="bounce-in">
				<div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-yellow-300 via-orange-400 to-pink-500 shadow-xl shadow-orange-400/40 mb-3 mascot">
					<span class="text-3xl">🏫</span>
				</div>
			</div>
			<h1 class="game-title text-3xl sm:text-4xl font-black tracking-tight leading-none whitespace-nowrap">
				MDRS (SC-32) BAHADDURGHATTA
			</h1>
			<p class="mt-2 text-4xl sm:text-5xl font-black text-white/90 font-kannada" style="text-shadow: 2px 2px 0 rgba(0,0,0,0.15)">ಪ್ರಗತಿ</p>
			<div class="flex items-center justify-center gap-2 mt-1.5">
				<span class="h-px w-6 bg-white/20"></span>
				<span class="text-xs font-bold text-white/50 uppercase tracking-widest">Learn. Play. Grow.</span>
				<span class="h-px w-6 bg-white/20"></span>
			</div>
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
		background: linear-gradient(160deg, #0c0a2a 0%, #1a1150 25%, #2d1b69 40%, #1e1145 60%, #0f0d30 80%, #0a0820 100%);
		background-size: 400% 400%;
		animation: gradientShift 20s ease infinite;
	}
	@keyframes gradientShift {
		0% { background-position: 0% 50%; }
		50% { background-position: 100% 50%; }
		100% { background-position: 0% 50%; }
	}
	.starfield {
		position: fixed; inset: 0; z-index: 0; pointer-events: none;
		background-image:
			radial-gradient(1px 1px at 10% 20%, rgba(255,255,255,0.5) 0%, transparent 100%),
			radial-gradient(1px 1px at 25% 45%, rgba(255,255,255,0.4) 0%, transparent 100%),
			radial-gradient(1px 1px at 40% 10%, rgba(255,255,255,0.3) 0%, transparent 100%),
			radial-gradient(1px 1px at 55% 60%, rgba(255,255,255,0.5) 0%, transparent 100%),
			radial-gradient(1px 1px at 70% 30%, rgba(255,255,255,0.4) 0%, transparent 100%),
			radial-gradient(1px 1px at 85% 55%, rgba(255,255,255,0.3) 0%, transparent 100%),
			radial-gradient(1px 1px at 15% 75%, rgba(255,255,255,0.4) 0%, transparent 100%),
			radial-gradient(1px 1px at 50% 85%, rgba(255,255,255,0.3) 0%, transparent 100%),
			radial-gradient(1px 1px at 80% 80%, rgba(255,255,255,0.5) 0%, transparent 100%),
			radial-gradient(1px 1px at 35% 35%, rgba(255,255,255,0.4) 0%, transparent 100%),
			radial-gradient(1.5px 1.5px at 60% 15%, rgba(255,255,255,0.6) 0%, transparent 100%),
			radial-gradient(1.5px 1.5px at 90% 45%, rgba(255,255,255,0.5) 0%, transparent 100%),
			radial-gradient(1px 1px at 5% 50%, rgba(255,255,255,0.3) 0%, transparent 100%),
			radial-gradient(1px 1px at 45% 70%, rgba(255,255,255,0.4) 0%, transparent 100%),
			radial-gradient(1.5px 1.5px at 75% 5%, rgba(255,255,255,0.5) 0%, transparent 100%),
			radial-gradient(1px 1px at 20% 90%, rgba(255,255,255,0.3) 0%, transparent 100%);
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
	.deco { position: fixed; z-index: 1; pointer-events: none; }
	.deco-moon { top: 6%; left: 6%; font-size: 3.5rem; opacity: 0.5; animation: floatA 8s ease-in-out infinite; filter: drop-shadow(0 0 12px rgba(255,255,200,0.4)); }
	.deco-planet { bottom: 12%; right: 6%; font-size: 3rem; opacity: 0.4; animation: floatB 10s ease-in-out infinite; filter: drop-shadow(0 0 10px rgba(200,180,255,0.3)); }
	.deco-rocket { top: 15%; right: 5%; font-size: 2.5rem; opacity: 0.5; animation: rocketFly 12s ease-in-out infinite; filter: drop-shadow(0 0 8px rgba(255,150,50,0.4)); }
	.deco-ufo { bottom: 25%; left: 4%; font-size: 2.2rem; opacity: 0.4; animation: ufoFloat 9s ease-in-out infinite 2s; filter: drop-shadow(0 0 10px rgba(100,255,200,0.3)); }
	.deco-star1 { top: 10%; right: 15%; font-size: 1.4rem; opacity: 0.6; animation: twinkle 3s ease-in-out infinite; }
	.deco-star2 { bottom: 22%; left: 12%; font-size: 1.1rem; opacity: 0.5; animation: twinkle 4s ease-in-out infinite 1s; }
	.deco-star3 { top: 38%; right: 22%; font-size: 1.5rem; opacity: 0.45; animation: twinkle 3.5s ease-in-out infinite 0.5s; }
	.deco-star4 { top: 55%; left: 8%; font-size: 1rem; opacity: 0.55; animation: twinkle 2.5s ease-in-out infinite 1.5s; }
	.deco-star5 { bottom: 10%; right: 18%; font-size: 1.2rem; opacity: 0.5; animation: twinkle 3s ease-in-out infinite 2s; }

	@keyframes floatA {
		0%, 100% { transform: translateY(0) translateX(0); }
		50% { transform: translateY(-18px) translateX(8px); }
	}
	@keyframes floatB {
		0%, 100% { transform: translateY(0) translateX(0) rotate(0deg); }
		50% { transform: translateY(12px) translateX(-10px) rotate(8deg); }
	}
	@keyframes rocketFly {
		0%, 100% { transform: translateY(0) translateX(0) rotate(-25deg); }
		25% { transform: translateY(-25px) translateX(12px) rotate(-18deg); }
		50% { transform: translateY(-10px) translateX(-8px) rotate(-30deg); }
		75% { transform: translateY(-30px) translateX(8px) rotate(-20deg); }
	}
	@keyframes ufoFloat {
		0%, 100% { transform: translateY(0) translateX(0) rotate(0deg); }
		25% { transform: translateY(-12px) translateX(15px) rotate(5deg); }
		50% { transform: translateY(5px) translateX(-10px) rotate(-3deg); }
		75% { transform: translateY(-8px) translateX(8px) rotate(2deg); }
	}
	@keyframes twinkle {
		0%, 100% { opacity: 0.3; transform: scale(1); }
		50% { opacity: 0.8; transform: scale(1.4); }
	}

	/* Shooting stars */
	.shooting-star {
		position: fixed; z-index: 1; pointer-events: none;
		width: 4px; height: 4px; border-radius: 50%;
		background: white;
		box-shadow: 0 0 8px 3px rgba(255,255,255,0.7), -25px 0 18px 2px rgba(255,255,255,0.4), -50px 0 30px 1px rgba(255,255,255,0.15);
	}
	.ss1 {
		top: 6%; left: -5%;
		animation: shootingStar 4s linear infinite 0s;
	}
	.ss2 {
		top: 22%; left: -5%;
		animation: shootingStar 5s linear infinite 2.5s;
	}
	.ss3 {
		top: 42%; left: -5%;
		animation: shootingStar 3.5s linear infinite 5s;
	}
	@keyframes shootingStar {
		0% { transform: translateX(0) translateY(0) rotate(-35deg); opacity: 0; }
		5% { opacity: 1; }
		70% { opacity: 1; }
		100% { transform: translateX(110vw) translateY(45vh) rotate(-35deg); opacity: 0; }
	}

	/* Comets */
	.comet {
		position: fixed; z-index: 1; pointer-events: none;
		width: 5px; height: 5px; border-radius: 50%;
		background: #fbbf24;
		box-shadow: 0 0 10px 4px rgba(251,191,36,0.6), -20px 0 25px 3px rgba(251,191,36,0.25), -40px 0 35px 2px rgba(251,191,36,0.1);
	}
	.comet1 {
		top: 12%; right: -5%;
		animation: cometLeft 6s linear infinite 1s;
	}
	.comet2 {
		bottom: 18%; right: -5%;
		animation: cometLeft 7s linear infinite 4s;
	}
	@keyframes cometLeft {
		0% { transform: translateX(0) translateY(0) rotate(25deg); opacity: 0; }
		5% { opacity: 1; }
		70% { opacity: 1; }
		100% { transform: translateX(-110vw) translateY(35vh) rotate(25deg); opacity: 0; }
	}

	:global(.animate-spin-slow) { animation: spin 4s linear infinite; }
	@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

	.bounce-in { animation: bounceIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55); }
	@keyframes bounceIn { 0% { transform: scale(0.3); opacity: 0; } 60% { transform: scale(1.1); } 100% { transform: scale(1); opacity: 1; } }
	.fade-in { animation: fadeUp 0.5s ease-out; }
	@keyframes fadeUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }

	/* Game title */
	.game-title {
		font-family: 'Permanent Marker', cursive;
		color: #fbbf24;
		text-shadow:
			0 0 10px rgba(251,191,36,0.5),
			0 2px 0 #b45309,
			0 3px 0 #92400e,
			0 4px 0 #78350f,
			0 5px 8px rgba(0,0,0,0.3);
		-webkit-text-stroke: 1px rgba(255,255,255,0.15);
		letter-spacing: 0.04em;
	}
</style>
