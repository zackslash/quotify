package main

import (
	"os"

	"github.com/zackslash/quotify"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// Name of the slack ID of the slack channel to post in
	channel = os.Getenv("SLACK_CHANNEL")

	// Slack webhook URL
	slackWebhook = os.Getenv("SLACK_WEBHOOK_URL")

	// Render endpoint (Using Fortifi URLGrab)
	imageGenEndpoint = os.Getenv("IMAGE_GEN_ENDPOINT")
)

// Handler handles lambda request
func Handler() {
	quotify.DeliverInspiration(slackWebhook, imageGenEndpoint, channel)
}

func main() {
	lambda.Start(Handler)
}
