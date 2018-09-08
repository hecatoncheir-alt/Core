package socket

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// EventData structure for communication with client
type EventData struct {
	ClientID,
	Message,
	APIVersion,
	Data string
}

// Client is a structure of connected client object
type Client struct {
	ID,
	APIVersion string
	Connection *websocket.Conn
	wmu        sync.Mutex
	Log        *log.Logger
	Channel    chan EventData
}

// NewConnectedClient for constructor for Client
func NewConnectedClient(clientConnection *websocket.Conn) *Client {
	clientID, _ := uuid.NewUUID()
	client := Client{
		ID:         clientID.String(),
		Channel:    make(chan EventData),
		Connection: clientConnection}

	client.Log = log.New(os.Stdout, "Connected client: ", 3)

	go func() {
		for {

			inputMessage := EventData{}
			_, messageBytes, err := clientConnection.ReadMessage()

			if err != nil {
				client.Log.Printf("Can't receive message from %s. %v \n", client.ID, err)
				client.Log.Printf("Closed connection of client %s \n", client.ID)
				client.Connection.Close()
				break
			}

			err = json.Unmarshal(messageBytes, &inputMessage)
			if err != nil {
				client.Log.Printf("Fail unmarshal event: %v", err)
			}

			inputMessage.ClientID = client.ID
			client.Channel <- inputMessage
		}
	}()

	return &client
}

// Write need for send event to client
func (client *Client) Write(message, APIVersion, data string) {

	event := EventData{
		ClientID:   client.ID,
		Message:    message,
		APIVersion: APIVersion,
		Data:       data}

	client.wmu.Lock()
	err := client.Connection.WriteJSON(event)
	if err != nil {
		client.Log.Printf("Fail write event: %v", err)
	}

	client.wmu.Unlock()
}
