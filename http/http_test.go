package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/hecatoncheir/Core/configuration"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpServer() {
	server := New("v1")
	goroutines.Done()
	config := configuration.New()

	err := server.SetUp("", config.HTTPServer.Host, config.HTTPServer.Port)
	if err != nil {
		server.log.Printf("Faild SetUp HTTP server with error: %v", err)
	}
}

func TestHttpServerCanSendVersionOfAPI(test *testing.T) {

	os.Setenv("HTTPServer-Port", "3000")

	goroutines.Add(1)
	go once.Do(SetUpServer)
	goroutines.Wait()

	config := configuration.New()

	iri := fmt.Sprintf("http://%v:%v/api/version", config.HTTPServer.Host, config.HTTPServer.Port)
	response, err := http.Get(iri)
	if err != nil {
		test.Fatal(err)
	}

	encodedBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		test.Fatal(err)
	}

	decodedBody := map[string]string{}

	err = json.Unmarshal(encodedBody, &decodedBody)
	if err != nil {
		test.Fatal(err)
	}

	if decodedBody["apiVersion"] != "v1" {
		test.Fatalf("The api version should be 'v1'.")
	}
}
