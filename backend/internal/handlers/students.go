package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pragati/backend/internal/middleware"
	"github.com/pragati/backend/internal/models"
	"github.com/rs/zerolog/log"
)

var satsRegex = regexp.MustCompile(`^\d{9}$`)

type StudentHandler struct {
	db *pgxpool.Pool
}

func NewStudentHandler(db *pgxpool.Pool) *StudentHandler {
	return &StudentHandler{db: db}
}

func (h *StudentHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	query := `SELECT s.id, s.school_id, s.user_id, s.sats_number, s.admission_no,
		s.roll_no, s.first_name, s.last_name, s.date_of_birth, s.gender,
		s.photo_url, s.blood_group, s.address, s.phone, s.email,
		s.class_id, s.section_id, s.house_id, s.academic_year_id,
		s.parent_name, s.parent_phone, s.parent_email, s.is_active,
		s.created_at, s.updated_at
		FROM students s
		WHERE s.school_id = $1 AND s.deleted_at IS NULL`

	args := []interface{}{claims.SchoolID}
	argIdx := 2

	if classID := r.URL.Query().Get("class_id"); classID != "" {
		query += fmt.Sprintf(" AND s.class_id = $%d", argIdx)
		args = append(args, classID)
		argIdx++
	}
	if sectionID := r.URL.Query().Get("section_id"); sectionID != "" {
		query += fmt.Sprintf(" AND s.section_id = $%d", argIdx)
		args = append(args, sectionID)
		argIdx++
	}
	if search := r.URL.Query().Get("search"); search != "" {
		query += fmt.Sprintf(` AND (s.sats_number ILIKE $%d OR s.first_name ILIKE $%d OR s.last_name ILIKE $%d)`, argIdx, argIdx+1, argIdx+2)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIdx += 3
	}
	if academicYearID := r.URL.Query().Get("academic_year_id"); academicYearID != "" {
		query += fmt.Sprintf(" AND s.academic_year_id = $%d", argIdx)
		args = append(args, academicYearID)
		argIdx++
	}

	var total int
	countQuery := strings.Replace(query, "SELECT s.id, s.school_id", "SELECT COUNT(*)", 1)
	err := h.db.QueryRow(r.Context(), countQuery, args...).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("count students failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to count students"},
		})
		return
	}

	query += fmt.Sprintf(" ORDER BY s.roll_no ASC, s.first_name ASC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("list students failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to fetch students"},
		})
		return
	}
	defer rows.Close()

	students := []models.Student{}
	for rows.Next() {
		var s models.Student
		err := rows.Scan(&s.ID, &s.SchoolID, &s.UserID, &s.SATSNumber, &s.AdmissionNo,
			&s.RollNo, &s.FirstName, &s.LastName, &s.DateOfBirth, &s.Gender,
			&s.PhotoURL, &s.BloodGroup, &s.Address, &s.Phone, &s.Email,
			&s.ClassID, &s.SectionID, &s.HouseID, &s.AcademicYearID,
			&s.ParentName, &s.ParentPhone, &s.ParentEmail, &s.IsActive,
			&s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			log.Error().Err(err).Msg("scan student row failed")
			continue
		}
		students = append(students, s)
	}

	renderJSON(w, http.StatusOK, models.APIResponse{
		Data: students,
		Meta: models.Pagination{
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
	})
}

