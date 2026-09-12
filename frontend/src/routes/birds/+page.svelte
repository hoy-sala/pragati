<script lang="ts">
  import { BIRDS, BIRD_COLORS, type Bird as BirdInfo } from "$lib/data/birds";
  import BirdShape from "$lib/components/BirdShape.svelte";
  import Button from "$lib/components/Button.svelte";
  import Select from "$lib/components/Select.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import { Check, Eye, Sparkles, Volume2, MapPin, ListChecks, Bird, ExternalLink, Telescope, Smartphone, Camera } from "lucide-svelte";
  import { onMount } from "svelte";

  const SPOTTED_KEY = "pragati:birds:spotted";

  let search = $state("");
  let filterSize = $state("");
  let filterColor = $state("");
  let spottedOnly = $state(false);
  let spotted = $state<string[]>([]);
  let selected = $state<BirdInfo | null>(null);
  let detailOpen = $state(false);

  const sizeOptions = [
    { id: "", name: "All sizes" },
    { id: "S", name: "Small (sparrow-sized)" },
    { id: "M", name: "Medium (myna-sized)" },
    { id: "L", name: "Large (crow-sized+)" },
  ];
  const colorOptions = [
    { id: "", name: "All colours" },
    ...BIRD_COLORS,
  ];

  onMount(() => {
    try {
      const raw = localStorage.getItem(SPOTTED_KEY);
      if (raw) spotted = JSON.parse(raw);
    } catch {
      spotted = [];
    }
  });

  function saveSpotted() {
    try {
      localStorage.setItem(SPOTTED_KEY, JSON.stringify(spotted));
    } catch {
      /* ignore */
    }
  }

  function toggleSpotted(id: string) {
    spotted = spotted.includes(id)
      ? spotted.filter((s) => s !== id)
      : [...spotted, id];
    saveSpotted();
  }

  let filtered = $derived(
    BIRDS.filter((b) => {
      if (filterSize && b.size !== filterSize) return false;
      if (filterColor && !b.colors.includes(filterColor)) return false;
      if (spottedOnly && !spotted.includes(b.id)) return false;
      if (!search.trim()) return true;
      const q = search.trim().toLowerCase();
      return (
        b.name_en.toLowerCase().includes(q) ||
        (b.name_kn ?? "").includes(search.trim()) ||
        b.scientific.toLowerCase().includes(q)
      );
    }),
  );

  function openDetail(b: BirdInfo) {
    selected = b;
    detailOpen = true;
  }

  function shapeProps(b: BirdInfo): { tail?: any; crest?: any } {
    if (b.id === "red-whiskered-bulbul") return { crest: "tuft" as const };
    if (b.id === "common-tailorbird") return { tail: "cock" as const };
    if (b.id === "oriental-magpie-robin") return { tail: "long" as const };
    return {};
  }
</script>

<svelte:head><title>Bird Catalog — Pragati</title></svelte:head>

