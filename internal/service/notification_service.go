package service

import (
	"context"
	"log"
)

// NotificationService 通知サービスインターフェース
type NotificationService interface {
	// NotifyApprovalResult 承認結果を通知
	NotifyApprovalResult(ctx context.Context, employeeNumber int, result, message string)

	// NotifyPendingApproval 承認待ちを通知
	NotifyPendingApproval(ctx context.Context, approverID int, requestType, requestID string)
}

type notificationService struct {
	// TODO: メール送信やLINE通知の実装
}

// NewNotificationService 通知サービスを作成
func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (s *notificationService) NotifyApprovalResult(ctx context.Context, employeeNumber int, result, message string) {
	// TODO: 実際の通知実装（メール、LINE等）
	log.Printf("通知: 従業員番号 %d に承認結果を通知 - %s: %s", employeeNumber, result, message)
}

func (s *notificationService) NotifyPendingApproval(ctx context.Context, approverID int, requestType, requestID string) {
	// TODO: 実際の通知実装（メール、LINE等）
	log.Printf("通知: 承認者 %d に承認待ちを通知 - %s: %s", approverID, requestType, requestID)
}
