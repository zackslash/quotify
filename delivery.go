package quotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
)

const (
	displayName = "'Quote Fortifi' Daily Inspiration BOT"
	displayIcon = ":party_parrot:"
)

// DeliverInspiration adds new inspiration to team slack channel
// it also preloads the next image in the chain ready for
// consecutive execution of this process
func DeliverInspiration(slackWebhook, generateEndpoint, slackChannel, inEntropy string) {
	var entropy []string
	err := json.Unmarshal([]byte(inEntropy), &entropy)
	if err != nil {
		fmt.Println(time.Now().UTC(), "failed to parse entropy")
	}

	now, next := getURLS(entropy, generateEndpoint)

	// Make sure current image is in cache
	http.Get(now)
	fmt.Println(time.Now().UTC(), "pre-loaded current inspiration")

	// Preload the 'next' image so that it is cached in URL grab or refreshes
	// old cache using the same entropy;
	http.Get(next)
	fmt.Println(time.Now().UTC(), "pre-loaded next inspiration")

	// Post the 'now' image to channel
	payload := slack.Payload{
		Username:  displayName,
		Channel:   slackChannel,
		IconEmoji: displayIcon,
		Attachments: []slack.Attachment{
			slack.Attachment{
				ImageUrl: &now,
			},
		},
	}
	slack.Send(slackWebhook, "", payload)
	fmt.Println(time.Now().UTC(), "posted team inspiration")
}

// Deterministically retrieve current and next image url using entropy provided
func getURLS(entropy []string, generateEndpoint string) (current, next string) {
	// random seed used for image on each day
	td := time.Now().UTC().Day()
	current = generateEndpoint
	next = generateEndpoint
	ci := 0

	for i := 0; i <= td; i++ {
		lastIndex := len(entropy) - 1
		ci++
		if ci > lastIndex {
			ci = 0
		}

		if i == td {
			currentSuff, nextSuff := getSuffixForPeriods()

			// Current url with current rotation code
			current = current + entropy[ci] + currentSuff

			if nextSuff != "" {
				// Next url with current rotation code + suffix
				// (image was generated for the same day as the prior image)
				next = current + entropy[ci] + nextSuff
			} else {
				// Next url with tomorrow's rotation code
				if ci+1 > lastIndex {
					ci = 0
				} else {
					ci++
				}
				next = next + entropy[ci]
			}
		}
	}
	return
}

// getSuffixForPeriods gets the link suffix for
// current period of the day (since two images per day)
// and the period following the current one
func getSuffixForPeriods() (now, next string) {
	const (
		amSuffix   = ""
		pmSuffix   = "B"
		stdPMTempl = "PM"
	)
	now = time.Now().UTC().Format(stdPMTempl)
	if now == stdPMTempl {
		now = pmSuffix
		next = amSuffix
		return
	}

	now = amSuffix
	next = pmSuffix
	return
}
