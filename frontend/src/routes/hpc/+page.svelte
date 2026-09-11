<script lang="ts">
  import { page } from "$app/stores";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import PageTabs from "$lib/components/PageTabs.svelte";
  import { HPC_TABS } from "$lib/utils/tabs";
  import { getAuthState } from "$lib/stores/auth.svelte";
  import { effectiveRole, hasRole } from "$lib/utils/roles";
  import GridTab from "$lib/components/hpc/GridTab.svelte";
  import AssessTab from "$lib/components/hpc/AssessTab.svelte";
  import ConfigTab from "$lib/components/hpc/ConfigTab.svelte";
  import ImportTab from "$lib/components/hpc/ImportTab.svelte";

  const auth = getAuthState();
  let role = $derived(effectiveRole(auth.currentUser));
  let isAdmin = $derived(hasRole(auth.currentUser, "admin"));

  const VALID_TABS = ["grid", "assess", "import", "config"];
  let tabParam = $derived($page.url.searchParams.get("tab") ?? "");
  let tab = $derived(VALID_TABS.includes(tabParam) ? tabParam : "grid");
  // Don't downgrade admin tabs before auth finishes loading.
  let activeTab = $derived(
    (tab === "import" || tab === "config") && !auth.isLoading && !isAdmin
      ? "grid"
      : tab,
  );

  let selectedClass = $state("");
  let selectedTerm = $state("Term 1");
</script>

<svelte:head><title>HPC — Pragati</title></svelte:head>

<div class="max-w-7xl mx-auto space-y-4">
  <PageHeader
    title="Holistic Progress Card (HPC)"
    subtitle="Assess learning outcomes per student"
  />

  <PageTabs tabs={HPC_TABS} role={role} />

  {#if activeTab === "grid"}
    <GridTab bind:cls={selectedClass} bind:term={selectedTerm} isAdmin={isAdmin} />
  {:else if activeTab === "assess"}
    <AssessTab bind:cls={selectedClass} bind:term={selectedTerm} />
  {:else if activeTab === "import"}
    <ImportTab bind:cls={selectedClass} isAdmin={isAdmin} />
  {:else if activeTab === "config"}
    <ConfigTab bind:cls={selectedClass} isAdmin={isAdmin} />
  {/if}
</div>
