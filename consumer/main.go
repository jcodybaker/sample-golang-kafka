package main

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

func main() {
	log.Println("Starting consumer")
	err := runConsumer()
	if err != nil {
		log.Fatalf("Error consuming messages: %v", err)
	}
	<-make(chan struct{})
}

func runConsumer() error {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	// init config, enable errors and notifications
	config := sarama.NewConfig()
	config.Metadata.Full = true
	config.ClientID = "test-client2"
	config.Producer.Return.Successes = true

	// Kafka SASL configuration
	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	// TLS configuration
	// https://docs.digitalocean.com/products/databases/kafka/how-to/connect/
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(os.Getenv("KAFKA_CA_CERT")))
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = tlsConfig

	brokers := []string{broker}
	consumer, err := sarama.NewConsumerGroup(brokers, "test-group", config)
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
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %d, offset = %d\n", string(message.Value), message.Timestamp, message.Topic, message.Partition, message.Offset)
		session.MarkMessage(message, "")
	}
	return nil
}

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
