// hotelservice/service/search_service.go

/*
Solr Setup:\
Asegúrate de tener un servidor Solr en ejecución y que los índices y campos necesarios
estén configurados según tu modelo de datos
*/
package services

import (
	"context"
	"log"
	"sync"

	solr "github.com/rtt/Go-Solr"
	"github.com/streadway/amqp"
)

// HotelSearchService maneja la búsqueda y sincronización de hoteles.
type HotelSearchService struct {
	solrClient   *solr.SolrInterface // ver por que sale error, en teoria es pq no se importo correctamente el "github.com/rtt/Go-Solr"
	rabbitMQConn *amqp.Connection
}

// NewHotelSearchService crea una nueva instancia de HotelSearchService.
// inicializa el servicio

func NewHotelSearchService(solrURL, rabbitMQURL string) (*HotelSearchService, error) {
	solrClient, err := solr.Init(solrURL, "hotels")
	if err != nil {
		return nil, err
	}

	rabbitMQConn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	return &HotelSearchService{
		solrClient:   &solrClient,
		rabbitMQConn: rabbitMQConn,
	}, nil
}

// StartSyncListener inicia el oyente de eventos de sincronización de RabbitMQ.
// Inicia un oyente de eventos en RabbitMQ para recibir mensajes de sincronización.
func (s *HotelSearchService) StartSyncListener() {
	channel, err := s.rabbitMQConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		// Cada vez que un hotel se crea o se actualiza,
		// se envía un mensaje a la cola hotel_sync_queue
		"hotel_sync_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		// Sincronizar el hotel en Solr al recibir un mensaje de sincronización
		hotelID := string(msg.Body)
		//El oyente recibe estos mensajes y sincroniza el hotel correspondiente
		s.SyncHotelByID(hotelID)
		//Obtiene la información del hotel por ID y la agrega al índice de Solr.
	}
}

// SyncHotelByID sincroniza un hotel específico en Solr por ID.

func (s *HotelSearchService) SyncHotelByID(hotelID string) {
	hotel, err := GetHotelByID(context.Background(), hotelID)
	if err != nil {
		log.Println("Error obteniendo el hotel:", err)
		return
	}

	if err := s.solrClient.Add(hotel); err != nil {
		log.Println("Error sincronizando el hotel en Solr:", err)
		return
	}

	if err := s.solrClient.Commit(); err != nil {
		log.Println("Error haciendo commit en Solr:", err)
		return
	}

	log.Printf("Hotel sincronizado en Solr: %s\n", hotel.Name)
}

// SearchHotels realiza una búsqueda de hoteles en Solr.

func (s *HotelSearchService) SearchHotels(query string) ([]HotelSearchResult, error) {

	// Realiza una búsqueda de hoteles en Solr utilizando el motor de búsqueda
	// Para cada resultado de la búsqueda, obtiene los detalles del hotel y verifica la disponibilidad

	results, err := s.solrClient.Search(query, 0, 10, nil)
	if err != nil {
		return nil, err
	}

	var searchResults []HotelSearchResult
	var wg sync.WaitGroup

	for _, result := range results.Results {
		wg.Add(1)
		go func(result solr.Document) {
			defer wg.Done()
			hotelID := result.Field("id").(string)
			hotel, err := GetHotelByID(context.Background(), hotelID)
			if err != nil {
				log.Println("Error obteniendo el hotel:", err)
				return
			}

			// Verificar la disponibilidad concurrentemente
			available := make(chan bool)
			go func(hotelID string) {
				defer close(available)
				available <- CheckHotelAvailability(context.Background(), hotelID)
			}(hotelID)

			searchResult := HotelSearchResult{
				ID:           hotel.ID,
				Name:         hotel.Name,
				Description:  hotel.Description,
				Thumbnail:    hotel.Thumbnail,
				Availability: <-available,
			}

			searchResults = append(searchResults, searchResult)
		}(result)
	}

	wg.Wait()
	return searchResults, nil
}

// CheckHotelAvailability verifica la disponibilidad de un hotel.
func CheckHotelAvailability(ctx context.Context, hotelID string) bool {
	// simulamos una verificación de disponibilidad
	// COMPLETAR?
	return true
}

// HotelSearchResult representa el resultado de la búsqueda de hoteles.
type HotelSearchResult struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	Availability bool   `json:"availability"`
}
