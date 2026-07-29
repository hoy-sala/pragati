<script lang="ts">
	import { getAuthState } from '$lib/stores/auth.svelte';
	import { api } from '$lib/api/client.svelte';
	import { onMount } from 'svelte';

	const auth = getAuthState();

	let stats = $state<{
		total_students: number;
		total_teachers: number;
		total_classes: number;
		total_assessments: number;
		students_by_class: { class: string; count: number }[];
	} | null>(null);
	let loading = $state(true);

	onMount(async () => {
		const res = await api<typeof stats>('GET', '/dashboard/stats');
		if (res.data) stats = res.data;
		loading = false;
	});
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Dashboard</h1>
			<p class="text-sm text-slate-500 mt-1">Welcome back, {auth.currentUser?.name}</p>
		</div>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">{loading ? '--' : stats?.total_students ?? 0}</div>
			<div class="text-sm text-slate-500 mt-1">Students</div>
		</div>
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">{loading ? '--' : stats?.total_teachers ?? 0}</div>
			<div class="text-sm text-slate-500 mt-1">Teachers</div>
		</div>
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">{loading ? '--' : stats?.total_classes ?? 0}</div>
			<div class="text-sm text-slate-500 mt-1">Classes</div>
		</div>
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">{loading ? '--' : stats?.total_assessments ?? 0}</div>
			<div class="text-sm text-slate-500 mt-1">Assessments</div>
		</div>
	</div>

	{#if stats?.students_by_class && stats.students_by_class.length > 0}
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<h2 class="text-lg font-semibold text-slate-900 mb-3">Students by Class</h2>
			<div class="space-y-2">
				{#each stats.students_by_class as item}
					<div class="flex items-center gap-3">
						<span class="text-sm font-medium text-slate-700 w-20">{item.class}</span>
						<div class="flex-1 h-5 bg-slate-100 rounded-full overflow-hidden">
							<div class="h-full bg-primary-500 rounded-full transition-all" style="width: {Math.max(4, (item.count / Math.max(...stats.students_by_class.map(s => s.count))) * 100)}%"></div>
						</div>
						<span class="text-sm text-slate-600 w-8 text-right">{item.count}</span>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>
