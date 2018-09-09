package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hecatoncheir/Core/faas"
	"log"
	"net/http"
	"os"
	"sync"
)

// Server is an object of socket server structure
type Server struct {
	APIVersion string
	HTTPServer *http.Server
	Log        *log.Logger
	Clients    map[string]*Client

	functions    faas.FunctionsInterface
	clientsMutex sync.Mutex
	headersUp    websocket.Upgrader
}

// New is constructor for socket server
func New(apiVersion string, faasFunctions faas.FunctionsInterface) *Server {
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	socketServer := Server{
		APIVersion: apiVersion,
		functions:  faasFunctions,
		Clients:    make(map[string]*Client),
		headersUp:  up}

	logPrefix := fmt.Sprintf("SocketServer ")
	socketServer.Log = log.New(os.Stdout, logPrefix, 3)

	return &socketServer
}

// SetUp is a method for listen server on port
func (server *Server) SetUp(host string, port int) error {
	server.HTTPServer = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	server.HTTPServer.Handler = http.HandlerFunc(server.ClientConnectedHandler)

	eventMessage := fmt.Sprintf("Socket server listen on %v, port:%v \n", host, port)
	server.Log.Println(eventMessage)

	err := server.HTTPServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// ClientConnectedHandler handler for connected client
func (server *Server) ClientConnectedHandler(response http.ResponseWriter, request *http.Request) {
	socketConnection, err := server.headersUp.Upgrade(response, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewConnectedClient(socketConnection)

	server.clientsMutex.Lock()
	server.Clients[client.ID] = client
	server.clientsMutex.Unlock()

	eventMessage := fmt.Sprintf("Client: %v connected. Connected clients: %v", client.ID, len(server.Clients))
	server.Log.Println(eventMessage)

	go server.listenConnectedClient(client)
}

// listenConnectedClient need for receive and broadcast client messages
func (server *Server) listenConnectedClient(client *Client) {
	for event := range client.Channel {
		event.ClientID = client.ID

		eventMessage := fmt.Sprintf("Received event: %v from connected client: %v", event, client.ID)

		server.Log.Println(eventMessage)

		switch event.Message {
		case "Need api version":
			server.Clients[event.ClientID].Write("Version of API", server.APIVersion, "")

		case "Need items by name":
			//TODO call FAASFunctions
			//eventData := broker.EventData{Message: event.Message, Data: event.Data}
			//eventData.ClientID = client.ID
			//server.Broker.Write(eventData)
		}
	}

	server.clientsMutex.Lock()
	delete(server.Clients, client.ID)
	server.clientsMutex.Unlock()
}

// WriteToAll send events to all connected clients
func (server *Server) WriteToAll(message string, data string) {
	for _, connection := range server.Clients {
		go connection.Write(message, server.APIVersion, data)
	}
}

// WriteToClient send events to all connected clients
func (server *Server) WriteToClient(clientID, message, APIVersion, data string) {
	for _, connection := range server.Clients {
		if connection.ID == clientID {

			eventMessage := fmt.Sprintf("Writing message: %v to connected client: %v", message, clientID)
			server.Log.Println(eventMessage)

			go connection.Write(message, server.APIVersion, data)
			break
		}
	}
}
