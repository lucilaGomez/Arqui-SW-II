## Arquitectura Software 2 (2023)

Objetivo
Implementar un sistema de microservicios, mediante el cual una cadena de hoteles
disponibiliza su oferta de forma local y se integra con un proveedor central para validar los
pedidos al momento de la reserva.
Arquitectura
La aplicación debe contar con al menos 4 microservicios:
1. Microservicio frontend
2. Microservicio de ficha de hotel
3. Microservicio de búsqueda de hotel
4. Microservicio de usuarios, reserva y disponibilidad
Microservicio frontend
Debe tener 4 pantallas:
● Pantalla inicial con 3 campos de búsqueda: ciudad, fecha desde y fecha hasta.
● Pantalla de resultado de la búsqueda con el listado de hoteles con disponibilidad
para esa ciudad: nombre, descripción y thumbnail.
● Pantalla de detalle del hotel con la información completa: nombre, descripción, fotos,
amenities y botón de reserva.
● Congrats, con confirmación exitosa o rechazo de la reserva.
Microservicio de ficha de hotel
Microservicio RESTful en Go que almacena la información de los hoteles con formato
documental en una base de datos no relacional (MongoDB) y disponibiliza esa información
de forma RESTful. Soporta y valida la información de creación y modificación de hoteles así
también como la obtención de la información por id de hotel. Notifica a una cola de
mensajes (RabbitMQ) cuando se crea o se modifica un hotel.
Universidad Católica de Córdoba
Microservicio de búsqueda de hotel
Microservicio RESTful en Go que contiene un motor de búsqueda (Solr) que disponibiliza la
información de los hoteles en base a los criterios de búsqueda que se le especifiquen.
Escucha la cola de mensajes de creación y actualización del servicio de ficha de hoteles y
sincroniza los cambios haciendo un GET por ID. Consulta concurrentemente la
disponibilidad de los resultados para filtrar en la búsqueda mediante un atributo dinámico
que indique si tiene disponibilidad o no, esto significa tener un atributo por ej “availability”,
pero que no se persista, sino que devuelva los resultados de la consulta al servicio de
disponibilidad de forma interna y concurrente para el listado de resultados.
Microservicio de usuarios, reserva y disponibilidad
Microservicio RESTtful en Go que contiene una base de datos relacional (MySQL) con la
información de los usuarios/clientes del sitio y sus reservas. El microservicio tiene que tener
la capacidad de retornar la disponibilidad de los hoteles con un caché distribuido
(Memcached) con TTL de 10 segundos que evita su cálculo subsiguiente. De esa forma, se
optimizan las consultas suponiendo una consistencia eventual de 10 segundos. Al momento
de realizar las reservas se tiene que validar externamente con el sitio Amadeus (anexo 1),
dado que la oferta puede estar distribuída y su validación debe estar concentrada por este
proveedor externo de servicios de booking. Este microservicio debe proveer todos los
endpoints necesarios para la operación del cliente y gestión de sus reservas. El mapping
entre el ID interno del hotel y el ID de Amadeus lo conoce este microservicio. Los nuevos
hoteles se
