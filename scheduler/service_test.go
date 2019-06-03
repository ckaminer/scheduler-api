package scheduler_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/ckaminer/schedule-api/scheduler"
)

var _ = Describe("Service", func() {
	Context("#ValidateAppointmentInput", func() {
		It("Should return true if the end time is greater than the start time and there are no existing appts", func() {
			a := Appointment{
				StartTime: 9,
				EndTime:   10,
			}

			valid := ValidateAppointmentInput(Schedule{}, a)

			Expect(valid).To(BeTrue())
		})

		It("Should return false if the end time is less than or equal to the start time", func() {
			a := Appointment{
				StartTime: 9,
				EndTime:   5,
			}

			valid := ValidateAppointmentInput(Schedule{}, a)

			Expect(valid).To(BeFalse())

			a = Appointment{
				StartTime: 5,
				EndTime:   5,
			}

			valid = ValidateAppointmentInput(Schedule{}, a)

			Expect(valid).To(BeFalse())
		})

		It("Should return false if the start time is 0", func() {
			a := Appointment{
				StartTime: 0,
				EndTime:   10,
			}

			valid := ValidateAppointmentInput(Schedule{}, a)

			Expect(valid).To(BeFalse())
		})

		Context("Should return true if there are no overlaps with existing appointments", func() {
			It("Beginning of schedule", func() {
				scheduledAppointments := map[int]Appointment{
					1: Appointment{
						StartTime: 4,
						EndTime:   8,
					},
					2: Appointment{
						StartTime: 11,
						EndTime:   13,
					},
				}

				s := Schedule{
					Appointments: scheduledAppointments,
				}

				a := Appointment{
					StartTime: 1,
					EndTime:   3,
				}

				valid := ValidateAppointmentInput(s, a)

				Expect(valid).To(BeTrue())
			})

			It("End of schedule", func() {
				scheduledAppointments := map[int]Appointment{
					1: Appointment{
						StartTime: 4,
						EndTime:   8,
					},
					2: Appointment{
						StartTime: 11,
						EndTime:   13,
					},
				}

				s := Schedule{
					Appointments: scheduledAppointments,
				}

				a := Appointment{
					StartTime: 14,
					EndTime:   18,
				}

				valid := ValidateAppointmentInput(s, a)

				Expect(valid).To(BeTrue())
			})

			It("Sandwiched bewteen two appointments", func() {
				scheduledAppointments := map[int]Appointment{
					1: Appointment{
						StartTime: 4,
						EndTime:   8,
					},
					2: Appointment{
						StartTime: 11,
						EndTime:   13,
					},
				}

				s := Schedule{
					Appointments: scheduledAppointments,
				}

				a := Appointment{
					StartTime: 9,
					EndTime:   10,
				}

				valid := ValidateAppointmentInput(s, a)

				Expect(valid).To(BeTrue())
			})
		})

		Context("Should return false if there is any overlap with existing appointments on the given schedule", func() {
			It("EndTime equal to existing StartTime", func() {
				scheduledAppointments := map[int]Appointment{
					1: Appointment{
						StartTime: 4,
						EndTime:   8,
					},
				}
				s := Schedule{
					Appointments: scheduledAppointments,
				}

				a := Appointment{
					StartTime: 1,
					EndTime:   4,
				}

				valid := ValidateAppointmentInput(s, a)

				Expect(valid).To(BeFalse())
			})

			It("EndTime in range of existing appointment", func() {
				scheduledAppointments := map[int]Appointment{
					1: Appointment{
						StartTime: 4,
						EndTime:   8,
					},
				}
				s := Schedule{
					Appointments: scheduledAppointments,
				}

				a := Appointment{
					StartTime: 1,
					EndTime:   5,
				}

				valid := ValidateAppointmentInput(s, a)

				Expect(valid).To(BeFalse())
			})

			It("StartTime equal to existing EndTime", func() {
				scheduledAppointments := map[int]Appointment{
					1: Appointment{
						StartTime: 4,
						EndTime:   8,
					},
				}
				s := Schedule{
					Appointments: scheduledAppointments,
				}

				a := Appointment{
					StartTime: 8,
					EndTime:   10,
				}

				valid := ValidateAppointmentInput(s, a)

				Expect(valid).To(BeFalse())
			})

			It("StartTime in range of existing appointment", func() {
				scheduledAppointments := map[int]Appointment{
					1: Appointment{
						StartTime: 4,
						EndTime:   8,
					},
				}
				s := Schedule{
					Appointments: scheduledAppointments,
				}

				a := Appointment{
					StartTime: 6,
					EndTime:   10,
				}

				valid := ValidateAppointmentInput(s, a)

				Expect(valid).To(BeFalse())
			})
		})
	})
})
