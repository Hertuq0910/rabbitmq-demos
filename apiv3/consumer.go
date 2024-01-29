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

	q, err := ch.QueueDeclare(
		"",    // nombre de la cola (RabbitMQ generará un nombre único)
		false, // duradera
		false, // eliminar cuando no se usa
		true,  // exclusiva
		false, // no espera mensajes de autorespuesta
		nil,   // argumentos adicionales
	)
	failOnError(err, "Failed to declare a queue")

	// Claves de enrutamiento a las que el consumidor está interesado
	keys := []string{"routing_key_1"}

	// Vincular la cola al intercambio utilizando las claves de enrutamiento
	for _, key := range keys {
		err = ch.QueueBind(
			q.Name,        // nombre de la cola
			key,           // routing key
			"direct_logs", // nombre del intercambio
			false,
			nil,
		)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // nombre de la cola
		"",     // consumer
		false,  // auto-ack (manual ack para controlar la confirmación)
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf(" [x] Received %s\n", d.Body)

			// Simular trabajo
			time.Sleep(time.Second * 2)

			fmt.Printf(" [x] Done\n")

			// Confirma el procesamiento del mensaje
			d.Ack(false)
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