<div class="max-w-6xl mx-auto space-y-6">
  <div>
    <p class="text-xs font-semibold tracking-widest text-primary-600 uppercase">MDRS Bahaddurghatta campus</p>
    <h1 class="text-2xl sm:text-3xl font-bold text-slate-900 mt-1">Campus Bird Catalog</h1>
    <p class="text-sm text-slate-500 mt-1">
      Know every feathered neighbour. Spot them, identify them, tick them off — and fall in love with birdwatching.
    </p>
    <div class="mt-3 flex items-center gap-3 max-w-md">
      <div class="flex-1 h-2.5 rounded-full bg-slate-200 overflow-hidden">
        <div
          class="h-full bg-gradient-to-r from-emerald-500 to-teal-500 rounded-full transition-all"
          style="width: {Math.round((spotted.length / BIRDS.length) * 100)}%"
        ></div>
      </div>
      <span class="text-xs font-semibold text-slate-600 whitespace-nowrap">
        {spotted.length}/{BIRDS.length} spotted
      </span>
    </div>
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
    <div class="flex flex-wrap gap-3 items-end">
      <div class="flex-1 min-w-48">
        <SearchFilter bind:value={search} placeholder="Search name, Kannada name, scientific name..." />
      </div>
      <div class="w-44">
        <Select bind:value={filterSize} options={sizeOptions} placeholder="All sizes" />
      </div>
      <div class="w-40">
        <Select bind:value={filterColor} options={colorOptions} placeholder="All colours" />
      </div>
      <Button
        variant={spottedOnly ? "primary" : "secondary"}
        icon={ListChecks}
        onclick={() => (spottedOnly = !spottedOnly)}
      >
        Spotted
      </Button>
    </div>
  </div>

  {#if filtered.length === 0}
    <div class="bg-white rounded-xl border border-slate-200">
      <EmptyState
        icon={Bird}
        title="No birds match."
        hint="Try a different name, size or colour."
      />
    </div>
  {:else}
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-4">
      {#each filtered as b (b.id)}
        <button
          onclick={() => openDetail(b)}
          class="group text-left bg-white rounded-2xl border border-slate-200 overflow-hidden hover:border-primary-300 hover:shadow-md active:scale-[0.99] transition-all"
        >
          <div
            class="relative h-32 sm:h-36 flex items-center justify-center"
            style="background: linear-gradient(135deg, {b.tile[0]}, {b.tile[1]})"
          >
            <BirdShape shape={b.shape} {...shapeProps(b)} class="w-20 h-20 sm:w-24 sm:h-24 text-slate-800/90 transition-transform group-hover:scale-105" />
            <span class="absolute top-2 left-2 text-[10px] font-bold px-1.5 py-0.5 rounded bg-white/80 text-slate-600">{b.size}</span>
            {#if spotted.includes(b.id)}
              <span class="absolute top-2 right-2 w-6 h-6 rounded-full bg-emerald-500 text-white flex items-center justify-center" aria-label="Spotted">
                <Check size={14} strokeWidth={3} />
              </span>
            {/if}
          </div>
          <div class="p-3">
            <div class="text-sm font-semibold text-slate-800 leading-tight">{b.name_en}</div>
            {#if b.name_kn}
              <div class="text-xs text-slate-500 font-kannada mt-0.5">{b.name_kn}</div>
            {/if}
            <div class="text-[11px] text-slate-400 italic truncate">{b.scientific}</div>
          </div>
        </button>
      {/each}
    </div>
  {/if}

  <div>
    <h2 class="text-sm font-semibold text-slate-800 mb-3">Go further with Cornell Lab</h2>
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
      <div class="bg-white rounded-2xl border border-slate-200 p-4 flex flex-col gap-2">
        <Telescope size={20} class="text-primary-600" />
        <div class="text-sm font-semibold text-slate-800">eBird Karnataka</div>
        <p class="text-xs text-slate-500 flex-1">See what birders across Karnataka report — real sightings near you, updated daily.</p>
        <a
          href="https://ebird.org/region/IN-KA"
          target="_blank"
          rel="noreferrer"
          class="inline-flex items-center gap-1.5 text-xs font-semibold text-primary-700 hover:text-primary-800"
        >
          Explore sightings <ExternalLink size={13} />
        </a>
      </div>
      <div class="bg-white rounded-2xl border border-slate-200 p-4 flex flex-col gap-2">
        <Smartphone size={20} class="text-primary-600" />
        <div class="text-sm font-semibold text-slate-800">Merlin Bird ID — free app</div>
        <p class="text-xs text-slate-500 flex-1">Point your phone at a bird or record its song — Merlin identifies it. Made by Cornell Lab.</p>
        <div class="flex flex-wrap gap-2">
          <a
            href="https://play.google.com/store/apps/details?id=com.labs.merlinbirdid.app"
            target="_blank"
            rel="noreferrer"
            class="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 rounded-lg bg-slate-900 text-white hover:bg-slate-700 transition-colors"
          >
            Android <ExternalLink size={13} />
          </a>
          <a
            href="https://apps.apple.com/us/app/merlin-bird-id-by-cornell-lab/id773457673"
            target="_blank"
            rel="noreferrer"
            class="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 rounded-lg border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors"
          >
            iPhone <ExternalLink size={13} />
          </a>
        </div>
      </div>
      <div class="bg-white rounded-2xl border border-slate-200 p-4 flex flex-col gap-2">
        <Camera size={20} class="text-primary-600" />
        <div class="text-sm font-semibold text-slate-800">Macaulay Library</div>
        <p class="text-xs text-slate-500 flex-1">The world's largest archive of bird photos, calls and videos — every eBird photo lives here.</p>
        <a
          href="https://www.macaulaylibrary.org"
          target="_blank"
          rel="noreferrer"
          class="inline-flex items-center gap-1.5 text-xs font-semibold text-primary-700 hover:text-primary-800"
        >
          Browse the archive <ExternalLink size={13} />
        </a>
      </div>
    </div>
    <p class="text-[11px] text-slate-400 mt-3 text-center">
      Species pages cross-checked with eBird · Cornell Lab of Ornithology
    </p>
  </div>
</div>

{#if selected}
  <Modal bind:open={detailOpen} title={selected.name_en} maxWidth="max-w-lg">
    <div class="space-y-4">
      <div
        class="rounded-2xl h-44 flex items-center justify-center relative overflow-hidden"
        style="background: linear-gradient(135deg, {selected.tile[0]}, {selected.tile[1]})"
      >
        {#if selected.photo_url}
          <img
            src={selected.photo_url}
            alt={selected.name_en}
            class="absolute inset-0 w-full h-full object-cover"
            loading="lazy"
          />
        {:else}
          <BirdShape shape={selected.shape} {...shapeProps(selected)} class="w-32 h-32 text-slate-800/90" />
        {/if}
        {#if selected.photo_credit}
          <span class="absolute bottom-2 right-2 text-[10px] px-1.5 py-0.5 rounded bg-black/55 text-white">
            © {selected.photo_credit}
          </span>
        {/if}
      {#if selected && spotted.includes(selected.id)}
          <span class="absolute top-3 right-3 flex items-center gap-1 text-xs font-semibold px-2.5 py-1 rounded-full bg-emerald-500 text-white">
            <Check size={13} strokeWidth={3} /> Spotted
          </span>
        {/if}
      </div>

      <div>
        {#if selected.name_kn}
          <p class="text-lg text-slate-700 font-kannada">{selected.name_kn}</p>
        {/if}
        <p class="text-xs text-slate-400 italic">{selected.scientific}</p>
        <div class="flex flex-wrap gap-1.5 mt-2">
          <span class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-slate-100 text-slate-600">
            {selected.size === "S" ? "Small" : selected.size === "M" ? "Medium" : "Large"}
          </span>
          {#each selected.colors as c}
            <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-primary-50 text-primary-700 capitalize">{c}</span>
          {/each}
        </div>
      </div>

      <div>
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-500 mb-1.5 flex items-center gap-1.5">
          <Eye size={13} /> How to identify
        </h3>
        <ul class="space-y-1.5">
          {#each selected.marks as m}
            <li class="flex items-start gap-2 text-sm text-slate-700">
              <Check size={15} strokeWidth={3} class="text-emerald-600 shrink-0 mt-0.5" />
              <span>{m}</span>
            </li>
          {/each}
        </ul>
      </div>

      {#if selected.audio_url}
        <div class="bg-slate-900 rounded-xl px-3 py-2.5 space-y-1">
          <div class="flex items-center gap-2 text-xs text-slate-300">
            <Volume2 size={14} />
            <span>Hear it{#if selected.audio_credit} · © {selected.audio_credit}{/if}</span>
          </div>
          <audio controls preload="none" src={selected.audio_url} class="w-full h-8"></audio>
        </div>
      {:else if selected.call}
        <div class="flex items-start gap-2 text-sm bg-amber-50 border border-amber-200 rounded-xl px-3 py-2.5">
          <Volume2 size={16} class="text-amber-600 shrink-0 mt-0.5" />
          <span class="text-slate-700"><strong>Listen for:</strong> {selected.call}</span>
        </div>
      {/if}

      <div class="flex items-start gap-2 text-sm text-slate-600">
        <MapPin size={16} class="text-slate-400 shrink-0 mt-0.5" />
        <span><strong>Look at:</strong> {selected.where}</span>
      </div>

      <div class="flex items-start gap-2 text-sm bg-violet-50 border border-violet-200 rounded-xl px-3 py-2.5">
        <Sparkles size={16} class="text-violet-600 shrink-0 mt-0.5" />
        <span class="text-slate-700">{selected.fact}</span>
      </div>

      <p class="text-sm text-slate-600 border-l-4 border-emerald-400 pl-3">
        <strong>Spotter tip:</strong> {selected.tip}
      </p>

      <a
        href="https://ebird.org/species/{selected.ebird}"
        target="_blank"
        rel="noreferrer"
        class="flex items-center justify-center gap-2 w-full px-3.5 py-2 text-sm font-medium rounded-lg bg-slate-900 text-white hover:bg-slate-700 active:bg-slate-800 transition-colors"
      >
        <ExternalLink size={16} />
        Photos, calls & range map on eBird
      </a>
      <p class="text-[11px] text-slate-400 text-center -mt-2">Species data: eBird · Cornell Lab of Ornithology</p>
    </div>
    {#snippet footer()}
      <Button variant="ghost" onclick={() => (detailOpen = false)}>Close</Button>
      {#if selected && spotted.includes(selected.id)}
        <Button variant="secondary" onclick={() => toggleSpotted(selected!.id)}>Unmark spotted</Button>
      {:else}
        <Button icon={Check} onclick={() => toggleSpotted(selected!.id)}>I spotted this!</Button>
      {/if}
    {/snippet}
  </Modal>
{/if}
