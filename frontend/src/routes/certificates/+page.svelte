<script lang="ts">
	import { api, apiUrl } from '$lib/api/client.svelte';
	import type { CertificateEvent, CertificateParticipant, CertificateSignatory, Student, AcademicYear } from '$lib/types';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Select from '$lib/components/Select.svelte';

	const CATEGORIES = [
		{ id: 'sports', name: 'Sports' },
		{ id: 'cultural', name: 'Cultural' },
		{ id: 'academic', name: 'Academic' },
		{ id: 'other', name: 'Other' },
	];

	const POSITIONS = [
		{ id: '1st', name: '1st (First Prize)' },
		{ id: '2nd', name: '2nd (Runner Up)' },
		{ id: '3rd', name: '3rd (Consolation)' },
		{ id: 'participation', name: 'Participation' },
	];

	const SIGNATORY_ROLES = [
		{ id: 'principal', name: 'Principal' },
		{ id: 'chief_guest', name: 'Chief Guest' },
		{ id: 'chief_judge', name: 'Chief Judge' },
		{ id: 'judge', name: 'Judge' },
		{ id: 'coordinator', name: 'Event Coordinator' },
	];

	let events: CertificateEvent[] = $state([]);
	let academicYears: AcademicYear[] = $state([]);
	let students: Student[] = $state([]);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	let showForm = $state(false);
	let newName = $state('');
	let newCategory = $state('sports');
	let newHeldDate = $state('');
	let newVenue = $state('');
	let newDescription = $state('');
	let newAcademicYearId = $state('');

	let expandedEventId = $state<string | null>(null);
	let eventDetails = $state<Record<string, { participants: CertificateParticipant[]; signatories: CertificateSignatory[] }>>({});
	let loadingDetails = $state(false);

	let studentOptions = $derived(
		students.map(s => ({ id: s.id, name: `${s.first_name} ${s.last_name || ''}`.trim() + (s.sats_number ? ` (${s.sats_number})` : '') }))
	);

	let academicYearOptions = $derived(
		academicYears.map(y => ({ id: y.id, name: y.name + (y.is_current ? ' (current)' : '') }))
	);

	let partForm: Record<string, { student_id: string; position: string; prize_title: string; issue_date: string }> = $state({});

	let signForm: Record<string, { name: string; role: string; title: string; signature_url: string; uploading: boolean; error: string }> = $state({});

	onMount(async () => {
		const [eventRes, yearRes, studentRes] = await Promise.all([
			api<CertificateEvent[]>('GET', '/certificates/events?limit=100'),
			api<AcademicYear[]>('GET', '/academic-years?limit=50'),
			api<Student[]>('GET', '/students?limit=500')
		]);
		if (eventRes.data) {
			events = eventRes.data;
			for (const e of events) initForms(e.id);
		}
		if (yearRes.data) {
			academicYears = yearRes.data;
			const current = yearRes.data.find(y => y.is_current);
			if (current) newAcademicYearId = current.id;
		}
		if (studentRes.data) students = studentRes.data;
		loading = false;
	});

	function initForms(eventId: string) {
		if (!partForm[eventId]) {
			partForm[eventId] = { student_id: '', position: 'participation', prize_title: '', issue_date: '' };
		}
		if (!signForm[eventId]) {
			signForm[eventId] = { name: '', role: 'principal', title: '', signature_url: '', uploading: false, error: '' };
		}
	}

	async function createEvent() {
		if (!newName.trim()) return;
		saving = true;
		error = '';
		const res = await api<{ id: string }>('POST', '/certificates/events', {
			name: newName.trim(),
			category: newCategory,
			held_date: newHeldDate || undefined,
			venue: newVenue.trim() || undefined,
			description: newDescription.trim() || undefined,
			academic_year_id: newAcademicYearId || undefined
		});
		saving = false;
		if (res.error) {
			error = res.error.message;
			return;
		}
		events = [...events, {
			id: res.data!.id,
			school_id: '',
			name: newName.trim(),
			category: newCategory,
			held_date: newHeldDate || undefined,
			venue: newVenue.trim() || undefined,
			description: newDescription.trim() || undefined,
			academic_year_id: newAcademicYearId || undefined,
			created_at: new Date().toISOString(),
			updated_at: new Date().toISOString()
		}];
		initForms(res.data!.id);
		newName = '';
		newHeldDate = '';
		newVenue = '';
		newDescription = '';
		showForm = false;
	}

	async function toggleEvent(id: string) {
		if (expandedEventId === id) {
			expandedEventId = null;
			return;
		}
		expandedEventId = id;
		loadingDetails = true;
		const res = await api<{ event: CertificateEvent; participants: CertificateParticipant[]; signatories: CertificateSignatory[] }>('GET', `/certificates/events/${id}`);
		if (res.data) {
			eventDetails[id] = { participants: res.data.participants, signatories: res.data.signatories };
		}
		loadingDetails = false;
	}

	async function addParticipant(eventId: string) {
		const f = partForm[eventId];
		if (!f?.student_id) return;
		saving = true;
		error = '';
		const res = await api('POST', `/certificates/events/${eventId}/participants`, {
			student_id: f.student_id,
			position: f.position,
			prize_title: f.prize_title || undefined,
			issue_date: f.issue_date || undefined
		});
		saving = false;
		if (res.error) {
			error = res.error.message;
			return;
		}
		// refresh details
		const detail = await api<{ event: CertificateEvent; participants: CertificateParticipant[]; signatories: CertificateSignatory[] }>('GET', `/certificates/events/${eventId}`);
		if (detail.data) eventDetails[eventId] = { participants: detail.data.participants, signatories: detail.data.signatories };
		f.student_id = '';
		f.prize_title = '';
		f.issue_date = '';
	}

	async function deleteParticipant(certId: string, eventId: string) {
		if (!confirm('Remove this participant?')) return;
		const res = await api('DELETE', `/certificates/${certId}`);
		if (!res.error) {
			eventDetails[eventId].participants = eventDetails[eventId].participants.filter(p => p.id !== certId);
		}
	}

	async function uploadSignature(eventId: string) {
		const f = signForm[eventId];
		const input = document.getElementById(`sig-file-${eventId}`) as HTMLInputElement;
		if (!input?.files?.length) {
			f.error = 'Please choose a signature image.';
			return;
		}
		f.uploading = true;
		f.error = '';
		const formData = new FormData();
		formData.append('file', input.files[0]);
		try {
			const token = localStorage.getItem('access_token');
			const res = await fetch(apiUrl('/certificates/signatures'), {
				method: 'POST',
				headers: { 'Authorization': 'Bearer ' + token },
				body: formData
			});
			const json = await res.json();
			if (json.data?.url) {
				f.signature_url = json.data.url;
			} else {
				f.error = json.error?.message || 'Upload failed';
			}
		} catch {
			f.error = 'Unable to reach server.';
		}
		f.uploading = false;
		input.value = '';
	}

	async function addSignatory(eventId: string) {
		const f = signForm[eventId];
		if (!f?.name.trim()) return;
		saving = true;
		error = '';
		const res = await api('POST', `/certificates/events/${eventId}/signatories`, {
			name: f.name.trim(),
			role: f.role,
			title: f.title.trim() || undefined,
			signature_url: f.signature_url || undefined
		});
		saving = false;
		if (res.error) {
			error = res.error.message;
			return;
		}
		const detail = await api<{ event: CertificateEvent; participants: CertificateParticipant[]; signatories: CertificateSignatory[] }>('GET', `/certificates/events/${eventId}`);
		if (detail.data) eventDetails[eventId] = { participants: detail.data.participants, signatories: detail.data.signatories };
		signForm[eventId] = { name: '', role: 'principal', title: '', signature_url: '', uploading: false, error: '' };
	}

	async function deleteSignatory(sigId: string, eventId: string) {
		if (!confirm('Remove this signatory?')) return;
		const res = await api('DELETE', `/certificates/signatories/${sigId}`);
		if (!res.error) {
			eventDetails[eventId].signatories = eventDetails[eventId].signatories.filter(s => s.id !== sigId);
		}
	}

	function categoryLabel(cat: string): string {
		return CATEGORIES.find(c => c.id === cat)?.name || cat;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Certificates</h1>
			<p class="text-sm text-slate-500 mt-1">Generate premium certificates for competition participants</p>
		</div>
		<button onclick={() => showForm = !showForm} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			{showForm ? 'Cancel' : 'Add Event'}
		</button>
	</div>

	{#if showForm}
		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
				<input bind:value={newName} placeholder="Event name (e.g. Kannada Elocution)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<Select bind:value={newCategory} options={CATEGORIES} placeholder="Category" />
				<Select bind:value={newAcademicYearId} options={academicYearOptions} placeholder="Academic year" />
				<input bind:value={newHeldDate} type="date" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newVenue} placeholder="Venue (optional)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newDescription} placeholder="Description (optional)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
			</div>
			{#if error}
				<div class="text-sm text-danger-600">{error}</div>
			{/if}
			<button onclick={createEvent} disabled={saving || !newName.trim()} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{saving ? 'Saving...' : 'Create Event'}
			</button>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
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
					<tr><td colspan="4" class="px-4 py-8 text-center text-slate-400">Loading...</td></tr>
				{:else if events.length === 0}
					<tr><td colspan="4" class="px-4 py-8 text-center text-slate-400">No events yet. Add one above.</td></tr>
				{:else}
					{#each events as e (e.id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50 cursor-pointer" onclick={() => toggleEvent(e.id)}>
							<td class="px-4 py-3">
								<div class="font-medium">{e.name}</div>
								<div class="text-xs text-slate-400 mt-0.5">{expandedEventId === e.id ? 'Click to collapse' : 'Click to manage participants & signatures'}</div>
							</td>
							<td class="px-4 py-3">
								<span class="px-2 py-0.5 rounded-full text-xs font-medium bg-primary-50 text-primary-700">{categoryLabel(e.category)}</span>
							</td>
							<td class="px-4 py-3 text-slate-600">{e.held_date ? new Date(e.held_date).toLocaleDateString() : '—'}</td>
							<td class="px-4 py-3 text-slate-500">{e.venue || '—'}</td>
						</tr>
						{#if expandedEventId === e.id}
							<tr>
								<td colspan="4" class="px-4 py-4 bg-slate-50/60">
									{#if loadingDetails && !eventDetails[e.id]}
										<div class="text-sm text-slate-400">Loading...</div>
									{:else if eventDetails[e.id]}
										<div class="space-y-6">
											<div>
												<h3 class="text-sm font-semibold text-slate-800 mb-2">Signatories (Principal & Judges)</h3>
												<div class="flex flex-wrap gap-3 mb-3">
													{#each eventDetails[e.id].signatories as sig}
														<div class="flex items-center gap-2 bg-white border border-slate-200 rounded-lg px-3 py-2 text-sm">
															<div>
																<div class="font-medium text-slate-800">{sig.name}</div>
																<div class="text-xs text-slate-500">{SIGNATORY_ROLES.find(r => r.id === sig.role)?.name || sig.role}{sig.title ? ` — ${sig.title}` : ''}</div>
															</div>
															<button onclick={() => deleteSignatory(sig.id, e.id)} class="text-slate-400 hover:text-danger-600 transition-colors ml-1">✕</button>
														</div>
													{/each}
												</div>

												<div class="bg-white border border-slate-200 rounded-lg p-3 space-y-2">
												<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-2">
													<input bind:value={signForm[e.id].name} placeholder="Signatory name" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
													<Select bind:value={signForm[e.id].role} options={SIGNATORY_ROLES} placeholder="Role" />
													<input bind:value={signForm[e.id].title} placeholder="Title (e.g. Principal)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
													<div>
														<input id="sig-file-{e.id}" type="file" accept="image/*" class="text-xs w-full" onchange={() => uploadSignature(e.id)} />
													</div>
													<button onclick={() => addSignatory(e.id)} disabled={saving || !signForm[e.id].name.trim()}
														class="px-3 py-2 bg-slate-800 text-white rounded-lg text-sm font-medium hover:bg-slate-900 disabled:opacity-50 transition-colors">
														{signForm[e.id].uploading ? 'Uploading...' : 'Add Signatory'}
													</button>
												</div>
												{#if signForm[e.id].signature_url}
													<div class="flex items-center gap-2 text-xs text-green-600">
														<span>✓ Signature uploaded</span>
														<img src={apiUrl(signForm[e.id].signature_url!)} alt="signature preview" class="h-8" />
													</div>
												{/if}
												{#if signForm[e.id].error}
													<div class="text-xs text-danger-600">{signForm[e.id].error}</div>
												{/if}
												</div>
											</div>

											<div>
												<h3 class="text-sm font-semibold text-slate-800 mb-2">Participants ({eventDetails[e.id].participants.length})</h3>
												<div class="overflow-x-auto bg-white border border-slate-200 rounded-lg">
													<table class="w-full text-sm">
														<thead>
															<tr class="bg-slate-50 text-slate-600">
																<th class="text-left px-3 py-2 font-medium">Student</th>
																<th class="text-left px-3 py-2 font-medium">SATS</th>
																<th class="text-left px-3 py-2 font-medium">Class</th>
																<th class="text-left px-3 py-2 font-medium">Position</th>
																<th class="text-left px-3 py-2 font-medium">Prize</th>
																<th class="text-right px-3 py-2 font-medium">Actions</th>
															</tr>
														</thead>
														<tbody>
															{#each eventDetails[e.id].participants as p (p.id)}
																<tr class="border-t border-slate-100">
																	<td class="px-3 py-2 font-medium">{p.student_name}</td>
																	<td class="px-3 py-2 text-slate-500">{p.sats_number}</td>
																	<td class="px-3 py-2 text-slate-500">{p.class_name || '—'}</td>
																	<td class="px-3 py-2">{POSITIONS.find(x => x.id === p.position)?.name.split(' ')[0] || p.position}</td>
																	<td class="px-3 py-2 text-slate-600">{p.prize_title || '—'}</td>
																	<td class="px-3 py-2">
																		<div class="flex justify-end gap-2">
																			<button onclick={() => goto(`/certificates/print/${p.id}`)}
																				class="px-2 py-1 bg-primary-600 text-white rounded text-xs font-medium hover:bg-primary-700 transition-colors">Print</button>
																			<button onclick={() => deleteParticipant(p.id, e.id)} class="px-2 py-1 border border-slate-200 text-slate-500 rounded text-xs hover:text-danger-600 transition-colors">Remove</button>
																		</div>
																	</td>
																</tr>
															{:else}
																<tr><td colspan="6" class="px-3 py-4 text-center text-slate-400">No participants yet.</td></tr>
															{/each}
														</tbody>
													</table>
												</div>

												<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-2 mt-3">
													<Select bind:value={partForm[e.id].student_id} options={studentOptions} placeholder="Select student" searchable />
													<Select bind:value={partForm[e.id].position} options={POSITIONS} placeholder="Position" />
													<input bind:value={partForm[e.id].prize_title} placeholder="Prize title (e.g. First Prize)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
													<input bind:value={partForm[e.id].issue_date} type="date" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
													<button onclick={() => addParticipant(e.id)} disabled={saving || !partForm[e.id].student_id}
														class="px-3 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
														{saving ? 'Saving...' : 'Add Participant'}
													</button>
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