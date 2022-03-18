package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"examples/models"
)

func main() {
	// Connect to NATS
	opt := nats.UserInfo("Tom", "123456")
	nc, err := nats.Connect("nats://localhost:4223", opt)
	if err != nil {
		log.Fatal(err)
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Create durable consumer monitor
	sub, _ := js.QueueSubscribe("student.Created", "queue-push", func(msg *nats.Msg) {
		err := msg.Ack()
		if err != nil {
			log.Fatal(err)
		}
		var student models.Student
		err = json.Unmarshal(msg.Data, &student)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Student with StudentID:%d has been processed\n", student.StudentID)
	}, nats.Durable("durable-push"), nats.ManualAck(), nats.MaxDeliver(2), nats.AckWait(time.Second))

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	sub.Unsubscribe()
	nc.Drain()
}
