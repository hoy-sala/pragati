<script lang="ts">
	import { ChevronDown, X } from 'lucide-svelte';

	let {
		value = $bindable(''),
		options,
		labelKey = 'name',
		valueKey = 'id',
		placeholder = 'Select...',
		label = '',
		icon,
		clearable = false,
		disabled = false,
		size = 'md',
		onselect,
		class: className = '',
	}: {
		value?: string;
		options: any[];
		labelKey?: string;
		valueKey?: string;
		placeholder?: string;
		label?: string;
		icon?: import('svelte').ComponentType;
		clearable?: boolean;
		disabled?: boolean;
		size?: 'sm' | 'md';
		onselect?: (val: string) => void;
		class?: string;
	} = $props();

	let open = $state(false);
	let selectedLabel = $derived(
		options.find((o) => o[valueKey] === value)?.[labelKey] ?? ''
	);
	let listId = $state('select-' + Math.random().toString(36).slice(2, 8));

	function toggle() {
		if (!disabled) open = !open;
	}

	function select(val: string) {
		value = val;
		open = false;
		onselect?.(val);
	}

	function clear(e: MouseEvent) {
		e.stopPropagation();
		value = '';
		open = false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') open = false;
		if (e.key === 'Enter' || e.key === ' ') toggle();
	}

	function handleBlur() {
		setTimeout(() => { open = false; }, 150);
	}
</script>

<div class="relative {className}" role="combobox" aria-expanded={open} aria-haspopup="listbox" tabindex="-1" onkeydown={handleKeydown}>
	{#if label}
		<label for={listId} class="block text-xs font-medium text-slate-500 mb-1.5">{label}</label>
	{/if}
	<button
		type="button"
		id={listId}
		onclick={toggle}
		onblur={handleBlur}
		disabled={disabled}
		aria-haspopup="listbox"
		aria-expanded={open}
		class="w-full flex items-center gap-2 rounded-lg border transition-all
			{size === 'sm' ? 'px-2 py-1 text-xs' : 'px-3 py-2 text-sm'}
			{open
				? 'border-primary-400 ring-2 ring-primary-100'
				: 'border-slate-300 hover:border-slate-400'}
			{disabled ? 'bg-slate-50 text-slate-400 cursor-not-allowed' : 'bg-white text-slate-900 cursor-pointer'}
			focus:outline-none"
	>
		{#if icon}
			<svelte:component this={icon} size={16} class="text-slate-400 shrink-0" />
		{/if}
		<span class="flex-1 text-left truncate {selectedLabel ? '' : 'text-slate-400'}">
			{selectedLabel || placeholder}
		</span>
		{#if clearable && value}
			<span role="button" tabindex="-1" onclick={clear} class="p-0.5 rounded hover:bg-slate-200 text-slate-400 hover:text-slate-600 cursor-pointer">
				<X size={14} />
			</span>
		{/if}
		<ChevronDown
			size={16}
			class="text-slate-400 transition-transform duration-200 {open ? 'rotate-180' : ''}"
		/>
	</button>

	{#if open}
		<ul
			role="listbox"
			tabindex="-1"
			class="absolute z-50 mt-1 w-full bg-white border border-slate-200 rounded-lg shadow-lg max-h-60 overflow-y-auto py-1"
		>
			{#if !options.length}
				<li class="px-3 py-2 text-sm text-slate-400 text-center">No options</li>
			{:else}
				{#each options as opt (opt[valueKey])}
					{@const isSelected = opt[valueKey] === value}
					<li
						role="option"
						aria-selected={isSelected}
						onmousedown={() => select(opt[valueKey])}
						class="cursor-pointer transition-colors flex items-center gap-2
							{size === 'sm' ? 'px-2 py-1 text-xs' : 'px-3 py-2 text-sm'}
							{isSelected
								? 'bg-primary-50 text-primary-700 font-medium'
								: 'text-slate-700 hover:bg-slate-100'}"
					>
						{opt[labelKey]}
					</li>
				{/each}
			{/if}
		</ul>
	{/if}
</div>
