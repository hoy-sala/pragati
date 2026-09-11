<script lang="ts">
	import { page } from '$app/stores';
	import type { PageTab } from '$lib/utils/tabs';

	let { tabs, role = '', active }: { tabs: PageTab[]; role?: string; active?: string } = $props();

	let visible = $derived(tabs.filter((t) => !t.roles || t.roles.includes(role)));

	function currentUrl(): string {
		return $page.url.pathname + $page.url.search;
	}

	function splitHref(href: string): { path: string; search: string } {
		const i = href.indexOf('?');
		return i === -1 ? { path: href, search: '' } : { path: href.slice(0, i), search: href.slice(i) };
	}

	function isActive(href: string): boolean {
		if (active !== undefined) return active === href;
		const cur = currentUrl();
		if (cur === href) return true;
		const { path, search } = splitHref(href);
		const curPath = $page.url.pathname;
		if (search) return false; // query tabs only match exactly
		if (curPath === path) return true;
		if (!curPath.startsWith(path + '/')) return false;
		// Longest-prefix wins so /mentors doesn't stay active on /mentors/roster.
		return !visible.some(
			(t) => {
				const tp = splitHref(t.href);
				return (
					t.href !== href &&
					tp.search === '' &&
					tp.path.length > path.length &&
					(curPath === tp.path || curPath.startsWith(tp.path + '/'))
				);
			}
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
