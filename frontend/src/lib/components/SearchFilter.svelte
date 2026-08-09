<script lang="ts">
	let {
		value = '',
		placeholder = 'Search...',
		onInput
	}: {
		value?: string;
		placeholder?: string;
		onInput: (value: string) => void;
	} = $props();

	let inputValue = $state(value);

	$effect(() => { inputValue = value; });

	let debounceTimer: ReturnType<typeof setTimeout>;

	function handleInput(e: Event) {
		inputValue = (e.target as HTMLInputElement).value;
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => onInput(inputValue), 300);
	}
</script>

<div class="relative">
	<svg class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
	<input
		type="text"
		bind:value={inputValue}
		{placeholder}
		oninput={handleInput}
		class="w-full pl-9 pr-3 py-2 rounded-lg border border-slate-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-400"
	/>
</div>
