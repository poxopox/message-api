package main

import (
	"container/ring"
	"encoding/json"
	"github.com/poxopox/text-buffer/message"
	"github.com/poxopox/text-buffer/redis"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)
var messageStore message.MessageStore


const stringToRead = "This is a string to read"

func clearScreen () {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func createStringRing (items []string) (returnRing *ring.Ring) {
	returnRing = ring.New(len(items))
	for i := 0; i < returnRing.Len(); i++ {
		returnRing.Value = items[i]
		returnRing = returnRing.Next()
	}
	return returnRing
}

func getMessageFromRequestBody (body io.ReadCloser) (msg *message.Message, err error) {
	msg = &message.Message{}
	bodyBytes := make([]byte, 100)
	n, err := body.Read(bodyBytes)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	if n > 0 {
		jsonBytes := bodyBytes[0:n]
		err = json.Unmarshal(jsonBytes, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func HandlePostMessage (writer http.ResponseWriter, request *http.Request, messageStore message.MessageStore) error {
	msg, err := getMessageFromRequestBody(request.Body)
	if err != nil {
		return err
	}
	msg.TimeStamp = time.Now()
	err = messageStore.PutMessage(msg)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("{ \"message\": \"wrote 1 message\" }"))
	return err
}

func HandleGetMessage (writer http.ResponseWriter, request *http.Request) error {

	messages, err := messageStore.ReadMessages("")
	if err != nil {
		return err
	}
	responseBody, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	_, err = writer.Write(responseBody)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	redisCli := redis.NewConnection("127.0.0.1:6379", "", 0)

	defer redisCli.Close()

	status := redisCli.Ping()
	resultString, err := status.Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Status: ", resultString)

	messageStore = redis.NewMessageStore(redisCli)

	http.HandleFunc("/messages", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			err := HandlePostMessage(writer, request, messageStore)
			if err != nil {
				log.Fatal(err)
			}
		}
		if request.Method == "GET" {
			err := HandleGetMessage(writer, request)
			if err != nil {
				log.Fatal(err)
			}
		}
		msgs, err := messageStore.ReadMessages("")
		if err != nil {
			log.Fatal(err)
		}
		for _, msg := range msgs {
			log.Println(msg.Message)
		}
	})

	log.Fatal(http.ListenAndServe(":9090", nil))

}