package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jamaalul/proyek1pbl/app/model"
)

type ClinicRepository interface {
	CreatePolyclinic(ctx context.Context, db DBTX, p *model.Polyclinic) error
	ListPolyclinics(ctx context.Context, db DBTX) ([]model.Polyclinic, error)
	GetPolyclinicByID(ctx context.Context, db DBTX, id int64) (*model.Polyclinic, error)

	CreateDoctorUser(ctx context.Context, db DBTX, fullName, email, passwordHash string) (int64, error)
	CreateDoctor(ctx context.Context, db DBTX, doc *model.Doctor) error
	ListDoctors(ctx context.Context, db DBTX, polyclinicID *int64) ([]model.Doctor, error)
	GetDoctorByID(ctx context.Context, db DBTX, id int64) (*model.Doctor, error)
	GetDoctorByUserID(ctx context.Context, db DBTX, userID int64) (*model.Doctor, error)

	CreateSchedule(ctx context.Context, db DBTX, s *model.DoctorSchedule) error
	GetScheduleByID(ctx context.Context, db DBTX, id int64) (*model.DoctorSchedule, error)
	GetScheduleForUpdate(ctx context.Context, db DBTX, id int64) (*model.DoctorSchedule, error)
	ListSchedules(ctx context.Context, db DBTX, doctorID, polyclinicID *int64, date, from, to string) ([]model.DoctorSchedule, error)
	UpdateSchedule(ctx context.Context, db DBTX, s *model.DoctorSchedule) error
	CountActiveBookingsBySchedule(ctx context.Context, db DBTX, scheduleID int64) (int, error)
}

type clinicPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewClinicRepository(pool *pgxpool.Pool) ClinicRepository {
	return &clinicPostgresRepository{pool: pool}
}

func (r *clinicPostgresRepository) getDB(db DBTX) DBTX {
	if db != nil {
		return db
	}
	return r.pool
}

// --- Polyclinics ---

func (r *clinicPostgresRepository) CreatePolyclinic(ctx context.Context, db DBTX, p *model.Polyclinic) error {
	query := `
		INSERT INTO polyclinics (name, description, created_at, updated_at)
		VALUES ($1, $2, now(), now())
		RETURNING id, created_at, updated_at
	`
	return r.getDB(db).QueryRow(ctx, query, p.Name, p.Description).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *clinicPostgresRepository) ListPolyclinics(ctx context.Context, db DBTX) ([]model.Polyclinic, error) {
	query := `
		SELECT id, name, COALESCE(description, ''), created_at, updated_at
		FROM polyclinics
		ORDER BY id ASC
	`
	rows, err := r.getDB(db).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list polyclinics: %w", err)
	}
	defer rows.Close()

	var result []model.Polyclinic
	for rows.Next() {
		var p model.Polyclinic
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan polyclinic: %w", err)
		}
		result = append(result, p)
	}
	if result == nil {
		result = []model.Polyclinic{}
	}
	return result, nil
}

func (r *clinicPostgresRepository) GetPolyclinicByID(ctx context.Context, db DBTX, id int64) (*model.Polyclinic, error) {
	query := `
		SELECT id, name, COALESCE(description, ''), created_at, updated_at
		FROM polyclinics
		WHERE id = $1
	`
	var p model.Polyclinic
	err := r.getDB(db).QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// --- Doctors ---

func (r *clinicPostgresRepository) CreateDoctorUser(ctx context.Context, db DBTX, fullName, email, passwordHash string) (int64, error) {
	query := `
		INSERT INTO users (full_name, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, 'dokter', TRUE, now(), now())
		RETURNING id
	`
	var userID int64
	err := r.getDB(db).QueryRow(ctx, query, fullName, email, passwordHash).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *clinicPostgresRepository) CreateDoctor(ctx context.Context, db DBTX, doc *model.Doctor) error {
	query := `
		INSERT INTO doctors (user_id, polyclinic_id, specialization, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		RETURNING id, created_at, updated_at
	`
	return r.getDB(db).QueryRow(ctx, query, doc.UserID, doc.PolyclinicID, doc.Specialization).Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)
}

