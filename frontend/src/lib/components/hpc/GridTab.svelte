<script lang="ts">
  import { api, apiUrl } from "$lib/api/client.svelte";
  import type { Class, AcademicYear } from "$lib/types";
  import { onMount } from "svelte";
  import Select from "$lib/components/Select.svelte";
  import Button from "$lib/components/Button.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import { toast } from "$lib/stores/toast.svelte";

  let {
    cls = $bindable(""),
    term = $bindable("Term 1"),
    isAdmin,
  }: { cls?: string; term?: string; isAdmin: boolean } = $props();

  interface HPCGridRow {
    student_id: string;
    sats_number: string;
    name: string;
    roll_no: number;
    entry_id: string;
    status: string;
    has_pdf: boolean;
  }

  let classes = $state<Class[]>([]);
  let years = $state<AcademicYear[]>([]);
  let grid = $state<HPCGridRow[]>([]);
  let selectedYear = $state("");
  let loading = $state(false);
  let publishing = $state(false);
  let statusMsg = $state("");
  let statusType = $state<"success" | "error">("success");
  let summary = $state({
    total_students: 0,
    published_count: 0,
    draft_count: 0,
  });

  const terms = ["Term 1", "Term 2"];
  const termOptions = terms.map((t) => ({ id: t, name: t }));

  let gridSearch = $state("");
  let filteredGrid = $derived(
    grid.filter((r) => {
      if (!gridSearch.trim()) return true;
      const q = gridSearch.toLowerCase();
      return (
        r.name?.toLowerCase().includes(q) ||
        r.sats_number?.includes(q)
      );
    }),
  );

  onMount(async () => {
    const [cRes, yRes] = await Promise.all([
      api<Class[]>("GET", "/classes?limit=100"),
      api<AcademicYear[]>("GET", "/academic-years?limit=10"),
    ]);
    if (cRes.data) classes = cRes.data;
    if (yRes.data) {
      years = yRes.data;
      const current = yRes.data.find(
        (y: { is_current: boolean }) => y.is_current,
      );
      if (current) selectedYear = current.id;
    }
  });

  async function loadGrid() {
    if (!cls || !selectedYear) return;
    loading = true;
    statusMsg = "";
    const params = `class_id=${cls}&academic_year_id=${selectedYear}&term=${term}`;
    const [gRes, rRes] = await Promise.all([
      api<HPCGridRow[]>("GET", "/hpc/grid?" + params),
      api("GET", "/hpc/reports/class?" + params),
    ]);
    loading = false;
    if (gRes.data) grid = gRes.data;
    if (rRes.data) summary = rRes.data as any;
  }

  $effect(() => {
    if (cls && selectedYear) loadGrid();
  });

  async function migrateFromMarks() {
    if (!cls || !selectedYear) return;
    statusMsg = "Migrating marks to HPC entries...";
    statusType = "success";
    const res = await api<{ migrated: number }>(
      "POST",
      "/hpc/migrate-from-marks?" +
        `class_id=${cls}&academic_year_id=${selectedYear}&term=${term}`,
    );
    if (res.data) {
      statusMsg = `Migrated ${res.data.migrated} entries from marks.`;
      statusType = "success";
      loadGrid();
    } else if (res.error) {
      statusMsg = res.error.message;
      statusType = "error";
      toast(res.error.message, "error");
    }
  }

  async function publishAll() {
    const unpublished = grid.filter((r) => r.entry_id && r.status === "draft");
    if (unpublished.length === 0) {
      statusMsg = "No draft entries to publish.";
      statusType = "error";
      return;
    }
    if (
      !confirm(
        `Publish all ${unpublished.length} draft HPC entries for this class?`,
      )
    )
      return;
    publishing = true;
    statusMsg = `Publishing ${unpublished.length} entries...`;
    statusType = "success";
    let published = 0;
    let failed = 0;
    for (const row of unpublished) {
      const res = await api("POST", "/hpc/entries/publish", {
        entry_id: row.entry_id,
      });
      if (res.data) published++;
      else failed++;
    }
    publishing = false;
    if (failed > 0) {
      statusMsg = `Published ${published}, failed ${failed}.`;
      statusType = "error";
      toast(
        `${failed} entr${failed === 1 ? "y" : "ies"} failed to publish`,
        "error",
      );
    } else {
      statusMsg = `Published ${published} entries.`;
      statusType = "success";
    }
    loadGrid();
  }
</script>

