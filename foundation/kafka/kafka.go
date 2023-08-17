package kafka

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func GetKafkaWriter(ch chan error, brokers string, topic string) *kafka.Writer {
	defer func() {
		if err := recover(); err != nil {
			ch <- errors.Wrap(err.(error), "panic occured while writing to kafka "+topic)
		}
	}()

	brkrs := trimEachInList(strings.Split(brokers, ","))
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brkrs,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	})
	return w
}

func GetKafkaReader(ch chan error, brokers string, topic, cg string) *kafka.Reader {
	defer func() {
		if err := recover(); err != nil {
			ch <- errors.Wrap(err.(error), "panic occured while writing to kafka "+topic)
		}
	}()
	brkrs := trimEachInList(strings.Split(brokers, ","))
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brkrs,
		GroupID:  cg,
		Topic:    topic,
		MinBytes: 10e1,
		MaxBytes: 10e6, // 10MB
	})
	return r
}

func trimEachInList(list []string) []string {
	var trimmed []string
	for _, item := range list {
		if strings.TrimSpace(item) != "" {
			trimmed = append(trimmed, strings.TrimSpace(item))
		}
	}
	return trimmed
}
