package rabittmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Logger struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewLogger(amqpURL string) (*Logger, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	err = channel.ExchangeDeclare(
		EXCHANGE,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %s", err)
	}

	return &Logger{
		conn:     conn,
		channel:  channel,
		exchange: EXCHANGE,
	}, nil
}

func (l *Logger) log(level, message string) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2023-07-01 01:02:03")

	err := l.channel.Publish(
		l.exchange,
		level,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("%s - [%s] %s", formattedTime, level, message)),
		},
	)
	if err != nil {
		log.Printf("failed to publish a message: %s", err)
	}
}

func (l *Logger) Error(message string) {
	l.log(ERROR, message)
}

func (l *Logger) Info(message string) {
	l.log(INFO, message)
}

func (l *Logger) Debug(message string) {
	l.log(DEBUG, message)
}

func (l *Logger) Close() {
	_ = l.channel.Close()
	_ = l.conn.Close()
}
