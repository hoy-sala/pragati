import fs from 'fs';
const API = 'http://192.168.1.2:9090/api/v1';

async function api(method, path, body, token) {
	const res = await fetch(`${API}${path}`, {
		method,
		headers: {
			'Content-Type': 'application/json',
			...(token ? { Authorization: `Bearer ${token}` } : {}),
		},
		body: body ? JSON.stringify(body) : undefined,
	});
	return res.json();
}

// 1. Login
const loginRes = await api('POST', '/auth/login', {
	email: 'admin@pragati.edu',
	password: 'pragati123',
});
if (loginRes.error) { console.error('Login failed:', loginRes.error); process.exit(1); }
const token = loginRes.data.access_token;
console.log('Logged in');

// 2. Get classes and academic year
const [classRes, yearRes] = await Promise.all([
	api('GET', '/classes?limit=50', null, token),
	api('GET', '/academic-years?limit=50', null, token),
]);
if (!classRes.data) { console.error('Failed to fetch classes'); process.exit(1); }
const classes = classRes.data;
const academicYear = yearRes.data?.find(y => y.is_current) || yearRes.data?.[0];
console.log(`Found ${classes.length} classes, academic year: ${academicYear?.name}`);

// Build class name -> id map
const classMap = {};
for (const c of classes) {
	classMap[c.name] = c.id;
}

// 3. Read CSV
const csvText = fs.readFileSync('C:\\Users\\MDRS Bahaddurghatta\\Downloads\\student_list.csv', 'utf-8');
const lines = csvText.trim().split('\n');
const headers = lines[0].split(',').map(h => h.replace(/"/g, '').trim());
console.log(`Headers: ${headers.join(', ')}`);

// Find column indices
const idxName = headers.findIndex(h => h === 'Student Name');
const idxFather = headers.findIndex(h => h === 'Father Name');
const idxSATS = headers.findIndex(h => h === 'Student Id');
const idxDOB = headers.findIndex(h => h === 'DOB');
const idxSex = headers.findIndex(h => h === 'Sex');
const idxClass = headers.findIndex(h => h === 'Class Studying');

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

function normalizeClass(val) {
	const m = val.match(/(\d+)/);
	return m ? `Class ${m[1]}` : val;
}

function normalizeGender(val) {
	if (!val) return '';
	const upper = val.toUpperCase();
	if (upper.includes('GIRL') || upper.includes('FEMALE') || val.startsWith('2')) return 'female';
	return 'male';
}

let imported = 0;
let errors = [];

for (let i = 1; i < lines.length; i++) {
	const cols = parseCSVLine(lines[i]);
	if (cols.length < 2 || !cols[idxName]) continue;

	const studentId = cols[idxSATS]?.replace(/\s+/g, '') || '';
	const name = cols[idxName] || '';
	const fatherName = cols[idxFather] || '';
	const className = normalizeClass(cols[idxClass] || '');
	const gender = normalizeGender(cols[idxSex] || '');
	const dob = cols[idxDOB] || '';
	const classId = classMap[className];
	const academicYearId = academicYear?.id;

	if (!studentId || !classId) {
		errors.push({ row: i, name, reason: !studentId ? 'missing SATS' : `unknown class: ${className}` });
		continue;
	}

	// Split name into first and last (last word = last name)
	const parts = name.trim().split(/\s+/);
	const firstName = parts[0];
	const lastName = parts.slice(1).join(' ') || undefined;

	const body = {
		sats_number: studentId,
		first_name: firstName,
		last_name: lastName,
		gender: gender || undefined,
		date_of_birth: dob || undefined,
		parent_name: fatherName || undefined,
		class_id: classId,
		academic_year_id: academicYearId,
	};

	const res = await api('POST', '/students', body, token);
	if (res.error) {
		errors.push({ row: i, name, reason: res.error.message });
	} else {
		imported++;
	}

	if (i % 50 === 0) console.log(`Progress: ${i}/${lines.length - 1} (${imported} imported, ${errors.length} errors)`);
}

console.log(`\nDone. Imported: ${imported}, Errors: ${errors.length}`);
if (errors.length > 0) {
	console.log('Errors:');
	for (const e of errors.slice(0, 20)) {
		console.log(`  Row ${e.row}: ${e.name} — ${e.reason}`);
	}
}