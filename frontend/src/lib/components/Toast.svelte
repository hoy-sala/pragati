<script lang="ts">
	import { toastStore } from '$lib/stores/toast.svelte';
	import { CheckCircle2, XCircle, Info, X } from 'lucide-svelte';

	const icons = { success: CheckCircle2, error: XCircle, info: Info };
	const colors = {
		success: 'text-emerald-600',
		error: 'text-red-600',
		info: 'text-blue-600'
	};
</script>

<div class="fixed bottom-4 right-4 z-[70] flex flex-col gap-2 no-print" aria-live="polite">
	{#each toastStore.toasts as t (t.id)}
		{@const Icon = icons[t.type]}
		<div class="flex items-center gap-2.5 bg-white border border-slate-200 rounded-lg shadow-lg px-3.5 py-2.5 min-w-64 max-w-sm toast-in">
			<Icon size={16} class="shrink-0 {colors[t.type]}" />
			<p class="text-sm text-slate-700 flex-1">{t.message}</p>
			<button onclick={() => toastStore.dismiss(t.id)} class="text-slate-300 hover:text-slate-500 shrink-0" aria-label="Dismiss">
				<X size={14} />
			</button>
		</div>
	{/each}
</div>

<style>
	.toast-in {
		animation: toastIn 0.18s ease-out;
	}
	@keyframes toastIn {
		from { opacity: 0; transform: translateY(8px); }
		to { opacity: 1; transform: translateY(0); }
	}
</style>
