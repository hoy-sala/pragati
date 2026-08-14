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

				<!-- rangoli / lotus background mandala -->
				<svg class="bg-mandala" viewBox="0 0 200 200" aria-hidden="true">
					<g transform="translate(100,100)" fill="none" stroke="#c9a227">
						<circle r="92" stroke-width="0.35" opacity="0.5"/>
						<circle r="84" stroke-width="0.2" opacity="0.4"/>
						<circle r="70" stroke-width="0.3" opacity="0.45"/>
						<circle r="62" stroke-width="0.15" opacity="0.35"/>
						<!-- radiating spokes -->
						{#each Array(72) as _, i}
							<line x1="0" y1="-62" x2="0" y2="-92" transform="rotate({i * 5})" stroke-width="0.12" opacity="0.35"/>
						{/each}
						<!-- outer lotus petals -->
						<g opacity="0.5">
							{#each Array(16) as _, i}
								<path d="M0 -92 Q12 -84 0 -70 Q-12 -84 0 -92" transform="rotate({i * 22.5})" stroke-width="0.35"/>
							{/each}
						</g>
						<!-- middle ring of dots (bindi) -->
						<g opacity="0.55">
							{#each Array(24) as _, i}
								<circle cx="0" cy="-46" r="1.1" transform="rotate({i * 15})" fill="#c9a227" stroke="none"/>
							{/each}
						</g>
						<!-- inner lotus -->
						<g opacity="0.55">
							{#each Array(12) as _, i}
								<path d="M0 -46 Q9 -40 0 -30 Q-9 -40 0 -46" transform="rotate({i * 30})" stroke-width="0.3"/>
							{/each}
						</g>
						<circle r="22" stroke-width="0.25" opacity="0.5"/>
						<g opacity="0.6">
							{#each Array(8) as _, i}
								<path d="M0 -22 Q6 -18 0 -12 Q-6 -18 0 -22" transform="rotate({i * 45})" stroke-width="0.3"/>
							{/each}
						</g>
						<circle r="6" stroke-width="0.3" opacity="0.6"/>
					</g>
				</svg>

				<!-- faint paisley corner sprays -->
				<svg class="bg-paisley tl" viewBox="0 0 120 120" aria-hidden="true">
					<g fill="none" stroke="#d8b45a" stroke-linecap="round">
						<path d="M20 100 C6 88 4 56 20 40 C32 52 32 84 20 100" stroke-width="0.5" opacity="0.5"/>
						<path d="M20 78 C14 70 14 52 20 46 C26 52 26 70 20 78" stroke-width="0.35" opacity="0.4"/>
						<path d="M36 60 C30 54 30 42 36 38 C42 42 42 54 36 60" stroke-width="0.3" opacity="0.35"/>
					</g>
				</svg>
				<svg class="bg-paisley tr" viewBox="0 0 120 120" aria-hidden="true">
					<g fill="none" stroke="#d8b45a" stroke-linecap="round">
						<path d="M100 100 C114 88 116 56 100 40 C88 52 88 84 100 100" stroke-width="0.5" opacity="0.5"/>
						<path d="M100 78 C106 70 106 52 100 46 C94 52 94 70 100 78" stroke-width="0.35" opacity="0.4"/>
						<path d="M84 60 C90 54 90 42 84 38 C78 42 78 54 84 60" stroke-width="0.3" opacity="0.35"/>
					</g>
				</svg>
				<svg class="bg-paisley bl" viewBox="0 0 120 120" aria-hidden="true">
					<g fill="none" stroke="#d8b45a" stroke-linecap="round">
						<path d="M20 20 C6 32 4 64 20 80 C32 68 32 36 20 20" stroke-width="0.5" opacity="0.5"/>
						<path d="M20 42 C14 50 14 68 20 74 C26 68 26 50 20 42" stroke-width="0.35" opacity="0.4"/>
					</g>
				</svg>
				<svg class="bg-paisley br" viewBox="0 0 120 120" aria-hidden="true">
					<g fill="none" stroke="#d8b45a" stroke-linecap="round">
						<path d="M100 20 C114 32 116 64 100 80 C88 68 88 36 100 20" stroke-width="0.5" opacity="0.5"/>
						<path d="M100 42 C106 50 106 68 100 74 C94 68 94 50 100 42" stroke-width="0.35" opacity="0.4"/>
					</g>
				</svg>

				<!-- ornamental woven border -->
				<div class="frame">
					<div class="frame-line outer"></div>
					<div class="frame-band zari"></div>
					<div class="frame-line inner"></div>
					<div class="frame-corner"></div>
					<svg class="corner tl" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#c9a227" stroke-linecap="round">
							<path d="M16 124 V28 Q16 16 28 16 H124" stroke-width="2"/>
							<path d="M28 112 V40 Q28 28 40 28 H112" stroke-width="0.8"/>
							<path d="M46 104 C34 94 32 66 46 52 C56 62 56 92 46 104" stroke-width="1"/>
							<path d="M46 86 C40 80 40 64 46 58 C52 64 52 80 46 86" stroke-width="0.6"/>
							<circle cx="46" cy="102" r="2.4" fill="#c9a227" stroke="none"/>
							<path d="M30 26 L34 22 L38 26 L34 30 Z" fill="#c9a227" stroke="none"/>
							<circle cx="108" cy="30" r="1.8" fill="#c9a227" stroke="none"/>
						</g>
					</svg>
					<svg class="corner tr" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#c9a227" stroke-linecap="round">
							<path d="M124 124 V28 Q124 16 112 16 H16" stroke-width="2"/>
							<path d="M112 112 V40 Q112 28 100 28 H28" stroke-width="0.8"/>
							<path d="M94 104 C106 94 108 66 94 52 C84 62 84 92 94 104" stroke-width="1"/>
							<path d="M94 86 C100 80 100 64 94 58 C88 64 88 80 94 86" stroke-width="0.6"/>
							<circle cx="94" cy="102" r="2.4" fill="#c9a227" stroke="none"/>
							<path d="M110 26 L114 22 L118 26 L114 30 Z" fill="#c9a227" stroke="none"/>
							<circle cx="32" cy="30" r="1.8" fill="#c9a227" stroke="none"/>
						</g>
					</svg>
					<svg class="corner bl" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#c9a227" stroke-linecap="round">
							<path d="M16 16 V112 Q16 124 28 124 H124" stroke-width="2"/>
							<path d="M28 28 V100 Q28 112 40 112 H112" stroke-width="0.8"/>
							<path d="M46 36 C34 46 32 74 46 88 C56 78 56 48 46 36" stroke-width="1"/>
							<path d="M46 54 C40 60 40 76 46 82 C52 76 52 60 46 54" stroke-width="0.6"/>
							<circle cx="46" cy="38" r="2.4" fill="#c9a227" stroke="none"/>
							<path d="M30 114 L34 110 L38 114 L34 118 Z" fill="#c9a227" stroke="none"/>
							<circle cx="108" cy="110" r="1.8" fill="#c9a227" stroke="none"/>
						</g>
					</svg>
					<svg class="corner br" viewBox="0 0 140 140" aria-hidden="true">
						<g fill="none" stroke="#c9a227" stroke-linecap="round">
							<path d="M124 16 V112 Q124 124 112 124 H16" stroke-width="2"/>
							<path d="M112 28 V100 Q112 112 100 112 H28" stroke-width="0.8"/>
							<path d="M94 36 C106 46 108 74 94 88 C84 78 84 48 94 36" stroke-width="1"/>
							<path d="M94 54 C100 60 100 76 94 82 C88 76 88 60 94 54" stroke-width="0.6"/>
							<circle cx="94" cy="38" r="2.4" fill="#c9a227" stroke="none"/>
							<path d="M110 114 L114 110 L118 114 L114 118 Z" fill="#c9a227" stroke="none"/>
							<circle cx="32" cy="110" r="1.8" fill="#c9a227" stroke="none"/>
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
			radial-gradient(ellipse at 50% 42%, rgba(255, 251, 240, 0) 0%, rgba(248, 235, 205, 0.7) 100%),
			#fdf6e7;
		overflow: hidden;
		box-shadow: 0 12px 48px rgba(0, 0, 0, 0.28);
	}

	/* rangoli / lotus mandala background */
	.bg-mandala {
		position: absolute;
		top: 50%;
		left: 50%;
		width: 76%;
		height: 76%;
		transform: translate(-50%, -50%);
		z-index: 0;
	}

	.bg-paisley {
		position: absolute;
		width: 42mm;
		height: 42mm;
		z-index: 0;
		opacity: 0.7;
	}
	.bg-paisley.tl { top: 20mm; left: 22mm; }
	.bg-paisley.tr { top: 20mm; right: 22mm; }
	.bg-paisley.bl { bottom: 16mm; left: 22mm; }
	.bg-paisley.br { bottom: 16mm; right: 22mm; }

	/* ornamental woven border */
	.frame {
		position: absolute;
		inset: 0;
		z-index: 1;
		pointer-events: none;
	}
	.frame-line.outer {
		position: absolute;
		inset: 3.2mm;
		border: 1.4px solid #a8842c;
		border-radius: 1.6mm;
	}
	.frame-band.zari {
		position: absolute;
		inset: 5.2mm;
		border-radius: 1.4mm;
		overflow: hidden;
		background:
			repeating-linear-gradient(45deg, #7b1f2d 0 0.9mm, #c9a227 0.9mm 1.8mm);
		mask:
			linear-gradient(#fff 0 0) content-box,
			linear-gradient(#fff 0 0);
		mask-composite: exclude;
		-webkit-mask:
			linear-gradient(#fff 0 0) content-box,
			linear-gradient(#fff 0 0);
		-webkit-mask-composite: xor;
		padding: 1.1mm;
		opacity: 0.9;
	}
	.frame-line.inner {
		position: absolute;
		inset: 8.2mm;
		border: 1px solid #a8842c;
		border-radius: 1.2mm;
	}
	.frame-corner {
		position: absolute;
		inset: 3.2mm;
		border-radius: 1.6mm;
		box-shadow:
			0 0 0 0.4mm #fdf6e7,
			0 0 0 0.75mm #7b1f2d inset;
		pointer-events: none;
	}

	.corner {
		position: absolute;
		width: 34mm;
		height: 34mm;
		z-index: 2;
		filter: drop-shadow(0 0.5mm 0.4mm rgba(0,0,0,0.15));
	}
	.corner.tl { top: 0.8mm; left: 0.8mm; }
	.corner.tr { top: 0.8mm; right: 0.8mm; }
	.corner.bl { bottom: 0.8mm; left: 0.8mm; }
	.corner.br { bottom: 0.8mm; right: 0.8mm; }

	.cert-inner {
		position: relative;
		z-index: 3;
		padding: 14mm 24mm 11mm;
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
		width: 36mm;
		display: flex;
		justify-content: center;
	}
	.logo {
		width: 24mm;
		height: 24mm;
		object-fit: contain;
		filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.14));
	}

	.header {
		flex: 1;
		text-align: center;
	}
	.kareis {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 700;
		font-size: 14.5pt;
		color: #7b1f2d;
		letter-spacing: 0.5px;
	}
	.kareis-en {
		font-family: 'Cinzel', serif;
		font-weight: 600;
		font-size: 8.5pt;
		color: #a8842c;
		letter-spacing: 2.5px;
		margin-top: 2px;
		text-transform: uppercase;
	}
	.school {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 600;
		font-size: 11.5pt;
		color: #4a3f2f;
		margin-top: 5px;
	}
	.school-en {
		font-family: 'Playfair Display', serif;
		font-weight: 500;
		font-size: 8pt;
		color: #7a6a4e;
		letter-spacing: 1px;
		margin-top: 2px;
		text-transform: uppercase;
	}

	.title {
		margin-top: 8.5mm;
		width: 100%;
	}
	.title-rule {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		margin: 2.6mm auto;
		width: 72%;
	}
	.title-rule span {
		flex: 1;
		height: 1px;
		background: linear-gradient(90deg, transparent, #a8842c);
	}
	.title-rule span:last-child {
		background: linear-gradient(90deg, #a8842c, transparent);
	}
	.title-rule i {
		width: 7px;
		height: 7px;
		transform: rotate(45deg);
		background: #c9a227;
		flex-shrink: 0;
	}
	.title-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 700;
		font-size: 20pt;
		color: #7b1f2d;
		letter-spacing: 1px;
	}
	.title-en {
		font-family: 'Cinzel', serif;
		font-weight: 700;
		font-size: 14pt;
		color: #a8842c;
		letter-spacing: 6px;
		margin-top: 2px;
		text-transform: uppercase;
	}

	.certify {
		margin-top: 6.5mm;
		display: flex;
		flex-direction: column;
		gap: 1px;
	}
	.certify-en {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 10.5pt;
		color: #4a3f2f;
	}
	.certify-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 9pt;
		color: #7a6a4e;
	}

	.student-name {
		margin-top: 4mm;
		font-family: 'Great Vibes', cursive;
		font-size: 27pt;
		color: #7b1f2d;
		padding: 0 24mm 3mm;
		border-bottom: 1.2px solid #c9a227;
		line-height: 1.5;
	}

	.details {
		margin-top: 6mm;
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		align-items: baseline;
		gap: 4px 7px;
		max-width: 238mm;
		line-height: 1.6;
	}
	.details-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 9.5pt;
		color: #7a6a4e;
	}
	.details-en {
		font-family: 'Playfair Display', serif;
		font-size: 10pt;
		color: #4a3f2f;
	}
	.category-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 9.5pt;
		color: #7a6a4e;
	}
	.category-en {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		font-size: 9.5pt;
		color: #7a6a4e;
	}
	.value {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		font-size: 11pt;
		color: #2d2a24;
	}
	.value.class { color: #7b1f2d; }
	.value.event {
		color: #7b1f2d;
		font-size: 12pt;
	}
	.value.prize {
		color: #d97706;
		font-size: 12pt;
	}
	.prize-kn {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 600;
		font-size: 9.5pt;
		color: #b45309;
	}

	.date-line {
		margin-top: 5.5mm;
		display: flex;
		align-items: baseline;
		gap: 6px;
		font-size: 9.5pt;
	}
	.date-en {
		font-family: 'Playfair Display', serif;
		font-weight: 600;
		color: #4a3f2f;
	}
	.date-value {
		font-family: 'Playfair Display', serif;
		font-style: italic;
		color: #2d2a24;
	}
	.date-sep { color: #a8842c; }

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
		color: #2d2a24;
		margin-top: 2.5mm;
	}
	.sig-role {
		font-family: 'Anek Kannada', sans-serif;
		font-weight: 500;
		font-size: 7.5pt;
		color: #7a6a4e;
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