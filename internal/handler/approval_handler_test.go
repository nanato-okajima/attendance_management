package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApprovalHandler_ApproveLeaveRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalService := mock.NewMockApprovalService(ctrl)
	mockCorrectionService := mock.NewMockAttendanceCorrectionService(ctrl)
	h := NewApprovalHandler(mockApprovalService, mockCorrectionService)

	approverID := 9999
	id := uint(1)
	comment := "Approved"

	tests := []struct {
		name       string
		input      ApproveLeaveRequestRequest
		setup      func(ma *mock.MockApprovalService)
		wantStatus int
	}{
		{
			name: "Success",
			input: ApproveLeaveRequestRequest{
				Comment: comment,
			},
			setup: func(ma *mock.MockApprovalService) {
				ma.EXPECT().ApproveLeaveRequest(gomock.Any(), id, approverID, comment).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "ServiceError",
			input: ApproveLeaveRequestRequest{
				Comment: comment,
			},
			setup: func(ma *mock.MockApprovalService) {
				ma.EXPECT().ApproveLeaveRequest(gomock.Any(), id, approverID, comment).Return(errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockApprovalService)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("employee_number", approverID)
			c.Params = []gin.Param{{Key: "id", Value: "1"}}

			jsonBytes, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest("POST", "/approvals/leave/1/approve", bytes.NewBuffer(jsonBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			h.ApproveLeaveRequest(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestApprovalHandler_GetPendingApprovals(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApprovalService := mock.NewMockApprovalService(ctrl)
	mockCorrectionService := mock.NewMockAttendanceCorrectionService(ctrl)
	h := NewApprovalHandler(mockApprovalService, mockCorrectionService)

	tests := []struct {
		name       string
		setup      func(ma *mock.MockApprovalService, mc *mock.MockAttendanceCorrectionService)
		wantStatus int
	}{
		{
			name: "Success",
			setup: func(ma *mock.MockApprovalService, mc *mock.MockAttendanceCorrectionService) {
				ma.EXPECT().GetPendingApprovals(gomock.Any()).Return([]*entity.LeaveRequest{
					{ID: 1},
				}, nil)
				mc.EXPECT().GetCorrectionRequests(gomock.Any(), 0, gomock.Any()).Return([]*entity.AttendanceCorrection{
					{ID: 2},
				}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockApprovalService, mockCorrectionService)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("GET", "/approvals/pending", nil)

			h.GetPendingApprovals(c)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
