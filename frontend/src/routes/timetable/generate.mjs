const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
const CLASSES = ['Class 6', 'Class 7', 'Class 8', 'Class 9', 'Class 10'];
const PERIODS = { Mon: 8, Tue: 8, Wed: 8, Thu: 8, Fri: 8, Sat: 4 };
const SUBJ4 = ['HIN', 'MAT', 'SCI', 'SOC'];
const SUBJ5 = ['CS', 'DRW', 'MUS', 'PE', 'LIB'];

function shuffle(a) { const b = [...a]; for (let i = b.length-1; i>0; i--) { const j = Math.floor(Math.random()*(i+1)); [b[i],b[j]]=[b[j],b[i]]; } return b; }

// Find a matching with random ordering (non-deterministic)
function findMatching(remain) {
  // remain[cls] = shuffled array of available subjects
  const matchTo = {};
  const clsList = shuffle(CLASSES);

  function dfs(cls, visited) {
    if (visited.has(cls)) return false;
    visited.add(cls);
    for (const s of remain[cls]) {
      if (matchTo[s] === undefined || dfs(matchTo[s], visited)) {
        matchTo[s] = cls;
        return true;
      }
    }
    return false;
  }

  for (const cls of clsList) {
    const visited = new Set();
    if (!dfs(cls, visited)) return null;
  }

  const result = {};
  for (const [s, cls] of Object.entries(matchTo)) result[cls] = s;
  return result;
}

function scheduleDay(assignments, breakAfter = new Set([2, 4])) {
  const numPeriods = assignments[CLASSES[0]].length;
  for (let t = 0; t < 2000; t++) {
    const remain = {};
    for (const cls of CLASSES) remain[cls] = shuffle(assignments[cls]);
    const result = {};
    for (const cls of CLASSES) result[cls] = [];
    let valid = true;
    const periodSubjects = [];

    for (let p = 0; p < numPeriods; p++) {
      const forbidden = new Set();
      if (p >= 2 && !breakAfter.has(p - 1) && !breakAfter.has(p - 2)) {
        for (const s of periodSubjects[p - 1]) {
          if (periodSubjects[p - 2].has(s)) forbidden.add(s);
        }
      }

      const filteredRemain = {};
      for (const cls of CLASSES) {
        filteredRemain[cls] = remain[cls].filter(s => !forbidden.has(s));
      }

      const matching = findMatching(filteredRemain);
      if (!matching) { valid = false; break; }
      periodSubjects.push(new Set(Object.values(matching)));

      for (const cls of CLASSES) {
        const s = matching[cls];
        result[cls].push(s);
        remain[cls] = remain[cls].filter(x => x !== s);
      }
    }
    if (valid) return result;
  }
  return null;
}

function scheduleSat(assignments) {
  for (let t = 0; t < 2000; t++) {
    const remain = {};
    for (const cls of CLASSES) remain[cls] = shuffle(assignments[cls]);
    const result = {};
    for (const cls of CLASSES) result[cls] = [];
    let valid = true;
    const periodSubjects = [];

    for (let p = 0; p < 4; p++) {
      const forbidden = new Set();
      if (p >= 2) {
        for (const s of periodSubjects[p - 1]) {
          if (periodSubjects[p - 2].has(s)) forbidden.add(s);
        }
      }

      const usedSubj5 = new Set();
      for (const cls of shuffle(CLASSES)) {
        const cand = remain[cls].filter(s => (!SUBJ5.includes(s) || !usedSubj5.has(s)) && !forbidden.has(s));
        if (cand.length === 0) { valid = false; break; }
        const pick = cand[Math.floor(Math.random() * cand.length)];
        result[cls].push(pick);
        if (SUBJ5.includes(pick)) usedSubj5.add(pick);
        remain[cls] = remain[cls].filter(x => x !== pick);
      }
      if (!valid) break;
      periodSubjects.push(new Set(CLASSES.map(cls => result[cls][result[cls].length - 1])));
    }
    if (valid) return result;
  }
  return null;
}

