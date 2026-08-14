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

	const PRIZE_KN: Record<string, string> = {
		'First Prize': 'ಪ್ರಥಮ ಬಹುಮಾನ',
		'Runner Up': 'ದ್ವಿತೀಯ ಬಹುಮಾನ',
		'Consolation': 'ತೃತೀಯ ಬಹುಮಾನ',
		'Participation': 'ಭಾಗವಹಿಸುವಿಕೆ',
	};

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
		<!-- background guilloché watermark -->
		<svg class="watermark" viewBox="0 0 100 100" aria-hidden="true">
			<g transform="translate(50,50)">
				<circle r="44" fill="none" stroke="#e2e8f0" stroke-width="0.3"/>
				<circle r="40" fill="none" stroke="#e2e8f0" stroke-width="0.2"/>
				<circle r="36" fill="none" stroke="#e2e8f0" stroke-width="0.25"/>
				<circle r="32" fill="none" stroke="#e2e8f0" stroke-width="0.2"/>
				<circle r="28" fill="none" stroke="#e2e8f0" stroke-width="0.3"/>
				<g stroke="#e2e8f0" stroke-width="0.15">
					{#each Array(36) as _, i}
						<line x1="0" y1="-28" x2="0" y2="-44" transform="rotate({i * 10})"/>
					{/each}
				</g>
			</g>
		</svg>

		<!-- guilloché corner rosettes -->
		<svg class="corner tl" viewBox="0 0 120 120" aria-hidden="true">
			<g transform="translate(60,60)">
				<circle r="52" fill="none" stroke="#c7a25a" stroke-width="0.6"/>
				<circle r="46" fill="none" stroke="#c7a25a" stroke-width="0.35"/>
				<circle r="40" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				<circle r="34" fill="none" stroke="#c7a25a" stroke-width="0.3"/>
				<circle r="28" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				{#each Array(48) as _, i}
					<path d="M 0 -28 L 0 -52" transform="rotate({i * 7.5})" stroke="#c7a25a" stroke-width="0.25"/>
				{/each}
			</g>
		</svg>
		<svg class="corner tr" viewBox="0 0 120 120" aria-hidden="true">
			<g transform="translate(60,60)">
				<circle r="52" fill="none" stroke="#c7a25a" stroke-width="0.6"/>
				<circle r="46" fill="none" stroke="#c7a25a" stroke-width="0.35"/>
				<circle r="40" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				<circle r="34" fill="none" stroke="#c7a25a" stroke-width="0.3"/>
				<circle r="28" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				{#each Array(48) as _, i}
					<path d="M 0 -28 L 0 -52" transform="rotate({i * 7.5})" stroke="#c7a25a" stroke-width="0.25"/>
				{/each}
			</g>
		</svg>
		<svg class="corner bl" viewBox="0 0 120 120" aria-hidden="true">
			<g transform="translate(60,60)">
				<circle r="52" fill="none" stroke="#c7a25a" stroke-width="0.6"/>
				<circle r="46" fill="none" stroke="#c7a25a" stroke-width="0.35"/>
				<circle r="40" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				<circle r="34" fill="none" stroke="#c7a25a" stroke-width="0.3"/>
				<circle r="28" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				{#each Array(48) as _, i}
					<path d="M 0 -28 L 0 -52" transform="rotate({i * 7.5})" stroke="#c7a25a" stroke-width="0.25"/>
				{/each}
			</g>
		</svg>
		<svg class="corner br" viewBox="0 0 120 120" aria-hidden="true">
			<g transform="translate(60,60)">
				<circle r="52" fill="none" stroke="#c7a25a" stroke-width="0.6"/>
				<circle r="46" fill="none" stroke="#c7a25a" stroke-width="0.35"/>
				<circle r="40" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				<circle r="34" fill="none" stroke="#c7a25a" stroke-width="0.3"/>
				<circle r="28" fill="none" stroke="#c7a25a" stroke-width="0.5"/>
				{#each Array(48) as _, i}
					<path d="M 0 -28 L 0 -52" transform="rotate({i * 7.5})" stroke="#c7a25a" stroke-width="0.25"/>
				{/each}
			</g>
		</svg>

		<div class="cert-inner">
			<div class="logo-row">
				<img src="/logos/karnataka-emblem.png" alt="Government of Karnataka" class="logo" />
				<img src="/logos/kreis-logo.png" alt="KREIS" class="logo" />
			</div>

			<div class="header">
				<div class="kareis">ಕರ್ನಾಟಕ ವಸತಿ ಶಿಕ್ಷಣ ಸಂಸ್ಥೆಗಳ ಸೊಸೈಟಿ</div>
				<div class="kareis-en">KARNATAKA RESIDENTIAL EDUCATIONAL INSTITUTIONS SOCIETY</div>
				<div class="school">ಮೊರಾರ್ಜಿ ದೇಸಾಯಿ ವಸತಿ ಶಾಲೆ (SC-32) ಬಹದ್ದೂರ್ಘಟ್ಟ, ಚಿತ್ರದುರ್ಗ</div>
				<div class="school-en">MORARJI DESAI RESIDENTIAL SCHOOL (SC-32) BAHADDURGHATTA, CHITRADURGA</div>
			</div>

			<div class="title">
				<div class="title-kn">ಸಾಧನೆ ಪ್ರಮಾಣಪತ್ರ</div>
				<div class="title-en">CERTIFICATE OF ACHIEVEMENT</div>
			</div>

			<div class="certify">
				<span class="certify-kn">ಇದರಿಂದ ಪ್ರಮಾಣೀಕರಿಸಲಾಗುತ್ತದೆ,</span>
				<span class="certify-en">This is to certify that</span>
			</div>

			<div class="student-name">{cert.student_name}</div>

			<div class="details">
				<span class="details-kn">ಇವರು</span>
				<span class="details-en">studying in</span>
				<span class="value">{cert.class_name || '—'}</span>
				<span class="details-kn">ಇವರು ಭಾಗವಹಿಸಿದರು</span>
				<span class="details-en">has participated in</span>
				<span class="value event">{cert.event?.name || 'Event'}</span>
				{#if cert.event?.category === 'sports'}
					<span class="details-kn">ಕ್ರೀಡಾ ಸ್ಪರ್ಧೆಯಲ್ಲಿ</span>
					<span class="details-en">Sports Competition</span>
				{:else if cert.event?.category === 'cultural'}
					<span class="details-kn">ಸಾಂಸ್ಕೃತಿಕ ಸ್ಪರ್ಧೆಯಲ್ಲಿ</span>
					<span class="details-en">Cultural Competition</span>
				{:else if cert.event?.category === 'academic'}
					<span class="details-kn">ಶೈಕ್ಷಣಿಕ ಸ್ಪರ್ಧೆಯಲ್ಲಿ</span>
					<span class="details-en">Academic Competition</span>
				{/if}
				<span class="details-kn">ಮತ್ತು ಗಳಿಸಿದರು</span>
				<span class="details-en">and secured</span>
				<span class="value prize">{positionLabel(cert.position)}</span>
				<span class="prize-kn">({positionKn(cert.position)})</span>
			</div>

			<div class="date-line">
				<span class="details-kn">ದಿನಾಂಕ</span>
				<span class="details-en">Date</span>
				<span class="value">{formatDate(cert.issue_date || cert.event?.held_date)}</span>
				{#if cert.event?.venue}
					<span class="venue"> • {cert.event.venue}</span>
				{/if}
			</div>

			<div class="signatures">
				{#each cert.signatories as sig, i}
					<div class="signature-block">
						{#if sig.signature_url}
							<img src={apiUrl(sig.signature_url)} alt="signature" class="signature-img" />
						{:else}
							<div class="signature-line"></div>
						{/if}
						<div class="sig-name">{sig.name}</div>
						<div class="sig-role">{sig.title || ROLE_LABELS[sig.role] || sig.role}</div>
					</div>
				{:else}
					<div class="signature-block">
						<div class="signature-line"></div>
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
	@import url('https://fonts.googleapis.com/css2?family=Cinzel:wght@600;700;800&family=Great+Vibes&family=Playfair+Display:ital,wght@0,500;0,600;0,700;1,500&family=Anek+Kannada:wght@400;500;600;700&family=Inter:wght@400;500;600;700&display=swap');

	:global(html), :global(body) {
		margin: 0;
		padding: 0;
		height: 100%;
		background: #e2e8f0;
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
		background: #fefcf6;
		overflow: hidden;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.25);
	}

	.watermark {
		position: absolute;
		top: 50%;
		left: 50%;
		width: 55%;
		height: 55%;
		transform: translate(-50%, -50%);
		opacity: 0.55;
		z-index: 0;
	}

	.corner {
		position: absolute;
		width: 28mm;
		height: 28mm;
		opacity: 0.9;
		z-index: 1;
	}
	.corner.tl { top: 6mm; left: 6mm; }
	.corner.tr { top: 6mm; right: 6mm; transform: rotate(90deg); }
	.corner.bl { bottom: 6mm; left: 6mm; transform: rotate(-90deg); }
	.corner.br { bottom: 6mm; right: 6mm; transform: rotate(180deg); }

	.cert-inner {
		position: relative;
		z-index: 2;
		padding: 10mm 18mm;
		border: 1.5px solid #c7a25a;
		margin: 3.5mm;
		height: calc(210mm - 7mm);
		box-sizing: border-box;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
	}

	.logo-row {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		width: 100%;
		position: absolute;
		top: 8mm;
		left: 0;
		right: 0;
		padding: 0 8mm;
	}
	.logo {
		width: 20mm;
		height: 20mm;
		object-fit: contain;
	}

	.header {
		margin-top: 24mm;
	}
	.kareis {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 700;
		font-size: 15pt;
		color: #1e3a8a;
		letter-spacing: 0.5px;
	}
	.kareis-en {
		font-family: 'Cinzel', serif;
		font-weight: 700;
		font-size: 9.5pt;
		color: #b45309;
		letter-spacing: 2px;
		margin-top: 2px;
	}
	.school {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 600;
		font-size: 12pt;
		color: #1e293b;
		margin-top: 6px;
	}
	.school-en {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		font-size: 8.5pt;
		color: #475569;
		letter-spacing: 0.5px;
		margin-top: 2px;
	}

	.title {
		margin-top: 10mm;
	}
	.title-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 700;
		font-size: 22pt;
		color: #b45309;
	}
	.title-en {
		font-family: 'Cinzel', serif;
		font-weight: 800;
		font-size: 17pt;
		color: #1e3a8a;
		letter-spacing: 4px;
		margin-top: 3px;
		text-transform: uppercase;
	}

	.certify {
		margin-top: 8mm;
		display: flex;
		flex-direction: column;
		gap: 2px;
	}
	.certify-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 11pt;
		color: #475569;
	}
	.certify-en {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 10.5pt;
		color: #475569;
	}

	.student-name {
		margin-top: 5mm;
		font-family: 'Great Vibes', cursive;
		font-size: 30pt;
		color: #1e3a8a;
		padding: 0 20mm;
		border-bottom: 1px solid #c7a25a;
		line-height: 1.6;
	}

	.details {
		margin-top: 7mm;
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		align-items: baseline;
		gap: 4px 8px;
		max-width: 240mm;
	}
	.details-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 10pt;
		color: #475569;
	}
	.details-en {
		font-family: 'Inter', sans-serif;
		font-size: 9pt;
		color: #475569;
	}
	.value {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 11.5pt;
		color: #1e293b;
	}
	.value.event {
		color: #1e3a8a;
		font-size: 12.5pt;
	}
	.value.prize {
		color: #b45309;
		font-size: 13pt;
	}
	.prize-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 600;
		font-size: 10pt;
		color: #b45309;
	}

	.date-line {
		margin-top: 5mm;
		display: flex;
		align-items: baseline;
		gap: 6px;
		font-size: 10pt;
		color: #475569;
	}
	.date-line .venue { color: #64748b; font-style: italic; }

	.signatures {
		margin-top: auto;
		display: flex;
		justify-content: space-around;
		align-items: flex-end;
		width: 100%;
		padding-top: 4mm;
		gap: 8mm;
	}
	.signature-block {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 55mm;
	}
	.signature-img {
		height: 12mm;
		width: auto;
		object-fit: contain;
		max-width: 50mm;
		margin-bottom: 1mm;
	}
	.signature-line {
		width: 100%;
		height: 12mm;
		border-bottom: 1px solid #94a3b8;
		margin-bottom: 1mm;
	}
	.sig-name {
		font-family: 'Playfair Display', serif;
		font-weight: 700;
		font-size: 9pt;
		color: #1e293b;
	}
	.sig-role {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 8pt;
		color: #64748b;
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