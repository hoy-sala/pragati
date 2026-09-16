<script lang="ts">
  import { KA_BIRDS, KA_FAMILIES, KA_REGION, type KaBird, type BirdShapeKind } from "$lib/data/kaBirds";
  import { KA_PHOTOS, KA_PHOTO_COUNT } from "$lib/data/kaPhotos";
  import { KA_KANNADA } from "$lib/data/kaKannada";
  import BirdShape from "$lib/components/BirdShape.svelte";
  import Button from "$lib/components/Button.svelte";
  import Select from "$lib/components/Select.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import IucnBadge from "$lib/components/IucnBadge.svelte";
  import { KA_DETAILS } from "$lib/data/kaDetails";
  import { Bird, Layers, Camera, ExternalLink } from "lucide-svelte";

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
  let kaPage = $state(1);
  let kaSelected = $state<KaBird | null>(null);
  let kaOpen = $state(false);
  const kaFamilyOptions = [{ id: "", name: "All families" }, ...KA_FAMILIES.map((f) => ({ id: f, name: f }))];
  let kaFiltered = $derived(
    KA_BIRDS.filter((b) => {
      if (kaFamily && b.family !== kaFamily) return false;
      if (kaShape && b.shape !== kaShape) return false;
      if (!kaSearch.trim()) return true;
      const q = kaSearch.trim().toLowerCase();
      const kn = KA_KANNADA[b.code] || "";
      return b.com.toLowerCase().includes(q) || b.sci.toLowerCase().includes(q) || b.family.toLowerCase().includes(q) || kn.includes(kaSearch.trim());
    }).sort((a, b) => a.com.localeCompare(b.com))
  );
  let kaShapes = $derived(
    (() => {
      const m = new Map<BirdShapeKind, number>();
      for (const b of KA_BIRDS) m.set(b.shape, (m.get(b.shape) || 0) + 1);
      return [...m.entries()].sort((a, b) => b[1] - a[1]);
    })()
  );
  let kaSimilar = $derived(
    (() => {
      const sel = kaSelected;
      if (!sel) return [] as KaBird[];
      return KA_BIRDS.filter((b) => b.code !== sel.code && b.family === sel.family && b.shape === sel.shape).slice(0, 6);
    })()
  );
  let kaPaged = $derived(kaFiltered.slice((kaPage - 1) * KA_PAGE_SIZE, kaPage * KA_PAGE_SIZE));
  $effect(() => {
    void kaSearch;
    void kaFamily;
    void kaShape;
    kaPage = 1;
  });
</script>

<svelte:head><title>Bird Catalog — Pragati</title></svelte:head>

