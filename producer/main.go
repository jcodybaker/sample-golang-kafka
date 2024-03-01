package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func logRequest(r *http.Request) {
	uri := r.RequestURI
	method := r.Method
	log.Println("Got request!", method, uri)
}

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	conf := sarama.NewConfig()
	conf.Metadata.Full = true
	conf.ClientID = "test-client"
	// conf.Version = sarama.V3_6_0_0
	conf.Producer.Return.Successes = true

	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = username
	conf.Net.SASL.Password = password
	conf.Net.SASL.Handshake = true
	conf.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	// https://docs.digitalocean.com/products/databases/kafka/how-to/connect/
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(os.Getenv("KAFKA_CA_CERT")))
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = tlsConfig

	brokers := []string{broker}
	producer, err := sarama.NewSyncProducer(brokers, conf)
	if err != nil {
		log.Panicf("Error creating producer: %v", err)
	}
	defer producer.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)
	})

	http.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder("hello"),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		fmt.Fprintf(w, "Message produced in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	bindAddr := fmt.Sprintf(":%s", port)
	fmt.Printf("==> Server listening at %s 🚀\n", bindAddr)
	if err := http.ListenAndServe(bindAddr, nil); err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}
