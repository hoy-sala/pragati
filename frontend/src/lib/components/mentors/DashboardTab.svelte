<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import {
    AlertTriangle,
    Users,
    FileText,
    Phone,
    CheckCircle,
  } from "lucide-svelte";
  import Button from "$lib/components/Button.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import { toast } from "$lib/stores/toast.svelte";

  let { year, canReview }: { year: string; canReview: boolean } = $props();

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
  let reviewOpen = $state(false);
  let reviewNotes = $state("");
  let reviewSaving = $state(false);

  async function load() {
    if (!year) return;
    loading = true;
    const [aRes, sRes] = await Promise.all([
      api<typeof alerts>("GET", `/mentors/principal/alerts`),
      api<typeof summary>(
        "GET",
        `/mentors/principal/summary?academic_year_id=${year}`,
      ),
    ]);
    if (aRes.data) alerts = aRes.data;
    if (sRes.data) summary = sRes.data;
    loading = false;
  }

  $effect(() => {
    if (year) load();
  });

  function openReview(log: AlertLog) {
    reviewingLog = log;
    reviewNotes = "";
    reviewOpen = true;
  }

  $effect(() => {
    if (!reviewOpen) reviewingLog = null;
  });

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
      reviewOpen = false;
      await load();
    }
    reviewSaving = false;
  }
</script>

<div class="space-y-6">
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
               {#if canReview}<Button
                 size="sm"
                 variant="secondary"
                 onclick={() => openReview(a)}
               >Review</Button>{/if}
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
              <tr>
                <td colspan="6" class="px-4">
                  <EmptyState title="No summary data for this year." />
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}

  {#if reviewingLog}
    <Modal
      bind:open={reviewOpen}
      title="Review Log"
      maxWidth="max-w-md"
    >
      <div class="space-y-4">
        <p class="text-xs text-slate-500">
          {reviewingLog.student_name} · {reviewingLog.severity} · {reviewingLog.category}
          · {reviewingLog.log_date}
        </p>
        <p class="text-sm text-slate-600">{reviewingLog.description}</p>
        <textarea
          bind:value={reviewNotes}
          rows="4"
          placeholder="Principal notes..."
          class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
        ></textarea>
      </div>
      {#snippet footer()}
        <Button
          variant="secondary"
          onclick={() => {
            if (!reviewSaving) reviewOpen = false;
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
      {/snippet}
    </Modal>
  {/if}
</div>
