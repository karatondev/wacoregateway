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
	ExchangePublish(ctx context.Context, exchange string, kind string, message any, options ...AMQPPublisherOption) error
}

type AMQPPublisherOption func(options *AMQPPublisherOptions)

type AMQPPublisher struct {
	pool    *sync.Pool
	channel *amqpx.Channel
}

func NewAMQPPublisher(conn amqpx.ChannelReader) AMQPPublisherInterface {
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &AMQPPublisher{
		pool: &sync.Pool{
			New: func() any {
				return &AMQPPublisherOptions{}
			},
		},
		channel: channel,
	}
}

func (p *AMQPPublisher) Publish(ctx context.Context, queue string, message any, options ...AMQPPublisherOption) error {
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

	_, err = p.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	return p.channel.PublishWithContext(
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

func (p *AMQPPublisher) ExchangePublish(ctx context.Context, exchange string, kind string, message any, options ...AMQPPublisherOption) error {
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

	err = p.channel.ExchangeDeclare(
		exchange, // name
		kind,     // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return p.channel.PublishWithContext(
		ctx,
		exchange,
		"",
		opts.Mandatory,
		opts.Immediate,
		amqp.Publishing{
			DeliveryMode: opts.Publishing.DeliveryMode,
			ContentType:  opts.Publishing.ContentType,
			Body:         body,
		})
}
