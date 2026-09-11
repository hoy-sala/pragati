<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import type { Class } from "$lib/types";
  import { onMount } from "svelte";
  import Select from "$lib/components/Select.svelte";
  import Button from "$lib/components/Button.svelte";

  let classes = $state<Class[]>([]);
  let selectedClass = $state("");
  let selectedStage = $state("middle");
  let loadedConfig = $state<any>(null);
  let gradingScheme = $state<any[]>([]);
  let proficiencyScale = $state<any[]>([]);
  let terms = $state(["Term1", "Term2"]);
  let loading = $state(true);
  let saving = $state(false);
  let statusMsg = $state("");
  let statusType = $state<"success" | "error">("success");

  const stages = [
    { value: "foundational", label: "Foundational (Class 1-2)" },
    { value: "preparatory", label: "Preparatory (Class 3-5)" },
    { value: "middle", label: "Middle (Class 6-8)" },
    { value: "secondary", label: "Secondary (Class 9-10)" },
  ];
  const stageOptions = stages.map((s) => ({ id: s.value, name: s.label }));

  onMount(async () => {
    const cRes = await api<Class[]>("GET", "/classes?limit=100");
    if (cRes.data) classes = cRes.data;
    gradingScheme = [
      { grade: "A+", min_pct: 90, max_pct: 100, descriptor: "Outstanding" },
      { grade: "A", min_pct: 75, max_pct: 89, descriptor: "Excellent" },
      { grade: "B+", min_pct: 60, max_pct: 74, descriptor: "Very Good" },
      { grade: "B", min_pct: 45, max_pct: 59, descriptor: "Good" },
      { grade: "C", min_pct: 33, max_pct: 44, descriptor: "Satisfactory" },
      { grade: "D", min_pct: 20, max_pct: 32, descriptor: "Needs Improvement" },
      { grade: "E", min_pct: 0, max_pct: 19, descriptor: "Requires Remedial" },
    ];
    proficiencyScale = [
      { level: 1, label: "Beginning", descriptor: "Needs significant support" },
      {
        level: 2,
        label: "Developing",
        descriptor: "Occasional support needed",
      },
      {
        level: 3,
        label: "Proficient",
        descriptor: "Meets expectations independently",
      },
      { level: 4, label: "Advanced", descriptor: "Exceeds expectations" },
    ];
    loading = false;
  });

  async function loadConfig() {
    const params = new URLSearchParams();
    if (selectedClass) params.set("class_id", selectedClass);
    params.set("stage", selectedStage);
    const res = await api<any>("GET", "/hpc/config?" + params.toString());
    if (res.data) {
      loadedConfig = res.data;
      gradingScheme = res.data.grading_scheme || gradingScheme;
      proficiencyScale = res.data.proficiency_scale || proficiencyScale;
      terms = res.data.terms || terms;
      statusMsg = "Configuration loaded.";
      statusType = "success";
    } else {
      statusMsg = "No existing config. Create a new one.";
      statusType = "success";
    }
  }

  async function saveConfig() {
    saving = true;
    statusMsg = "";
    const res = await api("PUT", "/hpc/config", {
      stage: selectedStage,
      class_id: selectedClass || undefined,
      academic_year_id: undefined,
      grading_scheme: gradingScheme,
      proficiency_scale: proficiencyScale,
      co_scholastic_areas: {},
      health_params: {},
      terms,
    });
    saving = false;
    if (res.data) {
      statusMsg = "Configuration saved.";
      statusType = "success";
    } else if (res.error) {
      statusMsg = res.error.message;
      statusType = "error";
    }
  }

  function addGrade() {
    gradingScheme = [
      ...gradingScheme,
      { grade: "", min_pct: 0, max_pct: 0, descriptor: "" },
    ];
  }
</script>

