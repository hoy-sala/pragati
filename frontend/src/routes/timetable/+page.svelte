<script lang="ts">
	import { WEEKLY_TIMETABLE, SUBJECT_INFO, WEEKDAY_TIMES, SAT_TIMES, DAY_LABELS, BREAK_CODES, ACTIVITY_CODES, BREAK_TIMES, getSegments } from './timetable.data';

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
							{@const nonBreakPeriods = schedule.days[showWeekday ? 0 : 5].periods.filter(p => !BREAK_CODES.has(p.code))}
							{@const cell = nonBreakPeriods[pi]}
							{@const isActivity = ACTIVITY_CODES.has(cell.code)}
							<th class="px-2 py-2.5 text-center font-semibold border-r border-slate-200 last:border-r-0 w-20 {isActivity ? 'text-slate-400' : 'text-slate-700'}">
								{#if isActivity}
									<div class="text-[11px]">{cell.name}</div>
								{:else}
									<div>P{nonBreakPeriods.slice(0, pi).filter(p => !ACTIVITY_CODES.has(p.code)).length + 1}</div>
								{/if}
								<div class="text-[10px] font-normal text-slate-400">{t}</div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each dayIndices as di}
						{@const day = schedule.days[di]}
						{@const segments = getSegments(day.periods)}
						{@const totalCols = times.length}
						{#each segments as segment, si}
							<tr class="border-t border-slate-200">
								{#if si === 0}
									<td class="sticky left-0 bg-white z-10 px-3 py-2 font-semibold text-slate-700 border-r border-slate-200" rowspan="{segments.length}">{DAY_LABELS[di]}</td>
								{/if}
								{#if segment.startCol > 0}
									<td colspan="{segment.startCol}"></td>
								{/if}
								{#each segment.periods as cell}
									{@const info = SUBJECT_INFO[cell.code]}
									{@const isActivity = ACTIVITY_CODES.has(cell.code)}
									<td class="px-2 py-2 text-center border-r border-slate-200 last:border-r-0 {isActivity ? 'bg-slate-50 text-slate-400 italic' : ''}"
										style="background-color: {isActivity ? '#F8FAFC' : (info?.color || '#fff')}">
										{#if isActivity}
											<div class="text-[11px]">{cell.name}</div>
										{:else}
											<div class="font-bold text-slate-800 text-xs">{cell.code}</div>
											<div class="text-[10px] text-slate-600 leading-tight">{cell.name}</div>
										{/if}
									</td>
								{/each}
									{#if totalCols - segment.startCol - segment.periods.length > 0}
									<td colspan="{totalCols - segment.startCol - segment.periods.length}" class="border-r border-slate-200"></td>
								{/if}
							</tr>
							{#if segment.breakAfter}
								{@const bt = BREAK_TIMES[segment.breakAfter.code]}
								<tr>
									<td colspan="{totalCols}" class="bg-slate-50 px-3 py-2 text-center border-b border-slate-200">
										<span class="font-medium text-slate-600 text-xs">{segment.breakAfter.name}</span>
										<span class="ml-2 text-slate-400 text-[10px]">{bt}</span>
									</td>
								</tr>
							{/if}
						{/each}
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
			<p>Friday P7–P8 reserved for Cultural Programme (CUL).</p>
			<p>P1 reserved for core academic subjects (Languages, Mathematics, Science, Social Studies).</p>
		</div>
	</div>
</div>
