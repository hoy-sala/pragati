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

<div class="mx-auto max-w-4xl space-y-6">
  <a href="/birds" class="bback">
    <ArrowLeft size={14} /> Back to all birds
  </a>

  {#if !bird}
    <div class="bpanel">
      <EmptyState icon={Bird} title="Species not found." hint="This code is not in the Karnataka checklist." />
    </div>
  {:else}
    <section class="bhero overflow-hidden">
      <div class="bphoto relative h-56 sm:h-72 overflow-hidden">
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
            <BirdShape shape={bird.shape} class="w-40 h-40 opacity-80" />
          </div>
        {/if}
      </div>
      <div class="p-5 sm:p-7">
        <p class="text-xs font-bold tracking-widest uppercase opacity-60">Karnataka checklist</p>
        <h1 class="btitle mt-1 text-2xl sm:text-3xl">{bird.com}</h1>
        {#if KA_KANNADA[bird.code]}
          <p class="mt-1 text-lg font-bold text-emerald-700">{KA_KANNADA[bird.code]}</p>
        {/if}
        <p class="mt-1 text-sm italic opacity-60">{bird.sci}</p>
        <div class="mt-3 flex flex-wrap gap-1.5">
          <span class="btag">{bird.family}</span>
          <span class="btag">{bird.order}</span>
          {#if detail?.iucn}
            <IucnBadge code={detail.iucn} />
          {/if}
        </div>
      </div>
      <div class="flex gap-1 overflow-x-auto px-5 sm:px-7 no-print" style="border-top: 2px solid var(--ink)">
        {#each tabs as t (t.id)}
          <button
            onclick={() => (tab = t.id)}
            class="btab {tab === t.id ? 'btab-on' : ''}"
          >
            {t.label}
          </button>
        {/each}
      </div>
    </section>

    {#if tab === "overview"}
      <section class="bpanel space-y-3 p-5 sm:p-6">
        {#if detail?.extract}
          <p class="text-[15px] leading-relaxed">{detail.extract}</p>
          <p class="text-[11px] font-semibold opacity-50">
            Description from <a href={KA_PHOTOS[bird.code]?.wiki ?? `https://en.wikipedia.org/wiki/Special:Search?search=${encodeURIComponent(bird.sci)}`} target="_blank" rel="noreferrer" class="underline">Wikipedia</a>
            (CC BY-SA) · taxonomy & checklist from eBird
          </p>
        {:else}
          <p class="text-sm font-semibold opacity-60">No description available yet for this species.</p>
        {/if}
      </section>
    {:else if tab === "similar"}
      <section class="bpanel p-5 sm:p-6">
        {#if similar.length === 0}
          <p class="text-sm font-semibold opacity-60">No close relatives in the Karnataka checklist.</p>
        {:else}
          <p class="mb-3 text-xs font-semibold opacity-60">Other {bird.family} recorded in Karnataka — same-shape species first.</p>
          <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
            {#each similar as s (s.code)}
              <a
                href="/birds/{s.code}"
                class="bcard group"
              >
                <div class="bphoto flex h-24 items-center justify-center overflow-hidden">
                  {#if KA_PHOTOS[s.code]}
                    <img src={KA_PHOTOS[s.code].photo} alt={s.com} loading="lazy" class="w-full h-full object-contain p-1 transition-transform group-hover:scale-105" />
                  {:else}
                    <BirdShape shape={s.shape} class="w-12 h-12 opacity-80" />
                  {/if}
                </div>
                <div class="p-2">
                  <div class="bname truncate text-xs leading-tight" title={s.com}>{s.com}</div>
                  {#if KA_KANNADA[s.code]}
                    <div class="truncate text-[11px] font-bold text-emerald-700">{KA_KANNADA[s.code]}</div>
                  {/if}
                </div>
              </a>
            {/each}
          </div>
        {/if}
      </section>
    {:else if tab === "facts"}
      <section class="bpanel p-5 sm:p-6">
        <dl class="grid grid-cols-1 gap-3 text-sm">
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 font-bold opacity-50">Common name</dt>
            <dd class="font-bold">{bird.com}</dd>
          </div>
          {#if KA_KANNADA[bird.code]}
            <div class="flex items-start gap-3">
              <dt class="w-32 shrink-0 font-bold opacity-50">Kannada name</dt>
              <dd class="text-emerald-700 font-bold">{KA_KANNADA[bird.code]} <span class="font-normal opacity-50">(eBird)</span></dd>
            </div>
          {/if}
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 font-bold opacity-50">Scientific name</dt>
            <dd class="italic">{bird.sci}</dd>
          </div>
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 font-bold opacity-50">Family</dt>
            <dd>{bird.family}</dd>
          </div>
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 font-bold opacity-50">Order</dt>
            <dd>{bird.order}</dd>
          </div>
          {#if KA_COLORS[bird.code]}
            <div class="flex items-start gap-3">
              <dt class="w-32 shrink-0 font-bold opacity-50">Colours</dt>
              <dd class="flex flex-wrap items-center gap-2">
                {#each KA_COLORS[bird.code] as cid}
                  {@const sw = KA_SWATCHES.find((s) => s.id === cid)}
                  {#if sw}
                    <span class="inline-flex items-center gap-1.5 text-xs font-semibold">
                      <span class="bsw-mini" style="background-color:{sw.hex}"></span>
                      {sw.label}
                    </span>
                  {/if}
                {/each}
                <span class="text-[11px] opacity-50">(approximate plumage)</span>
              </dd>
            </div>
          {/if}
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 font-bold opacity-50">eBird code</dt>
            <dd class="font-mono text-xs pt-0.5">{bird.code}</dd>
          </div>
          {#if detail?.iucn}
            <div class="flex items-start gap-3">
              <dt class="w-32 shrink-0 font-bold opacity-50">Conservation</dt>
              <dd><IucnBadge code={detail.iucn} /> <span class="text-[11px] opacity-50">IUCN Red List via Wikidata</span></dd>
            </div>
          {/if}
          <div class="flex items-start gap-3">
            <dt class="w-32 shrink-0 font-bold opacity-50">Recorded in</dt>
            <dd>Karnataka (IN-KA) — from the eBird state checklist</dd>
          </div>
        </dl>
      </section>
    {:else}
      <section class="bpanel space-y-3 p-5 sm:p-6">
        <a
          href="https://ebird.org/species/{bird.code}"
          target="_blank"
          rel="noreferrer"
          class="bbtn bbtn-ink w-full"
        >
          <ExternalLink size={16} /> Photos, calls & range map on eBird
        </a>
        <a
          href={KA_PHOTOS[bird.code]?.wiki ?? `https://en.wikipedia.org/wiki/Special:Search?search=${encodeURIComponent(bird.sci)}`}
          target="_blank"
          rel="noreferrer"
          class="bbtn w-full"
        >
          Read the Wikipedia species page <ExternalLink size={15} />
        </a>
        <p class="text-center text-[11px] font-semibold opacity-50">Photo © {KA_PHOTOS[bird.code]?.credit ?? 'Wikimedia Commons'}</p>
      </section>
    {/if}
  {/if}
</div>

<style>
  .bhero {
    background: var(--paper);
    border: 2.5px solid var(--ink);
    box-shadow: 6px 6px 0 var(--ink);
    border-radius: 18px;
    color: var(--ink);
  }
  .bpanel {
    background: var(--paper);
    border: 2px solid var(--ink);
    box-shadow: 4px 4px 0 var(--ink);
    border-radius: 16px;
    color: var(--ink);
  }
  .btitle {
    font-family: var(--font-display);
    font-weight: 800; letter-spacing: -0.02em; line-height: 1.05;
  }
  .bphoto { background: var(--cream); border-bottom: 2px solid var(--ink); }
  .bhero .bphoto { border-bottom: 2.5px solid var(--ink); }
  .btag {
    font-size: 0.7rem; font-weight: 700;
    padding: 0.2rem 0.55rem; border-radius: 999px;
    border: 1.5px solid var(--ink); background: var(--cream); color: var(--ink);
  }
  .btab {
    flex: none; padding: 0.8rem 1rem;
    font-size: 0.85rem; font-weight: 800; font-family: var(--font-display);
    color: var(--ink-soft); cursor: pointer;
    border-bottom: 3px solid transparent; margin-bottom: -2px;
  }
  .btab:hover { color: var(--ink); }
  .btab-on { color: var(--ink); border-bottom-color: var(--ink); }
  .bback {
    display: inline-flex; align-items: center; gap: 0.4rem;
    font-size: 0.78rem; font-weight: 800;
    padding: 0.45rem 0.8rem; border-radius: 10px;
    border: 2px solid var(--ink); background: var(--paper); color: var(--ink);
    box-shadow: 2px 2px 0 var(--ink);
    transition: box-shadow 120ms, transform 120ms;
  }
  .bback:hover { transform: translate(1px, 1px); box-shadow: 1px 1px 0 var(--ink); }
  .bback:active { transform: translate(2px, 2px); box-shadow: 0 0 0 var(--ink); }
  .bbtn {
    display: inline-flex; align-items: center; justify-content: center; gap: 0.5rem;
    font-family: var(--font-display); font-weight: 700; font-size: 0.9rem;
    padding: 0.65rem 1rem; border-radius: 12px;
    border: 2.5px solid var(--ink); background: var(--paper); color: var(--ink);
    box-shadow: 3px 3px 0 var(--ink); cursor: pointer;
    transition: box-shadow 120ms, transform 120ms, background 120ms;
  }
  .bbtn:hover { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--ink); background: var(--cream); }
  .bbtn:active { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
  .bbtn-ink { background: var(--ink); color: var(--paper); }
  .bbtn-ink:hover { background: var(--ink); }
  .bcard {
    background: var(--paper);
    border: 2px solid var(--ink);
    box-shadow: 3px 3px 0 var(--ink);
    border-radius: 14px; overflow: hidden;
    color: var(--ink);
    transition: box-shadow 120ms, transform 120ms;
  }
  .bcard:hover { transform: translate(1px, 1px); box-shadow: 2px 2px 0 var(--ink); }
  .bcard:active { transform: translate(3px, 3px); box-shadow: 0 0 0 var(--ink); }
  .bname { font-family: var(--font-display); font-weight: 700; }
  .bsw-mini {
    width: 1rem; height: 1rem; border-radius: 999px;
    border: 1.5px solid var(--ink); display: inline-block;
  }
  .bbtn:focus-visible, .bcard:focus-visible, .btab:focus-visible,
  .bback:focus-visible {
    outline: 3px solid var(--ink);
    outline-offset: 2px;
  }
</style>
