<script lang="ts">
  import { page } from "$app/stores";
  import { KA_BIRDS } from "$lib/data/kaBirds";
  import { KA_PHOTOS } from "$lib/data/kaPhotos";
  import { KA_KANNADA } from "$lib/data/kaKannada";
  import { KA_DETAILS } from "$lib/data/kaDetails";
  import { KA_COLORS, KA_SWATCHES } from "$lib/data/kaColors";
  import BirdShape from "$lib/components/BirdShape.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import IucnBadge from "$lib/components/IucnBadge.svelte";
  import { Bird, ExternalLink, ArrowLeft } from "lucide-svelte";

  const code = $page.params.code;
  const bird = KA_BIRDS.find((b) => b.code === code);
  const detail = (bird && KA_DETAILS[bird.code]) || null;
  const similar = bird
    ? KA_BIRDS.filter((b) => b.code !== bird.code && b.family === bird.family)
        .sort((a, b) => (a.shape === bird.shape ? 0 : 1) - (b.shape === bird.shape ? 0 : 1))
        .slice(0, 12)
    : [];

  type Tab = "overview" | "similar" | "facts" | "links";
  let tab = $state<Tab>("overview");
  const tabs: { id: Tab; label: string }[] = [
    { id: "overview", label: "Overview" },
    { id: "similar", label: `Similar (${similar.length})` },
    { id: "facts", label: "Facts" },
    { id: "links", label: "Links" }
  ];
</script>

<svelte:head><title>{bird ? `${bird.com} — Bird Catalog` : "Not found — Bird Catalog"}</title></svelte:head>

