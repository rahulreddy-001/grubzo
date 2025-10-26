package query

import "github.com/gofrs/uuid"

type FilesQuery struct {
	TenantID uint
	IDs      []uuid.UUID
	OwnerId  *uint
	Limit    int
	Offset   int
}
