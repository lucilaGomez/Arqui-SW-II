// hotelservice/app/app.go

package app

import (
	"flights-api/hotelservice/controllers"

	"github.com/gin-gonic/gin"
)

func StartApp() {
	router := gin.Default()

	// Configuración de rutas para el controlador de hoteles
	controllers.SetupRoutes(router)

	// Inicia el servidor en el puerto 8081
	_ = router.Run(":8081")
}
