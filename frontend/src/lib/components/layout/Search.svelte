<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import { goto } from '$app/navigation';
	import { Search, User, X } from 'lucide-svelte';
	import { onMount } from 'svelte';

	let query = $state('');
	let results = $state<{ id: string; name: string; sats_number: string; class_name: string }[]>([]);
	let open = $state(false);
	let loading = $state(false);
	let el = $state<HTMLDivElement>();
	let searchSeq = 0;

	async function search() {
		if (query.trim().length < 2) { results = []; open = false; return; }
		loading = true;
		const seq = ++searchSeq;
		const res = await api<{ id: string; name: string; sats_number: string; class_name: string }[]>('GET', `/students/search?q=${encodeURIComponent(query)}`);
		if (seq === searchSeq && res.data) {
			results = res.data;
			open = results.length > 0;
		}
		loading = false;
	}

	$effect(() => {
		if (query.trim().length >= 2) {
			const t = setTimeout(search, 300);
			return () => clearTimeout(t);
		} else {
			results = [];
			open = false;
		}
	});

	function select(id: string) {
		open = false;
		query = '';
		goto(`/students/${id}`);
	}

	function close() { open = false; }

	onMount(() => {
		function onClick(e: MouseEvent) {
			if (el && !el.contains(e.target as Node)) close();
		}
		document.addEventListener('click', onClick);
		return () => document.removeEventListener('click', onClick);
	});
</script>

<div class="relative px-3 py-2" bind:this={el}>
	<div class="flex items-center gap-2 px-3 py-2 rounded-lg bg-slate-100 border border-transparent focus-within:border-primary-300 focus-within:bg-white transition-colors">
		<Search size={14} class="text-slate-400 shrink-0" />
		<input bind:value={query} placeholder="Search students..." class="flex-1 bg-transparent text-sm text-slate-700 placeholder:text-slate-400 outline-none" />
		{#if query}
			<button onclick={() => { query = ''; results = []; open = false; }} class="text-slate-400 hover:text-slate-600"><X size={14} /></button>
		{/if}
	</div>

	{#if open && results.length > 0}
		<div class="absolute left-3 right-3 top-full mt-1 bg-white border border-slate-200 rounded-lg shadow-lg z-50 max-h-72 overflow-y-auto">
			{#each results as s}
				<button onclick={() => select(s.id)} class="w-full flex items-center gap-3 px-3 py-2.5 text-left hover:bg-slate-50 transition-colors">
					<div class="w-8 h-8 rounded-full bg-primary-50 flex items-center justify-center shrink-0">
						<User size={14} class="text-primary-600" />
					</div>
					<div class="min-w-0 flex-1">
						<div class="text-sm font-medium text-slate-800 truncate">{s.name}</div>
						<div class="text-xs text-slate-400">SATS: {s.sats_number} · {s.class_name}</div>
					</div>
				</button>
			{/each}
		</div>
	{:else if open && loading}
		<div class="absolute left-3 right-3 top-full mt-1 bg-white border border-slate-200 rounded-lg shadow-lg z-50 p-3 text-sm text-slate-400">Searching...</div>
	{/if}
</div>
