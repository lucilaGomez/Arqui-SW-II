// hotelservice/domain/hotels_domain.go

package domain

// Hotel representa la estructura de datos del hotel.
type Hotel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}
