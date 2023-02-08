package queue

import (
	"os"
	"strconv"

	"github.com/wiliehidayat87/rmqp"
)

var (
	Rabbit rmqp.AMQP

	username = os.Getenv("RMQ_USER")
	password = os.Getenv("RMQ_PASS")
	host     = os.Getenv("RMQ_HOST")
	port     = os.Getenv("RMQ_PORT")
)

func init() {
	port, _ := strconv.Atoi(port)

	Rabbit.SetAmqpURL(host, port, username, password)
	Rabbit.SetUpConnectionAmqp()
}
