package notification

type NotificationService interface {
	Send([]string, float64) error
}
