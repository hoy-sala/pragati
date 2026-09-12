<script lang="ts">
	import { getAuthState } from '$lib/stores/auth.svelte';
	import { page } from '$app/stores';
	import { LayoutGrid } from 'lucide-svelte';
	import { bottomTabs, isNavActive, allHrefs, effectiveRole } from '$lib/utils/nav';

	let { onmore }: { onmore: () => void } = $props();

	const auth = getAuthState();
	let role = $derived(effectiveRole(auth.currentUser));
	let tabs = $derived(bottomTabs(role));
	let hrefs = $derived(allHrefs());

	function isActive(href: string): boolean {
		return isNavActive(href, $page.url.pathname, $page.url.search, hrefs);
	}
</script>

<nav class="md:hidden bg-white border-t border-slate-200 no-print" aria-label="Primary">
	<div class="grid grid-cols-5 px-1 pt-1" style="padding-bottom: calc(0.375rem + env(safe-area-inset-bottom));">
		{#each tabs as t (t.href)}
			<a
				href={t.href}
				aria-current={isActive(t.href) ? 'page' : undefined}
				class="flex flex-col items-center gap-1 py-1.5 rounded-lg min-h-[52px] justify-center transition-colors {isActive(t.href) ? 'text-primary-700' : 'text-slate-400 active:text-slate-600'}"
			>
				<t.icon size={22} strokeWidth={isActive(t.href) ? 2.4 : 2} />
				<span class="text-[10px] font-medium leading-none {isActive(t.href) ? '' : ''}">{t.label}</span>
				<span class="h-1 w-8 rounded-full {isActive(t.href) ? 'bg-primary-600' : 'bg-transparent'}"></span>
			</a>
		{/each}
		<button
			onclick={onmore}
			aria-label="More navigation options"
			class="flex flex-col items-center gap-1 py-1.5 rounded-lg min-h-[52px] justify-center text-slate-400 active:text-slate-600 transition-colors"
		>
			<LayoutGrid size={22} />
			<span class="text-[10px] font-medium leading-none">More</span>
			<span class="h-1 w-8 rounded-full bg-transparent"></span>
		</button>
	</div>
</nav>
