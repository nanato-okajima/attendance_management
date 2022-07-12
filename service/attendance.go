package service

import (
	"github.com/nanato-okajima/attendance_management/service/repository"
)

type AttendanceService interface {
	Register(attendance *repository.Attendance) error
	List() (*[]repository.Attendance, error)
	Update(attendance *repository.Attendance, id string) error
	Delete(id string) error
}

type attendanceServiceImplementaion struct {
	ad repository.AttendanceRepository
}

func NewAttendanceService(ad repository.AttendanceRepository) AttendanceService {
	return &attendanceServiceImplementaion{ad}
}

func (asi attendanceServiceImplementaion) Register(attendance *repository.Attendance) error {
	if err := asi.ad.Insert(attendance); err != nil {
		return err
	}

	return nil
}

func (asi attendanceServiceImplementaion) List() (*[]repository.Attendance, error) {
	attendances := []repository.Attendance{}
	if err := asi.ad.Fetch(&attendances); err != nil {
		return nil, err
	}

	return &attendances, nil
}

func (asi attendanceServiceImplementaion) Update(attendance *repository.Attendance, id string) error {
	if err := asi.ad.Update(attendance, id); err != nil {
		return err
	}

	return nil
}

func (asi attendanceServiceImplementaion) Delete(id string) error {
	if err := asi.ad.Delete(id); err != nil {
		return err
	}

	return nil
}
