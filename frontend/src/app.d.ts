/// <reference types="@sveltejs/kit" />

declare namespace App {
	interface User {
		id: string;
		school_id: string;
		email: string;
		name: string;
		role: 'admin' | 'principal' | 'teacher' | 'special_educator' | 'student' | 'parent';
		phone: string;
		avatar_url?: string;
		is_active: boolean;
		created_at: string;
		updated_at: string;
	}

	interface Student {
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

	interface Class {
		id: string;
		school_id: string;
		academic_year_id?: string;
		name: string;
		code?: string;
		sort_order: number;
		created_at: string;
	}

	interface Subject {
		id: string;
		school_id: string;
		name: string;
		code?: string;
		is_language: boolean;
		is_core: boolean;
	}

	interface AcademicYear {
		id: string;
		school_id: string;
		name: string;
		start_date: string;
		end_date: string;
		is_current: boolean;
	}

	interface APIResponse<T = unknown> {
		data?: T;
		meta?: Pagination;
		error?: APIError;
	}

	interface Pagination {
		offset: number;
		limit: number;
		total: number;
	}

	interface APIError {
		code: string;
		message: string;
		details?: unknown;
	}
}
