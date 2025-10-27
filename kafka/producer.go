package kafka

import (
	"api-gateway-module/config"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	_allAcks = "all"
)

type Producer struct {
	cfg      config.Producer
	Producer *kafka.Producer
}

func NewProducer(cfg config.Producer) Producer {
	url := cfg.URL
	id := cfg.ClientID
	acks := cfg.Acks
	if acks == "" {
		acks = _allAcks
	}
	conf := &kafka.ConfigMap{
		"bootstrap.servers": url,  // kafka broker url
		"client.id":         id,   // produce client id
		"acks":              acks, // 0, 1, all
	}

	producer, err := kafka.NewProducer(conf)
	// 모듈 자체를 사용할 수 없는 상황이라 panic 처리
	if err != nil {
		panic(err.Error())
	}
	return Producer{
		cfg:      cfg,
		Producer: producer,
	}
}

// SendEvent 외부에서 global하게 사용하기 위해 byte로 받음.
// 직렬화 / 역직렬화를 처리
func (p Producer) SendEvent(v []byte) {
	topic := p.cfg.Topic

	err := p.Producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: v,
		}, nil)
	if err != nil {
		log.Println("Failed to produce message", string(v))
	} else {
		log.Println("Success to produce message", string(v))
	}
}
