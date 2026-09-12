<script lang="ts">
	import type { BirdShapeKind } from '$lib/data/birds';

	let {
		shape,
		tail,
		crest,
		class: className = ''
	}: {
		shape: BirdShapeKind;
		tail?: 'short' | 'medium' | 'long' | 'fork' | 'cock' | 'streamers';
		crest?: 'none' | 'tuft' | 'fan';
		class?: string;
	} = $props();

	type PerchCfg = {
		tail: 'short' | 'medium' | 'long' | 'fork' | 'cock' | 'streamers';
		bill: 'cone' | 'heavy' | 'longthin' | 'curve' | 'hook' | 'dagger';
		crest: 'none' | 'tuft' | 'fan';
		bulk: number;
		headR: number;
		legs: 'short' | 'long';
	};

	const DEFAULTS: Record<string, PerchCfg> = {
		songbird:   { tail: 'short',  bill: 'cone',     crest: 'none', bulk: 0,  headR: 10, legs: 'short' },
		crow:       { tail: 'medium', bill: 'heavy',    crest: 'none', bulk: 2,  headR: 10, legs: 'short' },
		parrot:     { tail: 'long',   bill: 'hook',     crest: 'none', bulk: 1,  headR: 10, legs: 'short' },
		dove:       { tail: 'medium', bill: 'cone',     crest: 'none', bulk: 2,  headR: 9,  legs: 'short' },
		longtail:   { tail: 'long',   bill: 'cone',     crest: 'none', bulk: 0,  headR: 10, legs: 'short' },
		hoopoe:     { tail: 'medium', bill: 'longthin', crest: 'fan',  bulk: 0,  headR: 9,  legs: 'short' },
		barbet:     { tail: 'short',  bill: 'heavy',    crest: 'none', bulk: 3,  headR: 11, legs: 'short' },
		drongo:     { tail: 'fork',   bill: 'cone',     crest: 'none', bulk: 0,  headR: 10, legs: 'short' },
		sunbird:    { tail: 'medium', bill: 'curve',    crest: 'none', bulk: -1, headR: 9,  legs: 'short' },
		kingfisher: { tail: 'short',  bill: 'dagger',   crest: 'none', bulk: 1,  headR: 12, legs: 'short' },
		slim:       { tail: 'streamers', bill: 'longthin', crest: 'none', bulk: -1, headR: 9, legs: 'short' },
		wader:      { tail: 'short',  bill: 'cone',     crest: 'none', bulk: 0,  headR: 9,  legs: 'long' },
	};

	let pose: 'perch' | 'soar' | 'front' | 'train' = $derived(
		shape === 'raptor' ? 'soar' : shape === 'owl' ? 'front' : shape === 'peafowl' ? 'train' : 'perch'
	);

	let cfg: PerchCfg = $derived({
		...(DEFAULTS[shape] ?? DEFAULTS.songbird),
		...(tail ? { tail } : {}),
		...(crest ? { crest } : {}),
	});

	let rx = $derived(20 + cfg.bulk * 2.2);
	let ry = $derived(15 + cfg.bulk * 1.8);
</script>

