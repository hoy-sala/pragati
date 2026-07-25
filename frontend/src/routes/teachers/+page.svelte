<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Class, Subject } from '$lib/types';
	import { onMount } from 'svelte';

	interface Teacher {
		id: string;
		school_id: string;
		email: string;
		name: string;
		role: string;
		phone: string;
		is_active: boolean;
		created_at: string;
		updated_at: string;
	}

	interface SubjectAssignment {
		teacher_id: string;
		subject_id: string;
		class_id: string;
		subject_name: string;
		class_name: string;
	}

	let teachers: Teacher[] = $state([]);
	let classes: Class[] = $state([]);
	let subjects: Subject[] = $state([]);
	let loading = $state(true);
	let showForm = $state(false);
	let saving = $state(false);
	let error = $state('');

	let newName = $state('');
	let newEmail = $state('');
	let newPassword = $state('');
	let newPhone = $state('');

	let editingTeacherId: string | null = $state(null);
	let editingSubjects: SubjectAssignment[] = $state([]);
	let savingSubjects = $state(false);

	onMount(async () => {
		const [tRes, cRes, sRes] = await Promise.all([
			api<Teacher[]>('GET', '/teachers?limit=100'),
			api<Class[]>('GET', '/classes?limit=50'),
			api<Subject[]>('GET', '/subjects?limit=100')
		]);
		if (tRes.data) teachers = tRes.data;
		if (cRes.data) classes = cRes.data;
		if (sRes.data) subjects = sRes.data;
		loading = false;
	});

	async function createTeacher() {
		if (!newName.trim() || !newEmail.trim() || !newPassword.trim()) return;
		saving = true;
		error = '';
		const res = await api<{ id: string }>('POST', '/teachers', {
			name: newName.trim(),
			email: newEmail.trim(),
			password: newPassword,
			phone: newPhone.trim() || undefined
		});
		saving = false;
		if (res.error) {
			error = res.error.message;
			return;
		}
		teachers = [...teachers, {
			id: res.data!.id,
			school_id: '',
			email: newEmail.trim(),
			name: newName.trim(),
			role: 'teacher',
			phone: newPhone.trim(),
			is_active: true,
			created_at: new Date().toISOString(),
			updated_at: new Date().toISOString()
		}];
		newName = '';
		newEmail = '';
		newPassword = '';
		newPhone = '';
		showForm = false;
	}

	async function toggleSubjects(teacher: Teacher) {
		if (editingTeacherId === teacher.id) {
			editingTeacherId = null;
			return;
		}
		editingTeacherId = teacher.id;
		const res = await api<SubjectAssignment[]>('GET', `/teachers/${teacher.id}/subjects`);
		editingSubjects = res.data || [];
	}

	function addSubjectRow() {
		editingSubjects = [...editingSubjects, {
			teacher_id: editingTeacherId!,
			subject_id: '',
			class_id: '',
			subject_name: '',
			class_name: ''
		}];
	}

	function removeSubjectRow(index: number) {
		editingSubjects = editingSubjects.filter((_, i) => i !== index);
	}

	async function saveSubjects() {
		savingSubjects = true;
		const valid = editingSubjects.filter(s => s.subject_id && s.class_id);
		const res = await api('PUT', `/teachers/${editingTeacherId}/subjects`, {
			subjects: valid.map(s => ({ subject_id: s.subject_id, class_id: s.class_id }))
		});
		savingSubjects = false;
		if (res.error) return;
		editingTeacherId = null;
	}

	function subjectName(id: string): string {
		return subjects.find(s => s.id === id)?.name || id;
	}

	function className(id: string): string {
		return classes.find(c => c.id === id)?.name || id;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Teachers</h1>
			<p class="text-sm text-slate-500 mt-1">{teachers.length} teachers</p>
		</div>
		<button onclick={() => showForm = !showForm} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			{showForm ? 'Cancel' : 'Add Teacher'}
		</button>
	</div>

	{#if showForm}
		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
			<div class="grid grid-cols-1 sm:grid-cols-4 gap-3">
				<input bind:value={newName} placeholder="Full name" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newEmail} type="email" placeholder="Email address" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newPassword} type="password" placeholder="Password" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
				<input bind:value={newPhone} placeholder="Phone (optional)" class="px-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
			</div>
			{#if error}
				<div class="text-sm text-danger-600">{error}</div>
			{/if}
			<button onclick={createTeacher} disabled={saving || !newName.trim() || !newEmail.trim() || !newPassword.trim()} class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{saving ? 'Saving...' : 'Create'}
			</button>
		</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		<table class="w-full text-sm">
			<thead>
				<tr class="bg-slate-50 text-slate-600">
					<th class="text-left px-4 py-3 font-medium">Name</th>
					<th class="text-left px-4 py-3 font-medium">Email</th>
					<th class="text-left px-4 py-3 font-medium">Phone</th>
					<th class="text-left px-4 py-3 font-medium">Subjects</th>
					<th class="text-left px-4 py-3 font-medium"></th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr><td colspan="5" class="px-4 py-8 text-center text-slate-400">Loading...</td></tr>
				{:else if teachers.length === 0}
					<tr><td colspan="5" class="px-4 py-8 text-center text-slate-400">No teachers yet. Add one above.</td></tr>
				{:else}
					{#each teachers as t (t.id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="px-4 py-3 font-medium">{t.name}</td>
							<td class="px-4 py-3 text-slate-500">{t.email}</td>
							<td class="px-4 py-3 text-slate-500">{t.phone || '—'}</td>
							<td class="px-4 py-3">
								<button onclick={() => toggleSubjects(t)} class="text-xs text-primary-600 hover:text-primary-700 underline underline-offset-2">
									{editingTeacherId === t.id ? 'Close' : 'Assign Subjects'}
								</button>
							</td>
							<td class="px-4 py-3">
								<span class="inline-block px-2 py-0.5 rounded-full text-xs {t.is_active ? 'bg-green-100 text-green-700' : 'bg-slate-100 text-slate-500'}">
									{t.is_active ? 'Active' : 'Inactive'}
								</span>
							</td>
						</tr>
						{#if editingTeacherId === t.id}
							<tr class="bg-slate-50">
								<td colspan="5" class="px-4 py-3">
									<div class="space-y-2">
										{#if editingSubjects.length === 0}
											<p class="text-xs text-slate-400">No subject assignments yet.</p>
										{:else}
											{#each editingSubjects as s, i (i)}
												<div class="flex items-center gap-2">
													<select bind:value={s.subject_id} class="flex-1 px-2 py-1.5 rounded border border-slate-300 text-xs focus:outline-none focus:ring-2 focus:ring-primary-500">
														<option value="">Select subject</option>
														{#each subjects as sub (sub.id)}
															<option value={sub.id}>{sub.name}</option>
														{/each}
													</select>
													<select bind:value={s.class_id} class="flex-1 px-2 py-1.5 rounded border border-slate-300 text-xs focus:outline-none focus:ring-2 focus:ring-primary-500">
														<option value="">Select class</option>
														{#each classes as c (c.id)}
															<option value={c.id}>{c.name}</option>
														{/each}
													</select>
													<button onclick={() => removeSubjectRow(i)} class="px-2 py-1 text-xs text-danger-600 hover:text-danger-700">Remove</button>
												</div>
											{/each}
										{/if}
										<div class="flex items-center gap-2 pt-1">
											<button onclick={addSubjectRow} class="px-2 py-1 text-xs bg-slate-200 text-slate-700 rounded hover:bg-slate-300 transition-colors">+ Add subject</button>
											<button onclick={saveSubjects} disabled={savingSubjects} class="px-3 py-1 text-xs bg-primary-600 text-white rounded hover:bg-primary-700 disabled:opacity-50 transition-colors">
												{savingSubjects ? 'Saving...' : 'Save'}
											</button>
										</div>
									</div>
								</td>
							</tr>
						{/if}
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>
