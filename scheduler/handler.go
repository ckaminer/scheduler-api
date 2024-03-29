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

func CreateScheduleHandler(w http.ResponseWriter, r *http.Request) {
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

func ScheduleDetailsHandler(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := convertIDParam(r, "scheduleID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("ScheduleDetailsHandler - no schedule found for ID: ", scheduleID)
		http_helpers.RespondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	http_helpers.RespondWithJSON(w, http.StatusOK, s)
}

func DeleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := convertIDParam(r, "scheduleID")
	if err != nil {
		http_helpers.RespondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var s Schedule
	s, found := ScheduleCollection[scheduleID]
	if !found {
		log.Println("DeleteScheduleHandler - no schedule found for ID: ", scheduleID)
		http_helpers.RespondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	delete(ScheduleCollection, scheduleID)
	http_helpers.RespondWithJSON(w, http.StatusOK, s)
}

func CreateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
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

	createdAppt, err := createAppointment(a, scheduleID)
	if err != nil {
		if httpErr, ok := err.(http_helpers.HttpError); ok {
			http_helpers.RespondWithError(w, httpErr.StatusCode, httpErr.Message)
		} else {
			http_helpers.RespondWithError(w, http.StatusServiceUnavailable, "Unable to create appointment")
		}
		return
	}

	http_helpers.RespondWithJSON(w, http.StatusCreated, createdAppt)
}

func AppointmentDetailsHandler(w http.ResponseWriter, r *http.Request) {
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
		http_helpers.RespondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	var a Appointment
	a, found = s.Appointments[appointmentID]
	if !found {
		log.Println("AppointmentDetailsHandler - no appointment found for ID: ", appointmentID)
		http_helpers.RespondWithError(w, http.StatusNotFound, "Appointment not found")
		return
	}

	http_helpers.RespondWithJSON(w, http.StatusOK, a)
}

func DeleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
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
		http_helpers.RespondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	var a Appointment
	a, found = s.Appointments[appointmentID]
	if !found {
		log.Println("DeleteAppointmentHandler - no appointment found for ID: ", appointmentID)
		http_helpers.RespondWithError(w, http.StatusNotFound, "Appointment not found")
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

type ScheduleResponse struct {
	ID           int           `json:"id"`
	OwnerName    string        `json:"owner_name"`
	Appointments []Appointment `json:"appointments"`
}

func (s Schedule) MarshalJSON() ([]byte, error) {
	sortedAppointments := sortAppointments(s)
	scheduleRes := ScheduleResponse{
		ID:           s.ID,
		OwnerName:    s.OwnerName,
		Appointments: sortedAppointments,
	}

	return json.Marshal(scheduleRes)
}
