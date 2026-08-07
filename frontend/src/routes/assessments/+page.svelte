<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';
	import { ClipboardCheck, Plus, BookOpen, Users, Filter, Eye } from 'lucide-svelte';
	import Select from '$lib/components/Select.svelte';
	import Button from '$lib/components/Button.svelte';
	import type { Assessment, AssessmentCategory, Class, Subject } from '$lib/types';

	let assessments = $state<Assessment[]>([]);
	let categories = $state<AssessmentCategory[]>([]);
	let classes = $state<Class[]>([]);
	let subjects = $state<Subject[]>([]);

	let selectedClass = $state('');
	let selectedSubject = $state('');
	let selectedCategory = $state('');

	let loading = $state(true);
	let total = $state(0);
	let page = $state(0);
	const pageSize = 20;
	let publishing = $state<string | null>(null);
	let statusMsg = $state('');
	let statusType = $state<'info' | 'error' | 'success'>('info');

	let filteredSubjects = $derived(
		selectedClass ? subjects.filter(s => s.is_core || s.is_language) : subjects
	);

	async function load() {
		loading = true;
		const params = new URLSearchParams();
		if (selectedClass) params.set('class_id', selectedClass);
		if (selectedSubject) params.set('subject_id', selectedSubject);
		if (selectedCategory) params.set('category_id', selectedCategory);
		params.set('limit', String(pageSize));
		params.set('offset', String(page * pageSize));
		const res = await api<Assessment[]>('GET', `/assessments?${params}`);
		if (res.data) {
			assessments = res.data;
			total = (res.meta as unknown as { total: number })?.total ?? res.data.length;
		}
		loading = false;
	}

	onMount(async () => {
		const [cr, clR, subR] = await Promise.all([
			api<AssessmentCategory[]>('GET', '/assessment-categories'),
			api<Class[]>('GET', '/classes'),
			api<Subject[]>('GET', '/subjects'),
		]);
		if (cr.data) categories = cr.data;
		if (clR.data) classes = clR.data;
		if (subR.data) subjects = subR.data;
		await load();
	});

	$effect(() => {
		if (classes.length) load();
	});

	function onFilterChange() {
		page = 0;
		load();
	}

	function goto(p: number) {
		page = p;
		load();
	}

	async function publish(a: Assessment) {
		if (a.is_published) return;
		publishing = a.id;
		statusMsg = '';
		const res = await api<unknown>('POST', `/assessments/${a.id}/publish`);
		if (res.data) {
			statusMsg = `Published "${a.name}"`;
			statusType = 'success';
			await load();
		} else if (res.error) {
			statusMsg = res.error.message;
			statusType = 'error';
		}
		publishing = null;
	}

	function statusBadge(a: Assessment) {
		if (a.is_published) return { label: 'Published', cls: 'bg-emerald-50 text-emerald-700' };
		return { label: 'Draft', cls: 'bg-amber-50 text-amber-700' };
	}

	function marksProgress(a: Assessment): { pct: number; label: string; cls: string } {
		const entered = a.marks_count ?? 0;
		const total = a.student_count ?? 0;
		if (total === 0) return { pct: 0, label: '0/0', cls: 'bg-slate-200' };
		const pct = Math.round((entered / total) * 100);
		const cls = pct >= 80 ? 'bg-emerald-500' : pct >= 40 ? 'bg-amber-500' : 'bg-red-400';
		return { pct, label: `${entered}/${total}`, cls };
	}

	const totalPages = $derived(Math.max(1, Math.ceil(total / pageSize)));
</script>

