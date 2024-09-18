package rocketmq

import "github.com/apache/rocketmq-client-go/v2/consumer"

type Conf struct {
	GroupName     string
	NsResolver    []string
	BookerAddr    string
	ConsumerModel consumer.MessageModel
}
