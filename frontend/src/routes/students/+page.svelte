<script lang="ts">
  import { api, apiUpload } from "$lib/api/client.svelte";
  import {
    Plus,
    Pencil,
    Trash2,
    User,
    Hash,
    Phone,
    Mail,
    MapPin,
    Calendar,
    Droplets,
    Users,
    Eye,
    EyeOff,
    Upload,
    X,
  } from "lucide-svelte";
  import Button from "$lib/components/Button.svelte";
  import Select from "$lib/components/Select.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import type { Student, Class, AcademicYear, ImportResult } from "$lib/types";
  import { toast } from "$lib/stores/toast.svelte";
  import { getAuthState } from "$lib/stores/auth.svelte";
  import { hasRole } from "$lib/utils/roles";
  import { onMount } from "svelte";

  const auth = getAuthState();
  let isAdmin = $derived(hasRole(auth.currentUser, "admin"));

  let allStudents: Student[] = $state([]);
  let classes: Class[] = $state([]);
  let academicYears: AcademicYear[] = $state([]);
  let loading = $state(true);
  let saving = $state(false);
  let error = $state("");

  let showForm = $state(false);
  let showImport = $state(false);
  let importFile = $state<File | null>(null);
  let importing = $state(false);
  let importResult = $state<ImportResult | null>(null);
  let importError = $state("");
  let editingId: string | null = $state(null);
  let formSATS = $state("");
  let formFirstName = $state("");
  let formLastName = $state("");
  let formRollNo = $state(0);
  let formGender = $state("");
  let formDOB = $state("");
  let formBloodGroup = $state("");
  let formPhone = $state("");
  let formEmail = $state("");
  let formAddress = $state("");
  let formClassId = $state("");
  let formAcademicYearId = $state("");
  let formFatherName = $state("");
  let formMotherName = $state("");
  let formParentName = $state("");
  let formParentPhone = $state("");
  let formParentEmail = $state("");

  // Search, filter, pagination
  let search = $state("");
  let filterClass = $state("");
  let page = $state(1);
  const pageSize = 20;

  let filteredStudents = $derived(
    allStudents.filter((s) => {
      if (filterClass && s.class_id !== filterClass) return false;
      if (!search.trim()) return true;
      const q = search.toLowerCase();
      return (
        s.first_name?.toLowerCase().includes(q) ||
        s.last_name?.toLowerCase().includes(q) ||
        s.sats_number?.includes(q)
      );
    }),
  );

  let paginatedStudents = $derived(
    filteredStudents.slice((page - 1) * pageSize, page * pageSize),
  );
  let totalStudents = $derived(filteredStudents.length);

  function onPageChange(p: number) {
    page = p;
  }
  function resetPage() {
    page = 1;
  }

  $effect(() => {
    void search;
    void filterClass;
    page = 1;
  });

  onMount(async () => {
    const [sRes, cRes, yRes] = await Promise.all([
      api<Student[]>("GET", "/students?limit=500"),
      api<Class[]>("GET", "/classes?limit=50"),
      api<AcademicYear[]>("GET", "/academic-years?limit=50"),
    ]);
    if (sRes.data) allStudents = sRes.data;
    if (cRes.data) classes = cRes.data;
    if (yRes.data) {
      academicYears = yRes.data;
      const current = yRes.data.find((y) => y.is_current);
      if (current) formAcademicYearId = current.id;
    }
    loading = false;
  });

  function resetForm() {
    formSATS = "";
    formFirstName = "";
    formLastName = "";
    formRollNo = 0;
    formGender = "";
    formDOB = "";
    formBloodGroup = "";
    formPhone = "";
    formEmail = "";
    formAddress = "";
    formClassId = "";
    formFatherName = "";
    formMotherName = "";
    formParentName = "";
    formParentPhone = "";
    formParentEmail = "";
    if (academicYears.length > 0) {
      const current = academicYears.find((y) => y.is_current);
      formAcademicYearId = current?.id ?? academicYears[0].id;
    }
    error = "";
  }

  function openCreate() {
    editingId = null;
    resetForm();
    showForm = true;
  }

  function openEdit(s: Student) {
    editingId = s.id;
    formSATS = s.sats_number;
    formFirstName = s.first_name;
    formLastName = s.last_name ?? "";
    formRollNo = s.roll_no ?? 0;
    formGender = s.gender ?? "";
    formDOB = s.date_of_birth ? s.date_of_birth.slice(0, 10) : "";
    formBloodGroup = s.blood_group ?? "";
    formPhone = s.phone ?? "";
    formEmail = s.email ?? "";
    formAddress = s.address ?? "";
    formClassId = s.class_id;
    formAcademicYearId = s.academic_year_id;
    formFatherName = s.father_name ?? "";
    formMotherName = s.mother_name ?? "";
    formParentName = s.parent_name ?? "";
    formParentPhone = s.parent_phone ?? "";
    formParentEmail = s.parent_email ?? "";
    error = "";
    showForm = true;
  }

  function cancelForm() {
    showForm = false;
    editingId = null;
    resetForm();
  }

  function openImport() {
    showImport = true;
    importFile = null;
    importResult = null;
    importError = "";
  }

  function closeImport() {
    showImport = false;
    importing = false;
    importFile = null;
    importResult = null;
    importError = "";
  }

  const IMPORT_HEADER =
    "sats_number,first_name,last_name,class,date_of_birth,gender,father_name,mother_name,parent_name";

  function csvEscape(v: string): string {
    return /["\n,]/.test(v) ? '"' + v.replace(/"/g, '""') + '"' : v;
  }

  function cleanRegisterField(v: string): string {
    return v === "" || /^N\/?A$/i.test(v) ? "" : v;
  }

  function registerToCsv(parsed: unknown): string {
    const maybeArr = parsed as Record<string, unknown>;
    const arr = Array.isArray(parsed)
      ? parsed
      : Array.isArray(maybeArr?.students)
        ? (maybeArr.students as unknown[])
        : Array.isArray(maybeArr?.data)
          ? (maybeArr.data as unknown[])
          : null;
    if (!arr) throw new Error("Unrecognized JSON. Expected an array of student records (register export).");
    const lines = [IMPORT_HEADER];
    for (const rec of arr) {
      if (!rec || typeof rec !== "object") continue;
      const f: Record<string, string> = {};
      for (const [k, v] of Object.entries(rec as Record<string, unknown>)) {
        f[k.toLowerCase()] = String(v ?? "").trim();
      }
      const sats = f.student_id.replace(/\s+/g, "");
      const name = f.student_name.trim();
      const first = name.split(/\s+/)[0] || name;
      const last = name.split(/\s+/).slice(1).join(" ");
      const clsMatch = f.class_studying.match(/(\d+)/);
      const cls = clsMatch ? "Class " + clsMatch[1] : "";
      let dob = "";
      const dc = cleanRegisterField(f.dob).split("/");
      if (
        dc.length === 3 &&
        /^\d{1,2}$/.test(dc[0]) &&
        /^\d{1,2}$/.test(dc[1]) &&
        /^\d{4}$/.test(dc[2])
      ) {
        dob = `${dc[2]}-${dc[1].padStart(2, "0")}-${dc[0].padStart(2, "0")}`;
      }
      let gender = "";
      if (/^1/.test(f.sex)) gender = "male";
      else if (/^2/.test(f.sex)) gender = "female";
      const father = cleanRegisterField(f.father_name);
      const mother = cleanRegisterField(f.mother_name);
      lines.push([sats, first, last, cls, dob, gender, father, mother, father || mother].map(csvEscape).join(","));
    }
    return lines.join("\n");
  }

  async function importStudents() {
    if (!importFile) {
      importError = "Select a CSV or JSON file.";
      return;
    }
    importing = true;
    importError = "";
    importResult = null;
    let uploadFile = importFile;
    if (importFile.name.toLowerCase().endsWith(".json")) {
      try {
        const parsed = JSON.parse(await importFile.text());
        const csv = registerToCsv(parsed);
        const rowCount = csv.trim().split(/\n/).length - 1;
        if (rowCount <= 0) throw new Error("No student records found in the JSON file.");
        uploadFile = new File([csv], "students_import.csv", { type: "text/csv" });
      } catch (e) {
        const msg = e instanceof Error ? e.message : "Invalid JSON file.";
        importError = msg;
        toast(msg, "error");
        importing = false;
        return;
      }
    }
    const formData = new FormData();
    formData.append("file", uploadFile);
    try {
      const json = await apiUpload<ImportResult>("/students/import", formData);
      if (json.data) {
        importResult = json.data;
        const sRes = await api<Student[]>("GET", "/students?limit=500");
        if (sRes.data) allStudents = sRes.data;
        toast("Student import finished", "success");
      } else if (json.error) {
        importError = json.error.message;
        toast(importError, "error");
      } else {
        importError = "Import failed with no response.";
        toast(importError, "error");
      }
    } catch {
      importError = "Import failed. Please try again.";
      toast(importError, "error");
    } finally {
      importing = false;
    }
  }

  async function save() {
    if (!formSATS.trim() || !formFirstName.trim() || !formClassId) return;
    if (formSATS.trim().length !== 9) {
      error = "SATS number must be exactly 9 characters";
      return;
    }
    saving = true;
    error = "";

    const body = {
      sats_number: formSATS.trim(),
      first_name: formFirstName.trim(),
      last_name: formLastName.trim() || undefined,
      roll_no: formRollNo || undefined,
      gender: formGender || undefined,
      date_of_birth: formDOB || undefined,
      blood_group: formBloodGroup || undefined,
      phone: formPhone.trim() || undefined,
      email: formEmail.trim() || undefined,
      address: formAddress.trim() || undefined,
      class_id: formClassId,
      academic_year_id: formAcademicYearId,
      father_name: formFatherName.trim() || undefined,
      mother_name: formMotherName.trim() || undefined,
      parent_name: formParentName.trim() || undefined,
      parent_phone: formParentPhone.trim() || undefined,
      parent_email: formParentEmail.trim() || undefined,
    };

    if (editingId) {
      const res = await api<Student>("PUT", `/students/${editingId}`, body);
      saving = false;
      if (res.error) {
        error = res.error.message;
        return;
      }
      allStudents = allStudents.map((s) =>
        s.id === editingId ? { ...s, ...body } : s,
      );
      toast("Student updated", "success");
    } else {
      const res = await api<{ id: string }>("POST", "/students", body);
      saving = false;
      if (res.error) {
        error = res.error.message;
        return;
      }
      allStudents = [
        ...allStudents,
        {
          id: res.data!.id,
          school_id: "",
          sats_number: formSATS.trim(),
          first_name: formFirstName.trim(),
          last_name: formLastName.trim() || undefined,
          roll_no: formRollNo || undefined,
          gender: formGender || undefined,
          blood_group: formBloodGroup || undefined,
          phone: formPhone.trim() || undefined,
          email: formEmail.trim() || undefined,
          address: formAddress.trim() || undefined,
          class_id: formClassId,
          academic_year_id: formAcademicYearId,
          father_name: formFatherName.trim() || undefined,
          mother_name: formMotherName.trim() || undefined,
          parent_name: formParentName.trim() || undefined,
          parent_phone: formParentPhone.trim() || undefined,
          parent_email: formParentEmail.trim() || undefined,
          is_active: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        } as Student,
      ];
      toast("Student created", "success");
    }
    showForm = false;
    editingId = null;
    resetForm();
  }

  async function removeStudent(id: string) {
    if (!confirm("Delete this student? This action cannot be undone.")) return;
    const res = await api("DELETE", `/students/${id}`);
    if (res.error) {
      toast(res.error.message, "error");
      return;
    }
    allStudents = allStudents.filter((s) => s.id !== id);
    toast("Student deleted", "success");
  }

  function className(id: string): string {
    return classes.find((c) => c.id === id)?.name ?? id;
  }

  function yearName(id: string): string {
    return academicYears.find((y) => y.id === id)?.name ?? id;
  }

  function genderLabel(v: string): string {
    if (v === "male") return "Male";
    if (v === "female") return "Female";
    return "—";
  }

  function genderBadge(v: string): string {
    if (v === "male") return "bg-blue-100 text-blue-700";
    if (v === "female") return "bg-pink-100 text-pink-700";
    return "bg-slate-100 text-slate-500";
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold text-slate-900">Students</h1>
      <p class="text-sm text-slate-500 mt-1">
        {allStudents.length} student{allStudents.length !== 1 ? "s" : ""} enrolled
      </p>
    </div>
    <div class="flex items-center gap-2">
      {#if isAdmin}
      <Button onclick={openImport} variant="secondary" icon={Upload}
        >Import CSV</Button
      >
      {/if}
      <Button onclick={openCreate} icon={Plus}>Add Student</Button>
    </div>
  </div>

  {#if showImport}
    <div
      class="bg-white rounded-xl border border-slate-200 shadow-sm"
      role="dialog"
      aria-label="Import students from CSV"
    >
      <div
        class="px-4 py-3 border-b border-slate-100 flex items-center justify-between"
      >
        <h3
          class="text-sm font-semibold text-slate-700 flex items-center gap-2"
        >
          <Upload size={16} class="text-primary-500" />
          Import Students
        </h3>
        <button
          onclick={closeImport}
          aria-label="Close import dialog"
          class="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-md hover:bg-slate-100"
        >
          <X size={16} />
        </button>
      </div>
      <div class="p-4 space-y-4">
        <p class="text-sm text-slate-500">
          Upload a CSV or JSON of students. Rows whose SATS number already
          exists are updated; new ones are inserted. Leave the academic year
          column empty to use the current year.
        </p>
        <div>
          <label for="st-import-csv" class="block text-xs font-medium text-slate-500 mb-1"
            >File</label
          >
          <input
            id="st-import-csv"
            type="file"
            accept=".csv,.json"
            onchange={(e: Event) => {
              const el = e.target as HTMLInputElement;
              importFile = el.files?.[0] || null;
            }}
            class="w-full text-sm file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-medium file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
          />
        </div>
        <details class="text-xs text-slate-500">
          <summary class="cursor-pointer hover:text-slate-700"
            >Accepted formats</summary
          >
          <p class="mt-2">
            <strong>JSON</strong> — register export (keys like
            <code>student_name</code>, <code>student_id</code>,
            <code>class_studying</code>, <code>dob</code>
            (DD/MM/YYYY), <code>sex</code> (1-BOY/2-GIRL),
            <code>father_name</code>, <code>mother_name</code>).
          </p>
          <pre
            class="mt-2 p-3 bg-slate-50 rounded-lg text-xs leading-relaxed overflow-x-auto"
          ><code>sats_number,first_name,last_name,class,date_of_birth,gender,father_name,mother_name,parent_name,parent_phone,parent_email,admission_no,roll_no,academic_year</code></pre>
          <p class="mt-2"><code>sats_number</code>, <code>first_name</code> and <code>class</code> are required. <code>class</code> is the class name or code, e.g. <code>Class 6</code>.</p>
          <p class="mt-2"><code>date_of_birth</code> uses <code>YYYY-MM-DD</code>. <code>gender</code> is <code>male</code> or <code>female</code>.</p>
        </details>

        {#if importError}
          <div
            class="flex items-center gap-2 text-sm px-3 py-2 rounded-lg bg-red-50 text-danger-600 border border-red-200"
          >
            <span>{importError}</span>
          </div>
        {/if}

        {#if importResult}
          <div
            class="text-sm bg-green-50 border border-green-200 rounded-lg p-4"
          >
            <p class="text-green-700 font-medium">
              Imported {importResult.imported} student{importResult.imported !== 1 ? "s" : ""}
              {importResult.skipped > 0 ? `, skipped ${importResult.skipped}` : ""}
            </p>
            {#if importResult.errors && importResult.errors.length > 0}
              <ul class="mt-2 space-y-1 max-h-48 overflow-y-auto">
                {#each importResult.errors as e}
                  <li class="text-amber-700">
                    Row {e.row} ({e.sats_number}): {e.message}
                  </li>
                {/each}
              </ul>
            {/if}
          </div>
        {/if}

        <div class="flex items-center gap-2 pt-1">
          <Button
            onclick={importStudents}
            disabled={importing || !importFile}
            loading={importing}
          >
            {importing ? "Importing..." : "Import"}
          </Button>
          <Button onclick={closeImport} variant="secondary">Cancel</Button>
        </div>
      </div>
    </div>
  {/if}

  {#if showForm}
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm">
      <div
        class="px-4 py-3 border-b border-slate-100 flex items-center justify-between"
      >
        <h3
          class="text-sm font-semibold text-slate-700 flex items-center gap-2"
        >
          <User size={16} class="text-primary-500" />
          {editingId ? "Edit Student" : "New Student"}
        </h3>
        <button
          onclick={cancelForm}
          aria-label="Close form"
          class="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-md hover:bg-slate-100"
        >
          <EyeOff size={16} />
        </button>
      </div>
      <div class="p-4 space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div>
            <label
              for="st-sats"
              class="block text-xs font-medium text-slate-500 mb-1"
              >SATS Number *</label
            >
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><Hash size={14} /></span
              >
              <input
                id="st-sats"
                bind:value={formSATS}
                placeholder="9-digit SATS number"
                maxlength={9}
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
          </div>
          <div>
            <label
              for="st-first"
              class="block text-xs font-medium text-slate-500 mb-1"
              >First Name *</label
            >
            <input
              id="st-first"
              bind:value={formFirstName}
              placeholder="First name"
              class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
            />
          </div>
          <div>
            <label
              for="st-last"
              class="block text-xs font-medium text-slate-500 mb-1"
              >Last Name</label
            >
            <input
              id="st-last"
              bind:value={formLastName}
              placeholder="Last name"
              class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-4 gap-3">
          <div>
            <label
              for="st-roll"
              class="block text-xs font-medium text-slate-500 mb-1"
              >Roll No</label
            >
            <input
              id="st-roll"
              bind:value={formRollNo}
              type="number"
              min="0"
              placeholder="Roll number"
              class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
            />
          </div>
          <div>
            <label
              for="st-gender"
              class="block text-xs font-medium text-slate-500 mb-1"
              >Gender</label
            >
            <Select
              id="st-gender"
              bind:value={formGender}
              options={[
                { id: "male", name: "Male" },
                { id: "female", name: "Female" },
              ]}
              placeholder="Select gender"
              icon={Users}
            />
          </div>
          <div>
            <label
              for="st-dob"
              class="block text-xs font-medium text-slate-500 mb-1"
              >Date of Birth</label
            >
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><Calendar size={14} /></span
              >
              <input
                id="st-dob"
                bind:value={formDOB}
                type="date"
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
          </div>
          <div>
            <label
              for="st-blood"
              class="block text-xs font-medium text-slate-500 mb-1"
              >Blood Group</label
            >
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><Droplets size={14} /></span
              >
              <input
                id="st-blood"
                bind:value={formBloodGroup}
                placeholder="e.g. O+"
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div>
            <label
              for="st-phone"
              class="block text-xs font-medium text-slate-500 mb-1">Phone</label
            >
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><Phone size={14} /></span
              >
              <input
                id="st-phone"
                bind:value={formPhone}
                placeholder="Phone number"
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
          </div>
          <div>
            <label
              for="st-email"
              class="block text-xs font-medium text-slate-500 mb-1">Email</label
            >
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><Mail size={14} /></span
              >
              <input
                id="st-email"
                bind:value={formEmail}
                type="email"
                placeholder="Email address"
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
          </div>
          <div>
            <label
              for="st-address"
              class="block text-xs font-medium text-slate-500 mb-1"
              >Address</label
            >
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><MapPin size={14} /></span
              >
              <input
                id="st-address"
                bind:value={formAddress}
                placeholder="Address"
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <Select
            bind:value={formClassId}
            label="Class *"
            options={classes}
            icon={Users}
            placeholder="Select class"
          />
        </div>

        <div class="border-t border-slate-100 pt-4">
          <p
            class="text-xs font-medium text-slate-500 mb-2 flex items-center gap-1.5"
          >
            <Users size={12} /> Parents / Guardian
          </p>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <input
              bind:value={formFatherName}
              placeholder="Father name"
              class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
            />
            <input
              bind:value={formMotherName}
              placeholder="Mother name"
              class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
            />
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 mt-3">
            <input
              bind:value={formParentName}
              placeholder="Guardian name"
              class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
            />
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><Phone size={14} /></span
              >
              <input
                bind:value={formParentPhone}
                placeholder="Guardian phone"
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
            <div class="relative">
              <span
                class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"
                ><Mail size={14} /></span
              >
              <input
                bind:value={formParentEmail}
                type="email"
                placeholder="Guardian email"
                class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
              />
            </div>
          </div>
        </div>

        {#if error}
          <div
            class="flex items-center gap-2 text-sm px-3 py-2 rounded-lg bg-red-50 text-danger-600 border border-red-200"
          >
            <span>{error}</span>
          </div>
        {/if}

        <div class="flex items-center gap-2 pt-1">
          <Button
            onclick={save}
            disabled={saving ||
              !formSATS.trim() ||
              !formFirstName.trim() ||
              !formClassId}
            loading={saving}
          >
            {editingId ? "Update Student" : "Create Student"}
          </Button>
          <Button onclick={cancelForm} variant="secondary">Cancel</Button>
        </div>
      </div>
    </div>
  {/if}

  <div
    class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-sm"
  >
    <div class="p-3 border-b border-slate-100 flex flex-col sm:flex-row gap-3">
      <div class="sm:flex-1">
        <SearchFilter
          bind:value={search}
          placeholder="Search by name or SATS..."
        />
      </div>
      <div class="sm:w-56">
        <Select
          bind:value={filterClass}
          options={classes}
          placeholder="All classes"
          clearable
        />
      </div>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-slate-50 text-slate-600">
            <th
              class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider"
              >SATS</th
            >
            <th
              class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider"
              >Name</th
            >
            <th
              class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider"
              >Class</th
            >
            <th
              class="text-right px-4 py-3 font-semibold text-xs uppercase tracking-wider"
              >Actions</th
            >
          </tr>
        </thead>
        <tbody>
          {#if loading}
            <tr
              ><td colspan="4" class="px-4 py-12 text-center text-slate-400"
                >Loading...</td
              ></tr
            >
          {:else if filteredStudents.length === 0}
            <tr
              ><td colspan="4" class="px-4 py-12 text-center text-slate-400"
                >No students match your search.</td
              ></tr
            >
          {:else}
            {#each paginatedStudents as s (s.id)}
              <tr
                class="border-t border-slate-100 hover:bg-slate-50 transition-colors"
              >
                <td class="px-4 py-3 font-mono text-xs text-slate-500"
                  >{s.sats_number}</td
                >
                <td class="px-4 py-3 font-medium"
                  >{s.first_name} {s.last_name}</td
                >
                <td class="px-4 py-3">{className(s.class_id)}</td>
                <td class="px-4 py-3 text-right">
                  <Button
                    onclick={() => openEdit(s)}
                    variant="ghost"
                    size="sm"
                    icon={Pencil}>Edit</Button
                  >
                  {#if isAdmin}
                  <Button
                    onclick={() => removeStudent(s.id)}
                    variant="ghost"
                    size="sm"
                    icon={Trash2}>Delete</Button
                  >
                  {/if}
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
    <Pagination
      {page}
      total={totalStudents}
      {pageSize}
      onChange={onPageChange}
    />
  </div>
</div>
