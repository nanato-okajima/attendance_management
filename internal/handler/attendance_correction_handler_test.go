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

func TestAttendanceCorrectionHandler_CreateCorrectionRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAttendanceCorrectionService(ctrl)
	h := NewAttendanceCorrectionHandler(mockService)

	employeeNumber := 1001
	attendanceID := uint(1)
	reason := "Forgot to clock out"
	now := time.Now()

	nowStr := now.Format(time.RFC3339)

	tests := []struct {
		name       string
		input      CreateCorrectionRequestRequest
		setup      func(m *mock.MockAttendanceCorrectionService)
		wantStatus int
	}{
		{
			name: "Success",
			input: CreateCorrectionRequestRequest{
				AttendanceID:         attendanceID,
				CorrectionType:       int(entity.CorrectionTypeClockOut),
				CorrectedClosingTime: nowStr,
				Reason:               reason,
			},
			setup: func(m *mock.MockAttendanceCorrectionService) {
				m.EXPECT().CreateCorrectionRequest(gomock.Any(), gomock.Any()).Return(&entity.AttendanceCorrection{
					ID:             1,
					EmployeeNumber: employeeNumber,
				}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "ServiceError",
			input: CreateCorrectionRequestRequest{
				AttendanceID:         attendanceID,
				CorrectionType:       int(entity.CorrectionTypeClockOut),
				CorrectedClosingTime: nowStr,
				Reason:               reason,
			},
			setup: func(m *mock.MockAttendanceCorrectionService) {
				m.EXPECT().CreateCorrectionRequest(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
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
			c.Request = httptest.NewRequest("POST", "/attendance-corrections/requests", bytes.NewBuffer(jsonBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			h.CreateCorrectionRequest(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAttendanceCorrectionHandler_GetCorrectionRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAttendanceCorrectionService(ctrl)
	h := NewAttendanceCorrectionHandler(mockService)

	employeeNumber := 1001

	tests := []struct {
		name       string
		role       string
		query      string
		setup      func(m *mock.MockAttendanceCorrectionService)
		wantStatus int
	}{
		{
			name:  "Success_User",
			role:  "user",
			query: "",
			setup: func(m *mock.MockAttendanceCorrectionService) {
				m.EXPECT().GetCorrectionRequests(gomock.Any(), employeeNumber, nil).Return([]*entity.AttendanceCorrection{}, nil)
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

			c.Request = httptest.NewRequest("GET", "/attendance-corrections/requests"+tt.query, nil)

			h.GetCorrectionRequests(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAttendanceCorrectionHandler_GetCorrectionRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAttendanceCorrectionService(ctrl)
	h := NewAttendanceCorrectionHandler(mockService)

	employeeNumber := 1001
	id := uint(1)

	tests := []struct {
		name       string
		role       string
		setup      func(m *mock.MockAttendanceCorrectionService)
		wantStatus int
	}{
		{
			name: "Success_Owner",
			role: "user",
			setup: func(m *mock.MockAttendanceCorrectionService) {
				m.EXPECT().GetCorrectionRequestByID(gomock.Any(), id).Return(&entity.AttendanceCorrection{
					ID:             id,
					EmployeeNumber: employeeNumber,
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Forbidden_OtherUser",
			role: "user",
			setup: func(m *mock.MockAttendanceCorrectionService) {
				m.EXPECT().GetCorrectionRequestByID(gomock.Any(), id).Return(&entity.AttendanceCorrection{
					ID:             id,
					EmployeeNumber: 1002, // Different user
				}, nil)
			},
			wantStatus: http.StatusBadRequest, // Validation error "権限がありません"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockService)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("employee_number", employeeNumber)
			c.Set("role", tt.role)
			c.Params = []gin.Param{{Key: "id", Value: "1"}}

			c.Request = httptest.NewRequest("GET", "/attendance-corrections/requests/1", nil)

			h.GetCorrectionRequest(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
