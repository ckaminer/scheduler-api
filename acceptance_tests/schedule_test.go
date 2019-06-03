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

			req, _ := http.NewRequest("POST", "http://localhost:8080/schedules", bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var s scheduler.Schedule
			err = json.NewDecoder(res.Body).Decode(&s)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))
			Expect(s.OwnerName).To(Equal("Tyrion Lannister"))
			Expect(s.ID).To(BeNumerically(">", 0))

			// Send additional request to confirm ID is being incremented
			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var s2 scheduler.Schedule
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

			req, _ := http.NewRequest("POST", "http://localhost:8080/schedules", bytes.NewReader(reqBody))

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

			req, _ := http.NewRequest("POST", "http://localhost:8080/schedules", bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var createdSchedule scheduler.Schedule
			err = json.NewDecoder(res.Body).Decode(&createdSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			// Find the schedule that was just created
			url := fmt.Sprintf("http://localhost:8080/schedules/%v", createdSchedule.ID)
			req, _ = http.NewRequest("GET", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var foundSchedule scheduler.Schedule
			err = json.NewDecoder(res.Body).Decode(&foundSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(foundSchedule).To(Equal(createdSchedule))
		})

		It("Should return a bad request for a non-numeric scheduleID", func() {
			req, _ := http.NewRequest("GET", "http://localhost:8080/schedules/blamo", nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("Should return a not found for a not found scheduleID", func() {
			req, _ := http.NewRequest("GET", "http://localhost:8080/schedules/48", nil)

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

			req, _ := http.NewRequest("POST", "http://localhost:8080/schedules", bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var createdSchedule scheduler.Schedule
			err = json.NewDecoder(res.Body).Decode(&createdSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			// Delete the schedule that was just created
			url := fmt.Sprintf("http://localhost:8080/schedules/%v", createdSchedule.ID)
			req, _ = http.NewRequest("DELETE", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var deletedSchedule scheduler.Schedule
			err = json.NewDecoder(res.Body).Decode(&deletedSchedule)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(deletedSchedule).To(Equal(createdSchedule))

			// Confirm that schedule is deleted
			url = fmt.Sprintf("http://localhost:8080/schedules/%v", createdSchedule.ID)
			req, _ = http.NewRequest("GET", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("Should return a bad request for a non-numeric scheduleID", func() {
			req, _ := http.NewRequest("DELETE", "http://localhost:8080/schedules/blamo", nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("Should return a not found for a not found scheduleID", func() {
			req, _ := http.NewRequest("DELETE", "http://localhost:8080/schedules/48", nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
