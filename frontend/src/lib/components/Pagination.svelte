<script lang="ts">
	let {
		page: currentPage = 1,
		total,
		pageSize = 20,
		onChange
	}: {
		page?: number;
		total: number;
		pageSize?: number;
		onChange: (page: number) => void;
	} = $props();

	let totalPages = $derived(Math.max(1, Math.ceil(total / pageSize)));
	let startItem = $derived((currentPage - 1) * pageSize + 1);
	let endItem = $derived(Math.min(currentPage * pageSize, total));

	function goto(p: number) {
		if (p < 1 || p > totalPages || p === currentPage) return;
		onChange(p);
	}

	function getPages(): (number | string)[] {
		const result: (number | string)[] = [];
		if (totalPages <= 7) {
			for (let i = 1; i <= totalPages; i++) result.push(i);
		} else {
			result.push(1);
			if (currentPage > 3) result.push('...');
			const start = Math.max(2, currentPage - 1);
			const end = Math.min(totalPages - 1, currentPage + 1);
			for (let i = start; i <= end; i++) result.push(i);
			if (currentPage < totalPages - 2) result.push('...');
			result.push(totalPages);
		}
		return result;
	}
</script>

{#if total > pageSize}
	<div class="flex items-center justify-between px-4 py-3 border-t border-slate-200">
		<span class="text-xs text-slate-500">
			{startItem}–{endItem} of {total}
		</span>
		<div class="flex items-center gap-1">
			<button onclick={() => goto(page - 1)} disabled={page === 1}
				class="px-2 py-1 text-xs rounded border border-slate-200 text-slate-500 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed">
				Prev
			</button>
			{#each getPages() as p}
				{#if p === '...'}
					<span class="px-1 text-xs text-slate-400">...</span>
				{:else}
					<button onclick={() => goto(p)}
						class="px-2.5 py-1 text-xs rounded {p === page ? 'bg-primary-600 text-white font-medium' : 'border border-slate-200 text-slate-600 hover:bg-slate-50'}">
						{p}
					</button>
				{/if}
			{/each}
			<button onclick={() => goto(page + 1)} disabled={page === totalPages}
				class="px-2 py-1 text-xs rounded border border-slate-200 text-slate-500 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed">
				Next
			</button>
		</div>
	</div>
{/if}
