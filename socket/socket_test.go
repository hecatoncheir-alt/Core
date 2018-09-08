package socket

import (
	"os"
	"sync"
	"testing"

	"fmt"

	"github.com/hecatoncheir/Core/configuration"
	"golang.org/x/net/websocket"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpSocketServer() {
	testServer := New("v1.0")
	goroutines.Done()

	config := configuration.New()
	err := testServer.SetUp(config.SocketServer.Host, config.SocketServer.Port)
	if err != nil {
		fmt.Println("SetUpSocketServer fail with: ", err)
	}

	defer testServer.HTTPServer.Close()
}

func TestSocketServerCanHandleEvents(test *testing.T) {

	os.Setenv("SocketServer-Port", "3001")
	os.Setenv("HTTPServer-Port", "3000")

	goroutines.Add(1)
	go once.Do(SetUpSocketServer)
	goroutines.Wait()

	config := configuration.New()

	iriOfWebSocketServer := fmt.Sprintf("ws://%v:%v", config.SocketServer.Host,
		config.SocketServer.Port)

	iriOfHTTPServer := fmt.Sprintf("http://%v:%v", config.HTTPServer.Host,
		config.HTTPServer.Port)

	socketConnection, err := websocket.Dial(iriOfWebSocketServer, "", iriOfHTTPServer)
	if err != nil {
		test.Error(err)
	}

	inputMessage := make(chan EventData)

	go func() {
		defer socketConnection.Close()
		defer close(inputMessage)

		for {
			messageFromServer := EventData{}
			err = websocket.JSON.Receive(socketConnection, &messageFromServer)
			if err != nil {
				test.Error(err)
				break
			}

			inputMessage <- messageFromServer
		}
	}()

	messageToServer := EventData{Message: "Need api version"}
	err = websocket.JSON.Send(socketConnection, messageToServer)

	if err != nil {
		test.Error(err)
	}

	for messageFromServer := range inputMessage {
		if messageFromServer.Message == "Version of API" {
			break
		}

		if messageFromServer.Message != "Version of API" {
			test.Fail()
		}
	}
}
