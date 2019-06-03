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
	Context("#CreateSchedule", func() {
		It("Should return a StatusCreated and the created entity upon successful creation of a schedule", func() {
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.CreateSchedule)

			s := Schedule{
				OwnerName: "Tyrion Lannister",
			}
			reqBody, _ := json.Marshal(s)

			r, _ := http.NewRequest("POST", "/schedules", bytes.NewReader(reqBody))

			handler.ServeHTTP(recorder, r)

			Expect(recorder.Code).To(Equal(http.StatusCreated))

			var resBody Schedule
			err := json.NewDecoder(recorder.Body).Decode(&resBody)
			if err != nil {
				Fail("Unable to decode response body")
			}

			Expect(resBody.OwnerName).To(Equal(s.OwnerName))
			Expect(resBody.ID).To(Equal(1))
		})

		It("Should increment the ID by one for each created schedule", func() {
			SchedulesCreatedCount = 4
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.CreateSchedule)

			s := Schedule{
				OwnerName: "Tyrion Lannister",
			}
			reqBody, _ := json.Marshal(s)

			r, _ := http.NewRequest("POST", "/schedules", bytes.NewReader(reqBody))

			handler.ServeHTTP(recorder, r)

			Expect(recorder.Code).To(Equal(http.StatusCreated))

			var resBody Schedule
			err := json.NewDecoder(recorder.Body).Decode(&resBody)
			if err != nil {
				Fail("Unable to decode response body")
			}

			Expect(resBody.OwnerName).To(Equal(s.OwnerName))
			Expect(resBody.ID).To(Equal(5))
		})

		It("Should return a StatusBadRequest if the reqBody is invalid", func() {
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.CreateSchedule)

			reqBody := []byte(`{"owner_name": true}`)

			r, _ := http.NewRequest("POST", "/schedules", bytes.NewReader(reqBody))

			handler.ServeHTTP(recorder, r)

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("#ScheduleDetails", func() {
		It("Should return schedule details for the given scheduleID", func() {
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.ScheduleDetails)

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

			var resBody Schedule
			err := json.NewDecoder(recorder.Body).Decode(&resBody)
			if err != nil {
				Fail("Unable to decode response body")
			}

			Expect(resBody).To(Equal(s))
		})

		It("Should return a StatusBadRequest for a non-numerical schedule ID", func() {
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.ScheduleDetails)

			r, _ := http.NewRequest("GET", "/schedules/blamo", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("scheduleID", "blamo")
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler.ServeHTTP(recorder, r)

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})

		It("Should return a StatusNotFound for a scheduleID that does not have an associated schedule", func() {
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.ScheduleDetails)

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
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.DeleteSchedule)

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

			var resBody Schedule
			err := json.NewDecoder(recorder.Body).Decode(&resBody)
			if err != nil {
				Fail("Unable to decode response body")
			}

			Expect(resBody).To(Equal(s))

			if _, found := ScheduleCollection[s.ID]; found {
				Fail("Schedule should be deleted from storage")
			}
		})

		It("Should return a StatusBadRequest for a non-numerical schedule ID", func() {
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.DeleteSchedule)

			r, _ := http.NewRequest("DELETE", "/schedules/blamo", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("scheduleID", "blamo")
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler.ServeHTTP(recorder, r)

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})

		It("Should return a StatusNotFound for a scheduleID that does not have an associated schedule", func() {
			recorder := httptest.NewRecorder()
			wrapper := SchedulerHandler{}
			handler := http.HandlerFunc(wrapper.DeleteSchedule)

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
