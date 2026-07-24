<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Subject } from '$lib/types';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	const studentId = $page.params.student_id;
	const term = $page.url.searchParams.get('term') || 'Term1';
	const year = $page.url.searchParams.get('year') || '';
	let studentName = $state('');

	let subjects = $state<Subject[]>([]);
	let scholastic = $state<any[]>([]);
	let entryId = $state('');
	let entryVersion = $state(1);
	let loading = $state(true);
	let saving = $state(false);
	let statusMsg = $state('');

	// Co-scholastic areas
	let coScholastic = $state<any>({
		life_skills: { thinking: { label: 'Thinking Skills', grade: '' }, social: { label: 'Social Skills', grade: '' }, emotional: { label: 'Emotional Skills', grade: '' } },
		attitudes: { towards_teachers: { label: 'Towards Teachers', grade: '' }, towards_schoolmates: { label: 'Towards Schoolmates', grade: '' }, towards_school: { label: 'Towards School Programs', grade: '' }, towards_environment: { label: 'Towards Environment', grade: '' } },
		values: { honesty: { label: 'Honesty', grade: '' }, responsibility: { label: 'Responsibility', grade: '' }, cooperation: { label: 'Cooperation', grade: '' }, discipline: { label: 'Discipline', grade: '' } },
		participation: { sports: { label: 'Sports & Games', grade: '' }, arts: { label: 'Arts & Crafts', grade: '' }, music: { label: 'Music & Dance', grade: '' }, clubs: { label: 'Clubs & Societies', grade: '' } }
	});

	// Health & PE
	let healthPE = $state<any>({
		height: { label: 'Height (cm)', value: '' },
		weight: { label: 'Weight (kg)', value: '' },
		bmi: { label: 'BMI', value: '' },
		vision: { label: 'Vision', value: '' },
		dental: { label: 'Dental Health', value: '' },
		physical_fitness: { label: 'Physical Fitness', value: '' }
	});

	// Work education
	let workEducation = $state<any>({ activity: '', description: '' });

	// Self / Peer / Parent
	let selfAssessment = $state<any>({ strong_areas: '', areas_for_improvement: '', interests: '' });
	let peerAssessment = $state<any>({ teamwork: '', behavior: '', helpfulness: '' });
	let parentFeedback = $state<any>({ academic_progress: '', behavior_at_home: '', suggestions: '' });

	// Attendance
	let attendance = $state({ present: 0, total: 0, percentage: 0 });

	let teacherRemarks = $state('');

	const gradeOptions = ['', 'A+', 'A', 'B+', 'B', 'C', 'D', 'E'];

	const coScholasticSections = $derived.by(() => ({
		life_skills: { label: 'Life Skills', items: Object.keys(coScholastic.life_skills || {}).map(k => ({ key: k })) },
		attitudes: { label: 'Attitudes & Behaviours', items: Object.keys(coScholastic.attitudes || {}).map(k => ({ key: k })) },
		values: { label: 'Values', items: Object.keys(coScholastic.values || {}).map(k => ({ key: k })) },
		participation: { label: 'Participation & Co-curricular', items: Object.keys(coScholastic.participation || {}).map(k => ({ key: k })) }
	}));

	onMount(async () => {
		const [sRes] = await Promise.all([
			api<Subject[]>('GET', '/subjects?limit=50')
		]);
		if (sRes.data) subjects = sRes.data;

		if (studentId) {
			const res = await api<any>('GET', `/hpc/entries?student_id=${studentId}&term=${term}&academic_year_id=${year}`);
			if (res.data) {
				const d = res.data;
				studentName = d.student_name || '';
				entryId = d.id;
				entryVersion = d.version;
				scholastic = d.scholastic || [];
				if (d.co_scholastic && Object.keys(d.co_scholastic).length > 0) coScholastic = d.co_scholastic;
				if (d.health_pe && Object.keys(d.health_pe).length > 0) healthPE = d.health_pe;
				if (d.work_education && Object.keys(d.work_education).length > 0) workEducation = d.work_education;
				if (d.self_assessment && Object.keys(d.self_assessment).length > 0) selfAssessment = d.self_assessment;
				if (d.peer_assessment && Object.keys(d.peer_assessment).length > 0) peerAssessment = d.peer_assessment;
				if (d.parent_feedback && Object.keys(d.parent_feedback).length > 0) parentFeedback = d.parent_feedback;
				teacherRemarks = d.teacher_remarks || '';
				attendance = d.attendance_summary || { present: 0, total: 0, percentage: 0 };
			} else {
				studentName = 'New Entry';
			}
		}
		loading = false;
	});

	function getSubjectName(id: string) {
		return subjects.find(s => s.id === id)?.name || id;
	}

	async function save() {
		saving = true;
		statusMsg = '';
		const payload = {
			id: entryId || undefined,
			version: entryVersion,
			entry: {
				student_id: studentId,
				academic_year_id: year,
				term,
				scholastic,
				co_scholastic: coScholastic,
				health_pe: healthPE,
				work_education: workEducation,
				self_assessment: selfAssessment,
				peer_assessment: peerAssessment,
				parent_feedback: parentFeedback,
				teacher_remarks: teacherRemarks,
				attendance_summary: attendance
			}
		};
		const res = await api('PUT', '/hpc/entries', payload);
		saving = false;
		if (res.data) {
			entryId = (res.data as any).id;
			statusMsg = 'Saved successfully.';
		} else if (res.error) {
			statusMsg = res.error.code === 'VERSION_CONFLICT'
				? 'Conflict: entry was modified. Please refresh.'
				: 'Error: ' + res.error.message;
		}
	}

	async function publish() {
		if (!entryId) return;
		const res = await api('POST', '/hpc/entries/publish', { entry_id: entryId });
		if (res.data) statusMsg = 'Published successfully.';
		else if (res.error) statusMsg = 'Error: ' + res.error.message;
	}
