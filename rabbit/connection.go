package rabbit

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/medfriend/shared-commons-go/util/consul"
	"github.com/rabbitmq/amqp091-go"
)

func ConnRabbitMQ(consulCLient *api.Client) (*amqp091.Connection, error) {

	rabbitInfo, _ := consul.GetKeyValue(consulCLient, "RABBIT")

	var resultRabbitmq map[string]string

	err := json.Unmarshal([]byte(rabbitInfo), &resultRabbitmq)

	if err != nil {
		return nil, err
	}

	s := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		resultRabbitmq["RABBIT_USER"],
		resultRabbitmq["RABBIT_PASSWORD"],
		resultRabbitmq["RABBIT_HOST"],
		resultRabbitmq["RABBIT_PORT"])

	return Connect(s)
}
