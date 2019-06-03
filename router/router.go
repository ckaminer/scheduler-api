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

	schedulerHandler := scheduler.SchedulerHandler{
		Service: scheduler.SchedulerService{},
	}
	r.Post("/schedules", schedulerHandler.CreateSchedule)
	r.Get("/schedules/{scheduleID}", schedulerHandler.ScheduleDetails)
	r.Delete("/schedules/{scheduleID}", schedulerHandler.DeleteSchedule)

	r.Post("/schedules/{scheduleID}/appointments", schedulerHandler.CreateAppointment)
	r.Get("/schedules/{scheduleID}/appointments/{appointmentID}", schedulerHandler.AppointmentDetails)
	r.Delete("/schedules/{scheduleID}/appointments/{appointmentID}", schedulerHandler.DeleteAppointment)
	return r
}
