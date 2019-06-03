package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ckaminer/go-utils/http_helpers"
	"github.com/go-chi/chi"
)

type SchedulerHandler struct {
	Service SchedulerService
}

func (handler *SchedulerHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var s Schedule
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println("CreateScheduleHandler Err: ", err.Error())
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}
	defer r.Body.Close()

	s.ID = SchedulesCreatedCount + 1
	s.Appointments = make(map[int]Appointment)
	ScheduleCollection[s.ID] = s
	SchedulesCreatedCount++

	http_helpers.RespondWithJSON(w, http.StatusCreated, s)
}

func (handler *SchedulerHandler) ScheduleDetails(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := convertIDParam(r, "scheduleID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("ScheduleDetailsHandler - no schedule found for ID: ", scheduleID)
		http_helpers.RespondWithParsedError(w, http_helpers.NotFoundError{EntityType: "Schedule"})
		return
	}

	http_helpers.RespondWithJSON(w, http.StatusOK, s)
}

func (handler *SchedulerHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := convertIDParam(r, "scheduleID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("DeleteScheduleHandler - no schedule found for ID: ", scheduleID)
		http_helpers.RespondWithParsedError(w, http_helpers.NotFoundError{EntityType: "Schedule"})
		return
	}

	delete(ScheduleCollection, scheduleID)
	http_helpers.RespondWithJSON(w, http.StatusOK, s)
}

func (handler *SchedulerHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := convertIDParam(r, "scheduleID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var a Appointment
	err = json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		log.Println("CreateAppointmentHandler Err: ", err.Error())
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}
	defer r.Body.Close()

	createdAppt, err := handler.Service.CreateAppointment(a, scheduleID)
	if err != nil {
		http_helpers.RespondWithParsedError(w, err)
		return
	}

	http_helpers.RespondWithJSON(w, http.StatusCreated, createdAppt)
}

func (handler *SchedulerHandler) AppointmentDetails(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := convertIDParam(r, "scheduleID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	appointmentID, err := convertIDParam(r, "appointmentID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("AppointmentDetailsHandler - no schedule found for ID: ", scheduleID)
		http_helpers.RespondWithParsedError(w, http_helpers.NotFoundError{EntityType: "Schedule"})
		return
	}

	var a Appointment
	a, found = s.Appointments[appointmentID]
	if !found {
		log.Println("AppointmentDetailsHandler - no appointment found for ID: ", appointmentID)
		http_helpers.RespondWithParsedError(w, http_helpers.NotFoundError{EntityType: "Appointment"})
		return
	}

	http_helpers.RespondWithJSON(w, http.StatusOK, a)
}

func (handler *SchedulerHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := convertIDParam(r, "scheduleID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	appointmentID, err := convertIDParam(r, "appointmentID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("DeleteScheduleHandler - no schedule found for ID: ", scheduleID)
		http_helpers.RespondWithParsedError(w, http_helpers.NotFoundError{EntityType: "Schedule"})
		return
	}

	var a Appointment
	a, found = s.Appointments[appointmentID]
	if !found {
		log.Println("DeleteAppointmentHandler - no appointment found for ID: ", appointmentID)
		http_helpers.RespondWithParsedError(w, http_helpers.NotFoundError{EntityType: "Appointment"})
		return
	}

	delete(s.Appointments, appointmentID)
	http_helpers.RespondWithJSON(w, http.StatusOK, a)
}

func convertIDParam(r *http.Request, paramName string) (int, error) {
	idParam := chi.URLParam(r, paramName)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		errorMessage := fmt.Sprintf("Scheduler Handler - Failed to convert ID param (%v) due to: %v", paramName, err.Error())
		log.Println(errorMessage)
	}
	return id, err
}