func (h *StudentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	var s models.Student
	err := h.db.QueryRow(r.Context(),
		`SELECT id, school_id, user_id, sats_number, admission_no,
			roll_no, first_name, last_name, date_of_birth, gender,
			photo_url, blood_group, address, phone, email,
			class_id, section_id, house_id, academic_year_id,
			parent_name, parent_phone, parent_email, is_active,
			created_at, updated_at
		 FROM students WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		id, claims.SchoolID,
	).Scan(&s.ID, &s.SchoolID, &s.UserID, &s.SATSNumber, &s.AdmissionNo,
		&s.RollNo, &s.FirstName, &s.LastName, &s.DateOfBirth, &s.Gender,
		&s.PhotoURL, &s.BloodGroup, &s.Address, &s.Phone, &s.Email,
		&s.ClassID, &s.SectionID, &s.HouseID, &s.AcademicYearID,
		&s.ParentName, &s.ParentPhone, &s.ParentEmail, &s.IsActive,
		&s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "student not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: s})
}

func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	var req models.CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	if !satsRegex.MatchString(req.SATSNumber) {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "SATS number must be exactly 9 digits"},
		})
		return
	}

	if req.FirstName == "" {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "first name is required"},
		})
		return
	}

	var dob *time.Time
	if req.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			dob = &t
		}
	}

	houseID := nullIfEmpty(req.HouseID)
	sectionID := nullIfEmpty(req.SectionID)

	id := uuid.New().String()
	_, err := h.db.Exec(r.Context(),
		`INSERT INTO students (id, school_id, sats_number, admission_no, roll_no,
			first_name, last_name, date_of_birth, gender, blood_group,
			address, phone, email, class_id, section_id, house_id,
			academic_year_id, parent_name, parent_phone, parent_email,
			is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,true,NOW(),NOW())`,
		id, claims.SchoolID, req.SATSNumber, req.AdmissionNo, req.RollNo,
		req.FirstName, req.LastName, dob, req.Gender, req.BloodGroup,
		req.Address, req.Phone, req.Email, req.ClassID, sectionID, houseID,
		req.AcademicYearID, req.ParentName, req.ParentPhone, req.ParentEmail,
	)
	if err != nil {
		if isDuplicateError(err) {
			renderJSON(w, http.StatusConflict, models.APIResponse{
				Error: &models.APIError{Code: "DUPLICATE", Message: "a student with this SATS number already exists"},
			})
			return
		}
		log.Error().Err(err).Msg("create student failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to create student"},
		})
		return
	}

	renderJSON(w, http.StatusCreated, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *StudentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	var req models.UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	addField := func(field string, value interface{}) {
		argIdx++
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIdx))
		args = append(args, value)
	}

	args = append(args, id, claims.SchoolID)

	if req.FirstName != nil {
		addField("first_name", *req.FirstName)
	}
	if req.LastName != nil {
		addField("last_name", *req.LastName)
	}
	if req.Gender != nil {
		addField("gender", *req.Gender)
	}
	if req.Phone != nil {
		addField("phone", *req.Phone)
	}
	if req.Email != nil {
		addField("email", *req.Email)
	}
	if req.Address != nil {
		addField("address", *req.Address)
	}
	if req.ClassID != nil {
		addField("class_id", *req.ClassID)
	}
	if req.SectionID != nil {
		addField("section_id", nullIfEmpty(*req.SectionID))
	}
	if req.BloodGroup != nil {
		addField("blood_group", *req.BloodGroup)
	}
	if req.IsActive != nil {
		addField("is_active", *req.IsActive)
	}
	if req.AdmissionNo != nil {
		addField("admission_no", *req.AdmissionNo)
	}
	if req.RollNo != nil {
		addField("roll_no", *req.RollNo)
	}
	if req.ParentName != nil {
		addField("parent_name", *req.ParentName)
	}
	if req.ParentPhone != nil {
		addField("parent_phone", *req.ParentPhone)
	}
	if req.ParentEmail != nil {
		addField("parent_email", *req.ParentEmail)
	}
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		if t, err := time.Parse("2006-01-02", *req.DateOfBirth); err == nil {
			addField("date_of_birth", t)
		}
	}

	if len(setClauses) == 0 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "no fields to update"},
		})
		return
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	setSQL := strings.Join(setClauses, ", ")

	query := fmt.Sprintf(
		`UPDATE students SET %s WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		setSQL,
	)

	result, err := h.db.Exec(r.Context(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("update student failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to update student"},
		})
		return
	}

	if result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "student not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]string{"id": id}})
}

