<script lang="ts">
	import { getAuthState } from '$lib/stores/auth.svelte';
	import { GraduationCap } from 'lucide-svelte';
	import { userDisplayName, userInitials, effectiveRole, roleTitle } from '$lib/utils/nav';

	let { onmenu }: { onmenu: () => void } = $props();

	const auth = getAuthState();
	let displayName = $derived(userDisplayName(auth.currentUser));
	let initials = $derived(userInitials(displayName));
	let title = $derived(roleTitle(effectiveRole(auth.currentUser)));
</script>

<header class="md:hidden sticky top-0 z-40 bg-white/95 backdrop-blur border-b border-slate-200 no-print">
	<div class="flex items-center justify-between px-4 h-14">
		<div class="flex items-center gap-2.5 min-w-0">
			<div class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center shrink-0">
				<GraduationCap size={16} class="text-white" />
			</div>
			<div class="min-w-0">
				<div class="text-sm font-bold text-slate-800 font-kannada leading-tight">ಪ್ರಗತಿ</div>
				<div class="text-[10px] text-slate-400 leading-tight">{title}</div>
			</div>
		</div>
		<button
			onclick={onmenu}
			aria-label="Open navigation menu"
			class="w-9 h-9 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 text-xs font-semibold shrink-0 active:scale-95 transition-transform"
		>
			{initials}
		</button>
	</div>
</header>
