package acceptance_tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ckaminer/schedule-api/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schedule Routes", func() {
	Context("POST /schedules", func() {
		It("Should return a created schedule with a generated ID upon success", func() {
			reqBody := []byte(`
				{
					"owner_name": "Tyrion Lannister"
				}
			`)

			url := fmt.Sprintf("%v/schedules", acceptanceUrl)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var s scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&s)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))
			Expect(s.OwnerName).To(Equal("Tyrion Lannister"))
			Expect(s.ID).To(BeNumerically(">", 0))
			Expect(s.Appointments).To(Equal([]scheduler.Appointment{}))

			// Send additional request to confirm ID is being incremented
			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var s2 scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&s2)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))
			Expect(s2.ID).To(Equal(s.ID + 1))
		})

		It("Should return a bad request if the request body is invalid", func() {
			reqBody := []byte(`
					{
						"owner_name": true
					}
				`)

			url := fmt.Sprintf("%v/schedules", acceptanceUrl)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	Context("GET /schedules/{scheduleID}", func() {
		It("Should return a schedule for the associated schedule ID provided in the req path", func() {
			// Create schedule
			reqBody := []byte(`
				{
					"owner_name": "Tyrion Lannister"
				}
			`)

			url := fmt.Sprintf("%v/schedules", acceptanceUrl)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var createdSchedule scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&createdSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			// Find the schedule that was just created
			url = fmt.Sprintf("%v/schedules/%v", acceptanceUrl, createdSchedule.ID)
			req, _ = http.NewRequest("GET", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var foundSchedule scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&foundSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(foundSchedule).To(Equal(createdSchedule))
		})

		It("Should return appointments sorted by start time", func() {
			// Create schedule
			reqBody := []byte(`
				{
					"owner_name": "Tyrion Lannister"
				}
			`)

			url := fmt.Sprintf("%v/schedules", acceptanceUrl)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var createdSchedule scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&createdSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			// Create appointments
			appointments := []scheduler.Appointment{
				scheduler.Appointment{
					StartTime: 10,
					EndTime:   12,
				},
				scheduler.Appointment{
					StartTime: 1,
					EndTime:   3,
				},
				scheduler.Appointment{
					StartTime: 5,
					EndTime:   9,
				},
			}

			for _, a := range appointments {
				url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, createdSchedule.ID)
				reqBody, _ := json.Marshal(a)
				req, _ = http.NewRequest("POST", url, bytes.NewReader(reqBody))

				res, err = client.Do(req)
				if err != nil {
					Fail("Failed to send request")
				}

				Expect(res.StatusCode).To(Equal(http.StatusCreated))
			}

			// Find the schedule that was just created
			url = fmt.Sprintf("%v/schedules/%v", acceptanceUrl, createdSchedule.ID)
			req, _ = http.NewRequest("GET", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var foundSchedule scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&foundSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusOK))
			// Check for start times from created appointments
			Expect(foundSchedule.Appointments[0].StartTime).To(Equal(1))
			Expect(foundSchedule.Appointments[1].StartTime).To(Equal(5))
			Expect(foundSchedule.Appointments[2].StartTime).To(Equal(10))
		})

		It("Should return a bad request for a non-numeric scheduleID", func() {
			url := fmt.Sprintf("%v/schedules/blamo", acceptanceUrl)
			req, _ := http.NewRequest("GET", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("Should return a not found for a not found scheduleID", func() {
			url := fmt.Sprintf("%v/schedules/-1", acceptanceUrl)
			req, _ := http.NewRequest("GET", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("DELETE /schedules/{scheduleID}", func() {
		It("Should return the deleted schedule for the associated schedule ID provided in the req path", func() {
			// Create schedule
			reqBody := []byte(`
				{
					"owner_name": "Tyrion Lannister"
				}
			`)

			url := fmt.Sprintf("%v/schedules", acceptanceUrl)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var createdSchedule scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&createdSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			// Delete the schedule that was just created
			url = fmt.Sprintf("%v/schedules/%v", acceptanceUrl, createdSchedule.ID)
			req, _ = http.NewRequest("DELETE", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var deletedSchedule scheduler.ScheduleResponse
			err = json.NewDecoder(res.Body).Decode(&deletedSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(deletedSchedule).To(Equal(createdSchedule))

			// Confirm that schedule is deleted
			url = fmt.Sprintf("%v/schedules/%v", acceptanceUrl, createdSchedule.ID)
			req, _ = http.NewRequest("GET", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("Should return a bad request for a non-numeric scheduleID", func() {
			url := fmt.Sprintf("%v/schedules/blamo", acceptanceUrl)
			req, _ := http.NewRequest("DELETE", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("Should return a not found for a not found scheduleID", func() {
			url := fmt.Sprintf("%v/schedules/-1", acceptanceUrl)
			req, _ := http.NewRequest("DELETE", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
