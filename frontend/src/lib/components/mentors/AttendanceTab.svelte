<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import {
    Check,
    X,
    Clock,
    Phone,
    LoaderCircle,
  } from "lucide-svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import { Users } from "lucide-svelte";
  import { toast } from "$lib/stores/toast.svelte";

  let { year }: { year: string } = $props();

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
  let updating = $state<Set<string>>(new Set());

  async function loadAttendance() {
    if (!year) return;
    loading = true;
    const res = await api<typeof roster>(
      "GET",
      `/mentors/attendance?date=${date}&academic_year_id=${year}`,
    );
    if (res.data) roster = res.data;
    loading = false;
  }

  $effect(() => {
    if (year) loadAttendance();
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
      toast("Parent contacted", "success");
    }
    updating = new Set([...updating].filter((id) => id !== studentId));
    loadAttendance();
  }

  function statusColor(s: string): string {
    if (s === "present") return "bg-emerald-500";
    if (s === "absent") return "bg-red-500";
    return "bg-amber-500";
  }
</script>

<div class="space-y-6">
  <div
    class="bg-white rounded-xl border border-slate-200 p-4 no-print flex flex-wrap gap-3"
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
  {:else if roster.length === 0}
    <div class="bg-white rounded-xl border border-slate-200">
      <EmptyState icon={Users} title="No students in your roster for this year." />
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
