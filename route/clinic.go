package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jamaalul/proyek1pbl/app/service"
	"github.com/jamaalul/proyek1pbl/helper"
	"github.com/jamaalul/proyek1pbl/middleware"
)

func RegisterClinicRoutes(r fiber.Router, svc *service.ClinicService, jwtMgr *helper.JWTManager, perms *helper.PermissionSet) {
	api := r.Group("", middleware.RequireAuth(jwtMgr))

	// Polyclinics
	api.Get("/polyclinics", svc.ListPolyclinics)
	api.Post("/polyclinics", middleware.RequirePermission(perms, "clinic:manage"), svc.CreatePolyclinic)

	// Doctors
	api.Get("/doctors", svc.ListDoctors)
	api.Post("/doctors", middleware.RequirePermission(perms, "clinic:manage"), svc.CreateDoctor)

	// Schedules
	api.Get("/schedules", svc.ListSchedules)
	api.Post("/schedules", middleware.RequirePermission(perms, "clinic:manage"), svc.CreateSchedule)
	api.Put("/schedules/:id", middleware.RequirePermission(perms, "clinic:manage"), svc.UpdateSchedule)
}
