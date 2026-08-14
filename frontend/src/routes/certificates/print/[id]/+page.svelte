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

	function positionKn(pos: string): string {
		if (pos === '1st') return 'ಪ್ರಥಮ ಬಹುಮಾನ';
		if (pos === '2nd') return 'ದ್ವಿತೀಯ ಬಹುಮಾನ';
		if (pos === '3rd') return 'ತೃತೀಯ ಬಹುಮಾನ';
		return 'ಭಾಗವಹಿಸುವಿಕೆ';
	}

	function categoryKn(cat?: string): string {
		if (cat === 'sports') return 'ಕ್ರೀಡಾ ಸ್ಪರ್ಧೆಯಲ್ಲಿ';
		if (cat === 'cultural') return 'ಸಾಂಸ್ಕೃತಿಕ ಸ್ಪರ್ಧೆಯಲ್ಲಿ';
		if (cat === 'academic') return 'ಶೈಕ್ಷಣಿಕ ಸ್ಪರ್ಧೆಯಲ್ಲಿ';
		return 'ಸ್ಪರ್ಧೆಯಲ್ಲಿ';
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
				<!-- subtle guilloché background -->
				<svg class="bg-pattern" viewBox="0 0 120 120" aria-hidden="true">
					<g transform="translate(60,60)">
						<circle r="56" fill="none" stroke="#e6d9b8" stroke-width="0.25"/>
						<circle r="49" fill="none" stroke="#e6d9b8" stroke-width="0.15"/>
						<circle r="42" fill="none" stroke="#e6d9b8" stroke-width="0.2"/>
						<circle r="35" fill="none" stroke="#e6d9b8" stroke-width="0.15"/>
						<circle r="28" fill="none" stroke="#e6d9b8" stroke-width="0.25"/>
						<circle r="21" fill="none" stroke="#e6d9b8" stroke-width="0.15"/>
						{#each Array(60) as _, i}
							<line x1="0" y1="-21" x2="0" y2="-56" transform="rotate({i * 6})" stroke="#e6d9b8" stroke-width="0.1"/>
						{/each}
					</g>
				</svg>

				<!-- ornamental double frame -->
				<div class="frame">
					<div class="frame-outer"></div>
					<div class="frame-inner"></div>
					<svg class="flourish tl" viewBox="0 0 120 120" aria-hidden="true">
						<g fill="none" stroke="#b8934c" stroke-linecap="round">
							<path d="M8 112 V24 Q8 8 24 8 H112" stroke-width="2"/>
							<path d="M20 100 V36 Q20 20 36 20 H100" stroke-width="1"/>
							<path d="M30 90 C18 80 18 42 30 30" stroke-width="0.8"/>
							<path d="M40 80 C32 72 32 48 40 40" stroke-width="0.6"/>
							<circle cx="30" cy="90" r="2.2" fill="#b8934c"/>
							<circle cx="30" cy="30" r="2.2" fill="#b8934c"/>
						</g>
					</svg>
					<svg class="flourish tr" viewBox="0 0 120 120" aria-hidden="true">
						<g fill="none" stroke="#b8934c" stroke-linecap="round">
							<path d="M112 112 V24 Q112 8 96 8 H8" stroke-width="2"/>
							<path d="M100 100 V36 Q100 20 84 20 H20" stroke-width="1"/>
							<path d="M90 90 C102 80 102 42 90 30" stroke-width="0.8"/>
							<path d="M80 80 C88 72 88 48 80 40" stroke-width="0.6"/>
							<circle cx="90" cy="90" r="2.2" fill="#b8934c"/>
							<circle cx="90" cy="30" r="2.2" fill="#b8934c"/>
						</g>
					</svg>
					<svg class="flourish bl" viewBox="0 0 120 120" aria-hidden="true">
						<g fill="none" stroke="#b8934c" stroke-linecap="round">
							<path d="M8 8 V96 Q8 112 24 112 H112" stroke-width="2"/>
							<path d="M20 20 V84 Q20 100 36 100 H100" stroke-width="1"/>
							<path d="M30 30 C18 40 18 78 30 90" stroke-width="0.8"/>
							<path d="M40 40 C32 48 32 72 40 80" stroke-width="0.6"/>
							<circle cx="30" cy="30" r="2.2" fill="#b8934c"/>
							<circle cx="30" cy="90" r="2.2" fill="#b8934c"/>
						</g>
					</svg>
					<svg class="flourish br" viewBox="0 0 120 120" aria-hidden="true">
						<g fill="none" stroke="#b8934c" stroke-linecap="round">
							<path d="M112 8 V96 Q112 112 96 112 H8" stroke-width="2"/>
							<path d="M100 20 V84 Q100 100 84 100 H20" stroke-width="1"/>
							<path d="M90 30 C102 40 102 78 90 90" stroke-width="0.8"/>
							<path d="M80 40 C88 48 88 72 80 80" stroke-width="0.6"/>
							<circle cx="90" cy="30" r="2.2" fill="#b8934c"/>
							<circle cx="90" cy="90" r="2.2" fill="#b8934c"/>
						</g>
					</svg>
				</div>

				<div class="cert-inner">
					<div class="logo-row">
						<div class="logo-side">
							<img src="/logos/karnataka-emblem.png" alt="Government of Karnataka" class="logo" />
						</div>
						<div class="header">
							<div class="kareis">ಕರ್ನಾಟಕ ವಸತಿ ಶಿಕ್ಷಣ ಸಂಸ್ಥೆಗಳ ಸೊಸೈಟಿ</div>
							<div class="kareis-en">KARNATAKA RESIDENTIAL EDUCATIONAL INSTITUTIONS SOCIETY</div>
							<div class="school">ಮೊರಾರ್ಜಿ ದೇಸಾಯಿ ವಸತಿ ಶಾಲೆ (SC-32) ಬಹದ್ದೂರ್ಘಟ್ಟ, ಚಿತ್ರದುರ್ಗ</div>
							<div class="school-en">MORARJI DESAI RESIDENTIAL SCHOOL (SC-32) BAHADDURGHATTA, CHITRADURGA</div>
						</div>
						<div class="logo-side">
							<img src="/logos/kreis-logo.png" alt="KREIS" class="logo" />
						</div>
					</div>

					<div class="title">
						<div class="title-rule"><span></span><i></i><span></span></div>
						<div class="title-kn">ಸಾಧನೆ ಪ್ರಮಾಣಪತ್ರ</div>
						<div class="title-en">CERTIFICATE OF ACHIEVEMENT</div>
						<div class="title-rule"><span></span><i></i><span></span></div>
					</div>

					<div class="certify">
						<span class="certify-en">This is to certify that</span>
						<span class="certify-kn">ಇದರಿಂದ ಪ್ರಮಾಣೀಕರಿಸಲಾಗುತ್ತದೆ</span>
					</div>

					<div class="student-name">{cert.student_name}</div>

					<div class="details">
						<span class="details-en">studying in</span>
						<span class="value class">{cert.class_name || '—'}</span>
						<span class="details-kn">ಇವರು ಭಾಗವಹಿಸಿದರು</span>
						<span class="details-en">has participated in</span>
						<span class="value event">{cert.event?.name || 'Event'}</span>
						<span class="category-kn">{categoryKn(cert.event?.category)}</span>
						<span class="category-en">{categoryEn(cert.event?.category)}</span>
						<span class="details-kn">ಮತ್ತು ಗಳಿಸಿದರು</span>
						<span class="details-en">and secured</span>
						<span class="value prize">{positionLabel(cert.position)}</span>
						<span class="prize-kn">({positionKn(cert.position)})</span>
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
	@import url('https://fonts.googleapis.com/css2?family=Cinzel:wght@500;600;700&family=Great+Vibes&family=Playfair+Display:ital,wght@0,400;0,500;0,600;1,400&family=Anek+Kannada:wght@400;500;600;700&family=Inter:wght@400;500;600&display=swap');

	:global(html), :global(body) {
		margin: 0;
		padding: 0;
		height: 100%;
		background: #dfe4ea;
		font-family: 'Inter', 'Anek Kannada', sans-serif;
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
			radial-gradient(ellipse at 50% 45%, rgba(255, 253, 246, 0) 0%, rgba(245, 237, 214, 0.55) 100%),
			#fdfaf1;
		overflow: hidden;
		box-shadow: 0 12px 48px rgba(0, 0, 0, 0.28);
	}

	.bg-pattern {
		position: absolute;
		top: 50%;
		left: 50%;
		width: 60%;
		height: 60%;
		transform: translate(-50%, -50%);
		opacity: 0.6;
		z-index: 0;
	}

	.frame {
		position: absolute;
		inset: 0;
		z-index: 1;
		pointer-events: none;
	}
	.frame-outer {
		position: absolute;
		inset: 6mm;
		border: 1.5px solid #b8934c;
		border-radius: 2mm;
	}
	.frame-inner {
		position: absolute;
		inset: 8.5mm;
		border: 0.8px solid #cdb379;
		border-radius: 1.5mm;
	}
	.flourish {
		position: absolute;
		width: 26mm;
		height: 26mm;
		z-index: 2;
	}
	.flourish.tl { top: 3.5mm; left: 3.5mm; }
	.flourish.tr { top: 3.5mm; right: 3.5mm; transform: scaleX(-1); }
	.flourish.bl { bottom: 3.5mm; left: 3.5mm; transform: scaleY(-1); }
	.flourish.br { bottom: 3.5mm; right: 3.5mm; transform: scale(-1, -1); }

	.cert-inner {
		position: relative;
		z-index: 3;
		padding: 12mm 22mm 10mm;
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
		width: 34mm;
		display: flex;
		justify-content: center;
	}
	.logo {
		width: 24mm;
		height: 24mm;
		object-fit: contain;
		filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.12));
	}

	.header {
		flex: 1;
		text-align: center;
	}
	.kareis {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 700;
		font-size: 14.5pt;
		color: #1e3a8a;
		letter-spacing: 0.5px;
	}
	.kareis-en {
		font-family: 'Cinzel', serif;
		font-weight: 600;
		font-size: 8.5pt;
		color: #b8934c;
		letter-spacing: 2.5px;
		margin-top: 2px;
		text-transform: uppercase;
	}
	.school {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 600;
		font-size: 11.5pt;
		color: #3b4252;
		margin-top: 5px;
	}
	.school-en {
		font-family: 'Playfair Display', serif;
		font-weight: 500;
		font-size: 8pt;
		color: #64748b;
		letter-spacing: 1px;
		margin-top: 2px;
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
		width: 75%;
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
		width: 7px;
		height: 7px;
		transform: rotate(45deg);
		background: #b8934c;
		flex-shrink: 0;
	}
	.title-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 700;
		font-size: 20pt;
		color: #1e3a8a;
		letter-spacing: 1px;
	}
	.title-en {
		font-family: 'Cinzel', serif;
		font-weight: 700;
		font-size: 14pt;
		color: #b8934c;
		letter-spacing: 6px;
		margin-top: 2px;
		text-transform: uppercase;
	}

	.certify {
		margin-top: 7mm;
		display: flex;
		flex-direction: column;
		gap: 1px;
	}
	.certify-en {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 10.5pt;
		color: #475569;
	}
	.certify-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 9pt;
		color: #64748b;
	}

	.student-name {
		margin-top: 4mm;
		font-family: 'Great Vibes', cursive;
		font-size: 27pt;
		color: #1e3a8a;
		padding: 0 24mm 3mm;
		border-bottom: 1.2px solid #b8934c;
		line-height: 1.5;
	}

	.details {
		margin-top: 6.5mm;
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		align-items: baseline;
		gap: 4px 7px;
		max-width: 235mm;
		line-height: 1.6;
	}
	.details-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 9.5pt;
		color: #64748b;
	}
	.details-en {
		font-family: 'Playfair Display', serif;
		font-size: 10pt;
		color: #475569;
	}
	.category-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 9.5pt;
		color: #64748b;
	}
	.category-en {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 9.5pt;
		color: #64748b;
	}
	.value {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		font-size: 11pt;
		color: #1e293b;
	}
	.value.class { color: #1e3a8a; }
	.value.event {
		color: #1e3a8a;
		font-size: 12pt;
	}
	.value.prize {
		color: #b45309;
		font-size: 12pt;
	}
	.prize-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 600;
		font-size: 9.5pt;
		color: #b45309;
	}

	.date-line {
		margin-top: 6mm;
		display: flex;
		align-items: baseline;
		gap: 6px;
		font-size: 9.5pt;
	}
	.date-en {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		color: #475569;
	}
	.date-value {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		color: #334155;
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
		background: #94a3b8;
	}
	.sig-name {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		font-size: 8.5pt;
		color: #1e293b;
		margin-top: 2.5mm;
	}
	.sig-role {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 7.5pt;
		color: #64748b;
		margin-top: 1mm;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	@media print {
		@page {
			size: A4 landscape;
			margin: 0;
		}
		:global(html), :global(body) {
			background: white;
			overflow: visible;
		}
		.stage, .stage.show {
			padding: 0;
			min-height: 0;
			display: block;
		}
		.toolbar, .no-print {
			display: none !important;
		}
		.cert-wrap {
			transform: none !important;
			height: 210mm !important;
		}
		.cert-page {
			margin: 0;
			width: 297mm;
			height: 210mm;
			box-shadow: none;
			-webkit-print-color-adjust: exact;
			print-color-adjust: exact;
		}
	}
</style>