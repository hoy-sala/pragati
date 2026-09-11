<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import {
    Plus,
    Heart,
    Brain,
    ShieldAlert,
    BookOpen,
    AlertTriangle,
    ClipboardList,
  } from "lucide-svelte";
  import Button from "$lib/components/Button.svelte";
  import Select from "$lib/components/Select.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import { toast } from "$lib/stores/toast.svelte";

  let { year }: { year: string } = $props();

  let logs = $state<
    {
      id: string;
      student_id: string;
      student_name: string;
      log_date: string;
      category: string;
      severity: string;
      description: string;
      action_taken: string;
      parent_informed: boolean;
      reviewed_by_principal: boolean;
    }[]
  >([]);
  let loading = $state(true);
  let showForm = $state(false);
  let saving = $state(false);
  let students = $state<{ id: string; name: string }[]>([]);
  let form = $state({
    student_id: "",
    category: "health",
    severity: "low",
    description: "",
    action_taken: "",
    parent_informed: false,
  });

  const categories = [
    { id: "health", name: "Health", icon: Heart },
    { id: "behavior", name: "Behavior", icon: Brain },
    { id: "grievance", name: "Grievance", icon: ShieldAlert },
    { id: "academic", name: "Academic", icon: BookOpen },
  ];
  const severities = [
    { id: "low", name: "Low", cls: "bg-emerald-100 text-emerald-700" },
    { id: "medium", name: "Medium", cls: "bg-blue-100 text-blue-700" },
    { id: "high", name: "High", cls: "bg-amber-100 text-amber-700" },
    { id: "urgent", name: "Urgent", cls: "bg-red-100 text-red-700" },
  ];

  async function loadStudents() {
    const roster = await api<{ id: string; name: string }[]>(
      "GET",
      "/mentors/roster?academic_year_id=current",
    );
    if (roster.data) students = roster.data;
  }

  async function loadLogs() {
    if (!year) return;
    loading = true;
    const res = await api<typeof logs>(
      "GET",
      `/mentors/logs?academic_year_id=${year}`,
    );
    if (res.data) logs = res.data;
    loading = false;
  }

  $effect(() => {
    loadStudents();
    if (year) loadLogs();
  });

  async function submitLog() {
    if (!form.student_id || !form.description) return;
    if (saving) return;
    saving = true;
    const res = await api("POST", "/mentors/logs", form);
    if (res.error) {
      toast(res.error.message, "error");
    } else {
      toast("Log submitted", "success");
      showForm = false;
      form = {
        student_id: "",
        category: "health",
        severity: "low",
        description: "",
        action_taken: "",
        parent_informed: false,
      };
      loadLogs();
    }
    saving = false;
  }

  function severityCls(sev: string) {
    return (
      severities.find((s) => s.id === sev)?.cls || "bg-slate-100 text-slate-600"
    );
  }
  function catName(cat: string) {
    return categories.find((c) => c.id === cat)?.name || cat;
  }
</script>

<div class="space-y-6">
  <div class="flex justify-end">
    <Button icon={Plus} onclick={() => (showForm = true)}>New Log</Button>
  </div>

  <Modal bind:open={showForm} title="New Log" maxWidth="max-w-2xl">
    <div class="space-y-4">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <Select
            id="log-student"
            label="Student *"
            bind:value={form.student_id}
            options={students}
            placeholder="Select student"
            searchable
          />
        </div>
        <div>
          <Select
            id="log-cat"
            label="Category"
            bind:value={form.category}
            options={categories}
          />
        </div>
        <div>
          <Select
            id="log-sev"
            label="Severity"
            bind:value={form.severity}
            options={severities}
          />
        </div>
        <div class="flex items-end pb-1">
          <label class="flex items-center gap-2 text-sm"
            ><input
              type="checkbox"
              bind:checked={form.parent_informed}
              class="rounded"
            /> Parent informed</label
          >
        </div>
      </div>
      <div>
        <label for="log-desc" class="block text-xs font-medium text-slate-600 mb-1"
          >Description *</label
        >
        <textarea
          id="log-desc"
          bind:value={form.description}
          rows="3"
          class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm"
          placeholder="Describe the issue..."
        ></textarea>
      </div>
      <div>
        <label for="log-action" class="block text-xs font-medium text-slate-600 mb-1"
          >Action Taken</label
        >
        <textarea
          id="log-action"
          bind:value={form.action_taken}
          rows="2"
          class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm"
          placeholder="Steps taken..."
        ></textarea>
      </div>
    </div>
    {#snippet footer()}
      <Button variant="ghost" onclick={() => (showForm = false)}>Cancel</Button>
      <Button onclick={submitLog} loading={saving} disabled={saving}
        >Submit Log</Button
      >
    {/snippet}
  </Modal>

  {#if loading}<div
      class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400"
    >
      Loading...
    </div>
  {:else}
    <div class="space-y-3">
      {#each logs as l}
        <div class="bg-white rounded-xl border border-slate-200 p-4">
          <div class="flex items-start justify-between mb-2">
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-slate-800"
                >{l.student_name}</span
              >
              <span
                class="text-xs px-2 py-0.5 rounded {severityCls(l.severity)}"
                >{l.severity}</span
              >
              <span class="text-xs text-slate-400">{catName(l.category)}</span>
            </div>
            <span class="text-xs text-slate-400">{l.log_date}</span>
          </div>
          <p class="text-sm text-slate-600 mb-2">{l.description}</p>
          {#if l.action_taken}<p class="text-xs text-slate-500">
              <strong>Action:</strong>
              {l.action_taken}
            </p>{/if}
          <div class="flex items-center gap-3 mt-2">
            {#if l.parent_informed}<span class="text-xs text-emerald-600"
                >Parent informed</span
              >{/if}
            {#if l.reviewed_by_principal}<span class="text-xs text-blue-600"
                >Reviewed by Principal</span
              >{/if}
            {#if l.severity === "urgent" || l.severity === "high"}<AlertTriangle
                size={14}
                class="text-amber-500"
              />{/if}
          </div>
        </div>
      {:else}
        <div class="bg-white rounded-xl border border-slate-200">
          <EmptyState icon={ClipboardList} title="No logs yet" />
        </div>
      {/each}
    </div>
  {/if}
</div>
