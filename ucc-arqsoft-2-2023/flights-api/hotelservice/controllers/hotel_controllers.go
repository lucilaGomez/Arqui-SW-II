// hotelservice/controllers/hotels_controllers.go

package controllers

import (
	"context"
	"flights-api/hotelservice/domain"
	"flights-api/hotelservice/services"

	// VER POR QUE NO SE PUEDE IMPORTAR ESTO
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura las rutas para el controlador de hoteles.
func SetupRoutes(router *gin.Engine) {
	router.POST("/hotel", createHotel)
	router.GET("/hotel/:id", getHotel)
}

// createHotel maneja la creación de un nuevo hotel.
func createHotel(c *gin.Context) {
	var hotel domain.Hotel
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateHotel(context.Background(), hotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el hotel"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hotel creado exitosamente"})
}

// getHotel obtiene los detalles de un hotel por ID.
func getHotel(c *gin.Context) {
	hotelID := c.Param("id")

	hotel, err := services.GetHotelByID(context.Background(), hotelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener detalles del hotel"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}
