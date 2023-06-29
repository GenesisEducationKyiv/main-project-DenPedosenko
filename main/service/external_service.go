package service

type ExternalService interface {
	CurrentBTCToUAHRate() (float64, error)
}
