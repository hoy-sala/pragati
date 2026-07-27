<script lang="ts">
	import { WEEKLY_TIMETABLE, SUBJECT_INFO, WEEKDAY_TIMES, SAT_TIMES, DAY_LABELS, BREAK_CODES, ACTIVITY_CODES, TEACHER_NAMES } from './timetable.data';

	let activeClass = $state(0);
	let showAll = $state(true);

	const legend = Object.entries(SUBJECT_INFO);
	const MON_FRI_INDICES = [0, 1, 2, 3, 4];
	const SAT_INDICES = [5];
	const DAY_BG = ['bg-blue-50', 'bg-green-50', 'bg-amber-50', 'bg-purple-50', 'bg-pink-50'];
</script>

<svelte:head>
	<title>Time Table 2026-27 — Morarji Desai Residential School</title>
</svelte:head>

<style>
	@media print {
		@page {
			size: A4 landscape;
			margin: 8mm;
		}
		body { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
		.no-print { display: none !important; }
		.min-h-screen { min-height: auto !important; background: white !important; }
		.max-w-7xl { max-width: 100% !important; margin: 0 !important; padding: 0 !important; }
		.overflow-x-auto { overflow: visible !important; border: 1px solid #999 !important; border-radius: 0 !important; box-shadow: none !important; }
		table { font-size: 8px !important; width: 100% !important; }
		th, td { padding: 1px 2px !important; border-color: #999 !important; }
		.shadow-sm { box-shadow: none !important; }
		.rounded-xl { border-radius: 0 !important; }
		.bg-slate-50, .bg-slate-100 { background: #f5f5f5 !important; }
		.sticky { position: static !important; }
		.space-y-6 > :not(:last-child) { margin-bottom: 4mm; }
		.space-y-1 > :not(:last-child) { margin-bottom: 0; }
		h1 { font-size: 11pt !important; }
		h2 { font-size: 10pt !important; }
		p { font-size: 8pt !important; }
		.text-xs { font-size: 7pt !important; }
		.text-sm { font-size: 8pt !important; }
		.text-\[11px\] { font-size: 7pt !important; }
		.text-\[10px\] { font-size: 6.5pt !important; }
		.text-\[9px\] { font-size: 6pt !important; }
		.leading-tight { line-height: 1.1 !important; }
		.truncate { overflow: visible !important; max-width: none !important; }
		.min-h-screen > div:last-child { display: none !important; }
		table + div { display: none !important; }
	}
</style>

<div class="min-h-screen bg-slate-50">
	<div class="max-w-7xl mx-auto px-4 py-6 space-y-6">
		<div class="text-center space-y-1">
			<p class="text-xs font-medium text-slate-500 uppercase tracking-wider">Karnataka Residential Educational Institutions Society</p>
			<h1 class="text-xl font-bold text-slate-900">Morarji Desai Residential School (SC-32) Bahaddurghatta, Chitradurga</h1>
			<h2 class="text-2xl font-bold text-primary-700">Time Table 2026-27</h2>
			<p class="text-sm text-slate-500">Monday – Friday: 9:40 AM – 4:20 PM &nbsp;|&nbsp; Saturday: 8:30 AM – 12:30 PM</p>
		</div>

		<div class="flex flex-wrap items-center justify-center gap-3 no-print">
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

		<!-- Mon–Fri Table -->
		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
			<table class="w-full text-[11px]">
				<thead>
					<tr class="bg-slate-100">
						<th class="sticky left-0 bg-slate-100 z-10 px-2 py-1.5 text-left font-semibold text-slate-700 border-r border-slate-200 w-16">Day</th>
						{#if showAll}
							<th class="sticky left-16 bg-slate-100 z-10 px-1 py-1.5 text-center font-semibold text-slate-700 border-r border-slate-200 w-10">Class</th>
						{/if}
						{#each WEEKDAY_TIMES as t, pi}
							{@const cell = WEEKLY_TIMETABLE[0].days[0].periods[pi]}
							{@const isBreak = BREAK_CODES.has(cell.code)}
							{@const isActivity = ACTIVITY_CODES.has(cell.code)}
							<th class="px-1 py-1.5 text-center font-semibold border-r border-slate-200 last:border-r-0 w-16 {isBreak || isActivity ? 'text-slate-400' : 'text-slate-700'}">
								{#if isBreak || isActivity}
									<div class="text-[10px]">{cell.name}</div>
								{:else}
									<div>P{WEEKLY_TIMETABLE[0].days[0].periods.slice(0, pi).filter(p => !BREAK_CODES.has(p.code) && !ACTIVITY_CODES.has(p.code)).length + 1}</div>
								{/if}
								<div class="text-[9px] font-normal text-slate-400">{t}</div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#if showAll}
						{#each MON_FRI_INDICES as di}
							{@const classes = WEEKLY_TIMETABLE.length}
							{#each WEEKLY_TIMETABLE as cls, ci}
								{@const day = cls.days[di]}
								<tr class="border-t border-slate-200 {ci === classes - 1 ? 'border-b-2 border-slate-300' : ''}">
									{#if ci === 0}
										<td class="sticky left-0 z-10 px-2 py-1 font-semibold text-slate-700 border-r border-slate-200 {DAY_BG[di]}" rowspan="{classes}">{DAY_LABELS[di]}</td>
									{/if}
									<td class="sticky left-16 bg-white z-10 px-1 py-1 text-center font-semibold text-slate-600 border-r border-slate-200 text-[11px]">{ci + 6}</td>
									{#each day.periods.slice(0, WEEKDAY_TIMES.length) as cell, pi}
										{@const info = SUBJECT_INFO[cell.code]}
										{@const isBreak = BREAK_CODES.has(cell.code)}
										{@const isActivity = ACTIVITY_CODES.has(cell.code)}
										{#if isBreak || isActivity}
											{#if ci === 0}
												<td class="px-1 py-1 text-center border-r border-slate-200 last:border-r-0 bg-slate-50 text-slate-400 italic" rowspan="{classes}">
													<div class="text-[10px]">{cell.name}</div>
												</td>
											{/if}
										{:else}
											<td class="px-1 py-1 text-center align-middle border-r border-slate-200 last:border-r-0"
												style={`background-color: ${info?.color || '#fff'}`}>
												<div class="flex flex-col items-center justify-center leading-tight">
													<div class="font-bold text-slate-800 text-[11px]">{cell.code}</div>
													<div class="text-[9px] text-slate-400 truncate max-w-[70px]">{TEACHER_NAMES[cell.code] || cell.name}</div>
												</div>
											</td>
										{/if}
									{/each}
								</tr>
							{/each}
						{/each}
					{:else}
						{#each MON_FRI_INDICES as di}
							{@const day = WEEKLY_TIMETABLE[activeClass].days[di]}
							<tr class="border-t border-b border-slate-200 {DAY_BG[di]}">
								<td class="sticky left-0 z-10 px-2 py-1 font-semibold text-slate-700 border-r border-slate-200 {DAY_BG[di]}">{DAY_LABELS[di]}</td>
								{#each day.periods.slice(0, WEEKDAY_TIMES.length) as cell, pi}
									{@const info = SUBJECT_INFO[cell.code]}
									{@const isBreak = BREAK_CODES.has(cell.code)}
									{@const isActivity = ACTIVITY_CODES.has(cell.code)}
									{#if isBreak || isActivity}
										<td class="px-1 py-1 text-center border-r border-slate-200 last:border-r-0 bg-slate-50 text-slate-400 italic">
											<div class="text-[10px]">{cell.name}</div>
										</td>
									{:else}
										<td class="px-1 py-1 text-center align-middle border-r border-slate-200 last:border-r-0 {isActivity ? 'bg-slate-50 text-slate-400 italic' : ''}"
											style={isActivity ? '' : `background-color: ${info?.color || '#fff'}`}>
											{#if isActivity}
												<div class="text-[10px]">{cell.name}</div>
											{:else}
												<div class="flex flex-col items-center justify-center leading-tight">
													<div class="font-bold text-slate-800 text-[11px]">{cell.code}</div>
													<div class="text-[9px] text-slate-400 truncate max-w-[70px]">{TEACHER_NAMES[cell.code] || cell.name}</div>
												</div>
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

		<!-- Saturday Table -->
		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
			<table class="w-full text-[11px]">
				<thead>
					<tr class="bg-slate-100">
						<th class="sticky left-0 bg-slate-100 z-10 px-2 py-1.5 text-left font-semibold text-slate-700 border-r border-slate-200 w-16">Day</th>
						{#if showAll}
							<th class="sticky left-16 bg-slate-100 z-10 px-1 py-1.5 text-center font-semibold text-slate-700 border-r border-slate-200 w-10">Class</th>
						{/if}
						{#each SAT_TIMES as t, pi}
							{@const cell = WEEKLY_TIMETABLE[0].days[5].periods[pi]}
							{@const isBreak = BREAK_CODES.has(cell.code)}
							{@const isActivity = ACTIVITY_CODES.has(cell.code)}
							<th class="px-1 py-1.5 text-center font-semibold border-r border-slate-200 last:border-r-0 w-16 {isBreak || isActivity ? 'text-slate-400' : 'text-slate-700'}">
								{#if isBreak || isActivity}
									<div class="text-[10px]">{cell.name}</div>
								{:else}
									<div>P{WEEKLY_TIMETABLE[0].days[5].periods.slice(0, pi).filter(p => !BREAK_CODES.has(p.code) && !ACTIVITY_CODES.has(p.code)).length + 1}</div>
								{/if}
								<div class="text-[9px] font-normal text-slate-400">{t}</div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#if showAll}
						{#each SAT_INDICES as di}
							{@const classes = WEEKLY_TIMETABLE.length}
							{#each WEEKLY_TIMETABLE as cls, ci}
								{@const day = cls.days[di]}
								<tr class="border-t border-slate-200 {ci === classes - 1 ? 'border-b-2 border-slate-300' : ''}">
									{#if ci === 0}
										<td class="sticky left-0 z-10 px-2 py-1 font-semibold text-slate-700 border-r border-slate-200 {DAY_BG[di]}" rowspan="{classes}">{DAY_LABELS[di]}</td>
									{/if}
									<td class="sticky left-16 bg-white z-10 px-1 py-1 text-center font-semibold text-slate-600 border-r border-slate-200 text-[11px]">{ci + 6}</td>
									{#each day.periods.slice(0, SAT_TIMES.length) as cell, pi}
										{@const info = SUBJECT_INFO[cell.code]}
										{@const isBreak = BREAK_CODES.has(cell.code)}
										{@const isActivity = ACTIVITY_CODES.has(cell.code)}
										{#if isBreak || isActivity}
											{#if ci === 0}
												<td class="px-1 py-1 text-center border-r border-slate-200 last:border-r-0 bg-slate-50 text-slate-400 italic" rowspan="{classes}">
													<div class="text-[10px]">{cell.name}</div>
												</td>
											{/if}
										{:else}
											<td class="px-1 py-1 text-center align-middle border-r border-slate-200 last:border-r-0"
												style={`background-color: ${info?.color || '#fff'}`}>
												<div class="flex flex-col items-center justify-center leading-tight">
													<div class="font-bold text-slate-800 text-[11px]">{cell.code}</div>
													<div class="text-[9px] text-slate-400 truncate max-w-[70px]">{TEACHER_NAMES[cell.code] || cell.name}</div>
												</div>
											</td>
										{/if}
									{/each}
								</tr>
							{/each}
						{/each}
					{:else}
						{#each SAT_INDICES as di}
							{@const day = WEEKLY_TIMETABLE[activeClass].days[di]}
							<tr class="border-t border-b border-slate-200 {DAY_BG[di]}">
								<td class="sticky left-0 z-10 px-2 py-1 font-semibold text-slate-700 border-r border-slate-200 {DAY_BG[di]}">{DAY_LABELS[di]}</td>
								{#each day.periods.slice(0, SAT_TIMES.length) as cell, pi}
									{@const info = SUBJECT_INFO[cell.code]}
									{@const isBreak = BREAK_CODES.has(cell.code)}
									{@const isActivity = ACTIVITY_CODES.has(cell.code)}
									{#if isBreak || isActivity}
										<td class="px-1 py-1 text-center border-r border-slate-200 last:border-r-0 bg-slate-50 text-slate-400 italic">
											<div class="text-[10px]">{cell.name}</div>
										</td>
									{:else}
										<td class="px-1 py-1 text-center align-middle border-r border-slate-200 last:border-r-0 {isActivity ? 'bg-slate-50 text-slate-400 italic' : ''}"
											style={isActivity ? '' : `background-color: ${info?.color || '#fff'}`}>
											{#if isActivity}
												<div class="text-[10px]">{cell.name}</div>
											{:else}
												<div class="flex flex-col items-center justify-center leading-tight">
													<div class="font-bold text-slate-800 text-[11px]">{cell.code}</div>
													<div class="text-[9px] text-slate-400 truncate max-w-[70px]">{TEACHER_NAMES[cell.code] || cell.name}</div>
												</div>
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

		<div class="bg-white rounded-xl border border-slate-200 p-4 no-print">
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
	</div>
</div>