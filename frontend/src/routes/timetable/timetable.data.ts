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
  KAN: { name: 'Kannada', color: '#DBEAFE' },
  ENG: { name: 'English', color: '#E0E7FF' },
  HIN: { name: 'Hindi', color: '#F3E8FF' },
  MAT: { name: 'Mathematics', color: '#FFF3E0' },
  SCI: { name: 'Science', color: '#FFF9C4' },
  SOC: { name: 'Social Studies', color: '#FFE0B2' },
  CS: { name: 'Computer Science', color: '#B2EBF2' },
  DRW: { name: 'Drawing & Visual Arts', color: '#FCE4EC' },
  MUS: { name: 'Music & Performing Arts', color: '#FFE0E0' },
  PE: { name: 'Physical Education', color: '#C8E6C9' },
  LIB: { name: 'Library & Reading', color: '#E0E0E0' },
  CUL: { name: 'Cultural Programme', color: '#D4EDDA' },
  BRK: { name: 'Short Break', color: '#F1F5F9' },
  LUN: { name: 'Lunch Break', color: '#F1F5F9' },
  ASM: { name: 'Morning Assembly', color: '#FFF9C4' },
  PTR: { name: 'Physical Training', color: '#C8E6C9' },
  BRF: { name: 'Breakfast', color: '#FFE0B2' },
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
    Tue: ['KAN', 'MAT', 'PE', 'ENG', 'HIN', 'LIB', 'SCI', 'SOC'],
    Wed: ['SOC', 'CS', 'KAN', 'HIN', 'ENG', 'MAT', 'SCI', 'LIB'],
    Thu: ['SCI', 'MUS', 'ENG', 'CS', 'SOC', 'MAT', 'KAN', 'DRW'],
    Fri: ['MAT', 'SCI', 'HIN', 'SOC', 'DRW', 'KAN', 'ENG', 'CUL'],
    Sat: ['KAN', 'MAT', 'SCI', 'SOC'],
  },
  'Class 8': {
    Mon: ['HIN', 'SOC', 'SCI', 'MAT', 'ENG', 'DRW', 'CS', 'KAN'],
    Tue: ['ENG', 'SCI', 'SOC', 'DRW', 'MAT', 'MUS', 'KAN', 'HIN'],
    Wed: ['KAN', 'ENG', 'PE', 'SCI', 'SOC', 'LIB', 'MAT', 'MUS'],
    Thu: ['SOC', 'SCI', 'HIN', 'MAT', 'LIB', 'KAN', 'ENG', 'PE'],
    Fri: ['SCI', 'CS', 'KAN', 'MAT', 'HIN', 'ENG', 'SOC', 'CUL'],
    Sat: ['MAT', 'SCI', 'SOC', 'KAN'],
  },
  'Class 9': {
    Mon: ['MAT', 'ENG', 'PE', 'KAN', 'LIB', 'SCI', 'HIN', 'SOC'],
    Tue: ['HIN', 'ENG', 'SCI', 'CS', 'LIB', 'KAN', 'SOC', 'MAT'],
    Wed: ['ENG', 'SCI', 'MAT', 'SOC', 'DRW', 'CS', 'KAN', 'HIN'],
    Thu: ['KAN', 'SOC', 'DRW', 'ENG', 'MAT', 'HIN', 'SCI', 'MUS'],
    Fri: ['KAN', 'SOC', 'SCI', 'ENG', 'PE', 'MUS', 'MAT', 'CUL'],
    Sat: ['SCI', 'SOC', 'KAN', 'MAT'],
  },
  'Class 10': {
    Mon: ['SCI', 'MAT', 'DRW', 'ENG', 'SOC', 'CS', 'KAN', 'MUS'],
    Tue: ['MAT', 'HIN', 'MUS', 'KAN', 'SOC', 'SCI', 'ENG', 'PE'],
    Wed: ['HIN', 'SOC', 'LIB', 'KAN', 'MAT', 'ENG', 'PE', 'SCI'],
    Thu: ['ENG', 'MAT', 'CS', 'SOC', 'KAN', 'LIB', 'HIN', 'SCI'],
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
