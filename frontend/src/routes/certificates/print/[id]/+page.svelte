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

				<!-- layered corner waves -->
				<svg class="corner tl" viewBox="0 0 320 320" preserveAspectRatio="xMinYMin meet" aria-hidden="true">
					<path d="M0 320 C90 320 190 250 190 160 C190 80 250 0 320 0 L0 0 Z" fill="#0B2545"/>
					<path d="M0 320 C120 320 230 240 230 150 C230 70 280 10 320 10 L0 10 Z" fill="#2563EB"/>
					<path d="M0 320 C150 320 265 225 265 135 C265 55 300 20 320 20 L0 20 Z" fill="#93C5FD"/>
					<path d="M0 320 C175 320 300 215 300 120 C300 45 315 30 320 30 L0 30 Z" fill="#E5E7EB"/>
				</svg>
				<svg class="corner br" viewBox="0 0 320 320" preserveAspectRatio="xMaxYMax meet" aria-hidden="true">
					<path d="M0 320 C90 320 190 250 190 160 C190 80 250 0 320 0 L0 0 Z" fill="#0B2545"/>
					<path d="M0 320 C120 320 230 240 230 150 C230 70 280 10 320 10 L0 10 Z" fill="#2563EB"/>
					<path d="M0 320 C150 320 265 225 265 135 C265 55 300 20 320 20 L0 20 Z" fill="#93C5FD"/>
					<path d="M0 320 C175 320 300 215 300 120 C300 45 315 30 320 30 L0 30 Z" fill="#E5E7EB"/>
				</svg>

				<!-- dotted accent lines -->
				<div class="dots dots-l"></div>
				<div class="dots dots-r"></div>

				<div class="cert-inner">
					<div class="logo-row">
						<div class="logo-side">
							<img src="/logos/karnataka-emblem.png" alt="Government of Karnataka" class="logo" />
						</div>
						<div class="header">
							<div class="kareis">KARNATAKA RESIDENTIAL EDUCATIONAL INSTITUTIONS SOCIETY</div>
							<div class="school">MORARJI DESAI RESIDENTIAL SCHOOL (SC-32)</div>
							<div class="school-en">BAHADDURGHATTA, CHITRADURGA</div>
						</div>
						<div class="logo-side">
							<img src="/logos/kreis-logo.png" alt="KREIS" class="logo" />
						</div>
					</div>

					<!-- golden seal -->
					<svg class="seal" viewBox="0 0 120 120" aria-hidden="true">
						<defs>
							<radialGradient id="sealGold" cx="35%" cy="30%" r="80%">
								<stop offset="0%" stop-color="#F7D774"/>
								<stop offset="45%" stop-color="#E6B325"/>
								<stop offset="100%" stop-color="#A67C1E"/>
							</radialGradient>
						</defs>
						<circle cx="60" cy="60" r="53" fill="none" stroke="#E6B325" stroke-width="13" stroke-dasharray="8 6.3"/>
						<circle cx="60" cy="60" r="53" fill="none" stroke="#C99B2A" stroke-width="13" stroke-dasharray="8 6.3" transform="rotate(10 60 60)"/>
						<circle cx="60" cy="60" r="42" fill="url(#sealGold)" stroke="#8a6a14" stroke-width="1.2"/>
						<circle cx="60" cy="60" r="33" fill="none" stroke="#A67C1E" stroke-width="1"/>
						<circle cx="60" cy="60" r="26" fill="rgba(255,255,255,0.35)"/>
						<path d="M60 22 L64.5 39 L82 39 L68 50 L73 68 L60 57 L47 68 L52 50 L38 39 L55.5 39 Z" fill="#8a6a14"/>
					</svg>

					<div class="title">
						<div class="title-main">CERTIFICATE</div>
						<div class="title-sub">OF ACHIEVEMENT</div>
					</div>

					<div class="intro">
						This certificate is proudly presented to
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
	@import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@500;600;700;800&family=Playfair+Display:ital,wght@0,500;0,600;0,700;1,400&family=Inter:wght@400;500;600&display=swap');

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
		background: #ffffff;
		overflow: hidden;
		box-shadow: 0 12px 48px rgba(0, 0, 0, 0.28);
	}

	/* layered corner waves */
	.corner {
		position: absolute;
		z-index: 1;
		pointer-events: none;
	}
	.corner.tl {
		top: 0;
		left: 0;
		width: 115mm;
		height: 115mm;
	}
	.corner.br {
		bottom: 0;
		right: 0;
		width: 115mm;
		height: 115mm;
		transform: rotate(180deg);
	}

	/* dotted accent lines */
	.dots {
		position: absolute;
		top: 70mm;
		bottom: 70mm;
		width: 6px;
		z-index: 2;
		background-image: radial-gradient(circle, rgba(37, 99, 235, 0.5) 2.2px, transparent 3px);
		background-size: 6px 15px;
		background-repeat: repeat-y;
		pointer-events: none;
	}
	.dots-l { left: 11mm; }
	.dots-r { right: 11mm; }

	.cert-inner {
		position: relative;
		z-index: 3;
		padding: 15mm 30mm 16mm;
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
	.logo {
		width: 20mm;
		height: 20mm;
		object-fit: contain;
		filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.12));
	}

	.header {
		flex: 1;
		text-align: center;
	}
	.kareis {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 9.5pt;
		color: #a8842c;
		letter-spacing: 2px;
		text-transform: uppercase;
	}
	.school {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 17pt;
		color: #0B2545;
		margin-top: 4px;
		letter-spacing: 0.5px;
		text-transform: uppercase;
	}
	.school-en {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 8.5pt;
		color: #7c8ba1;
		letter-spacing: 2px;
		margin-top: 2px;
		text-transform: uppercase;
	}

	/* golden seal */
	.seal {
		width: 30mm;
		height: 30mm;
		margin-top: 8mm;
		filter: drop-shadow(0 3px 6px rgba(166, 124, 30, 0.45));
	}

	.title {
		margin-top: 7mm;
		width: 100%;
	}
	.title-main {
		font-family: 'Montserrat', sans-serif;
		font-weight: 800;
		font-size: 26pt;
		color: #0B2545;
		letter-spacing: 12px;
		text-transform: uppercase;
		margin-right: -12px;
	}
	.title-sub {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 12pt;
		color: #1E293B;
		letter-spacing: 8px;
		text-transform: uppercase;
		margin-top: 3mm;
	}

	.intro {
		margin-top: 9mm;
		font-family: 'Inter', sans-serif;
		font-weight: 400;
		font-size: 10pt;
		color: #334155;
	}

	.student-name {
		margin-top: 5mm;
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 31pt;
		color: #1D4ED8;
		padding: 0 10mm 4mm;
		border-bottom: 1.2px solid #1D4ED8;
		line-height: 1.4;
		width: 100%;
	}

	.body {
		margin-top: 9mm;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 3.5mm;
		max-width: 150mm;
		line-height: 1.5;
	}
	.body-line {
		font-family: 'Inter', sans-serif;
		font-weight: 400;
		font-size: 9.5pt;
		color: #1E293B;
	}
	.event-name {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 13pt;
		color: #0B2545;
		text-transform: uppercase;
		letter-spacing: 1px;
	}
	.event-meta {
		font-family: 'Inter', sans-serif;
		font-weight: 500;
		font-size: 8.5pt;
		color: #64748B;
		text-transform: uppercase;
		letter-spacing: 1.5px;
	}
	.event-meta .dot { color: #1D4ED8; margin: 0 3px; }
	.prize {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 11.5pt;
		color: #2563EB;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.date-line {
		margin-top: 9mm;
		display: flex;
		align-items: baseline;
		gap: 5px;
		font-size: 9.5pt;
	}
	.date-en {
		font-family: 'Inter', sans-serif;
		font-weight: 600;
		color: #0F172A;
	}
	.date-value {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-weight: 600;
		color: #0F172A;
	}
	.date-sep { color: #1D4ED8; }

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
		background: #b3bfcc;
	}
	.sig-name {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 10pt;
		color: #0B2545;
		margin-top: 2.5mm;
	}
	.sig-role {
		font-family: 'Inter', sans-serif;
		font-weight: 500;
		font-size: 8pt;
		color: #64748B;
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
	}
</style>