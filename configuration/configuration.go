package configuration

import (
	"log"
	"os"
	"strconv"
)

// Configuration is a structure of environment settings
type Configuration struct {
	APIVersion,
	DatabaseGateway,
	FunctionsGateway string

	HTTPServer struct {
		Host                 string
		Port                 int
		StaticFilesDirectory string
	}

	SocketServer struct {
		Host string
		Port int
	}
}

var Logger = log.New(os.Stdout, "Configuration: ", log.Lshortfile)

// New a constructor for *Configuration
func New() *Configuration {
	configuration := Configuration{}

	// APIVersion
	ApiVersion := os.Getenv("APIVersion")
	if ApiVersion == "" {
		configuration.APIVersion = "1.0.0"
	} else {
		configuration.APIVersion = ApiVersion
	}

	// DatabaseGateway
	DatabaseGatewayFromEnvironment := os.Getenv("DatabaseGateway")
	if DatabaseGatewayFromEnvironment == "" {
		configuration.DatabaseGateway = "localhost:9080"
	} else {
		configuration.DatabaseGateway = DatabaseGatewayFromEnvironment
	}

	// FunctionsGateway
	FunctionsGatewayFromEnvironment := os.Getenv("FunctionsGateway")
	if FunctionsGatewayFromEnvironment == "" {
		configuration.FunctionsGateway = "localhost:8080"
	} else {
		configuration.FunctionsGateway = FunctionsGatewayFromEnvironment
	}

	// HTTPServer
	HTTPServerHostFromEnvironment := os.Getenv("HTTPServer-Host")
	if HTTPServerHostFromEnvironment == "" {
		configuration.HTTPServer.Host = "localhost"
	} else {
		configuration.HTTPServer.Host = HTTPServerHostFromEnvironment
	}

	HTTPServerPortFromEnvironment := os.Getenv("HTTPServer-Port")
	if HTTPServerPortFromEnvironment == "" {
		configuration.HTTPServer.Port = 80
	} else {
		port, err := strconv.Atoi(HTTPServerPortFromEnvironment)
		if err != nil {
			Logger.Fatal(err)
		}
		configuration.HTTPServer.Port = port
	}

	HTTPServerStaticFilesDirectoryFromEnvironment := os.Getenv("HTTPServer-StaticFilesDirectory")
	if HTTPServerStaticFilesDirectoryFromEnvironment == "" {
		configuration.HTTPServer.StaticFilesDirectory = "build/web"
	} else {
		configuration.HTTPServer.StaticFilesDirectory = HTTPServerStaticFilesDirectoryFromEnvironment
	}

	// SocketServer
	SocketServerHostFromEnvironment := os.Getenv("SocketServer-Host")
	if SocketServerHostFromEnvironment == "" {
		configuration.SocketServer.Host = "localhost"
	} else {
		configuration.SocketServer.Host = SocketServerHostFromEnvironment
	}

	SocketServerPortFromEnvironment := os.Getenv("SocketServer-Port")
	if SocketServerPortFromEnvironment == "" {
		configuration.SocketServer.Port = 81
	} else {
		port, err := strconv.Atoi(SocketServerPortFromEnvironment)
		if err != nil {
			log.Fatal(err)
		}
		configuration.SocketServer.Port = port
	}

	return &configuration
}
