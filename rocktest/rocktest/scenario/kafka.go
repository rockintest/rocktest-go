package scenario

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

type kafkaData struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func (module *Module) Kafka_connect(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	server, err := scenario.GetString(paramsEx, "server", nil)
	if err != nil {
		return err
	}

	group, _ := scenario.GetString(paramsEx, "group", "group")

	offset, _ := scenario.GetString(paramsEx, "offset", "earliest")

	name, _ := scenario.GetString(paramsEx, "name", "default")

	log.Infof("Connect to Kafka broker %s", server)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": server})
	if err != nil {
		return err
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"group.id":          group,
		"auto.offset.reset": offset,
	})
	if err != nil {
		return err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Warnf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Debugf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	scenario.PutStore("kafka."+name, kafkaData{producer: p, consumer: c})
	scenario.PutCleanup("kafka", closeKafka)

	return nil
}

func (module *Module) Kafka_consume(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	name, _ := scenario.GetString(paramsEx, "name", "default")
	data := scenario.GetStore("kafka." + name).(kafkaData)

	topic, err := scenario.GetString(paramsEx, "topic", nil)
	if err != nil {
		return err
	}

	timeout, _ := scenario.GetNumber(paramsEx, "timeout", 60)

	data.consumer.SubscribeTopics([]string{topic}, nil)

	var ret string = "["

	msg, err := data.consumer.ReadMessage(time.Duration(timeout) * time.Second)
	if err == nil {
		log.Debugf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		ret = fmt.Sprintf("%s\"%s\"", ret, msg.Value)
	} else {
		// If timeout : nothing to read more yet
		// The client will automatically try to recover from all errors.
		log.Errorf("Consumer error: %v (%v)\n", err, msg)
		return err
	}

	for {
		msg, err := data.consumer.ReadMessage(time.Duration(1) * time.Second)
		if err == nil {
			log.Debugf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			ret = fmt.Sprintf("%s,\"%s\"", ret, msg.Value)
		} else {
			// If timeout : nothing to read more yet
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				break
			}
			// The client will automatically try to recover from all errors.
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
			return err
		}
	}

	ret = fmt.Sprintf("%s]", ret)
	scenario.PutContextAs(paramsEx, "kafka", "consume", ret)

	return nil
}

func (module *Module) Kafka_check(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	name, _ := scenario.GetString(paramsEx, "name", "default")
	data := scenario.GetStore("kafka." + name).(kafkaData)

	topic, err := scenario.GetString(paramsEx, "topic", nil)
	if err != nil {
		return err
	}

	timeout, _ := scenario.GetNumber(paramsEx, "timeout", 60)

	pattern, err := scenario.GetString(paramsEx, "expect", nil)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	data.consumer.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := data.consumer.ReadMessage(time.Duration(timeout) * time.Second)
		if err == nil {
			log.Debugf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			if re.Match([]byte(string(msg.Value))) {
				log.Debugf("%s ~ %s => YES", string(msg.Value), pattern)
				scenario.PutContextAs(paramsEx, "kafka", "check", string(msg.Value))
				return nil
			} else {
				log.Debugf("%s ~ %s => NO", string(msg.Value), pattern)
			}
		} else {
			// If timeout : nothing to read more yet
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
			return err
		}
	}

}

func (module *Module) Kafka_produce(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	name, _ := scenario.GetString(paramsEx, "name", "default")
	data := scenario.GetStore("kafka." + name).(kafkaData)

	topic, err := scenario.GetString(paramsEx, "topic", nil)
	if err != nil {
		return err
	}

	msg, err := scenario.GetString(paramsEx, "msg", nil)
	if err != nil {
		return err
	}

	data.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          []byte(msg),
	}, nil)

	data.producer.Flush(15 * 1000)

	return nil
}

func closeKafka(scenario *Scenario) error {
	log.Info("Cleanup Kafka")

	for k, v := range scenario.Store {
		if strings.HasPrefix(k, "kafka.") {
			log.Debugf("Closing %s", k)
			data := v.(kafkaData)
			data.producer.Close()
			data.consumer.Close()
		}
	}

	return nil
}
