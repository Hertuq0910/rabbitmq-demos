package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
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

	// Consumir mensajes de la cola
	msgs, err := ch.Consume(
		q.Name, // Nombre de la cola
		"",     // Consumidor
		true,   // AutoAck (confirmación automática de mensajes)
		false,  // Exclusive
		false,  // NoLocal
		false,  // NoWait
		nil,    // Argumentos adicionales
	)
	failOnError(err, "Failed to register a consumer")

	// Canal para manejar señales de interrupción
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Canal para comunicar la finalización de la tarea
	taskCompleted := make(chan struct{})

	startTime := time.Now()

	// Crear una barra de progreso
	bar := progressbar.NewOptions(
		300,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Processing"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionOnCompletion(func() {
			close(taskCompleted)
		}),
	)

	// Procesar mensajes en segundo plano
	go func() {
		for {
			select {
			case msg := <-msgs:
				handleMessage(msg.Body)
				// Actualizar la barra de progreso
				bar.Add(1)
			}
		}
	}()

	// Esperar la señal de interrupción o la finalización de la tarea
	select {
	case <-interrupt:
		fmt.Println("Received interrupt signal. Exiting...")
	case <-taskCompleted:
		// Tarea completada, imprimir el tiempo total y salir
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		fmt.Printf("\nTotal processing time: %s\n", duration)
	}

	// Cerrar el programa
	os.Exit(0)
}

func handleMessage(body []byte) {
	var user User
	err := json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Agrega tu lógica de procesamiento aquí, por ejemplo, almacenar en una base de datos, etc.
	time.Sleep(1 * time.Second) // Simulando un procesamiento que lleva tiempo
}
