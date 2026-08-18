<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';
	import { Users, UserPlus, Trash2, Search, UserRound, GraduationCap, AlertCircle } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import Select from '$lib/components/Select.svelte';
	import type { Class, AcademicYear } from '$lib/types';

	let classes = $state<Class[]>([]);
	let years = $state<AcademicYear[]>([]);
	let selectedYear = $state('');
	let selectedClass = $state('');

	let assignments = $state<{ id: string; mentor_id: string; mentor_name: string; student_id: string; student_name: string; sats_number: string; class_name: string; academic_year_id: string }[]>([]);
	let stats = $state<{ unassigned: number; mentors: { id: string; name: string; student_count: number }[] } | null>(null);
	let unassignedStudents = $state<{ id: string; name: string; sats_number: string }[]>([]);

	let loading = $state(true);
	let selectedMentorId = $state('');
	let selectedStudentIds = $state<Set<string>>(new Set());
	let statusMsg = $state('');
	let errorMsg = $state('');
	let studentSearch = $state('');
	let assignmentSearch = $state('');
	let assignLoading = $state(false);

	const selectedMentor = $derived(stats?.mentors.find(m => m.id === selectedMentorId));

	const filteredUnassigned = $derived(
		unassignedStudents.filter(s =>
			!studentSearch.trim() ||
			s.name.toLowerCase().includes(studentSearch.toLowerCase()) ||
			s.sats_number.toLowerCase().includes(studentSearch.toLowerCase())
		)
	);

	const mentorAssignments = $derived(
		assignments.filter(a =>
			(!selectedMentorId || a.mentor_id === selectedMentorId) &&
			(!assignmentSearch.trim() ||
				a.student_name.toLowerCase().includes(assignmentSearch.toLowerCase()) ||
				a.sats_number.toLowerCase().includes(assignmentSearch.toLowerCase()) ||
				a.mentor_name.toLowerCase().includes(assignmentSearch.toLowerCase()))
		)
	);

	const totalAssigned = $derived(assignments.length);

	async function load() {
		if (!selectedYear) return;
		loading = true;
		errorMsg = '';
		const params = new URLSearchParams({ academic_year_id: selectedYear });
		if (selectedClass) params.set('class_id', selectedClass);
		const [aRes, sRes, uRes] = await Promise.all([
			api<typeof assignments>('GET', `/mentors/assignments?${params}`),
			api<typeof stats>('GET', `/mentors/stats?${params}`),
			api<typeof unassignedStudents>('GET', `/students?${params}&unassigned=true`),
		]);
		if (aRes.error) errorMsg = aRes.error.message;
		if (aRes.data) assignments = aRes.data;
		if (sRes.data) stats = sRes.data;
		if (uRes.data) unassignedStudents = uRes.data;
		loading = false;
	}

	onMount(async () => {
		const [cr, yr] = await Promise.all([
			api<Class[]>('GET', '/classes'),
			api<AcademicYear[]>('GET', '/academic-years'),
		]);
		if (cr.data) classes = cr.data;
		if (yr.data) {
			years = yr.data;
			const cur = yr.data.find(y => y.is_current);
			if (cur) { selectedYear = cur.id; load(); }
		}
	});

	async function assign() {
		if (!selectedMentorId || selectedStudentIds.size === 0) return;
		assignLoading = true;
		statusMsg = '';
		errorMsg = '';
		try {
			for (const sid of selectedStudentIds) {
				const res = await api('POST', '/mentors/assignments', { mentor_id: selectedMentorId, student_id: sid, academic_year_id: selectedYear });
				if (res.error) throw new Error(res.error.message);
			}
			statusMsg = `Assigned ${selectedStudentIds.size} student(s) to ${selectedMentor?.name ?? ''}`;
			selectedStudentIds = new Set();
			await load();
		} catch (e) {
			errorMsg = e instanceof Error ? e.message : 'Failed to assign students';
		} finally {
			assignLoading = false;
			setTimeout(() => (statusMsg = ''), 4000);
		}
	}

	async function unassign(id: string) {
		statusMsg = '';
		errorMsg = '';
		const res = await api('DELETE', `/mentors/assignments/${id}`);
		if (res.error) {
			errorMsg = res.error.message;
		} else {
			await load();
		}
	}

	function toggleStudent(id: string) {
		const next = new Set(selectedStudentIds);
		if (next.has(id)) next.delete(id); else next.add(id);
		selectedStudentIds = next;
	}

	function selectMentor(id: string) {
		selectedMentorId = id;
		selectedStudentIds = new Set();
	}

	$effect(() => { if (selectedYear) load(); });
