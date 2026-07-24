<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { QuizAssignment, Question, Subject } from '$lib/types';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	const quizId = $page.params.id;

	let subjects = $state<Subject[]>([]);
	let questions = $state<Question[]>([]);
	let saving = $state(false);

	let title = $state('');
	let description = $state('');
	let target_type = $state<'student' | 'staff'>('student');
	let target_id = $state('');
	let pass_pct = $state(40);
	let max_attempts = $state(1);
	let duration_min = $state<number | null>(null);
	let shuffle_questions = $state(false);
	let shuffle_options = $state(false);
	let show_result = $state(true);
	let start_at = $state('');
	let end_at = $state('');

	let selectedQuestions = $state<{ id: string; marks: number }[]>([]);
	let filterSubject = $state('');

	onMount(async () => {
		const [quizRes, subRes, qRes] = await Promise.all([
			api<QuizAssignment>('GET', `/quizzes/${quizId}`),
			api<Subject[]>('GET', '/subjects'),
			api<Question[]>('GET', '/questions'),
		]);
		if (subRes.data) subjects = subRes.data;
		if (qRes.data) questions = qRes.data;

		if (quizRes.data) {
			const q = quizRes.data;
			title = q.title;
			description = q.description;
			target_type = q.target_type;
			target_id = q.target_id || '';
			pass_pct = q.pass_pct;
			max_attempts = q.max_attempts;
			duration_min = q.duration_min ?? null;
			shuffle_questions = q.shuffle_questions;
			shuffle_options = q.shuffle_options;
			show_result = q.show_result;
			start_at = q.start_at ? q.start_at.slice(0, 16) : '';
			end_at = q.end_at ? q.end_at.slice(0, 16) : '';

			const sqRes = await api<{ question_id: string; marks: number }[]>('GET', `/quizzes/${quizId}/questions`);
			if (sqRes.data) {
				selectedQuestions = sqRes.data.map(sq => ({ id: sq.question_id, marks: sq.marks }));
			}
		}
	});

	let filteredQuestions = $derived(
		filterSubject
			? questions.filter(q => q.subject_id === filterSubject)
			: questions
	);

	let selectedIds = $derived(new Set(selectedQuestions.map(sq => sq.id)));

	function toggleQuestion(qid: string) {
		if (selectedIds.has(qid)) {
			selectedQuestions = selectedQuestions.filter(sq => sq.id !== qid);
		} else {
			selectedQuestions = [...selectedQuestions, { id: qid, marks: 1 }];
		}
	}

	function updateMarks(qid: string, marks: number) {
		selectedQuestions = selectedQuestions.map(sq =>
			sq.id === qid ? { ...sq, marks } : sq
		);
	}

	async function save() {
		saving = true;
		await api('PUT', `/quizzes/${quizId}`, {
			title, description, target_type, target_id: target_id || undefined,
			pass_pct, max_attempts, duration_min: duration_min || undefined,
			shuffle_questions, shuffle_options, show_result,
			start_at: start_at || undefined, end_at: end_at || undefined,
		});
		if (selectedQuestions.length > 0) {
			await api('POST', `/quizzes/${quizId}/questions`, {
				questions: selectedQuestions.map(sq => ({ question_id: sq.id, marks: sq.marks }))
			});
		}
		saving = false;
		goto('/quizzes');
	}
</script>

<div class="max-w-3xl mx-auto space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-slate-900">Edit Quiz</h1>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-6 space-y-5">
		<h2 class="text-lg font-semibold text-slate-900">Quiz Details</h2>

		<div>
			<label class="block text-sm font-medium text-slate-700 mb-1">Title</label>
			<input type="text" bind:value={title} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm" required>
		</div>

		<div>
			<label class="block text-sm font-medium text-slate-700 mb-1">Description</label>
			<textarea bind:value={description} rows={2} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm"></textarea>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Target</label>
				<select bind:value={target_type} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="student">Students</option>
					<option value="staff">Staff</option>
				</select>
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Target ID (optional)</label>
				<input type="text" bind:value={target_id} placeholder="Class ID or empty" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
		</div>

		<div class="grid grid-cols-3 gap-4">
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Pass %</label>
				<input type="number" bind:value={pass_pct} min="0" max="100" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Max Attempts</label>
				<input type="number" bind:value={max_attempts} min="1" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Duration (min)</label>
				<input type="number" bind:value={duration_min} min="0" placeholder="No limit" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
		</div>

		<div class="flex flex-wrap gap-6">
			<label class="flex items-center gap-2 text-sm text-slate-700">
				<input type="checkbox" bind:checked={shuffle_questions}>
				Shuffle questions
			</label>
			<label class="flex items-center gap-2 text-sm text-slate-700">
				<input type="checkbox" bind:checked={shuffle_options}>
				Shuffle options
			</label>
			<label class="flex items-center gap-2 text-sm text-slate-700">
				<input type="checkbox" bind:checked={show_result}>
				Show result after submission
			</label>
		</div>

		<div class="grid grid-cols-2 gap-4">
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Start At</label>
				<input type="datetime-local" bind:value={start_at} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">End At</label>
				<input type="datetime-local" bind:value={end_at} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
		</div>
	</div>

	<div class="bg-white rounded-xl border border-slate-200 p-6 space-y-4">
		<h2 class="text-lg font-semibold text-slate-900">Questions</h2>

		<div class="flex items-center gap-3">
			<select bind:value={filterSubject} class="px-3 py-1.5 rounded-lg border border-slate-300 text-sm">
				<option value="">All Subjects</option>
				{#each subjects as s}
					<option value={s.id}>{s.name}</option>
				{/each}
			</select>
			<span class="text-xs text-slate-400">{selectedQuestions.length} selected</span>
		</div>

		<div class="divide-y divide-slate-100 max-h-96 overflow-y-auto border border-slate-200 rounded-lg">
			{#each filteredQuestions as q (q.id)}
				<div class="flex items-center gap-3 p-3 hover:bg-slate-50">
					<input type="checkbox" checked={selectedIds.has(q.id)} onchange={() => toggleQuestion(q.id)} class="shrink-0">
					<div class="flex-1 min-w-0">
						<p class="text-sm text-slate-900 line-clamp-1">{q.question_text}</p>
						<div class="text-xs text-slate-400 mt-0.5">
							{q.question_type} &middot; {q.difficulty}
						</div>
					</div>
					{#if selectedIds.has(q.id)}
						<input type="number" value={selectedQuestions.find(sq => sq.id === q.id)?.marks || 1}
							oninput={(e) => updateMarks(q.id, Number((e.target as HTMLInputElement).value))}
							min="0" step="0.5" class="w-16 px-2 py-1 rounded border border-slate-300 text-xs text-center" title="Marks">
					{/if}
				</div>
			{/each}
		</div>
	</div>

	<div class="flex justify-end gap-3">
		<a href="/quizzes" class="px-4 py-2 border border-slate-300 text-slate-700 rounded-lg text-sm font-medium hover:bg-slate-50">Cancel</a>
		<button onclick={save} disabled={!title || saving}
			class="px-4 py-2 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
			{saving ? 'Saving...' : 'Save Changes'}
		</button>
	</div>
</div>
