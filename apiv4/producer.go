// 300 usuarios creados con faker, cada usuario posee de a 3 datos, se encolan con rabbitmq y el progreso se ve reflejado en consola en un barra de porcentaje.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/streadway/amqp"
)

// User struct to represent a user
type User struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	startTime := time.Now()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"user_queue", // Nombre de la cola
		false,        // Durable
		false,        // AutoDelete
		false,        // Exclusive
		false,        // NoWait
		nil,          // Argumentos adicionales
	)
	failOnError(err, "Failed to declare a queue")

	// Generar 300 usuarios y publicarlos en RabbitMQ
	var lastUserPosition int
	for i := 0; i < 300; i++ {
		user := generateRandomUser()
		body, err := json.Marshal(user)
		failOnError(err, "Failed to marshal user to JSON")

		err = ch.Publish(
			"",     // Exchange
			q.Name, // Nombre de la cola
			false,  // Mandatory
			false,  // Immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		failOnError(err, "Failed to publish a message")

		lastUserPosition = i + 1
		time.Sleep(1 * time.Millisecond) // Simular un retraso mínimo entre usuarios
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Último usuario generado: %d\n", lastUserPosition)
	fmt.Printf("Duración total de la ejecución: %s\n", duration)
}

func generateRandomUser() User {
	var user User
	err := faker.FakeData(&user)
	failOnError(err, "Failed to generate random user")
	user.Address = generateRandomAddress()
	return user
}

func generateRandomAddress() string {
	// Generar una dirección ficticia de manera manual
	// Puedes personalizar esta función según tus necesidades
	street := faker.FirstName()
	city := getRandomCity()
	state := getRandomState()
	zipCode := getRandomZipCode()

	return fmt.Sprintf("%s St, %s, %s %s", street, city, state, zipCode)
}

func getRandomCity() string {
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix"}
	return cities[rand.Intn(len(cities))]
}

func getRandomState() string {
	states := []string{"NY", "CA", "IL", "TX", "AZ"}
	return states[rand.Intn(len(states))]
}

func getRandomZipCode() string {
	// Simplemente generaremos un código postal de cinco dígitos
	return fmt.Sprintf("%05d", rand.Intn(99999))
}
