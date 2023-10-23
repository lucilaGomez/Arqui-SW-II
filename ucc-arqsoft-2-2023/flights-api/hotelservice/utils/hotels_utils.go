// hotelservice/utils/hotels_utils.go

package utils

import (
	"github.com/google/uuid"
)

// NewUUID genera un nuevo UUID.
func NewUUID() string {
	id := uuid.New()
	return id.String()
}
