package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"order-notification/internal/entity"
	"order-notification/internal/service"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.OrderService
}

func NewConsumer(brokers []string, topic string, groupID string, service *service.OrderService) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{
		reader:  r,
		service: service,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if ctx.Err() != nil {
				logrus.Errorf("Kafka read error: %v", err)
				return
			}

			if err != nil {
				logrus.Errorf("read error: %v", err)
				continue
			}

			var order entity.Order
			logrus.Debugf("Raw message: %+v", msg)
			err = json.Unmarshal(msg.Value, &order)
			if err != nil {
				logrus.Errorf("unmarshal error: %v", err)
				continue
			}

			logrus.Debugf("Order received:%s\n", order.OrderUID)

			err = c.service.SaveOrder(order, ctx)
			if err != nil {
				logrus.Errorf("Process order error: %v", err)
				continue
			}
		}
	}
}

/* Запуск ZOOKEPER и Kafka
bin/zookeeper-server-start.sh config/zookeeper.properties
bin/kafka-server-start.sh config/server.properties

*/
