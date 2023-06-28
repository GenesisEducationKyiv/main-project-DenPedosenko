package errormapper

import (
	"net/http"

	"ses.genesis.com/exchange-web-service/src/persistent"
)

type StorageErrorToHTTPMapper struct {
}

func NewStorageErrorToHTTPMapper() *StorageErrorToHTTPMapper {
	return &StorageErrorToHTTPMapper{}
}

func (mapper *StorageErrorToHTTPMapper) MapError(err persistent.StorageError) int {
	switch err.Code {
	case persistent.Conflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
