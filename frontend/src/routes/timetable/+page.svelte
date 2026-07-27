<script lang="ts">
	import { WEEKLY_TIMETABLE, SUBJECT_INFO, WEEKDAY_TIMES, SAT_TIMES, DAY_LABELS, BREAK_CODES, ACTIVITY_CODES, TEACHER_NAMES } from './timetable.data';

	let activeClass = $state(0);
	let showAll = $state(true);
	let viewMode = $state<'class' | 'subject'>('class');
	let selectedSubject = $state('KAN');

	const legend = Object.entries(SUBJECT_INFO);
	const MON_FRI_INDICES = [0, 1, 2, 3, 4];
	const SAT_INDICES = [5];
	const DAY_BG = ['bg-blue-50', 'bg-green-50', 'bg-amber-50', 'bg-purple-50', 'bg-pink-50'];

	const ACADEMIC_SUBJECTS = Object.entries(SUBJECT_INFO).filter(([c]) => !BREAK_CODES.has(c) && !ACTIVITY_CODES.has(c));

	const WEEKDAY_PERIOD_MAP = [1, 2, 3, 5, 6, 8, 9, 10];
	const SAT_PERIOD_MAP = [3, 4, 6, 7];
	const WEEKDAY_SLOT_LABELS = ['P1', 'P2', 'P3', 'P4', 'P5', 'P6', 'P7', 'P8'];
	const SAT_SLOT_LABELS = ['P1', 'P2', 'P3', 'P4'];

	type SubjectGrid = Record<string, (number | null)[]>;

	function buildSubjectGrid(): Record<string, SubjectGrid> {
		const grid: Record<string, SubjectGrid> = {};
		const allCodes = new Set(Object.keys(SUBJECT_INFO));
		for (const code of allCodes) {
			if (BREAK_CODES.has(code) || ACTIVITY_CODES.has(code)) continue;
			grid[code] = {};
			for (let d = 0; d < 6; d++) {
				const map = d < 5 ? WEEKDAY_PERIOD_MAP : SAT_PERIOD_MAP;
				grid[code]['d' + d] = new Array(map.length).fill(null);
			}
		}
		WEEKLY_TIMETABLE.forEach((cls, ci) => {
			const classNum = ci + 6;
			cls.days.forEach((day, di) => {
				const map = di < 5 ? WEEKDAY_PERIOD_MAP : SAT_PERIOD_MAP;
				map.forEach((pi, slot) => {
					const cell = day.periods[pi];
					if (cell && grid[cell.code]) {
						grid[cell.code]['d' + di][slot] = classNum;
					}
				});
			});
		});
		return grid;
	}

	const subjectGrid = buildSubjectGrid();
</script>

<svelte:head>
	<title>Time Table 2026-27 — Morarji Desai Residential School</title>
</svelte:head>

