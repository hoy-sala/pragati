<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Student, Pagination } from '$lib/types';
	import { onMount } from 'svelte';

	let students: Student[] = $state([]);
	let total = $state(0);
	let loading = $state(true);
	let classNameMap = $state<Record<string, string>>({});

	onMount(async () => {
		const [res, classRes] = await Promise.all([
			api<Student[]>('GET', '/students?limit=50'),
			api<{ id: string; name: string }[]>('GET', '/classes?limit=100')
		]);
		if (classRes.data) {
			for (const c of classRes.data) {
				classNameMap[c.id] = c.name;
			}
		}
		if (res.data) {
			students = res.data;
			total = (res.meta as Pagination)?.total ?? students.length;
		}
		loading = false;
	});
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Students</h1>
			<p class="text-sm text-slate-500 mt-1">{total} students enrolled</p>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
		<table class="w-full text-sm">
			<thead>
				<tr class="bg-slate-50 text-slate-600">
					<th class="text-left px-4 py-3 font-medium">SATS</th>
					<th class="text-left px-4 py-3 font-medium">Name</th>
					<th class="text-left px-4 py-3 font-medium">Class</th>
					<th class="text-left px-4 py-3 font-medium">Roll No</th>
				</tr>
			</thead>
			<tbody>
				{#if loading}
					<tr><td colspan="4" class="px-4 py-8 text-center text-slate-400">Loading...</td></tr>
				{:else if students.length === 0}
					<tr><td colspan="4" class="px-4 py-8 text-center text-slate-400">No students yet</td></tr>
				{:else}
					{#each students as s (s.id)}
						<tr class="border-t border-slate-100 hover:bg-slate-50">
							<td class="px-4 py-3 font-mono text-xs">{s.sats_number}</td>
							<td class="px-4 py-3">{s.first_name} {s.last_name}</td>
							<td class="px-4 py-3 text-slate-500">{classNameMap[s.class_id] || s.class_id}</td>
							<td class="px-4 py-3 text-slate-500">{s.roll_no}</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
</div>
