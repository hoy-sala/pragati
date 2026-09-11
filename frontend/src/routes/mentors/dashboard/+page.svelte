<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import { onMount } from "svelte";
  import {
    AlertTriangle,
    Users,
    FileText,
    Phone,
    CheckCircle,
  } from "lucide-svelte";
  import { toast } from "$lib/stores/toast.svelte";
  import Button from "$lib/components/Button.svelte";
  import Select from "$lib/components/Select.svelte";
  import type { AcademicYear } from "$lib/types";

  type AlertLog = {
    id: string;
    student_id: string;
    student_name: string;
    log_date: string;
    category: string;
    severity: string;
    description: string;
    action_taken: string;
    parent_informed: boolean;
    mentor_name: string;
  };

  let years = $state<AcademicYear[]>([]);
  let selectedYear = $state("");
  let alerts = $state<AlertLog[]>([]);
  let summary = $state<
    {
      mentor_id: string;
      mentor_name: string;
      student_count: number;
      avg_attendance_pct: number;
      log_count: number;
      urgent_count: number;
      parent_contacts: number;
    }[]
  >([]);
  let loading = $state(true);
  let reviewingLog = $state<AlertLog | null>(null);
  let reviewNotes = $state("");
  let reviewSaving = $state(false);

  onMount(async () => {
    const res = await api<AcademicYear[]>("GET", "/academic-years");
    if (res.data) {
      years = res.data;
      const cur = res.data.find((y) => y.is_current);
      if (cur) selectedYear = cur.id;
    }
  });

  async function load() {
    if (!selectedYear) return;
    loading = true;
    const [aRes, sRes] = await Promise.all([
      api<typeof alerts>("GET", `/mentors/principal/alerts`),
      api<typeof summary>(
        "GET",
        `/mentors/principal/summary?academic_year_id=${selectedYear}`,
      ),
    ]);
    if (aRes.data) alerts = aRes.data;
    if (sRes.data) summary = sRes.data;
    loading = false;
  }

  $effect(() => {
    if (selectedYear) load();
  });

  function openReview(log: AlertLog) {
    reviewingLog = log;
    reviewNotes = "";
  }

  async function saveReview() {
    if (!reviewingLog) return;
    if (reviewSaving) return;
    reviewSaving = true;
    const res = await api("PUT", `/mentors/logs/${reviewingLog.id}/review`, {
      principal_notes: reviewNotes,
    });
    if (res.error) {
      toast(res.error.message, "error");
    } else {
      toast("Review saved", "success");
      reviewingLog = null;
      await load();
    }
    reviewSaving = false;
  }
</script>

<svelte:head><title>Mentor Dashboard - Pragati</title></svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
  <div class="flex items-center gap-3">
    <div
      class="w-10 h-10 rounded-xl bg-gradient-to-br from-red-500 to-rose-700 flex items-center justify-center"
    >
      <AlertTriangle size={20} class="text-white" />
    </div>
    <div>
      <h1 class="text-2xl font-bold text-slate-900">Mentor Dashboard</h1>
      <p class="text-sm text-slate-500">Principal oversight</p>
    </div>
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4 no-print w-56">
    <Select
      bind:value={selectedYear}
      options={years}
      placeholder="Select year"
    />
  </div>

  {#if loading}<div
      class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400"
    >
      Loading...
    </div>
  {:else}
    {#if alerts.length > 0}
      <div class="bg-red-50 border border-red-200 rounded-xl p-5">
        <h2
          class="text-base font-semibold text-red-700 mb-3 flex items-center gap-2"
        >
          <AlertTriangle size={18} /> Urgent Alerts ({alerts.length})
        </h2>
        <div class="space-y-2">
          {#each alerts as a}
            <div
              class="bg-white rounded-lg p-3 border border-red-100 flex items-start justify-between"
            >
              <div>
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-sm font-medium text-slate-800"
                    >{a.student_name}</span
                  >
                  <span
                    class="text-xs px-2 py-0.5 rounded bg-red-100 text-red-600"
                    >{a.severity}</span
                  >
                  <span class="text-xs text-slate-400">{a.category}</span>
                </div>
                <p class="text-sm text-slate-600">{a.description}</p>
                <p class="text-xs text-slate-400 mt-1">
                  By {a.mentor_name} . {a.log_date}
                </p>
              </div>
               <Button
                 size="sm"
                 variant="secondary"
                 onclick={() => openReview(a)}
               >Review</Button>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <div class="bg-white rounded-xl border border-slate-200 p-5">
      <h2 class="text-base font-semibold text-slate-900 mb-4">
        Monthly Summary
      </h2>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200">
              <th class="px-4 py-2 text-left font-semibold text-slate-600"
                >Mentor</th
              >
              <th class="px-4 py-2 text-center font-semibold text-slate-600"
                >Students</th
              >
              <th class="px-4 py-2 text-center font-semibold text-slate-600"
                >Avg Attendance</th
              >
              <th class="px-4 py-2 text-center font-semibold text-slate-600"
                >Logs</th
              >
              <th class="px-4 py-2 text-center font-semibold text-slate-600"
                >Urgent</th
              >
              <th class="px-4 py-2 text-center font-semibold text-slate-600"
                >Parent Contacts</th
              >
            </tr>
          </thead>
          <tbody>
            {#each summary as m}
              <tr class="border-b border-slate-100">
                <td class="px-4 py-2 font-medium text-slate-700"
                  >{m.mentor_name}</td
                >
                <td class="px-4 py-2 text-center text-slate-600"
                  >{m.student_count}</td
                >
                <td class="px-4 py-2 text-center"
                  ><span
                    class="text-xs px-2 py-0.5 rounded {m.avg_attendance_pct >=
                    80
                      ? 'bg-emerald-100 text-emerald-600'
                      : 'bg-amber-100 text-amber-600'}"
                    >{m.avg_attendance_pct.toFixed(1)}%</span
                  ></td
                >
                <td class="px-4 py-2 text-center text-slate-600"
                  >{m.log_count}</td
                >
                <td class="px-4 py-2 text-center"
                  ><span
                    class={m.urgent_count > 0
                      ? "text-red-600 font-bold"
                      : "text-slate-400"}>{m.urgent_count}</span
                  ></td
                >
                <td class="px-4 py-2 text-center text-slate-600"
                  >{m.parent_contacts}</td
                >
              </tr>
            {:else}
              <tr
                ><td colspan="6" class="px-4 py-8 text-center text-slate-400"
                  >No data</td
                ></tr
              >
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}

  {#if reviewingLog}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4"
      role="dialog"
      aria-modal="true"
      aria-label="Review log"
    >
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-5 space-y-4">
        <div>
          <h3 class="text-base font-semibold text-slate-900">Review Log</h3>
          <p class="text-xs text-slate-500 mt-0.5">
            {reviewingLog.student_name} · {reviewingLog.severity} · {reviewingLog.category}
            · {reviewingLog.log_date}
          </p>
        </div>
        <p class="text-sm text-slate-600">{reviewingLog.description}</p>
        <textarea
          bind:value={reviewNotes}
          rows="4"
          placeholder="Principal notes..."
          class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
        ></textarea>
        <div class="flex justify-end gap-2">
          <Button
            variant="secondary"
            onclick={() => {
              if (!reviewSaving) reviewingLog = null;
            }}
            disabled={reviewSaving}
            >Cancel</Button
          >
          <Button
            onclick={saveReview}
            disabled={reviewSaving}
            loading={reviewSaving}
            >Save Review</Button
          >
        </div>
      </div>
    </div>
  {/if}
</div>
