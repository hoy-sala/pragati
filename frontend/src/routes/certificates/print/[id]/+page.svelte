<script lang="ts">
	import { api, apiUrl } from '$lib/api/client.svelte';
	import type { CertificateDetail } from '$lib/types';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	const certId = $page.params.id;

	let cert = $state<CertificateDetail | null>(null);
	let loading = $state(true);
	let error = $state('');

	let pageScale = $state(1);
	let ready = $state(false);

	function fitToScreen() {
		const padding = 48;
		const availW = window.innerWidth - padding;
		const availH = window.innerHeight - 140;
		const pageW = (210 / 25.4) * 96;
		const pageH = (297 / 25.4) * 96;
		pageScale = Math.min(1, availW / pageW, availH / pageH);
	}

	const ROLE_LABELS: Record<string, string> = {
		principal: 'Principal',
		chief_guest: 'Chief Guest',
		chief_judge: 'Chief Judge',
		judge: 'Judge',
		coordinator: 'Event Coordinator',
	};

	function positionLabel(pos: string): string {
		if (pos === '1st') return 'First Prize';
		if (pos === '2nd') return 'Second Prize';
		if (pos === '3rd') return 'Third Prize';
		return 'Certificate of Participation';
	}

	function categoryEn(cat?: string): string {
		if (cat === 'sports') return 'Sports Competition';
		if (cat === 'cultural') return 'Cultural Competition';
		if (cat === 'academic') return 'Academic Competition';
		return 'Competition';
	}

	function formatDate(d?: string): string {
		if (!d) return '';
		return new Date(d).toLocaleDateString('en-GB', { day: 'numeric', month: 'long', year: 'numeric' });
	}

	onMount(async () => {
		const res = await api<CertificateDetail>('GET', `/certificates/${certId}`);
		if (res.data) {
			cert = res.data;
		} else {
			error = res.error?.message || 'Certificate not found';
		}
		loading = false;
		fitToScreen();
		window.addEventListener('resize', fitToScreen);
		ready = true;
		setTimeout(() => window.print(), 600);
	});
</script>

<div class="no-print toolbar">
	<button onclick={() => window.print()} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">Print</button>
	<button onclick={() => history.back()} class="px-4 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50 transition-colors">Back</button>
</div>

