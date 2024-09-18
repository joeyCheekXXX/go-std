package rocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"go-std/log"
	"go.uber.org/zap"
)

type Producer struct {
	client rocketmq.Producer
}

func NewProducer(conf *Conf) (*Producer, error) {

	client, err := rocketmq.NewProducer(
		producer.WithGroupName(conf.GroupName),
		producer.WithNsResolver(primitive.NewPassthroughResolver(conf.NsResolver)),
		producer.WithRetry(3), // retry 3 times
	)
	if err != nil {
		return nil, err
	}

	log.Logger.Info("启动Rocketmq-producer服务", zap.String("group-name", conf.GroupName), zap.Strings("ns-resolver", conf.NsResolver))

	return &Producer{
		client: client,
	}, nil
}

func (p *Producer) Start() {
	err := p.client.Start()
	if err != nil {
		log.Logger.Error("rocketmq producer start error", zap.String("error", err.Error()))
		return
	}
}

func (p *Producer) Stop() {
	err := p.client.Shutdown()
	if err != nil {
		log.Logger.Error("rocketmq producer stop error", zap.String("error", err.Error()))
		return
	}
}

type sendMode string

const (
	SYNC   sendMode = "SYNC"   // 同步发送
	ASYNC  sendMode = "ASYNC"  // 异步发送 默认
	ONEWAY sendMode = "ONEWAY" // 单向发送
)

type ProducerMessage struct {
	Topic    string
	Body     any
	Tag      string
	SendMode sendMode
	// WithDelayTimeLevel set message delay time to consume.
	// reference delay level definition: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
	// delay level starts from 1. for example, if we set param level=1, then the delay time is 1s.
	DelayTimeLevel int
	AsyncCallback  func(ctx context.Context, result *primitive.SendResult, err error) // 异步发送回调
}

func (p *Producer) Send(ctx context.Context, message *ProducerMessage) (err error) {

	body, err := json.Marshal(message.Body)
	if err != nil {
		return errors.Wrap(err, "producer send json marshal error")
	}

	msg := primitive.NewMessage(message.Topic, body)
	msg.WithTag(message.Tag)
	msg.WithDelayTimeLevel(message.DelayTimeLevel)

	switch message.SendMode {
	case "SYNC":
		_, err = p.client.SendSync(ctx, msg)
		if err != nil {
			log.Logger.Error("rocketmq producer sync send error", zap.String("error", err.Error()))
			return
		}
	case "ONEWAY":
		err = p.client.SendOneWay(ctx, msg)
		if err != nil {
			log.Logger.Error("rocketmq producer oneway send error", zap.String("error", err.Error()))
			return
		}

	default:
		err = p.client.SendAsync(ctx, message.AsyncCallback, msg)
		if err != nil {
			log.Logger.Error("rocketmq producer async send error", zap.String("error", err.Error()))
			return
		}
	}

	return
}
