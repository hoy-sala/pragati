<script lang="ts">
	import { login, staffLogin, studentLogin } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { Phone, Key, Hash, Calendar, Shield, ArrowLeft, Gamepad2, CalendarDays, LogIn, Building2, Users, GraduationCap, Info } from 'lucide-svelte';

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
	let aboutOpen = $state(false);

	async function handleStudentSubmit() {
		error = ''; loading = true;
		let result;
		try { result = await studentLogin(satsNumber, dateOfBirth); }
		catch { result = { error: 'Unable to reach server.' }; }
		loading = false;
		if ('error' in result) error = result.error; else goto('/reports');
	}
	async function handleStaffSubmit() {
		error = ''; loading = true;
		let result;
		try { result = await staffLogin(mobile, staffPassword); }
		catch { result = { error: 'Unable to reach server.' }; }
		loading = false;
		if ('error' in result) error = result.error; else goto('/reports');
	}
	async function handleAdminSubmit() {
		error = ''; loading = true;
		let result;
		try { result = await login(adminEmail, adminPassword); }
		catch { result = { error: 'Unable to reach server.' }; }
		loading = false;
		if ('error' in result) error = result.error; else goto('/reports');
	}
</script>

<svelte:head>
	<title>Pragati — MDRS Bahaddurghatta</title>
</svelte:head>

