<script lang="ts">
	import { LoaderCircle } from 'lucide-svelte';

	let {
		children,
		variant = 'primary' as 'primary' | 'secondary' | 'danger' | 'ghost',
		size = 'md' as 'sm' | 'md' | 'lg',
		icon,
		disabled = false,
		loading = false,
		type = 'button' as 'button' | 'submit' | 'reset',
		class: className = '',
		onclick,
		...rest
	}: {
		children?: import('svelte').Snippet;
		variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
		size?: 'sm' | 'md' | 'lg';
		icon?: import('svelte').ComponentType;
		disabled?: boolean;
		loading?: boolean;
		type?: 'button' | 'submit' | 'reset';
		class?: string;
		onclick?: (e: MouseEvent) => void;
		[key: string]: unknown;
	} = $props();

	const base = 'inline-flex items-center justify-center gap-2 font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-1 disabled:opacity-50 disabled:pointer-events-none';

	const variants: Record<string, string> = {
		primary: 'bg-primary-600 text-white hover:bg-primary-700 active:bg-primary-800',
		secondary: 'bg-slate-100 text-slate-700 hover:bg-slate-200 active:bg-slate-300 border border-slate-300',
		danger: 'bg-danger-500 text-white hover:bg-red-600 active:bg-red-700',
		ghost: 'text-slate-600 hover:bg-slate-100 hover:text-slate-900 active:bg-slate-200',
	};

	const sizes: Record<string, string> = {
		sm: 'px-2.5 py-1.5 text-xs rounded-md',
		md: 'px-3.5 py-2 text-sm rounded-lg',
		lg: 'px-5 py-2.5 text-base rounded-lg',
	};
</script>

<button
	{type}
	{disabled}
	{onclick}
	class="{base} {variants[variant]} {sizes[size]} {className}"
	{...rest}
>
	{#if loading}
		<LoaderCircle class="animate-spin" size={size === 'sm' ? 14 : size === 'lg' ? 20 : 16} />
	{:else if icon}
		{@const Icon = icon}
		<Icon size={size === 'sm' ? 14 : size === 'lg' ? 20 : 16} />
	{/if}
	{@render children?.()}
</button>