func (r *clinicPostgresRepository) ListDoctors(ctx context.Context, db DBTX, polyclinicID *int64) ([]model.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, u.full_name, u.email, d.polyclinic_id, p.name as polyclinic_name,
		       COALESCE(d.specialization, ''), d.created_at, d.updated_at
		FROM doctors d
		JOIN users u ON u.id = d.user_id
		JOIN polyclinics p ON p.id = d.polyclinic_id
	`
	var args []any
	if polyclinicID != nil {
		query += " WHERE d.polyclinic_id = $1"
		args = append(args, *polyclinicID)
	}
	query += " ORDER BY d.id ASC"

	rows, err := r.getDB(db).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list doctors: %w", err)
	}
	defer rows.Close()

	var result []model.Doctor
	for rows.Next() {
		var d model.Doctor
		if err := rows.Scan(&d.ID, &d.UserID, &d.FullName, &d.Email, &d.PolyclinicID, &d.PolyclinicName, &d.Specialization, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan doctor: %w", err)
		}
		result = append(result, d)
	}
	if result == nil {
		result = []model.Doctor{}
	}
	return result, nil
}

func (r *clinicPostgresRepository) GetDoctorByID(ctx context.Context, db DBTX, id int64) (*model.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, u.full_name, u.email, d.polyclinic_id, p.name as polyclinic_name,
		       COALESCE(d.specialization, ''), d.created_at, d.updated_at
		FROM doctors d
		JOIN users u ON u.id = d.user_id
		JOIN polyclinics p ON p.id = d.polyclinic_id
		WHERE d.id = $1
	`
	var d model.Doctor
	err := r.getDB(db).QueryRow(ctx, query, id).Scan(&d.ID, &d.UserID, &d.FullName, &d.Email, &d.PolyclinicID, &d.PolyclinicName, &d.Specialization, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *clinicPostgresRepository) GetDoctorByUserID(ctx context.Context, db DBTX, userID int64) (*model.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, u.full_name, u.email, d.polyclinic_id, p.name as polyclinic_name,
		       COALESCE(d.specialization, ''), d.created_at, d.updated_at
		FROM doctors d
		JOIN users u ON u.id = d.user_id
		JOIN polyclinics p ON p.id = d.polyclinic_id
		WHERE d.user_id = $1
	`
	var d model.Doctor
	err := r.getDB(db).QueryRow(ctx, query, userID).Scan(&d.ID, &d.UserID, &d.FullName, &d.Email, &d.PolyclinicID, &d.PolyclinicName, &d.Specialization, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// --- Schedules ---

func (r *clinicPostgresRepository) CreateSchedule(ctx context.Context, db DBTX, s *model.DoctorSchedule) error {
	query := `
		INSERT INTO doctor_schedules (doctor_id, practice_date, start_time, end_time, quota, created_at, updated_at)
		VALUES ($1, $2::date, $3::time, $4::time, $5, now(), now())
		RETURNING id, to_char(practice_date, 'YYYY-MM-DD'), to_char(start_time, 'HH24:MI'), to_char(end_time, 'HH24:MI'), created_at, updated_at
	`
	return r.getDB(db).QueryRow(ctx, query, s.DoctorID, s.PracticeDate, s.StartTime, s.EndTime, s.Quota).
		Scan(&s.ID, &s.PracticeDate, &s.StartTime, &s.EndTime, &s.CreatedAt, &s.UpdatedAt)
}

func (r *clinicPostgresRepository) GetScheduleByID(ctx context.Context, db DBTX, id int64) (*model.DoctorSchedule, error) {
	query := `
		SELECT s.id, s.doctor_id, u.full_name as doctor_name, d.polyclinic_id, p.name as polyclinic_name,
		       to_char(s.practice_date, 'YYYY-MM-DD'),
		       to_char(s.start_time, 'HH24:MI'),
		       to_char(s.end_time, 'HH24:MI'),
		       s.quota,
		       COALESCE(COUNT(a.id) FILTER (WHERE a.status <> 'cancelled'), 0) as booked_count,
		       s.quota - COALESCE(COUNT(a.id) FILTER (WHERE a.status <> 'cancelled'), 0) as remaining_quota,
		       s.created_at, s.updated_at
		FROM doctor_schedules s
		JOIN doctors d ON d.id = s.doctor_id
		JOIN users u ON u.id = d.user_id
		JOIN polyclinics p ON p.id = d.polyclinic_id
		LEFT JOIN appointments a ON a.schedule_id = s.id
		WHERE s.id = $1
		GROUP BY s.id, u.full_name, d.polyclinic_id, p.name
	`
	var s model.DoctorSchedule
	err := r.getDB(db).QueryRow(ctx, query, id).Scan(
		&s.ID, &s.DoctorID, &s.DoctorName, &s.PolyclinicID, &s.PolyclinicName,
		&s.PracticeDate, &s.StartTime, &s.EndTime, &s.Quota,
		&s.BookedCount, &s.RemainingQuota, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *clinicPostgresRepository) GetScheduleForUpdate(ctx context.Context, db DBTX, id int64) (*model.DoctorSchedule, error) {
	query := `
		SELECT id, doctor_id, to_char(practice_date, 'YYYY-MM-DD'),
		       to_char(start_time, 'HH24:MI'), to_char(end_time, 'HH24:MI'), quota,
		       created_at, updated_at
		FROM doctor_schedules
		WHERE id = $1
		FOR UPDATE
	`
	var s model.DoctorSchedule
	err := r.getDB(db).QueryRow(ctx, query, id).Scan(
		&s.ID, &s.DoctorID, &s.PracticeDate, &s.StartTime, &s.EndTime, &s.Quota, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *clinicPostgresRepository) ListSchedules(ctx context.Context, db DBTX, doctorID, polyclinicID *int64, date, from, to string) ([]model.DoctorSchedule, error) {
	query := `
		SELECT s.id, s.doctor_id, u.full_name as doctor_name, d.polyclinic_id, p.name as polyclinic_name,
		       to_char(s.practice_date, 'YYYY-MM-DD'),
		       to_char(s.start_time, 'HH24:MI'),
		       to_char(s.end_time, 'HH24:MI'),
		       s.quota,
		       COALESCE(COUNT(a.id) FILTER (WHERE a.status <> 'cancelled'), 0) as booked_count,
		       s.quota - COALESCE(COUNT(a.id) FILTER (WHERE a.status <> 'cancelled'), 0) as remaining_quota,
		       s.created_at, s.updated_at
		FROM doctor_schedules s
		JOIN doctors d ON d.id = s.doctor_id
		JOIN users u ON u.id = d.user_id
		JOIN polyclinics p ON p.id = d.polyclinic_id
		LEFT JOIN appointments a ON a.schedule_id = s.id
	`
	var conditions []string
	var args []any
	argIndex := 1

	if doctorID != nil {
		conditions = append(conditions, "s.doctor_id = $"+strconv.Itoa(argIndex))
		args = append(args, *doctorID)
		argIndex++
	}
	if polyclinicID != nil {
		conditions = append(conditions, "d.polyclinic_id = $"+strconv.Itoa(argIndex))
		args = append(args, *polyclinicID)
		argIndex++
	}
	if date != "" {
		conditions = append(conditions, "s.practice_date = $"+strconv.Itoa(argIndex)+"::date")
		args = append(args, date)
		argIndex++
	}
	if from != "" {
		conditions = append(conditions, "s.practice_date >= $"+strconv.Itoa(argIndex)+"::date")
		args = append(args, from)
		argIndex++
	}
	if to != "" {
		conditions = append(conditions, "s.practice_date <= $"+strconv.Itoa(argIndex)+"::date")
		args = append(args, to)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY s.id, u.full_name, d.polyclinic_id, p.name
		ORDER BY s.practice_date ASC, s.start_time ASC
	`

	rows, err := r.getDB(db).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list schedules: %w", err)
	}
	defer rows.Close()

	var result []model.DoctorSchedule
	for rows.Next() {
		var s model.DoctorSchedule
		if err := rows.Scan(
			&s.ID, &s.DoctorID, &s.DoctorName, &s.PolyclinicID, &s.PolyclinicName,
			&s.PracticeDate, &s.StartTime, &s.EndTime, &s.Quota,
			&s.BookedCount, &s.RemainingQuota, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}
		result = append(result, s)
	}
	if result == nil {
		result = []model.DoctorSchedule{}
	}
	return result, nil
}

func (r *clinicPostgresRepository) UpdateSchedule(ctx context.Context, db DBTX, s *model.DoctorSchedule) error {
	query := `
		UPDATE doctor_schedules
		SET start_time = $1::time, end_time = $2::time, quota = $3, updated_at = now()
		WHERE id = $4
		RETURNING to_char(start_time, 'HH24:MI'), to_char(end_time, 'HH24:MI'), updated_at
	`
	return r.getDB(db).QueryRow(ctx, query, s.StartTime, s.EndTime, s.Quota, s.ID).
		Scan(&s.StartTime, &s.EndTime, &s.UpdatedAt)
}

func (r *clinicPostgresRepository) CountActiveBookingsBySchedule(ctx context.Context, db DBTX, scheduleID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM appointments
		WHERE schedule_id = $1 AND status <> 'cancelled'
	`
	var count int
	err := r.getDB(db).QueryRow(ctx, query, scheduleID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count active bookings: %w", err)
	}
	return count, nil
}
