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

var _ = Describe("Appointment Routes", func() {
	var scheduleID int

	BeforeEach(func() {
		// Create Schedule
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

		scheduleID = s.ID
	})

	Context("POST /schedules/{scheduleID}/appointments", func() {
		It("Should return a created appointment with a generated ID upon success", func() {
			reqBody := []byte(`
				{
					"start_time": 5,
					"end_time": 9
				}
			`)

			url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, scheduleID)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var a scheduler.Appointment
			err = json.NewDecoder(res.Body).Decode(&a)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))
			Expect(a.ID).To(BeNumerically(">", 0))
			Expect(a.ScheduleID).To(Equal(scheduleID))
			Expect(a.StartTime).To(Equal(5))
			Expect(a.EndTime).To(Equal(9))

			// Send additional request to confirm ID is being incremented
			reqBody = []byte(`
				{
					"start_time": 10,
					"end_time": 11
				}
			`)

			newReq, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))
			res, err = client.Do(newReq)
			if err != nil {
				Fail("Failed to send request")
			}

			var a2 scheduler.Appointment
			err = json.NewDecoder(res.Body).Decode(&a2)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))
			Expect(a2.ID).To(Equal(a.ID + 1))
		})

		It("Should return a not found if the schedule is not found", func() {
			reqBody := []byte(`
					{
						"start_time": 5,
						"end_time": 9
					}
				`)

			url := fmt.Sprintf("%v/schedules/-1/appointments", acceptanceUrl)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("Should return a bad request if the request body is invalid", func() {
			reqBody := []byte(`
					{
						"start_time": true
					}
				`)

			url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, scheduleID)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})

		Context("Should return a unprocessable entity if a time validation fails", func() {
			It("start time greater than end time", func() {
				reqBody := []byte(`
					{
						"start_time": 5,
						"end_time": 3
					}
				`)

				url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, scheduleID)
				req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

				client := http.Client{}

				res, err := client.Do(req)
				if err != nil {
					Fail("Failed to send request")
				}

				Expect(res.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			})

			It("start time is zero", func() {
				reqBody := []byte(`
					{
						"start_time": 0,
						"end_time": 3
					}
				`)

				url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, scheduleID)
				req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

				client := http.Client{}

				res, err := client.Do(req)
				if err != nil {
					Fail("Failed to send request")
				}

				Expect(res.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			})

			It("range overlaps existing appointment", func() {
				reqBody := []byte(`
					{
						"start_time": 1,
						"end_time": 4
					}
				`)

				url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, scheduleID)
				req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

				client := http.Client{}

				res, err := client.Do(req)
				if err != nil {
					Fail("Failed to send request")
				}

				Expect(res.StatusCode).To(Equal(http.StatusCreated))

				reqBody = []byte(`
					{
						"start_time": 3,
						"end_time": 6
					}
				`)

				req, _ = http.NewRequest("POST", url, bytes.NewReader(reqBody))

				res, err = client.Do(req)
				if err != nil {
					Fail("Failed to send request")
				}

				Expect(res.StatusCode).To(Equal(http.StatusUnprocessableEntity))
			})
		})
	})

	Context("GET /schedules/{scheduleID}/appointments/{appointmentID}", func() {
		It("Should return a appointment for the associated appointment ID provided in the req path", func() {
			// Create appt
			reqBody := []byte(`
				{
					"start_time": 5,
					"end_time": 9
				}
			`)

			url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, scheduleID)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var createdAppointment scheduler.Appointment
			err = json.NewDecoder(res.Body).Decode(&createdAppointment)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			// Find the appointment that was just created
			url = fmt.Sprintf("%v/schedules/%v/appointments/%v", acceptanceUrl, scheduleID, createdAppointment.ID)
			req, _ = http.NewRequest("GET", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var foundAppointment scheduler.Appointment
			err = json.NewDecoder(res.Body).Decode(&foundAppointment)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(foundAppointment).To(Equal(createdAppointment))
		})

		It("Should return a bad request for a non-numeric scheduleID or appointmentID", func() {
			url := fmt.Sprintf("%v/schedules/blamo/appointments/1", acceptanceUrl)
			req, _ := http.NewRequest("GET", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))

			url = fmt.Sprintf("%v/schedules/1/appointments/blamo", acceptanceUrl)
			req, _ = http.NewRequest("GET", url, nil)
			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("Should return a not found for a not found scheduleID", func() {
			url := fmt.Sprintf("%v/schedules/-1/appointments/1", acceptanceUrl)
			req, _ := http.NewRequest("GET", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("Should return a not found for a not found appointmentID", func() {
			url := fmt.Sprintf("%v/schedules/%v/appointments/-1", acceptanceUrl, scheduleID)
			req, _ := http.NewRequest("GET", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("DELETE /appointments/{scheduleID}", func() {
		It("Should return the deleted schedule for the associated schedule ID provided in the req path", func() {
			// create appt
			reqBody := []byte(`
				{
					"start_time": 5,
					"end_time": 9
				}
			`)

			url := fmt.Sprintf("%v/schedules/%v/appointments", acceptanceUrl, scheduleID)
			req, _ := http.NewRequest("POST", url, bytes.NewReader(reqBody))

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var createdAppointment scheduler.Appointment
			err = json.NewDecoder(res.Body).Decode(&createdAppointment)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusCreated))

			// Delete the appointment that was just created
			url = fmt.Sprintf("%v/schedules/%v/appointments/%v", acceptanceUrl, scheduleID, createdAppointment.ID)
			req, _ = http.NewRequest("DELETE", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			var deletedAppointment scheduler.Appointment
			err = json.NewDecoder(res.Body).Decode(&deletedAppointment)
			if err != nil {
				Fail("Failed to decode response body")
			}

			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(deletedAppointment).To(Equal(createdAppointment))

			// Confirm that appt is deleted
			url = fmt.Sprintf("%v/schedules/%v/appointments/%v", acceptanceUrl, scheduleID, createdAppointment.ID)
			req, _ = http.NewRequest("GET", url, nil)

			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("Should return a bad request for a non-numeric scheduleID or appointmentID", func() {
			url := fmt.Sprintf("%v/schedules/blamo/appointments/2", acceptanceUrl)
			req, _ := http.NewRequest("DELETE", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))

			url = fmt.Sprintf("%v/schedules/2/appointments/blamo", acceptanceUrl)
			req, _ = http.NewRequest("DELETE", url, nil)
			res, err = client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("Should return a not found for a not found scheduleID", func() {
			url := fmt.Sprintf("%v/schedules/-1/appointments/2", acceptanceUrl)
			req, _ := http.NewRequest("DELETE", url, nil)

			client := http.Client{}

			res, err := client.Do(req)
			if err != nil {
				Fail("Failed to send request")
			}

			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("Should return a not found for a not found scheduleID", func() {
			url := fmt.Sprintf("%v/schedules/%v/appointments/-1", acceptanceUrl, scheduleID)
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
