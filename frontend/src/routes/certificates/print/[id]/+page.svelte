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
					<div class="head-divider"><span></span><i></i><span></span></div>

					<!-- golden seal -->
					<svg class="seal" viewBox="0 0 120 120" aria-hidden="true">
						<defs>
							<radialGradient id="sealGold" cx="35%" cy="30%" r="80%">
								<stop offset="0%" stop-color="#FBE9B0"/>
								<stop offset="35%" stop-color="#F5D061"/>
								<stop offset="70%" stop-color="#E6B325"/>
								<stop offset="100%" stop-color="#A67C1E"/>
							</radialGradient>
						</defs>
						<circle cx="60" cy="60" r="55" fill="none" stroke="#D4AF37" stroke-width="1" opacity="0.6"/>
						<circle cx="60" cy="60" r="52" fill="none" stroke="#E6B325" stroke-width="12" stroke-dasharray="7 6.4"/>
						<circle cx="60" cy="60" r="52" fill="none" stroke="#C99B2A" stroke-width="12" stroke-dasharray="7 6.4" transform="rotate(10 60 60)"/>
						<circle cx="60" cy="60" r="41" fill="url(#sealGold)" stroke="#8a6a14" stroke-width="1"/>
						<circle cx="60" cy="60" r="35" fill="none" stroke="#A67C1E" stroke-width="0.8"/>
						<circle cx="60" cy="60" r="28" fill="rgba(255,255,255,0.4)"/>
						<circle cx="60" cy="60" r="18" fill="none" stroke="#A67C1E" stroke-width="0.8"/>
						<path d="M60 24 L64 38 L79 38 L67 47 L72 62 L60 53 L48 62 L53 47 L41 38 L56 38 Z" fill="#A67C1E"/>
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
		background:
			radial-gradient(ellipse 85% 55% at 50% 28%, rgba(37, 99, 235, 0.04), transparent 72%),
			#ffffff;
		overflow: hidden;
		box-shadow: 0 12px 48px rgba(0, 0, 0, 0.28);
	}
	/* elegant double frame */
	.cert-page::before {
		content: '';
		position: absolute;
		inset: 5mm;
		border: 1px solid rgba(184, 147, 76, 0.55);
		border-radius: 1.5mm;
		z-index: 1;
		pointer-events: none;
	}
	.cert-page::after {
		content: '';
		position: absolute;
		inset: 6.2mm;
		border: 0.5px solid rgba(11, 37, 69, 0.18);
		border-radius: 1.5mm;
		z-index: 1;
		pointer-events: none;
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
		width: 110mm;
		height: 110mm;
	}
	.corner.br {
		bottom: 0;
		right: 0;
		width: 110mm;
		height: 110mm;
		transform: rotate(180deg);
	}

	/* dotted accent lines */
	.dots {
		position: absolute;
		top: 72mm;
		bottom: 72mm;
		width: 4px;
		z-index: 2;
		background-image: radial-gradient(circle, rgba(37, 99, 235, 0.55) 1.6px, transparent 2.2px);
		background-size: 4px 16px;
		background-repeat: repeat-y;
		pointer-events: none;
	}
	.dots-l { left: 8mm; }
	.dots-r { right: 8mm; }

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

	.head-divider {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		width: 82%;
		margin-top: 6mm;
	}
	.head-divider span {
		flex: 1;
		height: 1px;
		background: linear-gradient(90deg, transparent, rgba(184, 147, 76, 0.7));
	}
	.head-divider span:last-child {
		background: linear-gradient(90deg, rgba(184, 147, 76, 0.7), transparent);
	}
	.head-divider i {
		width: 6px;
		height: 6px;
		transform: rotate(45deg);
		background: #b8934c;
		flex-shrink: 0;
		box-shadow: 0 0 0 2px rgba(184, 147, 76, 0.15);
	}

	/* golden seal */
	.seal {
		width: 27mm;
		height: 27mm;
		margin-top: 7mm;
		filter: drop-shadow(0 3px 8px rgba(166, 124, 30, 0.4));
	}

	.title {
		margin-top: 6mm;
		width: 100%;
	}
	.title-main {
		font-family: 'Montserrat', sans-serif;
		font-weight: 800;
		font-size: 24pt;
		color: #0B2545;
		letter-spacing: 10px;
		text-transform: uppercase;
		margin-right: -10px;
	}
	.title-sub {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 11pt;
		color: #1E293B;
		letter-spacing: 7px;
		text-transform: uppercase;
		margin-top: 2.5mm;
	}

	.intro {
		margin-top: 8mm;
		font-family: 'Inter', sans-serif;
		font-weight: 400;
		font-size: 10pt;
		color: #334155;
		letter-spacing: 0.5px;
	}

	.student-name {
		margin-top: 4.5mm;
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 30pt;
		color: #1D4ED8;
		padding: 0 10mm 4mm;
		border-bottom: 1.5px solid #1D4ED8;
		line-height: 1.4;
		width: 100%;
	}

	.body {
		margin-top: 8mm;
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
		color: #1E293B;
		letter-spacing: 0.3px;
	}
	.event-name {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 12.5pt;
		color: #0B2545;
		text-transform: uppercase;
		letter-spacing: 1.5px;
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
		letter-spacing: 1.5px;
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
		color: #0F172A;
	}
	.date-value {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-weight: 600;
		color: #0F172A;
	}
	.date-sep { color: #1D4ED8; }

	.sig-sep {
		margin-top: 10mm;
		width: 46%;
		height: 1px;
		background: linear-gradient(90deg, transparent, rgba(184, 147, 76, 0.6), transparent);
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
		background: #c3cdd9;
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