<div class="max-w-4xl mx-auto space-y-6">
  <a href="/birds" class="inline-flex items-center gap-1.5 text-xs font-semibold text-primary-700 hover:text-primary-800">
    <ArrowLeft size={14} /> Back to all birds
  </a>

  {#if !bird}
    <div class="bg-white rounded-2xl border border-slate-200 shadow-sm">
      <EmptyState icon={Bird} title="Species not found." hint="This code is not in the Karnataka checklist." />
    </div>
  {:else}
    <section class="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
      <div class="relative h-56 sm:h-72 bg-gradient-to-br from-emerald-50 to-teal-100 overflow-hidden">
        {#if KA_PHOTOS[bird.code]}
          <img
            src={KA_PHOTOS[bird.code].photo}
            alt={bird.com}
            class="absolute inset-0 w-full h-full object-contain p-3"
          />
          <span class="absolute bottom-2 right-2 text-[10px] px-1.5 py-0.5 rounded bg-black/60 text-white truncate max-w-[80%]">
            © {KA_PHOTOS[bird.code].credit}
          </span>
        {:else}
          <div class="absolute inset-0 flex items-center justify-center">
            <BirdShape shape={bird.shape} class="w-40 h-40 text-slate-700/80" />
          </div>
        {/if}
      </div>
      <div class="p-5 sm:p-7">
        <p class="text-xs font-semibold tracking-widest text-primary-600 uppercase">Karnataka checklist</p>
        <h1 class="text-2xl sm:text-3xl font-bold text-slate-900 mt-1">{bird.com}</h1>
        {#if KA_KANNADA[bird.code]}
          <p class="text-lg font-semibold text-emerald-700 mt-1">{KA_KANNADA[bird.code]}</p>
        {/if}
        <p class="text-sm text-slate-400 italic mt-1">{bird.sci}</p>
        <div class="flex flex-wrap gap-1.5 mt-3">
          <span class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-slate-100 text-slate-600">{bird.family}</span>
          <span class="text-[11px] font-medium px-2 py-0.5 rounded-full bg-primary-50 text-primary-700">{bird.order}</span>
          {#if detail?.iucn}
            <IucnBadge code={detail.iucn} />
          {/if}
        </div>
      </div>
      <div class="flex gap-1 px-5 sm:px-7 border-t border-slate-100 overflow-x-auto no-print">
        {#each tabs as t (t.id)}
          <button
            onclick={() => (tab = t.id)}
            class="shrink-0 px-4 py-3 text-sm font-semibold border-b-2 -mb-px transition-colors {tab === t.id
              ? 'border-primary-600 text-primary-700'
              : 'border-transparent text-slate-500 hover:text-slate-800'}"
          >
            {t.label}
          </button>
        {/each}
      </div>
    </section>

    {#if tab === "overview"}
      <section class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5 sm:p-6 space-y-3">
        {#if detail?.extract}
          <p class="text-[15px] text-slate-700 leading-relaxed">{detail.extract}</p>
          <p class="text-[11px] text-slate-400">
            Description from <a href={KA_PHOTOS[bird.code]?.wiki ?? `https://en.wikipedia.org/wiki/Special:Search?search=${encodeURIComponent(bird.sci)}`} target="_blank" rel="noreferrer" class="underline hover:text-slate-600">Wikipedia</a>
            (CC BY-SA) · taxonomy & checklist from eBird
          </p>
        {:else}
          <p class="text-sm text-slate-500">No description available yet for this species.</p>
        {/if}
      </section>
    {:else if tab === "similar"}
      <section class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5 sm:p-6">
        {#if similar.length === 0}
          <p class="text-sm text-slate-500">No close relatives in the Karnataka checklist.</p>
        {:else}
          <p class="text-xs text-slate-500 mb-3">Other {bird.family} recorded in Karnataka — same-shape species first.</p>
          <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
            {#each similar as s (s.code)}
              <a
                href="/birds/{s.code}"
                class="group text-left bg-white rounded-2xl border border-slate-200 overflow-hidden hover:border-primary-300 hover:shadow-md transition-all"
              >
                <div class="h-24 flex items-center justify-center bg-gradient-to-br from-emerald-50 to-teal-100 overflow-hidden">
                  {#if KA_PHOTOS[s.code]}
                    <img src={KA_PHOTOS[s.code].photo} alt={s.com} loading="lazy" class="w-full h-full object-contain p-1 transition-transform group-hover:scale-105" />
                  {:else}
                    <BirdShape shape={s.shape} class="w-12 h-12 text-slate-700/80" />
                  {/if}
                </div>
                <div class="p-2">
                  <div class="text-xs font-semibold text-slate-800 leading-tight truncate" title={s.com}>{s.com}</div>
                  {#if KA_KANNADA[s.code]}
                    <div class="text-[11px] font-medium text-emerald-700 truncate">{KA_KANNADA[s.code]}</div>
                  {/if}
                </div>
              </a>
            {/each}
          </div>
        {/if}
      </section>
    {:else if tab === "facts"}
      <section class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5 sm:p-6">
        <dl class="grid grid-cols-1 gap-3 text-sm">
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 text-slate-400 font-medium">Common name</dt>
            <dd class="text-slate-800 font-semibold">{bird.com}</dd>
          </div>
          {#if KA_KANNADA[bird.code]}
            <div class="flex items-start gap-3">
              <dt class="w-32 shrink-0 text-slate-400 font-medium">Kannada name</dt>
              <dd class="text-emerald-700 font-semibold">{KA_KANNADA[bird.code]} <span class="text-slate-400 font-normal">(eBird)</span></dd>
            </div>
          {/if}
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 text-slate-400 font-medium">Scientific name</dt>
            <dd class="text-slate-700 italic">{bird.sci}</dd>
          </div>
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 text-slate-400 font-medium">Family</dt>
            <dd class="text-slate-700">{bird.family}</dd>
          </div>
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 text-slate-400 font-medium">Order</dt>
            <dd class="text-slate-700">{bird.order}</dd>
          </div>
          {#if KA_COLORS[bird.code]}
            <div class="flex items-start gap-3">
              <dt class="w-32 shrink-0 text-slate-400 font-medium">Colours</dt>
              <dd class="flex flex-wrap items-center gap-2">
                {#each KA_COLORS[bird.code] as cid}
                  {@const sw = KA_SWATCHES.find((s) => s.id === cid)}
                  {#if sw}
                    <span class="inline-flex items-center gap-1.5 text-xs text-slate-600">
                      <span class="w-4 h-4 rounded-full border border-slate-300" style="background-color:{sw.hex}"></span>
                      {sw.label}
                    </span>
                  {/if}
                {/each}
                <span class="text-[11px] text-slate-400">(approximate plumage)</span>
              </dd>
            </div>
          {/if}
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 text-slate-400 font-medium">eBird code</dt>
            <dd class="text-slate-700 font-mono text-xs pt-0.5">{bird.code}</dd>
          </div>
          {#if detail?.iucn}
            <div class="flex items-start gap-3">
              <dt class="w-32 shrink-0 text-slate-400 font-medium">Conservation</dt>
              <dd><IucnBadge code={detail.iucn} /> <span class="text-[11px] text-slate-400">IUCN Red List via Wikidata</span></dd>
            </div>
          {/if}
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 text-slate-400 font-medium">Recorded in</dt>
            <dd class="text-slate-700">Karnataka (IN-KA) — from the eBird state checklist</dd>
          </div>
        </dl>
      </section>
    {:else}
      <section class="bg-white rounded-2xl border border-slate-200 shadow-sm p-5 sm:p-6 space-y-3">
        <a
          href="https://ebird.org/species/{bird.code}"
          target="_blank"
          rel="noreferrer"
          class="flex items-center justify-center gap-2 w-full px-3.5 py-2.5 text-sm font-medium rounded-lg bg-slate-900 text-white hover:bg-slate-700 transition-colors"
        >
          <ExternalLink size={16} /> Photos, calls & range map on eBird
        </a>
        <a
          href={KA_PHOTOS[bird.code]?.wiki ?? `https://en.wikipedia.org/wiki/Special:Search?search=${encodeURIComponent(bird.sci)}`}
          target="_blank"
          rel="noreferrer"
          class="flex items-center justify-center gap-2 w-full px-3.5 py-2.5 text-sm font-medium rounded-lg border border-slate-300 text-slate-700 hover:bg-slate-50 transition-colors"
        >
          Read the Wikipedia species page <ExternalLink size={15} />
        </a>
        <p class="text-[11px] text-slate-400 text-center">Photo © {KA_PHOTOS[bird.code]?.credit ?? 'Wikimedia Commons'}</p>
      </section>
    {/if}
  {/if}
</div>
