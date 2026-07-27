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
  LUN: { name: 'Lunch Break', color: '#F1F5F9' },
  ASM: { name: 'Morning Assembly', color: '#FEF9C3' },
  PTR: { name: 'Physical Training', color: '#D1FAE5' },
  BRF: { name: 'Breakfast', color: '#FFF7ED' },
};

export const BREAK_CODES = new Set(['LUN']);
export const ACTIVITY_CODES = new Set(['ASM', 'PTR', 'BRF']);

export const BREAK_TIMES: Record<string, string> = {
  LUN: '1:30 – 2:20',
};

const DAYS = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as const;

const RAW: Record<string, Record<string, string[]>> = {
  'Class 6': {
    Mon: ['MAT', 'HIN', 'SOC', 'KAN', 'MUS', 'ENG', 'SCI', 'CS'],
    Tue: ['MAT', 'KAN', 'SOC', 'LIB', 'SCI', 'HIN', 'ENG', 'PE'],
    Wed: ['MAT', 'CS', 'KAN', 'HIN', 'SOC', 'ENG', 'DRW', 'SCI'],
    Thu: ['MAT', 'LIB', 'SOC', 'MUS', 'ENG', 'HIN', 'SCI', 'KAN'],
    Fri: ['KAN', 'SOC', 'PE', 'MAT', 'SCI', 'ENG', 'DRW', 'CUL'],
    Sat: ['KAN', 'SOC', 'SCI', 'MAT'],
  },
  'Class 7': {
    Mon: ['PE', 'ENG', 'DRW', 'SCI', 'MAT', 'SOC', 'HIN', 'KAN'],
    Tue: ['ENG', 'LIB', 'DRW', 'SCI', 'MAT', 'KAN', 'HIN', 'SOC'],
    Wed: ['CS', 'KAN', 'HIN', 'ENG', 'PE', 'MAT', 'SCI', 'SOC'],
    Thu: ['SCI', 'ENG', 'KAN', 'LIB', 'SOC', 'MUS', 'MAT', 'HIN'],
    Fri: ['MAT', 'KAN', 'CS', 'ENG', 'SOC', 'MUS', 'SCI', 'CUL'],
    Sat: ['SOC', 'MAT', 'KAN', 'SCI'],
  },
  'Class 8': {
    Mon: ['KAN', 'SOC', 'MUS', 'DRW', 'ENG', 'HIN', 'MAT', 'SCI'],
    Tue: ['SCI', 'MUS', 'HIN', 'PE', 'ENG', 'SOC', 'KAN', 'MAT'],
    Wed: ['ENG', 'SCI', 'SOC', 'CS', 'KAN', 'HIN', 'MAT', 'LIB'],
    Thu: ['KAN', 'DRW', 'MAT', 'HIN', 'SCI', 'SOC', 'ENG', 'CS'],
    Fri: ['SCI', 'LIB', 'KAN', 'SOC', 'MAT', 'PE', 'ENG', 'CUL'],
    Sat: ['SOC', 'MAT', 'SCI', 'KAN'],
  },
  'Class 9': {
    Mon: ['SCI', 'MAT', 'HIN', 'LIB', 'SOC', 'KAN', 'DRW', 'ENG'],
    Tue: ['SOC', 'HIN', 'ENG', 'MAT', 'KAN', 'SCI', 'PE', 'CS'],
    Wed: ['HIN', 'ENG', 'MUS', 'MAT', 'SCI', 'SOC', 'CS', 'KAN'],
    Thu: ['HIN', 'SOC', 'ENG', 'DRW', 'PE', 'MAT', 'KAN', 'SCI'],
    Fri: ['LIB', 'MAT', 'SOC', 'KAN', 'ENG', 'SCI', 'MUS', 'CUL'],
    Sat: ['MAT', 'SCI', 'SOC', 'KAN'],
  },
  'Class 10': {
    Mon: ['ENG', 'SCI', 'KAN', 'MAT', 'HIN', 'CS', 'SOC', 'PE'],
    Tue: ['LIB', 'SCI', 'KAN', 'SOC', 'HIN', 'MAT', 'MUS', 'ENG'],
    Wed: ['SCI', 'SOC', 'MAT', 'PE', 'CS', 'KAN', 'ENG', 'HIN'],
    Thu: ['MUS', 'HIN', 'SCI', 'MAT', 'KAN', 'ENG', 'DRW', 'SOC'],
    Fri: ['SOC', 'ENG', 'MAT', 'SCI', 'KAN', 'DRW', 'LIB', 'CUL'],
    Sat: ['KAN', 'SOC', 'MAT', 'SCI'],
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
        ...rawPeriods.slice(0, 4),
        { code: 'LUN', name: SUBJECT_INFO['LUN'].name },
        ...rawPeriods.slice(4),
      ],
    };
  }),
}));

export const WEEKDAY_TIMES = [
  '9:45 – 9:55',
  '10:00 – 10:40',
  '10:40 – 11:20',
  '11:20 – 12:00',
  '12:00 – 12:40',
  '12:40 – 1:30',
  '1:30 – 2:10',
  '2:10 – 2:50',
  '2:50 – 3:30',
  '3:30 – 4:10',
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
