package http

import (
	"context"

	"github.com/mrclmr/icm/internal/cont"
)

// OwnersGetter downloads owners.
type OwnersGetter interface {
	GetOwners(context.Context) ([]cont.Owner, error)
}
