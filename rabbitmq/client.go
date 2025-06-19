package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

// Client - структура клиента RabbitMQ
type Client struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   string
}

// Config - конфигурация подключения
type Config struct {
    URL       string
    QueueName string
}

// New создает новый клиент RabbitMQ
func New(cfg Config) (*Client, error) {
    conn, err := amqp.Dial(cfg.URL)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    _, err = ch.QueueDeclare(
        cfg.QueueName,
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        conn.Close()
        return nil, err
    }

    return &Client{
        conn:    conn,
        channel: ch,
        queue:   cfg.QueueName,
    }, nil
}

// Publish публикует сообщение в очередь
func (c *Client) Publish(body []byte) error {
    return c.channel.Publish(
        "",     // exchange
        c.queue, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        body,
        },
    )
}

// Consume подписывается на сообщения из очереди
func (c *Client) Consume(handler func(body []byte)) error {
    msgs, err := c.channel.Consume(
        c.queue, // queue
        "",      // consumer
        true,    // auto-ack
        false,   // exclusive
        false,   // no-local
        false,   // no-wait
        nil,     // args
    )
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            handler(msg.Body)
        }
    }()

    return nil
}

// Close закрывает соединение с RabbitMQ
func (c *Client) Close() error {
    if err := c.channel.Close(); err != nil {
        return err
    }
    if err := c.conn.Close(); err != nil {
        return err
    }
    return nil
}