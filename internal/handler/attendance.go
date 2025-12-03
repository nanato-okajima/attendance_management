package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nanato-okajima/attendance_management/internal/service"
)

type AttendanceHandler interface {
	ClockIn(c *gin.Context)
	ClockOut(c *gin.Context)
	GetToday(c *gin.Context)
	GetMonthly(c *gin.Context)
}

type attendanceHandler struct {
	attendanceService service.AttendanceService
}

func NewAttendanceHandler(attendanceService service.AttendanceService) AttendanceHandler {
	return &attendanceHandler{
		attendanceService: attendanceService,
	}
}

type ClockRequest struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (h *attendanceHandler) ClockIn(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")

	var req ClockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	attendance, err := h.attendanceService.ClockIn(employeeNumber, req.Latitude, req.Longitude, 1)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message":    "clocked in successfully",
		"attendance": attendance,
	})
}

func (h *attendanceHandler) ClockOut(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")

	var req ClockRequest
	_ = c.ShouldBindJSON(&req)

	attendance, err := h.attendanceService.ClockOut(employeeNumber, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":    "clocked out successfully",
		"attendance": attendance,
	})
}

func (h *attendanceHandler) GetToday(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")

	attendance, err := h.attendanceService.GetTodayAttendance(employeeNumber)
	if err != nil {
		c.JSON(404, gin.H{"error": "no attendance record found"})
		return
	}

	c.JSON(200, attendance)
}

func (h *attendanceHandler) GetMonthly(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(400, gin.H{"error": "year and month are required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid year"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(400, gin.H{"error": "invalid month"})
		return
	}

	attendances, err := h.attendanceService.GetMonthlyAttendances(employeeNumber, year, month)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch attendances"})
		return
	}

	c.JSON(200, gin.H{
		"year":        year,
		"month":       month,
		"attendances": attendances,
	})
}