<div class="page">
	<header class="site-header">
		<a href="/" class="identity" aria-label="Pragati home">
			<div class="logo-box" aria-hidden="true">
				<Building2 size={22} strokeWidth={2.2} />
			</div>
			<div class="identity-text">
				<p class="wordmark wordmark-kannada-main">ಪ್ರಗತಿ</p>
				<p class="tagline">MDRS (SC-32) Bahaddurghatta — where curiosity becomes progress</p>
			</div>
		</a>
		<button class="about-link" onclick={() => aboutOpen = true}>
			<Info size={16} strokeWidth={2} aria-hidden="true" />
			<span class="about-label">About</span>
		</button>
	</header>

	{#if view === 'home'}
		<div class="intro intro-kids">
			<div class="kids-art" aria-hidden="true"><span class="ka ka1">🎮</span><span class="ka ka2">⭐</span><span class="ka ka3">📚</span></div>
			<h1 class="intro-title"><em>Let's play</em> & learn!</h1>
			<p class="intro-sub">Pick a quiz, check your timetable, or sign in — all yours.</p>
		</div>

		<div class="card-grid">
			<!-- Quiz -->
			<a href="/play" class="card card-quiz">
				<div class="card-head">
					<span class="eyebrow">Classes 6–10 · Your level</span>
					<span class="icon-badge badge-amber"><Gamepad2 size={20} /></span>
				</div>
				<h2 class="card-title">Quizzes</h2>
				<p class="card-desc">Pick your class & subject — then play. Fast, fun, no pressure.</p>
				<span class="card-cta">Let's play <span aria-hidden="true">→</span></span>
			</a>

			<!-- Timetable -->
			<a href="/timetable" class="card card-tt">
				<div class="card-head">
					<span class="eyebrow">Know what's next</span>
					<span class="icon-badge badge-mint"><CalendarDays size={20} /></span>
				</div>
				<h2 class="card-title">Timetable</h2>
				<p class="card-desc">Your week at a glance — simple and clear.</p>
				<span class="card-cta cta-ghost">View week <span aria-hidden="true">→</span></span>
			</a>

			<!-- Sign in -->
			<button onclick={() => view = 'login'} class="card card-interactive card-signin">
				<div class="card-head">
					<span class="eyebrow">For you</span>
					<span class="icon-badge badge-plum"><LogIn size={20} /></span>
				</div>
				<h2 class="card-title">Sign in</h2>
				<p class="card-desc">Students, teachers & parents — jump in here.</p>
				<div class="role-row">
					<span class="role"><GraduationCap size={13} /> Student</span>
					<span class="role"><Users size={13} /> Staff</span>
				</div>
				<span class="card-cta">Enter <span aria-hidden="true">→</span></span>
			</button>
		</div>

		<p class="footer-note">Developed by <strong>Hoysala T</strong>, Computer Science Teacher — MDRS (SC-32) Bahaddurghatta</p>

		<!-- About dialog -->
		{#if aboutOpen}
			<div class="scrim" role="presentation" onclick={() => aboutOpen = false}></div>
			<div class="about" role="dialog" aria-modal="true" aria-label="About Pragati">
				<button class="about-close" onclick={() => aboutOpen = false} aria-label="Close">✕</button>
				<h2 class="about-title">About ಪ್ರಗತಿ</h2>
				<div class="about-body">
					<p>Dear Students,</p>
					<p>We — the staff of <strong>MDRS (SC-32) Bahaddurghatta</strong> — built <span class="font-kannada">ಪ್ರಗತಿ</span> — Pragati — for you. A simple, safe place for your timetable, your quizzes and your progress.</p>
					<p>Our quizzes are practice, not pressure. They come straight from what you learn in class. Pick your level, play at your pace, and learn a little every single day.</p>
					<p>Let's build a learning revolution — one question, one day at a time. It starts with you. <a href="/play">Try a quiz</a> or <button class="link-btn" onclick={() => { aboutOpen = false; view = 'login'; }}>sign in</button> with your class.</p>
				</div>
			</div>
		{/if}

	{:else}
		<div class="login-wrap">
			<button onclick={() => { view = 'home'; error = ''; }} class="back-link">
				<ArrowLeft size={16} /> Back to home
			</button>

			<div class="login-card">
				<div class="login-head">
					<div class="logo-box small" aria-hidden="true"><Building2 size={18} /></div>
					<div>
						<h1 class="login-title">Welcome back</h1>
						<p class="login-sub">Choose your role to continue</p>
					</div>
				</div>

				<div class="tabs" role="tablist" aria-label="Login role">
					<button role="tab" aria-selected={activeTab === 'student'} onclick={() => { activeTab = 'student'; error = ''; }} class="tab {activeTab === 'student' ? 'tab-active' : ''}">Student</button>
					<button role="tab" aria-selected={activeTab === 'staff'} onclick={() => { activeTab = 'staff'; error = ''; }} class="tab {activeTab === 'staff' ? 'tab-active' : ''}">Staff</button>
					<button role="tab" aria-selected={activeTab === 'admin'} onclick={() => { activeTab = 'admin'; error = ''; }} class="tab {activeTab === 'admin' ? 'tab-active' : ''}">Admin</button>
				</div>

				{#if activeTab === 'student'}
					<form onsubmit={(e) => { e.preventDefault(); handleStudentSubmit(); }} class="form">
						<label class="field">
							<span class="field-label">SATS number</span>
							<span class="control">
								<Hash size={16} class="field-icon" />
								<input type="text" bind:value={satsNumber} maxlength="9" placeholder="9-digit SATS number" class="input" autocomplete="off" />
							</span>
						</label>
						<label class="field">
							<span class="field-label">Date of birth</span>
							<span class="control">
								<Calendar size={16} class="field-icon" />
								<input type="date" bind:value={dateOfBirth} class="input" />
							</span>
						</label>
						{#if error}<div class="error" role="alert">{error}</div>{/if}
						<button type="submit" disabled={loading || !satsNumber || !dateOfBirth} class="btn-primary"> {loading ? 'Signing in…' : 'Sign in →'} </button>
						<p class="field-hint">No password needed.</p>
					</form>
				{:else if activeTab === 'staff'}
					<form onsubmit={(e) => { e.preventDefault(); handleStaffSubmit(); }} class="form">
						<label class="field">
							<span class="field-label">Mobile number</span>
							<span class="control">
								<Phone size={16} class="field-icon" />
								<input type="tel" bind:value={mobile} maxlength="10" placeholder="10-digit mobile" class="input" />
							</span>
						</label>
						<label class="field">
							<span class="field-label">Password</span>
							<span class="control">
								<Key size={16} class="field-icon" />
								<input type="password" bind:value={staffPassword} placeholder="Password" class="input" />
							</span>
						</label>
						{#if error}<div class="error" role="alert">{error}</div>{/if}
						<button type="submit" disabled={loading || !mobile || !staffPassword} class="btn-primary"> {loading ? 'Signing in…' : 'Sign in →'} </button>
					</form>
				{:else}
					<form onsubmit={(e) => { e.preventDefault(); handleAdminSubmit(); }} class="form">
						<label class="field">
							<span class="field-label">Email address</span>
							<span class="control">
								<Shield size={16} class="field-icon" />
								<input type="email" bind:value={adminEmail} placeholder="you@school.edu" class="input" autocomplete="email" />
							</span>
						</label>
						<label class="field">
							<span class="field-label">Password</span>
							<span class="control">
								<Key size={16} class="field-icon" />
								<input type="password" bind:value={adminPassword} placeholder="Password" class="input" />
							</span>
						</label>
						{#if error}<div class="error" role="alert">{error}</div>{/if}
						<button type="submit" disabled={loading || !adminEmail || !adminPassword} class="btn-primary"> {loading ? 'Signing in…' : 'Sign in →'} </button>
					</form>
				{/if}
			</div>


		</div>
	{/if}
</div>

<style>
	.page {
		max-width: 1040px;
		margin: 0 auto;
		padding: 1rem clamp(1rem, 4vw, 1.5rem) 3.5rem;
	}
	.site-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
		margin: 1.25rem 0 2rem;
	}
	.identity {
		display: flex;
		align-items: center;
		gap: 0.9rem;
		text-decoration: none;
		color: inherit;
		min-width: 0;
	}
	.logo-box {
		width: 48px;
		height: 48px;
		flex: none;
		background: var(--paper);
		border: 2px solid var(--ink);
		box-shadow: 3px 3px 0 var(--ink);
		border-radius: 12px;
		display: grid;
		place-items: center;
		color: var(--ink);
	}
	.logo-box.small { width: 40px; height: 40px; border-radius: 10px; box-shadow: 2px 2px 0 var(--ink); }
	.identity-text { min-width: 0; }
	.wordmark {
		font-family: var(--font-display);
		font-weight: 800;
		font-size: clamp(1.55rem, 1.2rem + 1.6vw, 2rem);
		letter-spacing: -0.03em;
		color: var(--ink);
		margin: 0;
		line-height: 1;
	}
	.wordmark-kannada-main { font-family: 'Anek Kannada', var(--font-display), system-ui, sans-serif; font-weight: 700; letter-spacing: -0.01em; }
	.tagline {
		color: var(--ink-soft);
		font-size: 0.92rem;
		margin: 0.22rem 0 0;
		line-height: 1.3;
	}
	.wordmark-kannada { font-family: 'Anek Kannada', system-ui, sans-serif; font-weight: 700; color: var(--ink-soft); font-size: 0.88em; margin-left: 0.4rem; letter-spacing: -0.01em; }
	.about-link {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		background: transparent;
		border: 2px solid transparent;
		border-radius: 10px;
		padding: 0.35rem 0.6rem;
		font: inherit;
		font-weight: 600;
		color: var(--ink-soft);
		cursor: pointer;
		margin-top: 0.2rem;
		text-decoration: underline;
		text-underline-offset: 3px;
		text-decoration-thickness: 1.5px;
	}
	.about-link:hover { color: var(--ink); background: var(--paper); border-color: var(--ink); }
	.about-icon { display: inline-grid; place-items: center; width: 18px; height: 18px; border: 1.5px solid currentColor; border-radius: 50%; font-size: 0.7rem; font-weight: 700; }

	.intro { margin-bottom: 1.5rem; max-width: 60ch; }
	.intro-kids { max-width: 48ch; }
	.kids-art { display: flex; gap: 0.6rem; margin-bottom: 0.8rem; }
	.ka { width: 42px; height: 42px; display: grid; place-items: center; border-radius: 12px; border: 1.5px solid var(--ink); font-size: 1.35rem; box-shadow: 2px 2px 0 var(--ink); animation: kaFloat 3s ease-in-out infinite; }
	.ka1 { background: var(--amber); --r: -6deg; transform: rotate(var(--r)); animation-delay: 0s; }
	.ka2 { background: var(--mint); --r: 5deg; transform: rotate(var(--r)); animation-delay: 0.4s; }
	.ka3 { background: var(--paper); --r: -3deg; transform: rotate(var(--r)); animation-delay: 0.8s; }
	@keyframes kaFloat { 0%,100% { transform: translateY(0) rotate(var(--r)); } 50% { transform: translateY(-4px) rotate(var(--r)); } }
	.intro-title {
		font-family: var(--font-display);
		font-size: clamp(1.7rem, 1.2rem + 1.8vw, 2.45rem);
		font-weight: 800;
		line-height: 1.08;
		color: var(--ink);
		margin: 0;
		letter-spacing: -0.025em;
	}
	.intro-title em { font-style: normal; text-decoration: underline; text-decoration-color: var(--amber); text-decoration-thickness: 6px; text-underline-offset: 3px; }
	.intro-sub {
		color: var(--ink-soft);
		font-size: 1rem;
		line-height: 1.65;
		margin: 0.6rem 0 0;
	}
	.intro-meta { display: flex; flex-wrap: wrap; gap: 0.45rem; margin-top: 0.9rem; }
	.pill {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		font-size: 0.76rem;
		font-weight: 600;
		border: 1.5px solid var(--ink);
		background: var(--paper);
		border-radius: 999px;
		padding: 0.24rem 0.6rem;
		color: var(--ink);
	}
	.mono { font-family: var(--font-mono); }
	.pill-dot { width: 8px; height: 8px; border-radius: 50%; border: 1.5px solid var(--ink); display: inline-block; }

	.card-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 1rem;
		align-items: stretch;
	}
	@media (max-width: 860px) { .card-grid { grid-template-columns: 1fr; } .site-header { flex-wrap: wrap; } }
	.card {
		background: var(--paper);
		border: 2px solid var(--ink);
		box-shadow: 5px 5px 0 var(--ink);
		border-radius: 18px;
		padding: 1.25rem 1.25rem 1.1rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		text-decoration: none;
		color: var(--ink);
		transition: box-shadow 120ms, transform 120ms;
	}
	.card:hover { transform: translate(1px, 1px); box-shadow: 4px 4px 0 var(--ink); }
	.card:active { transform: translate(3px, 3px); box-shadow: 2px 2px 0 var(--ink); }
	.card-interactive { cursor: pointer; width: 100%; text-align: left; font: inherit; }
	.card-head { display: flex; align-items: center; justify-content: space-between; gap: 0.75rem; }
	.eyebrow {
		font-family: var(--font-mono);
		font-size: 0.68rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--ink-soft);
	}
	.icon-badge {
		width: 40px; height: 40px; border-radius: 11px;
		border: 1.5px solid var(--ink); background: var(--cream);
		display: grid; place-items: center; color: var(--ink); flex: none;
	}
	.badge-amber { background: var(--amber); }
	.badge-mint { background: var(--mint); }
	.badge-plum { background: #EDE9FE; }
	.card-quiz { background: var(--paper); }
	.card-tt { background: var(--paper); }
	.card-title {
		font-family: var(--font-display);
		font-size: 1.28rem;
		font-weight: 800;
		line-height: 1.15;
		margin: 0;
	}
	.card-desc {
		color: var(--ink-soft);
		font-size: 0.9rem;
		line-height: 1.6;
		margin: 0;
		flex: 1;
	}
	.tag-row { display: flex; flex-wrap: wrap; gap: 0.4rem; }
	.tag {
		font-size: 0.72rem;
		font-weight: 600;
		border: 1.5px solid var(--ink);
		border-radius: 999px;
		padding: 0.18rem 0.5rem;
		background: var(--paper);
	}
	.tag-mint { background: var(--mint); }
	.role-row { display: flex; flex-wrap: wrap; gap: 0.35rem; }
	.role {
		display: inline-flex; align-items: center; gap: 0.28rem;
		font-size: 0.72rem; font-weight: 600;
		border: 1.5px solid var(--ink); border-radius: 999px; padding: 0.18rem 0.5rem; background: var(--cream);
	}
	.card-cta {
		margin-top: auto;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.3rem;
		min-height: 42px;
		padding: 0.55rem 1rem;
		border-radius: 12px;
		border: 2px solid var(--ink);
		background: var(--amber);
		color: var(--ink);
		box-shadow: 3px 3px 0 var(--ink);
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 0.94rem;
	}
	.cta-ghost { background: var(--paper); }
	.cta-dark { background: var(--ink); color: var(--paper); }
	.footer-note {
		margin-top: 1.75rem;
		text-align: center;
		color: var(--ink-soft);
		font-size: 0.88rem;
	}
	.footer-note a { color: var(--plum); font-weight: 700; }

	.scrim { position: fixed; inset: 0; background: rgba(31,26,46,0.45); backdrop-filter: blur(2px); z-index: 60; }
	.about {
		position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%);
		z-index: 61; background: var(--paper); border: 2.5px solid var(--ink);
		box-shadow: 6px 6px 0 var(--ink); border-radius: 18px;
		padding: 1.4rem; width: min(560px, calc(100vw - 2rem)); max-height: calc(100dvh - 2rem); overflow: auto;
	}
	.about-close {
		position: absolute; top: 0.6rem; right: 0.6rem;
		width: 36px; height: 36px; border-radius: 10px; border: 2.5px solid var(--ink);
		background: var(--paper); cursor: pointer; display: grid; place-items: center; color: var(--ink);
	}
	.about-title { font-family: var(--font-display); font-size: 1.45rem; font-weight: 800; margin: 0 2rem 0 0; }
	.about-body { display: flex; flex-direction: column; gap: 0.85rem; line-height: 1.6; color: var(--ink); margin-top: 0.85rem; }
	.about-body p { margin: 0; }
	.about-body a, .link-btn { color: var(--plum); font-weight: 700; background: none; border: 0; padding: 0; cursor: pointer; text-decoration: underline; font: inherit; }

	.login-wrap { max-width: 520px; margin: 0 auto; }
	.back-link {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		border: 2px solid transparent;
		background: transparent;
		color: var(--ink-soft);
		font-weight: 700;
		font-size: 0.9rem;
		cursor: pointer;
		padding: 0.3rem 0.5rem;
		border-radius: 10px;
		margin-bottom: 1rem;
	}
	.back-link:hover { background: var(--paper); border-color: var(--ink); color: var(--ink); }
	.login-card {
		background: var(--paper);
		border: 2.5px solid var(--ink);
		box-shadow: 6px 6px 0 var(--ink);
		border-radius: 18px;
		padding: clamp(1.15rem, 1rem + 1vw, 1.6rem);
	}
	.login-head { display: flex; gap: 0.85rem; align-items: center; margin-bottom: 1.25rem; }
	.login-title { font-family: var(--font-display); font-size: 1.6rem; font-weight: 800; margin: 0; line-height: 1; }
	.login-sub { color: var(--ink-soft); font-size: 0.95rem; margin: 0.15rem 0 0; }
	.tabs {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 0.4rem;
		background: var(--cream);
		border: 2.5px solid var(--ink);
		border-radius: 14px;
		padding: 0.35rem;
		margin-bottom: 1.25rem;
	}
	.tab {
		min-height: 40px;
		border-radius: 10px;
		border: 2px solid transparent;
		background: transparent;
		font-weight: 700;
		font-size: 0.92rem;
		cursor: pointer;
		color: var(--ink-soft);
	}
	.tab-active {
		background: var(--amber);
		border-color: var(--ink);
		color: var(--ink);
		box-shadow: 2px 2px 0 var(--ink);
	}
	.form { display: flex; flex-direction: column; gap: 0.9rem; }
	.field { display: flex; flex-direction: column; gap: 0.35rem; }
	.field-label {
		font-size: 0.78rem;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: var(--ink);
	}
	.control { position: relative; display: flex; align-items: center; }
	:global(.field-icon) { position: absolute; left: 0.75rem; color: var(--ink-soft); pointer-events: none; }
	.input {
		width: 100%;
		min-height: 46px;
		padding: 0.6rem 0.85rem 0.6rem 2.4rem;
		border-radius: 12px;
		border: 2.5px solid var(--ink);
		background: var(--cream);
		color: var(--ink);
		font-size: 0.98rem;
		font-weight: 500;
		outline: none;
	}
	.input:focus { background: var(--paper); border-color: var(--plum); box-shadow: 0 0 0 3px rgba(107,63,160,0.18); }
	.error {
		background: var(--coral-tint);
		border: 2px solid var(--coral);
		color: var(--ink);
		border-radius: 12px;
		padding: 0.6rem 0.85rem;
		font-size: 0.92rem;
		font-weight: 600;
	}
	.btn-primary {
		min-height: 46px;
		border-radius: 12px;
		border: 2.5px solid var(--ink);
		background: var(--amber);
		color: var(--ink);
		box-shadow: 3px 3px 0 var(--ink);
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 1rem;
		cursor: pointer;
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}
	.btn-primary:hover { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--ink); }
	.btn-primary:active { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
	.btn-primary:disabled { opacity: 0.52; cursor: not-allowed; transform: none; box-shadow: 3px 3px 0 var(--ink); }
	.field-hint { color: var(--ink-soft); font-size: 0.85rem; margin: 0; text-align: center; }
	.login-foot { text-align: center; color: var(--ink-soft); font-size: 0.85rem; margin-top: 1rem; }

	@media (prefers-reduced-motion: reduce) { .card, .btn-primary { transition: none; } }
</style>
