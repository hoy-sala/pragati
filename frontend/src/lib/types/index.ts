export type UserRole = 'admin' | 'principal' | 'teacher' | 'special_educator' | 'student' | 'parent';

export interface User {
	id: string;
	school_id: string;
	email: string;
	name: string;
	role: UserRole;
	phone: string;
	avatar_url?: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface Student {
	id: string;
	school_id: string;
	user_id?: string;
	sats_number: string;
	admission_no?: string;
	roll_no?: number;
	first_name: string;
	last_name?: string;
	date_of_birth?: string;
	gender?: string;
	photo_url?: string;
	blood_group?: string;
	address?: string;
	phone?: string;
	email?: string;
	class_id: string;
	section_id?: string;
	house_id?: string;
	academic_year_id: string;
	parent_name?: string;
	parent_phone?: string;
	parent_email?: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface Class {
	id: string;
	school_id: string;
	academic_year_id?: string;
	name: string;
	code?: string;
	sort_order: number;
	created_at: string;
}

export interface Subject {
	id: string;
	school_id: string;
	name: string;
	code?: string;
	is_language: boolean;
	is_core: boolean;
}

export interface AcademicYear {
	id: string;
	school_id: string;
	name: string;
	start_date: string;
	end_date: string;
	is_current: boolean;
}

export interface AssessmentCategory {
	id: string;
	school_id: string;
	name: string;
	code?: string;
	weightage: number;
	sort_order: number;
	is_active: boolean;
}

export interface Assessment {
	id: string;
	school_id: string;
	category_id: string;
	subject_id: string;
	teacher_id: string;
	class_id: string;
	section_id?: string;
	name?: string;
	max_marks: number;
	weightage: number;
	date?: string;
	chapters: string[];
	academic_year_id: string;
	is_published: boolean;
	is_locked: boolean;
	version: number;
	created_at: string;
}

export interface MarkGridRow {
	student_id: string;
	sats_number: string;
	name: string;
	roll_no: number;
	mark_id?: string;
	marks_obtained: number;
	is_absent: boolean;
	remarks: string;
}

export interface MarkInput {
	student_id: string;
	marks_obtained: number;
	is_absent: boolean;
	remarks: string;
}

export interface Question {
	id: string;
	school_id: string;
	subject_id: string;
	teacher_id: string;
	question_type: 'mcq' | 'true_false' | 'fill_blank' | 'short_answer';
	question_text: string;
	question_image?: string;
	options?: Option[];
	answer: string;
	marks: number;
	difficulty: string;
	chapters: string[];
	tags: string[];
	explanation?: string;
	is_active: boolean;
	created_at: string;
}

export interface Option {
	key: string;
	value: string;
	correct: boolean;
}

export interface QuizAssignment {
	id: string;
	school_id: string;
	title: string;
	description: string;
	target_type: 'student' | 'staff';
	target_id?: string;
	pass_pct: number;
	max_attempts: number;
	duration_min?: number | null;
	shuffle_questions: boolean;
	shuffle_options: boolean;
	show_result: boolean;
	start_at?: string | null;
	end_at?: string | null;
	is_published: boolean;
	created_by: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface QuizCreateInput {
	title: string;
	description: string;
	target_type: 'student' | 'staff';
	target_id?: string;
	pass_pct: number;
	max_attempts: number;
	duration_min?: number | null;
	shuffle_questions: boolean;
	shuffle_options: boolean;
	show_result: boolean;
	start_at?: string;
	end_at?: string;
}

export interface QuizListItem {
	id: string;
	title: string;
	description: string;
	target_type: string;
	target_id?: string;
	pass_pct: number;
	max_attempts: number;
	duration_min?: number | null;
	shuffle_questions: boolean;
	shuffle_options: boolean;
	show_result: boolean;
	start_at?: string | null;
	end_at?: string | null;
	is_published: boolean;
	question_count: number;
	attempt_count: number;
	created_by_name: string;
	created_at: string;
}

export interface AvailableQuizItem {
	id: string;
	title: string;
	description: string;
	duration_min?: number | null;
	pass_pct: number;
	max_attempts: number;
	question_count: number;
	start_at?: string | null;
	end_at?: string | null;
	attempts_used: number;
	last_status?: string;
	last_score?: number | null;
	last_passed?: boolean | null;
}

export interface QuizAttempt {
	id: string;
	quiz_id: string;
	user_id: string;
	attempt_no: number;
	status: 'in_progress' | 'submitted' | 'graded';
	score?: number | null;
	percentage?: number | null;
	passed?: boolean | null;
	started_at: string;
	submitted_at?: string | null;
	graded_at?: string | null;
}

export interface QuizResponse {
	id: string;
	attempt_id: string;
	question_id: string;
	selected_options?: string[];
	text_answer: string;
	is_correct?: boolean | null;
	marks_awarded: number;
	marks_total: number;
	graded_at?: string | null;
	graded_by?: string;
}

export interface QuizResponseDetail extends QuizResponse {
	question_text: string;
	question_type: string;
	options: Option[];
	correct_answer: string;
}

export interface QuizResultData {
	attempt: QuizAttempt;
	quiz: QuizAssignment;
	responses: QuizResponseDetail[];
	total_marks: number;
	total_awarded: number;
}

export interface Pagination {
	offset: number;
	limit: number;
	total: number;
}

export interface APIError {
	code: string;
	message: string;
	details?: unknown;
}

export interface APIResponse<T = unknown> {
	data?: T;
	meta?: Pagination;
	error?: APIError;
}
