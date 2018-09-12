package faas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type ParserPrice struct {
	Value    float64
	DateTime time.Time
}

type ParsedProduct struct {
	Name,
	IRI,
	PreviewImageLink string
	Price ParserPrice
}

type FunctionsInterface interface {
	MVideoPageParser(pageForParseIRI string, instruction storage.PageInstruction) []ParsedProduct
	MVideoPagesCountParser(pageForParseIRI string, instruction storage.PageInstruction) int

	//TODO

	//ReadProductsByName(string, string) []storage.Product
	//ReadProductByID(string, string) storage.Product
}

var Logger = log.New(os.Stdout, "FAASFunctions: ", log.Lshortfile)

type Functions struct {
	APIVersion,
	FunctionsGateway,
	DatabaseGateway string
}

func New(APIVersion, FunctionsGateway, DatabaseGateway string) Functions {
	return Functions{APIVersion, FunctionsGateway, DatabaseGateway}
}

func (functions Functions) MVideoPageParser(
	pageForParseIRI string, instruction storage.PageInstruction) (products []ParsedProduct) {

	functionPath := fmt.Sprintf("%v/%v", functions.FunctionsGateway, "mvideo-page-parser")

	body := struct {
		IRI          string
		Instructions storage.PageInstruction
	}{
		IRI:          pageForParseIRI,
		Instructions: instruction}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		Logger.Printf(err.Error())
		return nil
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		Logger.Printf(err.Error())
		return nil
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Logger.Printf(err.Error())
		return nil
	}

	err = json.Unmarshal(decodedResponse, &products)
	if err != nil {
		Logger.Printf(err.Error())
		return nil
	}

	return products
}

func (functions Functions) MVideoPagesCountParser(
	pageForParseIRI string, instruction storage.PageInstruction) (count int) {

	functionPath := fmt.Sprintf("%v/%v", functions.FunctionsGateway, "mvideo-pages-count-parser")

	body := struct {
		IRI          string
		Instructions storage.PageInstruction
	}{
		IRI:          pageForParseIRI,
		Instructions: instruction}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		Logger.Printf(err.Error())
		return count
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		Logger.Printf(err.Error())
		return count
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Logger.Printf(err.Error())
		return count
	}

	err = json.Unmarshal(decodedResponse, &count)
	if err != nil {
		Logger.Printf(err.Error())
		return count
	}

	return count
}
