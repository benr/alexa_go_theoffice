package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/adafruit/io-client-go"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
)

// AlexaRequest The Go Struct representing the Alexa JSON Request Payload
type AlexaRequest struct {
	Version string `json:"version"`
	Request struct {
		Type   string `json:"type"`
		Time   string `json:"timestamp"`
		Intent struct {
			Name               string `json:"name"`
			ConfirmationStatus string `json:"confirmationstatus"`
		} `json:"intent"`
	} `json:"request"`
}

// AlexaResponse The Go Struct representing the Alexa JSON Response Payload
type AlexaResponse struct {
	Version  string `json:"version"`
	Response struct {
		OutputSpeech struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"outputSpeech"`
	} `json:"response"`
}

// CreateResponse A Constructor for an AlexaResponse Object
func CreateResponse() *AlexaResponse {
	var resp AlexaResponse
	resp.Version = "1.0"
	resp.Response.OutputSpeech.Type = "PlainText"
	resp.Response.OutputSpeech.Text = "Hello.  Please override this default output."
	return &resp
}

// Say A convience function for the AlexaReponse Object
func (resp *AlexaResponse) Say(text string) {
	resp.Response.OutputSpeech.Text = text
}

// GetTemp Get Office Temp form Adafruit.IO
func GetTemp() string {
	fmt.Println("Starting...")
	aio := adafruitio.NewClient(os.Getenv("ADAFRUIT_IO_KEY"))
	feeds, resp, _ := aio.Feed.All()
	fmt.Printf("Feeds: %v\nResponse: %v\n", feeds, resp.StatusCode)
	if resp.StatusCode == 200 {
		for i, d := range feeds {
			fmt.Printf("(%d) %v %v %v\n", i, d.Name, d.LastValue, d.UnitType)
			if d.Name == "office_temp" {
				t, _ := strconv.ParseFloat(d.LastValue, 32)
				return fmt.Sprintf("%.1f", t)
			}
		}
	}
	return "blarf!"
}

// HandleRequest The Lambda Handler
func HandleRequest(ctx context.Context, i AlexaRequest) (AlexaResponse, error) {
	// Use Spew to output the request for debugging purposes:
	fmt.Println("---- Dumping Input Map: ----")
	spew.Dump(i)
	fmt.Println("---- Done. ----")

	// Example of accessing map value via index:
	log.Printf("Request type is %s", i.Request.Intent.Name)

	// Create a response object
	resp := CreateResponse()

	// Customize the response for each Alexa Intent
	switch i.Request.Intent.Name {
	case "officetemp":
		temp := GetTemp()
		s := fmt.Sprintf("The current temperature is %s degrees.", temp)
		//resp.Say("The current temperature is 68 degrees.")
		resp.Say(s)
	case "hello":
		resp.Say("Hello there, Lambda appears to be working properly.")
	case "AMAZON.HelpIntent":
		resp.Say("This app is easy to use, just say: ask the office how warm it is")
	default:
		resp.Say("I'm sorry, the input does not look like something I understand.")
	}

	return *resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
