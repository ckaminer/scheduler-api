package scheduler

import (
	"log"

	"github.com/ckaminer/go-utils/http_helpers"
)

type SchedulerService struct{}

func (service *SchedulerService) CreateAppointment(a Appointment, scheduleID int) (Appointment, error) {
	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("ScheduleDetailsService - no schedule found for ID: ", scheduleID)
		return a, http_helpers.NotFoundError{EntityType: "Schedule"}
	}

	a.ScheduleID = s.ID

	a.ID = AppointmentsCreatedCount + 1
	s.Appointments[a.ID] = a
	AppointmentsCreatedCount++

	return a, nil
}
