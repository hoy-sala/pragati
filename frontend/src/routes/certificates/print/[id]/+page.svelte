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
		document.documentElement.classList.add('cert-print');
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

				<!-- hairline frame -->
				<div class="frame">
					<div class="frame-line outer"></div>
					<div class="frame-line inner"></div>
					<i class="bracket tl"></i><i class="bracket tr"></i><i class="bracket bl"></i><i class="bracket br"></i>

					<svg class="corner tl" viewBox="0 0 80 80" aria-hidden="true">
						<path d="M10 72 V26 Q10 10 26 10 H72" fill="none" stroke="#B8860B" stroke-width="1"/>
						<path d="M18 72 V32 Q18 18 32 18 H72" fill="none" stroke="#D4AF37" stroke-width="0.8"/>
						<path d="M42 3 L44 8 L49 8 L45 12 L47 17 L42 14 L37 17 L39 12 L35 8 L40 8 Z" fill="#D4AF37"/>
					</svg>
					<svg class="corner tr" viewBox="0 0 80 80" aria-hidden="true">
						<path d="M70 72 V26 Q70 10 54 10 H8" fill="none" stroke="#B8860B" stroke-width="1"/>
						<path d="M62 72 V32 Q62 18 48 18 H8" fill="none" stroke="#D4AF37" stroke-width="0.8"/>
						<path d="M38 3 L40 8 L45 8 L41 12 L43 17 L38 14 L33 17 L35 12 L31 8 L36 8 Z" fill="#D4AF37"/>
					</svg>
					<svg class="corner bl" viewBox="0 0 80 80" aria-hidden="true">
						<path d="M10 8 V54 Q10 70 26 70 H72" fill="none" stroke="#B8860B" stroke-width="1"/>
						<path d="M18 8 V48 Q18 62 32 62 H72" fill="none" stroke="#D4AF37" stroke-width="0.8"/>
						<path d="M42 77 L44 72 L49 72 L45 68 L47 63 L42 66 L37 63 L39 68 L35 72 L40 72 Z" fill="#D4AF37"/>
					</svg>
					<svg class="corner br" viewBox="0 0 80 80" aria-hidden="true">
						<path d="M70 8 V54 Q70 70 54 70 H8" fill="none" stroke="#B8860B" stroke-width="1"/>
						<path d="M62 8 V48 Q62 62 48 62 H8" fill="none" stroke="#D4AF37" stroke-width="0.8"/>
						<path d="M38 77 L40 72 L45 72 L41 68 L43 63 L38 66 L33 63 L35 68 L31 72 L36 72 Z" fill="#D4AF37"/>
					</svg>

					<i class="dot-tick top"></i><i class="dot-tick right"></i><i class="dot-tick bottom"></i><i class="dot-tick left"></i>
				</div>

				<!-- watermark -->
				<img src="/logos/kreis-logo.png" alt="" class="watermark" aria-hidden="true" />

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

					<!-- medallion -->
					<div class="medallion-row">
						<svg class="laurel left" viewBox="0 0 110 44" aria-hidden="true">
							<path d="M6 22 C 40 6, 74 10, 104 22" fill="none" stroke="#B8860B" stroke-width="1"/>
							<g fill="#D4AF37">
								<path d="M18 18 q 8 -3 10 3 q -7 4 -10 -3 Z"/>
								<path d="M22 24 q 8 -2 12 3 q -7 5 -12 -3 Z"/>
								<path d="M38 16 q 8 -3 11 2 q -7 5 -11 -2 Z"/>
								<path d="M44 25 q 8 -1 12 4 q -8 4 -12 -4 Z"/>
								<path d="M60 15 q 8 -2 10 3 q -7 4 -10 -3 Z"/>
								<path d="M66 26 q 8 0 12 4 q -8 4 -12 -4 Z"/>
								<path d="M82 16 q 7 -2 9 3 q -6 3 -9 -3 Z"/>
								<path d="M88 24 q 7 0 10 3 q -7 3 -10 -3 Z"/>
							</g>
							<path d="M104 22 l 4 -2 l -1 4 l -3 -2 Z" fill="#D4AF37"/>
						</svg>
						<svg class="medallion" viewBox="0 0 120 120" aria-hidden="true">
							<circle cx="60" cy="60" r="58" fill="none" stroke="#0B2545" stroke-width="1"/>
							<circle cx="60" cy="60" r="52" fill="none" stroke="#D4AF37" stroke-width="1.4"/>
							<circle cx="60" cy="60" r="47" fill="none" stroke="#0B2545" stroke-width="0.5" opacity="0.6"/>
							<path d="M60 16 L63 26 L73 26 L65 32 L68 42 L60 36 L52 42 L55 32 L47 26 L57 26 Z" fill="#D4AF37"/>
							<text x="60" y="70" text-anchor="middle" font-family="Playfair Display, serif" font-weight="700" font-size="34" fill="#0B2545">MDRS</text>
							<circle cx="60" cy="100" r="2.5" fill="#D4AF37"/>
						</svg>
						<svg class="laurel right" viewBox="0 0 110 44" aria-hidden="true">
							<path d="M6 22 C 40 6, 74 10, 104 22" fill="none" stroke="#B8860B" stroke-width="1"/>
							<g fill="#D4AF37">
								<path d="M18 18 q 8 -3 10 3 q -7 4 -10 -3 Z"/>
								<path d="M22 24 q 8 -2 12 3 q -7 5 -12 -3 Z"/>
								<path d="M38 16 q 8 -3 11 2 q -7 5 -11 -2 Z"/>
								<path d="M44 25 q 8 -1 12 4 q -8 4 -12 -4 Z"/>
								<path d="M60 15 q 8 -2 10 3 q -7 4 -10 -3 Z"/>
								<path d="M66 26 q 8 0 12 4 q -8 4 -12 -4 Z"/>
								<path d="M82 16 q 7 -2 9 3 q -6 3 -9 -3 Z"/>
								<path d="M88 24 q 7 0 10 3 q -7 3 -10 -3 Z"/>
							</g>
							<path d="M104 22 l 4 -2 l -1 4 l -3 -2 Z" fill="#D4AF37"/>
						</svg>
					</div>

					<div class="title">
						<div class="title-main">Certificate</div>
						<div class="title-sub">of Achievement</div>
						<div class="title-rule"><span></span><i></i><span></span></div>
					</div>

					<div class="intro">
						This is to certify that
					</div>

					<div class="student-name">{cert.student_name}</div>

					<div class="body">
						<span class="body-line">has successfully participated in</span>
						<span class="event-name">{cert.event?.name || 'Event'}</span>
						<span class="event-meta">
							{categoryEn(cert.event?.category)}
							{#if cert.class_name}<span class="dot">•</span> Class {cert.class_name}{/if}
						</span>
						<span class="body-line">and has been awarded the</span>
						<span class="prize">{positionLabel(cert.position)}</span>
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
	@import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@500;600;700;800&family=Playfair+Display:ital,wght@0,500;0,600;0,700;1,400;1,500&family=Inter:wght@400;500;600&display=swap');

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
			radial-gradient(130% 100% at 50% 0%, #ffffff 0%, #FDFBF6 50%, #F5F1E6 100%);
		overflow: hidden;
		box-shadow: 0 12px 48px rgba(0, 0, 0, 0.25);
	}

	/* hairline double frame with corner brackets */
	.frame {
		position: absolute;
		inset: 0;
		z-index: 1;
		pointer-events: none;
	}
	.frame-line.outer {
		position: absolute;
		inset: 4mm;
		border: 0.8px solid #0B2545;
		opacity: 0.85;
	}
	.frame-line.inner {
		position: absolute;
		inset: 5.6mm;
		border: 0.5px solid #D4AF37;
		opacity: 0.9;
	}
	.bracket {
		position: absolute;
		width: 7mm;
		height: 7mm;
		border: 0.9px solid #D4AF37;
	}
	.bracket.tl { top: 5.4mm; left: 5.4mm; border-right: none; border-bottom: none; }
	.bracket.tr { top: 5.4mm; right: 5.4mm; border-left: none; border-bottom: none; }
	.bracket.bl { bottom: 5.4mm; left: 5.4mm; border-right: none; border-top: none; }
	.bracket.br { bottom: 5.4mm; right: 5.4mm; border-left: none; border-top: none; }

	/* decorative corner flourishes */
	.corner {
		position: absolute;
		width: 30mm;
		height: 30mm;
		z-index: 2;
	}
	.corner.tl { top: 5.6mm; left: 5.6mm; }
	.corner.tr { top: 5.6mm; right: 5.6mm; }
	.corner.bl { bottom: 5.6mm; left: 5.6mm; }
	.corner.br { bottom: 5.6mm; right: 5.6mm; }

	/* gold diamonds on the frame midpoints */
	.dot-tick {
		position: absolute;
		width: 5px;
		height: 5px;
		background: #D4AF37;
		transform: rotate(45deg);
		z-index: 2;
	}
	.dot-tick.top { top: 4mm; left: 50%; margin-left: -2.5px; }
	.dot-tick.bottom { bottom: 4mm; left: 50%; margin-left: -2.5px; }
	.dot-tick.left { left: 4mm; top: 50%; margin-top: -2.5px; }
	.dot-tick.right { right: 4mm; top: 50%; margin-top: -2.5px; }

	/* KREIS watermark */
	.watermark {
		position: absolute;
		top: 50%;
		left: 50%;
		width: 104mm;
		height: auto;
		transform: translate(-50%, -50%);
		z-index: 0;
		opacity: 0.07;
	}

	.cert-inner {
		position: relative;
		z-index: 3;
		padding: 15mm 32mm 16mm;
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
		width: 19mm;
		height: 19mm;
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
		font-size: 8.5pt;
		color: #A67C1E;
		letter-spacing: 2.5px;
		text-transform: uppercase;
	}
	.school {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 16pt;
		color: #0B2545;
		margin-top: 4px;
		letter-spacing: 0.5px;
		text-transform: uppercase;
	}
	.school-en {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 8pt;
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
		margin-top: 5mm;
	}
	.head-divider span {
		flex: 1;
		height: 1px;
		background: linear-gradient(90deg, transparent, rgba(184, 147, 76, 0.75));
	}
	.head-divider span:last-child {
		background: linear-gradient(90deg, rgba(184, 147, 76, 0.75), transparent);
	}
	.head-divider i {
		width: 5px;
		height: 5px;
		transform: rotate(45deg);
		background: #D4AF37;
		flex-shrink: 0;
	}

	/* medallion with laurel branches */
	.medallion-row {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 1mm;
		margin-top: 5.5mm;
		width: 100%;
	}
	.medallion {
		width: 24mm;
		height: 24mm;
		margin: 0;
		filter: drop-shadow(0 2px 6px rgba(11, 37, 69, 0.18));
	}
	.laurel {
		width: 42mm;
		height: 17mm;
		flex-shrink: 0;
		opacity: 0.85;
	}
	.laurel.right {
		transform: scaleX(-1);
	}

	.title {
		margin-top: 5mm;
		width: 100%;
	}
	.title-main {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 22pt;
		color: #0B2545;
		letter-spacing: 3px;
		text-transform: uppercase;
	}
	.title-sub {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		font-size: 10pt;
		color: #A67C1E;
		letter-spacing: 5px;
		text-transform: uppercase;
		margin-top: 1.5mm;
	}
	.title-rule {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		margin: 3mm auto 0;
		width: 60%;
	}
	.title-rule span {
		flex: 1;
		height: 1px;
		background: linear-gradient(90deg, transparent, rgba(184, 147, 76, 0.8));
	}
	.title-rule span:last-child {
		background: linear-gradient(90deg, rgba(184, 147, 76, 0.8), transparent);
	}
	.title-rule i {
		width: 6px;
		height: 6px;
		transform: rotate(45deg);
		background: #D4AF37;
		flex-shrink: 0;
	}

	.intro {
		margin-top: 7mm;
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 11pt;
		color: #42536B;
	}

	.student-name {
		margin-top: 4mm;
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 32pt;
		color: #0B2545;
		padding: 0 10mm 4mm;
		border-bottom: 1px solid #D4AF37;
		line-height: 1.4;
		width: 100%;
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
		color: #33415c;
		letter-spacing: 0.3px;
	}
	.event-name {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 14pt;
		color: #0B2545;
		letter-spacing: 0.5px;
	}
	.event-meta {
		font-family: 'Montserrat', sans-serif;
		font-weight: 500;
		font-size: 8pt;
		color: #8a97ad;
		text-transform: uppercase;
		letter-spacing: 1.5px;
	}
	.event-meta .dot { color: #D4AF37; margin: 0 3px; }
	.prize {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 13pt;
		color: #A67C1E;
		letter-spacing: 0.5px;
	}

	.date-line {
		margin-top: 8mm;
		display: flex;
		align-items: baseline;
		gap: 5px;
		font-size: 9.5pt;
	}
	.date-en {
		font-family: 'Montserrat', sans-serif;
		font-weight: 600;
		color: #33415c;
	}
	.date-value {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-weight: 500;
		color: #33415c;
	}
	.date-sep { color: #D4AF37; }

	.sig-sep {
		margin-top: 9mm;
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
		:global(html.cert-print), :global(html.cert-print body) {
			width: 210mm;
			height: 297mm;
			margin: 0 !important;
			padding: 0 !important;
			background: white;
			overflow: hidden;
		}
		:global(html.cert-print .flex.h-screen),
		:global(html.cert-print .flex),
		:global(html.cert-print .h-screen),
		:global(html.cert-print .overflow-hidden),
		:global(html.cert-print .overflow-y-auto),
		:global(html.cert-print main),
		:global(html.cert-print .p-0) {
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