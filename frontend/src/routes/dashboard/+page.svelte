<script lang="ts">
	import { getAuthState } from '$lib/stores/auth.svelte';

	const auth = getAuthState();
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-slate-900">Dashboard</h1>
			<p class="text-sm text-slate-500 mt-1">
				Welcome back, {auth.currentUser?.name}
			</p>
		</div>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">--</div>
			<div class="text-sm text-slate-500 mt-1">Students</div>
		</div>
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">--</div>
			<div class="text-sm text-slate-500 mt-1">Teachers</div>
		</div>
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">--</div>
			<div class="text-sm text-slate-500 mt-1">Classes</div>
		</div>
		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<div class="text-2xl font-bold text-primary-600">--</div>
			<div class="text-sm text-slate-500 mt-1">Assessments</div>
		</div>
	</div>
</div>
