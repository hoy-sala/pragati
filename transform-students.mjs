import fs from 'fs';

const csvText = fs.readFileSync('C:\\Users\\MDRS Bahaddurghatta\\Downloads\\student_list.csv', 'utf-8');
const lines = csvText.trim().split('\n');

function parseCSVLine(line) {
	const result = [];
	let current = '';
	let inQuotes = false;
	for (let i = 0; i < line.length; i++) {
		const ch = line[i];
		if (ch === '"') { inQuotes = !inQuotes; continue; }
		if (ch === ',' && !inQuotes) { result.push(current.trim()); current = ''; continue; }
		current += ch;
	}
	result.push(current.trim());
	return result;
}

const headers = parseCSVLine(lines[0]);
const idxName = headers.findIndex(h => h === 'Student Name');
const idxFather = headers.findIndex(h => h === 'Father Name');
const idxSATS = headers.findIndex(h => h === 'Student Id');
const idxDOB = headers.findIndex(h => h === 'DOB');
const idxSex = headers.findIndex(h => h === 'Sex');
const idxClass = headers.findIndex(h => h === 'Class Studying');

function normalizeClass(val) {
	const m = val.match(/(\d+)/);
	return m ? `Class ${m[1]}` : val;
}

function normalizeGender(val) {
	if (!val) return '';
	if (val.startsWith('2')) return 'female';
	return 'male';
}

const outHeaders = ['sats_number', 'first_name', 'last_name', 'date_of_birth', 'gender', 'class', 'academic_year', 'parent_name'];
const outLines = [outHeaders.join(',')];

for (let i = 1; i < lines.length; i++) {
	const cols = parseCSVLine(lines[i]);
	if (cols.length < 2 || !cols[idxName]) continue;

	const studentId = cols[idxSATS]?.replace(/\s+/g, '') || '';
	const name = cols[idxName] || '';
	const fatherName = cols[idxFather] || '';
	const className = normalizeClass(cols[idxClass] || '');
	const gender = normalizeGender(cols[idxSex] || '');
	const dob = cols[idxDOB] || '';

	const parts = name.trim().split(/\s+/);
	const firstName = parts[0];
	const lastName = parts.slice(1).join(' ');

	const row = [
		studentId,
		firstName,
		lastName || '',
		dob,
		gender,
		className,
		'2026-27',
		fatherName.replace(/,/g, ';'),
	];
	outLines.push(row.join(','));
}

const outputPath = 'C:\\Users\\MDRS Bahaddurghatta\\Downloads\\students_for_import.csv';
fs.writeFileSync(outputPath, outLines.join('\n'), 'utf-8');
console.log(`Written ${outLines.length - 1} rows to ${outputPath}`);