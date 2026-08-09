<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';
	import { Users, UserPlus, Trash2, AlertCircle } from 'lucide-svelte';
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

	async function load() {
		if (!selectedYear) return;
		loading = true;
		const params = new URLSearchParams({ academic_year_id: selectedYear });
		if (selectedClass) params.set('class_id', selectedClass);
		const [aRes, sRes, uRes] = await Promise.all([
			api<typeof assignments>('GET', `/mentors/assignments?${params}`),
			api<typeof stats>('GET', `/mentors/stats?${params}`),
			api<typeof unassignedStudents>('GET', `/students?${params}&unassigned=true`),
		]);
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
		statusMsg = '';
		for (const sid of selectedStudentIds) {
			await api('POST', '/mentors/assignments', { mentor_id: selectedMentorId, student_id: sid, academic_year_id: selectedYear });
		}
		statusMsg = `Assigned ${selectedStudentIds.size} student(s)`;
		selectedStudentIds = new Set();
		load();
	}

	async function unassign(id: string) {
		await api('DELETE', `/mentors/assignments/${id}`);
		load();
	}

	function toggleStudent(id: string) {
		const next = new Set(selectedStudentIds);
		if (next.has(id)) next.delete(id); else next.add(id);
		selectedStudentIds = next;
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
			<p class="text-sm text-slate-500">Assign students to mentors (max 25 per mentor)</p>
		</div>
	</div>

	{#if statusMsg}<div class="text-sm px-4 py-2 rounded-lg bg-emerald-50 text-emerald-700">{statusMsg}</div>{/if}

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
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
			<!-- Mentor Stats -->
			<div class="bg-white rounded-xl border border-slate-200 p-5">
				<h2 class="text-base font-semibold text-slate-900 mb-3">Mentors</h2>
				{#if stats}
					<div class="mb-3 p-3 rounded-lg bg-amber-50 border border-amber-200">
						<span class="text-sm font-medium text-amber-700">{stats.unassigned} unassigned</span>
					</div>
					<div class="space-y-2">
						{#each stats.mentors as m}
							<button onclick={() => selectedMentorId = m.id}
								class="w-full flex items-center justify-between p-2.5 rounded-lg border transition-colors {selectedMentorId === m.id ? 'border-primary-300 bg-primary-50' : 'border-slate-100 hover:bg-slate-50'}">
								<span class="text-sm font-medium text-slate-700">{m.name}</span>
								<span class="text-xs px-2 py-0.5 rounded-full {m.student_count >= 25 ? 'bg-red-100 text-red-600' : m.student_count >= 20 ? 'bg-amber-100 text-amber-600' : 'bg-emerald-100 text-emerald-600'}">
									{m.student_count}/25
								</span>
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Unassigned Students -->
			<div class="lg:col-span-2 bg-white rounded-xl border border-slate-200 p-5">
				<div class="flex items-center justify-between mb-3">
					<h2 class="text-base font-semibold text-slate-900">Students</h2>
					{#if selectedMentorId && selectedStudentIds.size > 0}
						<Button icon={UserPlus} onclick={assign}>Assign {selectedStudentIds.size}</Button>
					{/if}
				</div>
				{#if !selectedMentorId}
					<div class="p-8 text-center text-sm text-slate-400">Select a mentor to assign students</div>
				{:else if unassignedStudents.length === 0}
					<div class="p-8 text-center text-sm text-slate-400">All students assigned</div>
				{:else}
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-2 max-h-96 overflow-y-auto">
						{#each unassignedStudents as s}
							<button onclick={() => toggleStudent(s.id)}
								class="flex items-center gap-3 p-2.5 rounded-lg border text-left transition-colors {selectedStudentIds.has(s.id) ? 'border-primary-300 bg-primary-50' : 'border-slate-100 hover:bg-slate-50'}">
								<div class="w-5 h-5 rounded border {selectedStudentIds.has(s.id) ? 'bg-primary-600 border-primary-600' : 'border-slate-300'} flex items-center justify-center">
									{#if selectedStudentIds.has(s.id)}<span class="text-white text-xs">?</span>{/if}
								</div>
								<div>
									<div class="text-sm font-medium text-slate-700">{s.name}</div>
									<div class="text-xs text-slate-400">SATS: {s.sats_number}</div>
								</div>
							</button>
						{/each}
					</div>
				{/if}
			</div>
		</div>

		<!-- Current Assignments -->
		<div class="bg-white rounded-xl border border-slate-200 p-5">
			<h2 class="text-base font-semibold text-slate-900 mb-3">Current Assignments</h2>
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200">
							<th class="px-4 py-2 text-left font-semibold text-slate-600">Mentor</th>
							<th class="px-4 py-2 text-left font-semibold text-slate-600">Student</th>
							<th class="px-4 py-2 text-left font-semibold text-slate-600">SATS</th>
							<th class="px-4 py-2 text-left font-semibold text-slate-600">Class</th>
							<th class="px-4 py-2 w-10"></th>
						</tr>
					</thead>
					<tbody>
						{#each assignments as a}
							<tr class="border-b border-slate-100 hover:bg-slate-50/50">
								<td class="px-4 py-2 text-slate-700">{a.mentor_name}</td>
								<td class="px-4 py-2 text-slate-700">{a.student_name}</td>
								<td class="px-4 py-2 text-slate-500">{a.sats_number}</td>
								<td class="px-4 py-2 text-slate-500">{a.class_name}</td>
								<td class="px-4 py-2">
									<button onclick={() => unassign(a.id)} class="p-1 text-slate-400 hover:text-red-500"><Trash2 size={14} /></button>
								</td>
							</tr>
						{:else}
							<tr><td colspan="5" class="px-4 py-8 text-center text-slate-400">No assignments yet</td></tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>
