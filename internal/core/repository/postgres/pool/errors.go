package core_postgres_pool

import "errors"

var (
	ErrNotRows            = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
)
