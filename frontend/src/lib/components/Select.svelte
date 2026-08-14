<script lang="ts">
	import { ChevronDown, X, Search } from 'lucide-svelte';

	let {
		value = $bindable(''),
		options,
		labelKey = 'name',
		valueKey = 'id',
		placeholder = 'Select...',
		label = '',
		id,
		icon,
		clearable = false,
		disabled = false,
		size = 'md',
		searchable = false,
		onselect,
		class: className = '',
	}: {
		value?: string;
		options: any[];
		labelKey?: string;
		valueKey?: string;
		placeholder?: string;
		label?: string;
		id?: string;
		icon?: import('svelte').ComponentType;
		clearable?: boolean;
		disabled?: boolean;
		size?: 'sm' | 'md';
		searchable?: boolean;
		onselect?: (val: string) => void;
		class?: string;
	} = $props();

	let open = $state(false);
	let search = $state('');
	let filteredOptions = $derived(
		searchable && search.trim()
			? options.filter((o) => String(o[labelKey]).toLowerCase().includes(search.toLowerCase()))
			: options
	);
	let selectedLabel = $derived(
		options.find((o) => o[valueKey] === value)?.[labelKey] ?? ''
	);
	let listId = $derived(id || 'select-' + Math.random().toString(36).slice(2, 8));
	let listboxId = $derived(id ? id + '-listbox' : 'listbox-' + Math.random().toString(36).slice(2, 8));

	function toggle() {
		if (!disabled) {
			open = !open;
			if (open && searchable) {
				search = '';
			}
		}
	}

	function select(val: string) {
		value = val;
		open = false;
		onselect?.(val);
	}

	function clear(e: Event) {
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

	let triggerEl: HTMLButtonElement | undefined = $state();

	function portalDropdown(node: HTMLElement, getTrigger: () => HTMLButtonElement | undefined) {
		const reposition = () => {
			const trigger = getTrigger();
			if (!trigger) return;
			const rect = trigger.getBoundingClientRect();
			node.style.position = 'fixed';
			node.style.top = `${rect.bottom + 4}px`;
			node.style.left = `${rect.left}px`;
			node.style.width = `${trigger.offsetWidth}px`;
		};
		const onScroll = () => {
			if (document.body.contains(node)) reposition();
		};
		document.body.appendChild(node);
		reposition();
		window.addEventListener('scroll', onScroll, true);
		window.addEventListener('resize', onScroll);
		return {
			update() {
				reposition();
			},
			destroy() {
				node.remove();
				window.removeEventListener('scroll', onScroll, true);
				window.removeEventListener('resize', onScroll);
			}
		};
	}
</script>

<div class="relative {className}" role="combobox" aria-controls={listboxId} aria-expanded={open} aria-haspopup="listbox" tabindex="-1" onkeydown={handleKeydown}>
	{#if label}
		<label for={listId} class="block text-xs font-medium text-slate-500 mb-1.5">{label}</label>
	{/if}
	<button
		type="button"
		bind:this={triggerEl}
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
			{@const Icon = icon}
			<Icon size={16} class="text-slate-400 shrink-0" />
		{/if}
		<span class="flex-1 text-left truncate {selectedLabel ? '' : 'text-slate-400'}">
			{selectedLabel || placeholder}
		</span>
		{#if clearable && value}
			<span role="button" tabindex="0" onclick={clear} onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); clear(e); } }} class="p-0.5 rounded hover:bg-slate-200 text-slate-400 hover:text-slate-600 cursor-pointer" aria-label="Clear selection">
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
			id={listboxId}
			role="listbox"
			tabindex="-1"
			use:portalDropdown={() => triggerEl}
			class="z-50 mt-1 bg-white border border-slate-200 rounded-lg shadow-lg max-h-60 overflow-y-auto py-1"
		>
			{#if searchable}
				<li class="px-2 pb-1 sticky top-0 bg-white border-b border-slate-100">
					<div class="relative">
						<Search size={14} class="absolute left-2 top-1/2 -translate-y-1/2 text-slate-400" />
						<input
							type="text"
							bind:value={search}
							placeholder="Search..."
							class="w-full pl-7 pr-2 py-1.5 text-sm rounded-md border border-slate-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
						/>
					</div>
				</li>
			{/if}
			{#if !filteredOptions.length}
				<li class="px-3 py-2 text-sm text-slate-400 text-center">No options</li>
			{:else}
				{#each filteredOptions as opt (opt[valueKey])}
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
