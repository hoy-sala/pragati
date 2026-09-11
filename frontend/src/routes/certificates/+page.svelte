<script lang="ts">
  import { api, apiUrl, apiUpload } from "$lib/api/client.svelte";
  import type {
    CertificateEvent,
    CertificateParticipant,
    CertificateSignatory,
    Student,
  } from "$lib/types";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { Trash2, Plus } from "lucide-svelte";
  import Button from "$lib/components/Button.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import Select from "$lib/components/Select.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import { toast } from "$lib/stores/toast.svelte";

  const CATEGORIES = [
    { id: "sports", name: "Sports" },
    { id: "cultural", name: "Cultural" },
    { id: "academic", name: "Academic" },
    { id: "other", name: "Other" },
  ];

  const POSITIONS = [
    { id: "1st", name: "1st (First Prize)" },
    { id: "2nd", name: "2nd (Runner Up)" },
    { id: "3rd", name: "3rd (Consolation)" },
    { id: "participation", name: "Participation" },
  ];

  const SIGNATORY_ROLES = [
    { id: "principal", name: "Principal" },
    { id: "chief_guest", name: "Chief Guest" },
    { id: "chief_judge", name: "Chief Judge" },
    { id: "judge", name: "Judge" },
    { id: "coordinator", name: "Event Coordinator" },
  ];

  let events: CertificateEvent[] = $state([]);
  let students: Student[] = $state([]);
  let loading = $state(true);
  let savingEvent = $state(false);
  let savingParticipant = $state(false);
  let savingSignatory = $state(false);
  let error = $state("");

  let showForm = $state(false);
  let search = $state("");
  let newName = $state("");
  let newCategory = $state("sports");
  let newHeldDate = $state("");
  let newVenue = $state("");
  let newDescription = $state("");
  let newAcademicYearId = $state("");

  let expandedEventId = $state<string | null>(null);
  let eventDetails = $state<
    Record<
      string,
      {
        participants: CertificateParticipant[];
        signatories: CertificateSignatory[];
      }
    >
  >({});
  let loadingDetails = $state(false);

  let filteredEvents = $derived(
    events.filter((e) => {
      if (!search.trim()) return true;
      const q = search.toLowerCase();
      return (
        e.name?.toLowerCase().includes(q) ||
        e.venue?.toLowerCase().includes(q) ||
        categoryLabel(e.category).toLowerCase().includes(q)
      );
    }),
  );

  let studentOptions = $derived(
    students.map((s) => ({
      id: s.id,
      name:
        `${s.first_name} ${s.last_name || ""}`.trim() +
        (s.sats_number ? ` (${s.sats_number})` : ""),
    })),
  );

  let partForm: Record<
    string,
    {
      student_id: string;
      position: string;
      prize_title: string;
      issue_date: string;
    }
  > = $state({});

  let signForm: Record<
    string,
    {
      name: string;
      role: string;
      title: string;
      signature_url: string;
      uploading: boolean;
      error: string;
    }
  > = $state({});

  onMount(async () => {
    const [eventRes, yearRes, studentRes] = await Promise.all([
      api<CertificateEvent[]>("GET", "/certificates/events?limit=100"),
      api<{ id: string; is_current: boolean }[]>(
        "GET",
        "/academic-years?limit=50",
      ),
      api<Student[]>("GET", "/students?limit=500"),
    ]);
    if (eventRes.data) {
      events = eventRes.data;
      for (const e of events) initForms(e.id);
    }
    if (yearRes.data) {
      const current = yearRes.data.find((y) => y.is_current);
      if (current) newAcademicYearId = current.id;
    }
    if (studentRes.data) students = studentRes.data;
    loading = false;
  });

  function initForms(eventId: string) {
    if (!partForm[eventId]) {
      partForm[eventId] = {
        student_id: "",
        position: "participation",
        prize_title: "",
        issue_date: "",
      };
    }
    if (!signForm[eventId]) {
      signForm[eventId] = {
        name: "",
        role: "principal",
        title: "",
        signature_url: "",
        uploading: false,
        error: "",
      };
    }
  }

  async function createEvent() {
    if (!newName.trim()) return;
    savingEvent = true;
    error = "";
    const res = await api<{ id: string }>("POST", "/certificates/events", {
      name: newName.trim(),
      category: newCategory,
      held_date: newHeldDate || undefined,
      venue: newVenue.trim() || undefined,
      description: newDescription.trim() || undefined,
      academic_year_id: newAcademicYearId || undefined,
    });
    savingEvent = false;
    if (res.error) {
      error = res.error.message;
      return;
    }
    toast("Event created", "success");
    events = [
      ...events,
      {
        id: res.data!.id,
        school_id: "",
        name: newName.trim(),
        category: newCategory,
        held_date: newHeldDate || undefined,
        venue: newVenue.trim() || undefined,
        description: newDescription.trim() || undefined,
        academic_year_id: newAcademicYearId || undefined,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      },
    ];
    initForms(res.data!.id);
    newName = "";
    newHeldDate = "";
    newVenue = "";
    newDescription = "";
    showForm = false;
  }

  async function toggleEvent(id: string) {
    if (expandedEventId === id) {
      expandedEventId = null;
      return;
    }
    expandedEventId = id;
    loadingDetails = true;
    const res = await api<{
      event: CertificateEvent;
      participants: CertificateParticipant[];
      signatories: CertificateSignatory[];
    }>("GET", `/certificates/events/${id}`);
    if (res.data) {
      eventDetails[id] = {
        participants: res.data.participants,
        signatories: res.data.signatories,
      };
    }
    loadingDetails = false;
  }

  async function addParticipant(eventId: string) {
    const f = partForm[eventId];
    if (!f?.student_id) return;
    savingParticipant = true;
    error = "";
    const res = await api(
      "POST",
      `/certificates/events/${eventId}/participants`,
      {
        student_id: f.student_id,
        position: f.position,
        prize_title: f.prize_title || undefined,
        issue_date: f.issue_date || undefined,
      },
    );
    savingParticipant = false;
    if (res.error) {
      error = res.error.message;
      return;
    }
    toast("Participant added", "success");
    // refresh details
    const detail = await api<{
      event: CertificateEvent;
      participants: CertificateParticipant[];
      signatories: CertificateSignatory[];
    }>("GET", `/certificates/events/${eventId}`);
    if (detail.data)
      eventDetails[eventId] = {
        participants: detail.data.participants,
        signatories: detail.data.signatories,
      };
    f.student_id = "";
    f.prize_title = "";
    f.issue_date = "";
  }

  async function deleteParticipant(certId: string, eventId: string) {
    if (!confirm("Remove this participant?")) return;
    const res = await api("DELETE", `/certificates/${certId}`);
    if (res.error) {
      toast(res.error.message, "error");
      return;
    }
    eventDetails[eventId].participants = eventDetails[
      eventId
    ].participants.filter((p) => p.id !== certId);
    toast("Participant removed", "success");
  }

  async function uploadSignature(eventId: string) {
    const f = signForm[eventId];
    const input = document.getElementById(
      `sig-file-${eventId}`,
    ) as HTMLInputElement;
    if (!input?.files?.length) {
      f.error = "Please choose a signature image.";
      return;
    }
    f.uploading = true;
    f.error = "";
    const formData = new FormData();
    formData.append("file", input.files[0]);
    const json = await apiUpload<{ url: string }>(
      "/certificates/signatures",
      formData,
    );
    if (json.data?.url) {
      f.signature_url = json.data.url;
    } else {
      f.error = json.error?.message || "Upload failed";
    }
    f.uploading = false;
    input.value = "";
  }

  async function addSignatory(eventId: string) {
    const f = signForm[eventId];
    if (!f?.name.trim()) return;
    savingSignatory = true;
    error = "";
    const res = await api(
      "POST",
      `/certificates/events/${eventId}/signatories`,
      {
        name: f.name.trim(),
        role: f.role,
        title: f.title.trim() || undefined,
        signature_url: f.signature_url || undefined,
      },
    );
    savingSignatory = false;
    if (res.error) {
      error = res.error.message;
      return;
    }
    toast("Signatory added", "success");
    const detail = await api<{
      event: CertificateEvent;
      participants: CertificateParticipant[];
      signatories: CertificateSignatory[];
    }>("GET", `/certificates/events/${eventId}`);
    if (detail.data)
      eventDetails[eventId] = {
        participants: detail.data.participants,
        signatories: detail.data.signatories,
      };
    signForm[eventId] = {
      name: "",
      role: "principal",
      title: "",
      signature_url: "",
      uploading: false,
      error: "",
    };
  }

  async function deleteSignatory(sigId: string, eventId: string) {
    if (!confirm("Remove this signatory?")) return;
    const res = await api("DELETE", `/certificates/signatories/${sigId}`);
    if (res.error) {
      toast(res.error.message, "error");
      return;
    }
    eventDetails[eventId].signatories = eventDetails[
      eventId
    ].signatories.filter((s) => s.id !== sigId);
    toast("Signatory removed", "success");
  }

  function categoryLabel(cat: string): string {
    return CATEGORIES.find((c) => c.id === cat)?.name || cat;
  }