<svg viewBox="0 0 100 100" class={className} fill="currentColor" aria-hidden="true">
	{#if pose === 'soar'}
		<!-- soaring raptor -->
		<path d="M6 58 Q30 32 48 50 Q52 54 56 50 Q74 32 94 58 Q74 50 56 60 L52 62 L48 60 Q26 50 6 58 Z" />
		<ellipse cx="52" cy="57" rx="6" ry="4.5" />
		<path d="M46 60 L40 68 L44 68 L48 62 Z" />
	{:else if pose === 'front'}
		<!-- owl, front view -->
		<path d="M38 34 L32 18 L44 28 Z" />
		<path d="M62 34 L68 18 L56 28 Z" />
		<ellipse cx="50" cy="60" rx="21" ry="26" />
		<path d="M44 84 L44 92 M56 84 L56 92" stroke="currentColor" stroke-width="3" stroke-linecap="round" />
		<circle cx="42" cy="52" r="5.5" fill="white" fill-opacity="0.92" />
		<circle cx="58" cy="52" r="5.5" fill="white" fill-opacity="0.92" />
		<circle cx="42" cy="52" r="2.4" />
		<circle cx="58" cy="52" r="2.4" />
		<path d="M47 62 L53 62 L50 67 Z" />
	{:else if pose === 'train'}
		<!-- peacock -->
		<path d="M30 78 L14 96 M38 80 L28 98 M62 80 L72 98 M70 78 L86 96" stroke="currentColor" stroke-width="3" stroke-linecap="round" />
		<circle cx="14" cy="95" r="2.6" />
		<circle cx="28" cy="97" r="2.6" />
		<circle cx="72" cy="97" r="2.6" />
		<circle cx="86" cy="95" r="2.6" />
		<ellipse cx="52" cy="62" rx="19" ry="15" />
		<circle cx="70" cy="42" r="10" />
		<path d="M79 37 L92 41 L79 44 Z" />
		<path d="M66 32 L62 20 M70 31 L70 18 M74 32 L78 20" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" />
		<circle cx="62" cy="18" r="2.2" />
		<circle cx="70" cy="16" r="2.2" />
		<circle cx="78" cy="18" r="2.2" />
		<path d="M46 76 L46 88 M56 76 L56 88" stroke="currentColor" stroke-width="3" stroke-linecap="round" />
	{:else}
		<!-- perched bird, side profile facing right -->
		{#if cfg.tail === 'short'}
			<path d="M38 58 L14 66 L16 73 L40 67 Z" />
		{:else if cfg.tail === 'medium'}
			<path d="M38 56 L10 68 L13 75 L40 68 Z" />
		{:else if cfg.tail === 'long'}
			<path d="M38 55 L6 74 L10 80 L40 68 Z" />
		{:else if cfg.tail === 'fork'}
			<path d="M38 55 L8 63 L10 69 L40 63 Z" />
			<path d="M38 61 L10 75 L14 80 L40 68 Z" />
		{:else if cfg.tail === 'cock'}
			<path d="M40 58 L46 24 L53 26 L46 60 Z" />
		{:else if cfg.tail === 'streamers'}
			<path d="M38 57 L12 68 L14 73 L40 66 Z" />
			<path d="M14 70 L2 82 M16 71 L8 84" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
		{/if}
		<ellipse cx="56" cy="62" rx={rx} ry={ry} />
		<circle cx="74" cy="40" r={cfg.headR} />
		{#if cfg.bill === 'cone'}
			<path d="M82 35 L94 40 L82 44 Z" />
		{:else if cfg.bill === 'heavy'}
			<path d="M81 32 L97 40 L81 46 Z" />
		{:else if cfg.bill === 'longthin'}
			<path d="M82 38 L100 42 L100 44 L82 41 Z" />
		{:else if cfg.bill === 'curve'}
			<path d="M82 35 Q94 39 97 48 L93 49 Q89 43 81 39 Z" />
		{:else if cfg.bill === 'hook'}
			<path d="M81 32 Q94 33 94 41 Q89 41 88 39 L88 45 Q83 44 81 40 Z" />
		{:else if cfg.bill === 'dagger'}
			<path d="M83 34 L102 41 L83 46 Z" />
		{/if}
		{#if cfg.crest === 'tuft'}
			<path d="M67 32 L71 21 L76 31 Z" />
		{:else if cfg.crest === 'fan'}
			<path d="M70 31 L60 15 M70 31 L66 13 M70 31 L72 12 M70 31 L78 14 M70 31 L83 17" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" />
		{/if}
		{#if cfg.legs === 'long'}
			<path d="M50 75 L50 92 M60 75 L60 92 M46 92 L54 92 M56 92 L64 92" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" />
		{:else}
			<path d="M51 75 L51 88 M61 75 L61 88 M47 88 L55 88 M57 88 L65 88" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" />
		{/if}
	{/if}
</svg>
