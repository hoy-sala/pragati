<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		open = $bindable(false),
		title = '',
		maxWidth = 'max-w-md',
		children,
		footer
	}: {
		open?: boolean;
		title?: string;
		maxWidth?: string;
		children: Snippet;
		footer?: Snippet;
	} = $props();

	function close() {
		open = false;
	}

	function onkeydown(e: KeyboardEvent) {
		if (open && e.key === 'Escape') close();
	}
</script>

<svelte:window onkeydown={onkeydown} />

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<button class="absolute inset-0 bg-black/40 cursor-default" onclick={close} aria-label="Close dialog" tabindex="-1"></button>
		<div role="dialog" aria-modal="true" aria-label={title} class="relative bg-white rounded-xl shadow-xl w-full {maxWidth} p-6">
			{#if title}
				<h3 class="text-base font-semibold text-slate-900 mb-4">{title}</h3>
			{/if}
			{@render children()}
			{#if footer}
				<div class="flex justify-end gap-2 mt-6">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}
