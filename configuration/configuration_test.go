package configuration

import (
	"os"
	"testing"
)

func TestConfigurationHasAPIVersion(test *testing.T) {
	defaultValues := New()

	if defaultValues.APIVersion != "1.0.0" {
		test.Fail()
	}

	os.Setenv("APIVersion", "1.0.1")

	notDefaultValues := New()

	if notDefaultValues.APIVersion != "1.0.1" {
		test.Fail()
	}
}

func TestConfigurationHasDatabaseGateway(test *testing.T) {
	defaultValues := New()

	if defaultValues.DatabaseGateway != "localhost:9080" {
		test.Fail()
	}

	updatedEnvironmentValue := "192.168.99.100:9080"
	os.Setenv("DatabaseGateway", updatedEnvironmentValue)

	notDefaultValues := New()

	if notDefaultValues.DatabaseGateway != updatedEnvironmentValue {
		test.Fail()
	}
}

func TestConfigurationHasFunctionsGateway(test *testing.T) {
	defaultValues := New()

	if defaultValues.FunctionsGateway != "localhost:8080" {
		test.Fail()
	}

	updatedEnvironmentValue := "192.168.99.100:8080"
	os.Setenv("FunctionsGateway", updatedEnvironmentValue)

	notDefaultValues := New()

	if notDefaultValues.FunctionsGateway != updatedEnvironmentValue {
		test.Fail()
	}
}

func TestConfigurationHasHTTPServer(test *testing.T) {
	defaultValues := New()

	if defaultValues.HTTPServer.Host != "localhost" {
		test.Fail()
	}

	if defaultValues.HTTPServer.Port != 80 {
		test.Fail()
	}

	if defaultValues.HTTPServer.StaticFilesDirectory != "build/web" {
		test.Fail()
	}

	updatedHTTPServerHostEnvironmentValue := "192.168.99.100"
	os.Setenv("HTTPServer-Host", updatedHTTPServerHostEnvironmentValue)

	updatedHTTPServerPortEnvironmentValue := "82"
	os.Setenv("HTTPServer-Port", updatedHTTPServerPortEnvironmentValue)

	updatedHTTPServerStaticFileDirectoryEnvironmentValue := "web"
	os.Setenv("HTTPServer-StaticFilesDirectory", updatedHTTPServerStaticFileDirectoryEnvironmentValue)

	notDefaultValues := New()

	if notDefaultValues.HTTPServer.Host != updatedHTTPServerHostEnvironmentValue {
		test.Fail()
	}

	if notDefaultValues.HTTPServer.Port != 82 {
		test.Fail()
	}

	if notDefaultValues.HTTPServer.StaticFilesDirectory != updatedHTTPServerStaticFileDirectoryEnvironmentValue {
		test.Fail()
	}
}

func TestConfigurationHasSocketServer(test *testing.T) {
	defaultValues := New()

	if defaultValues.SocketServer.Host != "localhost" {
		test.Fail()
	}

	if defaultValues.SocketServer.Port != 81 {
		test.Fail()
	}

	updatedSocketServerHostEnvironmentValue := "192.168.99.100"
	os.Setenv("SocketServer-Host", updatedSocketServerHostEnvironmentValue)

	updatedSocketServerPortEnvironmentValue := "83"
	os.Setenv("SocketServer-Port", updatedSocketServerPortEnvironmentValue)

	notDefaultValues := New()

	if notDefaultValues.SocketServer.Host != updatedSocketServerHostEnvironmentValue {
		test.Fail()
	}

	if notDefaultValues.SocketServer.Port != 83 {
		test.Fail()
	}
}
