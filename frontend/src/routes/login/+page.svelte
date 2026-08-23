<script lang="ts">
	import { login, staffLogin, studentLogin } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { Phone, Key, Hash, Calendar, Shield, ArrowLeft } from 'lucide-svelte';

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

<svelte:head>
	<title>Pragati — MDRS Bahaddurghatta</title>
</svelte:head>

<div class="page">
	<!-- Header — a11y.quest style -->
	<header class="site-header">
		<a href="/" class="identity">
			<div class="logo-box" aria-hidden="true">
				<span style="font-size:1.65rem">🏫</span>
			</div>
			<div class="identity-text">
				<p class="wordmark">PRAGATI<span class="dot">.</span><span class="wordmark-kannada">quest</span></p>
				<p class="tagline"><span class="tagline-strong">MDRS (SC-32) Bahaddurghatta</span> <span class="tagline-sep">·</span> KREIS, Karnataka</p>
			</div>
		</a>
		<div class="header-right">
			<span class="header-badge">Learn · Play · Grow</span>
		</div>
	</header>

	{#if view === 'home'}
		<!-- Intro -->
		<div class="intro">
			<h1 class="intro-title">A school that <em>feels</em> like a playground.</h1>
			<p class="intro-sub">Quizzes, timetables &amp; progress — all in one place. Pick a card to begin. Built for students, teachers &amp; parents of MDRS Bahaddurghatta.</p>
			<div class="intro-meta">
				<span class="pill"><span class="pill-dot" style="background:var(--teal)"></span> 3 spaces</span>
				<span class="pill"><span class="pill-dot" style="background:var(--amber)"></span> NEP 2020</span>
				<span class="pill">ಪ್ರಗತಿ — progress</span>
			</div>
		</div>

		<!-- Cards — 3-up brutalist -->
		<div class="card-grid">
			<!-- Quiz Arena -->
			<a href="/play" class="card">
				<div class="card-top">
					<div class="icon-wrap icon-amber">🎮</div>
					<span class="eyebrow">Practice · GK · States</span>
				</div>
				<h2 class="card-title">Quiz Arena</h2>
				<p class="card-desc">Battle through questions, earn points, keep a streak and climb the leaderboard. Works off-line, works on phones.</p>
				<div class="tag-row">
					<span class="tag tag-amber">⭐ Points</span>
					<span class="tag tag-mint">🔥 Streak</span>
					<span class="tag tag-coral">🏆 Rank</span>
				</div>
				<span class="card-cta">Play now <span aria-hidden="true">→</span></span>
			</a>

			<!-- Timetable -->
			<a href="/timetable" class="card">
				<div class="card-top">
					<div class="icon-wrap icon-mint">🗓️</div>
					<span class="eyebrow">Class-wise · Subject-wise</span>
				</div>
				<h2 class="card-title">Timetable</h2>
				<p class="card-desc">Your weekly schedule at a glance — know what's next, never miss a period.</p>
				<div class="mini-schedule" aria-hidden="true">
					<div class="mini-row"><span class="mini-dot dot-amber"></span><span class="mini-bar" style="width:78%"></span><span class="mini-label">8AM</span></div>
					<div class="mini-row"><span class="mini-dot dot-teal"></span><span class="mini-bar" style="width:58%"></span><span class="mini-label">10AM</span></div>
					<div class="mini-row"><span class="mini-dot dot-coral"></span><span class="mini-bar" style="width:88%"></span><span class="mini-label">2PM</span></div>
				</div>
				<span class="card-cta cta-ghost">View schedule <span aria-hidden="true">→</span></span>
			</a>

			<!-- Sign In -->
			<button onclick={() => view = 'login'} class="card card-interactive">
				<div class="card-top">
					<div class="avatar-row" aria-hidden="true">
						<span class="ava ava-blue">🧑‍🎓</span>
						<span class="ava ava-amber">👩‍🏫</span>
						<span class="ava ava-plum">👨‍💼</span>
					</div>
					<span class="eyebrow">Secure sign-in</span>
				</div>
				<h2 class="card-title">Sign in</h2>
				<p class="card-desc">Students use SATS + DOB. Staff use mobile + password. Admin uses email.</p>
				<div class="tag-row">
					<span class="tag tag-paper">Student</span>
					<span class="tag tag-paper">Teacher</span>
					<span class="tag tag-paper">Admin</span>
				</div>
				<span class="card-cta cta-dark">Enter <span aria-hidden="true">→</span></span>
			</button>
		</div>

		<!-- Footer note -->
		<p class="footer-note">Karnataka Residential Educational Institutions Society · MDRS (SC-32) Bahaddurghatta · Holiday? Check the <a href="/timetable">timetable</a>.</p>

	{:else}
		<!-- Login view -->
		<div class="login-wrap">
			<button onclick={() => { view = 'home'; error = ''; }} class="back-link">
				<ArrowLeft size={16} /> Back to home
			</button>

			<div class="login-card">
				<div class="login-head">
					<div class="logo-box small" aria-hidden="true"><span style="font-size:1.25rem">🏫</span></div>
					<div>
						<h1 class="login-title">Welcome back</h1>
						<p class="login-sub">Choose your role to enter</p>
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
						<p class="field-hint">No password needed. Use your school SATS number and DOB.</p>
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

			<p class="login-foot">Protected by KREIS. <a href="/">Back to site</a></p>
		</div>
	{/if}
</div>

<style>
	.page {
		max-width: 1040px;
		margin: 0 auto;
		padding: 1rem clamp(1rem, 4vw, 1.5rem) 3rem;
	}
	/* header */
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
		width: 56px;
		height: 56px;
		flex: none;
		background: var(--paper);
		border: 3px solid var(--ink);
		box-shadow: 4px 4px 0 var(--ink);
		border-radius: 16px;
		display: grid;
		place-items: center;
	}
	.logo-box.small { width: 44px; height: 44px; border-radius: 12px; box-shadow: 3px 3px 0 var(--ink); }
	.identity-text { min-width: 0; }
	.wordmark {
		font-family: var(--font-display);
		font-weight: 800;
		font-size: clamp(1.5rem, 1.1rem + 2vw, 2rem);
		letter-spacing: -0.03em;
		color: var(--ink);
		margin: 0;
		line-height: 1;
	}
	.wordmark-kannada { font-weight: 700; color: var(--plum); }
	.dot { color: var(--coral); }
	.tagline {
		color: var(--ink-soft);
		font-size: 0.95rem;
		margin: 0.2rem 0 0;
		line-height: 1.3;
	}
	.tagline-strong { font-weight: 700; color: var(--ink); }
	.tagline-sep { opacity: 0.4; margin: 0 0.2rem; }
	.header-right { display: flex; align-items: center; padding-top: 0.35rem; }
	.header-badge {
		font-family: var(--font-mono);
		font-size: 0.75rem;
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		background: var(--paper);
		border: 2px solid var(--ink);
		border-radius: 999px;
		padding: 0.35rem 0.7rem;
		color: var(--ink);
	}

	/* intro */
	.intro { margin-bottom: 1.75rem; }
	.intro-title {
		font-family: var(--font-display);
		font-size: clamp(1.9rem, 1.4rem + 2.4vw, 3rem);
		font-weight: 800;
		line-height: 1.05;
		color: var(--ink);
		margin: 0;
		max-width: 22ch;
	}
	.intro-title em { font-style: normal; background: var(--amber); padding: 0 0.22em; border-radius: 6px; box-decoration-break: clone; }
	.intro-sub {
		color: var(--ink-soft);
		font-size: clamp(1rem, 0.95rem + 0.4vw, 1.15rem);
		line-height: 1.5;
		margin: 0.75rem 0 0;
		max-width: 60ch;
	}
	.intro-meta { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-top: 1rem; }
	.pill {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.82rem;
		font-weight: 700;
		border: 2px solid var(--ink);
		background: var(--paper);
		border-radius: 999px;
		padding: 0.28rem 0.65rem;
		color: var(--ink);
	}
	.pill-dot { width: 8px; height: 8px; border-radius: 50%; border: 1.5px solid var(--ink); display: inline-block; }

	/* cards */
	.card-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 1.25rem;
	}
	@media (max-width: 860px) { .card-grid { grid-template-columns: 1fr; } .site-header { flex-direction: column; } }
	.card {
		background: var(--paper);
		border: 3px solid var(--ink);
		box-shadow: 8px 8px 0 var(--ink);
		border-radius: 22px;
		padding: 1.35rem 1.35rem 1.15rem;
		display: flex;
		flex-direction: column;
		gap: 0.85rem;
		text-decoration: none;
		color: var(--ink);
		transition: transform 80ms, box-shadow 80ms;
	}
	.card:hover { transform: translate(2px, 2px); box-shadow: 6px 6px 0 var(--ink); }
	.card:active { transform: translate(6px, 6px); box-shadow: 2px 2px 0 var(--ink); }
	.card-interactive { cursor: pointer; width: 100%; text-align: left; font: inherit; }
	.card-top { display: flex; align-items: center; justify-content: space-between; gap: 0.75rem; }
	.eyebrow {
		font-family: var(--font-mono);
		font-size: 0.7rem;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--ink-soft);
	}
	.icon-wrap {
		width: 46px;
		height: 46px;
		border-radius: 14px;
		border: 2.5px solid var(--ink);
		display: grid;
		place-items: center;
		font-size: 1.35rem;
		flex: none;
	}
	.icon-amber { background: var(--amber); }
	.icon-mint { background: var(--mint); }
	.avatar-row { display: flex; gap: -0.5rem; }
	.ava {
		width: 40px;
		height: 40px;
		border-radius: 12px;
		border: 2.5px solid var(--ink);
		display: grid;
		place-items: center;
		font-size: 1.1rem;
		margin-right: -8px;
	}
	.ava-blue { background: #DBEAFE; }
	.ava-amber { background: var(--amber-tint); }
	.ava-plum { background: #EDE9FE; }
	.card-title {
		font-family: var(--font-display);
		font-size: 1.55rem;
		font-weight: 800;
		line-height: 1.1;
		margin: 0;
	}
	.card-desc {
		color: var(--ink-soft);
		font-size: 0.96rem;
		line-height: 1.5;
		margin: 0;
		flex: 1;
	}
	.tag-row { display: flex; flex-wrap: wrap; gap: 0.4rem; }
	.tag {
		font-size: 0.78rem;
		font-weight: 700;
		border: 2px solid var(--ink);
		border-radius: 999px;
		padding: 0.22rem 0.55rem;
	}
	.tag-amber { background: var(--amber-tint); }
	.tag-mint { background: var(--mint); }
	.tag-coral { background: var(--coral-tint); }
	.tag-paper { background: var(--cream); }
	.mini-schedule { display: flex; flex-direction: column; gap: 0.45rem; padding: 0.2rem 0; }
	.mini-row { display: flex; align-items: center; gap: 0.5rem; }
	.mini-dot { width: 10px; height: 10px; border-radius: 50%; border: 2px solid var(--ink); flex: none; }
	.dot-amber { background: var(--amber); } .dot-teal { background: #0E7C71; } .dot-coral { background: #FBDAD3; }
	.mini-bar { height: 8px; border-radius: 999px; border: 2px solid var(--ink); background: var(--cream-deep); }
	.mini-label { font-family: var(--font-mono); font-size: 0.7rem; font-weight: 700; color: var(--ink-soft); width: 2.2rem; text-align: right; }
	.card-cta {
		margin-top: 0.2rem;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.35rem;
		min-height: 44px;
		padding: 0.6rem 1rem;
		border-radius: 14px;
		border: 3px solid var(--ink);
		background: var(--amber);
		color: var(--ink);
		box-shadow: 4px 4px 0 var(--ink);
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 1rem;
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

	/* login view */
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
		border: 3px solid var(--ink);
		box-shadow: 8px 8px 0 var(--ink);
		border-radius: 22px;
		padding: clamp(1.2rem, 1rem + 1.5vw, 1.9rem);
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
		min-height: 48px;
		border-radius: 14px;
		border: 3px solid var(--ink);
		background: var(--amber);
		color: var(--ink);
		box-shadow: 5px 5px 0 var(--ink);
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 1.05rem;
		cursor: pointer;
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}
	.btn-primary:hover { transform: translate(2px, 2px); box-shadow: 3px 3px 0 var(--ink); }
	.btn-primary:active { transform: translate(5px, 5px); box-shadow: 0 0 0 var(--ink); }
	.btn-primary:disabled { opacity: 0.55; cursor: not-allowed; transform: none; box-shadow: 5px 5px 0 var(--ink); }
	.field-hint { color: var(--ink-soft); font-size: 0.85rem; margin: 0; text-align: center; }
	.login-foot { text-align: center; color: var(--ink-soft); font-size: 0.85rem; margin-top: 1rem; }
	.login-foot a { color: var(--plum); font-weight: 700; }

	@media (prefers-reduced-motion: reduce) { .card, .btn-primary { transition: none; } }
</style>
