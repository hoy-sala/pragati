<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { Phone, MapPin, AlertTriangle } from 'lucide-svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import { Users } from 'lucide-svelte';

	let { year }: { year: string } = $props();

	let roster = $state<{ id: string; sats_number: string; name: string; roll_no: number; gender: string; father_name: string; mother_name: string; parent_phone: string; address: string; class_name: string }[]>([]);
	let loading = $state(true);

	async function loadRoster() {
		if (!year) return;
		loading = true;
		const res = await api<typeof roster>('GET', `/mentors/roster?academic_year_id=${year}`);
		if (res.data) roster = res.data;
		loading = false;
	}

	$effect(() => { if (year) loadRoster(); });
</script>

<div class="space-y-6">
	{#if loading}<div class="bg-white rounded-xl border border-slate-200 p-12 text-center text-sm text-slate-400">Loading...</div>
	{:else if roster.length === 0}<div class="bg-white rounded-xl border border-slate-200"><EmptyState icon={Users} title="No students assigned to you" /></div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each roster as s}
				<div class="bg-white rounded-xl border border-slate-200 p-4">
					<div class="flex items-start justify-between mb-3">
						<div><h3 class="text-sm font-semibold text-slate-800">{s.name}</h3><p class="text-xs text-slate-400">SATS: {s.sats_number} . Roll {s.roll_no} . {s.class_name}</p></div>
						{#if !s.parent_phone}<AlertTriangle size={16} class="text-amber-500 shrink-0" aria-label="No phone number" />{/if}
					</div>
					<div class="space-y-2">
						{#if s.father_name}<div class="text-xs text-slate-600"><span class="text-slate-400">Father:</span> {s.father_name}</div>{/if}
						{#if s.mother_name}<div class="text-xs text-slate-600"><span class="text-slate-400">Mother:</span> {s.mother_name}</div>{/if}
						{#if s.parent_phone}<a href="tel:{s.parent_phone}" class="flex items-center gap-1.5 text-xs text-blue-600 hover:text-blue-700 font-medium"><Phone size={12} /> {s.parent_phone}</a>
						{:else}<div class="text-xs text-red-500 flex items-center gap-1"><Phone size={12} /> No phone</div>{/if}
						{#if s.address}<div class="flex items-start gap-1.5 text-xs text-slate-500 mt-1"><MapPin size={12} class="shrink-0 mt-0.5" /><span class="line-clamp-2">{s.address}</span></div>{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
