# Proyecto AMQP2 - API FAKER

Este proyecto consta de un productor y un consumidor que utilizan RabbitMQ para procesar usuarios falsos generados aleatoriamente con datos ficticios.

## Requisitos previos

Asegúrate de tener instalados los siguientes requisitos antes de ejecutar el proyecto:

- [Go](https://golang.org/dl/) (versión 1.15 o superior)
- [RabbitMQ](https://www.rabbitmq.com/download.html) (servidor en ejecución)

## Instalación

1. Clona el repositorio:

   ```bash
   - git clone https://github.com/Hertuq0910/rabbitmq-demos.git
   - cd tuproyecto
   
2. Instala dependencias:

   - go get -u github.com/schollz/progressbar/v3
   - go get -u github.com/streadway/amqp

3. Ejecución:

   - go run producer.go
   - go run consumer.go (activa los consumidores que sean necesarios)

> ¡Las contribuciones son bienvenidas! Si encuentras errores o mejoras posibles, por favor, abre un issue o envía un pull request.

