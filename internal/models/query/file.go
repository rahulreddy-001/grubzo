package query

import "github.com/gofrs/uuid"

type FilesQuery struct {
	TenantID uint64
	IDs      []uuid.UUID
	OwnerId  *uint64
	Limit    int
	Offset   int
}
