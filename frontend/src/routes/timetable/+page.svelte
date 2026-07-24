<script lang="ts">
	import { WEEKLY_TIMETABLE, SUBJECT_INFO, WEEKDAY_TIMES, SAT_TIMES, DAY_LABELS, BREAK_CODES, ACTIVITY_CODES } from './timetable.data';

	let activeClass = $state(0);
	let showAll = $state(false);
	let showWeekday = $state(true);

	const legend = Object.entries(SUBJECT_INFO);
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
			<button onclick={() => showAll = true}
				class="px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors
					{showAll ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-slate-600 border-slate-200 hover:border-primary-300'}">
				All Classes
			</button>
			{#each WEEKLY_TIMETABLE as _, i}
				<button onclick={() => { activeClass = i; showAll = false; }}
					class="px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors
						{!showAll && activeClass === i ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-slate-600 border-slate-200 hover:border-primary-300'}">
					Class {i + 6}
				</button>
			{/each}
		</div>

		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
			<table class="w-full text-xs">
				<thead>
					<tr class="bg-slate-100">
						<th class="sticky left-0 bg-slate-100 z-10 px-3 py-2.5 text-left font-semibold text-slate-700 border-r border-slate-200 w-20">Day</th>
						{#if showAll}
							<th class="sticky left-20 bg-slate-100 z-10 px-2 py-2.5 text-center font-semibold text-slate-700 border-r border-slate-200 w-12">Cls</th>
						{/if}
						{#each times as t, pi}
							{@const refPeriods = showWeekday ? 0 : 5}
							{@const cell = WEEKLY_TIMETABLE[0].days[refPeriods].periods[pi]}
							{@const isBreak = BREAK_CODES.has(cell.code)}
							{@const isActivity = ACTIVITY_CODES.has(cell.code)}
							<th class="px-2 py-2.5 text-center font-semibold border-r border-slate-200 last:border-r-0 w-20 {isBreak || isActivity ? 'text-slate-400' : 'text-slate-700'}">
								{#if isBreak || isActivity}
									<div class="text-[11px]">{cell.name}</div>
								{:else}
									<div>P{WEEKLY_TIMETABLE[0].days[refPeriods].periods.slice(0, pi).filter(p => !BREAK_CODES.has(p.code) && !ACTIVITY_CODES.has(p.code)).length + 1}</div>
								{/if}
								<div class="text-[10px] font-normal text-slate-400">{t}</div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#if showAll}
						{#each dayIndices as di}
							{@const classes = WEEKLY_TIMETABLE.length}
							{#each WEEKLY_TIMETABLE as cls, ci}
								{@const day = cls.days[di]}
								<tr class="border-t border-slate-200">
									{#if ci === 0}
										<td class="sticky left-0 bg-white z-10 px-3 py-2 font-semibold text-slate-700 border-r border-slate-200" rowspan="{classes}">{DAY_LABELS[di]}</td>
									{/if}
									<td class="sticky left-20 bg-white z-10 px-2 py-2 text-center font-semibold text-slate-600 border-r border-slate-200 text-xs">{ci + 6}</td>
									{#each day.periods.slice(0, times.length) as cell, pi}
										{@const info = SUBJECT_INFO[cell.code]}
										{@const isBreak = BREAK_CODES.has(cell.code)}
										{@const isActivity = ACTIVITY_CODES.has(cell.code)}
										{#if isBreak}
											{#if ci === 0}
												<td class="px-2 py-2 text-center border-r border-slate-200 last:border-r-0 bg-slate-50 text-slate-400 italic" rowspan="{classes}">
													<div class="text-[11px]">{cell.name}</div>
												</td>
											{/if}
										{:else}
											<td class="px-2 py-2 text-center border-r border-slate-200 last:border-r-0 {isActivity ? 'bg-slate-50 text-slate-400 italic' : ''}"
												style={isActivity ? '' : `background-color: ${info?.color || '#fff'}`}>
												{#if isActivity}
													<div class="text-[11px]">{cell.name}</div>
												{:else}
													<div class="font-bold text-slate-800 text-xs">{cell.code}</div>
													<div class="text-[10px] text-slate-600 leading-tight">{cell.name}</div>
												{/if}
											</td>
										{/if}
									{/each}
								</tr>
							{/each}
						{/each}
					{:else}
						{#each dayIndices as di, dayIdx}
							{@const day = WEEKLY_TIMETABLE[activeClass].days[di]}
							<tr class="border-t border-slate-200">
								{#if dayIdx === 0}
									<td class="sticky left-0 bg-white z-10 px-3 py-2 font-semibold text-slate-700 border-r border-slate-200" rowspan="{dayIndices.length}">{DAY_LABELS[di]}</td>
								{/if}
								{#each day.periods.slice(0, times.length) as cell, pi}
									{@const info = SUBJECT_INFO[cell.code]}
									{@const isBreak = BREAK_CODES.has(cell.code)}
									{@const isActivity = ACTIVITY_CODES.has(cell.code)}
									{#if isBreak}
										{#if dayIdx === 0}
											<td class="px-2 py-2 text-center border-r border-slate-200 last:border-r-0 bg-slate-50 text-slate-400 italic" rowspan="{dayIndices.length}">
												<div class="text-[11px]">{cell.name}</div>
											</td>
										{/if}
									{:else}
										<td class="px-2 py-2 text-center border-r border-slate-200 last:border-r-0 {isActivity ? 'bg-slate-50 text-slate-400 italic' : ''}"
											style={isActivity ? '' : `background-color: ${info?.color || '#fff'}`}>
											{#if isActivity}
												<div class="text-[11px]">{cell.name}</div>
											{:else}
												<div class="font-bold text-slate-800 text-xs">{cell.code}</div>
												<div class="text-[10px] text-slate-600 leading-tight">{cell.name}</div>
											{/if}
										</td>
									{/if}
								{/each}
							</tr>
						{/each}
					{/if}
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
