package main

import (
	"encoding/json"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ryananyangu/gojsrunner/controllers"
	"github.com/ryananyangu/gojsrunner/models"
	"github.com/ryananyangu/gojsrunner/utils"
)

func main() {
	utils.Log.Info(os.Getenv("AMQP_SERVER_URL"))
	connectRabbitMQ, err := amqp.Dial("amqp://myuser:mypassword@172.17.0.1:5672/")
	if err != nil {
		utils.Log.Error(err)
		// return err
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		utils.Log.Error(err)
		// return err
	}
	defer channelRabbitMQ.Close()

	// Have deadletter Q for escalations
	msgs, err := channelRabbitMQ.Consume(
		"paymentsTpQ", // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)

	utils.Log.Info(err)

	var forever chan struct{}

	for d := range msgs {
		go func(d amqp.Delivery) {
			request := models.Request{}
			if err := json.Unmarshal(d.Body, &request); err != nil {
				// Will not have acknowledge hence messages will still be on queue
				utils.Log.Error(err)
			}
			controllers.RequestTransformation(&request)
			d.Ack(false)
			// d.Acknowledger.Ack(true)
		}(d)
	}

	utils.Log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
