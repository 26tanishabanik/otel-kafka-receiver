package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

type Netrounds struct {
	ResourceName    string                            `json:"resourceName"`
	MetricAttribute string                            `json:"metricAttribute"`
	Values          map[string]map[string]interface{} `json:"values"`
	Timestamp       string                            `json:"timestamp"`
}

const (
	BootstrapServersDefault = "127.0.0.1:9094"
	TopicDefault            = "customdata"
	DelayDefault            = 1000
	MessageDefault          = "Hello from Go Kafka Sarama test-2"
	MessageCountDefault     = 10
	ProducerAcksDefault     = int16(1)
)

type ProducerConfig struct {
	BootstrapServers string
	Topic            string
	Delay            int
	Message          string
	MessageCount     int64
	ProducerAcks     int16
}

func NewProducerConfig() *ProducerConfig {
	config := ProducerConfig{
		BootstrapServers: BootstrapServersDefault,
		Topic:            TopicDefault,
		Delay:            DelayDefault,
		Message:          MessageDefault,
		MessageCount:     MessageCountDefault,
		ProducerAcks:     ProducerAcksDefault,
	}
	return &config
}

func main() {

	config := NewProducerConfig()
	log.Printf("Go producer starting with config=%+v\n", config)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.ProducerAcks)
	producerConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{config.BootstrapServers}, producerConfig)
	if err != nil {
		log.Printf("Error creating the Sarama sync producer: %v", err)
		os.Exit(1)
	}

	end := make(chan int, 1)

	values := map[string]map[string]interface{}{
		"cpu_usage_percent": {
			"floatVal": 10.3,
		},
		"memory_usage_bytes": {
			"floatVal": 8.73,
		},
	}
	cmt := Netrounds{
		ResourceName:    "node-1",
		MetricAttribute: "data_centre",
		Values:          values,
		Timestamp:       "2023-07-31T16:34:42Z",
	}
	value, err := json.Marshal(cmt)
	if err != nil {
		fmt.Println("error in marshalling: ", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: config.Topic,
		Value: sarama.StringEncoder(value),
	}
	log.Printf("Sending message: value=%s\n", msg.Value)
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Errors sending message: %v\n", err)
	} else {
		log.Printf("Message sent: partition=%d, offset=%d\n", partition, offset)
	}
	select {
	case <-end:
		log.Printf("Finished sending %d messages\n", config.MessageCount)
	case sig := <-signals:
		log.Printf("Got signal: %v\n", sig)
	}

	err = producer.Close()
	if err != nil {
		log.Printf("Error closing the Sarama sync producer: %v", err)
		os.Exit(1)
	}
	log.Printf("Producer closed")

}