</script>

<svelte:head><title>Mentor Assignments ? Pragati</title></svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<div class="flex items-center gap-3">
		<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-teal-500 to-emerald-700 flex items-center justify-center">
			<Users size={20} class="text-white" />
		</div>
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Mentor Assignments</h1>
			<p class="text-sm text-slate-500">Manage which students each mentor is responsible for</p>
		</div>
	</div>

	{#if statusMsg}
		<div class="text-sm px-4 py-2.5 rounded-lg bg-emerald-50 text-emerald-700 border border-emerald-200">{statusMsg}</div>
	{/if}
	{#if errorMsg}
		<div class="flex items-center gap-2 text-sm px-4 py-2.5 rounded-lg bg-red-50 text-red-700 border border-red-200"><AlertCircle size={16} /> {errorMsg}</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
		<div class="flex flex-wrap gap-3">
			<div class="w-44">
				<Select bind:value={selectedYear} options={[{ id: '', name: 'Select year' }, ...years.map(y => ({ id: y.id, name: y.name }))]} />
			</div>
			<div class="w-44">
				<Select bind:value={selectedClass} options={[{ id: '', name: 'All Classes' }, ...classes.map(c => ({ id: c.id, name: c.name }))]} />
			</div>
		</div>
	</div>

	{#if loading}
		<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Loading...</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
			<div class="bg-white rounded-xl border border-slate-200 p-4 flex items-center gap-3">
				<div class="w-10 h-10 rounded-lg bg-teal-50 flex items-center justify-center"><Users size={18} class="text-teal-600" /></div>
				<div><div class="text-2xl font-bold text-slate-900">{stats?.mentors.length ?? 0}</div><div class="text-xs text-slate-500">Mentors</div></div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-4 flex items-center gap-3">
				<div class="w-10 h-10 rounded-lg bg-blue-50 flex items-center justify-center"><GraduationCap size={18} class="text-blue-600" /></div>
				<div><div class="text-2xl font-bold text-slate-900">{totalAssigned}</div><div class="text-xs text-slate-500">Students assigned</div></div>
			</div>
			<div class="bg-white rounded-xl border border-slate-200 p-4 flex items-center gap-3">
				<div class="w-10 h-10 rounded-lg bg-amber-50 flex items-center justify-center"><UserRound size={18} class="text-amber-600" /></div>
				<div><div class="text-2xl font-bold text-amber-600">{stats?.unassigned ?? 0}</div><div class="text-xs text-slate-500">Unassigned</div></div>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Mentor Stats -->
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<h2 class="text-base font-semibold text-slate-900 mb-3">Mentors</h2>
				{#if stats}
					<div class="space-y-2 max-h-[480px] overflow-y-auto pr-1">
						{#each stats.mentors as m}
							<button onclick={() => selectMentor(m.id)}
								class="w-full flex items-center justify-between gap-2 p-2.5 rounded-lg border transition-colors {selectedMentorId === m.id ? 'border-primary-300 bg-primary-50 ring-1 ring-primary-200' : 'border-slate-100 hover:bg-slate-50'}">
								<div class="flex items-center gap-2 min-w-0">
									<div class="w-7 h-7 rounded-full bg-slate-100 flex items-center justify-center shrink-0">
										<UserRound size={14} class="text-slate-500" />
									</div>
									<span class="text-sm font-medium text-slate-700 truncate">{m.name}</span>
								</div>
								<span class="text-xs px-2 py-0.5 rounded-full shrink-0 {m.student_count >= 25 ? 'bg-red-100 text-red-600' : m.student_count >= 20 ? 'bg-amber-100 text-amber-600' : 'bg-emerald-100 text-emerald-600'}">
									{m.student_count}/25
								</span>
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Unassigned Students -->
			<div class="lg:col-span-2 bg-white rounded-xl border border-slate-200 p-5">
				<div class="flex flex-wrap items-center justify-between gap-3 mb-3">
					<div>
						<h2 class="text-base font-semibold text-slate-900">Unassigned Students</h2>
						<p class="text-xs text-slate-500">
							{#if selectedMentor}
								Assigning to <span class="font-medium text-teal-600">{selectedMentor.name}</span>
							{:else}
								Select a mentor to begin
							{/if}
						</p>
					</div>
					{#if selectedMentorId && selectedStudentIds.size > 0}
						<Button icon={UserPlus} onclick={assign} loading={assignLoading}>Assign {selectedStudentIds.size}</Button>
					{/if}
				</div>

				{#if !selectedMentorId}
					<div class="p-10 text-center text-sm text-slate-400">
						<UserRound size={32} class="mx-auto mb-2 text-slate-300" />
						Select a mentor on the left to assign students
					</div>
				{:else}
					<div class="relative mb-3">
						<Search size={15} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							bind:value={studentSearch}
							placeholder="Search by name or SATS..."
							class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
						/>
					</div>
					{#if unassignedStudents.length === 0}
						<div class="p-8 text-center text-sm text-slate-400">All students assigned</div>
					{:else if filteredUnassigned.length === 0}
						<div class="p-8 text-center text-sm text-slate-400">No students match your search</div>
					{:else}
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-2 max-h-72 overflow-y-auto pr-1">
							{#each filteredUnassigned as s}
								<button onclick={() => toggleStudent(s.id)}
									class="flex items-center gap-3 p-2.5 rounded-lg border text-left transition-colors {selectedStudentIds.has(s.id) ? 'border-primary-300 bg-primary-50' : 'border-slate-100 hover:bg-slate-50'}">
									<div class="w-5 h-5 rounded border shrink-0 {selectedStudentIds.has(s.id) ? 'bg-primary-600 border-primary-600' : 'border-slate-300'} flex items-center justify-center">
										{#if selectedStudentIds.has(s.id)}<span class="text-white text-[11px] leading-none">✓</span>{/if}
									</div>
									<div class="min-w-0">
										<div class="text-sm font-medium text-slate-700 truncate">{s.name}</div>
										<div class="text-xs text-slate-400">SATS: {s.sats_number}</div>
									</div>
								</button>
							{/each}
						</div>
					{/if}
				{/if}
			</div>
		</div>

		<!-- Current Assignments -->
		<div class="bg-white rounded-xl border border-slate-200 p-5">
			<div class="flex flex-wrap items-center justify-between gap-3 mb-3">
				<div>
					<h2 class="text-base font-semibold text-slate-900">
						{#if selectedMentor}{selectedMentor.name}'s Students{:else}All Assignments{/if}
					</h2>
					<p class="text-xs text-slate-500">{mentorAssignments.length} student(s)</p>
				</div>
				<div class="flex items-center gap-2">
					<div class="relative">
						<Search size={15} class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							bind:value={assignmentSearch}
							placeholder="Search students..."
							class="w-56 pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
						/>
					</div>
					{#if selectedMentorId}
						<Button variant="ghost" size="sm" onclick={() => { selectedMentorId = ''; assignmentSearch = ''; }}>Clear mentor</Button>
					{/if}
				</div>
			</div>
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200">
							<th class="px-4 py-2 text-left font-semibold text-slate-600">Student</th>
							<th class="px-4 py-2 text-left font-semibold text-slate-600">SATS</th>
							<th class="px-4 py-2 text-left font-semibold text-slate-600">Class</th>
							{#if !selectedMentorId}<th class="px-4 py-2 text-left font-semibold text-slate-600">Mentor</th>{/if}
							<th class="px-4 py-2 w-10"></th>
						</tr>
					</thead>
					<tbody>
						{#each mentorAssignments as a}
							<tr class="border-b border-slate-100 hover:bg-slate-50/50">
								<td class="px-4 py-2 text-slate-700">{a.student_name}</td>
								<td class="px-4 py-2 text-slate-500">{a.sats_number}</td>
								<td class="px-4 py-2 text-slate-500">{a.class_name}</td>
								{#if !selectedMentorId}<td class="px-4 py-2 text-slate-500">{a.mentor_name}</td>{/if}
								<td class="px-4 py-2">
									<button onclick={() => unassign(a.id)} title="Unassign" class="p-1 text-slate-400 hover:text-red-500"><Trash2 size={14} /></button>
								</td>
							</tr>
						{:else}
							<tr><td colspan="6" class="px-4 py-8 text-center text-slate-400">No assignments match</td></tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>