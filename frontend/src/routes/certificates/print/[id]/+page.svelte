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
		const pageW = (297 / 25.4) * 96;
		const pageH = (210 / 25.4) * 96;
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
		<div class="cert-wrap" style="transform: scale({pageScale}); height: {210 * pageScale}mm;">
			<div class="cert-page">

				<!-- globe / concentric arcs watermark -->
				<svg class="bg-globe" viewBox="0 0 200 200" aria-hidden="true">
					<defs>
						<radialGradient id="globeFill" cx="38%" cy="32%" r="75%">
							<stop offset="0%" stop-color="#f7f9fc"/>
							<stop offset="100%" stop-color="#eef2f7"/>
						</radialGradient>
					</defs>
					<circle cx="100" cy="100" r="96" fill="url(#globeFill)"/>
					<g transform="translate(100,100)" fill="none">
						<circle r="92" stroke="#b8934c" stroke-width="0.6" opacity="0.55"/>
						<circle r="84" stroke="#dfe6ee" stroke-width="0.4"/>
						<circle r="74" stroke="#dfe6ee" stroke-width="0.35"/>
						<circle r="56" stroke="#dfe6ee" stroke-width="0.4"/>
						<circle r="38" stroke="#dfe6ee" stroke-width="0.3"/>
						<ellipse rx="92" ry="38" stroke="#dfe6ee" stroke-width="0.3"/>
						<ellipse rx="38" ry="92" stroke="#dfe6ee" stroke-width="0.3"/>
						<ellipse rx="74" ry="56" stroke="#dfe6ee" stroke-width="0.25"/>
						<line x1="-92" y1="0" x2="92" y2="0" stroke="#dfe6ee" stroke-width="0.25"/>
						{#each Array(24) as _, i}
							<circle cx="0" cy="-56" r="1.3" fill="#b8934c" stroke="none" opacity="0.5"/>
						{/each}
					</g>
				</svg>

				<!-- soft side glows -->
				<div class="bg-glow bg-glow-l"></div>
				<div class="bg-glow bg-glow-r"></div>

				<!-- modern frame -->
				<div class="frame">
					<div class="frame-line outer"></div>
					<div class="frame-line gold"></div>
					<div class="accent-stripe top"></div>
					<div class="accent-stripe bottom"></div>
					<svg class="corner tl" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#a8842c" stroke-linecap="round" stroke-linejoin="round">
							<path d="M16 124 V34 Q16 16 34 16 H124" stroke-width="2.6"/>
							<path d="M30 110 V48 Q30 30 48 30 H110" stroke-width="1"/>
							<circle cx="30" cy="110" r="2.6" fill="#a8842c" stroke="none"/>
							<path d="M44 44 L52 36 L60 44 L52 52 Z" fill="#a8842c" stroke="none"/>
						</g>
					</svg>
					<svg class="corner tr" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#a8842c" stroke-linecap="round" stroke-linejoin="round">
							<path d="M124 124 V34 Q124 16 106 16 H16" stroke-width="2.6"/>
							<path d="M110 110 V48 Q110 30 92 30 H30" stroke-width="1"/>
							<circle cx="110" cy="110" r="2.6" fill="#a8842c" stroke="none"/>
							<path d="M96 44 L88 36 L80 44 L88 52 Z" fill="#a8842c" stroke="none"/>
						</g>
					</svg>
					<svg class="corner bl" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#a8842c" stroke-linecap="round" stroke-linejoin="round">
							<path d="M16 16 V106 Q16 124 34 124 H124" stroke-width="2.6"/>
							<path d="M30 30 V92 Q30 110 48 110 H110" stroke-width="1"/>
							<circle cx="30" cy="30" r="2.6" fill="#a8842c" stroke="none"/>
							<path d="M44 96 L52 104 L60 96 L52 88 Z" fill="#a8842c" stroke="none"/>
						</g>
					</svg>
					<svg class="corner br" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#a8842c" stroke-linecap="round" stroke-linejoin="round">
							<path d="M124 16 V106 Q124 124 106 124 H16" stroke-width="2.6"/>
							<path d="M110 30 V92 Q110 110 92 110 H30" stroke-width="1"/>
							<circle cx="110" cy="30" r="2.6" fill="#a8842c" stroke="none"/>
							<path d="M96 96 L88 104 L80 96 L88 88 Z" fill="#a8842c" stroke="none"/>
						</g>
					</svg>
				</div>

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

					<div class="title">
						<div class="title-rule"><span></span><i></i><span></span></div>
						<div class="title-en">CERTIFICATE OF ACHIEVEMENT</div>
						<div class="title-rule"><span></span><i></i><span></span></div>
					</div>

					<div class="certify">
						<span class="certify-en">This is to certify that</span>
					</div>

					<div class="student-name">{cert.student_name}</div>

					<div class="details">
						<span class="details-en">studying in</span>
						<span class="value class">{cert.class_name || '—'}</span>
						<span class="details-en">has participated in</span>
						<span class="value event">{cert.event?.name || 'Event'}</span>
						<span class="category-en">{categoryEn(cert.event?.category)}</span>
						<span class="details-en">and secured</span>
						<span class="value prize">{positionLabel(cert.position)}</span>
					</div>

					<div class="date-line">
						<span class="date-en">Date:</span>
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
	@import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@500;600;700&family=Great+Vibes&family=Playfair+Display:ital,wght@0,400;0,500;0,600;1,400&family=Inter:wght@400;500;600&display=swap');

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
		width: 297mm;
		height: 210mm;
		margin: 0;
		background:
			radial-gradient(ellipse 120% 90% at 50% 0%, #ffffff 0%, #fbfcfe 55%, #f4f7fb 100%);
		overflow: hidden;
		box-shadow: 0 12px 48px rgba(0, 0, 0, 0.28);
	}

	/* globe / international watermark */
	.bg-globe {
		position: absolute;
		top: 52%;
		left: 50%;
		width: 76%;
		height: 76%;
		transform: translate(-50%, -50%);
		z-index: 0;
		opacity: 0.9;
	}

	/* soft radial glows behind content */
	.bg-glow {
		position: absolute;
		top: 50%;
		width: 55mm;
		height: 55mm;
		border-radius: 50%;
		transform: translateY(-50%);
		z-index: 0;
		pointer-events: none;
	}
	.bg-glow-l {
		left: -18mm;
		background: radial-gradient(circle, rgba(184, 147, 76, 0.10) 0%, transparent 70%);
	}
	.bg-glow-r {
		right: -18mm;
		background: radial-gradient(circle, rgba(224, 93, 51, 0.07) 0%, transparent 70%);
	}

	/* modern frame */
	.frame {
		position: absolute;
		inset: 0;
		z-index: 1;
		pointer-events: none;
	}
	.frame-line.outer {
		position: absolute;
		inset: 4mm;
		border: 1px solid #d8dfe8;
		border-radius: 1mm;
	}
	.frame-line.gold {
		position: absolute;
		inset: 5.6mm;
		border: 1.4px solid #b8934c;
		border-radius: 0.8mm;
		opacity: 0.9;
	}
	.accent-stripe {
		position: absolute;
		left: 5.6mm;
		right: 5.6mm;
		height: 2.4mm;
		background: linear-gradient(90deg, #b8934c, #d8b25e 20%, #dfe6ee 38%, #dfe6ee 62%, #d8b25e 80%, #b8934c);
		border-radius: 1mm;
	}
	.accent-stripe.top { top: 6.8mm; }
	.accent-stripe.bottom { bottom: 6.8mm; }

	.corner {
		position: absolute;
		width: 34mm;
		height: 34mm;
		z-index: 2;
	}
	.corner.tl { top: 1.6mm; left: 1.6mm; }
	.corner.tr { top: 1.6mm; right: 1.6mm; }
	.corner.bl { bottom: 1.6mm; left: 1.6mm; }
	.corner.br { bottom: 1.6mm; right: 1.6mm; }

	.cert-inner {
		position: relative;
		z-index: 3;
		padding: 17mm 26mm 14mm;
		height: 210mm;
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
		width: 38mm;
		display: flex;
		justify-content: center;
	}
	.logo {
		width: 25mm;
		height: 25mm;
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
		font-size: 10.5pt;
		color: #a8842c;
		letter-spacing: 2px;
		text-transform: uppercase;
	}
	.school {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 19pt;
		color: #123b5c;
		margin-top: 5px;
		letter-spacing: 0.5px;
		text-transform: uppercase;
	}
	.school-en {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 9pt;
		color: #7c8ba1;
		letter-spacing: 2px;
		margin-top: 3px;
		text-transform: uppercase;
	}

	.title {
		margin-top: 9mm;
		width: 100%;
	}
	.title-rule {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		margin: 3mm auto;
		width: 68%;
	}
	.title-rule span {
		flex: 1;
		height: 1px;
		background: linear-gradient(90deg, transparent, #b8934c);
	}
	.title-rule span:last-child {
		background: linear-gradient(90deg, #b8934c, transparent);
	}
	.title-rule i {
		width: 6px;
		height: 6px;
		transform: rotate(45deg);
		background: #b8934c;
		flex-shrink: 0;
	}
	.title-en {
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
		font-size: 13.5pt;
		color: #a8842c;
		letter-spacing: 7px;
		margin-top: 3px;
		text-transform: uppercase;
	}

	.certify {
		margin-top: 7.5mm;
		display: flex;
		flex-direction: column;
		gap: 1px;
	}
	.certify-en {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 11pt;
		color: #42536b;
	}

	.student-name {
		margin-top: 4.5mm;
		font-family: 'Great Vibes', cursive;
		font-size: 29pt;
		color: #123b5c;
		padding: 0 26mm 3mm;
		border-bottom: 1.2px solid #b8934c;
		line-height: 1.5;
	}

	.details {
		margin-top: 7mm;
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		align-items: baseline;
		gap: 4px 8px;
		max-width: 240mm;
		line-height: 1.6;
	}
	.details-en {
		font-family: 'Playfair Display', serif;
		font-size: 10.5pt;
		color: #42536b;
	}
	.category-en {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 9.5pt;
		color: #7c8ba1;
	}
	.value {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		font-size: 11.5pt;
		color: #26364a;
	}
	.value.class { color: #123b5c; }
	.value.event {
		color: #123b5c;
		font-size: 12.5pt;
	}
	.value.prize {
		color: #e05d33;
		font-size: 12.5pt;
	}

	.date-line {
		margin-top: 6mm;
		display: flex;
		align-items: baseline;
		gap: 6px;
		font-size: 10pt;
	}
	.date-en {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		color: #42536b;
	}
	.date-value {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		color: #26364a;
	}
	.date-sep { color: #b8934c; }

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
		width: 52mm;
	}
	.signature-area {
		height: 12mm;
		width: 100%;
		display: flex;
		align-items: flex-end;
		justify-content: center;
	}
	.signature-img {
		height: 12mm;
		width: auto;
		object-fit: contain;
		max-width: 46mm;
	}
	.signature-line {
		width: 100%;
		height: 1.2px;
		background: #b3bfcc;
	}
	.sig-name {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		font-size: 9pt;
		color: #26364a;
		margin-top: 2.5mm;
	}
	.sig-role {
		font-family: 'Montserrat', sans-serif;
		font-weight: 500;
		font-size: 7.5pt;
		color: #7c8ba1;
		margin-top: 1mm;
		text-transform: uppercase;
		letter-spacing: 1.5px;
	}

	@media print {
		@page {
			size: A4 landscape;
			margin: 0;
		}
		:global(html), :global(body) {
			width: 297mm;
			height: 210mm;
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
			width: 297mm !important;
			height: 210mm !important;
			min-height: 0 !important;
			max-height: 210mm !important;
			overflow: hidden !important;
			padding: 0 !important;
			margin: 0 !important;
		}
		.stage, .stage.show {
			padding: 0;
			min-height: 0;
			display: block;
			width: 297mm;
			height: 210mm;
			overflow: hidden;
		}
		.toolbar, .no-print {
			display: none !important;
		}
		.cert-wrap {
			transform: none !important;
			width: 297mm !important;
			height: 210mm !important;
			margin: 0 !important;
		}
		.cert-page {
			margin: 0;
			width: 297mm;
			height: 210mm;
			overflow: hidden;
			box-shadow: none;
			-webkit-print-color-adjust: exact;
			print-color-adjust: exact;
		}
	}
</style>