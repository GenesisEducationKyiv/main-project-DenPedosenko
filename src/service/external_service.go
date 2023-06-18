package service

type ExternalService interface {
	GetCurrentBTCToUAHRate() (float64, error)
}
