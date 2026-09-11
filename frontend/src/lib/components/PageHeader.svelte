<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		title,
		subtitle = '',
		icon,
		actions
	}: {
		title: string;
		subtitle?: string;
		icon?: Snippet;
		actions?: Snippet;
	} = $props();
</script>

<div class="flex items-center justify-between gap-3">
	<div class="flex items-center gap-3 min-w-0">
		{#if icon}
			{@render icon()}
		{/if}
		<div class="min-w-0">
			<h1 class="text-2xl font-bold text-slate-900">{title}</h1>
			{#if subtitle}
				<p class="text-sm text-slate-500 mt-1">{subtitle}</p>
			{/if}
		</div>
	</div>
	{#if actions}
		<div class="flex items-center gap-2 shrink-0">
			{@render actions()}
		</div>
	{/if}
</div>