func (h *StudentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	claims := middleware.GetUserClaims(r.Context())

	result, err := h.db.Exec(r.Context(),
		`UPDATE students SET deleted_at = NOW(), is_active = false
		 WHERE id = $1 AND school_id = $2 AND deleted_at IS NULL`,
		id, claims.SchoolID,
	)
	if err != nil {
		log.Error().Err(err).Msg("delete student failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to delete student"},
		})
		return
	}

	if result.RowsAffected() == 0 {
		renderJSON(w, http.StatusNotFound, models.APIResponse{
			Error: &models.APIError{Code: "NOT_FOUND", Message: "student not found"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]bool{"success": true}})
}

func (h *StudentHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "file is required"},
		})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "failed to parse CSV"},
		})
		return
	}

	if len(records) < 2 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "CSV must have a header row and at least one data row"},
		})
		return
	}

	headers := records[0]
	colMap := make(map[string]int)
	for i, h := range headers {
		colMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	requiredCols := []string{"sats_number", "first_name", "class"}
	for _, col := range requiredCols {
		if _, ok := colMap[col]; !ok {
			renderJSON(w, http.StatusBadRequest, models.APIResponse{
				Error: &models.APIError{Code: "VALIDATION_ERROR", Message: fmt.Sprintf("missing required column: %s", col)},
			})
			return
		}
	}

	result := models.ImportResult{}

	for i := 1; i < len(records); i++ {
		row := records[i]
		getCol := func(col string) string {
			if idx, ok := colMap[col]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		satsNo := getCol("sats_number")
		firstName := getCol("first_name")

		if !satsRegex.MatchString(satsNo) {
			result.Errors = append(result.Errors, models.ImportRowError{
				Row: i + 1, SATS: satsNo, Field: "sats_number",
				Message: "SATS number must be exactly 9 digits",
			})
			result.Skipped++
			continue
		}

		if firstName == "" {
			result.Errors = append(result.Errors, models.ImportRowError{
				Row: i + 1, SATS: satsNo, Field: "first_name",
				Message: "first name is required",
			})
			result.Skipped++
			continue
		}

		classCode := getCol("class")
		var classID string
		err := h.db.QueryRow(r.Context(),
			`SELECT id FROM classes WHERE code = $1 OR name = $1 AND school_id = $2 AND deleted_at IS NULL`,
			classCode, claims.SchoolID,
		).Scan(&classID)
		if err != nil {
			result.Errors = append(result.Errors, models.ImportRowError{
				Row: i + 1, SATS: satsNo, Field: "class",
				Message: fmt.Sprintf("class not found: %s", classCode),
			})
			result.Skipped++
			continue
		}

		academicYearID := getCol("academic_year")
		if academicYearID == "" {
			h.db.QueryRow(r.Context(),
				`SELECT id FROM academic_years WHERE is_current = true AND school_id = $1 AND deleted_at IS NULL`,
				claims.SchoolID,
			).Scan(&academicYearID)
		}
		if academicYearID == "" {
			result.Errors = append(result.Errors, models.ImportRowError{
				Row: i + 1, SATS: satsNo, Field: "academic_year",
				Message: "academic year not found and no current year set",
			})
			result.Skipped++
			continue
		}

		sectionID := ""
		if sectionName := getCol("section"); sectionName != "" {
			h.db.QueryRow(r.Context(),
				`SELECT id FROM sections WHERE class_id = $1 AND name = $2 AND deleted_at IS NULL`,
				classID, sectionName,
			).Scan(&sectionID)
		}

		rollNo := 0
		if rn := getCol("roll_no"); rn != "" {
			rollNo, _ = strconv.Atoi(rn)
		}

		var dob *time.Time
		if dobStr := getCol("date_of_birth"); dobStr != "" {
			if t, err := time.Parse("2006-01-02", dobStr); err == nil {
				dob = &t
			}
		}

		id := uuid.New().String()
		_, err = h.db.Exec(r.Context(),
			`INSERT INTO students (id, school_id, sats_number, admission_no, roll_no,
				first_name, last_name, date_of_birth, gender, phone, email,
				class_id, section_id, academic_year_id,
				parent_name, parent_phone, parent_email,
				is_active, created_at, updated_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,true,NOW(),NOW())
			 ON CONFLICT (sats_number) DO UPDATE SET
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name,
				class_id = EXCLUDED.class_id,
				section_id = EXCLUDED.section_id,
				updated_at = NOW()`,
			id, claims.SchoolID, satsNo, getCol("admission_no"), rollNo,
			firstName, getCol("last_name"), dob, getCol("gender"),
			getCol("phone"), getCol("email"),
			classID, nullIfEmpty(sectionID), academicYearID,
			getCol("parent_name"), getCol("parent_phone"), getCol("parent_email"),
		)
		if err != nil {
			result.Errors = append(result.Errors, models.ImportRowError{
				Row: i + 1, SATS: satsNo, Field: "_db",
				Message: fmt.Sprintf("database error: %v", err),
			})
			result.Skipped++
			continue
		}

		result.Imported++
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: result})
}

