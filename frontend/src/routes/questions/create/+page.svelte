<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Subject } from '$lib/types';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let subjects = $state<Subject[]>([]);
	let questionType = $state('mcq');
	let subjectId = $state('');
	let questionText = $state('');
	let answer = $state('');
	let marks = $state(1);
	let difficulty = $state('medium');
	let chapters = $state('');

	// MCQ fields
	let options = $state([{ key: 'A', value: '', correct: false }, { key: 'B', value: '', correct: false }]);

	let saving = $state(false);
	let error = $state('');

	onMount(async () => {
		const res = await api<Subject[]>('GET', '/subjects');
		if (res.data) subjects = res.data;
	});

	function addOption() {
		const key = String.fromCharCode(65 + options.length);
		options = [...options, { key, value: '', correct: false }];
	}

	function removeOption(idx: number) {
		if (options.length <= 2) return;
		options = options.filter((_, i) => i !== idx).map((o, i) => ({ ...o, key: String.fromCharCode(65 + i) }));
	}

	function setCorrect(key: string) {
		options = options.map(o => ({ ...o, correct: o.key === key }));
	}

	async function handleSubmit() {
		error = '';
		if (!subjectId || !questionText || !answer) {
			error = 'Subject, question text, and answer are required.';
			return;
		}

		saving = true;
		const payload: any = {
			subject_id: subjectId,
			question_type: questionType,
			question_text: questionText,
			answer,
			marks,
			difficulty,
			chapters: chapters ? chapters.split(',').map(s => s.trim()) : [],
		};

		if (questionType === 'mcq') {
			payload.options = options;
		}

		const res = await api('POST', '/questions', payload);
		saving = false;

		if (res.error) {
			error = res.error.message;
		} else {
			goto('/questions');
		}
	}
</script>

<div class="max-w-2xl mx-auto space-y-6">
	<div class="flex items-center gap-4">
		<a href="/questions" class="text-slate-400 hover:text-slate-600">&larr; Back</a>
		<h1 class="text-xl font-bold text-slate-900">Create Question</h1>
	</div>

	<form class="bg-white rounded-xl border border-slate-200 p-6 space-y-4" onsubmit={handleSubmit}>
		<div class="grid grid-cols-2 gap-4">
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Subject *</label>
				<select bind:value={subjectId} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					{#each subjects as s}
						<option value={s.id}>{s.name}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Type *</label>
				<select bind:value={questionType} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="mcq">Multiple Choice</option>
					<option value="true_false">True / False</option>
					<option value="fill_blank">Fill in the Blank</option>
					<option value="short_answer">Short Answer</option>
				</select>
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Marks</label>
				<input type="number" bind:value={marks} min="0.5" step="0.5" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
			<div>
				<label class="block text-sm font-medium text-slate-700 mb-1">Difficulty</label>
				<select bind:value={difficulty} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="easy">Easy</option>
					<option value="medium">Medium</option>
					<option value="hard">Hard</option>
				</select>
			</div>
		</div>

		<div>
			<label class="block text-sm font-medium text-slate-700 mb-1">Question *</label>
			<textarea bind:value={questionText} rows="3" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm resize-none"></textarea>
		</div>

		<div>
			<label class="block text-sm font-medium text-slate-700 mb-1">Chapters</label>
			<input type="text" bind:value={chapters} placeholder="separate with commas" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
		</div>

		{#if questionType === 'mcq'}
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<label class="text-sm font-medium text-slate-700">Options</label>
					<button type="button" onclick={addOption} class="text-xs text-primary-600 hover:text-primary-700">+ Add option</button>
				</div>
				{#each options as opt, i}
					<div class="flex items-center gap-3">
						<button type="button" onclick={() => setCorrect(opt.key)}
							class="w-5 h-5 rounded-full border-2 flex items-center justify-center text-xs {opt.correct ? 'bg-primary-600 border-primary-600 text-white' : 'border-slate-300'}"
							title="Mark as correct">
							{opt.correct ? '✓' : ''}
						</button>
						<span class="text-xs font-mono text-slate-500 w-4">{opt.key}</span>
						<input type="text" bind:value={opt.value} placeholder="Option {opt.key}" class="flex-1 px-3 py-2 rounded-lg border border-slate-300 text-sm">
						<button type="button" onclick={() => removeOption(i)} class="text-xs text-danger-600 hover:text-danger-700">Remove</button>
					</div>
				{/each}
				{#if questionType === 'mcq'}
					<input type="hidden" bind:value={answer} />
				{/if}
				<p class="text-xs text-slate-400">Click the circle to mark the correct answer</p>
			</div>
		{/if}

		<div>
			<label class="block text-sm font-medium text-slate-700 mb-1">
				{questionType === 'true_false' ? 'Answer *' : questionType === 'fill_blank' ? 'Correct Answer(s) *' : 'Answer *'}
			</label>
			{#if questionType === 'true_false'}
				<select bind:value={answer} class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
					<option value="">Select</option>
					<option value="TRUE">TRUE</option>
					<option value="FALSE">FALSE</option>
				</select>
			{:else}
				<textarea bind:value={answer} rows="2" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm resize-none"></textarea>
			{/if}
		</div>

		{#if error}
			<div class="text-sm text-danger-600 bg-danger-50 rounded-lg p-3">{error}</div>
		{/if}

		<button type="submit" disabled={saving} class="w-full py-2 px-4 bg-primary-600 text-white rounded-lg text-sm font-medium hover:bg-primary-700 disabled:opacity-50 transition-colors">
			{saving ? 'Creating...' : 'Create Question'}
		</button>
	</form>
</div>