<div class="max-w-6xl mx-auto space-y-8">
  <section class="bg-gradient-to-br from-emerald-600 via-teal-600 to-cyan-600 text-white rounded-3xl px-6 py-8 sm:px-10 sm:py-12 shadow-lg">
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

  <div class="bg-white rounded-2xl border border-slate-200 shadow-sm p-4 no-print">
    <div class="flex flex-wrap gap-3 items-end">
      <div class="flex-1 min-w-48">
        <SearchFilter bind:value={kaSearch} placeholder="Search English, Kannada, scientific name, family..." />
      </div>
      <div class="w-60">
        <Select bind:value={kaFamily} options={kaFamilyOptions} placeholder="All families" />
      </div>
    </div>
    <div class="flex gap-2 overflow-x-auto pt-3 pb-1">
      <button
        onclick={() => (kaShape = "")}
        class="shrink-0 text-[11px] font-semibold px-3 py-1.5 rounded-full border transition-colors {kaShape === ''
          ? 'bg-slate-900 text-white border-slate-900'
          : 'bg-white text-slate-600 border-slate-300 hover:bg-slate-50'}"
      >
        All shapes
      </button>
      {#each kaShapes as [shape, n]}
        <button
          onclick={() => (kaShape = kaShape === shape ? "" : shape)}
          class="shrink-0 inline-flex items-center gap-1.5 text-[11px] font-semibold px-3 py-1.5 rounded-full border transition-colors {kaShape === shape
            ? 'bg-slate-900 text-white border-slate-900'
            : 'bg-white text-slate-600 border-slate-300 hover:bg-slate-50'}"
        >
          <BirdShape shape={shape} class="w-4 h-4" />
          {kaShapeLabels[shape] ?? shape} · {n}
        </button>
      {/each}
    </div>
  </div>

  <p class="text-xs font-medium text-slate-500">
    Showing {KA_PAGE_SIZE >= kaFiltered.length ? kaFiltered.length : KA_PAGE_SIZE} of {kaFiltered.length} species
  </p>

  {#if kaFiltered.length === 0}
    <div class="bg-white rounded-2xl border border-slate-200 shadow-sm">
      <EmptyState icon={Bird} title="No species match." hint="Try a different name or family." />
    </div>
  {:else}
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
      {#each kaPaged as b (b.code)}
        <button
          onclick={() => {
            kaSelected = b;
            kaOpen = true;
          }}
          class="group text-left bg-white rounded-2xl border border-slate-200 overflow-hidden transition-all hover:border-primary-300 hover:shadow-md active:scale-[0.99] flex flex-col"
        >
          <div class="h-28 sm:h-32 flex items-center justify-center bg-gradient-to-br from-emerald-50 to-teal-100 overflow-hidden">
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
        </button>
      {/each}
    </div>
    <div class="bg-white rounded-2xl border border-slate-200 shadow-sm no-print">
      <Pagination page={kaPage} total={kaFiltered.length} pageSize={KA_PAGE_SIZE} onChange={(p) => (kaPage = p)} />
    </div>
    <p class="text-[11px] text-slate-400 text-center">Photos by Wikimedia Commons contributors · click any bird for credit and species links</p>
  {/if}
</div>

{#if kaSelected}
  <Modal bind:open={kaOpen} title={kaSelected.com} maxWidth="max-w-lg">
    <div class="space-y-4">
      {#if KA_PHOTOS[kaSelected.code]}
        <div class="relative rounded-2xl h-48 sm:h-60 bg-gradient-to-br from-emerald-50 to-teal-100 overflow-hidden">
          <img
            src={KA_PHOTOS[kaSelected.code].photo}
            alt={kaSelected.com}
            class="absolute inset-0 w-full h-full object-contain p-2"
          />
          <span class="absolute bottom-2 right-2 text-[10px] px-1.5 py-0.5 rounded bg-black/60 text-white truncate max-w-[80%]">
            © {KA_PHOTOS[kaSelected.code].credit}
          </span>
        </div>
      {:else}
        <div class="relative rounded-2xl h-44 flex items-center justify-center bg-gradient-to-br from-emerald-50 to-teal-100 overflow-hidden">
          <BirdShape shape={kaSelected.shape} class="w-32 h-32 text-slate-700/80" />
        </div>
      {/if}

      <div>
        {#if KA_KANNADA[kaSelected.code]}
          <p class="text-base font-semibold text-emerald-700">{KA_KANNADA[kaSelected.code]}</p>
        {/if}
        <p class="text-xs text-slate-400 italic">{kaSelected.sci}</p>
        <div class="flex flex-wrap gap-1.5 mt-2">
          <span class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-slate-100 text-slate-600">{kaSelected.family}</span>
          <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-primary-50 text-primary-700">{kaSelected.order}</span>
          {#if KA_DETAILS[kaSelected.code]?.iucn}
            <IucnBadge code={KA_DETAILS[kaSelected.code].iucn} />
          {/if}
        </div>
      </div>

      {#if KA_DETAILS[kaSelected.code]?.extract}
        <p class="text-sm text-slate-600 leading-relaxed">{KA_DETAILS[kaSelected.code].extract}</p>
        <p class="text-[11px] text-slate-400 -mt-2">Description from Wikipedia (CC BY-SA) · full article linked below</p>
      {/if}

      <dl class="grid grid-cols-1 gap-2 text-sm">
        <div class="flex items-start gap-2">
          <dt class="w-28 shrink-0 text-slate-400 font-medium">eBird code</dt>
          <dd class="text-slate-700 font-mono text-xs leading-relaxed pt-0.5">{kaSelected.code}</dd>
        </div>
        <div class="flex items-start gap-2">
          <dt class="w-28 shrink-0 text-slate-400 font-medium">Recorded in</dt>
          <dd class="text-slate-700">{KA_REGION} — from the eBird checklist for the state</dd>
        </div>
      </dl>

      <div class="text-sm text-slate-600 bg-slate-50 border border-slate-200 rounded-xl px-3 py-2.5">
        Call recordings, range maps and full species history are available on this species' eBird and Wikipedia pages.
      </div>

      <a
        href="https://ebird.org/species/{kaSelected.code}"
        target="_blank"
        rel="noreferrer"
        class="flex items-center justify-center gap-2 w-full px-3.5 py-2 text-sm font-medium rounded-lg bg-slate-900 text-white hover:bg-slate-700 active:bg-slate-800 transition-colors"
      >
        <ExternalLink size={16} />
        Photos, calls & range map on eBird
      </a>
      {#if KA_PHOTOS[kaSelected.code]?.wiki}
        <a
          href={KA_PHOTOS[kaSelected.code].wiki}
          target="_blank"
          rel="noreferrer"
          class="flex items-center justify-center gap-2 w-full px-3.5 py-2 text-sm font-medium rounded-lg border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors"
        >
          Read the Wikipedia species page <ExternalLink size={15} />
        </a>
      {/if}
      <a
        href="/birds/{kaSelected.code}"
        class="flex items-center justify-center gap-2 w-full px-3.5 py-2 text-sm font-medium rounded-lg border border-primary-300 text-primary-700 hover:bg-primary-50 transition-colors"
      >
        Open full species details
      </a>
      {#if kaSimilar.length > 0}
        <div>
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wide mb-2">Similar species</p>
          <div class="grid grid-cols-3 gap-2">
            {#each kaSimilar as s (s.code)}
              <button
                onclick={() => (kaSelected = s)}
                class="group text-left bg-white rounded-xl border border-slate-200 overflow-hidden hover:border-primary-300 transition-colors"
              >
                <div class="h-16 flex items-center justify-center bg-gradient-to-br from-emerald-50 to-teal-100 overflow-hidden">
                  {#if KA_PHOTOS[s.code]}
                    <img src={KA_PHOTOS[s.code].photo} alt={s.com} loading="lazy" class="w-full h-full object-contain p-0.5" />
                  {:else}
                    <BirdShape shape={s.shape} class="w-8 h-8 text-slate-700/80" />
                  {/if}
                </div>
                <div class="px-1.5 py-1 text-[10px] font-semibold text-slate-700 leading-tight truncate" title={s.com}>{s.com}</div>
              </button>
            {/each}
          </div>
        </div>
      {/if}
      <p class="text-[11px] text-slate-400 text-center -mt-2">Photo © {KA_PHOTOS[kaSelected.code]?.credit ?? 'Wikimedia Commons'}</p>
    </div>
    {#snippet footer()}
      <Button variant="ghost" onclick={() => (kaOpen = false)}>Close</Button>
    {/snippet}
  </Modal>
{/if}