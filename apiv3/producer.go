// Cola por consumidor (mensaje al tiempo)
package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"direct_logs", // nombre del intercambio
		"direct",      // tipo de intercambio (direct para enrutamiento directo)
		false,         // duradera
		false,         // eliminar cuando no se usa
		false,         // exclusiva
		false,         // no espera mensajes de autorespuesta
		nil,           // argumentos adicionales
	)
	failOnError(err, "Failed to declare an exchange")

	body := "Hello, RabbitMQ!"
	err = ch.Publish(
		"direct_logs",   // nombre del intercambio
		"routing_key_1", // clave de enrutamiento
		false,           // mandar a la cola si no hay consumidores
		false,           // mandar a la cola si no puede ser enrutado
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf(" [x] Sent %s\n", body)
}
