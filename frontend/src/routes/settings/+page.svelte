<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';
	import { Settings, Calendar, Users, BookOpen, ClipboardCheck, Plus, Star, Shield, UserPlus, Power, RotateCw } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import type { AcademicYear, Class, Subject, AssessmentCategory } from '$lib/types';

	type Tab = 'users' | 'years' | 'classes' | 'subjects' | 'categories';
	let activeTab = $state<Tab>('users');

	let years = $state<AcademicYear[]>([]);
	let classes = $state<Class[]>([]);
	let subjects = $state<Subject[]>([]);
	let categories = $state<AssessmentCategory[]>([]);
	let users = $state<{ id: string; email: string; name: string; role: string; phone: string; is_active: boolean }[]>([]);

	let showUserForm = $state(false);
	let userForm = $state({ email: '', password: '', name: '', role: 'teacher', phone: '' });
	let resettingId = $state<string | null>(null);
	let resetPassword = $state('');

	let loading = $state(true);
	let statusMsg = $state('');
	let statusType = $state<'info' | 'error' | 'success'>('info');

	let showYearForm = $state(false);
	let yearForm = $state({ name: '', start_date: '', end_date: '', is_current: false });

	let showClassForm = $state(false);
	let classForm = $state({ name: '', code: '', sort_order: 0 });

	let showSubjectForm = $state(false);
	let subjectForm = $state({ name: '', code: '', is_language: false, is_core: false });

	let showCatForm = $state(false);
	let catForm = $state({ name: '', code: '', weightage: 0, sort_order: 0 });

	const roleLabels: Record<string, string> = {
		admin: 'Administrator', principal: 'Principal', teacher: 'Teacher',
		special_educator: 'Special Educator', student: 'Student', parent: 'Parent',
	};

	async function loadUsers() {
		const res = await api<typeof users>('GET', '/users');
		if (res.data) users = res.data;
	}

	async function createUser() {
		if (!userForm.name || !userForm.role) { msg('Name and role are required', 'error'); return; }
		const res = await api<{ data: { id: string } }>('POST', '/users', userForm);
		if (res.data) {
			msg(`User "${userForm.name}" created`);
			showUserForm = false;
			userForm = { email: '', password: '', name: '', role: 'teacher', phone: '' };
			loadUsers();
		} else if (res.error) {
			msg(res.error.message, 'error');
		}
	}

	async function toggleUser(id: string) {
		const res = await api<unknown>('PATCH', `/users/${id}/toggle`);
		if (res.data) { msg('User status updated'); loadUsers(); }
		else if (res.error) msg(res.error.message, 'error');
	}

	async function resetPwd(id: string) {
		if (!resetPassword || resetPassword.length < 6) { msg('Password must be at least 6 characters', 'error'); return; }
		const res = await api<unknown>('POST', `/users/${id}/reset-password`, { password: resetPassword });
		if (res.data) { msg('Password reset'); resettingId = null; resetPassword = ''; }
		else if (res.error) msg(res.error.message, 'error');
	}

	const tabItems: { id: Tab; label: string; icon: typeof Settings }[] = [
		{ id: 'users', label: 'Users', icon: Shield },
		{ id: 'years', label: 'Academic Years', icon: Calendar },
		{ id: 'classes', label: 'Classes', icon: Users },
		{ id: 'subjects', label: 'Subjects', icon: BookOpen },
		{ id: 'categories', label: 'Categories', icon: ClipboardCheck },
	];

	function msg(text: string, type: 'success' | 'error' | 'info' = 'success') {
		statusMsg = text;
		statusType = type;
		setTimeout(() => { statusMsg = ''; }, 4000);
	}

	async function loadAll() {
		loading = true;
		const [yr, cl, sub, cat, uR] = await Promise.all([
			api<AcademicYear[]>('GET', '/academic-years'),
			api<Class[]>('GET', '/classes'),
			api<Subject[]>('GET', '/subjects'),
			api<AssessmentCategory[]>('GET', '/assessment-categories'),
			api<typeof users>('GET', '/users'),
		]);
		if (yr.data) years = yr.data;
		if (cl.data) classes = cl.data;
		if (sub.data) subjects = sub.data;
		if (cat.data) categories = cat.data;
		if (uR.data) users = uR.data;
		loading = false;
	}

	onMount(loadAll);

	async function createYear() {
		if (!yearForm.name || !yearForm.start_date || !yearForm.end_date) {
			msg('Fill all fields', 'error');
			return;
		}
		const res = await api<{ data: { id: string } }>('POST', '/academic-years', yearForm);
		if (res.data) {
			msg(`Academic year "${yearForm.name}" created`);
			showYearForm = false;
			yearForm = { name: '', start_date: '', end_date: '', is_current: false };
			loadAll();
		} else if (res.error) {
			msg(res.error.message, 'error');
		}
	}

	async function setCurrentYear(id: string) {
		const res = await api<unknown>('POST', `/academic-years/${id}/set-current`);
		if (res.data) {
			msg('Current academic year updated');
			loadAll();
		} else if (res.error) {
			msg(res.error.message, 'error');
		}
	}

	async function createClass() {
		if (!classForm.name) { msg('Name is required', 'error'); return; }
		const res = await api<{ data: { id: string } }>('POST', '/classes', classForm);
		if (res.data) {
			msg(`Class "${classForm.name}" created`);
			showClassForm = false;
			classForm = { name: '', code: '', sort_order: classes.length + 1 };
			loadAll();
		} else if (res.error) {
			msg(res.error.message, 'error');
		}
	}

	async function createSubject() {
		if (!subjectForm.name) { msg('Name is required', 'error'); return; }
		const res = await api<{ data: { id: string } }>('POST', '/subjects', subjectForm);
		if (res.data) {
			msg(`Subject "${subjectForm.name}" created`);
			showSubjectForm = false;
			subjectForm = { name: '', code: '', is_language: false, is_core: false };
			loadAll();
		} else if (res.error) {
			msg(res.error.message, 'error');
		}
	}

	async function createCategory() {
		if (!catForm.name) { msg('Name is required', 'error'); return; }
		const res = await api<{ data: { id: string } }>('POST', '/assessment-categories', catForm);
		if (res.data) {
			msg(`Category "${catForm.name}" created`);
			showCatForm = false;
			catForm = { name: '', code: '', weightage: 0, sort_order: categories.length + 1 };
			loadAll();
		} else if (res.error) {
			msg(res.error.message, 'error');
		}
	}
