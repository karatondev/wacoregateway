package messaging

import (
	"context"
	"encoding/json"

	"sync"

	"wacoregateway/internal/provider/amqpx"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	defaultAMQPPublishing = &amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
	}
)

type AMQPPublisherOptions struct {
	Exchange   string
	Mandatory  bool
	Immediate  bool
	Publishing *amqp.Publishing
}

func (o *AMQPPublisherOptions) reset() {
	o.Exchange = ""
	o.Mandatory = false
	o.Immediate = false
	o.Publishing = nil
}

type AMQPPublisherInterface interface {
	Publish(ctx context.Context, queue string, message any, options ...AMQPPublisherOption) error
}

type AMQPPublisherOption func(options *AMQPPublisherOptions)

type AMQPPublisher struct {
	pool *sync.Pool
	conn amqpx.ChannelReader
}

func NewAMQPPublisher(conn amqpx.ChannelReader) AMQPPublisherInterface {

	return &AMQPPublisher{
		pool: &sync.Pool{
			New: func() any {
				return &AMQPPublisherOptions{}
			},
		},
		conn: conn,
	}
}

func (p *AMQPPublisher) Publish(ctx context.Context, queue string, message any, options ...AMQPPublisherOption) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return errors.WithStack(err)
	}
	defer channel.Close()

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	opts := p.pool.Get().(*AMQPPublisherOptions)
	defer func() {
		opts.reset()
		p.pool.Put(opts)
	}()

	for _, option := range options {
		option(opts)
	}
	if opts.Publishing == nil {
		opts.Publishing = defaultAMQPPublishing
	}

	_, err = channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	return channel.PublishWithContext(
		ctx,
		opts.Exchange,
		queue,
		opts.Mandatory,
		opts.Immediate,
		amqp.Publishing{
			DeliveryMode: opts.Publishing.DeliveryMode,
			ContentType:  opts.Publishing.ContentType,
			Body:         body,
		})
}
