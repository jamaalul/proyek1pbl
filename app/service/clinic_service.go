package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jamaalul/proyek1pbl/app/model"
	"github.com/jamaalul/proyek1pbl/app/repository"
	"github.com/jamaalul/proyek1pbl/helper"
)

type ClinicService struct {
	clinicRepo repository.ClinicRepository
	txMgr      *repository.TxManager
}

func NewClinicService(clinicRepo repository.ClinicRepository, txMgr *repository.TxManager) *ClinicService {
	return &ClinicService{
		clinicRepo: clinicRepo,
		txMgr:      txMgr,
	}
}

// --- Polyclinic Handlers ---

func (s *ClinicService) CreatePolyclinic(c *fiber.Ctx) error {
	var req model.CreatePolyclinicRequest
	if err := helper.ParseAndValidate(c, &req); err != nil {
		return err
	}

	poly := &model.Polyclinic{
		Name:        req.Name,
		Description: req.Description,
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	if err := s.clinicRepo.CreatePolyclinic(ctx, nil, poly); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return helper.ErrConflict("CONFLICT", "Nama poliklinik sudah ada")
		}
		return helper.ErrInternal(err)
	}

	location := fmt.Sprintf("/api/v1/polyclinics/%d", poly.ID)
	return helper.Created(c, "Poliklinik berhasil dibuat", poly, location)
}

func (s *ClinicService) ListPolyclinics(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	list, err := s.clinicRepo.ListPolyclinics(ctx, nil)
	if err != nil {
		return helper.ErrInternal(err)
	}

	return helper.Success(c, fiber.StatusOK, "Daftar poliklinik berhasil diambil", list)
}

// --- Doctor Handlers ---

func (s *ClinicService) CreateDoctor(c *fiber.Ctx) error {
	var req model.CreateDoctorRequest
	if err := helper.ParseAndValidate(c, &req); err != nil {
		return err
	}

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.ErrInternal(err)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	doctor := &model.Doctor{
		FullName:       req.FullName,
		Email:          req.Email,
		PolyclinicID:   req.PolyclinicID,
		Specialization: req.Specialization,
	}

	// Buat user dan profile dokter dalam satu transaksi atomik (PRD §9.6)
	err = s.txMgr.WithTx(ctx, func(tx repository.DBTX) error {
		userID, userErr := s.clinicRepo.CreateDoctorUser(ctx, tx, req.FullName, req.Email, hash)
		if userErr != nil {
			return userErr
		}

		doctor.UserID = userID
		return s.clinicRepo.CreateDoctor(ctx, tx, doctor)
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return helper.ErrConflict("EMAIL_TAKEN", "Email sudah terdaftar")
			}
			if pgErr.Code == "23503" {
				return helper.ErrNotFound("Poliklinik tidak ditemukan")
			}
		}
		return helper.ErrInternal(err)
	}

	location := fmt.Sprintf("/api/v1/doctors/%d", doctor.ID)
	return helper.Created(c, "Dokter berhasil dibuat", doctor, location)
}

func (s *ClinicService) ListDoctors(c *fiber.Ctx) error {
	var polyclinicID *int64
	if polyIDStr := c.Query("polyclinic_id"); polyIDStr != "" {
		if id, err := strconv.ParseInt(polyIDStr, 10, 64); err == nil && id > 0 {
			polyclinicID = &id
		}
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	list, err := s.clinicRepo.ListDoctors(ctx, nil, polyclinicID)
	if err != nil {
		return helper.ErrInternal(err)
	}

	return helper.Success(c, fiber.StatusOK, "Daftar dokter berhasil diambil", list)
}

// --- Schedule Handlers ---

func (s *ClinicService) CreateSchedule(c *fiber.Ctx) error {
	var req model.CreateScheduleRequest
	if err := helper.ParseAndValidate(c, &req); err != nil {
		return err
	}

	if req.Quota <= 0 {
		return helper.ErrUnprocessable("VALIDATION_ERROR", "Kuota harus lebih besar dari 0")
	}

	if req.EndTime <= req.StartTime {
		return helper.ErrUnprocessable("VALIDATION_ERROR", "Jam selesai harus setelah jam mulai")
	}

	// Validasi tanggal tidak boleh di masa lalu berdasarkan waktu Asia/Jakarta (PRD Decision #14 & FR-MASTER-3)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("Asia/Jakarta", 7*3600)
	}
	todayJakarta := time.Now().In(loc).Format("2006-01-02")
	if req.PracticeDate < todayJakarta {
		return helper.ErrUnprocessable("VALIDATION_ERROR", "Tanggal praktik tidak boleh di masa lalu")
	}

	schedule := &model.DoctorSchedule{
		DoctorID:     req.DoctorID,
		PracticeDate: req.PracticeDate,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Quota:        req.Quota,
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	if err := s.clinicRepo.CreateSchedule(ctx, nil, schedule); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23P01" || pgErr.Code == "23505" {
				return helper.ErrConflict("SCHEDULE_OVERLAP", "Jadwal dokter bentrok dengan jadwal yang sudah ada")
			}
			if pgErr.Code == "23503" {
				return helper.ErrNotFound("Dokter tidak ditemukan")
			}
			if pgErr.Code == "23514" {
				return helper.ErrUnprocessable("VALIDATION_ERROR", "Validasi aturan jadwal gagal")
			}
		}
		return helper.ErrInternal(err)
	}

	location := fmt.Sprintf("/api/v1/schedules/%d", schedule.ID)
	return helper.Created(c, "Jadwal dokter berhasil dibuat", schedule, location)
}

