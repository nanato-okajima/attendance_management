package repository

type Attendance struct {
	AttendanceId     int64  `json:"attendance_id"`
	EmployeeId       int64  `json:"employee_id"`
	OpeningTime      string `json:"opening_time"`
	ClosingTime      string `json:"closing_time"`
	AttendanceStatus int64  `json:"attendance_status"`
}

type AttendanceRepository interface {
	Insert(attendance *Attendance) error
	Fetch(attendances *[]Attendance) error
	Update(attendance *Attendance, id string) error
	Delete(id string) error
}
