package persistent

type PersistentStorage interface {
	AllEmails() ([]string, error)
	SaveEmailToStorage(email string) (int, error)
	isEmailAlreadyExists(newEmail string) bool
}
