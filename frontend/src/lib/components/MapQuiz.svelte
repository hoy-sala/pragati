<script lang="ts">
	export type MapPin = {
		key: string;
		label: string;
		svg_x: number;
		svg_y: number;
		correct?: boolean;
	};

	let {
		mapName = 'india',
		options = [],
		selectedKey = '',
		answered = false,
		revealed = false,
		hiddenKeys = [],
		disabled = false,
		onPick = (_key: string) => {}
	}: {
		mapName?: 'india' | 'karnataka' | string;
		options?: MapPin[];
		selectedKey?: string;
		answered?: boolean;
		revealed?: boolean;
		hiddenKeys?: string[];
		disabled?: boolean;
		onPick?: (key: string) => void;
	} = $props();

	let outlines = $state<{ id: string; d: string }[]>([]);
	let viewBox = $state('0 0 760 860');
	let baseVB = $state({ x: 0, y: 0, w: 760, h: 860 });
	let zoom = $state(1);
	let schematic = $state(false);
	let loadError = $state(false);

	$effect(() => {
		const m = mapName || 'india';
		loadError = false;
		zoom = 1;
		fetch(`/maps/${m}.json`)
			.then((r) => {
				if (!r.ok) throw new Error('map missing');
				return r.json();
			})
			.then((j) => {
				outlines = j.outlines ?? [];
				const vb = j.viewBox ?? [0, 0, 760, 860];
				baseVB = { x: vb[0], y: vb[1], w: vb[2], h: vb[3] };
				viewBox = `${vb[0]} ${vb[1]} ${vb[2]} ${vb[3]}`;
				schematic = !!j.schematic;
			})
			.catch(() => { loadError = true; outlines = []; });
	});

	function applyZoom() {
		const cx = baseVB.x + baseVB.w / 2;
		const cy = baseVB.y + baseVB.h / 2;
		const w = baseVB.w / zoom;
		const h = baseVB.h / zoom;
		viewBox = `${cx - w / 2} ${cy - h / 2} ${w} ${h}`;
	}
	function zoomIn() { zoom = Math.min(3, +(zoom + 0.4).toFixed(2)); applyZoom(); }
	function zoomOut() { zoom = Math.max(1, +(zoom - 0.4).toFixed(2)); applyZoom(); }
	function zoomReset() { zoom = 1; applyZoom(); }

	// Spread out pins that would otherwise sit on top of each other
	let placed = $derived.by(() => {
		const vis = options.filter((o) => !hiddenKeys.includes(o.key));
		const pts = vis.map((o) => ({ ...o, px: o.svg_x, py: o.svg_y }));
		const R = 30;
		for (let i = 0; i < pts.length; i++) {
			for (let j = i + 1; j < pts.length; j++) {
				const dx = pts[j].px - pts[i].px;
				const dy = pts[j].py - pts[i].py;
				const d = Math.hypot(dx, dy);
				if (d < R && d > 0.01) {
					const push = (R - d) / 2;
					const ux = dx / d, uy = dy / d;
					pts[i].px -= ux * push; pts[i].py -= uy * push;
					pts[j].px += ux * push; pts[j].py += uy * push;
				}
			}
		}
		return pts;
	});

	function pinState(o: MapPin): 'idle' | 'picked-right' | 'picked-wrong' | 'answer' {
		const isSel = answered && o.key === selectedKey;
		if (isSel) return o.correct ? 'picked-right' : 'picked-wrong';
		if ((answered || revealed) && o.correct) return 'answer';
		return 'idle';
	}

	function pick(key: string) {
		if (answered || disabled) return;
		onPick(key);
	}
</script>

