// Ejecutar producer por cada consumidor
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func createExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"direct_logs", // nombre del intercambio
		"direct",      // tipo de intercambio (direct para enrutamiento directo)
		false,         // duradera
		false,         // eliminar cuando no se usa
		false,         // exclusiva
		false,         // no espera mensajes de autorespuesta
		nil,           // argumentos adicionales
	)
	return err
}

func main() {
	var conn *amqp.Connection
	var ch *amqp.Channel

	for {
		var err error

		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ. Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			fmt.Println("Failed to open a channel. Retrying in 5 seconds...")
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		err = createExchange(ch)
		if err != nil {
			fmt.Println("Failed to declare exchange. Retrying in 5 seconds...")
			conn.Close()
			ch.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	body := "Hello, RabbitMQ!"
	err := ch.Publish(
		"direct_logs", // nombre del intercambio
		"routing_key", // clave de enrutamiento común para todos los consumidores
		false,         // mandar a la cola si no hay consumidores
		false,         // mandar a la cola si no puede ser enrutado
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf(" [x] Sent %s\n", body)

	defer conn.Close()
	defer ch.Close()
}
