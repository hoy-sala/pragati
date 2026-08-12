export interface TimetableCell {
  code: string;
  name: string;
}

export const TEACHER_NAMES: Record<string, string> = {
  KAN: 'Siddappa N',
  ENG: 'Manoj Kumar M J',
  HIN: 'Manoranjan E',
  MAT: 'Vedamurthy D',
  SCI: 'Shivakumar K',
  SOC: 'Gopi K',
  CS: 'Hoysala T',
  DRW: 'Shivayogi Kasavannavar',
  MUS: 'Guruswamy E',
  PE: 'Kallesh K G',
  LIB: 'Hoysala T',
  CUL: '',
};

export interface DaySchedule {
  label: string;
  periods: TimetableCell[];
}

export interface ClassSchedule {
  name: string;
  days: DaySchedule[];
}

export const SUBJECT_INFO: Record<string, { name: string; color: string }> = {
  KAN: { name: 'Kannada', color: '#FECACA' },
  ENG: { name: 'English', color: '#BFDBFE' },
  HIN: { name: 'Hindi', color: '#BBF7D0' },
  MAT: { name: 'Mathematics', color: '#FEF08A' },
  SCI: { name: 'Science', color: '#E9D5FF' },
  SOC: { name: 'Social Studies', color: '#A5F3FC' },
  CS: { name: 'Computer Science', color: '#FBCFE8' },
  DRW: { name: 'Drawing & Visual Arts', color: '#FED7AA' },
  MUS: { name: 'Music & Performing Arts', color: '#99F6E4' },
  PE: { name: 'Physical Education', color: '#ECFCCB' },
  LIB: { name: 'Library & Reading', color: '#BAE6FD' },
  CUL: { name: 'Cultural Programme', color: '#C7D2FE' },
  BRK: { name: 'Short Break', color: '#F1F5F9' },
  LUN: { name: 'Lunch Break', color: '#F1F5F9' },
  ASM: { name: 'Morning Assembly', color: '#FEF9C3' },
  PTR: { name: 'Physical Training', color: '#DCFCE7' },
  BRF: { name: 'Breakfast', color: '#FEF3C7' },
};

export const BREAK_CODES = new Set(['BRK', 'LUN']);
export const ACTIVITY_CODES = new Set(['ASM', 'PTR', 'BRF']);

export const BREAK_TIMES: Record<string, string> = {
  BRK: '12:00 - 12:10',
  LUN: '1:30 - 2:20',
};

const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as const;