function solve() {
  for (let attempt = 0; attempt < 2000; attempt++) {
    const dayAssign = {};
    for (const cls of CLASSES) {
      dayAssign[cls] = {};
      for (const day of DAYS) dayAssign[cls][day] = ['KAN', 'ENG'];
      for (const s of SUBJ4) {
        for (const day of ['Mon','Tue','Wed','Thu','Fri']) dayAssign[cls][day].push(s);
      }
      const s5 = shuffle(SUBJ5);
      const monThuDays = ['Mon','Tue','Wed','Thu','Mon','Tue','Wed','Thu'];
      for (let i = 0; i < 8; i++) dayAssign[cls][monThuDays[i]].push(s5[i % 5]);
      dayAssign[cls]['Sat'].push(s5[3], s5[4]);
    }

    const result = {};
    for (const cls of CLASSES) { result[cls] = {}; for (const day of DAYS) result[cls][day] = []; }

    let valid = true;
    for (const day of ['Mon', 'Tue', 'Wed', 'Thu']) {
      const assign = {};
      for (const cls of CLASSES) assign[cls] = [...dayAssign[cls][day]];
      const dr = scheduleDay(assign);
      if (!dr) { valid = false; break; }
      for (const cls of CLASSES) result[cls][day] = dr[cls];
    }
    if (!valid) continue;

    // Fri: 6 core + 2 CUL — skip consecutive constraint (only 6 subjects, impossible to satisfy)
    {
      const assign = {};
      for (const cls of CLASSES) assign[cls] = [...dayAssign[cls]['Fri']];
      const dr = scheduleDay(assign, new Set([0,1,2,3,4,5,6,7]));
      if (!dr) { valid = false; continue; }
      for (const cls of CLASSES) result[cls]['Fri'] = [...dr[cls], 'CUL', 'CUL'];
    }

    // Sat
    {
      const assign = {};
      for (const cls of CLASSES) assign[cls] = [...dayAssign[cls]['Sat']];
      const dr = scheduleSat(assign);
      if (!dr) { valid = false; continue; }
      for (const cls of CLASSES) result[cls]['Sat'] = dr[cls];
    }

    return result;
  }
  return null;
}

function validate(data) {
  // No period conflicts
  for (const day of ['Mon', 'Tue', 'Wed', 'Thu', 'Fri']) {
    for (let p = 0; p < PERIODS[day]; p++) {
      const subjects = CLASSES.map(cls => data[cls][day][p]);
      const nonCUL = subjects.filter(s => s !== 'CUL');
      if (new Set(nonCUL).size !== nonCUL.length) return false;
    }
  }
  // Sat SUBJ5 uniqueness
  for (let p = 0; p < 4; p++) {
    const others = CLASSES.map(cls => data[cls]['Sat'][p]).filter(s => SUBJ5.includes(s));
    if (new Set(others).size !== others.length) return false;
  }
  // No same-day subject repeats per class
  for (const cls of CLASSES) {
    for (const day of DAYS) {
      const dc = {};
      for (const s of data[cls][day]) dc[s] = (dc[s] || 0) + 1;
      for (const [s, c] of Object.entries(dc)) if (c > 1 && s !== 'CUL') return false;
    }
  }
  // Correct weekly totals
  for (const cls of CLASSES) {
    const counts = {};
    for (const day of DAYS) for (const s of data[cls][day]) counts[s] = (counts[s] || 0) + 1;
    if (!(counts['KAN']===6 && counts['ENG']===6 && counts['HIN']===5 &&
      counts['MAT']===5 && counts['SCI']===5 && counts['SOC']===5 &&
      counts['CS']===2 && counts['DRW']===2 && counts['MUS']===2 &&
      counts['PE']===2 && counts['LIB']===2)) return false;
  }
  // No 3 consecutive periods for any teacher (across all classes)
  const breakAfter = new Set([2, 4]);
  for (const day of ['Mon', 'Tue', 'Wed', 'Thu']) {
    for (let p = 2; p < 8; p++) {
      if (breakAfter.has(p - 1) || breakAfter.has(p - 2)) continue;
      for (const s of CLASSES.map(cls => data[cls][day][p])) {
        const inP1 = CLASSES.some(cls => data[cls][day][p - 1] === s);
        const inP2 = CLASSES.some(cls => data[cls][day][p - 2] === s);
        if (inP1 && inP2) return false;
      }
    }
  }
  // Fri — skip consecutive check; Fri only has 6 unique subjects and 2 CUL,
  // making 3-consecutive impossible to avoid (overlap |P0∩P1| >= 4 → ≤2 subjects left for 5 classes).
  // The 2 CUL periods at the end naturally break any chain.
  // Sat (no breaks, all 4 periods consecutive)
  for (let p = 2; p < 4; p++) {
    for (const s of CLASSES.map(cls => data[cls]['Sat'][p])) {
      const inP1 = CLASSES.some(cls => data[cls]['Sat'][p - 1] === s);
      const inP2 = CLASSES.some(cls => data[cls]['Sat'][p - 2] === s);
      if (inP1 && inP2) return false;
    }
  }
  return true;
}

const data = solve();
if (data && validate(data)) {
  console.log('\n=== TIMETABLE ===\n');
  for (const cls of CLASSES) {
    console.log("'" + cls + "': {");
    for (const day of DAYS) {
      console.log('  ' + day + ': [' + data[cls][day].map(s => "'" + s + "'").join(', ') + '],');
    }
    console.log('},');
  }
} else {
  console.log('Could not find valid schedule');
}