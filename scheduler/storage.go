package scheduler

type Schedule struct {
	ID           int                 `json:"id"`
	OwnerName    string              `json:"owner_name"`
	Appointments map[int]Appointment `json:"appointments"`
}

type Appointment struct {
	ID         int `json:"id"`
	ScheduleID int `json:"schedule_id"`
	StartTime  int `json:"start_time"`
	EndTime    int `json:"end_time"`
}

var SchedulesCreatedCount int
var ScheduleCollection = make(map[int]Schedule)

var AppointmentsCreatedCount int
