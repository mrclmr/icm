package http

import (
	"github.com/mrclmr/icm/internal/cont"
	"golang.org/x/net/context"
)

// OwnersGetter downloads owners.
type OwnersGetter interface {
	GetOwners(context.Context) ([]cont.Owner, error)
}
