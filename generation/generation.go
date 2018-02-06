package main

import (
	"os"

	"github.com/zackslash/quotify"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// slack API 'token'
	token = os.Getenv("SLACK_TOKEN")

	// The Slack ID of 'quote' channel "Read channel" (Note: ID not name)
	channelID = os.Getenv("SLACK_QUOTE_CHANNEL_ID")
)

// Handler handles lambda request
func Handler() (string, error) {
	return quotify.GenerateSlackQuoteDisplay(channelID, token)
}

func main() {
	lambda.Start(Handler)
}
