//API V1 (ERROR EN TOKEN)
/*package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

var receiverCount int

type Receiver struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	queue         amqp.Queue
	lastMessage   string
	lastMessageMu sync.Mutex
}

type Sender struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type Notification struct {
	ID      string `json:"ID"`
	Message string `json:"message"`
}

func NewReceiver(rabbitMQURL, queueName string) *Receiver {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	receiverCount++

	return &Receiver{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}
}

func (r *Receiver) processMessage(body []byte) {
	var notification Notification
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Printf("Error al decodificar el mensaje: %s", err)
		return
	}

	id := notification.ID
	switch id {
	case "1":
		fmt.Println("Procesando mensaje con ID 1:", notification.Message)
		// Lógica para procesar mensajes con ID 1
	case "2":
		fmt.Println("Procesando mensaje con ID 2:", notification.Message)
		// Lógica para procesar mensajes con ID 2
	default:
		fmt.Printf("ID no reconocido: %s\n", id)
		return // Ignorar mensajes con ID no reconocido
	}
}

func (r *Receiver) Run() error {
	defer r.conn.Close()
	defer r.channel.Close()

	messages, err := r.channel.Consume(
		r.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %s", err)
	}

	for msg := range messages {
		r.lastMessageMu.Lock()
		r.lastMessage = string(msg.Body)
		r.lastMessageMu.Unlock()

		r.processMessage(msg.Body)

		if err := msg.Ack(false); err != nil {
			return fmt.Errorf("error al confirmar el mensaje: %s", err)
		}
	}

	return nil
}

func (r *Receiver) getLastMessage() (string, error) {
	r.lastMessageMu.Lock()
	defer r.lastMessageMu.Unlock()

	if r.lastMessage == "" {
		return "No messages received", nil
	}
	return fmt.Sprintf("Last message from queue %s: %s", r.queue.Name, r.lastMessage), nil
}

func (r *Receiver) StatusHandler(w http.ResponseWriter, _ *http.Request) {
	lastMessage, err := r.getLastMessage()
	if err != nil {
		http.Error(w, "Failed to get last message", http.StatusInternalServerError)
		return
	}

	if lastMessage == "No messages received" {
		http.Error(w, "No messages received", http.StatusNotFound)
		return
	}

	response := fmt.Sprintf("OK\nLast Message: %s", lastMessage)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func (s *Sender) SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(notification)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	err = s.channel.Publish(
		"",
		s.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %s", err)
		http.Error(w, "Failed to publish a message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func NewSender(rabbitMQURL, queueName string) *Sender {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	return &Sender{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}
}

func main() {
	rabbitMQURL := "amqp://guest:guest@localhost:5672/"
	queueName := "notification_queue"

	var wg sync.WaitGroup

	// Crear e iniciar el receptor
	receiver := NewReceiver(rabbitMQURL, queueName)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := receiver.Run(); err != nil {
			fmt.Printf("Error en el receptor: %s\n", err)
		}
	}()

	// Crear e iniciar el emisor
	sender := NewSender(rabbitMQURL, queueName)

	// Crear y configurar el enrutador HTTP
	router := mux.NewRouter()
	router.HandleFunc("/send", sender.SendNotificationHandler).Methods("POST")
	router.HandleFunc("/status", receiver.StatusHandler).Methods("GET")

	// Configurar y ejecutar el servidor HTTP
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error en el servidor HTTP: %s\n", err)
		}
	}()

	// Esperar a que todas las goroutines (receptores) terminen antes de salir
	wg.Wait()

	// Imprimir el número total de receptores creados
	fmt.Printf("Total de receptores creados: %d\n", receiverCount)
}
*/