</script>

<div class="max-w-4xl mx-auto space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-xl font-bold text-slate-900">HPC Entry</h1>
			<p class="text-sm text-slate-500">{studentName}</p>
		</div>
		<div class="flex gap-2">
			<button onclick={save} disabled={saving || loading}
				class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
				{saving ? 'Saving...' : 'Save'}
			</button>
			<button onclick={publish} disabled={!entryId}
				class="px-4 py-2 border border-slate-300 rounded-lg text-sm hover:bg-slate-50 disabled:opacity-50 transition-colors">
				Publish
			</button>
		</div>
	</div>

	{#if statusMsg}
		<div class="text-sm px-4 py-2 rounded-lg bg-slate-100 text-slate-700">{statusMsg}</div>
	{/if}

	{#if loading}
		<div class="p-8 text-center text-slate-400">Loading...</div>
	{:else}
		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
			<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Scholastic Achievement</h2>
			<table class="w-full text-sm">
				<thead>
					<tr class="bg-slate-50">
						<th class="text-left px-3 py-2 font-medium">Subject</th>
						<th class="text-center px-3 py-2 font-medium" style="width:100px;">Marks Scored</th>
						<th class="text-center px-3 py-2 font-medium" style="width:80px;">Max Marks</th>
						<th class="text-center px-3 py-2 font-medium" style="width:60px;">Grade</th>
					</tr>
				</thead>
				<tbody>
					{#each scholastic as s, i}
						<tr class="border-t border-slate-100">
							<td class="px-3 py-2">{getSubjectName(s.subject_id)}</td>
							<td class="px-3 py-2">
								<input type="number" bind:value={scholastic[i].marks_scored}
									class="w-20 px-2 py-1 border border-slate-300 rounded text-sm text-center" step="0.5" />
							</td>
							<td class="px-3 py-2 text-center">{s.max_marks}</td>
							<td class="px-3 py-2 text-center font-semibold">{s.grade}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
			<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Co-Scholastic Areas</h2>
			{#each Object.entries(coScholasticSections) as [sectionKey, section]}
				<div class="mb-3">
					<h3 class="text-xs font-medium text-slate-500 mb-2 uppercase tracking-wider">{section.label}</h3>
					<div class="grid grid-cols-2 gap-2">
						{#each section.items as item}
							<div class="flex items-center justify-between px-3 py-1.5 bg-slate-50 rounded">
								<span class="text-sm">{coScholastic[sectionKey][item.key].label}</span>
								<select bind:value={coScholastic[sectionKey][item.key].grade}
									class="text-xs px-2 py-0.5 border border-slate-300 rounded w-16 text-center">
									{#each gradeOptions as g}
										<option value={g}>{g || '—'}</option>
									{/each}
								</select>
							</div>
						{/each}
					</div>
				</div>
			{/each}
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
			<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Health & Physical Education</h2>
			<div class="grid grid-cols-2 gap-3">
				{#each Object.entries(healthPE) as [key, param] (key)}
					{@const p = param as { label: string; value: string }}
					<div>
						<label class="block text-xs text-slate-500 mb-0.5">{p.label}</label>
						<input type="text" bind:value={healthPE[key].value}
							class="w-full px-2 py-1 border border-slate-300 rounded text-sm" />
					</div>
				{/each}
			</div>
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
			<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Work Education</h2>
			<div>
				<label class="block text-xs text-slate-500 mb-0.5">Activity</label>
				<input type="text" bind:value={workEducation.activity}
					class="w-full px-2 py-1 border border-slate-300 rounded text-sm" />
			</div>
			<div>
				<label class="block text-xs text-slate-500 mb-0.5">Description</label>
				<textarea bind:value={workEducation.description}
					class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
			</div>
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
			<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Attendance</h2>
			<div class="flex gap-6 items-center">
				<div>
					<label class="block text-xs text-slate-500">Days Present</label>
					<input type="number" bind:value={attendance.present}
						class="w-24 px-2 py-1 border border-slate-300 rounded text-sm" />
				</div>
				<div>
					<label class="block text-xs text-slate-500">Total Days</label>
					<input type="number" bind:value={attendance.total}
						class="w-24 px-2 py-1 border border-slate-300 rounded text-sm" />
				</div>
				<div class="text-sm">
					<strong>Percentage:</strong> {attendance.total > 0 ? ((attendance.present / attendance.total) * 100).toFixed(1) : '0.0'}%
				</div>
			</div>
		</div>

		<div class="grid grid-cols-3 gap-4">
			<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
				<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Self Assessment</h2>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Strong Areas</label>
					<textarea bind:value={selfAssessment.strong_areas} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Areas for Improvement</label>
					<textarea bind:value={selfAssessment.areas_for_improvement} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Interests</label>
					<textarea bind:value={selfAssessment.interests} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
			</div>

			<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
				<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Peer Assessment</h2>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Teamwork</label>
					<textarea bind:value={peerAssessment.teamwork} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Behaviour</label>
					<textarea bind:value={peerAssessment.behavior} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Helpfulness</label>
					<textarea bind:value={peerAssessment.helpfulness} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
			</div>

			<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-3">
				<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Parent Feedback</h2>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Academic Progress</label>
					<textarea bind:value={parentFeedback.academic_progress} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Behaviour at Home</label>
					<textarea bind:value={parentFeedback.behavior_at_home} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
				<div>
					<label class="block text-xs text-slate-500 mb-0.5">Suggestions</label>
					<textarea bind:value={parentFeedback.suggestions} class="w-full px-2 py-1 border border-slate-300 rounded text-sm min-h-[60px]"></textarea>
				</div>
			</div>
		</div>

		<div class="bg-white rounded-xl border border-slate-200 p-4 space-y-4">
			<h2 class="text-sm font-semibold text-slate-700 border-b pb-2">Teacher's Holistic Remarks</h2>
			<textarea bind:value={teacherRemarks}
				class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm min-h-[100px]"
				placeholder="Enter holistic remarks about the student's overall progress..."></textarea>
		</div>
	{/if}
</div>
