<script lang="ts">
	import { getAuthState, logout } from '$lib/stores/auth.svelte';
	import { page } from '$app/stores';
	import { X, LogOut, GraduationCap } from 'lucide-svelte';
	import {
		visibleSections, isNavActive, allHrefs,
		userDisplayName, userInitials, roleTitle, effectiveRole
	} from '$lib/utils/nav';

	let { open = $bindable(false) }: { open?: boolean } = $props();

	const auth = getAuthState();
	let role = $derived(effectiveRole(auth.currentUser));
	let sections = $derived(visibleSections(role));
	let hrefs = $derived(allHrefs());

	function isActive(href: string): boolean {
		return isNavActive(href, $page.url.pathname, $page.url.search, hrefs);
	}

	function close() {
		open = false;
	}

	let lastUrl = $state('');
	$effect(() => {
		// Close the sheet whenever navigation happens (backstop; links close it directly).
		const url = $page.url.pathname + $page.url.search;
		if (lastUrl !== '' && url !== lastUrl) open = false;
		lastUrl = url;
	});
</script>

{#if open}
	<div class="md:hidden fixed inset-0 z-50 no-print" role="dialog" aria-modal="true" aria-label="Navigation menu">
		<button class="absolute inset-0 bg-slate-900/50 cursor-default" onclick={close} aria-label="Close menu" tabindex="-1"></button>
		<div class="absolute inset-x-0 bottom-0 bg-white rounded-t-3xl shadow-xl max-h-[85dvh] flex flex-col">
			<div class="flex items-center justify-between px-5 pt-4 pb-2">
				<div class="flex items-center gap-3 min-w-0">
					<div class="w-10 h-10 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 text-sm font-semibold shrink-0">
						{userInitials(userDisplayName(auth.currentUser))}
					</div>
					<div class="min-w-0">
						<div class="text-sm font-semibold text-slate-800 truncate">{userDisplayName(auth.currentUser)}</div>
						<div class="text-xs text-slate-500">{roleTitle(role)}</div>
					</div>
				</div>
				<button
					onclick={close}
					aria-label="Close menu"
					class="w-9 h-9 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 active:scale-95 transition-transform"
				>
					<X size={18} />
				</button>
			</div>
			<nav class="overflow-y-auto px-3 pb-3" style="padding-bottom: calc(0.75rem + env(safe-area-inset-bottom));">
				{#each sections as section}
					{#if section.label}
						<div class="px-2 mt-3 mb-1 text-[10px] font-semibold text-slate-400 uppercase tracking-wider">
							{section.label}
						</div>
					{/if}
					<div class="space-y-0.5">
						{#each section.items as item (item.href)}
							<a
								href={item.href}
								onclick={close}
								class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-[15px] transition-colors {isActive(item.href)
									? 'bg-primary-50 text-primary-700 font-medium'
									: 'text-slate-700 active:bg-slate-100'}"
							>
								<item.icon size={20} class={isActive(item.href) ? 'text-primary-600' : 'text-slate-400'} />
								<span class="truncate">{item.label}</span>
							</a>
						{/each}
					</div>
				{/each}
				<button
					onclick={() => { close(); logout(); }}
					class="mt-3 w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-[15px] text-danger-600 active:bg-danger-50 transition-colors"
				>
					<LogOut size={20} />
					<span>Sign out</span>
				</button>
				<div class="flex items-center justify-center gap-1.5 px-3 pt-3 pb-1 text-slate-400">
					<GraduationCap size={14} />
					<span class="text-[11px] font-kannada">ಪ್ರಗತಿ · PRAGATI</span>
				</div>
			</nav>
		</div>
	</div>
{/if}