{#if loading}
	<div class="flex h-screen items-center justify-center text-slate-400">Loading certificate...</div>
{:else if error}
	<div class="flex h-screen items-center justify-center text-danger-600">{error}</div>
{:else if cert}
	<div class="stage" class:show={ready}>
		<div class="cert-wrap" style="transform: scale({pageScale}); height: {297 * pageScale}mm;">
			<div class="cert-page">

				<!-- isometric grid -->
				<svg class="bg-grid" aria-hidden="true">
					<defs>
						<pattern id="iso" width="60" height="103.92" patternUnits="userSpaceOnUse" patternTransform="scale(0.7)">
							<g fill="none" stroke="#1B2A52" stroke-width="0.8">
								<path d="M60 0 L30 51.96 L0 0"/>
								<path d="M60 103.92 L30 51.96 L0 103.92"/>
								<path d="M30 0 L30 103.92"/>
							</g>
						</pattern>
					</defs>
					<rect width="100%" height="100%" fill="url(#iso)" opacity="0.5"/>
				</svg>

				<!-- low-poly overlay: top -->
				<svg class="poly poly-t" viewBox="0 0 210 100" preserveAspectRatio="none" aria-hidden="true">
					<polygon points="0,0 0,60 75,22 130,52 210,12 210,0" fill="#0D1B3E"/>
					<polygon points="0,0 0,40 65,16 115,36 210,5 210,0" fill="#16224A"/>
					<polygon points="0,0 0,24 55,10 100,24 210,0 210,0" fill="#1E2E5C"/>
					<polygon points="0,0 210,0 210,6 150,2 0,15" fill="#22D3EE" opacity="0.22"/>
				</svg>

				<!-- low-poly overlay: bottom-right -->
				<svg class="poly poly-br" viewBox="0 0 210 100" preserveAspectRatio="none" aria-hidden="true">
					<polygon points="0,0 0,50 60,22 120,42 210,0 210,0" fill="#0D1B3E"/>
					<polygon points="0,0 0,32 55,12 105,26 210,0 210,0" fill="#16224A"/>
					<polygon points="0,0 0,18 50,7 95,16 210,0 210,0" fill="#1E2E5C"/>
					<polygon points="210,0 210,6 150,2 0,12 0,0" fill="#22D3EE" opacity="0.18"/>
				</svg>

				<!-- faint corner glows -->
				<div class="glow glow-tl"></div>
				<div class="glow glow-br"></div>

				<div class="cert-inner">
					<div class="logo-row">
						<div class="logo-side">
							<div class="logo-disc"><img src="/logos/karnataka-emblem.png" alt="Government of Karnataka" class="logo" /></div>
						</div>
						<div class="header">
							<div class="kareis">KARNATAKA RESIDENTIAL EDUCATIONAL INSTITUTIONS SOCIETY</div>
							<div class="school">MORARJI DESAI RESIDENTIAL SCHOOL (SC-32)</div>
							<div class="school-en">BAHADDURGHATTA, CHITRADURGA</div>
						</div>
						<div class="logo-side">
							<div class="logo-disc"><img src="/logos/kreis-logo.png" alt="KREIS" class="logo" /></div>
						</div>
					</div>
					<div class="head-divider"><span></span><i></i><span></span></div>

					<!-- digital emblem -->
					<svg class="emblem" viewBox="0 0 120 120" aria-hidden="true">
						<defs>
							<radialGradient id="emGlow" cx="50%" cy="50%" r="50%">
								<stop offset="0%" stop-color="#22D3EE" stop-opacity="0.5"/>
								<stop offset="55%" stop-color="#38BDF8" stop-opacity="0.16"/>
								<stop offset="100%" stop-color="#38BDF8" stop-opacity="0"/>
							</radialGradient>
							<linearGradient id="platinum" x1="0" y1="0" x2="1" y2="1">
								<stop offset="0%" stop-color="#F8FAFE"/>
								<stop offset="40%" stop-color="#D3DBEA"/>
								<stop offset="100%" stop-color="#93A5C4"/>
							</linearGradient>
							<linearGradient id="hexBody" x1="0" y1="0" x2="0" y2="1">
								<stop offset="0%" stop-color="#182548"/>
								<stop offset="100%" stop-color="#0B1226"/>
							</linearGradient>
						</defs>
						<circle cx="60" cy="60" r="55" fill="url(#emGlow)"/>
						<polygon points="60,13 102,36.5 102,83.5 60,107 18,83.5 18,36.5" fill="url(#hexBody)" stroke="#38BDF8" stroke-width="1.2"/>
						<polygon points="60,22 94,41.5 94,78.5 60,98 26,78.5 26,41.5" fill="none" stroke="#22D3EE" stroke-width="0.7" opacity="0.5"/>
						<g stroke="#2E4B8C" stroke-width="0.6" opacity="0.9">
							<line x1="60" y1="60" x2="102" y2="36.5"/>
							<line x1="60" y1="60" x2="102" y2="83.5"/>
							<line x1="60" y1="60" x2="60" y2="107"/>
							<line x1="60" y1="60" x2="18" y2="83.5"/>
							<line x1="60" y1="60" x2="18" y2="36.5"/>
							<line x1="60" y1="60" x2="60" y2="13"/>
						</g>
						<path d="M60 40 L72 60 L60 80 L48 60 Z" fill="url(#platinum)"/>
						<path d="M60 40 L60 60 L72 60 Z" fill="#FFFFFF" opacity="0.9"/>
						<text x="60" y="57" text-anchor="middle" font-family="Montserrat, sans-serif" font-weight="700" font-size="13" fill="#F1F5FB">MDRS</text>
					</svg>

					<div class="title">
						<div class="title-main">CERTIFICATE</div>
						<div class="title-sub">OF ACHIEVEMENT</div>
					</div>

					<div class="intro">
						Acknowledging the outstanding achievement of
					</div>

					<div class="student-name">{cert.student_name}</div>

					<div class="body">
						<span class="body-line">for successfully fulfilling the requirements of the</span>
						<span class="event-name">{cert.event?.name || 'Event'}</span>
						<span class="event-meta">
							{categoryEn(cert.event?.category)}
							{#if cert.class_name}<span class="dot">•</span> Class {cert.class_name}{/if}
						</span>
						<span class="body-line">and has been awarded the</span>
						<span class="prize">{positionLabel(cert.position)}</span>
					</div>

					<div class="date-line">
						<span class="date-en">Awarded on</span>
						<span class="date-value">{formatDate(cert.issue_date || cert.event?.held_date)}</span>
						{#if cert.event?.venue}
							<span class="date-sep">•</span>
							<span class="date-en">Venue:</span>
							<span class="date-value">{cert.event.venue}</span>
						{/if}
					</div>

					<div class="sig-sep"></div>

					<div class="signatures">
						{#each cert.signatories as sig, i}
							<div class="signature-block">
								<div class="signature-area">
									{#if sig.signature_url}
										<img src={apiUrl(sig.signature_url)} alt="signature" class="signature-img" />
									{:else}
										<div class="signature-line"></div>
									{/if}
								</div>
								<div class="sig-name">{sig.name}</div>
								<div class="sig-role">{sig.title || ROLE_LABELS[sig.role] || sig.role}</div>
							</div>
						{:else}
							<div class="signature-block">
								<div class="signature-area"><div class="signature-line"></div></div>
								<div class="sig-name">Principal</div>
								<div class="sig-role">Principal</div>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	@import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@500;600;700;800&family=Cinzel+Decorative:wght@700&family=Inter:wght@400;500;600&display=swap');

	:global(html), :global(body) {
		margin: 0;
		padding: 0;
		height: 100%;
		background: #dfe4ea;
		font-family: 'Inter', sans-serif;
		overflow: hidden;
	}

	.stage {
		display: flex;
		justify-content: center;
		align-items: flex-start;
		padding: 20px;
		box-sizing: border-box;
		min-height: 100vh;
	}

	.stage.show {
		opacity: 1;
		transition: opacity 0.3s;
	}

	.cert-wrap {
		transform-origin: top center;
	}

	.toolbar {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 12px;
		padding: 12px;
		position: sticky;
		top: 0;
		z-index: 50;
		background: #0f172a;
	}

	.cert-page {
		position: relative;
		width: 210mm;
		height: 297mm;
		margin: 0;
		background:
			radial-gradient(ellipse 90% 60% at 50% 26%, #101c3d 0%, #0A122A 70%),
			#0A122A;
		overflow: hidden;
		box-shadow: 0 12px 48px rgba(0, 0, 0, 0.55);
	}

	/* isometric grid */
	.bg-grid {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		z-index: 0;
		opacity: 0.7;
	}

	/* low-poly overlays */
	.poly {
		position: absolute;
		z-index: 1;
		pointer-events: none;
	}
	.poly-t {
		top: 0;
		left: 0;
		width: 100%;
		height: 80mm;
	}
	.poly-br {
		bottom: 0;
		right: 0;
		width: 100%;
		height: 80mm;
		transform: rotate(180deg);
	}

	/* corner glows */
	.glow {
		position: absolute;
		z-index: 1;
		pointer-events: none;
		border-radius: 50%;
		filter: blur(30px);
	}
	.glow-tl {
		top: -40mm;
		left: -40mm;
		width: 110mm;
		height: 110mm;
		background: radial-gradient(circle, rgba(34, 211, 238, 0.18), transparent 65%);
	}
	.glow-br {
		bottom: -45mm;
		right: -45mm;
		width: 120mm;
		height: 120mm;
		background: radial-gradient(circle, rgba(56, 189, 248, 0.15), transparent 65%);
	}

	.cert-inner {
		position: relative;
		z-index: 3;
		padding: 14mm 30mm 16mm;
		height: 297mm;
		box-sizing: border-box;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
	}

	.logo-row {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		width: 100%;
		gap: 8mm;
	}
	.logo-side {
		flex-shrink: 0;
		width: 30mm;
		display: flex;
		justify-content: center;
	}
	.logo-disc {
		width: 22mm;
		height: 22mm;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(148, 197, 255, 0.25);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2mm;
	}
	.logo {
		width: 18mm;
		height: 18mm;
		object-fit: contain;
		filter: drop-shadow(0 1px 3px rgba(0, 0, 0, 0.5));
	}

	.header {
		flex: 1;
		text-align: center;
	}
	.kareis {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 8.5pt;
		color: #7DD3FC;
		letter-spacing: 2px;
		text-transform: uppercase;
	}
	.school {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 15pt;
		color: #E8ECF5;
		margin-top: 4px;
		letter-spacing: 1px;
		text-transform: uppercase;
	}
	.school-en {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 8pt;
		color: #7C8FB0;
		letter-spacing: 2px;
		margin-top: 2px;
		text-transform: uppercase;
	}

	.head-divider {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		width: 82%;
		margin-top: 5mm;
	}
	.head-divider span {
		flex: 1;
		height: 1px;
		background: linear-gradient(90deg, transparent, rgba(125, 211, 252, 0.5));
	}
	.head-divider span:last-child {
		background: linear-gradient(90deg, rgba(125, 211, 252, 0.5), transparent);
	}
	.head-divider i {
		width: 5px;
		height: 5px;
		transform: rotate(45deg);
		background: #22D3EE;
		flex-shrink: 0;
		box-shadow: 0 0 8px rgba(34, 211, 238, 0.7);
	}

	/* digital emblem */
	.emblem {
		width: 30mm;
		height: 30mm;
		margin-top: 6mm;
		animation: emblem-breathe 3s ease-in-out infinite;
	}
	@keyframes emblem-breathe {
		0%, 100% { filter: drop-shadow(0 0 4px rgba(34, 211, 238, 0.35)); }
		50% { filter: drop-shadow(0 0 12px rgba(56, 189, 248, 0.65)); }
	}

	.title {
		margin-top: 5.5mm;
		width: 100%;
	}
	.title-main {
		font-family: 'Montserrat', sans-serif;
		font-weight: 800;
		font-size: 23pt;
		letter-spacing: 9px;
		text-transform: uppercase;
		margin-right: -9px;
		background: linear-gradient(180deg, #F7F9FD 0%, #B9C6DE 48%, #7DD3FC 100%);
		-webkit-background-clip: text;
		background-clip: text;
		-webkit-text-fill-color: transparent;
		color: transparent;
	}
	.title-sub {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 10.5pt;
		color: #8EA3C4;
		letter-spacing: 7px;
		text-transform: uppercase;
		margin-top: 2.5mm;
	}

	.intro {
		margin-top: 7mm;
		font-family: 'Inter', sans-serif;
		font-weight: 400;
		font-size: 9.5pt;
		color: #8A9BB8;
		letter-spacing: 0.5px;
	}

	.student-name {
		margin-top: 4mm;
		font-family: 'Cinzel Decorative', serif;
		font-weight: 700;
		font-size: 26pt;
		color: #F1F5FB;
		padding: 0 10mm 4mm;
		border-bottom: 1.5px solid rgba(125, 211, 252, 0.55);
		line-height: 1.4;
		width: 100%;
		text-shadow: 0 0 18px rgba(125, 211, 252, 0.35);
	}

	.body {
		margin-top: 7mm;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 3mm;
		max-width: 150mm;
		line-height: 1.5;
	}
	.body-line {
		font-family: 'Inter', sans-serif;
		font-weight: 400;
		font-size: 9.5pt;
		color: #9AA9C2;
		letter-spacing: 0.3px;
	}
	.event-name {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 12.5pt;
		color: #7DD3FC;
		text-transform: uppercase;
		letter-spacing: 1.5px;
		text-shadow: 0 0 14px rgba(125, 211, 252, 0.45);
	}
	.event-meta {
		font-family: 'Inter', sans-serif;
		font-weight: 500;
		font-size: 8.5pt;
		color: #647CA0;
		text-transform: uppercase;
		letter-spacing: 1.5px;
	}
	.event-meta .dot { color: #22D3EE; margin: 0 3px; }
	.prize {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 11.5pt;
		color: #38BDF8;
		text-transform: uppercase;
		letter-spacing: 1.5px;
		text-shadow: 0 0 14px rgba(56, 189, 248, 0.4);
	}

	.date-line {
		margin-top: 8mm;
		display: flex;
		align-items: baseline;
		gap: 5px;
		font-size: 9.5pt;
	}
	.date-en {
		font-family: 'Inter', sans-serif;
		font-weight: 600;
		color: #9AA9C2;
	}
	.date-value {
		font-family: 'Inter', sans-serif;
		font-style: italic;
		font-weight: 500;
		color: #CBD5E1;
	}
	.date-sep { color: #22D3EE; }

	.sig-sep {
		margin-top: 9mm;
		width: 46%;
		height: 1px;
		background: linear-gradient(90deg, transparent, rgba(125, 211, 252, 0.45), transparent);
	}

	.signatures {
		margin-top: auto;
		display: flex;
		justify-content: space-evenly;
		align-items: flex-end;
		width: 100%;
		padding-top: 5mm;
	}
	.signature-block {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 48mm;
	}
	.signature-area {
		height: 13mm;
		width: 100%;
		display: flex;
		align-items: flex-end;
		justify-content: center;
	}
	.signature-img {
		height: 13mm;
		width: auto;
		object-fit: contain;
		max-width: 44mm;
	}
	.signature-line {
		width: 100%;
		height: 1.2px;
		background: #3E5480;
	}
	.sig-name {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 10pt;
		color: #D9E2F0;
		margin-top: 2.5mm;
	}
	.sig-role {
		font-family: 'Inter', sans-serif;
		font-weight: 500;
		font-size: 8pt;
		color: #647CA0;
		margin-top: 1mm;
		text-transform: uppercase;
		letter-spacing: 1.5px;
	}

	@media print {
		@page {
			size: A4 portrait;
			margin: 0;
		}
		:global(html), :global(body) {
			width: 210mm;
			height: 297mm;
			margin: 0 !important;
			padding: 0 !important;
			background: white;
			overflow: hidden;
		}
		:global(.flex.h-screen),
		:global(.flex),
		:global(.h-screen),
		:global(.overflow-hidden),
		:global(.overflow-y-auto),
		:global(main),
		:global(.p-0) {
			width: 210mm !important;
			height: 297mm !important;
			min-height: 0 !important;
			max-height: 297mm !important;
			overflow: hidden !important;
			padding: 0 !important;
			margin: 0 !important;
		}
		.stage, .stage.show {
			padding: 0;
			min-height: 0;
			display: block;
			width: 210mm;
			height: 297mm;
			overflow: hidden;
		}
		.toolbar, .no-print {
			display: none !important;
		}
		.cert-wrap {
			transform: none !important;
			width: 210mm !important;
			height: 297mm !important;
			margin: 0 !important;
		}
		.cert-page {
			margin: 0;
			width: 210mm;
			height: 297mm;
			overflow: hidden;
			box-shadow: none;
			-webkit-print-color-adjust: exact;
			print-color-adjust: exact;
		}
		.emblem {
			animation: none;
		}
	}
</style>