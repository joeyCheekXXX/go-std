package rocketmq

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/joeyCheek888/go-std/log"
	"go.uber.org/zap"
)

func CreateTopic(conf *Conf, topics ...string) {

	a, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(conf.NsResolver)),
	)
	if err != nil {
		log.Logger.Error("rocker-mq admin new error", zap.Error(err))
		return
	}

	defer a.Close()

	for _, topic := range topics {
		err = a.CreateTopic(
			context.Background(),
			admin.WithBrokerAddrCreate(conf.BookerAddr),
			admin.WithTopicCreate(topic),
		)
		if err != nil {
			log.Logger.Error("rocker-mq create topic error", zap.Error(err))
		}
	}

	return
}
