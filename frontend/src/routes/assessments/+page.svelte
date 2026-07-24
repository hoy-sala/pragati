<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { AssessmentCategory } from '$lib/types';
	import { onMount } from 'svelte';

	let categories = $state<AssessmentCategory[]>([]);
	let loading = $state(true);

	onMount(async () => {
		const res = await api<AssessmentCategory[]>('GET', '/assessment-categories');
		if (res.data) categories = res.data;
		loading = false;
	});
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Assessments</h1>
			<p class="text-sm text-slate-500 mt-1">Manage assessment categories and create assessments</p>
		</div>
		<a href="/assessments/create" class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 transition-colors">
			New Assessment
		</a>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
		<div class="bg-white rounded-xl border border-slate-200 p-6">
			<h2 class="text-lg font-semibold text-slate-900 mb-4">Categories</h2>
			{#if loading}
				<div class="text-sm text-slate-400">Loading...</div>
			{:else if categories.length === 0}
				<div class="text-sm text-slate-400">No categories configured</div>
			{:else}
				<div class="space-y-2">
					{#each categories as cat (cat.id)}
						<div class="flex items-center justify-between py-2 border-b border-slate-100 last:border-0">
							<div>
								<span class="text-sm font-medium text-slate-900">{cat.name}</span>
								<span class="text-xs text-slate-400 ml-2">{cat.code}</span>
							</div>
							<div class="text-xs text-slate-500">
								Weight: {cat.weightage}% | Order: {cat.sort_order}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-6">
			<h2 class="text-lg font-semibold text-slate-900 mb-4">Quick Actions</h2>
			<div class="space-y-3">
				<a href="/marks" class="block px-4 py-3 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors">
					<div class="text-sm font-medium text-slate-900">Enter Marks</div>
					<div class="text-xs text-slate-500 mt-0.5">Select assessment and enter marks using spreadsheet</div>
				</a>
				<a href="/assessments/create" class="block px-4 py-3 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors">
					<div class="text-sm font-medium text-slate-900">Create Assessment</div>
					<div class="text-xs text-slate-500 mt-0.5">Set up a new test, exam, or quiz</div>
				</a>
			</div>
		</div>
	</div>
</div>