func (s *ClinicService) ListSchedules(c *fiber.Ctx) error {
	var doctorID *int64
	if docStr := c.Query("doctor_id"); docStr != "" {
		if id, err := strconv.ParseInt(docStr, 10, 64); err == nil && id > 0 {
			doctorID = &id
		}
	}

	var polyclinicID *int64
	if polyStr := c.Query("polyclinic_id"); polyStr != "" {
		if id, err := strconv.ParseInt(polyStr, 10, 64); err == nil && id > 0 {
			polyclinicID = &id
		}
	}

	date := c.Query("date")
	from := c.Query("from")
	to := c.Query("to")

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	list, err := s.clinicRepo.ListSchedules(ctx, nil, doctorID, polyclinicID, date, from, to)
	if err != nil {
		return helper.ErrInternal(err)
	}

	return helper.Success(c, fiber.StatusOK, "Daftar jadwal dokter berhasil diambil", list)
}

func (s *ClinicService) UpdateSchedule(c *fiber.Ctx) error {
	id, ok := helper.ParamID64(c)
	if !ok {
		return helper.ErrBadRequest("BAD_REQUEST", "ID jadwal tidak valid")
	}

	var req model.UpdateScheduleRequest
	if err := helper.ParseAndValidate(c, &req); err != nil {
		return err
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	existing, err := s.clinicRepo.GetScheduleByID(ctx, nil, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return helper.ErrNotFound("Jadwal tidak ditemukan")
		}
		return helper.ErrInternal(err)
	}

	// Validasi penurunan kuota tidak boleh di bawah antrean aktif (PRD §6.3 & FR-MASTER-3)
	if req.Quota != nil {
		if *req.Quota <= 0 {
			return helper.ErrUnprocessable("VALIDATION_ERROR", "Kuota harus lebih besar dari 0")
		}

		activeBookings, countErr := s.clinicRepo.CountActiveBookingsBySchedule(ctx, nil, id)
		if countErr != nil {
			return helper.ErrInternal(countErr)
		}

		if *req.Quota < activeBookings {
			return helper.ErrConflict("QUOTA_BELOW_BOOKED", "Kuota tidak boleh lebih kecil dari jumlah reservasi aktif")
		}
		existing.Quota = *req.Quota
	}

	if req.StartTime != nil {
		existing.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		existing.EndTime = *req.EndTime
	}

	if existing.EndTime <= existing.StartTime {
		return helper.ErrUnprocessable("VALIDATION_ERROR", "Jam selesai harus setelah jam mulai")
	}

	if err := s.clinicRepo.UpdateSchedule(ctx, nil, existing); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23P01" || pgErr.Code == "23505" {
				return helper.ErrConflict("SCHEDULE_OVERLAP", "Jadwal dokter bentrok dengan jadwal yang sudah ada")
			}
			if pgErr.Code == "23514" {
				return helper.ErrUnprocessable("VALIDATION_ERROR", "Validasi aturan jadwal gagal")
			}
		}
		return helper.ErrInternal(err)
	}

	return helper.Success(c, fiber.StatusOK, "Jadwal dokter berhasil diperbarui", existing)
}
