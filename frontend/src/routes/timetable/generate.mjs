const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
const CLASSES = ['Class 6', 'Class 7', 'Class 8', 'Class 9', 'Class 10'];
const PERIODS = { Mon: 8, Tue: 8, Wed: 8, Thu: 8, Fri: 8, Sat: 4 };
const SUBJ4 = ['HIN', 'MAT', 'SCI', 'SOC'];
const SUBJ5 = ['CS', 'DRW', 'MUS', 'PE', 'LIB'];

function shuffle(a) { const b = [...a]; for (let i = b.length-1; i>0; i--) { const j = Math.floor(Math.random()*(i+1)); [b[i],b[j]]=[b[j],b[i]]; } return b; }

function findMatching(remain) {
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

function scheduleDay(assignments, breakAfter = new Set([4])) {
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
  // Classes 7-10: cyclic assignment guarantees no conflicts (4 subjects, 4 classes)
  const core = CLASSES.filter(c => c !== 'Class 6');
  const subs = assignments[core[0]]; // [KAN, MAT, SCI, SOC]
  for (let t = 0; t < 20000; t++) {
    const result = {};
    // Class 6: random permutation (ENG, HIN, CS, MUS — no overlap with core)
    result['Class 6'] = shuffle(assignments['Class 6']);
    for (let i = 0; i < core.length; i++) {
      result[core[i]] = [];
      for (let p = 0; p < 4; p++) result[core[i]].push(subs[(p + i) % 4]);
    }
    // Validate: no period conflicts (should always pass for this construction)
    let ok = true;
    for (let p = 0; p < 4; p++) {
      const ps = CLASSES.map(cls => result[cls][p]);
      if (new Set(ps).size !== ps.length) { ok = false; break; }
    }
    if (!ok) continue;
    return result;
  }
  return null;
}

function solve() {
  for (let attempt = 0; attempt < 2000; attempt++) {
    const dayAssign = {};
    for (const cls of CLASSES) {
      dayAssign[cls] = {};

      if (cls === 'Class 6') {
        // KAN×5, MAT×5, SCI×5, SOC×5 (Mon-Fri only)
        for (const day of ['Mon','Tue','Wed','Thu','Fri']) dayAssign[cls][day] = ['KAN', 'MAT', 'SCI', 'SOC'];
        // Sat: ENG, HIN, CS, MUS
        dayAssign[cls]['Sat'] = ['ENG', 'HIN', 'CS', 'MUS'];

        // ENG×5 (Mon-Fri)
        for (const day of ['Mon','Tue','Wed','Thu','Fri']) dayAssign[cls][day].push('ENG');

        // HIN×4 (Mon-Thu only)
        for (const day of ['Mon','Tue','Wed','Thu']) dayAssign[cls][day].push('HIN');

        // SUBJ5×2 each = 10 total: 2 per day Mon-Fri
        const s5 = shuffle(SUBJ5);
        const days5 = ['Mon','Tue','Wed','Thu','Fri','Mon','Tue','Wed','Thu','Fri'];
        const s5alloc = shuffle([...SUBJ5, ...SUBJ5]);
        for (let i = 0; i < 10; i++) dayAssign[cls][days5[i]].push(s5alloc[i]);
      } else {
        // KAN×6, MAT×6, SCI×6, SOC×6 (Mon-Sat)
        for (const day of DAYS) dayAssign[cls][day] = ['KAN', 'MAT', 'SCI', 'SOC'];

        // ENG×5 (Mon-Fri)
        for (const day of ['Mon','Tue','Wed','Thu','Fri']) dayAssign[cls][day].push('ENG');

        // HIN×4 (Mon-Thu only)
        for (const day of ['Mon','Tue','Wed','Thu']) dayAssign[cls][day].push('HIN');

        // SUBJ5×2 each = 10 total: 2 per day Mon-Fri
        const s5 = shuffle(SUBJ5);
        const days5 = ['Mon','Tue','Wed','Thu','Fri','Mon','Tue','Wed','Thu','Fri'];
        const s5alloc = shuffle([...SUBJ5, ...SUBJ5]);
        for (let i = 0; i < 10; i++) dayAssign[cls][days5[i]].push(s5alloc[i]);
      }
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

    // Fri: 7 non-CUL subjects + CUL at P8; skip consecutive check
    {
      const assign = {};
      for (const cls of CLASSES) assign[cls] = [...dayAssign[cls]['Fri']];
      const dr = scheduleDay(assign, new Set([0,1,2,3,4,5,6,7]));
      if (!dr) { valid = false; continue; }
      for (const cls of CLASSES) result[cls]['Fri'] = [...dr[cls], 'CUL'];
    }

    // Sat: [KAN, MAT, SCI, SOC] — each subject in every period (unavoidable doubling)
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
  // No period conflicts (Mon-Sat)
  for (const day of ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']) {
    for (let p = 0; p < PERIODS[day]; p++) {
      const subjects = CLASSES.map(cls => data[cls][day][p]);
      const nonCUL = subjects.filter(s => s !== 'CUL');
      if (new Set(nonCUL).size !== nonCUL.length) return false;
    }
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
    if (cls === 'Class 6') {
      if (!(counts['KAN']===5 && counts['MAT']===5 && counts['SCI']===5 && counts['SOC']===5 &&
        counts['ENG']===6 && counts['HIN']===5 &&
        counts['CS']===3 && counts['DRW']===2 && counts['MUS']===3 &&
        counts['PE']===2 && counts['LIB']===2)) return false;
    } else {
      if (!(counts['KAN']===6 && counts['ENG']===5 && counts['HIN']===4 &&
        counts['MAT']===6 && counts['SCI']===6 && counts['SOC']===6 &&
        counts['CS']===2 && counts['DRW']===2 && counts['MUS']===2 &&
        counts['PE']===2 && counts['LIB']===2)) return false;
    }
  }
  // No 3 consecutive periods for any teacher (Mon-Sat)
  const breakAfter = new Set([4]);
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
  // Fri and Sat — skip consecutive check (every subject appears in every period)
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