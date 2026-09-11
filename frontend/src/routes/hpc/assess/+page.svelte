<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import type { Class, Subject } from "$lib/types";
  import { onMount } from "svelte";
  import Select from "$lib/components/Select.svelte";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import PageTabs from "$lib/components/PageTabs.svelte";
  import { HPC_TABS } from "$lib/utils/tabs";
  import { getAuthState } from "$lib/stores/auth.svelte";
  import { effectiveRole } from "$lib/utils/roles";
  import { toast } from "$lib/stores/toast.svelte";

  const auth = getAuthState();
  let role = $derived(effectiveRole(auth.currentUser));

  let classes = $state<Class[]>([]);
  let subjects = $state<Subject[]>([]);
  let selectedClass = $state("");
  let selectedSubject = $state("");
  let selectedTerm = $state("Term 1");

  let loColumns = $state<any[]>([]);
  let students = $state<any[]>([]);
  let gridData = $state<any[]>([]);
  let loading = $state(false);
  let statusMsg = $state("");
  let statusType = $state<"success" | "error">("success");

  const terms = ["Term 1", "Term 2"];
  const termOptions = terms.map((t) => ({ id: t, name: t }));
  const proficiencyOptions = [
    { id: "0", name: "—" },
    { id: "1", name: "1-Beginning" },
    { id: "2", name: "2-Developing" },
    { id: "3", name: "3-Proficient" },
    { id: "4", name: "4-Advanced" },
  ];

  onMount(async () => {
    const [cRes, sRes] = await Promise.all([
      api<Class[]>("GET", "/classes?limit=100"),
      api<Subject[]>("GET", "/subjects?limit=50"),
    ]);
    if (cRes.data) classes = cRes.data;
    if (sRes.data) subjects = sRes.data;
  });

  async function loadGrid() {
    if (!selectedClass || !selectedSubject) return;
    loading = true;
    statusMsg = "";
    const params = `class_id=${selectedClass}&subject_id=${selectedSubject}&term=${selectedTerm}`;
    const res = await api<any>("GET", "/hpc/assessments?" + params);
    loading = false;
    if (res.data) {
      loColumns = res.data.columns || [];
      students = res.data.students || [];
      gridData = res.data.grid || [];
    } else if (res.error) {
      statusMsg = res.error.message;
      statusType = "error";
    }
  }

  $effect(() => {
    if (selectedClass && selectedSubject) loadGrid();
  });

  function getCellValue(studentId: string, loId: string): number {
    const row = gridData.find((r: any) => r.student.student_id === studentId);
    if (row && row.cells[loId]) return row.cells[loId].level || 0;
    return 0;
  }

  async function updateCell(studentId: string, loId: string, level: number) {
    const assessment = {
      learning_outcome_id: loId,
      proficiency_level: level,
    };
    const res = await api("POST", "/hpc/assess", {
      student_id: studentId,
      subject_id: selectedSubject,
      term: selectedTerm,
      assessments: [assessment],
    });
    if (res.error) {
      toast(`Failed to save assessment: ${res.error.message}`, "error");
    }
  }

</script>

<div class="space-y-4">
  <PageHeader
    title="Learning Outcome Assessment Grid"
    subtitle="Changes save automatically when you pick a level"
  />

  <PageTabs tabs={HPC_TABS} role={role} />

  <div class="bg-white rounded-xl border border-slate-200 p-4">
    <div class="flex flex-wrap gap-3 items-end">
      <div>
        <label
          for="assess-class"
          class="block text-xs font-medium text-slate-600 mb-1">Class</label
        >
        <Select
          id="assess-class"
          bind:value={selectedClass}
          options={classes}
          placeholder="Select"
        />
      </div>
      <div>
        <label
          for="assess-subject"
          class="block text-xs font-medium text-slate-600 mb-1">Subject</label
        >
        <Select
          id="assess-subject"
          bind:value={selectedSubject}
          options={subjects}
          placeholder="Select"
        />
      </div>
      <div>
        <label
          for="assess-term"
          class="block text-xs font-medium text-slate-600 mb-1">Term</label
        >
        <Select
          id="assess-term"
          bind:value={selectedTerm}
          options={termOptions}
        />
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

  <div class="bg-white rounded-xl border border-slate-200 overflow-x-auto">
    {#if loading}
      <div class="p-8 text-center text-sm text-slate-400">Loading...</div>
    {:else if loColumns.length === 0}
      <div class="p-8 text-center text-sm text-slate-400">
        {selectedSubject
          ? "No learning outcomes configured for this subject. Import them first."
          : "Select a subject and class to load."}
      </div>
    {:else}
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-slate-50">
            <th
              class="sticky left-0 bg-slate-50 z-10 px-3 py-2 text-left font-medium text-xs"
              style="min-width:180px;">Student</th
            >
            {#each loColumns as lo}
              <th
                class="px-2 py-2 text-center font-medium text-xs"
                style="min-width:80px; max-width:100px;"
                title={lo.description}
              >
                {lo.code}
              </th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each students as student}
            <tr class="border-t border-slate-100 hover:bg-slate-50">
              <td
                class="sticky left-0 bg-white z-10 px-3 py-1.5 font-medium text-xs"
                >{student.name}</td
              >
              {#each loColumns as lo}
                <td class="px-1 py-1.5 text-center">
                  <Select
                    value={String(getCellValue(student.student_id, lo.id))}
                    options={proficiencyOptions}
                    size="sm"
                    onselect={(v) =>
                      updateCell(student.student_id, lo.id, parseInt(v))}
                  />
                </td>
              {/each}
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </div>
</div>
