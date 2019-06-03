package scheduler

import (
	"log"
	"net/http"
	"sort"

	"github.com/ckaminer/go-utils/http_helpers"
)

func createAppointment(a Appointment, scheduleID int) (Appointment, error) {
	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("ScheduleDetailsService - no schedule found for ID: ", scheduleID)
		return a, http_helpers.HttpError{
			StatusCode: http.StatusNotFound,
			Message:    "Schedule not found",
		}
	}

	validAppt := ValidateAppointmentInput(s, a)
	if !validAppt {
		return a, http_helpers.HttpError{
			Message:    "Invalid appointment time",
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	a.ScheduleID = s.ID

	a.ID = AppointmentsCreatedCount + 1
	s.Appointments[a.ID] = a
	AppointmentsCreatedCount++

	return a, nil
}

func ValidateAppointmentInput(s Schedule, a Appointment) bool {
	if a.StartTime >= a.EndTime || a.StartTime == 0 {
		return false
	}

	for _, scheduledAppt := range s.Appointments {
		outsideRange := a.StartTime > scheduledAppt.EndTime || a.EndTime < scheduledAppt.StartTime
		if !outsideRange {
			return false
		}
	}

	return true
}

func sortAppointments(s Schedule) (appointments []Appointment) {
	if len(s.Appointments) == 0 {
		appointments = []Appointment{}
	}
	for _, appt := range s.Appointments {
		appointments = append(appointments, appt)
	}

	sort.SliceStable(appointments, func(i, j int) bool {
		return appointments[i].StartTime < appointments[j].StartTime
	})

	return
}
