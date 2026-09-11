<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { Users } from "lucide-svelte";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import PageTabs from "$lib/components/PageTabs.svelte";
  import Select from "$lib/components/Select.svelte";
  import { MENTOR_TABS } from "$lib/utils/tabs";
  import { getAuthState } from "$lib/stores/auth.svelte";
  import { effectiveRole, hasRole } from "$lib/utils/roles";
  import type { AcademicYear } from "$lib/types";
  import AssignTab from "$lib/components/mentors/AssignTab.svelte";
  import RosterTab from "$lib/components/mentors/RosterTab.svelte";
  import AttendanceTab from "$lib/components/mentors/AttendanceTab.svelte";
  import LogsTab from "$lib/components/mentors/LogsTab.svelte";
  import DashboardTab from "$lib/components/mentors/DashboardTab.svelte";

  const auth = getAuthState();
  let role = $derived(effectiveRole(auth.currentUser));
  let isAdmin = $derived(hasRole(auth.currentUser, "admin"));
  let canReview = $derived(hasRole(auth.currentUser, "admin", "principal"));

  const VALID_TABS = ["assignments", "roster", "attendance", "logs", "dashboard"];
  let tabParam = $derived($page.url.searchParams.get("tab") ?? "");
  let tab = $derived(VALID_TABS.includes(tabParam) ? tabParam : "assignments");
  // Don't downgrade the dashboard tab before auth finishes loading.
  let activeTab = $derived(
    tab === "dashboard" && !auth.isLoading && !canReview ? "assignments" : tab,
  );

  let years = $state<AcademicYear[]>([]);
  let selectedYear = $state("");

  onMount(async () => {
    const yr = await api<AcademicYear[]>("GET", "/academic-years");
    if (yr.data) {
      years = yr.data;
      const cur = yr.data.find((y) => y.is_current);
      if (cur) selectedYear = cur.id;
    }
  });
</script>

<svelte:head><title>Mentors — Pragati</title></svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
  <PageHeader
    title="Mentors"
    subtitle="Assignments, roster, attendance, logs and oversight"
  >
    {#snippet icon()}
      <div
        class="w-10 h-10 rounded-xl bg-gradient-to-br from-teal-500 to-emerald-700 flex items-center justify-center shrink-0"
      >
        <Users size={20} class="text-white" />
      </div>
    {/snippet}
  </PageHeader>

  <PageTabs tabs={MENTOR_TABS} role={role} />

  <div class="bg-white rounded-xl border border-slate-200 p-4 no-print w-56">
    <Select
      bind:value={selectedYear}
      options={years}
      placeholder="Select year"
    />
  </div>

  {#if activeTab === "assignments"}
    <AssignTab year={selectedYear} isAdmin={isAdmin} />
  {:else if activeTab === "roster"}
    <RosterTab year={selectedYear} />
  {:else if activeTab === "attendance"}
    <AttendanceTab year={selectedYear} />
  {:else if activeTab === "logs"}
    <LogsTab year={selectedYear} />
  {:else if activeTab === "dashboard"}
    <DashboardTab year={selectedYear} canReview={canReview} />
  {/if}
</div>