</script>

<div class="space-y-6">
  <PageHeader
    title="Certificates"
    subtitle="Generate premium certificates for competition participants"
  >
    {#snippet actions()}
      <Button
        icon={Plus}
        onclick={() => (showForm = true)}
      >
        Add Event
      </Button>
    {/snippet}
  </PageHeader>

  <Modal bind:open={showForm} title="Add Event" maxWidth="max-w-2xl">
    <div class="space-y-3">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <input
          bind:value={newName}
          placeholder="Event name (e.g. Kannada Elocution)"
          class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
        />
        <Select
          bind:value={newCategory}
          options={CATEGORIES}
          placeholder="Category"
        />
        <input
          bind:value={newHeldDate}
          type="date"
          class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
        />
        <input
          bind:value={newVenue}
          placeholder="Venue (optional)"
          class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
        />
        <input
          bind:value={newDescription}
          placeholder="Description (optional)"
          class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 sm:col-span-2"
        />
      </div>
      {#if error}
        <div class="text-sm text-danger-600">{error}</div>
      {/if}
    </div>
    {#snippet footer()}
      <Button variant="ghost" onclick={() => (showForm = false)}>Cancel</Button>
      <Button
        onclick={createEvent}
        disabled={savingEvent || !newName.trim()}
        loading={savingEvent}
      >
        {savingEvent ? "Saving..." : "Create Event"}
      </Button>
    {/snippet}
  </Modal>

  <div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
    <div class="p-4 border-b border-slate-200">
      <SearchFilter
        bind:value={search}
        placeholder="Search events by name, venue, or category..."
      />
    </div>
    <table class="w-full text-sm">
      <thead>
        <tr class="bg-slate-50 text-slate-600">
          <th class="text-left px-4 py-3 font-medium">Event</th>
          <th class="text-left px-4 py-3 font-medium">Category</th>
          <th class="text-left px-4 py-3 font-medium">Date</th>
          <th class="text-left px-4 py-3 font-medium">Venue</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          <tr
            ><td colspan="4" class="px-4 py-8 text-center text-slate-400"
              >Loading...</td
            ></tr
          >
        {:else if events.length === 0}
          <tr
            ><td colspan="4" class="px-4"
              ><EmptyState title="No events yet. Add one above." /></td
            ></tr
          >
        {:else if filteredEvents.length === 0}
          <tr
            ><td colspan="4" class="px-4"
              ><EmptyState
                title="No events match your search."
                hint="Try a different name, venue, or category."
              /></td
            ></tr
          >
        {:else}
          {#each filteredEvents as e (e.id)}
            <tr
              class="border-t border-slate-100 hover:bg-slate-50 cursor-pointer"
              onclick={() => toggleEvent(e.id)}
            >
              <td class="px-4 py-3">
                <div class="font-medium">{e.name}</div>
                <div class="text-xs text-slate-400 mt-0.5">
                  {expandedEventId === e.id
                    ? "Click to collapse"
                    : "Click to manage participants & signatures"}
                </div>
              </td>
              <td class="px-4 py-3">
                <span
                  class="px-2 py-0.5 rounded-full text-xs font-medium bg-primary-50 text-primary-700"
                  >{categoryLabel(e.category)}</span
                >
              </td>
              <td class="px-4 py-3 text-slate-600"
                >{e.held_date
                  ? new Date(e.held_date).toLocaleDateString()
                  : "—"}</td
              >
              <td class="px-4 py-3 text-slate-500">{e.venue || "—"}</td>
            </tr>
            {#if expandedEventId === e.id}
              <tr>
                <td colspan="4" class="px-4 py-4 bg-slate-50/60">
                  {#if loadingDetails && !eventDetails[e.id]}
                    <div class="text-sm text-slate-400">Loading...</div>
                  {:else if eventDetails[e.id]}
                    <div class="space-y-6">
                      <div>
                        <h3 class="text-sm font-semibold text-slate-800 mb-2">
                          Signatories (Principal & Judges)
                        </h3>
                        <div class="flex flex-wrap gap-3 mb-3">
                          {#each eventDetails[e.id].signatories as sig}
                            <div
                              class="flex items-center gap-2 bg-white border border-slate-200 rounded-lg px-3 py-2 text-sm"
                            >
                              <div>
                                <div class="font-medium text-slate-800">
                                  {sig.name}
                                </div>
                                <div class="text-xs text-slate-500">
                                  {SIGNATORY_ROLES.find(
                                    (r) => r.id === sig.role,
                                  )?.name || sig.role}{sig.title
                                    ? ` — ${sig.title}`
                                    : ""}
                                </div>
                              </div>
                              <Button
                                variant="ghost"
                                size="sm"
                                icon={Trash2}
                                onclick={() => deleteSignatory(sig.id, e.id)}
                                aria-label={`Remove signatory ${sig.name}`}
                              >Remove</Button>
                            </div>
                          {/each}
                        </div>

                        <div
                          class="bg-white border border-slate-200 rounded-lg p-3 space-y-2"
                        >
                          <div
                            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-2"
                          >
                            <input
                              bind:value={signForm[e.id].name}
                              placeholder="Signatory name"
                              class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                            />
                            <Select
                              bind:value={signForm[e.id].role}
                              options={SIGNATORY_ROLES}
                              placeholder="Role"
                            />
                            <input
                              bind:value={signForm[e.id].title}
                              placeholder="Title (e.g. Principal)"
                              class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                            />
                            <div>
                              <input
                                id="sig-file-{e.id}"
                                type="file"
                                accept="image/*"
                                class="text-xs w-full"
                                onchange={() => uploadSignature(e.id)}
                              />
                            </div>
                            <Button
                              variant="secondary"
                              onclick={() => addSignatory(e.id)}
                              disabled={savingSignatory ||
                                !signForm[e.id].name.trim()}
                              loading={signForm[e.id].uploading ||
                                savingSignatory}
                            >
                              {signForm[e.id].uploading
                                ? "Uploading..."
                                : savingSignatory
                                  ? "Saving..."
                                  : "Add Signatory"}
                            </Button>
                          </div>
                          {#if signForm[e.id].signature_url}
                            <div
                              class="flex items-center gap-2 text-xs text-green-600"
                            >
                              <span>✓ Signature uploaded</span>
                              <img
                                src={apiUrl(signForm[e.id].signature_url!)}
                                alt="signature preview"
                                class="h-8"
                              />
                            </div>
                          {/if}
                          {#if signForm[e.id].error}
                            <div class="text-xs text-danger-600">
                              {signForm[e.id].error}
                            </div>
                          {/if}
                        </div>
                      </div>

                      <div>
                        <h3 class="text-sm font-semibold text-slate-800 mb-2">
                          Participants ({eventDetails[e.id].participants
                            .length})
                        </h3>
                        <div
                          class="overflow-x-auto bg-white border border-slate-200 rounded-lg"
                        >
                          <table class="w-full text-sm">
                            <thead>
                              <tr class="bg-slate-50 text-slate-600">
                                <th class="text-left px-3 py-2 font-medium"
                                  >Student</th
                                >
                                <th class="text-left px-3 py-2 font-medium"
                                  >SATS</th
                                >
                                <th class="text-left px-3 py-2 font-medium"
                                  >Class</th
                                >
                                <th class="text-left px-3 py-2 font-medium"
                                  >Position</th
                                >
                                <th class="text-left px-3 py-2 font-medium"
                                  >Prize</th
                                >
                                <th class="text-right px-3 py-2 font-medium"
                                  >Actions</th
                                >
                              </tr>
                            </thead>
                            <tbody>
                              {#each eventDetails[e.id].participants as p (p.id)}
                                <tr class="border-t border-slate-100">
                                  <td class="px-3 py-2 font-medium"
                                    >{p.student_name}</td
                                  >
                                  <td class="px-3 py-2 text-slate-500"
                                    >{p.sats_number}</td
                                  >
                                  <td class="px-3 py-2 text-slate-500"
                                    >{p.class_name || "—"}</td
                                  >
                                  <td class="px-3 py-2"
                                    >{POSITIONS.find(
                                      (x) => x.id === p.position,
                                    )?.name.split(" ")[0] || p.position}</td
                                  >
                                  <td class="px-3 py-2 text-slate-600"
                                    >{p.prize_title || "—"}</td
                                  >
                                  <td class="px-3 py-2">
                                    <div class="flex justify-end gap-2">
                                      <Button
                                        size="sm"
                                        onclick={() =>
                                          goto(`/certificates/print/${p.id}`)}
                                        >Print</Button
                                      >
                                      <Button
                                        size="sm"
                                        variant="danger"
                                        onclick={() =>
                                          deleteParticipant(p.id, e.id)}
                                        >Remove</Button
                                      >
                                    </div>
                                  </td>
                                </tr>
                              {:else}
                                <tr
                                  ><td
                                    colspan="6"
                                    class="px-3 py-4 text-center text-slate-400"
                                    >No participants yet.</td
                                  ></tr
                                >
                              {/each}
                            </tbody>
                          </table>
                        </div>

                        <div
                          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-2 mt-3"
                        >
                          <Select
                            bind:value={partForm[e.id].student_id}
                            options={studentOptions}
                            placeholder="Select student"
                            searchable
                          />
                          <Select
                            bind:value={partForm[e.id].position}
                            options={POSITIONS}
                            placeholder="Position"
                          />
                          <input
                            bind:value={partForm[e.id].prize_title}
                            placeholder="Prize title (e.g. First Prize)"
                            class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                          />
                          <input
                            bind:value={partForm[e.id].issue_date}
                            type="date"
                            class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                          />
                          <Button
                            onclick={() => addParticipant(e.id)}
                            disabled={savingParticipant ||
                              !partForm[e.id].student_id}
                            loading={savingParticipant}
                          >
                            {savingParticipant
                              ? "Saving..."
                              : "Add Participant"}
                          </Button>
                        </div>
                      </div>
                    </div>
                  {/if}
                </td>
              </tr>
            {/if}
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
