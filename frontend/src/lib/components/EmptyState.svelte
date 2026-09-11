<script lang="ts">
	import type { Snippet, ComponentType } from 'svelte';

	let {
		icon,
		title,
		hint = '',
		action
	}: {
		icon?: ComponentType;
		title: string;
		hint?: string;
		action?: Snippet;
	} = $props();
</script>

<div class="p-12 text-center">
	{#if icon}
		{@const Icon = icon}
		<Icon size={36} class="mx-auto text-slate-300 mb-3" />
	{/if}
	<p class="text-slate-500 text-sm">{title}</p>
	{#if hint}
		<p class="text-xs text-slate-400 mt-1">{hint}</p>
	{/if}
	{#if action}
		<div class="mt-3">
			{@render action()}
		</div>
	{/if}
</div>
