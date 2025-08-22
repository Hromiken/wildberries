package test

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func Produser() {
	broker := "localhost:9092"
	topic := "orders"

	writer := kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()
	orderJSON := `{
		"order_uid": "test-001",
		"track_number": "WBTEST001",
		"entry": "WBIL",
		"delivery": {
			"name": "Test User",
			"phone": "+79991234567",
			"zip": "123456",
			"city": "TestCity",
			"address": "Test Street 1",
			"region": "Test Region",
			"email": "test@example.com"
		},
		"payment": {
			"transaction": "test-001",
			"currency": "RUB",
			"provider": "wbpay",
			"amount": 1500,
			"payment_dt": 1637907727,
			"bank": "alpha",
			"delivery_cost": 300,
			"goods_total": 1200,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 123,
				"track_number": "WBTEST001",
				"price": 1200,
				"rid": "12345",
				"name": "Test Product",
				"sale": 0,
				"size": "M",
				"total_price": 1200,
				"nm_id": 111,
				"brand": "TestBrand",
				"status": 202
			}
		],
		"locale": "ru",
		"internal_signature": "",
		"customer_id": "test-user",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
	}`
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test"),
		Value: []byte(orderJSON),
	})
	if err != nil {
		log.Fatal("Ошибка отправки в Kafka:", err)
	}

	fmt.Println("Заказ test-001 отправлен в Kafka → topic:", topic)

	time.Sleep(1 * time.Second)
}
