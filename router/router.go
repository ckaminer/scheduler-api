package router

import (
	"time"

	"github.com/ckaminer/schedule-api/scheduler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func InitializeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/schedules", scheduler.CreateScheduleHandler)
	r.Get("/schedules/{scheduleID}", scheduler.ScheduleDetailsHandler)
	r.Delete("/schedules/{scheduleID}", scheduler.DeleteScheduleHandler)

	r.Post("/schedules/{scheduleID}/appointments", scheduler.CreateAppointmentHandler)
	r.Get("/schedules/{scheduleID}/appointments/{appointmentID}", scheduler.AppointmentDetailsHandler)
	r.Delete("/schedules/{scheduleID}/appointments/{appointmentID}", scheduler.DeleteAppointmentHandler)
	return r
}
