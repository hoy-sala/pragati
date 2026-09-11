<script lang="ts">
	import { page } from '$app/stores';
	import type { PageTab } from '$lib/utils/tabs';

	let { tabs, role = '' }: { tabs: PageTab[]; role?: string } = $props();

	let visible = $derived(tabs.filter((t) => !t.roles || t.roles.includes(role)));

	function isActive(href: string): boolean {
		const path = $page.url.pathname;
		if (path === href) return true;
		if (!path.startsWith(href + '/')) return false;
		// Longest-prefix wins so /mentors doesn't stay active on /mentors/roster.
		return !visible.some(
			(t) =>
				t.href !== href &&
				t.href.length > href.length &&
				(path === t.href || path.startsWith(t.href + '/'))
		);
	}
</script>

{#if visible.length > 1}
	<nav class="flex gap-1 bg-slate-100 rounded-lg p-1 no-print" aria-label="Section navigation">
		{#each visible as t (t.href)}
			<a
				href={t.href}
				aria-current={isActive(t.href) ? 'page' : undefined}
				class="px-4 py-2 text-sm font-medium rounded-md transition-colors {isActive(t.href)
					? 'bg-white text-slate-900 shadow-sm'
					: 'text-slate-500 hover:text-slate-700'}"
			>
				{t.label}
			</a>
		{/each}
	</nav>
{/if}