<style>
	.print-only { display: none; }
	@media print {
		@page { size: A4 landscape; margin: 6mm; }
		@page portrait-page { size: A4 portrait; margin: 5mm; }
		:global(body) { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
		.no-print, .legend-print { display: none !important; }
		.print-only { display: inline !important; }
		.subject-print { page: portrait-page; }
		.min-h-screen { min-height: auto !important; }
		.max-w-7xl { max-width: 100% !important; margin: 0 !important; padding: 0 !important; }
		.overflow-x-auto { box-shadow: none !important; margin-bottom: 3px !important; }
		table { font-size: 7.5px !important; width: 100% !important; }
		th, td { padding: 1.5px 3px !important; }
		.space-y-6 > :not(:last-child) { margin-bottom: 3px; }
		h1 { font-size: 10pt !important; margin: 0 0 1px 0 !important; }
		h2 { font-size: 10pt !important; margin: 0 0 1px 0 !important; }
		.text-xs { font-size: 6.5px !important; }
		.text-sm { font-size: 7.5px !important; }
		.text-\[11px\] { font-size: 7.5px !important; }
		.text-\[10px\] { font-size: 7px !important; }
		.text-\[9px\] { font-size: 6.5px !important; }
		.leading-tight { line-height: 1.1 !important; }
		.w-16 { width: auto !important; min-width: 22px !important; }
		.w-10 { width: auto !important; min-width: 14px !important; }
		.w-14 { width: auto !important; min-width: 18px !important; padding: 1px 2px !important; }
		.sticky { position: static !important; }
		.border-r { border-right-width: 1px !important; }
	}
</style>

<div class="min-h-screen bg-slate-50">
	<div class="max-w-7xl mx-auto px-4 py-6 space-y-6 {viewMode === 'subject' ? 'subject-print' : ''}">
		<div class="text-center space-y-1">
			<p class="text-xs font-medium text-slate-500 uppercase tracking-wider">Karnataka Residential Educational Institutions Society</p>
			<h1 class="text-xl font-bold text-slate-900">Morarji Desai Residential School (SC-32) Bahaddurghatta, Chitradurga</h1>
			<h2 class="text-2xl font-bold text-primary-700">Time Table 2026-27 <span class="print-only text-slate-500 inline">
				- {#if viewMode === 'class'}
					{showAll ? 'All Classes' : 'Class ' + (activeClass + 6)} - Class Wise
				{:else}
					{SUBJECT_INFO[selectedSubject]?.name} - Subject Wise
				{/if}
			</span></h2>
		</div>

		<div class="flex flex-wrap items-center justify-center gap-3 no-print">
			<button onclick={() => viewMode = 'class'}
				class="px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors
					{viewMode === 'class' ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-slate-600 border-slate-200 hover:border-primary-300'}">
				Class Wise
			</button>
			<button onclick={() => viewMode = 'subject'}
				class="px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors
					{viewMode === 'subject' ? 'bg-primary-600 text-white border-primary-600' : 'bg-white text-slate-600 border-slate-200 hover:border-primary-300'}">
				Subject Wise
			</button>
			<span class="w-px h-5 bg-slate-300"></span>
			{#if viewMode === 'class'}
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
			{:else}
				<select bind:value={selectedSubject}
					class="px-3 py-1.5 rounded-lg border border-slate-300 text-xs font-medium bg-white text-slate-700">
					{#each ACADEMIC_SUBJECTS as [code, info]}
						<option value={code}>{code} — {info.name}</option>
					{/each}
				</select>
			{/if}
		</div>

		{#if viewMode === 'class'}
		<!-- Mon-Fri Table -->
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

		{:else}
		{@const grid = subjectGrid[selectedSubject]}
		{@const info = SUBJECT_INFO[selectedSubject]}
		{@const teacherName = TEACHER_NAMES[selectedSubject]}

		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
			<table class="w-full text-[11px]">
				<thead>
					<tr class="bg-slate-100">
						<th class="px-2 py-1.5 text-left font-semibold text-slate-700 border-r border-slate-200 w-16">Day</th>
						{#each WEEKDAY_SLOT_LABELS as label, si}
							<th class="px-1 py-1.5 text-center font-semibold text-slate-700 border-r border-slate-200 last:border-r-0 w-14">
								<div>{label}</div>
								<div class="text-[9px] font-normal text-slate-400">{WEEKDAY_TIMES[WEEKDAY_PERIOD_MAP[si]]}</div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each MON_FRI_INDICES as di}
						<tr class="border-t border-slate-200 {DAY_BG[di]}">
							<td class="px-2 py-2 font-semibold text-slate-700 border-r border-slate-200 {DAY_BG[di]}">{DAY_LABELS[di]}</td>
							{#each WEEKDAY_SLOT_LABELS as _, si}
								{@const cls = grid['d' + di][si]}
								<td class="px-1 py-2 text-center align-middle border-r border-slate-200 last:border-r-0" class:bg-slate-50={!cls}>
									{#if cls}
										<span class="font-bold text-slate-800 text-sm">{cls}</span>
									{:else}
										<span class="text-slate-300">&ndash;</span>
									{/if}
								</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div class="overflow-x-auto rounded-xl border border-slate-200 bg-white shadow-sm">
			<table class="w-full text-[11px]">
				<thead>
					<tr class="bg-slate-100">
						<th class="px-2 py-1.5 text-left font-semibold text-slate-700 border-r border-slate-200 w-16">Day</th>
						{#each SAT_SLOT_LABELS as label, si}
							<th class="px-1 py-1.5 text-center font-semibold text-slate-700 border-r border-slate-200 last:border-r-0 w-14">
								<div>{label}</div>
								<div class="text-[9px] font-normal text-slate-400">{SAT_TIMES[SAT_PERIOD_MAP[si]]}</div>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					<tr class="border-t border-slate-200 bg-amber-50">
						<td class="px-2 py-2 font-semibold text-slate-700 border-r border-slate-200">{DAY_LABELS[5]}</td>
						{#each SAT_SLOT_LABELS as _, si}
							{@const cls = grid['d5'][si]}
							<td class="px-1 py-2 text-center align-middle border-r border-slate-200 last:border-r-0" class:bg-slate-50={!cls}>
								{#if cls}
									<span class="font-bold text-slate-800 text-sm">{cls}</span>
								{:else}
									<span class="text-slate-300">&ndash;</span>
								{/if}
							</td>
						{/each}
					</tr>
				</tbody>
			</table>
		</div>

		<div class="flex items-center gap-4 text-xs text-slate-500">
			<span class="font-semibold text-slate-700">{selectedSubject}</span>
			<span>{info?.name}</span>
			{#if teacherName}
				<span class="text-slate-400">Teacher: {teacherName}</span>
			{/if}
		</div>
		{/if}

		<div class="bg-white rounded-xl border border-slate-200 p-4 legend-print">
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