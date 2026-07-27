export interface TimetableCell {
  code: string;
  name: string;
}

export const TEACHER_NAMES: Record<string, string> = {
  KAN: 'KAN Teacher',
  ENG: 'ENG Teacher',
  HIN: 'HIN Teacher',
  MAT: 'MAT Teacher',
  SCI: 'SCI Teacher',
  SOC: 'SOC Teacher',
  CS: 'CS Teacher',
  DRW: 'DRW Teacher',
  MUS: 'MUS Teacher',
  PE: 'PE Teacher',
  LIB: 'LIB Teacher',
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
    Mon: ['PE', 'MAT', 'HIN', 'KAN', 'SOC', 'ENG', 'MUS', 'SCI'],
    Tue: ['MAT', 'SOC', 'MUS', 'KAN', 'SCI', 'ENG', 'HIN', 'DRW'],
    Wed: ['DRW', 'SCI', 'LIB', 'HIN', 'KAN', 'ENG', 'SOC', 'MAT'],
    Thu: ['MAT', 'SOC', 'SCI', 'ENG', 'KAN', 'LIB', 'HIN', 'CS'],
    Fri: ['MAT', 'KAN', 'ENG', 'SOC', 'SCI', 'HIN', 'CUL', 'CUL'],
    Sat: ['ENG', 'CS', 'PE', 'KAN'],
  },
  'Class 7': {
    Mon: ['HIN', 'KAN', 'LIB', 'SCI', 'ENG', 'SOC', 'MAT', 'MUS'],
    Tue: ['ENG', 'KAN', 'SOC', 'CS', 'HIN', 'MUS', 'MAT', 'SCI'],
    Wed: ['MAT', 'SOC', 'CS', 'SCI', 'HIN', 'DRW', 'ENG', 'KAN'],
    Thu: ['SOC', 'ENG', 'MAT', 'KAN', 'HIN', 'PE', 'DRW', 'SCI'],
    Fri: ['SCI', 'SOC', 'MAT', 'HIN', 'ENG', 'KAN', 'CUL', 'CUL'],
    Sat: ['KAN', 'PE', 'ENG', 'LIB'],
  },
  'Class 8': {
    Mon: ['SOC', 'PE', 'ENG', 'MAT', 'DRW', 'HIN', 'SCI', 'KAN'],
    Tue: ['KAN', 'ENG', 'CS', 'MAT', 'SOC', 'SCI', 'PE', 'HIN'],
    Wed: ['ENG', 'MAT', 'HIN', 'CS', 'SOC', 'KAN', 'SCI', 'MUS'],
    Thu: ['KAN', 'HIN', 'MUS', 'SCI', 'LIB', 'ENG', 'SOC', 'MAT'],
    Fri: ['SOC', 'ENG', 'HIN', 'SCI', 'KAN', 'MAT', 'CUL', 'CUL'],
    Sat: ['LIB', 'DRW', 'KAN', 'ENG'],
  },
  'Class 9': {
    Mon: ['MAT', 'SCI', 'MUS', 'LIB', 'HIN', 'KAN', 'ENG', 'SOC'],
    Tue: ['MUS', 'HIN', 'SCI', 'ENG', 'CS', 'KAN', 'SOC', 'MAT'],
    Wed: ['SOC', 'ENG', 'PE', 'KAN', 'SCI', 'HIN', 'MAT', 'CS'],
    Thu: ['PE', 'DRW', 'KAN', 'SOC', 'ENG', 'SCI', 'MAT', 'HIN'],
    Fri: ['HIN', 'SCI', 'KAN', 'MAT', 'SOC', 'ENG', 'CUL', 'CUL'],
    Sat: ['KAN', 'LIB', 'ENG', 'DRW'],
  },
  'Class 10': {
    Mon: ['SCI', 'SOC', 'CS', 'ENG', 'KAN', 'MAT', 'HIN', 'DRW'],
    Tue: ['HIN', 'MAT', 'LIB', 'SCI', 'ENG', 'CS', 'KAN', 'SOC'],
    Wed: ['SCI', 'PE', 'KAN', 'ENG', 'MAT', 'SOC', 'LIB', 'HIN'],
    Thu: ['ENG', 'PE', 'HIN', 'MAT', 'SCI', 'SOC', 'KAN', 'MUS'],
    Fri: ['KAN', 'HIN', 'SCI', 'ENG', 'MAT', 'SOC', 'CUL', 'CUL'],
    Sat: ['DRW', 'MUS', 'KAN', 'ENG'],
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
