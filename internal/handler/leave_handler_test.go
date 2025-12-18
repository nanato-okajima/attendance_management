package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLeaveHandler_CreateLeaveRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockLeaveService(ctrl)
	h := NewLeaveHandler(mockService)

	employeeNumber := 1001
	startDate := "2023-01-01"
	endDate := "2023-01-02"
	reason := "Vacation"

	tests := []struct {
		name       string
		input      CreateLeaveRequestRequest
		setup      func(m *mock.MockLeaveService)
		wantStatus int
	}{
		{
			name: "Success",
			input: CreateLeaveRequestRequest{
				LeaveType: 1,
				StartDate: startDate,
				EndDate:   endDate,
				Reason:    reason,
			},
			setup: func(m *mock.MockLeaveService) {
				m.EXPECT().CreateLeaveRequest(gomock.Any(), gomock.Any()).Return(&entity.LeaveRequest{
					ID:             1,
					EmployeeNumber: employeeNumber,
					LeaveType:      entity.LeaveTypePaidLeave,
					StartDate:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:        time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
					Reason:         reason,
				}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "ValidationError_InvalidDate",
			input: CreateLeaveRequestRequest{
				LeaveType: 1,
				StartDate: "invalid",
				EndDate:   endDate,
				Reason:    reason,
			},
			setup: func(m *mock.MockLeaveService) {
				// No call expected
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "ServiceError",
			input: CreateLeaveRequestRequest{
				LeaveType: 1,
				StartDate: startDate,
				EndDate:   endDate,
				Reason:    reason,
			},
			setup: func(m *mock.MockLeaveService) {
				m.EXPECT().CreateLeaveRequest(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockService)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("employee_number", employeeNumber)

			jsonBytes, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest("POST", "/leave/requests", bytes.NewBuffer(jsonBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			h.CreateLeaveRequest(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLeaveHandler_GetLeaveRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockLeaveService(ctrl)
	h := NewLeaveHandler(mockService)

	employeeNumber := 1001

	tests := []struct {
		name       string
		role       string
		query      string
		setup      func(m *mock.MockLeaveService)
		wantStatus int
	}{
		{
			name:  "Success_User",
			role:  "user",
			query: "",
			setup: func(m *mock.MockLeaveService) {
				m.EXPECT().GetLeaveRequests(gomock.Any(), employeeNumber, nil).Return([]*entity.LeaveRequest{}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "Success_Admin_WithQuery",
			role:  "admin",
			query: "?employee_number=1002",
			setup: func(m *mock.MockLeaveService) {
				m.EXPECT().GetLeaveRequests(gomock.Any(), 1002, nil).Return([]*entity.LeaveRequest{}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockService)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("employee_number", employeeNumber)
			c.Set("role", tt.role)

			c.Request = httptest.NewRequest("GET", "/leave/requests"+tt.query, nil)

			h.GetLeaveRequests(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLeaveHandler_GetRemainingPaidLeaveDays(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockLeaveService(ctrl)
	h := NewLeaveHandler(mockService)

	employeeNumber := 1001

	tests := []struct {
		name       string
		setup      func(m *mock.MockLeaveService)
		wantStatus int
	}{
		{
			name: "Success",
			setup: func(m *mock.MockLeaveService) {
				m.EXPECT().GetRemainingPaidLeaveDays(gomock.Any(), gomock.Any()).Return(15.0, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockService)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("employee_number", employeeNumber)

			c.Request = httptest.NewRequest("GET", "/leave/remaining-days", nil)

			h.GetRemainingPaidLeaveDays(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