const RAW: Record<string, Record<string, string[]>> = {
  'Class 6': {
    Mon: ['KAN', 'CS', 'HIN', 'SCI', 'MAT', 'LIB', 'SOC', 'ENG'],
    Tue: ['SOC', 'DRW', 'KAN', 'MAT', 'SCI', 'HIN', 'CS', 'ENG'],
    Wed: ['SCI', 'MAT', 'DRW', 'ENG', 'HIN', 'SOC', 'MUS', 'KAN'],
    Thu: ['MAT', 'KAN', 'MUS', 'SCI', 'HIN', 'SOC', 'PE', 'ENG'],
    Fri: ['ENG', 'MAT', 'LIB', 'SCI', 'SOC', 'HIN', 'KAN', 'CUL'],
    Sat: ['PE', 'CS', 'ENG', 'MUS'],
  },
  'Class 7': {
    Mon: ['ENG', 'PE', 'SOC', 'HIN', 'KAN', 'MAT', 'MUS', 'SCI'],
    Tue: ['KAN', 'MAT', 'LIB', 'ENG', 'HIN', 'PE', 'SCI', 'SOC'],
    Wed: ['SOC', 'CS', 'KAN', 'HIN', 'ENG', 'MAT', 'SCI', 'LIB'],
    Thu: ['SCI', 'MUS', 'ENG', 'CS', 'SOC', 'MAT', 'KAN', 'DRW'],
    Fri: ['MAT', 'CS', 'HIN', 'SOC', 'DRW', 'KAN', 'SCI', 'CUL'],
    Sat: ['ENG', 'MAT', 'SCI', 'SOC'],
  },
  'Class 8': {
    Mon: ['HIN', 'SOC', 'SCI', 'MAT', 'ENG', 'KAN', 'CS', 'DRW'],
    Tue: ['ENG', 'SCI', 'SOC', 'DRW', 'MAT', 'MUS', 'KAN', 'HIN'],
    Wed: ['KAN', 'ENG', 'MUS', 'SCI', 'SOC', 'LIB', 'MAT', 'PE'],
    Thu: ['SOC', 'SCI', 'HIN', 'MAT', 'LIB', 'KAN', 'ENG', 'PE'],
    Fri: ['SCI', 'KAN', 'CS', 'MAT', 'HIN', 'ENG', 'SOC', 'CUL'],
    Sat: ['MAT', 'SCI', 'SOC', 'KAN'],
  },
  'Class 9': {
    Mon: ['MAT', 'ENG', 'LIB', 'KAN', 'PE', 'SCI', 'HIN', 'SOC'],
    Tue: ['HIN', 'ENG', 'SCI', 'CS', 'LIB', 'KAN', 'SOC', 'MAT'],
    Wed: ['ENG', 'SCI', 'CS', 'SOC', 'MAT', 'DRW', 'KAN', 'HIN'],
    Thu: ['KAN', 'SOC', 'DRW', 'ENG', 'MAT', 'HIN', 'SCI', 'MUS'],
    Fri: ['KAN', 'SOC', 'SCI', 'ENG', 'PE', 'MUS', 'MAT', 'CUL'],
    Sat: ['SCI', 'SOC', 'KAN', 'MAT'],
  },
  'Class 10': {
    Mon: ['SCI', 'MAT', 'DRW', 'ENG', 'SOC', 'CS', 'KAN', 'MUS'],
    Tue: ['MAT', 'HIN', 'MUS', 'KAN', 'SOC', 'SCI', 'ENG', 'PE'],
    Wed: ['HIN', 'SOC', 'MAT', 'KAN', 'LIB', 'SCI', 'PE', 'ENG'],
    Thu: ['ENG', 'MAT', 'CS', 'SOC', 'KAN', 'SCI', 'HIN', 'LIB'],
    Fri: ['SOC', 'ENG', 'DRW', 'KAN', 'MAT', 'SCI', 'HIN', 'CUL'],
    Sat: ['SOC', 'KAN', 'MAT', 'SCI'],
  },
};

export const WEEKLY_TIMETABLE: ClassSchedule[] = Object.entries(RAW).map(([name, days]) => ({
  name,
  days: DAYS.map(day => {
    const rawPeriods = days[day].map(code => ({ code, name: SUBJECT_INFO[code]?.name ?? code }));
    if (day === 'Sat') {
      return {
        label: day,
        periods: [
          { code: 'ASM', name: SUBJECT_INFO['ASM'].name },
          { code: 'PTR', name: SUBJECT_INFO['PTR'].name },
          { code: 'BRF', name: SUBJECT_INFO['BRF'].name },
          ...rawPeriods.slice(0, 2),
          { code: 'BRK', name: SUBJECT_INFO['BRK'].name },
          ...rawPeriods.slice(2),
        ],
      };
    }
    return {
      label: day,
      periods: [
        { code: 'ASM', name: SUBJECT_INFO['ASM'].name },
        ...rawPeriods.slice(0, 3),
        { code: 'BRK', name: SUBJECT_INFO['BRK'].name },
        ...rawPeriods.slice(3, 5),
        { code: 'LUN', name: SUBJECT_INFO['LUN'].name },
        ...rawPeriods.slice(5),
      ],
    };
  }),
}));

export const WEEKDAY_TIMES = [
  '9:40–10:00',
  '10:00–10:40',
  '10:40–11:20',
  '11:20–12:00',
  '12:00–12:10',
  '12:10–12:50',
  '12:50–1:30',
  '1:30–2:20',
  '2:20–3:00',
  '3:00–3:40',
  '3:40–4:20',
];

export const SAT_TIMES = [
  '8:30–8:40',
  '8:40–9:10',
  '9:10–9:40',
  '9:40–10:20',
  '10:20–11:00',
  '11:00–11:10',
  '11:10–11:50',
  '11:50–12:30',
];

export const DAY_LABELS = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
