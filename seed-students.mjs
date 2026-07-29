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

const login = await api('POST', '/auth/login', { email: 'admin@pragati.edu', password: 'pragati123' });
if (login.error) { console.error('Login failed:', login.error); process.exit(1); }
const token = login.data.access_token;
console.log('Logged in');

const classRes = await api('GET', '/classes?limit=50', null, token);
if (!classRes.data) { console.error('No classes returned'); process.exit(1); }
const classMap = {};
for (const c of classRes.data) classMap[c.name] = c.id;
console.log('Classes:', JSON.stringify(classMap));

const yearId = 'b2586083-8c59-40d6-af5b-5c91cf6903a7';

const csvText = fs.readFileSync('C:\\Users\\MDRS Bahaddurghatta\\Downloads\\student_list.csv', 'utf-8');
const lines = csvText.trim().split('\n');
const rawHeaders = lines[0].split(',').map(h => h.replace(/"/g, '').trim());
const idxName = rawHeaders.indexOf('Student Name');
const idxSATS = rawHeaders.indexOf('Student Id');
const idxDOB = rawHeaders.indexOf('DOB');
const idxSex = rawHeaders.indexOf('Sex');
const idxClass = rawHeaders.indexOf('Class Studying');

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

let imported = 0;
let errors = [];

for (let i = 1; i < lines.length; i++) {
	const cols = parseCSVLine(lines[i]);
	if (!cols[idxName]) continue;

	const studentId = (cols[idxSATS] || '').replace(/\s+/g, '');
	const m = (cols[idxClass] || '').match(/(\d+)/);
	const className = m ? 'Class ' + m[1] : '';
	const classId = classMap[className];

	if (!studentId || !classId) {
		errors.push({ row: i, name: cols[idxName], reason: !studentId ? 'missing SATS' : 'unknown class: ' + className });
		continue;
	}

	const parts = cols[idxName].trim().split(/\s+/);
	const body = {
		sats_number: studentId,
		first_name: parts[0],
		last_name: parts.slice(1).join(' ') || undefined,
		gender: (cols[idxSex] || '').startsWith('2') ? 'female' : 'male',
		date_of_birth: cols[idxDOB] || undefined,
		class_id: classId,
		academic_year_id: yearId,
	};

	const res = await api('POST', '/students', body, token);
	if (res.error) {
		errors.push({ row: i, name: cols[idxName], reason: res.error.message });
	} else {
		imported++;
	}
	if (i % 50 === 0) console.log('Progress:', i, '/', lines.length - 1);
}

console.log('Done. Imported:', imported, 'Errors:', errors.length);
if (errors.length) {
	console.log('First errors:');
	for (const e of errors.slice(0, 5)) console.log(`  Row ${e.row}: ${e.name} — ${e.reason}`);
}
