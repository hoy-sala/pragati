<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';
	import { Plus, Heart, Brain, ShieldAlert, BookOpen, AlertTriangle } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import type { AcademicYear } from '$lib/types';

	let years = $state<AcademicYear[]>([]);
	let selectedYear = $state('');
	let logs = $state<{ id: string; student_id: string; student_name: string; log_date: string; category: string; severity: string; description: string; action_taken: string; parent_informed: boolean; reviewed_by_principal: boolean }[]>([]);
	let loading = $state(true);
	let showForm = $state(false);
	let students = $state<{ id: string; name: string }[]>([]);
	let form = $state({ student_id: '', category: 'health', severity: 'low', description: '', action_taken: '', parent_informed: false });

	const categories = [{ id: 'health', name: 'Health', icon: Heart }, { id: 'behavior', name: 'Behavior', icon: Brain }, { id: 'grievance', name: 'Grievance', icon: ShieldAlert }, { id: 'academic', name: 'Academic', icon: BookOpen }];
	const severities = [{ id: 'low', name: 'Low', cls: 'bg-emerald-100 text-emerald-700' }, { id: 'medium', name: 'Medium', cls: 'bg-blue-100 text-blue-700' }, { id: 'high', name: 'High', cls: 'bg-amber-100 text-amber-700' }, { id: 'urgent', name: 'Urgent', cls: 'bg-red-100 text-red-700' }];

	onMount(async () => {
		const [yr, roster] = await Promise.all([
			api<AcademicYear[]>('GET', '/academic-years'),
			api<{ id: string; name: string }[]>('GET', '/mentors/roster?academic_year_id=current'),
		]);
		if (yr.data) { years = yr.data; const cur = yr.data.find(y => y.is_current); if (cur) selectedYear = cur.id; }
		if (roster.data) students = roster.data;
	});

	async function loadLogs() {
		loading = true;
		const res = await api<typeof logs>('GET', `/mentors/logs`);
		if (res.data) logs = res.data;
		loading = false;
	}

	$effect(() => { if (selectedYear) loadLogs(); });

	async function submitLog() {
		if (!form.student_id || !form.description) return;
		await api('POST', '/mentors/logs', form);
		showForm = false;
		form = { student_id: '', category: 'health', severity: 'low', description: '', action_taken: '', parent_informed: false };
		loadLogs();
	}

	function severityCls(sev: string) {
		return severities.find(s => s.id === sev)?.cls || 'bg-slate-100 text-slate-600';
	}
	function catName(cat: string) {
		return categories.find(c => c.id === cat)?.name || cat;
	}
</script>

<svelte:head><title>Mentor Logs - Pragati</title></svelte:head>

<div class="max-w-5xl mx-auto space-y-6">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-3">
			<div class="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500 to-pink-700 flex items-center justify-center"><ShieldAlert size={20} class="text-white" /></div>
			<div><h1 class="text-2xl font-bold text-slate-900">Mentor Logs</h1><p class="text-sm text-slate-500">Health, behavior, grievance tracking</p></div>
		</div>
		<Button icon={Plus} onclick={() => showForm = !showForm}>New Log</Button>
	</div>

	{#if showForm}
		<div class="bg-white rounded-xl border border-slate-200 p-5 space-y-4 no-print">
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
				<div>
					<label for="log-student" class="block text-xs font-medium text-slate-600 mb-1">Student *</label>
					<select id="log-student" bind:value={form.student_id} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
						<option value="">Select student</option>
						{#each students as s}<option value={s.id}>{s.name}</option>{/each}
					</select>
				</div>
				<div>
					<label for="log-cat" class="block text-xs font-medium text-slate-600 mb-1">Category</label>
					<select id="log-cat" bind:value={form.category} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
						{#each categories as c}<option value={c.id}>{c.name}</option>{/each}
					</select>
				</div>
				<div>
					<label for="log-sev" class="block text-xs font-medium text-slate-600 mb-1">Severity</label>
					<select id="log-sev" bind:value={form.severity} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
						{#each severities as s}<option value={s.id}>{s.name}</option>{/each}
					</select>
				</div>
				<div class="flex items-end pb-1"><label class="flex items-center gap-2 text-sm"><input type="checkbox" bind:checked={form.parent_informed} class="rounded" /> Parent informed</label></div>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Description *</label>
				<textarea bind:value={form.description} rows="3" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" placeholder="Describe the issue..."></textarea>
			</div>
			<div>
				<label class="block text-xs font-medium text-slate-600 mb-1">Action Taken</label>
				<textarea bind:value={form.action_taken} rows="2" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" placeholder="Steps taken..."></textarea>
			</div>
			<Button onclick={submitLog}>Submit Log</Button>
		</div>
	{/if}

	{#if loading}<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Loading...</div>
	{:else}
		<div class="space-y-3">
			{#each logs as l}
				<div class="bg-white rounded-xl border border-slate-200 p-4">
					<div class="flex items-start justify-between mb-2">
						<div class="flex items-center gap-2">
							<span class="text-sm font-medium text-slate-800">{l.student_name}</span>
							<span class="text-xs px-2 py-0.5 rounded {severityCls(l.severity)}">{l.severity}</span>
							<span class="text-xs text-slate-400">{catName(l.category)}</span>
						</div>
						<span class="text-xs text-slate-400">{l.log_date}</span>
					</div>
					<p class="text-sm text-slate-600 mb-2">{l.description}</p>
					{#if l.action_taken}<p class="text-xs text-slate-500"><strong>Action:</strong> {l.action_taken}</p>{/if}
					<div class="flex items-center gap-3 mt-2">
						{#if l.parent_informed}<span class="text-xs text-emerald-600">Parent informed</span>{/if}
						{#if l.reviewed_by_principal}<span class="text-xs text-blue-600">Reviewed by Principal</span>{/if}
						{#if l.severity === 'urgent' || l.severity === 'high'}<AlertTriangle size={14} class="text-amber-500" />{/if}
					</div>
				</div>
			{:else}
				<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">No logs yet</div>
			{/each}
		</div>
	{/if}
</div>
