// hotelservice/service/hotels_service.go

package services

import (
	"context"
	"flights-api/hotelservice/domain"

	// "flights-api/hotelservice/utils"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	// Inicializar la conexión a MongoDB

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	mongoClient = client
}

const hotelCollection = "hotels"

// CreateHotel crea un nuevo hotel en la base de datos.

func CreateHotel(ctx context.Context, hotel domain.Hotel) error {
	// Insertar el hotel en MongoDB

	collection := mongoClient.Database("hotelsdb").Collection(hotelCollection)
	_, err := collection.InsertOne(ctx, hotel)
	if err != nil {
		return err
	}

	// Notificar a RabbitMQ sobre la creación del hotel
	NotifyHotelCreation(hotel)

	return nil
}

// GetHotelByID obtiene los detalles de un hotel por su ID.

func GetHotelByID(ctx context.Context, hotelID string) (domain.Hotel, error) {
	// Obtener el hotel desde MongoDB por ID

	collection := mongoClient.Database("hotelsdb").Collection(hotelCollection)
	filter := bson.D{{"id", hotelID}}
	var result domain.Hotel
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Hotel{}, err
	}

	return result, nil
}

// NotifyHotelCreation notifica a RabbitMQ sobre la creación de un hotel.
func NotifyHotelCreation(hotel domain.Hotel) {
	// \notificar a RabbitMQ sobre la creación de un hotel
	// conexión a RabbitMQ para enviar un mensaje de creación
	// Ejemplo de mensaje dummy
	message := fmt.Sprintf("Nuevo hotel creado: %s", hotel.Name)
	fmt.Println(message)
}

/*
	biblioteca oficial de MongoDB para Go (go.mongodb.org/mongo-driver/mongo)
	y la biblioteca de RabbitMQ (github.com/streadway/amqp).
*/
