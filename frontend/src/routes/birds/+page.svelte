<script lang="ts">
  import { KA_BIRDS, KA_FAMILIES, KA_REGION, type KaBird } from "$lib/data/kaBirds";
  import { KA_PHOTOS } from "$lib/data/kaPhotos";
  import BirdShape from "$lib/components/BirdShape.svelte";
  import Button from "$lib/components/Button.svelte";
  import Select from "$lib/components/Select.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { Bird, ExternalLink, BarChart3 } from "lucide-svelte";

  const KA_PAGE_SIZE = 40;
  let kaSearch = $state("");
  let kaFamily = $state("");
  let kaPage = $state(1);
  let kaSelected = $state<KaBird | null>(null);
  let kaOpen = $state(false);
  const kaFamilyOptions = [{ id: "", name: "All families" }, ...KA_FAMILIES.map((f) => ({ id: f, name: f }))];
  let kaFiltered = $derived(
    KA_BIRDS.filter((b) => {
      if (kaFamily && b.family !== kaFamily) return false;
      if (!kaSearch.trim()) return true;
      const q = kaSearch.trim().toLowerCase();
      return b.com.toLowerCase().includes(q) || b.sci.toLowerCase().includes(q) || b.family.toLowerCase().includes(q);
    }).sort((a, b) => a.com.localeCompare(b.com))
  );
  let kaFamilyStats = $derived(
    (() => {
      const m = new Map<string, number>();
      for (const b of KA_BIRDS) m.set(b.family, (m.get(b.family) || 0) + 1);
      return [...m.entries()].sort((a, b) => b[1] - a[1]).slice(0, 12);
    })()
  );
  let kaMax = $derived(kaFamilyStats.reduce((mx, [, v]) => Math.max(mx, v), 1));
  let kaPaged = $derived(kaFiltered.slice((kaPage - 1) * KA_PAGE_SIZE, kaPage * KA_PAGE_SIZE));
  $effect(() => {
    void kaSearch;
    void kaFamily;
    kaPage = 1;
  });
</script>

<svelte:head><title>Bird Catalog — Pragati</title></svelte:head>

<div class="max-w-6xl mx-auto space-y-6">
  <div>
    <p class="text-xs font-semibold tracking-widest text-primary-600 uppercase">Karnataka checklist</p>
    <h1 class="text-2xl sm:text-3xl font-bold text-slate-900 mt-1">Birds of Karnataka</h1>
    <p class="text-sm text-slate-500 mt-1">
      {KA_BIRDS.length} species recorded for {KA_REGION} — mirrored once from the eBird API into this app. Nothing is sent back to eBird. Photographs by Wikimedia Commons contributors.
    </p>
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4">
    <h2 class="text-sm font-semibold text-slate-800 flex items-center gap-1.5 mb-3">
      <BarChart3 size={15} class="text-primary-600" /> Species per family — top 12
    </h2>
    <div class="space-y-2">
      {#each kaFamilyStats as [fam, n]}
        <div class="flex items-center gap-3">
          <div class="w-44 sm:w-56 shrink-0 text-xs text-slate-600 truncate" title={fam}>{fam}</div>
          <div class="flex-1 h-4 bg-slate-100 rounded overflow-hidden">
            <div
              class="h-full bg-gradient-to-r from-emerald-500 to-teal-500 rounded"
              style="width:{(n / kaMax) * 100}%"
            ></div>
          </div>
          <div class="w-8 shrink-0 text-xs font-semibold text-slate-600 text-right">{n}</div>
        </div>
      {/each}
    </div>
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
    <div class="flex flex-wrap gap-3 items-end">
      <div class="flex-1 min-w-48">
        <SearchFilter bind:value={kaSearch} placeholder="Search name, scientific name, family..." />
      </div>
      <div class="w-60">
        <Select bind:value={kaFamily} options={kaFamilyOptions} placeholder="All families" />
      </div>
    </div>
  </div>

  <p class="text-xs font-medium text-slate-500">{kaFiltered.length} species</p>

  {#if kaFiltered.length === 0}
    <div class="bg-white rounded-xl border border-slate-200">
      <EmptyState icon={Bird} title="No species match." hint="Try a different name or family." />
    </div>
  {:else}
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-4">
      {#each kaPaged as b (b.code)}
        <button
          onclick={() => {
            kaSelected = b;
            kaOpen = true;
          }}
          class="group text-left bg-white rounded-2xl border border-slate-200 overflow-hidden hover:border-primary-300 hover:shadow-md active:scale-[0.99] transition-all"
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
          <div class="p-3">
            <div class="text-sm font-semibold text-slate-800 leading-tight">{b.com}</div>
            <div class="text-[11px] text-slate-400 italic truncate">{b.sci}</div>
            <div class="text-[10px] text-slate-500 mt-0.5 truncate">{b.family}</div>
          </div>
        </button>
      {/each}
    </div>
    <div class="bg-white rounded-xl border border-slate-200 no-print">
      <Pagination page={kaPage} total={kaFiltered.length} pageSize={KA_PAGE_SIZE} onChange={(p) => (kaPage = p)} />
    </div>
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
        <p class="text-xs text-slate-400 italic">{kaSelected.sci}</p>
        <div class="flex flex-wrap gap-1.5 mt-2">
          <span class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-slate-100 text-slate-600">{kaSelected.family}</span>
          <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-primary-50 text-primary-700">{kaSelected.order}</span>
        </div>
      </div>

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
      <p class="text-[11px] text-slate-400 text-center -mt-2">Photo © {KA_PHOTOS[kaSelected.code]?.credit ?? 'Wikimedia Commons'}</p>
    </div>
    {#snippet footer()}
      <Button variant="ghost" onclick={() => (kaOpen = false)}>Close</Button>
    {/snippet}
  </Modal>
{/if}