package persistent

type Storage interface {
	AllEmails() ([]string, error)
	SaveEmailToStorage(email string) *StorageError
	IsEmailAlreadyExists(newEmail string) bool
}
