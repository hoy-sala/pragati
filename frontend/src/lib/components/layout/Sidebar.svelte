<script lang="ts">
	import { getAuthState, logout } from '$lib/stores/auth.svelte';
	import { page } from '$app/stores';
	import { LogOut, GraduationCap } from 'lucide-svelte';
	import {
		NAV_SECTIONS, visibleSections, isNavActive, allHrefs,
		userDisplayName, userInitials, roleTitle, effectiveRole,
		type NavItem
	} from '$lib/utils/nav';

	const auth = getAuthState();

	let role = $derived(effectiveRole(auth.currentUser));
	let displayName = $derived(userDisplayName(auth.currentUser));
	let initials = $derived(userInitials(displayName));
	let displayTitle = $derived(roleTitle(role));

	let roleBadgeColor = $derived(
		role === 'admin' ? 'bg-purple-100 text-purple-700' :
		role === 'principal' ? 'bg-blue-100 text-blue-700' :
		role === 'teacher' ? 'bg-amber-100 text-amber-700' :
		role === 'student' ? 'bg-emerald-100 text-emerald-700' :
		'bg-slate-100 text-slate-600'
	);

	let sections = $derived(visibleSections(role));
	let hrefs = $derived(allHrefs(NAV_SECTIONS));

	function isActive(href: string): boolean {
		return isNavActive(href, $page.url.pathname, $page.url.search, hrefs);
	}

	function isSectionActive(items: NavItem[]): boolean {
		return items.some(item => isActive(item.href));
	}
</script>

<aside class="w-60 bg-white border-r border-slate-200 hidden md:flex flex-col h-full no-print">
	<div class="px-4 py-4 border-b border-slate-200">
		<div class="flex items-center gap-3">
			<div class="w-9 h-9 rounded-lg bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shrink-0">
				<GraduationCap size={18} class="text-white" />
			</div>
			<div class="min-w-0">
				<h1 class="text-sm font-bold text-slate-800 font-kannada tracking-wide leading-tight">ಪ್ರಗತಿ</h1>
				<p class="text-[10px] text-slate-400 leading-tight">PRAGATI v1.0</p>
			</div>
		</div>
	</div>

	<nav class="flex-1 overflow-y-auto px-2 py-3 space-y-5 scrollbar-none">
		{#each sections as section}
			<div>
				{#if section.label}
					<div class="px-2 mb-1.5 text-[10px] font-semibold text-slate-400 uppercase tracking-wider">
						{section.label}
					</div>
				{/if}
				<div class="space-y-0.5">
					{#each section.items as item (item.href)}
						<a
							href={item.href}
							class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-md text-[13px] transition-colors {isActive(item.href)
								? 'bg-primary-50 text-primary-700 font-medium'
								: 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
						>
							<item.icon size={16} class={isActive(item.href) ? 'text-primary-600' : 'text-slate-400'} />
							<span class="truncate">{item.label}</span>
						</a>
					{/each}
				</div>
			</div>
		{/each}
	</nav>

	<div class="border-t border-slate-200 p-2.5">
		<div class="flex items-center gap-2.5 px-2 py-1.5">
			<div class="w-8 h-8 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 text-xs font-semibold shrink-0">
				{initials}
			</div>
			<div class="min-w-0 flex-1">
				<div class="text-sm font-medium text-slate-800 truncate leading-tight">{displayName}</div>
				<div class="flex items-center gap-1.5 mt-0.5">
					<span class="inline-block w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
					<span class="text-[11px] text-slate-500">{displayTitle}</span>
				</div>
			</div>
			<button
				onclick={logout}
				title="Sign out"
				class="p-1.5 rounded-md text-slate-400 hover:text-danger-600 hover:bg-danger-50 transition-colors shrink-0"
			>
				<LogOut size={15} />
			</button>
		</div>
	</div>
</aside>
