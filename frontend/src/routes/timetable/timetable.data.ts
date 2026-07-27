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
    Mon: ['MAT', 'SCI', 'DRW', 'MAT', 'ENG', 'KAN', 'MUS', 'DRW'],
    Tue: ['LIB', 'CS', 'CS', 'MUS', 'PE', 'PE', 'SCI', 'SOC'],
    Wed: ['KAN', 'SCI', 'ENG', 'ENG', 'LIB', 'HIN', 'MAT', 'HIN'],
    Thu: ['KAN', 'KAN', 'SCI', 'HIN', 'HIN', 'SOC', 'ENG', 'SCI'],
    Fri: ['MAT', 'MAT', 'KAN', 'SOC', 'KAN', 'HIN', 'CUL', 'CUL'],
    Sat: ['SOC', 'SOC', 'ENG', 'ENG'],
  },
  'Class 7': {
    Mon: ['HIN', 'LIB', 'KAN', 'ENG', 'HIN', 'MUS', 'DRW', 'KAN'],
    Tue: ['SOC', 'KAN', 'HIN', 'SCI', 'MUS', 'KAN', 'PE', 'ENG'],
    Wed: ['ENG', 'DRW', 'HIN', 'LIB', 'HIN', 'PE', 'SOC', 'CS'],
    Thu: ['SOC', 'ENG', 'MAT', 'ENG', 'KAN', 'SCI', 'MAT', 'CS'],
    Fri: ['KAN', 'ENG', 'MAT', 'MAT', 'SCI', 'SCI', 'CUL', 'CUL'],
    Sat: ['MAT', 'SCI', 'SOC', 'SOC'],
  },
  'Class 8': {
    Mon: ['CS', 'PE', 'MAT', 'SOC', 'MUS', 'CS', 'SOC', 'SCI'],
    Tue: ['KAN', 'ENG', 'KAN', 'SOC', 'SCI', 'MUS', 'ENG', 'DRW'],
    Wed: ['LIB', 'KAN', 'MAT', 'HIN', 'ENG', 'KAN', 'ENG', 'SCI'],
    Thu: ['MAT', 'DRW', 'SOC', 'KAN', 'ENG', 'LIB', 'HIN', 'ENG'],
    Fri: ['PE', 'SOC', 'SCI', 'SCI', 'HIN', 'MAT', 'CUL', 'CUL'],
    Sat: ['KAN', 'MAT', 'HIN', 'HIN'],
  },
  'Class 9': {
    Mon: ['LIB', 'DRW', 'SCI', 'KAN', 'PE', 'SOC', 'PE', 'ENG'],
    Tue: ['SCI', 'DRW', 'MUS', 'LIB', 'ENG', 'SOC', 'KAN', 'SCI'],
    Wed: ['HIN', 'MAT', 'SCI', 'MAT', 'SOC', 'ENG', 'MUS', 'KAN'],
    Thu: ['CS', 'HIN', 'KAN', 'CS', 'SCI', 'HIN', 'KAN', 'SOC'],
    Fri: ['SOC', 'KAN', 'HIN', 'ENG', 'MAT', 'ENG', 'CUL', 'CUL'],
    Sat: ['HIN', 'ENG', 'MAT', 'MAT'],
  },
  'Class 10': {
    Mon: ['PE', 'MAT', 'CS', 'PE', 'SOC', 'SCI', 'SCI', 'MUS'],
    Tue: ['ENG', 'HIN', 'LIB', 'MAT', 'MAT', 'CS', 'SOC', 'MUS'],
    Wed: ['MAT', 'LIB', 'KAN', 'DRW', 'KAN', 'DRW', 'HIN', 'ENG'],
    Thu: ['HIN', 'SOC', 'ENG', 'SCI', 'MAT', 'KAN', 'SOC', 'KAN'],
    Fri: ['ENG', 'HIN', 'ENG', 'KAN', 'SOC', 'KAN', 'CUL', 'CUL'],
    Sat: ['ENG', 'HIN', 'SCI', 'SCI'],
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