<div class="max-w-3xl mx-auto space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold text-slate-900">HPC Configuration</h1>
    <Button
      onclick={saveConfig}
      disabled={saving}
      loading={saving}
    >
      {saving ? "Saving..." : "Save Config"}
    </Button>
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

  <div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
    <div class="flex flex-wrap gap-3 items-end">
      <div>
        <label
          for="config-stage"
          class="block text-xs font-medium text-slate-600 mb-1">Stage</label
        >
        <Select
          id="config-stage"
          bind:value={selectedStage}
          options={stageOptions}
        />
      </div>
      <div>
        <label
          for="config-class"
          class="block text-xs font-medium text-slate-600 mb-1"
          >Class (optional)</label
        >
        <Select
          id="config-class"
          bind:value={selectedClass}
          options={classes}
          placeholder="All Classes"
        />
      </div>
      <Button variant="secondary" onclick={loadConfig}>
        Load Existing
      </Button>
    </div>
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
    <h2 class="text-sm font-semibold text-slate-700 border-b pb-2">
      Grading Scheme
    </h2>
    <table class="w-full text-sm">
      <thead>
        <tr class="bg-slate-50">
          <th class="text-left px-3 py-2 font-medium">Grade</th>
          <th class="text-center px-3 py-2 font-medium" style="width:80px;"
            >Min %</th
          >
          <th class="text-center px-3 py-2 font-medium" style="width:80px;"
            >Max %</th
          >
          <th class="text-left px-3 py-2 font-medium">Descriptor</th>
        </tr>
      </thead>
      <tbody>
        {#each gradingScheme as g, i}
          <tr class="border-t border-slate-100">
            <td class="px-3 py-2"
              ><input
                type="text"
                bind:value={gradingScheme[i].grade}
                class="w-12 px-2 py-1 border border-slate-300 rounded text-sm"
              /></td
            >
            <td class="px-3 py-2"
              ><input
                type="number"
                bind:value={gradingScheme[i].min_pct}
                class="w-20 px-2 py-1 border border-slate-300 rounded text-sm text-center"
              /></td
            >
            <td class="px-3 py-2"
              ><input
                type="number"
                bind:value={gradingScheme[i].max_pct}
                class="w-20 px-2 py-1 border border-slate-300 rounded text-sm text-center"
              /></td
            >
            <td class="px-3 py-2"
              ><input
                type="text"
                bind:value={gradingScheme[i].descriptor}
                class="w-full px-2 py-1 border border-slate-300 rounded text-sm"
              /></td
            >
          </tr>
        {/each}
      </tbody>
    </table>
    <button
      onclick={addGrade}
      class="text-xs text-primary-600 hover:text-primary-800"
      >+ Add Grade</button
    >
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
    <h2 class="text-sm font-semibold text-slate-700 border-b pb-2">
      Proficiency Scale
    </h2>
    <table class="w-full text-sm">
      <thead>
        <tr class="bg-slate-50">
          <th class="text-center px-3 py-2 font-medium" style="width:60px;"
            >Level</th
          >
          <th class="text-left px-3 py-2 font-medium">Label</th>
          <th class="text-left px-3 py-2 font-medium">Descriptor</th>
        </tr>
      </thead>
      <tbody>
        {#each proficiencyScale as p, i}
          <tr class="border-t border-slate-100">
            <td class="px-3 py-2 text-center font-medium">{p.level}</td>
            <td class="px-3 py-2"
              ><input
                type="text"
                bind:value={proficiencyScale[i].label}
                class="w-32 px-2 py-1 border border-slate-300 rounded text-sm"
              /></td
            >
            <td class="px-3 py-2"
              ><input
                type="text"
                bind:value={proficiencyScale[i].descriptor}
                class="w-full px-2 py-1 border border-slate-300 rounded text-sm"
              /></td
            >
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
    <h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Terms</h2>
    <div class="flex gap-3">
      {#each terms as t, i}
        <input
          type="text"
          bind:value={terms[i]}
          class="w-28 px-2 py-1 border border-slate-300 rounded text-sm"
        />
      {/each}
    </div>
  </div>
</div>
