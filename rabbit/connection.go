package rabbit

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/medfriend/shared-commons-go/util/consul"
	"github.com/rabbitmq/amqp091-go"
	"strconv"
	"strings"
)

func ConnRabbitMQ(consulCLient *api.Client) (*amqp091.Connection, map[string]int) {

	rabbitInfo, _ := consul.GetKeyValue(consulCLient, "RABBIT")

	var resultRabbitmq map[string]string

	err := json.Unmarshal([]byte(rabbitInfo), &resultRabbitmq)

	if err != nil {
		return nil, nil
	}

	s := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		resultRabbitmq["RABBIT_USER"],
		resultRabbitmq["RABBIT_PASSWORD"],
		resultRabbitmq["RABBIT_HOST"],
		resultRabbitmq["RABBIT_PORT"])

	services := strings.Split(resultRabbitmq["RABBIT_SERVICES"], ",")

	serviceZero := map[string]int{}

	for _, service := range services {
		serviceInfo, _ := consul.GetKeyValue(consulCLient, strings.Trim(service, " "))
		var resultServiceInfo map[string]string

		err := json.Unmarshal([]byte(serviceInfo), &resultServiceInfo)
		if err != nil {
		}

		port, _ := strconv.Atoi(resultServiceInfo["SERVICE_PORT"])
		serviceZero[strings.ToLower(service)] = port

	}

	conn, err := Connect(s)

	if err != nil {
	}

	return conn, serviceZero
}
