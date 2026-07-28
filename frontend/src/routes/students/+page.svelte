<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Student, Class, AcademicYear, Pagination } from '$lib/types';
	import { onMount } from 'svelte';

	let students: Student[] = $state([]);
	let classes: Class[] = $state([]);
	let academicYears: AcademicYear[] = $state([]);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	let showForm = $state(false);
	let editingId: string | null = $state(null);
	let formSATS = $state('');
	let formFirstName = $state('');
	let formLastName = $state('');
	let formRollNo = $state(0);
	let formGender = $state('');
	let formDOB = $state('');
	let formBloodGroup = $state('');
	let formPhone = $state('');
	let formEmail = $state('');
	let formAddress = $state('');
	let formClassId = $state('');
	let formAcademicYearId = $state('');
	let formParentName = $state('');
	let formParentPhone = $state('');
	let formParentEmail = $state('');

	onMount(async () => {
		const [sRes, cRes, yRes] = await Promise.all([
			api<Student[]>('GET', '/students?limit=200'),
			api<Class[]>('GET', '/classes?limit=50'),
			api<AcademicYear[]>('GET', '/academic-years?limit=50')
		]);
		if (sRes.data) students = sRes.data;
		if (cRes.data) classes = cRes.data;
		if (yRes.data) {
			academicYears = yRes.data;
			const current = yRes.data.find(y => y.is_current);
			if (current) formAcademicYearId = current.id;
		}
		loading = false;
	});

	function resetForm() {
		formSATS = '';
		formFirstName = '';
		formLastName = '';
		formRollNo = 0;
		formGender = '';
		formDOB = '';
		formBloodGroup = '';
		formPhone = '';
		formEmail = '';
		formAddress = '';
		formClassId = '';
		formParentName = '';
		formParentPhone = '';
		formParentEmail = '';
		if (academicYears.length > 0) {
			const current = academicYears.find(y => y.is_current);
			formAcademicYearId = current?.id ?? academicYears[0].id;
		}
		error = '';
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
		formLastName = s.last_name ?? '';
		formRollNo = s.roll_no ?? 0;
		formGender = s.gender ?? '';
		formDOB = s.date_of_birth ? s.date_of_birth.slice(0, 10) : '';
		formBloodGroup = s.blood_group ?? '';
		formPhone = s.phone ?? '';
		formEmail = s.email ?? '';
		formAddress = s.address ?? '';
		formClassId = s.class_id;
		formAcademicYearId = s.academic_year_id;
		formParentName = s.parent_name ?? '';
		formParentPhone = s.parent_phone ?? '';
		formParentEmail = s.parent_email ?? '';
		error = '';
		showForm = true;
	}

	function cancelForm() {
		showForm = false;
		editingId = null;
		resetForm();
	}

	async function save() {
		if (!formSATS.trim() || !formFirstName.trim() || !formClassId) return;
		if (formSATS.trim().length !== 9) { error = 'SATS number must be exactly 9 characters'; return; }
		saving = true;
		error = '';

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
			parent_name: formParentName.trim() || undefined,
			parent_phone: formParentPhone.trim() || undefined,
			parent_email: formParentEmail.trim() || undefined,
		};

		if (editingId) {
			const res = await api<Student>('PUT', `/students/${editingId}`, body);
			saving = false;
			if (res.error) { error = res.error.message; return; }
			students = students.map(s => s.id === editingId ? { ...s, ...body } : s);
		} else {
			const res = await api<{ id: string }>('POST', '/students', body);
			saving = false;
			if (res.error) { error = res.error.message; return; }
			students = [...students, {
				id: res.data!.id,
				school_id: '',
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
				parent_name: formParentName.trim() || undefined,
				parent_phone: formParentPhone.trim() || undefined,
				parent_email: formParentEmail.trim() || undefined,
				is_active: true,
				created_at: new Date().toISOString(),
				updated_at: new Date().toISOString(),
			} as Student];
		}
		showForm = false;
		editingId = null;
		resetForm();
	}

	async function removeStudent(id: string) {
		if (!confirm('Delete this student? This action cannot be undone.')) return;
		const res = await api('DELETE', `/students/${id}`);
		if (res.error) { error = res.error.message; return; }
		students = students.filter(s => s.id !== id);
	}

	function className(id: string): string {
		return classes.find(c => c.id === id)?.name ?? id;
	}

	function yearName(id: string): string {
		return academicYears.find(y => y.id === id)?.name ?? id;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Students</h1>
			<p class="text-sm text-slate-500 mt-1">{students.length} students enrolled</p>
		</div>
		<button onclick={openCreate} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			Add Student
		</button>
	</div>

	{#if showForm}
		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
			<h3 class="text-sm font-semibold text-slate-700">{editingId ? 'Edit Student' : 'New Student'}</h3>
			<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
				<input bind:value={formSATS} placeholder="SATS number (9 digits)" maxlength={9} class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formFirstName} placeholder="First name *" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formLastName} placeholder="Last name" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formRollNo} type="number" min="0" placeholder="Roll no" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<select bind:value={formGender} class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500">
					<option value="">Gender</option>
					<option value="male">Male</option>
					<option value="female">Female</option>
				</select>
				<input bind:value={formDOB} type="date" placeholder="Date of birth" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formBloodGroup} placeholder="Blood group" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formPhone} placeholder="Phone" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formEmail} type="email" placeholder="Email" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
			</div>
			<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
				<select bind:value={formClassId} class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500">
					<option value="">Select class *</option>
					{#each classes as c (c.id)}
						<option value={c.id}>{c.name}</option>
					{/each}
				</select>
				<select bind:value={formAcademicYearId} class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500">
					<option value="">Academic year</option>
					{#each academicYears as y (y.id)}
						<option value={y.id}>{y.name}</option>
					{/each}
				</select>
			</div>
			<input bind:value={formAddress} placeholder="Address" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
			<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
				<input bind:value={formParentName} placeholder="Parent name" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formParentPhone} placeholder="Parent phone" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={formParentEmail} type="email" placeholder="Parent email" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
			</div>
			{#if error}
				<div class="text-sm text-danger-600">{error}</div>
			{/if}
			<div class="flex items-center gap-2">
				<button onclick={save} disabled={saving || !formSATS.trim() || !formFirstName.trim() || !formClassId} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
					{saving ? 'Saving...' : editingId ? 'Update' : 'Create'}
				</button>
				<button onclick={cancelForm} class="px-4 py-2 bg-slate-200 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-300 transition-colors">Cancel</button>
			</div>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		<table class="w-full text-sm">
			<thead>
				<tr class="bg-slate-50 text-slate-600">
					<th class="text-left px-4 py-3 font-medium">SATS</th>
					<th class="text-left px-4 py-3 font-medium">Name</th>
					<th class="text-left px-4 py-3 font-medium">Class</th>
					<th class="text-left px-4 py-3 font-medium">Roll No</th>
					<th class="text-left px-4 py-3 font-medium">Gender</th>
					<th class="text-left px-4 py-3 font-medium">Year</th>
					<th class="text-left px-4 py-3 font-medium"></th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr><td colspan="7" class="px-4 py-8 text-center text-slate-400">Loading...</td></tr>
				{:else if students.length === 0}
					<tr><td colspan="7" class="px-4 py-8 text-center text-slate-400">No students yet. Add one above.</td></tr>
				{:else}
					{#each students as s (s.id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="px-4 py-3 font-mono text-xs">{s.sats_number}</td>
							<td class="px-4 py-3 font-medium">{s.first_name} {s.last_name}</td>
							<td class="px-4 py-3 text-slate-500">{className(s.class_id)}</td>
							<td class="px-4 py-3 text-slate-500">{s.roll_no ?? '—'}</td>
							<td class="px-4 py-3 text-slate-500">{s.gender ? (s.gender === 'male' ? 'M' : 'F') : '—'}</td>
							<td class="px-4 py-3 text-slate-500 text-xs">{yearName(s.academic_year_id)}</td>
							<td class="px-4 py-3">
								<div class="flex items-center gap-2">
									<button onclick={() => openEdit(s)} class="text-xs text-primary-600 hover:text-primary-700 underline underline-offset-2">Edit</button>
									<button onclick={() => removeStudent(s.id)} class="text-xs text-danger-600 hover:text-danger-700 underline underline-offset-2">Delete</button>
								</div>
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>