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

export const BREAK_CODES = new Set(['BRK', 'LUN']);
export const ACTIVITY_CODES = new Set(['ASM', 'PTR', 'BRF']);

export const BREAK_TIMES: Record<string, string> = {
  BRK: '12:00 - 12:10',
  LUN: '1:30 - 2:20',
};

const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as const;

const RAW: Record<string, Record<string, string[]>> = {
  'Class 6': {
    Mon: ['LIB', 'KAN', 'SOC', 'MAT', 'DRW', 'ENG', 'SCI', 'HIN'],
    Tue: ['LIB', 'KAN', 'HIN', 'MAT', 'ENG', 'SOC', 'SCI', 'PE'],
    Wed: ['SCI', 'MAT', 'SOC', 'ENG', 'HIN', 'CS', 'KAN', 'DRW'],
    Thu: ['ENG', 'SOC', 'SCI', 'KAN', 'HIN', 'PE', 'MUS', 'MAT'],
    Fri: ['SCI', 'CS', 'MAT', 'ENG', 'MUS', 'SOC', 'KAN', 'CUL'],
    Sat: ['CS', 'ENG', 'MUS', 'HIN'],
  },
  'Class 7': {
    Mon: ['HIN', 'ENG', 'SCI', 'CS', 'SOC', 'KAN', 'LIB', 'MAT'],
    Tue: ['ENG', 'SOC', 'SCI', 'KAN', 'MAT', 'HIN', 'PE', 'MUS'],
    Wed: ['DRW', 'MUS', 'KAN', 'SCI', 'ENG', 'SOC', 'HIN', 'MAT'],
    Thu: ['KAN', 'HIN', 'PE', 'MAT', 'DRW', 'SOC', 'SCI', 'ENG'],
    Fri: ['KAN', 'ENG', 'LIB', 'SOC', 'CS', 'MAT', 'SCI', 'CUL'],
    Sat: ['KAN', 'MAT', 'SCI', 'SOC'],
  },
  'Class 8': {
    Mon: ['PE', 'CS', 'MAT', 'HIN', 'ENG', 'SCI', 'KAN', 'SOC'],
    Tue: ['HIN', 'ENG', 'LIB', 'SCI', 'KAN', 'MAT', 'MUS', 'SOC'],
    Wed: ['ENG', 'HIN', 'MUS', 'KAN', 'SOC', 'MAT', 'DRW', 'SCI'],
    Thu: ['CS', 'MAT', 'KAN', 'ENG', 'SOC', 'SCI', 'HIN', 'PE'],
    Fri: ['SOC', 'LIB', 'SCI', 'MAT', 'KAN', 'DRW', 'ENG', 'CUL'],
    Sat: ['MAT', 'SCI', 'SOC', 'KAN'],
  },
  'Class 9': {
    Mon: ['CS', 'MAT', 'ENG', 'KAN', 'SCI', 'SOC', 'HIN', 'LIB'],
    Tue: ['KAN', 'PE', 'MAT', 'SOC', 'DRW', 'SCI', 'HIN', 'ENG'],
    Wed: ['PE', 'ENG', 'LIB', 'HIN', 'MAT', 'SCI', 'SOC', 'KAN'],
    Thu: ['SCI', 'ENG', 'MAT', 'SOC', 'CS', 'HIN', 'KAN', 'MUS'],
    Fri: ['ENG', 'SOC', 'MUS', 'KAN', 'DRW', 'SCI', 'MAT', 'CUL'],
    Sat: ['SCI', 'SOC', 'KAN', 'MAT'],
  },
  'Class 10': {
    Mon: ['SOC', 'SCI', 'HIN', 'ENG', 'KAN', 'DRW', 'MAT', 'MUS'],
    Tue: ['SCI', 'MAT', 'SOC', 'LIB', 'HIN', 'KAN', 'ENG', 'DRW'],
    Wed: ['HIN', 'SCI', 'MAT', 'SOC', 'KAN', 'ENG', 'LIB', 'PE'],
    Thu: ['HIN', 'PE', 'SOC', 'SCI', 'KAN', 'MAT', 'ENG', 'CS'],
    Fri: ['MAT', 'KAN', 'CS', 'SCI', 'ENG', 'MUS', 'SOC', 'CUL'],
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
  '9:40 - 10:00',
  '10:00 - 10:40',
  '10:40 - 11:20',
  '11:20 - 12:00',
  '12:00 - 12:10',
  '12:10 - 12:50',
  '12:50 - 1:30',
  '1:30 - 2:20',
  '2:20 - 3:00',
  '3:00 - 3:40',
  '3:40 - 4:20',
];

export const SAT_TIMES = [
  '8:30 - 8:40',
  '8:40 - 9:10',
  '9:10 - 9:40',
  '9:40 - 10:20',
  '10:20 - 11:00',
  '11:00 - 11:10',
  '11:10 - 11:50',
  '11:50 - 12:30',
];

export const DAY_LABELS = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