<svelte:head>
	<title>Assessments — Pragati School Management</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Assessments</h1>
			<p class="text-sm text-slate-500 mt-1">View and manage all assessments across classes</p>
		</div>
		<a href="/assessments/create">
			<Button icon={Plus}>New Assessment</Button>
		</a>
	</div>

	{#if statusMsg}
		<div class="text-sm px-4 py-2 rounded-lg {statusType === 'success' ? 'bg-emerald-50 text-emerald-700' : 'bg-red-50 text-red-700'}">{statusMsg}</div>
	{/if}

	<div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
		<div class="flex items-center gap-2 mb-3">
			<Filter size={16} class="text-slate-400" />
			<span class="text-sm font-medium text-slate-600">Filters</span>
		</div>
		<div class="flex flex-wrap gap-3 items-end">
			<div class="w-44">
				<Select bind:value={selectedClass} options={[{ id: '', name: 'All Classes' }, ...classes.map(c => ({ id: c.id, name: c.name }))]} placeholder="All Classes" />
			</div>
			<div class="w-44">
				<Select bind:value={selectedSubject} options={[{ id: '', name: 'All Subjects' }, ...filteredSubjects.map(s => ({ id: s.id, name: s.name }))]} placeholder="All Subjects" />
			</div>
			<div class="w-44">
				<Select bind:value={selectedCategory} options={[{ id: '', name: 'All Categories' }, ...categories.map(c => ({ id: c.id, name: c.name }))]} placeholder="All Categories" />
			</div>
			<Button variant="secondary" onclick={onFilterChange}>Apply</Button>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		<div class="px-6 py-3 border-b border-slate-200 flex items-center justify-between">
			<span class="text-sm text-slate-500">{total} assessment{total !== 1 ? 's' : ''}</span>
			<div class="flex items-center gap-1">
				<Button variant="ghost" size="sm" disabled={page === 0} onclick={() => goto(page - 1)}>Prev</Button>
				<span class="text-xs text-slate-500 px-2">Page {page + 1}/{totalPages}</span>
				<Button variant="ghost" size="sm" disabled={page >= totalPages - 1} onclick={() => goto(page + 1)}>Next</Button>
			</div>
		</div>

		{#if loading}
			<div class="p-12 text-center text-sm text-slate-400">Loading...</div>
		{:else if assessments.length === 0}
			<div class="p-12 text-center">
				<BookOpen size={40} class="mx-auto text-slate-300 mb-3" />
				<p class="text-slate-500 text-sm">No assessments found</p>
				<a href="/assessments/create" class="text-sm text-primary-600 hover:text-primary-700 mt-2 inline-block">Create your first assessment</a>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200">
							<th class="px-4 py-2.5 text-left font-semibold text-slate-600">Assessment</th>
							<th class="px-4 py-2.5 text-left font-semibold text-slate-600">Class</th>
							<th class="px-4 py-2.5 text-left font-semibold text-slate-600">Subject</th>
							<th class="px-4 py-2.5 text-left font-semibold text-slate-600">Category</th>
							<th class="px-4 py-2.5 text-center font-semibold text-slate-600 w-20">Max</th>
							<th class="px-4 py-2.5 text-center font-semibold text-slate-600 w-28">Marks</th>
							<th class="px-4 py-2.5 text-center font-semibold text-slate-600 w-24">Status</th>
							<th class="px-4 py-2.5 text-right font-semibold text-slate-600 w-40">Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each assessments as a (a.id)}
							<tr class="border-b border-slate-100 hover:bg-slate-50/50">
								<td class="px-4 py-3">
									<div class="font-medium text-slate-800">{a.name || '—'}</div>
									{#if a.date}
										<div class="text-xs text-slate-400">{a.date}</div>
									{/if}
								</td>
								<td class="px-4 py-3 text-slate-600">{a.class_name || '—'}</td>
								<td class="px-4 py-3 text-slate-600">{a.subject_name || '—'}</td>
								<td class="px-4 py-3 text-slate-600">{a.category_name || '—'}</td>
								<td class="px-4 py-3 text-center text-slate-600">{a.max_marks}</td>
								<td class="px-4 py-3">
									{#if a.student_count !== undefined}
									{@const mp = marksProgress(a)}
									<div class="flex items-center gap-2">
										<div class="flex-1 h-1.5 bg-slate-200 rounded-full overflow-hidden">
											<div class="h-full {mp.cls}" style="width: {mp.pct}%"></div>
										</div>
										<span class="text-xs text-slate-500 w-10 text-right">{mp.label}</span>
									</div>
									{:else}
									<span class="text-xs text-slate-400">—</span>
									{/if}
								</td>
								<td class="px-4 py-3 text-center">
									{#if true}
									{@const sb = statusBadge(a)}
									<span class="text-xs font-medium px-2 py-0.5 rounded {sb.cls}">{sb.label}</span>
									{/if}
								</td>
								<td class="px-4 py-3">
									<div class="flex items-center justify-end gap-2">
										<a href="/marks?assessment_id={a.id}" title="Enter marks">
											<Button variant="ghost" size="sm" icon={ClipboardCheck}>Marks</Button>
										</a>
										{#if !a.is_published}
											<Button variant="ghost" size="sm" icon={Eye} disabled={publishing === a.id} onclick={() => publish(a)}>
												{publishing === a.id ? '...' : 'Publish'}
											</Button>
										{/if}
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
