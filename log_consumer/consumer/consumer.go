package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/streadway/amqp"
)

const (
	host     = "RABBITMQ_HOST"
	port     = "RABBITMQ_PORT"
	user     = "RABBITMQ_USERNAME"
	password = "RABBITMQ_PASSWORD"
	queue    = "RABBITMQ_QUEUE"
	exchange = "RABBITMQ_EXCHANGE"
	logLevel = "RABBITMQ_LOG_LEVEL"
)

type rabbitMQConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	Exchange  string
	QueueName string
	LogLevel  string
}

type LoggerConsumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	queue          amqp.Queue
	rabbitMQConfig rabbitMQConfig
}

func NewLoggerConsumer() (*LoggerConsumer, error) {
	consumer := &LoggerConsumer{
		rabbitMQConfig: rabbitMQConfig{
			Host:      os.Getenv(host),
			Port:      os.Getenv(port),
			Username:  os.Getenv(user),
			Password:  os.Getenv(password),
			Exchange:  os.Getenv(exchange),
			QueueName: os.Getenv(queue),
			LogLevel:  os.Getenv(logLevel),
		},
	}

	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		consumer.rabbitMQConfig.Username,
		consumer.rabbitMQConfig.Password,
		consumer.rabbitMQConfig.Host,
		consumer.rabbitMQConfig.Port,
	))

	if err != nil {
		return nil, err
	}

	consumer.conn = conn
	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	consumer.channel = channel
	err = channel.ExchangeDeclare(
		consumer.rabbitMQConfig.Exchange,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		consumer.rabbitMQConfig.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		conn.Close()
		return nil, err
	}

	consumer.queue = queue

	return consumer, nil
}

func (c *LoggerConsumer) ConsumeLogs() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	err := c.channel.QueueBind(
		c.rabbitMQConfig.QueueName,
		c.rabbitMQConfig.LogLevel,
		c.rabbitMQConfig.Exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	messages, err := c.channel.Consume(
		c.rabbitMQConfig.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for message := range messages {
			log.Printf("Received log (level: %s): %s", c.rabbitMQConfig.LogLevel, message.Body)
		}
	}()

	<-interrupt

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
