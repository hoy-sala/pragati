<script lang="ts">
	import { WEEKLY_TIMETABLE, SUBJECT_INFO, WEEKDAY_TIMES, SAT_TIMES, DAY_LABELS } from './timetable.data';

	let activeClass = $state(0);
	let showWeekday = $state(true);

	const legend = Object.entries(SUBJECT_INFO);
	let schedule = $derived(WEEKLY_TIMETABLE[activeClass]);
	let times = $derived(showWeekday ? WEEKDAY_TIMES : SAT_TIMES);
	let dayIndices = $derived(showWeekday ? [0, 1, 2, 3, 4] : [5]);
</script>

<svelte:head>
	<title>Master Timetable 2026-27 — Morarji Desai Residential School</title>
</svelte:head>

<div class="min-h-screen bg-slate-50">
	<div class="max-w-7xl mx-auto px-4 py-6 space-y-6">
		<div class="text-center space-y-1">
			<p class="text-xs font-medium text-slate-500 uppercase tracking-wider">Karnataka Residential Educational Institutions Society</p>
			<h1 class="text-xl font-bold text-slate-900">Morarji Desai Residential School (SC-32) Bahaddurghatta, Chitradurga</h1>
			<h2 class="text-2xl font-bold text-primary-700">Master School Time Table 2026-27</h2>
			<p class="text-sm text-slate-500">Monday – Friday: 10:00 AM – 4:20 PM &nbsp;|&nbsp; Saturday: 9:50 AM – 12:30 PM</p>
		</div>

		<div class="flex flex-wrap items-center justify-center gap-3">
			<div class="flex rounded-lg border border-slate-200 overflow-hidden bg-white">
				<button onclick={() => showWeekday = true}
					class="px-4 py-1.5 text-xs font-medium transition-colors {showWeekday ? 'bg-primary-600 text-white' : 'text-slate-600 hover:bg-slate-50'}">Mon–Fri</button>
				<button onclick={() => showWeekday = false}
					class="px-4 py-1.5 text-xs font-medium transition-colors {!showWeekday ? 'bg-primary-600 text-white' : 'text-slate-600 hover:bg-slate-50'}">Saturday</button>
			</div>
			{#each WEEKLY_TIMETABLE as _, i}
				<button onclick={() => activeClass = i}
					class="px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors
						{activeClass === i ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-slate-600 border-slate-200 hover:border-primary-300'}">
					Class {i + 6}
				</button>
			{/each}
		</div>

		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
			<table class="w-full text-xs">
				<thead>
					<tr class="bg-slate-100">
						<th class="sticky left-0 bg-slate-100 z-10 px-3 py-2.5 text-left font-semibold text-slate-700 border-r border-slate-200 w-20">Day</th>
						{#each times as t, pi}
							<th class="px-2 py-2.5 text-center font-semibold text-slate-700 border-r border-slate-200 last:border-r-0 w-20">
								<div>P{pi + 1}</div>
								<div class="text-[10px] font-normal text-slate-400">{t}</div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each dayIndices as di}
						{@const day = schedule.days[di]}
						<tr class="border-t border-slate-200">
							<td class="sticky left-0 bg-white z-10 px-3 py-2 font-semibold text-slate-700 border-r border-slate-200">{DAY_LABELS[di]}</td>
							{#each day.periods.slice(0, times.length) as cell, pi}
								{@const info = SUBJECT_INFO[cell.code]}
								<td class="px-2 py-2 text-center border-r border-slate-200 last:border-r-0 {pi === 3 && showWeekday ? 'border-b-2 border-amber-300' : ''} {pi === 5 && showWeekday ? 'border-b-2 border-amber-300' : ''}"
									style="background-color: {info?.color || '#fff'}">
									<div class="font-bold text-slate-800 text-xs">{cell.code}</div>
									<div class="text-[10px] text-slate-600 leading-tight">{cell.name}</div>
								</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-4">
			<h3 class="text-sm font-semibold text-slate-700 mb-2">Subject Legend</h3>
			<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-1.5">
				{#each legend as [code, info]}
					<div class="flex items-center gap-2 px-2 py-1 rounded text-xs" style="background-color: {info.color}">
						<span class="font-bold text-slate-800">{code}</span>
						<span class="text-slate-600">{info.name}</span>
					</div>
				{/each}
			</div>
		</div>

		<div class="text-center text-xs text-slate-400 pb-4">
			<p class="font-medium text-slate-500 mb-1">Notes</p>
			<p>P1 reserved for core academic subjects (Languages, Mathematics, Science, Social Studies).</p>
			<p>Friday P7–P8 reserved for Cultural Programme (CUL).</p>
			<p>Short Break: 12:00–12:10 PM &nbsp;|&nbsp; Lunch Break: 1:30–2:20 PM</p>
			<p class="mt-1">Saturday: Morning Assembly 8:30 AM &nbsp;|&nbsp; PT 8:40–9:10 AM &nbsp;|&nbsp; Breakfast 9:10–9:50 AM</p>
		</div>
	</div>
</div>
