<script lang="ts">
  import { KA_BIRDS, KA_FAMILIES, KA_REGION, type BirdShapeKind } from "$lib/data/kaBirds";
  import { KA_PHOTOS, KA_PHOTO_COUNT } from "$lib/data/kaPhotos";
  import { KA_KANNADA } from "$lib/data/kaKannada";
  import { KA_COLORS, KA_SWATCHES } from "$lib/data/kaColors";
  import { KA_DETAILS } from "$lib/data/kaDetails";
  import BirdShape from "$lib/components/BirdShape.svelte";
  import Select from "$lib/components/Select.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { Bird, Layers, Camera, SlidersHorizontal, X } from "lucide-svelte";

  const KA_PAGE_SIZE = 40;
  const kaFamilyCount = new Set(KA_BIRDS.map((b) => b.family)).size;
  const kaShapeLabels: Record<string, string> = {
    songbird: 'Songbirds', crow: 'Crows & jays', parrot: 'Parrots', dove: 'Doves & pigeons',
    longtail: 'Long-tailed', hoopoe: 'Hoopoes', barbet: 'Barbets', drongo: 'Drongos',
    sunbird: 'Sunbirds', owl: 'Owls', raptor: 'Birds of prey', wader: 'Waders & waterbirds',
    kingfisher: 'Kingfishers', slim: 'Swifts & swallows', peafowl: 'Peafowl'
  };
  let kaSearch = $state("");
  let kaFamily = $state("");
  let kaShape = $state("");
  let kaColors = $state<string[]>([]);
  let kaPage = $state(1);
  let kaFiltersOpen = $state(false);
  const kaFamilyOptions = [{ id: "", name: "All families" }, ...KA_FAMILIES.map((f) => ({ id: f, name: f }))];
  let kaFiltered = $derived(
    KA_BIRDS.filter((b) => {
      if (kaFamily && b.family !== kaFamily) return false;
      if (kaShape && b.shape !== kaShape) return false;
      if (kaColors.length > 0) {
        const cols = KA_COLORS[b.code] || [];
        if (!cols.some((c) => kaColors.includes(c))) return false;
      }
      if (!kaSearch.trim()) return true;
      const q = kaSearch.trim().toLowerCase();
      const kn = KA_KANNADA[b.code] || "";
      return b.com.toLowerCase().includes(q) || b.sci.toLowerCase().includes(q) || b.family.toLowerCase().includes(q) || kn.includes(kaSearch.trim());
    }).sort((a, b) => a.com.localeCompare(b.com))
  );
  const kaThreatened: Record<string, string> = { VU: "Vulnerable", EN: "Endangered", CR: "Critically Endangered" };
  const kaThreatStyle: Record<string, string> = { VU: "bg-orange-500", EN: "bg-red-500", CR: "bg-rose-700" };
  let kaActiveCount = $derived(
    (kaSearch.trim() !== "" ? 1 : 0) + (kaFamily !== "" ? 1 : 0) + (kaShape !== "" ? 1 : 0) + kaColors.length
  );
  function kaClearFilters() {
    kaSearch = "";
    kaFamily = "";
    kaShape = "";
    kaColors = [];
  }
  let kaShapes = $derived(
    (() => {
      const m = new Map<BirdShapeKind, number>();
      for (const b of KA_BIRDS) m.set(b.shape, (m.get(b.shape) || 0) + 1);
      return [...m.entries()].sort((a, b) => b[1] - a[1]);
    })()
  );
  let kaPaged = $derived(kaFiltered.slice((kaPage - 1) * KA_PAGE_SIZE, kaPage * KA_PAGE_SIZE));
  $effect(() => {
    void kaSearch;
    void kaFamily;
    void kaShape;
    void kaColors;
    kaPage = 1;
  });
</script>

<svelte:head><title>Bird Catalog — Pragati</title></svelte:head>

