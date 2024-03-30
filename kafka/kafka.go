package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type FixedPartitioner struct{}

func (p *FixedPartitioner) Partition(msg *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	return 0, nil
}
func (p *FixedPartitioner) RequiresConsistency() bool {
	return true
}
func NewFixedPartitioner() sarama.PartitionerConstructor {
	return func(topic string) sarama.Partitioner {
		return &FixedPartitioner{}
	}
}

func ProduceWarningSubscriptionEndingMessage(email string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = NewFixedPartitioner()

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Print("failed to create producer ", err)
		return err
	}
	defer producer.Close()
	topic := "SubscriptionEnding"
	msg := fmt.Sprintf(`{"Email":"%s"}`, email)
	message := []byte(msg)
	p, o, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
		Partition: 0,
	})
	fmt.Println("partition ", p, "offset ", o)
	fmt.Println("message sent", msg)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	return nil

}
func SubscribedMessage(email, duration string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = NewFixedPartitioner()

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Print("failed to create producer ", err)
		return err
	}
	defer producer.Close()
	topic := "Subscribed"
	msg := fmt.Sprintf(`{"Email":"%s","Duration":"%s"}`, email, duration)
	message := []byte(msg)
	p, o, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
		Partition: 0,
	})
	fmt.Println("partition ", p, "offset ", o)
	fmt.Println("message sent", msg)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	return nil
}
