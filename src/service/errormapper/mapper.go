package errormapper

type StorageErrorMapper[T any, R any] interface {
	MapError(code T) R
}
