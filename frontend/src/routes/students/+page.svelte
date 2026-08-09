<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { Plus, Pencil, Trash2, User, Hash, Phone, Mail, MapPin, Calendar, Droplets, Users, Eye, EyeOff } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import Select from '$lib/components/Select.svelte';
	import SearchFilter from '$lib/components/SearchFilter.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import type { Student, Class, AcademicYear } from '$lib/types';
	import { onMount } from 'svelte';

	let allStudents: Student[] = $state([]);
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
	let formFatherName = $state('');
	let formMotherName = $state('');
	let formParentName = $state('');
	let formParentPhone = $state('');
	let formParentEmail = $state('');

	// Search, filter, pagination
	let search = $state('');
	let filterClass = $state('');
	let page = $state(1);
	const pageSize = 20;

	let filteredStudents = $derived(allStudents.filter(s => {
		if (filterClass && s.class_id !== filterClass) return false;
		if (!search.trim()) return true;
		const q = search.toLowerCase();
		return s.first_name?.toLowerCase().includes(q) ||
			s.last_name?.toLowerCase().includes(q) ||
			s.sats_number?.includes(q);
	}));

	let paginatedStudents = $derived(filteredStudents.slice((page - 1) * pageSize, page * pageSize));
	let totalStudents = $derived(filteredStudents.length);

	function onPageChange(p: number) { page = p; }
	function resetPage() { page = 1; }

	onMount(async () => {
		const [sRes, cRes, yRes] = await Promise.all([
			api<Student[]>('GET', '/students?limit=500'),
			api<Class[]>('GET', '/classes?limit=50'),
			api<AcademicYear[]>('GET', '/academic-years?limit=50')
		]);
		if (sRes.data) allStudents = sRes.data;
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
		formFatherName = '';
		formMotherName = '';
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
		formFatherName = s.father_name ?? '';
		formMotherName = s.mother_name ?? '';
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
			father_name: formFatherName.trim() || undefined,
			mother_name: formMotherName.trim() || undefined,
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
				father_name: formFatherName.trim() || undefined,
				mother_name: formMotherName.trim() || undefined,
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

	function genderLabel(v: string): string {
		if (v === 'male') return 'Male';
		if (v === 'female') return 'Female';
		return '—';
	}

	function genderBadge(v: string): string {
		if (v === 'male') return 'bg-blue-100 text-blue-700';
		if (v === 'female') return 'bg-pink-100 text-pink-700';
		return 'bg-slate-100 text-slate-500';
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Students</h1>
			<p class="text-sm text-slate-500 mt-1">{students.length} student{students.length !== 1 ? 's' : ''} enrolled</p>
		</div>
		<Button onclick={openCreate} icon={Plus}>Add Student</Button>
	</div>

	{#if showForm}
		<div class="bg-white rounded-xl border border-slate-200 shadow-sm">
			<div class="px-4 py-3 border-b border-slate-100 flex items-center justify-between">
				<h3 class="text-sm font-semibold text-slate-700 flex items-center gap-2">
					<User size={16} class="text-primary-500" />
					{editingId ? 'Edit Student' : 'New Student'}
				</h3>
				<button onclick={cancelForm} class="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-md hover:bg-slate-100">
					<EyeOff size={16} />
				</button>
			</div>
			<div class="p-4 space-y-4">
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
					<div>
						<label for="st-sats" class="block text-xs font-medium text-slate-500 mb-1">SATS Number *</label>
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><Hash size={14} /></span>
							<input id="st-sats" bind:value={formSATS} placeholder="9-digit SATS number" maxlength={9} class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
					</div>
					<div>
						<label for="st-first" class="block text-xs font-medium text-slate-500 mb-1">First Name *</label>
						<input id="st-first" bind:value={formFirstName} placeholder="First name" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
					</div>
					<div>
						<label for="st-last" class="block text-xs font-medium text-slate-500 mb-1">Last Name</label>
						<input id="st-last" bind:value={formLastName} placeholder="Last name" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
					</div>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-4 gap-3">
					<div>
						<label for="st-roll" class="block text-xs font-medium text-slate-500 mb-1">Roll No</label>
						<input id="st-roll" bind:value={formRollNo} type="number" min="0" placeholder="Roll number" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
					</div>
					<div>
						<label for="st-gender" class="block text-xs font-medium text-slate-500 mb-1">Gender</label>
						<Select id="st-gender" bind:value={formGender} options={[{ id: 'male', name: 'Male' }, { id: 'female', name: 'Female' }]} placeholder="Select gender" icon={Users} />
					</div>
					<div>
						<label for="st-dob" class="block text-xs font-medium text-slate-500 mb-1">Date of Birth</label>
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><Calendar size={14} /></span>
							<input id="st-dob" bind:value={formDOB} type="date" class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
					</div>
					<div>
						<label for="st-blood" class="block text-xs font-medium text-slate-500 mb-1">Blood Group</label>
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><Droplets size={14} /></span>
							<input id="st-blood" bind:value={formBloodGroup} placeholder="e.g. O+" class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
					</div>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
					<div>
						<label for="st-phone" class="block text-xs font-medium text-slate-500 mb-1">Phone</label>
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><Phone size={14} /></span>
							<input id="st-phone" bind:value={formPhone} placeholder="Phone number" class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
					</div>
					<div>
						<label for="st-email" class="block text-xs font-medium text-slate-500 mb-1">Email</label>
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><Mail size={14} /></span>
							<input id="st-email" bind:value={formEmail} type="email" placeholder="Email address" class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
					</div>
					<div>
						<label for="st-address" class="block text-xs font-medium text-slate-500 mb-1">Address</label>
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><MapPin size={14} /></span>
							<input id="st-address" bind:value={formAddress} placeholder="Address" class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
					</div>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
					<Select bind:value={formClassId} label="Class *" options={classes} icon={Users} placeholder="Select class" />
					<Select bind:value={formAcademicYearId} label="Academic Year" options={academicYears} icon={Calendar} placeholder="Select academic year" />
				</div>

				<div class="border-t border-slate-100 pt-4">
					<p class="text-xs font-medium text-slate-500 mb-2 flex items-center gap-1.5"><Users size={12} /> Parents / Guardian</p>
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
						<input bind:value={formFatherName} placeholder="Father name" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						<input bind:value={formMotherName} placeholder="Mother name" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
					</div>
					<div class="grid grid-cols-1 sm:grid-cols-3 gap-3 mt-3">
						<input bind:value={formParentName} placeholder="Guardian name" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><Phone size={14} /></span>
							<input bind:value={formParentPhone} placeholder="Guardian phone" class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
						<div class="relative">
							<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-slate-400"><Mail size={14} /></span>
							<input bind:value={formParentEmail} type="email" placeholder="Guardian email" class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400" />
						</div>
					</div>
				</div>

				{#if error}
					<div class="flex items-center gap-2 text-sm px-3 py-2 rounded-lg bg-red-50 text-danger-600 border border-red-200">
						<span>{error}</span>
					</div>
				{/if}

				<div class="flex items-center gap-2 pt-1">
					<Button onclick={save} disabled={saving || !formSATS.trim() || !formFirstName.trim() || !formClassId} loading={saving}>
						{editingId ? 'Update Student' : 'Create Student'}
					</Button>
					<Button onclick={cancelForm} variant="secondary">Cancel</Button>
				</div>
			</div>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-sm">
		<div class="p-4 border-b border-slate-200 flex flex-wrap gap-3">
			<div class="flex-1 min-w-48">
				<SearchFilter bind:value={search} placeholder="Search by name or SATS..." onInput={resetPage} />
			</div>
			<div class="w-40">
				<Select bind:value={filterClass} options={[{ id: '', name: 'All Classes' }, ...classes.map(c => ({ id: c.id, name: c.name }))]} placeholder="All Classes" onselect={resetPage} />
			</div>
		</div>
		<table class="w-full text-sm">
			<thead>
				<tr class="bg-slate-50 text-slate-600">
					<th class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider">SATS</th>
					<th class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider">Name</th>
					<th class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider">Class</th>
					<th class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider">Roll No</th>
					<th class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider">Gender</th>
					<th class="text-left px-4 py-3 font-semibold text-xs uppercase tracking-wider">Year</th>
					<th class="text-right px-4 py-3 font-semibold text-xs uppercase tracking-wider">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr><td colspan="7" class="px-4 py-12 text-center text-slate-400">Loading...</td></tr>
				{:else if totalStudents === 0}
					<tr><td colspan="7" class="px-4 py-12 text-center text-slate-400">No students found.</td></tr>
				{:else}
					{#each paginatedStudents as s (s.id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50 transition-colors">
							<td class="px-4 py-3.5 font-mono text-xs text-slate-500">{s.sats_number}</td>
							<td class="px-4 py-3.5 font-medium">
								<div class="flex items-center gap-2.5">
									<div class="w-8 h-8 rounded-full bg-primary-100 text-primary-600 flex items-center justify-center text-xs font-bold uppercase shrink-0">
										{s.first_name[0]}{s.last_name?.[0] ?? ''}
									</div>
									<div>
										<div>{s.first_name} {s.last_name}</div>
										{#if s.phone}
											<div class="text-xs text-slate-400">{s.phone}</div>
										{/if}
									</div>
								</div>
							</td>
							<td class="px-4 py-3.5"><span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-slate-100 text-slate-600">{className(s.class_id)}</span></td>
							<td class="px-4 py-3.5 text-slate-500">{s.roll_no ?? '—'}</td>
							<td class="px-4 py-3.5">
								{#if s.gender}
									<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium {genderBadge(s.gender)}">{genderLabel(s.gender)}</span>
								{:else}
									<span class="text-slate-400">—</span>
								{/if}
							</td>
							<td class="px-4 py-3.5 text-xs text-slate-500">{yearName(s.academic_year_id)}</td>
							<td class="px-4 py-3.5 text-right">
								<div class="flex items-center justify-end gap-1">
									<Button onclick={() => openEdit(s)} variant="ghost" size="sm" icon={Pencil}>Edit</Button>
									<Button onclick={() => removeStudent(s.id)} variant="ghost" size="sm" icon={Trash2}>Delete</Button>
								</div>
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
		<Pagination {total} pageSize={pageSize} {page} onChange={onPageChange} />
	</div>
</div>
