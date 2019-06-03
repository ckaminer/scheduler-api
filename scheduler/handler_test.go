package scheduler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/ckaminer/schedule-api/scheduler"
)

var _ = Describe("Handlers", func() {
	Context("Schedule Handlers", func() {
		Context("#CreateSchedule", func() {
			It("Should return a StatusCreated and the created entity upon successful creation of a schedule", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateScheduleHandler)

				reqBody := []byte(`
					{
						"owner_name": "Tyrion Lannister"
					}
				`)

				r, _ := http.NewRequest("POST", "/schedules", bytes.NewReader(reqBody))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusCreated))

				var resBody ScheduleResponse
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody.OwnerName).To(Equal("Tyrion Lannister"))
				Expect(resBody.ID).To(Equal(1))
				Expect(resBody.Appointments).To(Equal([]Appointment{}))
			})

			It("Should increment the ID by one for each created schedule", func() {
				SchedulesCreatedCount = 4
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateScheduleHandler)

				reqBody := []byte(`
					{
						"owner_name": "Tyrion Lannister"
					}
				`)

				r, _ := http.NewRequest("POST", "/schedules", bytes.NewReader(reqBody))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusCreated))

				var resBody ScheduleResponse
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody.OwnerName).To(Equal("Tyrion Lannister"))
				Expect(resBody.ID).To(Equal(5))

				Expect(SchedulesCreatedCount).To(Equal(5))
			})

			It("Should return a StatusBadRequest if the reqBody is invalid", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateScheduleHandler)

				reqBody := []byte(`{"owner_name": true}`)

				r, _ := http.NewRequest("POST", "/schedules", bytes.NewReader(reqBody))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("#ScheduleDetails", func() {
			It("Should return schedule details for the given scheduleID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(ScheduleDetailsHandler)

				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
				}
				ScheduleCollection[s.ID] = s

				r, _ := http.NewRequest("GET", "/schedules/31", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "31")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				var resBody ScheduleResponse
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody.ID).To(Equal(s.ID))
				Expect(resBody.OwnerName).To(Equal(s.OwnerName))
				Expect(resBody.Appointments).To(Equal([]Appointment{}))
			})

			It("Should return a StatusBadRequest for a non-numerical schedule ID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(ScheduleDetailsHandler)

				r, _ := http.NewRequest("GET", "/schedules/blamo", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "blamo")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusNotFound for a scheduleID that does not have an associated schedule", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(ScheduleDetailsHandler)

				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
				}
				ScheduleCollection[s.ID] = s

				r, _ := http.NewRequest("GET", "/schedules/32", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "32")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("#DeleteSchedule", func() {
			It("Should return a 200 and the deleted entity upon success", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteScheduleHandler)

				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
				}
				ScheduleCollection[s.ID] = s

				r, _ := http.NewRequest("DELETE", "/schedules/31", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "31")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				var resBody ScheduleResponse
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody.ID).To(Equal(s.ID))
				Expect(resBody.OwnerName).To(Equal(s.OwnerName))
				Expect(resBody.Appointments).To(Equal([]Appointment{}))

				if _, found := ScheduleCollection[s.ID]; found {
					Fail("Schedule should be deleted from storage")
				}
			})

			It("Should return a StatusBadRequest for a non-numerical schedule ID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteScheduleHandler)

				r, _ := http.NewRequest("DELETE", "/schedules/blamo", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "blamo")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusNotFound for a scheduleID that does not have an associated schedule", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteScheduleHandler)

				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
				}
				ScheduleCollection[s.ID] = s

				r, _ := http.NewRequest("DELETE", "/schedules/32", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "32")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Context("Appointment Handlers", func() {
		Context("#CreateAppointment", func() {
			It("Should return a StatusCreated and the created entity upon successful creation of a appointment", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateAppointmentHandler)

				s := Schedule{
					ID:           12,
					OwnerName:    "Tyrion Lannister",
					Appointments: make(map[int]Appointment),
				}
				ScheduleCollection[s.ID] = s

				a := Appointment{
					StartTime: 5,
					EndTime:   9,
				}
				reqBody, _ := json.Marshal(a)

				r, _ := http.NewRequest("POST", "/schedules/12/appointments", bytes.NewReader(reqBody))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "12")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusCreated))

				var resBody Appointment
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody.StartTime).To(Equal(a.StartTime))
				Expect(resBody.EndTime).To(Equal(a.EndTime))
				Expect(resBody.ScheduleID).To(Equal(s.ID))
				Expect(resBody.ID).To(Equal(1))

				scheduleAppts := ScheduleCollection[s.ID].Appointments
				Expect(len(scheduleAppts)).To(Equal(1))
				Expect(scheduleAppts[resBody.ID]).To(Equal(resBody))
			})

			It("Should increment the ID by one for each created appointment", func() {
				AppointmentsCreatedCount = 7
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateAppointmentHandler)

				s := Schedule{
					ID:           12,
					OwnerName:    "Tyrion Lannister",
					Appointments: make(map[int]Appointment),
				}
				ScheduleCollection[s.ID] = s

				a := Appointment{
					StartTime: 5,
					EndTime:   9,
				}
				reqBody, _ := json.Marshal(a)

				r, _ := http.NewRequest("POST", "/schedules/12/appointments", bytes.NewReader(reqBody))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "12")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusCreated))

				var resBody Appointment
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody.StartTime).To(Equal(a.StartTime))
				Expect(resBody.EndTime).To(Equal(a.EndTime))
				Expect(resBody.ScheduleID).To(Equal(s.ID))
				Expect(resBody.ID).To(Equal(8))

				Expect(AppointmentsCreatedCount).To(Equal(8))
			})

			It("Should return a StatusBadRequest for a non-numerical schedule ID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateAppointmentHandler)

				r, _ := http.NewRequest("POST", "/schedules/blamo/appointments", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "blamo")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusBadRequest if the reqBody is invalid", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateAppointmentHandler)

				reqBody := []byte(`{"start_time": true}`)

				r, _ := http.NewRequest("POST", "/schedules/1/appointments", bytes.NewReader(reqBody))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "1")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusNotFound for a scheduleID that does not have an associated schedule", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateAppointmentHandler)

				a := Appointment{
					StartTime: 5,
					EndTime:   9,
				}
				reqBody, _ := json.Marshal(a)

				r, _ := http.NewRequest("POST", "/schedules/32/appointments", bytes.NewReader(reqBody))

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "32")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})

			It("Should return a StatusUnprocessableEntity for invalid appointment times", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(CreateAppointmentHandler)

				scheduledAppts := map[int]Appointment{
					5: Appointment{
						StartTime: 7,
						EndTime:   15,
					},
				}
				s := Schedule{
					ID:           12,
					OwnerName:    "Tyrion Lannister",
					Appointments: scheduledAppts,
				}
				ScheduleCollection[s.ID] = s

				a := Appointment{
					StartTime: 5,
					EndTime:   9,
				}
				reqBody, _ := json.Marshal(a)

				r, _ := http.NewRequest("POST", "/schedules/12/appointments", bytes.NewReader(reqBody))
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "12")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			})
		})

		Context("#AppointmentDetails", func() {
			It("Should return appointment details for the given appointmentID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(AppointmentDetailsHandler)

				a := Appointment{
					ID:         12,
					ScheduleID: 29,
					StartTime:  5,
					EndTime:    90,
				}
				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
					Appointments: map[int]Appointment{
						a.ID: a,
					},
				}
				ScheduleCollection[s.ID] = s

				r, _ := http.NewRequest("GET", "/schedules/31/appointments/12", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("appointmentID", "12")
				rctx.URLParams.Add("scheduleID", "31")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				var resBody Appointment
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody).To(Equal(a))
			})

			It("Should return a StatusBadRequest for a non-numerical schedule ID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(AppointmentDetailsHandler)

				r, _ := http.NewRequest("GET", "schedules/blamo/appointments/12", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("appointmentID", "12")
				rctx.URLParams.Add("scheduleID", "blamo")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusBadRequest for a non-numerical appointment ID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(AppointmentDetailsHandler)

				r, _ := http.NewRequest("GET", "schedules/13/appointments/blamo", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("appointmentID", "blamo")
				rctx.URLParams.Add("scheduleID", "13")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusNotFound if the schedule is not found", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(AppointmentDetailsHandler)

				r, _ := http.NewRequest("GET", "/schedule/12/appointments/13", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "13")
				rctx.URLParams.Add("appointmentID", "13")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})

			It("Should return a StatusNotFound if the appointment is not found on the provided schedule", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(AppointmentDetailsHandler)

				a := Appointment{
					ID:         12,
					ScheduleID: 29,
					StartTime:  5,
					EndTime:    90,
				}
				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
					Appointments: map[int]Appointment{
						a.ID: a,
					},
				}
				ScheduleCollection[s.ID] = s

				r, _ := http.NewRequest("GET", "/schedule/31/appointments/12", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "31")
				rctx.URLParams.Add("appointmentID", "13")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("#DeleteAppointment", func() {
			It("Should return a 200 and the deleted entity upon success", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteAppointmentHandler)

				a := Appointment{
					ID:         12,
					ScheduleID: 29,
					StartTime:  5,
					EndTime:    90,
				}
				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
					Appointments: map[int]Appointment{
						a.ID: a,
					},
				}
				ScheduleCollection[s.ID] = s
				r, _ := http.NewRequest("DELETE", "/schedules/31/appointments/12", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "31")
				rctx.URLParams.Add("appointmentID", "12")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusOK))

				var resBody Appointment
				err := json.NewDecoder(recorder.Body).Decode(&resBody)
				if err != nil {
					Fail("Unable to decode response body")
				}

				Expect(resBody).To(Equal(a))

				if foundSchedule, found := ScheduleCollection[s.ID]; found {
					if _, foundAppt := foundSchedule.Appointments[a.ID]; foundAppt {
						Fail("Schedule should be deleted from storage")
					}
				} else {
					Fail("Should have found schedule")
				}
			})

			It("Should return a StatusBadRequest for a non-numerical schedule ID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteAppointmentHandler)

				r, _ := http.NewRequest("DELETE", "schedules/blamo/appointments/12", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("appointmentID", "12")
				rctx.URLParams.Add("scheduleID", "blamo")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusBadRequest for a non-numerical appointment ID", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteAppointmentHandler)

				r, _ := http.NewRequest("DELETE", "schedules/13/appointments/blamo", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("appointmentID", "blamo")
				rctx.URLParams.Add("scheduleID", "13")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("Should return a StatusNotFound if the schedule is not found", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteAppointmentHandler)

				r, _ := http.NewRequest("DELETE", "/schedule/12/appointments/13", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "13")
				rctx.URLParams.Add("appointmentID", "13")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})

			It("Should return a StatusNotFound if the appointment is not found on the provided schedule", func() {
				recorder := httptest.NewRecorder()
				handler := http.HandlerFunc(DeleteAppointmentHandler)

				a := Appointment{
					ID:         12,
					ScheduleID: 29,
					StartTime:  5,
					EndTime:    90,
				}
				s := Schedule{
					ID:        31,
					OwnerName: "Tyrion Lannister",
					Appointments: map[int]Appointment{
						a.ID: a,
					},
				}
				ScheduleCollection[s.ID] = s

				r, _ := http.NewRequest("DELETE", "/schedule/31/appointments/12", nil)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("scheduleID", "31")
				rctx.URLParams.Add("appointmentID", "13")
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				handler.ServeHTTP(recorder, r)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
