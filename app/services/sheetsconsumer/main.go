package main

import (
	sheetsmock "collect/app/services/sheetsconsumer/sheets-mock"
	k "collect/foundation/kafka"
	"collect/foundation/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Config struct {
	DOMAIN_KAFKA_BROKER   string
	DOMAIN_KAFKA_TOPIC    string
	SHEETS_CONSUMER_GROUP string
	INTERNAL_TOKEN        string
}

func main() {
	cfg := loadConfiguration("config.json")
	serverErrors := make(chan error, 1)
	r := k.GetKafkaReader(serverErrors, cfg.DOMAIN_KAFKA_BROKER, cfg.DOMAIN_KAFKA_TOPIC, cfg.SHEETS_CONSUMER_GROUP)
	log, err := logger.New("Sheets-consumer")
	if err != nil {
		log.Fatal(err)
	}
	go SheetsConsumer(r, log, cfg.INTERNAL_TOKEN)
	log.Info("sheets consumer started")
	err = <-serverErrors
	log.Fatal(err, "server error")
}

func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
		return Config{}
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func SheetsConsumer(r *kafka.Reader, log *zap.SugaredLogger, token string) {
	for {
		msg, err := r.ReadMessage(context.Background())
		if errors.Is(err, io.EOF) {
			continue
		}
		if err != nil {
			log.Error("kafka reader, could not read message " + err.Error())
			continue
		}
		handleBySource(log, string(msg.Key), msg.Value, token)
	}
}

func handleBySource(log *zap.SugaredLogger, source string, data []byte, token string) {
	switch source {
	case "form:FormCreated":
		var params sheetsmock.Form
		err := json.Unmarshal(data, &params)
		if err != nil {
			log.Errorln(err)
		}
		ques := getQuestions(log, params.ID, token)
		params.CreateCSV(log)
		for _, q := range ques {
			q.UpdateFirstRow(log)
		}
	case "question:QuestionCreated":
		var params sheetsmock.Question
		err := json.Unmarshal(data, &params)
		if err != nil {
			log.Errorln(err)
		}
		params.UpdateFirstRow(log)
	case "response:ResponseCreated":
		var params sheetsmock.Response
		err := json.Unmarshal(data, &params)
		if err != nil {
			log.Errorln(err)
		}
		ques := getQuestions(log, params.FormID, token)
		for _, q := range ques {
			q.UpdateFirstRow(log)
		}
		log.Infof("response recorded")
	case "answer:AnswerCreated":
		var params sheetsmock.AnswerEvent
		err := json.Unmarshal(data, &params)
		if err != nil {
			log.Errorln(err)
		}
		params.UpdateCSV(log)
	}
}

func getQuestions(log *zap.SugaredLogger, formId int64, token string) []sheetsmock.Question {
	client := http.Client{}
	url := fmt.Sprintf("%s/%d/questions", "http://localhost:3000/v1/forms", formId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorln(err)
		return []sheetsmock.Question{}
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + token},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorln(err)
		return []sheetsmock.Question{}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
		return []sheetsmock.Question{}
	}
	var ques []sheetsmock.Question
	err = json.Unmarshal(body, &ques)
	if err != nil {
		log.Fatalln(err)
		return []sheetsmock.Question{}
	}
	return ques
}
