<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import { onMount } from "svelte";
  import {
    Check,
    X,
    Clock,
    Phone,
    AlertTriangle,
    LoaderCircle,
  } from "lucide-svelte";
  import { toast } from "$lib/stores/toast.svelte";
  import PageTabs from "$lib/components/PageTabs.svelte";
  import { MENTOR_TABS } from "$lib/utils/tabs";
  import { getAuthState } from "$lib/stores/auth.svelte";
  import { effectiveRole } from "$lib/utils/roles";
  import type { AcademicYear } from "$lib/types";

  const auth = getAuthState();
  let role = $derived(effectiveRole(auth.currentUser));

  let years = $state<AcademicYear[]>([]);
  let selectedYear = $state("");
  let date = $state(new Date().toISOString().slice(0, 10));
  let roster = $state<
    {
      student_id: string;
      name: string;
      status: string;
      parent_contacted: boolean;
      remarks: string;
    }[]
  >([]);
  let loading = $state(true);
  let saving = $state(false);
  let updating = $state<Set<string>>(new Set());

  onMount(async () => {
    const res = await api<AcademicYear[]>("GET", "/academic-years");
    if (res.data) {
      years = res.data;
      const cur = res.data.find((y) => y.is_current);
      if (cur) selectedYear = cur.id;
    }
  });

  async function loadAttendance() {
    if (!selectedYear) return;
    loading = true;
    const res = await api<typeof roster>(
      "GET",
      `/mentors/attendance?date=${date}&academic_year_id=${selectedYear}`,
    );
    if (res.data) roster = res.data;
    loading = false;
  }

  $effect(() => {
    if (selectedYear) loadAttendance();
  });

  async function mark(studentId: string, status: string) {
    if (updating.has(studentId)) return;
    updating = new Set(updating).add(studentId);
    const res = await api("PUT", "/mentors/attendance", {
      student_id: studentId,
      date,
      status,
      remarks: "",
    });
    if (res.error) {
      toast(res.error.message, "error");
    } else {
      toast(`Marked ${status}`, "success");
    }
    updating = new Set([...updating].filter((id) => id !== studentId));
    loadAttendance();
  }

  async function contactParent(studentId: string) {
    if (updating.has(studentId)) return;
    updating = new Set(updating).add(studentId);
    const res = await api("POST", "/mentors/contact-parent", {
      student_id: studentId,
      date,
    });
    if (res.error) {
      toast(res.error.message, "error");
    } else {
      toast("Parent contact marked", "success");
    }
    updating = new Set([...updating].filter((id) => id !== studentId));
    loadAttendance();
  }

  function statusColor(s: string) {
    if (s === "present") return "bg-emerald-500";
    if (s === "absent") return "bg-red-500";
    return "bg-amber-500";
  }
</script>

<svelte:head><title>Daily Attendance - Pragati</title></svelte:head>

<div class="max-w-5xl mx-auto space-y-6">
  <div class="flex items-center gap-3">
    <div
      class="w-10 h-10 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-700 flex items-center justify-center"
    >
      <Check size={20} class="text-white" />
    </div>
    <div>
      <h1 class="text-2xl font-bold text-slate-900">Daily Attendance</h1>
      <p class="text-sm text-slate-500">{roster.length} students</p>
    </div>
  </div>
  <PageTabs tabs={MENTOR_TABS} role={role} />
  <div
    class="bg-white rounded-xl border border-slate-200 p-4 no-print flex gap-3"
  >
    <input
      type="date"
      bind:value={date}
      onchange={() => loadAttendance()}
      class="px-3 py-2 rounded-lg border border-slate-300 text-sm"
    />
  </div>
  {#if loading}<div
      class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400"
    >
      Loading...
    </div>
  {:else}
    <div
      class="bg-white rounded-xl border border-slate-200 divide-y divide-slate-100"
    >
      {#each roster as s}
        <div class="px-4 py-3 flex items-center gap-4">
          <div class="w-2 h-2 rounded-full {statusColor(s.status)}"></div>
          <span class="text-sm font-medium text-slate-700 flex-1">{s.name}</span
          >
          <div class="flex gap-2">
            <button
              disabled={updating.has(s.student_id)}
              onclick={() => mark(s.student_id, "present")}
              class="p-2 rounded-lg disabled:opacity-40 disabled:cursor-not-allowed {s.status ===
              'present'
                ? 'bg-emerald-100 text-emerald-600'
                : 'bg-slate-50 text-slate-400 hover:bg-emerald-50'}"
              ><Check size={16} /></button
            >
            <button
              disabled={updating.has(s.student_id)}
              onclick={() => mark(s.student_id, "absent")}
              class="p-2 rounded-lg disabled:opacity-40 disabled:cursor-not-allowed {s.status ===
              'absent'
                ? 'bg-red-100 text-red-600'
                : 'bg-slate-50 text-slate-400 hover:bg-red-50'}"
              ><X size={16} /></button
            >
            <button
              disabled={updating.has(s.student_id)}
              onclick={() => mark(s.student_id, "late")}
              class="p-2 rounded-lg disabled:opacity-40 disabled:cursor-not-allowed {s.status ===
              'late'
                ? 'bg-amber-100 text-amber-600'
                : 'bg-slate-50 text-slate-400 hover:bg-amber-50'}"
              ><Clock size={16} /></button
            >
          </div>
          {#if s.status === "absent" && !s.parent_contacted}
            <button
              disabled={updating.has(s.student_id)}
              onclick={() => contactParent(s.student_id)}
              class="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-blue-50 text-blue-600 text-xs font-medium hover:bg-blue-100 disabled:opacity-50 disabled:cursor-not-allowed"
              >{#if updating.has(s.student_id)}<LoaderCircle
                  size={12}
                  class="animate-spin"
                />{:else}<Phone size={12} />{/if} Call</button
            >
          {:else if s.parent_contacted}
            <span class="text-xs text-emerald-600 flex items-center gap-1"
              ><Phone size={12} /> Called</span
            >
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>
