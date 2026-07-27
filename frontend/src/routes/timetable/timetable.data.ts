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
    Mon: ['MAT', 'SOC', 'LIB', 'SCI', 'ENG', 'HIN', 'KAN', 'CS'],
    Tue: ['LIB', 'HIN', 'SCI', 'MUS', 'KAN', 'MAT', 'ENG', 'SOC'],
    Wed: ['SOC', 'SCI', 'PE', 'MAT', 'KAN', 'HIN', 'ENG', 'DRW'],
    Thu: ['MAT', 'PE', 'KAN', 'HIN', 'ENG', 'SOC', 'SCI', 'DRW'],
    Fri: ['CS', 'KAN', 'SCI', 'MUS', 'ENG', 'MAT', 'SOC', 'CUL'],
    Sat: ['CS', 'MUS', 'ENG', 'HIN'],
  },
  'Class 7': {
    Mon: ['ENG', 'MAT', 'DRW', 'HIN', 'KAN', 'SCI', 'SOC', 'LIB'],
    Tue: ['MAT', 'KAN', 'CS', 'ENG', 'SOC', 'HIN', 'MUS', 'SCI'],
    Wed: ['SCI', 'ENG', 'LIB', 'CS', 'HIN', 'SOC', 'KAN', 'MAT'],
    Thu: ['SOC', 'ENG', 'HIN', 'SCI', 'PE', 'MAT', 'DRW', 'KAN'],
    Fri: ['KAN', 'SCI', 'MUS', 'MAT', 'SOC', 'PE', 'ENG', 'CUL'],
    Sat: ['MAT', 'SOC', 'SCI', 'KAN'],
  },
  'Class 8': {
    Mon: ['HIN', 'ENG', 'KAN', 'LIB', 'MAT', 'SOC', 'SCI', 'DRW'],
    Tue: ['SOC', 'CS', 'ENG', 'MAT', 'HIN', 'KAN', 'SCI', 'MUS'],
    Wed: ['HIN', 'MUS', 'ENG', 'SCI', 'MAT', 'PE', 'SOC', 'KAN'],
    Thu: ['KAN', 'DRW', 'CS', 'SOC', 'MAT', 'HIN', 'ENG', 'SCI'],
    Fri: ['SCI', 'SOC', 'KAN', 'ENG', 'MAT', 'LIB', 'PE', 'CUL'],
    Sat: ['SCI', 'MAT', 'KAN', 'SOC'],
  },
  'Class 9': {
    Mon: ['LIB', 'SCI', 'SOC', 'MAT', 'HIN', 'ENG', 'CS', 'KAN'],
    Tue: ['DRW', 'MAT', 'HIN', 'KAN', 'SCI', 'ENG', 'SOC', 'PE'],
    Wed: ['MAT', 'HIN', 'KAN', 'SOC', 'ENG', 'SCI', 'MUS', 'CS'],
    Thu: ['HIN', 'SCI', 'PE', 'ENG', 'SOC', 'KAN', 'MAT', 'MUS'],
    Fri: ['ENG', 'LIB', 'MAT', 'DRW', 'KAN', 'SOC', 'SCI', 'CUL'],
    Sat: ['SOC', 'SCI', 'KAN', 'MAT'],
  },
  'Class 10': {
    Mon: ['SCI', 'KAN', 'CS', 'ENG', 'SOC', 'MAT', 'HIN', 'MUS'],
    Tue: ['ENG', 'SCI', 'MUS', 'SOC', 'MAT', 'CS', 'KAN', 'HIN'],
    Wed: ['KAN', 'MAT', 'SOC', 'HIN', 'LIB', 'ENG', 'DRW', 'SCI'],
    Thu: ['SCI', 'SOC', 'MAT', 'KAN', 'DRW', 'ENG', 'PE', 'HIN'],
    Fri: ['PE', 'MAT', 'LIB', 'SOC', 'SCI', 'ENG', 'KAN', 'CUL'],
    Sat: ['MAT', 'SOC', 'SCI', 'KAN'],
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
