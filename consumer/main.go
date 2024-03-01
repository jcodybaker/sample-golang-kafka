package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/IBM/sarama"
)

func main() {
	log.Println("Starting consumer")
	err := runConsumer()
	if err != nil {
		log.Fatalf("Error consuming messages: %v", err)
	}
}

func runConsumer() error {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	// https://docs.digitalocean.com/products/databases/kafka/how-to/connect/
	config := sarama.NewConfig()
	config.Metadata.Full = true
	config.ClientID = "sample-consumer-client"
	config.Producer.Return.Successes = true

	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(os.Getenv("KAFKA_CA_CERT")))
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = tlsConfig

	brokers := []string{broker}
	consumer, err := sarama.NewConsumerGroup(brokers, "sample-group", config)
	if err != nil {
		return err
	}
	defer consumer.Close()
	kafkaConsumer := &KafkaConsumer{}
	err = consumer.Consume(context.Background(), []string{topic}, kafkaConsumer)
	if err != nil {
		return err
	}
	return nil
}

type KafkaConsumer struct{}

func (consumer *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message consumed: value = %s, timestamp = %v, topic = %s, partition = %d, offset = %d\n", string(message.Value), message.Timestamp, message.Topic, message.Partition, message.Offset)
		session.MarkMessage(message, "")
	}
	return nil
}
