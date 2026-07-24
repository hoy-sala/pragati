package models

import "time"

type User struct {
	ID            string     `json:"id" db:"id"`
	SchoolID      string     `json:"school_id" db:"school_id"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	Name          string     `json:"name" db:"name"`
	Role          string     `json:"role" db:"role"`
	Phone         string     `json:"phone" db:"phone"`
	AvatarURL     string     `json:"avatar_url,omitempty" db:"avatar_url"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	AuditInfo
}

const (
	RoleAdmin          = "admin"
	RolePrincipal      = "principal"
	RoleTeacher        = "teacher"
	RoleSpecialEducator = "special_educator"
	RoleStudent        = "student"
	RoleParent         = "parent"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	User           *User  `json:"user"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	ExpiresIn      int64  `json:"expires_in"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin principal teacher special_educator student parent"`
	Phone    string `json:"phone"`
}