func (h *StudentHandler) BulkUpdate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentIDs   []string `json:"student_ids"`
		ClassID      *string  `json:"class_id"`
		SectionID    *string  `json:"section_id"`
		AcademicYearID *string `json:"academic_year_id"`
		IsActive     *bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "INVALID_INPUT", Message: "invalid request body"},
		})
		return
	}

	if len(req.StudentIDs) == 0 {
		renderJSON(w, http.StatusBadRequest, models.APIResponse{
			Error: &models.APIError{Code: "VALIDATION_ERROR", Message: "student_ids is required"},
		})
		return
	}

	claims := middleware.GetUserClaims(r.Context())

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("begin transaction failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to process update"},
		})
		return
	}
	defer tx.Rollback(r.Context())

	for _, sid := range req.StudentIDs {
		if req.ClassID != nil {
			_, err := tx.Exec(r.Context(),
				`UPDATE students SET class_id = $1, updated_at = NOW() WHERE id = $2 AND school_id = $3 AND deleted_at IS NULL`,
				*req.ClassID, sid, claims.SchoolID,
			)
			if err != nil {
				log.Error().Err(err).Str("student_id", sid).Msg("bulk update failed")
			}
		}
		if req.SectionID != nil {
			tx.Exec(r.Context(),
				`UPDATE students SET section_id = $1, updated_at = NOW() WHERE id = $2 AND school_id = $3 AND deleted_at IS NULL`,
				nullIfEmpty(*req.SectionID), sid, claims.SchoolID,
			)
		}
		if req.AcademicYearID != nil {
			tx.Exec(r.Context(),
				`UPDATE students SET academic_year_id = $1, updated_at = NOW() WHERE id = $2 AND school_id = $3 AND deleted_at IS NULL`,
				*req.AcademicYearID, sid, claims.SchoolID,
			)
		}
		if req.IsActive != nil {
			tx.Exec(r.Context(),
				`UPDATE students SET is_active = $1, updated_at = NOW() WHERE id = $2 AND school_id = $3 AND deleted_at IS NULL`,
				*req.IsActive, sid, claims.SchoolID,
			)
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		log.Error().Err(err).Msg("commit bulk update failed")
		renderJSON(w, http.StatusInternalServerError, models.APIResponse{
			Error: &models.APIError{Code: "INTERNAL_ERROR", Message: "failed to commit update"},
		})
		return
	}

	renderJSON(w, http.StatusOK, models.APIResponse{Data: map[string]interface{}{
		"updated": len(req.StudentIDs),
	}})
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "23505")
}
