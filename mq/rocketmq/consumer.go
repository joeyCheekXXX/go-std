package rocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/joeyCheek888/go-std/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	ConsumerModelBroadcasting = consumer.BroadCasting
	ConsumerModelClustering   = consumer.Clustering
)

type PushConsumer struct {
	rocketmq.PushConsumer
}

func NewConsumer(conf *Conf) (pc *PushConsumer, err error) {
	rlog.SetLogLevel("error")

	// 创建消费者
	pc = &PushConsumer{}

	pc.PushConsumer, err = rocketmq.NewPushConsumer(
		consumer.WithGroupName(conf.GroupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(conf.NsResolver)),
		consumer.WithConsumerModel(conf.ConsumerModel),
	)
	if err != nil {
		err = errors.Wrapf(err, "new push consumer failed, groupName: %s ,ns: %+v", conf.GroupName, conf.NsResolver)
		return
	}

	log.Logger.Info("启动Rocketmq-consumer", zap.String("group-name", conf.GroupName), zap.Strings("ns-resolver", conf.NsResolver), zap.Any("consumer-model", conf.ConsumerModel))

	return
}

func (pc *PushConsumer) Start() {
	err := pc.PushConsumer.Start()
	if err != nil {
		log.Logger.Error("start push consumer failed", zap.Error(err))
		return
	}

	return
}

func (pc *PushConsumer) Shutdown() (err error) { return pc.PushConsumer.Shutdown() }

func (pc *PushConsumer) Stop() {
	err := pc.PushConsumer.Shutdown()
	if err != nil {
		log.Logger.Error("shutdown push consumer failed", zap.Error(err))
		return
	}
	return
}

type Handler func(msg *primitive.MessageExt) (consumer.ConsumeResult, error)

func (pc *PushConsumer) Subscribe(topic string, handler Handler) {
	err := pc.PushConsumer.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		for _, messageExt := range ext {
			// 处理消息
			consumerResult, err := handler(messageExt)
			if err != nil {
				// 处理失败，重试
				return consumer.ConsumeRetryLater, nil
			}

			// 处理成功，返回消费成功
			return consumerResult, nil
		}

		return consumer.ConsumeRetryLater, nil
	})
	if err != nil {
		log.Logger.Error("subscribe topic failed", zap.String("topic", topic), zap.Error(err))
		return
	}
	return
}