<div class="max-w-6xl mx-auto space-y-6">
  <section class="bg-gradient-to-br from-emerald-600 via-teal-600 to-cyan-600 text-white rounded-3xl px-6 py-8 sm:px-10 sm:py-10 shadow-lg">
    <p class="text-xs font-semibold tracking-widest text-emerald-100 uppercase">Karnataka checklist</p>
    <h1 class="text-3xl sm:text-4xl font-bold mt-2">Birds of Karnataka</h1>
    <p class="text-sm sm:text-base text-emerald-50 mt-3 max-w-2xl leading-relaxed">
      {KA_BIRDS.length} species recorded for {KA_REGION} — mirrored once from the eBird API into this app. Nothing is sent back to eBird. Photographs by Wikimedia Commons contributors.
    </p>
    <div class="flex flex-wrap gap-2.5 mt-6">
      <span class="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 rounded-full bg-white/15">
        <Bird size={13} /> {KA_BIRDS.length} species
      </span>
      <span class="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 rounded-full bg-white/15">
        <Layers size={13} /> {kaFamilyCount} families
      </span>
      <span class="inline-flex items-center gap-1.5 text-xs font-semibold px-3 py-1.5 rounded-full bg-white/15">
        <Camera size={13} /> {KA_PHOTO_COUNT} photos
      </span>
    </div>
  </section>

  <div class="flex flex-col lg:flex-row gap-6 items-start">
    <aside class="w-full lg:w-64 lg:shrink-0 lg:sticky lg:top-3 lg:max-h-[calc(100dvh-2rem)] lg:overflow-y-auto no-print">
      <button
        onclick={() => (kaFiltersOpen = !kaFiltersOpen)}
        class="lg:hidden w-full inline-flex items-center justify-center gap-2 text-sm font-semibold px-4 py-2.5 rounded-xl bg-white border border-slate-200 shadow-sm text-slate-700"
        aria-expanded={kaFiltersOpen}
      >
        <SlidersHorizontal size={15} />
        Filters
        {#if kaActiveCount > 0}
          <span class="text-[11px] font-bold px-1.5 py-0.5 rounded-full bg-slate-900 text-white">{kaActiveCount}</span>
        {/if}
      </button>

      <div class="{kaFiltersOpen ? 'block' : 'hidden'} lg:block bg-white rounded-2xl border border-slate-200 shadow-sm p-4 space-y-5 mt-3 lg:mt-0">
        <div class="flex items-center justify-between">
          <p class="text-xs font-bold text-slate-500 uppercase tracking-wide inline-flex items-center gap-1.5">
            <SlidersHorizontal size={13} /> Filters
          </p>
          {#if kaActiveCount > 0}
            <button
              onclick={kaClearFilters}
              class="inline-flex items-center gap-1 text-[11px] font-semibold text-slate-500 hover:text-slate-800"
            >
              <X size={12} /> Clear all
            </button>
          {/if}
        </div>

        <div>
          <p class="text-[11px] font-semibold text-slate-500 uppercase tracking-wide mb-1.5">Search</p>
          <SearchFilter bind:value={kaSearch} placeholder="English, Kannada, family..." />
        </div>

        <div>
          <p class="text-[11px] font-semibold text-slate-500 uppercase tracking-wide mb-1.5">Family</p>
          <Select bind:value={kaFamily} options={kaFamilyOptions} placeholder="All families" />
        </div>

        <div>
          <p class="text-[11px] font-semibold text-slate-500 uppercase tracking-wide mb-1.5">Shape</p>
          <div class="space-y-1">
            <button
              onclick={() => (kaShape = "")}
              class="w-full flex items-center justify-between text-xs font-semibold px-2.5 py-1.5 rounded-lg transition-colors {kaShape === ''
                ? 'bg-slate-900 text-white'
                : 'text-slate-600 hover:bg-slate-100'}"
            >
              All shapes
              <span class="{kaShape === '' ? 'text-slate-300' : 'text-slate-400'}">{KA_BIRDS.length}</span>
            </button>
            {#each kaShapes as [shape, n]}
              <button
                onclick={() => (kaShape = kaShape === shape ? "" : shape)}
                class="w-full flex items-center gap-2 text-xs font-semibold px-2.5 py-1.5 rounded-lg transition-colors {kaShape === shape
                  ? 'bg-slate-900 text-white'
                  : 'text-slate-600 hover:bg-slate-100'}"
              >
                <BirdShape shape={shape} class="w-5 h-5 shrink-0" />
                <span class="flex-1 text-left truncate">{kaShapeLabels[shape] ?? shape}</span>
                <span class="{kaShape === shape ? 'text-slate-300' : 'text-slate-400'}">{n}</span>
              </button>
            {/each}
          </div>
        </div>

        <div>
          <p
            class="text-[11px] font-semibold text-slate-500 uppercase tracking-wide mb-1.5"
            title="Approximate plumage colours (breeding male where sexes differ)"
          >
            Colour
          </p>
          <div class="flex flex-wrap gap-2">
            {#each KA_SWATCHES as sw}
              {@const active = kaColors.includes(sw.id)}
              <button
                onclick={() => (kaColors = active ? kaColors.filter((c) => c !== sw.id) : [...kaColors, sw.id])}
                title={sw.label}
                aria-label="Filter by {sw.label}"
                aria-pressed={active}
                class="w-8 h-8 rounded-full border transition-all {active
                  ? 'border-slate-900 ring-2 ring-slate-900 ring-offset-2 scale-110'
                  : 'border-slate-300 hover:scale-110 hover:border-slate-500'}"
                style="background-color:{sw.hex}"
              ></button>
            {/each}
          </div>
          {#if kaColors.length > 0}
            <p class="text-[11px] text-slate-500 mt-1.5">
              {kaColors.map((c) => KA_SWATCHES.find((s) => s.id === c)?.label).join(", ")}
            </p>
          {/if}
        </div>
      </div>
    </aside>

    <div class="flex-1 min-w-0 w-full space-y-4">
      <p class="text-xs font-medium text-slate-500">
        Showing {KA_PAGE_SIZE >= kaFiltered.length ? kaFiltered.length : KA_PAGE_SIZE} of {kaFiltered.length} species
        {#if kaActiveCount > 0}
          <button onclick={kaClearFilters} class="ml-2 font-semibold text-primary-700 hover:text-primary-800">
            Clear filters ✕
          </button>
        {/if}
      </p>

      {#if kaFiltered.length === 0}
        <div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 text-center">
          <EmptyState icon={Bird} title="No species match." hint="Try a different name, colour, shape or family." />
          <button
            onclick={kaClearFilters}
            class="mt-2 text-xs font-semibold px-4 py-2 rounded-lg bg-slate-900 text-white hover:bg-slate-700 transition-colors"
          >
            Clear all filters
          </button>
        </div>
      {:else}
        <div class="grid grid-cols-2 md:grid-cols-3 2xl:grid-cols-4 gap-4">
          {#each kaPaged as b (b.code)}
            <a
              href="/birds/{b.code}"
              class="group text-left bg-white rounded-2xl border border-slate-200 overflow-hidden transition-all hover:border-primary-300 hover:shadow-md active:scale-[0.99] flex flex-col"
            >
              <div class="relative h-28 sm:h-32 flex items-center justify-center bg-gradient-to-br from-emerald-50 to-teal-100 overflow-hidden">
                {#if KA_DETAILS[b.code]?.iucn && kaThreatened[KA_DETAILS[b.code].iucn]}
                  <span
                    title="IUCN: {kaThreatened[KA_DETAILS[b.code].iucn]}"
                    class="absolute top-2 left-2 z-10 text-[10px] font-bold px-1.5 py-0.5 rounded text-white {kaThreatStyle[KA_DETAILS[b.code].iucn]}"
                  >
                    {KA_DETAILS[b.code].iucn}
                  </span>
                {/if}
                {#if KA_PHOTOS[b.code]}
                  <img
                    src={KA_PHOTOS[b.code].photo}
                    alt={b.com}
                    loading="lazy"
                    class="w-full h-full object-contain p-1 transition-transform group-hover:scale-105"
                  />
                {:else}
                  <BirdShape shape={b.shape} class="w-16 h-16 sm:w-20 sm:h-20 text-slate-700/80 transition-transform group-hover:scale-105" />
                {/if}
              </div>
              <div class="p-3 flex flex-col gap-0.5">
                <div class="text-sm font-semibold text-slate-800 leading-tight">{b.com}</div>
                {#if KA_KANNADA[b.code]}
                  <div class="text-xs font-medium text-emerald-700">{KA_KANNADA[b.code]}</div>
                {/if}
                <div class="text-[11px] text-slate-400 italic truncate">{b.sci}</div>
                <div class="text-[10px] text-slate-500 mt-0.5 truncate">{b.family}</div>
              </div>
            </a>
          {/each}
        </div>
        <div class="bg-white rounded-2xl border border-slate-200 shadow-sm no-print">
          <Pagination page={kaPage} total={kaFiltered.length} pageSize={KA_PAGE_SIZE} onChange={(p) => (kaPage = p)} />
        </div>
        <p class="text-[11px] text-slate-400 text-center">Photos by Wikimedia Commons contributors · open any bird for details, photos and links</p>
      {/if}
    </div>
  </div>
</div>
