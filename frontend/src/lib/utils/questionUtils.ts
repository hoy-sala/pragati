const typeColors: Record<string, string> = {
	mcq: 'bg-blue-100 text-blue-700',
	true_false: 'bg-purple-100 text-purple-700',
	fill_blank: 'bg-green-100 text-green-700',
	short_answer: 'bg-amber-100 text-amber-700'
};

const typeLabels: Record<string, string> = {
	mcq: 'MCQ',
	true_false: 'T/F',
	fill_blank: 'Fill',
	short_answer: 'Short'
};

export { typeColors, typeLabels };
