package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type getItemsRequest struct {
	SortBy     string
	SortOrder  string
	ItemsToGet int
}

type getItemsResponseError struct {
	Message string `json:"message"`
}

type getItemsResponseData struct {
	Item string `json:"item"`
}

type getItemsResponseBody struct {
	Result string                 `json:"result"`
	Data   []getItemsResponseData `json:"data"`
	Error  getItemsResponseError  `json:"error"`
}

type getItemsResponseHeaders struct {
	ContentType string `json:"Content-Type"`
}

type getItemsResponse struct {
	StatusCode int                     `json:"statusCode"`
	Headers    getItemsResponseHeaders `json:"headers"`
	Body       getItemsResponseBody    `json:"body"`
}

func main() {
	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String("ap-south-1")})

	// Get the 10 most recent items
	request := getItemsRequest{"time", "descending", 10}

	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling MyGetItemsFunction request")
		os.Exit(0)
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("optimize-source"), Payload: payload})
	log.Println(result)
	log.Println("PAYLOAD:", string(result.Payload))
	if err != nil {
		fmt.Println("Error calling optimize-source")
		os.Exit(0)
	}

	var resp getItemsResponse

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		fmt.Println("Error unmarshalling optimize-source response")
		os.Exit(0)
	}

	// If the status code is NOT 200, the call failed
	if resp.StatusCode != 200 {
		fmt.Println("Error getting items, StatusCode: " + strconv.Itoa(resp.StatusCode))
		os.Exit(0)
	}

	// If the result is failure, we got an error
	if resp.Body.Result == "failure" {
		fmt.Println("Failed to get items")
		os.Exit(0)
	}

	// Print out items
	if len(resp.Body.Data) > 0 {
		for i := range resp.Body.Data {
			fmt.Println(resp.Body.Data[i].Item)
		}
	} else {
		fmt.Println("There were no items")
	}
}
