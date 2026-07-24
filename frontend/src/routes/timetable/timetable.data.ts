export interface TimetableCell {
  code: string;
  name: string;
}

export const TEACHER_NAMES: Record<string, string> = {
  KAN: '',
  ENG: '',
  HIN: '',
  MAT: '',
  SCI: '',
  SOC: '',
  CS: '',
  DRW: '',
  MUS: '',
  PE: '',
  LIB: '',
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
  KAN: { name: 'Kannada', color: '#E0F2FE' },
  ENG: { name: 'English', color: '#E0E7FF' },
  HIN: { name: 'Hindi', color: '#F3E8FF' },
  MAT: { name: 'Mathematics', color: '#FFEDD5' },
  SCI: { name: 'Science', color: '#FEF08A' },
  SOC: { name: 'Social Studies', color: '#FED7AA' },
  CS: { name: 'Computer Science', color: '#CFFAFE' },
  DRW: { name: 'Drawing & Visual Arts', color: '#FCE7F3' },
  MUS: { name: 'Music & Performing Arts', color: '#FEE2E2' },
  PE: { name: 'Physical Education', color: '#D1FAE5' },
  LIB: { name: 'Library & Reading', color: '#E2E8F0' },
  CUL: { name: 'Cultural Programme', color: '#DCFCE7' },
  BRK: { name: 'Short Break', color: '#F1F5F9' },
  LUN: { name: 'Lunch Break', color: '#F1F5F9' },
  ASM: { name: 'Morning Assembly', color: '#FEF9C3' },
  PTR: { name: 'Physical Training', color: '#D1FAE5' },
  BRF: { name: 'Breakfast', color: '#FFF7ED' },
};

export const BREAK_CODES = new Set(['BRK', 'LUN', 'ASM']);
export const ACTIVITY_CODES = new Set(['ASM', 'PTR', 'BRF']);

export const BREAK_TIMES: Record<string, string> = {
  BRK: '12:00 – 12:10',
  LUN: '1:30 – 2:20',
};

const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as const;

const RAW: Record<string, Record<string, string[]>> = {
  'Class 6': {
    Mon: ['KAN', 'MAT', 'ENG', 'SCI', 'SOC', 'HIN', 'CS', 'DRW'],
    Tue: ['MAT', 'SCI', 'KAN', 'ENG', 'HIN', 'SOC', 'PE', 'MUS'],
    Wed: ['ENG', 'KAN', 'MAT', 'SOC', 'SCI', 'CS', 'LIB', 'HIN'],
    Thu: ['SCI', 'SOC', 'ENG', 'KAN', 'MAT', 'DRW', 'MUS', 'PE'],
    Fri: ['HIN', 'KAN', 'SCI', 'MAT', 'SOC', 'ENG', 'CUL', 'CUL'],
    Sat: ['MAT', 'LIB', 'KAN', 'SCI'],
  },
  'Class 7': {
    Mon: ['MAT', 'KAN', 'SCI', 'ENG', 'SOC', 'PE', 'HIN', 'SOC'],
    Tue: ['ENG', 'MAT', 'HIN', 'KAN', 'SCI', 'DRW', 'PE', 'CS'],
    Wed: ['KAN', 'SCI', 'ENG', 'SOC', 'MAT', 'MUS', 'LIB', 'DRW'],
    Thu: ['SOC', 'HIN', 'MAT', 'KAN', 'ENG', 'CS', 'MUS', 'SCI'],
    Fri: ['SCI', 'ENG', 'KAN', 'MAT', 'SOC', 'HIN', 'CUL', 'CUL'],
    Sat: ['KAN', 'LIB', 'MAT', 'SCI'],
  },
  'Class 8': {
    Mon: ['ENG', 'SCI', 'MAT', 'KAN', 'SOC', 'MUS', 'SOC', 'LIB'],
    Tue: ['KAN', 'MAT', 'SCI', 'SOC', 'ENG', 'CS', 'DRW', 'PE'],
    Wed: ['MAT', 'ENG', 'KAN', 'SCI', 'HIN', 'DRW', 'SOC', 'MUS'],
    Thu: ['HIN', 'KAN', 'MAT', 'ENG', 'SCI', 'LIB', 'PE', 'CS'],
    Fri: ['SCI', 'SOC', 'KAN', 'ENG', 'MAT', 'HIN', 'CUL', 'CUL'],
    Sat: ['KAN', 'HIN', 'SCI', 'MAT'],
  },
  'Class 9': {
    Mon: ['SCI', 'SOC', 'ENG', 'KAN', 'MAT', 'LIB', 'HIN', 'SOC'],
    Tue: ['MAT', 'SCI', 'ENG', 'HIN', 'KAN', 'MUS', 'CS', 'DRW'],
    Wed: ['KAN', 'MAT', 'SOC', 'SCI', 'HIN', 'PE', 'DRW', 'LIB'],
    Thu: ['MAT', 'KAN', 'MAT', 'SOC', 'ENG', 'CS', 'MUS', 'SCI'],
    Fri: ['ENG', 'KAN', 'SCI', 'MAT', 'SOC', 'HIN', 'CUL', 'CUL'],
    Sat: ['SCI', 'PE', 'KAN', 'ENG'],
  },
  'Class 10': {
    Mon: ['MAT', 'SCI', 'KAN', 'ENG', 'HIN', 'CS', 'SOC', 'DRW'],
    Tue: ['SCI', 'MAT', 'ENG', 'KAN', 'SOC', 'LIB', 'MUS', 'PE'],
    Wed: ['HIN', 'ENG', 'SOC', 'MAT', 'KAN', 'SCI', 'CS', 'MUS'],
    Thu: ['ENG', 'KAN', 'HIN', 'SOC', 'MAT', 'PE', 'LIB', 'SCI'],
    Fri: ['KAN', 'ENG', 'SCI', 'MAT', 'HIN', 'SOC', 'CUL', 'CUL'],
    Sat: ['SCI', 'DRW', 'MAT', 'KAN'],
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
          ...rawPeriods,
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
  '9:45 – 9:55',
  '10:00 – 10:40',
  '10:40 – 11:20',
  '11:20 – 12:00',
  '12:00 – 12:10',
  '12:10 – 12:50',
  '12:50 – 1:30',
  '1:30 – 2:20',
  '2:20 – 3:00',
  '3:00 – 3:40',
  '3:40 – 4:20',
];

export const SAT_TIMES = [
  '8:30 – 8:40',
  '8:40 – 9:10',
  '9:10 – 9:50',
  '9:50 – 10:30',
  '10:30 – 11:10',
  '11:10 – 11:50',
  '11:50 – 12:30',
];

export const DAY_LABELS = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
