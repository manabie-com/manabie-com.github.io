package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"examples/models"
)

const (
	streamName     = "student"
	streamSubjects = "student.*"
	subjectName    = "student.Created"
)

func main() {
	// Connect to NATS
	opt := nats.UserInfo("Bob", "123456")
	nc, err := nats.Connect("nats://localhost:4223", opt)
	checkErr(err)

	// Creates JetStreamContext
	js, err := nc.JetStream()
	checkErr(err)

	// Creates stream
	err = createStream(js)
	checkErr(err)

	// Create students by publishing messages
	err = createStudent(js)
	checkErr(err)

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	nc.Drain()
}

// createStudent publishes stream of events
// with subject "student.Created"
func createStudent(js nats.JetStreamContext) error {
	var student models.Student
	for i := 1; i <= 3; i++ {
		student = models.Student{
			StudentID: i,
			ParentID:  "Parent-" + strconv.Itoa(i),
			Status:    "created",
		}
		studentJSON, _ := json.Marshal(student)
		_, err := js.Publish(subjectName, studentJSON)
		if err != nil {
			return err
		}
		log.Printf("Student with StudentID:%d has been published\n", i)
	}
	return nil
}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext) error {
	// Check if the student stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  []string{streamSubjects},
			Retention: nats.LimitsPolicy,
			MaxAge:    time.Minute * 10000,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("check error: ", err)
	}
}