</script>

<svelte:head>
	<title>Settings — Pragati School Management</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<div class="flex items-center gap-3">
		<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-slate-600 to-slate-800 flex items-center justify-center">
			<Settings size={20} class="text-white" />
		</div>
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Settings</h1>
			<p class="text-sm text-slate-500">School configuration and preferences</p>
		</div>
	</div>

	{#if statusMsg}
		<div class="text-sm px-4 py-2 rounded-lg {statusType === 'success' ? 'bg-emerald-50 text-emerald-700' : statusType === 'error' ? 'bg-red-50 text-red-700' : 'bg-slate-100 text-slate-700'}">{statusMsg}</div>
	{/if}

	<div class="flex gap-1 bg-slate-100 rounded-lg p-1 no-print">
		{#each tabItems as item}
			<button onclick={() => activeTab = item.id}
				class="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab === item.id ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
				<item.icon size={16} />
				{item.label}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Loading...</div>

	{:else if activeTab === 'users'}
		<div class="bg-white rounded-xl border border-slate-200">
			<div class="px-6 py-4 border-b border-slate-200 flex items-center justify-between">
				<h2 class="text-base font-semibold text-slate-900">User Management</h2>
				<Button icon={UserPlus} onclick={() => showUserForm = !showUserForm}>
					{showUserForm ? 'Cancel' : 'New User'}
				</Button>
			</div>

			{#if showUserForm}
				<div class="px-6 py-4 bg-slate-50 border-b border-slate-200">
					<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
						<div>
							<label for="u-name" class="block text-xs font-medium text-slate-600 mb-1">Full Name *</label>
							<input id="u-name" bind:value={userForm.name} placeholder="John Doe" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="u-email" class="block text-xs font-medium text-slate-600 mb-1">Email</label>
							<input id="u-email" bind:value={userForm.email} placeholder="email@example.com" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="u-phone" class="block text-xs font-medium text-slate-600 mb-1">Phone</label>
							<input id="u-phone" bind:value={userForm.phone} placeholder="10-digit mobile" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="u-role" class="block text-xs font-medium text-slate-600 mb-1">Role *</label>
							<select id="u-role" bind:value={userForm.role} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
								<option value="teacher">Teacher</option>
								<option value="principal">Principal</option>
								<option value="special_educator">Special Educator</option>
								<option value="parent">Parent</option>
							</select>
						</div>
						<div>
							<label for="u-pass" class="block text-xs font-medium text-slate-600 mb-1">Password *</label>
							<input id="u-pass" type="password" bind:value={userForm.password} placeholder="Min 6 characters" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div class="flex items-end"><Button onclick={createUser}>Create User</Button></div>
					</div>
				</div>
			{/if}

			<div class="divide-y divide-slate-100">
				{#each users as u (u.id)}
					<div class="px-6 py-3 flex items-center justify-between hover:bg-slate-50/50">
						<div class="flex items-center gap-3">
							<div class="w-9 h-9 rounded-full bg-slate-100 flex items-center justify-center">
								<Users size={16} class="text-slate-500" />
							</div>
							<div>
								<div class="text-sm font-medium text-slate-800">{u.name}</div>
								<div class="text-xs text-slate-400">{u.email || '—'} {#if u.phone} · {u.phone}{/if}</div>
							</div>
						</div>
						<div class="flex items-center gap-2">
							<span class="text-xs px-2 py-0.5 rounded bg-slate-100 text-slate-600">{roleLabels[u.role] || u.role}</span>
							{#if u.is_active}
								<span class="text-xs px-2 py-0.5 rounded bg-emerald-50 text-emerald-600">Active</span>
							{:else}
								<span class="text-xs px-2 py-0.5 rounded bg-slate-100 text-slate-400">Inactive</span>
							{/if}
							{#if resettingId === u.id}
								<div class="flex items-center gap-2">
									<input type="password" bind:value={resetPassword} placeholder="New password" class="px-2 py-1 text-xs rounded border border-slate-300 w-32" />
									<Button size="sm" onclick={() => resetPwd(u.id)}>Save</Button>
									<Button size="sm" variant="ghost" onclick={() => { resettingId = null; resetPassword = ''; }}>Cancel</Button>
								</div>
							{:else}
								<button onclick={() => { resettingId = u.id; resetPassword = ''; }} title="Reset password" class="p-1.5 rounded text-slate-400 hover:text-slate-600 hover:bg-slate-100"><RotateCw size={14} /></button>
								<button onclick={() => toggleUser(u.id)} title={u.is_active ? 'Deactivate' : 'Activate'} class="p-1.5 rounded {u.is_active ? 'text-amber-500 hover:text-amber-600' : 'text-emerald-500 hover:text-emerald-600'} hover:bg-slate-100"><Power size={14} /></button>
							{/if}
						</div>
					</div>
				{:else}
					<div class="px-6 py-8 text-center text-sm text-slate-400">No users found</div>
				{/each}
			</div>
		</div>

	{:else if activeTab === 'years'}
		<div class="bg-white rounded-xl border border-slate-200">
			<div class="px-6 py-4 border-b border-slate-200 flex items-center justify-between">
				<h2 class="text-base font-semibold text-slate-900">Academic Years</h2>
				<Button icon={Plus} onclick={() => showYearForm = !showYearForm}>
					{showYearForm ? 'Cancel' : 'New Year'}
				</Button>
			</div>

			{#if showYearForm}
				<div class="px-6 py-4 bg-slate-50 border-b border-slate-200">
					<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
						<div>
							<label for="yr-name" class="block text-xs font-medium text-slate-600 mb-1">Name *</label>
							<input id="yr-name" bind:value={yearForm.name} placeholder="2026-27" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="yr-start" class="block text-xs font-medium text-slate-600 mb-1">Start Date *</label>
							<input id="yr-start" type="date" bind:value={yearForm.start_date} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="yr-end" class="block text-xs font-medium text-slate-600 mb-1">End Date *</label>
							<input id="yr-end" type="date" bind:value={yearForm.end_date} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div class="flex items-end gap-3">
							<label class="flex items-center gap-2 text-sm text-slate-700 pb-2">
								<input type="checkbox" bind:checked={yearForm.is_current} class="rounded" />
								Set as current
							</label>
							<Button onclick={createYear}>Create</Button>
						</div>
					</div>
				</div>
			{/if}

			<div class="divide-y divide-slate-100">
				{#each years as y (y.id)}
					<div class="px-6 py-3 flex items-center justify-between hover:bg-slate-50/50">
						<div class="flex items-center gap-3">
							{#if y.is_current}
								<span class="text-xs font-medium px-2 py-0.5 rounded bg-emerald-50 text-emerald-700 flex items-center gap-1"><Star size={10} /> Current</span>
							{:else}
								<button onclick={() => setCurrentYear(y.id)} class="text-xs font-medium px-2 py-0.5 rounded border border-slate-200 text-slate-500 hover:border-primary-300 hover:text-primary-600">
									Set Current
								</button>
							{/if}
							<span class="text-sm font-medium text-slate-800">{y.name}</span>
						</div>
						<span class="text-xs text-slate-400">{y.start_date} &mdash; {y.end_date}</span>
					</div>
				{:else}
					<div class="px-6 py-8 text-center text-sm text-slate-400">No academic years configured</div>
				{/each}
			</div>
		</div>

	{:else if activeTab === 'classes'}
		<div class="bg-white rounded-xl border border-slate-200">
			<div class="px-6 py-4 border-b border-slate-200 flex items-center justify-between">
				<h2 class="text-base font-semibold text-slate-900">Classes</h2>
				<Button icon={Plus} onclick={() => showClassForm = !showClassForm}>
					{showClassForm ? 'Cancel' : 'New Class'}
				</Button>
			</div>

			{#if showClassForm}
				<div class="px-6 py-4 bg-slate-50 border-b border-slate-200">
					<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
						<div>
							<label for="cl-name" class="block text-xs font-medium text-slate-600 mb-1">Name *</label>
							<input id="cl-name" bind:value={classForm.name} placeholder="Class 11" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="cl-code" class="block text-xs font-medium text-slate-600 mb-1">Code</label>
							<input id="cl-code" bind:value={classForm.code} placeholder="C11" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="cl-order" class="block text-xs font-medium text-slate-600 mb-1">Sort Order</label>
							<input id="cl-order" type="number" bind:value={classForm.sort_order} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
					</div>
					<div class="mt-3"><Button onclick={createClass}>Create</Button></div>
				</div>
			{/if}

			<div class="divide-y divide-slate-100">
				{#each classes as c (c.id)}
					<div class="px-6 py-3 flex items-center justify-between">
						<div class="flex items-center gap-3">
							<span class="text-sm font-medium text-slate-800">{c.name}</span>
							{#if c.code}<span class="text-xs text-slate-400">({c.code})</span>{/if}
						</div>
						<span class="text-xs text-slate-400">Order: {c.sort_order}</span>
					</div>
				{:else}
					<div class="px-6 py-8 text-center text-sm text-slate-400">No classes configured</div>
				{/each}
			</div>
		</div>

	{:else if activeTab === 'subjects'}
		<div class="bg-white rounded-xl border border-slate-200">
			<div class="px-6 py-4 border-b border-slate-200 flex items-center justify-between">
				<h2 class="text-base font-semibold text-slate-900">Subjects</h2>
				<Button icon={Plus} onclick={() => showSubjectForm = !showSubjectForm}>
					{showSubjectForm ? 'Cancel' : 'New Subject'}
				</Button>
			</div>

			{#if showSubjectForm}
				<div class="px-6 py-4 bg-slate-50 border-b border-slate-200">
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div>
							<label for="sb-name" class="block text-xs font-medium text-slate-600 mb-1">Name *</label>
							<input id="sb-name" bind:value={subjectForm.name} placeholder="Physics" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="sb-code" class="block text-xs font-medium text-slate-600 mb-1">Code</label>
							<input id="sb-code" bind:value={subjectForm.code} placeholder="PHY" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
					</div>
					<div class="mt-3 flex items-center gap-4">
						<label class="flex items-center gap-2 text-sm text-slate-700">
							<input type="checkbox" bind:checked={subjectForm.is_language} class="rounded" />
							Language
						</label>
						<label class="flex items-center gap-2 text-sm text-slate-700">
							<input type="checkbox" bind:checked={subjectForm.is_core} class="rounded" />
							Core Subject
						</label>
						<Button onclick={createSubject}>Create</Button>
					</div>
				</div>
			{/if}

			<div class="divide-y divide-slate-100">
				{#each subjects as s (s.id)}
					<div class="px-6 py-3 flex items-center justify-between">
						<div class="flex items-center gap-3">
							<span class="text-sm font-medium text-slate-800">{s.name}</span>
							{#if s.code}<span class="text-xs text-slate-400">({s.code})</span>{/if}
						</div>
						<div class="flex items-center gap-2">
							{#if s.is_language}
								<span class="text-xs px-2 py-0.5 rounded bg-blue-50 text-blue-600">Language</span>
							{/if}
							{#if s.is_core}
								<span class="text-xs px-2 py-0.5 rounded bg-purple-50 text-purple-600">Core</span>
							{/if}
						</div>
					</div>
				{:else}
					<div class="px-6 py-8 text-center text-sm text-slate-400">No subjects configured</div>
				{/each}
			</div>
		</div>

	{:else if activeTab === 'categories'}
		<div class="bg-white rounded-xl border border-slate-200">
			<div class="px-6 py-4 border-b border-slate-200 flex items-center justify-between">
				<h2 class="text-base font-semibold text-slate-900">Assessment Categories</h2>
				<Button icon={Plus} onclick={() => showCatForm = !showCatForm}>
					{showCatForm ? 'Cancel' : 'New Category'}
				</Button>
			</div>

			{#if showCatForm}
				<div class="px-6 py-4 bg-slate-50 border-b border-slate-200">
					<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
						<div>
							<label for="ct-name" class="block text-xs font-medium text-slate-600 mb-1">Name *</label>
							<input id="ct-name" bind:value={catForm.name} placeholder="Formative Assessment" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="ct-code" class="block text-xs font-medium text-slate-600 mb-1">Code</label>
							<input id="ct-code" bind:value={catForm.code} placeholder="FA" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="ct-weight" class="block text-xs font-medium text-slate-600 mb-1">Weightage %</label>
							<input id="ct-weight" type="number" bind:value={catForm.weightage} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
						<div>
							<label for="ct-order" class="block text-xs font-medium text-slate-600 mb-1">Sort Order</label>
							<input id="ct-order" type="number" bind:value={catForm.sort_order} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" />
						</div>
					</div>
					<div class="mt-3"><Button onclick={createCategory}>Create</Button></div>
				</div>
			{/if}

			<div class="divide-y divide-slate-100">
				{#each categories as cat (cat.id)}
					<div class="px-6 py-3 flex items-center justify-between">
						<div class="flex items-center gap-3">
							<span class="text-sm font-medium text-slate-800">{cat.name}</span>
							{#if cat.code}<span class="text-xs text-slate-400">({cat.code})</span>{/if}
						</div>
						<div class="flex items-center gap-3 text-xs text-slate-500">
							<span>Weight: {cat.weightage}%</span>
							<span>Order: {cat.sort_order}</span>
						</div>
					</div>
				{:else}
					<div class="px-6 py-8 text-center text-sm text-slate-400">No categories configured</div>
				{/each}
			</div>
		</div>
	{/if}
</div>