<div class="space-y-4">
  <div class="flex justify-end gap-2">
    {#if isAdmin}
    <Button
      variant="secondary"
      onclick={migrateFromMarks}
      disabled={!cls || publishing}
    >
      Migrate from Marks
    </Button>
    {/if}
    <Button
      onclick={publishAll}
      disabled={!cls || publishing}
      loading={publishing}
    >
      {publishing ? "Publishing..." : "Publish All"}
    </Button>
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4">
    <div class="flex flex-wrap gap-3 items-end">
      <div>
        <label
          for="hpc-class"
          class="block text-xs font-medium text-slate-600 mb-1">Class</label
        >
        <Select
          id="hpc-class"
          bind:value={cls}
          options={classes}
          placeholder="Select Class"
        />
      </div>
      <div>
        <label
          for="hpc-term"
          class="block text-xs font-medium text-slate-600 mb-1">Term</label
        >
        <Select id="hpc-term" bind:value={term} options={termOptions} />
      </div>
    </div>
  </div>

  {#if statusMsg}
    <div
      class="text-sm px-4 py-2 rounded-lg border {statusType === 'error'
        ? 'bg-red-50 text-red-700 border-red-200'
        : 'bg-emerald-50 text-emerald-700 border-emerald-200'}"
    >
      {statusMsg}
    </div>
  {/if}

  {#if summary.total_students > 0}
    <div class="flex gap-4 text-sm">
      <span class="px-3 py-1 bg-slate-100 rounded-full"
        >Total: <strong>{summary.total_students}</strong></span
      >
      <span class="px-3 py-1 bg-green-100 text-green-800 rounded-full"
        >Published: <strong>{summary.published_count}</strong></span
      >
      <span class="px-3 py-1 bg-amber-100 text-amber-800 rounded-full"
        >Draft: <strong>{summary.draft_count}</strong></span
      >
    </div>
  {/if}

  <div
    class="bg-white rounded-xl border border-slate-200 overflow-hidden overflow-x-auto"
  >
    <div class="p-4 border-b border-slate-200">
      <SearchFilter
        bind:value={gridSearch}
        placeholder="Search by name or SATS..."
      />
    </div>
    <table class="w-full text-sm">
      <thead>
        <tr class="bg-slate-50 text-slate-600">
          <th class="text-left px-4 py-3 font-medium">#</th>
          <th class="text-left px-4 py-3 font-medium">SATS</th>
          <th class="text-left px-4 py-3 font-medium">Name</th>
          <th class="text-left px-4 py-3 font-medium">Status</th>
          <th class="text-left px-4 py-3 font-medium">Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          <tr
            ><td colspan="5" class="px-4 py-8 text-center text-slate-400"
              >Loading...</td
            ></tr
          >
        {:else if grid.length === 0}
          <tr
            ><td colspan="5" class="px-4"
              ><EmptyState
                title={cls
                  ? 'No HPC entries found. Click "Migrate from Marks" to create.'
                  : "Select a class above."}
              /></td
            ></tr
          >
        {:else if filteredGrid.length === 0}
          <tr
            ><td colspan="5" class="px-4"
              ><EmptyState
                title="No students match your search."
              /></td
            ></tr
          >
        {:else}
          {#each filteredGrid as row (row.student_id)}
            <tr class="border-t border-slate-100 hover:bg-slate-50">
              <td class="px-4 py-3">{row.roll_no}</td>
              <td class="px-4 py-3 font-mono text-xs">{row.sats_number}</td>
              <td class="px-4 py-3 font-medium">{row.name}</td>
              <td class="px-4 py-3">
                {#if !row.entry_id}
                  <span class="text-xs text-slate-400">—</span>
                {:else if row.status === "published"}
                  <span
                    class="px-2 py-0.5 text-xs bg-green-100 text-green-700 rounded-full"
                    >Published</span
                  >
                {:else}
                  <span
                    class="px-2 py-0.5 text-xs bg-amber-100 text-amber-700 rounded-full"
                    >Draft</span
                  >
                {/if}
              </td>
              <td class="px-4 py-3">
                <div class="flex gap-2">
                  <a
                    href="/hpc/entry/{row.student_id}?term={term}&year={selectedYear}"
                    class="text-xs text-primary-600 hover:text-primary-800 font-medium"
                  >
                    {row.entry_id ? "Edit" : "Create"}
                  </a>
                  {#if row.has_pdf}
                    <a
                      href={apiUrl(`/hpc/entries/${row.entry_id}/pdf`)}
                      target="_blank"
                      class="text-xs text-slate-500 hover:text-slate-700">PDF</a
                    >
                  {/if}
                </div>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
