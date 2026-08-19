<script lang="ts">
	import { api } from '$lib/api/client.svelte';
	import type { Subject } from '$lib/types';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Select from '$lib/components/Select.svelte';
	import MathText from '$lib/components/MathText.svelte';

	let subjects = $state<Subject[]>([]);
	let questionType = $state('mcq');
	let subjectId = $state('');
	let questionText = $state('');
	let answer = $state('');
	let marks = $state(1);
	let difficulty = $state('medium');
	let chapters = $state('');

	const typeOptions = [
		{ id: 'mcq', name: 'Multiple Choice' },
		{ id: 'true_false', name: 'True / False' },
		{ id: 'fill_blank', name: 'Fill in the Blank' },
		{ id: 'short_answer', name: 'Short Answer' },
	];
	const difficultyOptions = [
		{ id: 'easy', name: 'Easy' },
		{ id: 'medium', name: 'Medium' },
		{ id: 'hard', name: 'Hard' },
	];
	const tfOptions = [
		{ id: 'TRUE', name: 'TRUE' },
		{ id: 'FALSE', name: 'FALSE' },
	];

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
				<label for="subject" class="block text-sm font-medium text-slate-700 mb-1">Subject *</label>
				<Select id="subject" bind:value={subjectId} options={subjects} placeholder="Select" />
			</div>
			<div>
				<label for="type" class="block text-sm font-medium text-slate-700 mb-1">Type *</label>
				<Select id="type" bind:value={questionType} options={typeOptions} />
			</div>
			<div>
				<label for="marks" class="block text-sm font-medium text-slate-700 mb-1">Marks</label>
				<input id="marks" type="number" bind:value={marks} min="0.5" step="0.5" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
			</div>
			<div>
				<label for="difficulty" class="block text-sm font-medium text-slate-700 mb-1">Difficulty</label>
				<Select id="difficulty" bind:value={difficulty} options={difficultyOptions} />
			</div>
		</div>

		<div>
			<label for="question" class="block text-sm font-medium text-slate-700 mb-1">Question *</label>
			<textarea id="question" bind:value={questionText} rows="3" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm resize-none"></textarea>
		</div>

		{#if questionText.trim()}
			<div class="p-3 rounded-lg bg-slate-50 border border-slate-200">
				<p class="text-xs text-slate-500 mb-1">Preview</p>
				<div class="text-sm text-slate-900 leading-relaxed"><MathText text={questionText} /></div>
			</div>
		{/if}

		<div>
			<label for="chapters" class="block text-sm font-medium text-slate-700 mb-1">Chapters</label>
			<input id="chapters" type="text" bind:value={chapters} placeholder="separate with commas" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm">
		</div>

		{#if questionType === 'mcq'}
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium text-slate-700">Options</span>
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
			<label for="answer" class="block text-sm font-medium text-slate-700 mb-1">
				{questionType === 'true_false' ? 'Answer *' : questionType === 'fill_blank' ? 'Correct Answer(s) *' : 'Answer *'}
			</label>
			{#if questionType === 'true_false'}
				<Select id="answer" bind:value={answer} options={tfOptions} placeholder="Select" />
			{:else}
				<textarea id="answer" bind:value={answer} rows="2" class="w-full px-3 py-2 rounded-lg border border-slate-300 text-sm resize-none"></textarea>
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