<div class="map-wrap">
	<div class="map-frame">
		<svg viewBox={viewBox} class="map-svg" role="img" aria-label="Outline map quiz">
			<defs>
				<linearGradient id="mapOcean" x1="0" y1="0" x2="0" y2="1">
					<stop offset="0%" stop-color="#D8EAF7" />
					<stop offset="100%" stop-color="#F2F8FD" />
				</linearGradient>
				<pattern id="mapDots" width="18" height="18" patternUnits="userSpaceOnUse">
					<circle cx="2" cy="2" r="1.1" fill="#1F1A2E" opacity="0.06" />
				</pattern>
				<filter id="landShadow" x="-20%" y="-20%" width="140%" height="140%">
					<feDropShadow dx="4" dy="4" stdDeviation="0" flood-color="#1F1A2E" flood-opacity="0.9" />
				</filter>
				<filter id="pinGlow" x="-80%" y="-80%" width="260%" height="260%">
					<feDropShadow dx="0" dy="0" stdDeviation="5" flood-color="#0E7C71" flood-opacity="0.85" />
				</filter>
			</defs>

			<rect x="-2000" y="-2000" width="6000" height="6000" fill="url(#mapOcean)" />
			<rect x="-2000" y="-2000" width="6000" height="6000" fill="url(#mapDots)" />

			{#if loadError}
				<text x="50%" y="50%" text-anchor="middle" font-size="22" fill="#4A4458">Map outline unavailable</text>
			{:else}
				{#each outlines as o (o.id)}
					<path d={o.d} fill="#FFFBEC" stroke="#1F1A2E" stroke-width="3" stroke-linejoin="round" filter="url(#landShadow)" />
				{/each}
			{/if}

			<!-- compass -->
			<g transform="translate({baseVB.x + baseVB.w - 46},{baseVB.y + 30})" opacity="0.75">
				<circle r="16" fill="#FFFCF5" stroke="#1F1A2E" stroke-width="2" />
				<text y="-2" text-anchor="middle" font-size="13" font-weight="800" fill="#1F1A2E">N</text>
				<text y="10" text-anchor="middle" font-size="11" fill="#1F1A2E">↑</text>
			</g>

			{#each placed as p, i (p.key)}
				{@const st = pinState(p)}
				<g
					class="pin pin-{st}"
					style="animation-delay:{i * 90}ms"
					transform="translate({p.px},{p.py})"
					role="button"
					tabindex={answered || disabled ? -1 : 0}
					aria-label="Pin {p.key}{(answered || revealed) ? `: ${p.label}` : ''}"
					onclick={() => pick(p.key)}
					onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); pick(p.key); } }}
				>
					{#if st === 'picked-right' || st === 'answer'}
						<circle class="ring" r="17" fill="none" stroke="#0E7C71" stroke-width="2.5" />
					{/if}
					<circle class="dot" r="15" />
					<text y="6" text-anchor="middle" class="pin-letter">{p.key}</text>
					{#if answered || revealed}
						<text x="20" y="5" class="pin-label">{p.label}</text>
					{/if}
				</g>
			{/each}
		</svg>

		<div class="map-tools">
			<button onclick={zoomIn} class="tool" aria-label="Zoom in">+</button>
			<button onclick={zoomOut} class="tool" aria-label="Zoom out">−</button>
			<button onclick={zoomReset} class="tool" aria-label="Reset zoom">⤾</button>
		</div>
	</div>

	{#if schematic}
		<p class="hint schematic-note">Schematic outline for practice — not to survey scale.</p>
	{/if}
	<p class="hint pin-hint">Tap a pin on the map{#if !(answered || revealed)} — names appear after you answer{/if}.</p>
</div>

<style>
	.map-wrap { display: flex; flex-direction: column; gap: 0.6rem; }
	.map-frame {
		position: relative;
		background: var(--paper);
		border: 2.5px solid var(--ink);
		box-shadow: 6px 6px 0 var(--ink);
		border-radius: 18px;
		overflow: hidden;
	}
	.map-svg { display: block; width: 100%; height: auto; max-height: 460px; }
	.map-tools {
		position: absolute; right: 0.6rem; bottom: 0.6rem;
		display: flex; gap: 0.4rem;
	}
	.tool {
		width: 36px; height: 36px; border-radius: 10px;
		border: 2px solid var(--ink); background: var(--paper);
		font-weight: 800; font-size: 1.05rem; cursor: pointer; color: var(--ink);
		box-shadow: 2px 2px 0 var(--ink);
	}
	.tool:active { transform: translate(2px, 2px); box-shadow: 0 0 0 var(--ink); }
	.schematic-note { margin: 0; }

	.pin { cursor: pointer; animation: dropIn 0.35s ease-out backwards; }
	.pin .dot {
		fill: #FFFCF5; stroke: #1F1A2E; stroke-width: 2.5;
		transition: fill 120ms, transform 120ms;
		transform-box: fill-box; transform-origin: center;
	}
	.pin:hover .dot { fill: #FDE9C2; transform: scale(1.12); }
	.pin:focus-visible { outline: 3px solid var(--plum); outline-offset: 3px; }
	.pin-letter {
		font-family: var(--font-display); font-weight: 800; font-size: 15px;
		fill: #1F1A2E; pointer-events: none;
	}
	.pin-picked-right .dot, .pin-answer .dot { fill: #D8F3E3; stroke: #0E7C71; }
	.pin-picked-right .pin-letter, .pin-answer .pin-letter { fill: #0E7C71; }
	.pin-picked-right .dot { filter: url(#pinGlow); }
	.pin-picked-wrong .dot { fill: #FBDAD3; stroke: #C2381B; }
	.pin-picked-wrong .pin-letter { fill: #C2381B; }
	.pin-picked-wrong { animation: shakeIt 0.45s ease-out; }

	.ring {
		transform-box: fill-box; transform-origin: center;
		animation: ringPulse 1.4s ease-out infinite;
	}
	@keyframes ringPulse {
		0% { transform: scale(0.7); opacity: 1; }
		100% { transform: scale(1.5); opacity: 0; }
	}
	@keyframes dropIn {
		0% { opacity: 0; transform: translateY(-10px) scale(0.7); }
		100% { opacity: 1; transform: translateY(0) scale(1); }
	}
	@keyframes shakeIt { 0%,100% { transform: translateX(0); } 20% { transform: translateX(-5px); } 40% { transform: translateX(5px); } 60% { transform: translateX(-3px); } 80% { transform: translateX(3px); } }

	.pin-label {
		font-family: var(--font-body); font-weight: 700; font-size: 13px;
		fill: #1F1A2E; stroke: #FFFCF5; stroke-width: 4px; paint-order: stroke;
	}
	.pin-hint { margin: 0; font-size: 0.82rem; text-align: center; }
</style>
